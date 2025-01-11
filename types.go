package main

func (a *App) TypeExporter(calibrationPoint CalibrationPoint, calibratedCalibrationPoint CalibratedCalibrationPoint, fixture Fixture, sacnConfig SACNConfig, dmxData DMXData) {
	// Explicitly export all types to the frontend, this should be done automatically by wails but when a type is "wrapped" in a map it doesn't seem to work
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
