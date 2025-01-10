package main

import (
	"context"

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

	opts := sacn.SenderOptions{
		SourceName: "Följe",
	}
	sender, err := sacn.NewSender("192.168.68.77", &opts)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	a.sender = sender
}

func (a *App) shutdown(ctx context.Context) {
	for uni := range a.activeUniverses {
		a.DeactiveUniverse(uni)
	}
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

func (a *App) TypeExporter(point Point, calibrationPoint CalibrationPoint, calibratedCalibrationPoint CalibratedCalibrationPoint, fixture Fixture) {
}

func (a *App) SetCalibrationPoints(calibrationPoints map[string]CalibrationPoint) {
	a.calibrationPoints = calibrationPoints
}

func (a *App) SetFixtures(fixtures map[string]Fixture) {
	a.fixtures = fixtures
	a.universeDMXData = make(map[uint16]DMXData)

	for uni := range a.universeDMXData {
		inUse := false
		for _, fixture := range fixtures {
			if uni == fixture.Universe {
				inUse = true
				break
			}
		}

		if inUse {
			continue
		}

		a.DeactiveUniverse(uni)
	}

	for _, fixture := range fixtures {
		a.ActiveUniverse(fixture.Universe)
	}
}

func (a *App) SendForAllFromMouse(x float32, y float32) {
	a.mouseX = x
	a.mouseY = y
}

func (a *App) SendPanTilt(fixtureId string, pan int, tilt int) {
	fixture := a.fixtures[fixtureId]
	packet := packet.NewDataPacket()

	data := a.universeDMXData[fixture.Universe]

	data[fixture.PanAddress] = byte(pan / 256)
	data[fixture.FinePanAddress] = byte(pan % 256)
	data[fixture.TiltAddress] = byte(tilt / 256)
	data[fixture.FineTiltAddress] = byte(tilt % 256)

	a.universeDMXData[fixture.Universe] = data

	packet.SetData(data[:])

	a.activeUniverses[fixture.Universe] <- packet
}

func (a *App) ActiveUniverse(uni uint16) {
	if a.activeUniverses[uni] != nil {
		return
	}

	universe, err := a.sender.StartUniverse(uni)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return
	}
	a.sender.SetMulticast(uni, true)
	a.sender.AddDestination(uni, "192.168.68.76")
	a.activeUniverses[uni] = universe
}

func (a *App) DeactiveUniverse(uni uint16) {
	if a.activeUniverses[uni] == nil {
		return
	}

	a.sender.StopUniverse(uni)
	delete(a.activeUniverses, uni)
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
