package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gitlab.com/patopest/go-sacn"
	"gitlab.com/patopest/go-sacn/packet"
)

// App struct
type App struct {
	ctx context.Context

	fixtures          map[string]Fixture
	universeDMXData   map[uint16]DMXData
	calibrationPoints map[string]CalibrationPoint
	mouseX            float32
	mouseY            float32
	sender            *sacn.Sender
	activeUniverses   map[uint16]chan<- packet.SACNPacket

	sacnStopLoop      chan bool
	sacnUpdatedConfig chan bool
	sacnConfig        *SACNConfig
	sacnWorkerWG      sync.WaitGroup
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.activeUniverses = make(map[uint16]chan<- packet.SACNPacket)
	a.fixtures = make(map[string]Fixture)
	a.calibrationPoints = make(map[string]CalibrationPoint)
	a.universeDMXData = make(map[uint16]DMXData)
	a.sacnConfig = &SACNConfig{
		Fps:          25,
		Multicast:    true,
		Destinations: []string{},
	}
	a.sacnWorkerWG = sync.WaitGroup{}
	a.sacnStopLoop = make(chan bool)
	a.sacnUpdatedConfig = make(chan bool)

	a.findPossibleIPAddresses()

	a.sacnWorkerWG.Add(1)
	go a.sacnWorker()
}

func (a *App) shutdown(ctx context.Context) {
	a.sacnStopLoop <- true
	a.sacnWorkerWG.Wait()
}

func (a *App) findPossibleIPAddresses() {
	interfaces, err := net.Interfaces()
	if err != nil {
		a.AlertDialog("Error finding IP addresses", err.Error())
		return
	}

	possibleAddresses := make([]string, 0)

	// Iterate through the interfaces to find the first non-loopback address
	for _, iface := range interfaces {
		// Skip loopback interfaces (127.0.0.1, etc.)
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Get the addresses for the interface
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}

		// Iterate over the addresses and look for an IPv4 address
		for _, addr := range addrs {
			// Check if the address is an IPv4 address
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
				possibleAddresses = append(possibleAddresses, ipNet.IP.String())
			}
		}
	}

	a.sacnConfig.PossibleIpAddresses = possibleAddresses
	if len(possibleAddresses) == 0 {
		a.AlertDialog("No IP addresses found", "This should not happen. Make sure you have a network interface active.")
		return
	}
	a.sacnConfig.IpAddress = possibleAddresses[0]

}

func (a *App) ensureSACNSender() error {
	if a.sender != nil {
		return nil
	}

	opts := sacn.SenderOptions{
		SourceName: "Följe",
	}
	sender, err := sacn.NewSender(a.sacnConfig.IpAddress, &opts)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		a.sender = nil
		a.AlertDialog("sACN sender error", err.Error())
		return err
	}

	a.sender = sender

	return nil
}

func (a *App) AlertDialog(title string, message string) {
	options := runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	}

	_, err := runtime.MessageDialog(a.ctx, options)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return
	}
}

func (a *App) ConfirmDialog(title string, message string) string {
	options := runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         title,
		Message:       message,
		Buttons:       []string{"Ok", "Cancel"},
		DefaultButton: "Ok",
		CancelButton:  "Cancel",
	}
	res, err := runtime.MessageDialog(a.ctx, options)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return ""
	}

	return res
}

func (a *App) TypeExporter(point Point, calibrationPoint CalibrationPoint, calibratedCalibrationPoint CalibratedCalibrationPoint, fixture Fixture, sacnConfig SACNConfig) {
}

func (a *App) SetCalibrationPoints(calibrationPoints map[string]CalibrationPoint) {
	a.calibrationPoints = calibrationPoints
}

func (a *App) SetFixtures(fixtures map[string]Fixture) {
	a.fixtures = fixtures
	a.universeDMXData = make(map[uint16]DMXData)

}

func (a *App) ensureSACNUniverses() error {
	err := a.ensureSACNSender()
	if err != nil {
		return err
	}

	for uni := range a.universeDMXData {
		inUse := false
		for _, fixture := range a.fixtures {
			if uni == fixture.Universe {
				inUse = true
				break
			}
		}

		if inUse {
			continue
		}

		a.deactiveUniverse(uni)
	}

	for _, fixture := range a.fixtures {
		a.activateUniverse(fixture.Universe)
	}

	return nil
}

