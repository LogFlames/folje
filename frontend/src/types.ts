export interface Fixture {
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
}

export interface CalibrationPoint {
    uid: string;
    x: number;
    y: number;
}