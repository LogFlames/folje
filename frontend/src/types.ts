export interface Fixture {
    id: string;
    name: string;
    universe: number;
    panAddress: number;
    finePanAddress: number;
    tiltAddress: number;
    fineTiltAddress: number;
    minPan: number;
    maxPan: number;
    minTilt: number;
    maxTilt: number;
    calibration: { [id: string]: CalibratedCalibrationPoint }
}

export interface CalibrationPoint {
    id: string;
    x: number;
    y: number;
}

export interface CalibratedCalibrationPoint {
    id: string;
    pan: number;
    tilt: number;
}