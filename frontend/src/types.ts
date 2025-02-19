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
    name: string;
    x: number;
    y: number;
}

export interface CalibratedCalibrationPoint {
    id: string;
    pan: number;
    tilt: number;
}

export interface Point {
    x: number;
    y: number;
}

export interface MousePos {
    x: number;
    y: number;
}

export type Fixtures = { [id: string]: Fixture };
export type CalibrationPoints = { [id: string]: CalibrationPoint };

export interface SACNConfig {
    ipAddress: string;
    possibleIdAddresses: string[];
    fps: number;
    multicast: boolean;
    destinations: string[];
}

export interface CalibratingFixture {
        fixture_id: string;
        calibration_point_id: string;
    }