func (a *App) SetMouseForAllFixtures(x float32, y float32) {
	// TODO
	a.mouseX = x
	a.mouseY = y
}

func (a *App) SetPanTiltForFixture(fixtureId string, pan int, tilt int) {
	fixture, exists := a.fixtures[fixtureId]
	if !exists {
		runtime.LogError(a.ctx, "Tried to set pan/tilt for non-existing fixture")
		return
	}

	data := a.universeDMXData[fixture.Universe]

	if fixture.PanAddress >= 0 && fixture.PanAddress < 512 {
		data[fixture.PanAddress] = byte(pan / 256)
	}
	if fixture.FinePanAddress >= 0 && fixture.FinePanAddress < 512 {
		data[fixture.FinePanAddress] = byte(pan % 256)
	}
	if fixture.TiltAddress >= 0 && fixture.TiltAddress < 512 {
		data[fixture.TiltAddress] = byte(tilt / 256)
	}
	if fixture.FineTiltAddress >= 0 && fixture.FineTiltAddress < 512 {
		data[fixture.FineTiltAddress] = byte(tilt % 256)
	}

	a.universeDMXData[fixture.Universe] = data
}

func (a *App) sacnWorker() {
	ticker := time.NewTicker(time.Second)
	if a.sacnConfig.Fps != 0 {
		ticker = time.NewTicker(time.Second / time.Duration(a.sacnConfig.Fps))
	}

	work := func() {
		a.ensureSACNSender()
		a.ensureSACNUniverses()

		for uni := range a.activeUniverses {
			p := packet.NewDataPacket()
			data := a.universeDMXData[uni]
			p.SetData(data[:])
			a.activeUniverses[uni] <- p
		}
	}

	stop := func() {
		for uni := range a.activeUniverses {
			a.deactiveUniverse(uni)
		}

		if a.sender != nil {
			a.closeSACNSender()
		}

		a.sacnWorkerWG.Done()
	}

	for {
		select {
		case <-a.sacnUpdatedConfig:
			ticker = time.NewTicker(time.Second)
			if a.sacnConfig.Fps != 0 {
				ticker = time.NewTicker(time.Second / time.Duration(a.sacnConfig.Fps))
			}
		case <-a.sacnStopLoop:
			stop()
			return
		case <-ticker.C:
			work()
		}
	}
}

func (a *App) activateUniverse(uni uint16) {
	if a.activeUniverses[uni] != nil {
		return
	}

	universe, err := a.sender.StartUniverse(uni)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return
	}
	a.sender.SetMulticast(uni, a.sacnConfig.Multicast)
	for _, dest := range a.sacnConfig.Destinations {
		a.sender.AddDestination(uni, dest)
	}
	a.activeUniverses[uni] = universe
}

func (a *App) deactiveUniverse(uni uint16) {
	if a.activeUniverses[uni] == nil {
		return
	}

	a.sender.StopUniverse(uni)
	delete(a.activeUniverses, uni)
}

func (a *App) SetSACNConfig(sacnConfig SACNConfig) {
	a.sacnConfig = &sacnConfig

	if a.sender != nil {
		a.closeSACNSender()
	}

	a.ensureSACNSender()
	a.ensureSACNUniverses()

	a.sacnUpdatedConfig <- true
}

func (a *App) closeSACNSender() {
	if a.sender == nil {
		return
	}

	a.sender.Close()
	a.sender = nil

	a.activeUniverses = make(map[uint16]chan<- packet.SACNPacket)
}

func (a *App) GetSACNConfig() SACNConfig {
	return *a.sacnConfig
}

type Point struct {
	X float32
	Y float32
}

type CalibrationPoint struct {
	Id   string
	Name string
	X    float32
	Y    float32
}

type CalibratedCalibrationPoint struct {
	Id   string
	Pan  int
	Tilt int
}

type Fixture struct {
	Id              string
	Name            string
	Universe        uint16
	PanAddress      int
	FinePanAddress  int
	TiltAddress     int
	FineTiltAddress int
	Calibration     map[string]CalibratedCalibrationPoint
}

type DMXData [512]byte

type SACNConfig struct {
	IpAddress           string
	PossibleIpAddresses []string
	Fps                 int
	Multicast           bool
	Destinations        []string
}
