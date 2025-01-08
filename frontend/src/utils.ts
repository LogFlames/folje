import type { Fixture, Point } from "./types";

export function convexHull(p_list: Point[]): Point[] {
    if (p_list.length < 3) return p_list;

    let hull = [];
    let tmp;

    // Find leftmost point
    tmp = p_list[0];
    for (const p of p_list) if (p.x < tmp.x) tmp = p;

    hull[0] = tmp;

    let endpoint: Point, secondlast: Point;
    let min_angle: number, new_end: Point;

    endpoint = hull[0];
    secondlast = { x: endpoint.x, y: endpoint.y + 10 };

    do {
        min_angle = Math.PI; // Initial value. Any angle must be lower that 2PI
        for (const p of p_list) {
            tmp = polarAngle(secondlast, endpoint, p);

            if (tmp <= min_angle) {
                new_end = p;
                min_angle = tmp;
            }
        }

        if (new_end != hull[0]) {
            hull.push(new_end);
            secondlast = endpoint;
            endpoint = new_end;
        }
    } while (new_end != hull[0]);

    return hull;
}

function polarAngle(a: Point, b: Point, c: Point): number {
    let x = (a.x - b.x) * (c.x - b.x) + (a.y - b.y) * (c.y - b.y);
    let y = (a.x - b.x) * (c.y - b.y) - (c.x - b.x) * (a.y - b.y);
    return Math.atan2(y, x);
}

export function calcTilt(fixture: Fixture, y: number): number {
    return fixture.minTilt + y * (fixture.maxTilt - fixture.minTilt);
}

export function calcPan(fixture: Fixture, x: number): number {
    return fixture.minPan + x * (fixture.maxPan - fixture.minPan);
}