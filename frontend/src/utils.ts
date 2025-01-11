import { main } from "../wailsjs/go/models";
import type { CalibrationPoint, CalibrationPoints, Fixture, Fixtures, MousePos, Point } from "./types";

export function convexHull(points: CalibrationPoint[]): Point[] {
    if (points.length < 3) return points;

    // Sort points lexicographically by x, then by y
    points.sort((a, b) => a.x - b.x || a.y - b.y);

    const crossProduct = (o, a, b) => {
        return (a.x - o.x) * (b.y - o.y) - (a.y - o.y) * (b.x - o.x);
    };

    // Build the lower hull
    const lower: Point[] = [];
    for (const p of points) {
        while (lower.length >= 2 && crossProduct(lower[lower.length - 2], lower[lower.length - 1], p) < 0) {
            lower.pop();
        }
        lower.push(p);
    }

    // Build the upper hull
    const upper: Point[] = [];
    for (const p of [...points].reverse()) {
        while (upper.length >= 2 && crossProduct(upper[upper.length - 2], upper[upper.length - 1], p) < 0) {
            upper.pop();
        }
        upper.push(p);
    }

    // Remove the last point of each half because it is repeated
    upper.pop();
    lower.pop();

    // Concatenate lower and upper hull to form the convex hull
    return lower.concat(upper);
}

export function calcTilt(fixture: Fixture, mousePos: MousePos, mouseDragStart: MousePos): number {
    let y = mousePos.y;
    if (mouseDragStart !== null) {
        y = mouseDragStart.y + (mousePos.y - mouseDragStart.y) * 0.05;
    }
    return fixture.minTilt + y * (fixture.maxTilt - fixture.minTilt);
}

export function calcPan(fixture: Fixture, mousePos: MousePos, mouseDragStart: MousePos): number {
    let x = mousePos.x;
    if (mouseDragStart !== null) {
        x = mouseDragStart.x + (mousePos.x - mouseDragStart.x) * 0.05;
    }
    return fixture.minPan + x * (fixture.maxPan - fixture.minPan);
}

export function convertFixturesToGo(fixtures: Fixtures, calibrationPoints: CalibrationPoints): { [id: string]: main.Fixture } {
    let goFixtures: { [id: string]: main.Fixture } = {};
    for (let fixture of Object.values(fixtures)) {
        let goCalibration: {
            [id: string]: main.CalibratedCalibrationPoint;
        } = {};

        for (let calibratedRalibrationPointId in fixture.calibration) {
            goCalibration[calibratedRalibrationPointId] = new main.CalibratedCalibrationPoint({
                Id: calibratedRalibrationPointId,
                Pan: Math.floor(fixture.calibration[calibratedRalibrationPointId].pan),
                Tilt: Math.floor(fixture.calibration[calibratedRalibrationPointId].tilt)
            });
        }

        goFixtures[fixture.id] = new main.Fixture({
            Id: fixture.id,
            Name: fixture.name,
            Universe: fixture.universe,
            PanAddress: fixture.panAddress - 1,
            FinePanAddress: fixture.finePanAddress - 1,
            TiltAddress: fixture.tiltAddress - 1,
            FineTiltAddress: fixture.fineTiltAddress - 1,
            Calibration: goCalibration
        });
    }

    return goFixtures;
}

export function convertCalibrationPointsToGo(calibrationPoints: CalibrationPoints): { [id: string]: main.CalibrationPoint } {
    let goCalibrationPoints: { [id: string]: main.CalibrationPoint } = {};

    for (let calibrationPoint of Object.values(calibrationPoints)) {
        goCalibrationPoints[calibrationPoint.id] = new main.CalibrationPoint({
            Id: calibrationPoint.id,
            Name: calibrationPoint.name,
            X: calibrationPoint.x,
            Y: calibrationPoint.y
        });
    }

    return goCalibrationPoints;
}
