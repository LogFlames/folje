package main

import (
	"context"
	"sync"

	"github.com/fogleman/delaunay"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gitlab.com/patopest/go-sacn"
	"gitlab.com/patopest/go-sacn/packet"
)

type App struct {
	ctx context.Context

	fixtures          map[string]Fixture
	universeDMXData   map[uint16]DMXData
	calibrationPoints map[string]CalibrationPoint
	sender            *sacn.Sender
	activeUniverses   map[uint16]chan<- packet.SACNPacket

	sacnStopLoop        chan bool
	sacnUpdatedConfig   chan bool
	sacnConfig          *SACNConfig
	sacnWorkerWG        sync.WaitGroup
	linearInterpolators map[string]*Linear2DPanTiltInterpolator
}

func NewApp() *App {
	return &App{}
}

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

	a.linearInterpolators = make(map[string]*Linear2DPanTiltInterpolator)

	a.findPossibleIPAddresses()

	a.sacnWorkerWG.Add(1)
	go a.sacnWorker()
}

func (a *App) shutdown(ctx context.Context) {
	a.sacnStopLoop <- true
	a.sacnWorkerWG.Wait()
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

func (a *App) SetCalibrationPoints(calibrationPoints map[string]CalibrationPoint) {
	a.calibrationPoints = calibrationPoints
	a.calculateLinearInterpolator()
}

func (a *App) SetFixtures(fixtures map[string]Fixture) {
	a.fixtures = fixtures
	a.universeDMXData = make(map[uint16]DMXData)
	a.calculateLinearInterpolator()
}

func (a *App) calculateLinearInterpolator() {
	a.linearInterpolators = make(map[string]*Linear2DPanTiltInterpolator)

	points := make([]Point, len(a.calibrationPoints))
	pointsIndexMap := make(map[string]int)
	index := 0
	for _, calibrationPoint := range a.calibrationPoints {
		points[index] = Point{X: calibrationPoint.X, Y: calibrationPoint.Y}
		pointsIndexMap[calibrationPoint.Id] = index
		index++
	}

	for _, fixture := range a.fixtures {
		calibrated := true
		for _, calibrationPoint := range a.calibrationPoints {
			_, exists := fixture.Calibration[calibrationPoint.Id]
			if !exists {
				calibrated = false
				break
			}
		}

		if !calibrated {
			continue
		}

		panValues := make([]float64, len(a.calibrationPoints))
		tiltValues := make([]float64, len(a.calibrationPoints))
		for _, calibrationPoint := range a.calibrationPoints {
			panValues[pointsIndexMap[calibrationPoint.Id]] = float64(fixture.Calibration[calibrationPoint.Id].Pan)
			tiltValues[pointsIndexMap[calibrationPoint.Id]] = float64(fixture.Calibration[calibrationPoint.Id].Tilt)
		}

		interp, err := NewLinear2DPanTiltInterpolator(points, panValues, tiltValues, -1.0)
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
			continue
		}

		a.linearInterpolators[fixture.Id] = interp
	}
}

func (a *App) SetMouseForAllFixtures(x float64, y float64) {
	for _, fixture := range a.fixtures {
		interp, exists := a.linearInterpolators[fixture.Id]
		if !exists {
			continue
		}

		pan, tilt, err := interp.Interpolate(delaunay.Point{X: x, Y: y})
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
			continue
		}

		if pan == interp.fillValue || tilt == interp.fillValue {
			continue
		}

		a.SetPanTiltForFixture(fixture.Id, int(pan), int(tilt))
	}
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
