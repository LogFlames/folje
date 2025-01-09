import type { Fixture, Point, MousePos, CalibrationPoint } from "./types";

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
        y = mouseDragStart.y + (mousePos.y - mouseDragStart.y) * 0.002;
    }
    return fixture.minTilt + y * (fixture.maxTilt - fixture.minTilt);
}

export function calcPan(fixture: Fixture, mousePos: MousePos, mouseDragStart: MousePos): number {
    let x = mousePos.x;
    if (mouseDragStart !== null) {
        x = mouseDragStart.x + (mousePos.x - mouseDragStart.x) * 0.002;
    }
    return fixture.minPan + x * (fixture.maxPan - fixture.minPan);
}