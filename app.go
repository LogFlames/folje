package main

import (
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type Folje struct {
	fixtures          []Fixture
	calibrationPoints map[string]CalibrationPoint
}

func (f *Folje) TypeExporter(point Point, calibrationPoint CalibrationPoint, calibratedCalibrationPoint CalibratedCalibrationPoint, fixture Fixture) {

}

func (f *Folje) SetCalibrationPoints(calibrationPoints map[string]CalibrationPoint) {
	f.calibrationPoints = calibrationPoints
}

func (f *Folje) SetFixtures(fixtures []Fixture) {
	f.fixtures = fixtures
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
	Universe        int
	PanAddress      int
	FinePanAddress  int
	TiltAddress     int
	FineTiltAddress int
	Calibration     map[string]CalibratedCalibrationPoint
}
