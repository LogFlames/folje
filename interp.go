package main

import (
	"errors"
	"math"

	"github.com/fogleman/delaunay"
)

type Linear2DPanTiltInterpolator struct {
	points     []delaunay.Point
	tri        *delaunay.Triangulation
	panValues  []float64
	tiltValues []float64
	fillValue  float64
}

func area(a, b, c delaunay.Point) float64 {
	return math.Abs((a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y)) / 2.0)
}

func isPointInTriangle(a, b, c, p delaunay.Point) bool {
	A := area(a, b, c)

	A1 := area(p, b, c)
	A2 := area(a, p, c)
	A3 := area(a, b, p)

	return math.Abs(A-(A1+A2+A3)) < 1e-9
}

func (interp *Linear2DPanTiltInterpolator) LocatePoint(p delaunay.Point) (int, error) {
	// TODO: Increase efficiency by moving by the edge towards the point instead of brute forcing all of the points
	if len(interp.tri.Triangles) == 0 {
		return -1, errors.New("no triangles in the triangulation")
	}

	for t := 0; t < len(interp.tri.Triangles); t += 3 {
		// Bounds check for triangle indices
		if t+2 >= len(interp.tri.Triangles) {
			LogError("Triangle index out of bounds: %d >= %d", t+2, len(interp.tri.Triangles))
			return -1, errors.New("triangle index out of bounds")
		}
		idx0, idx1, idx2 := interp.tri.Triangles[t], interp.tri.Triangles[t+1], interp.tri.Triangles[t+2]
		if idx0 >= len(interp.points) || idx1 >= len(interp.points) || idx2 >= len(interp.points) {
			LogError("Point index out of bounds in triangle: indices [%d, %d, %d], points len %d", idx0, idx1, idx2, len(interp.points))
			return -1, errors.New("point index out of bounds")
		}
		if isPointInTriangle(interp.points[idx0], interp.points[idx1], interp.points[idx2], p) {
			return t / 3, nil
		}
	}

	return -1, errors.New("point is outside the convex hull")
}

func (interp *Linear2DPanTiltInterpolator) Barycentric(point delaunay.Point, triangle int) (float64, float64, float64, error) {
	// Bounds check for triangle array access
	triIdx := 3 * triangle
	if triIdx+2 >= len(interp.tri.Triangles) {
		LogError("Barycentric triangle index out of bounds: %d >= %d", triIdx+2, len(interp.tri.Triangles))
		return 0, 0, 0, errors.New("triangle index out of bounds")
	}
	idx0, idx1, idx2 := interp.tri.Triangles[triIdx], interp.tri.Triangles[triIdx+1], interp.tri.Triangles[triIdx+2]
	if idx0 >= len(interp.points) || idx1 >= len(interp.points) || idx2 >= len(interp.points) {
		LogError("Barycentric point index out of bounds: indices [%d, %d, %d], points len %d", idx0, idx1, idx2, len(interp.points))
		return 0, 0, 0, errors.New("point index out of bounds")
	}
	p1 := interp.points[idx0]
	p2 := interp.points[idx1]
	p3 := interp.points[idx2]

	// Denominator is twice the area of the triangle
	denominator := (p2.X-p3.X)*(p1.Y-p3.Y) - (p2.Y-p3.Y)*(p1.X-p3.X)

	if denominator == 0 {
		return 0.0, 0.0, 0.0, errors.New("triangle vertices are collinear; denominator is zero")
	}

	l1 := ((p2.X-p3.X)*(point.Y-p3.Y) - (p2.Y-p3.Y)*(point.X-p3.X)) / denominator
	l2 := ((p3.X-p1.X)*(point.Y-p3.Y) - (p3.Y-p1.Y)*(point.X-p3.X)) / denominator
	l3 := 1.0 - l1 - l2

	return l1, l2, l3, nil
}

func NewLinear2DPanTiltInterpolator(points []Point, panValues []float64, tiltValues []float64, fillValue float64) (*Linear2DPanTiltInterpolator, error) {
	if len(points) == 0 || len(panValues) == 0 || len(tiltValues) == 0 || len(points) != len(panValues) || len(points) != len(tiltValues) {
		return nil, errors.New("points and values must have the same non-zero length")
	}

	delaunayPoints := make([]delaunay.Point, len(points))
	for i, point := range points {
		delaunayPoints[i] = delaunay.Point{X: point.X, Y: point.Y}
	}

	tri, err := delaunay.Triangulate(delaunayPoints)
	if err != nil {
		return nil, err
	}

	return &Linear2DPanTiltInterpolator{
		points:     delaunayPoints,
		tri:        tri,
		panValues:  panValues,
		tiltValues: tiltValues,
		fillValue:  fillValue,
	}, nil
}

func (interp *Linear2DPanTiltInterpolator) Interpolate(point delaunay.Point) (float64, float64, error) {
	triangle, err := interp.LocatePoint(point)
	if err != nil {
		return interp.fillValue, interp.fillValue, nil // point is outside the convex hull
	}

	baryDist1, baryDist2, baryDist3, err := interp.Barycentric(point, triangle)
	if err != nil {
		return interp.fillValue, interp.fillValue, errors.New("triangle vertices are collinear; denominator is zero")
	}

	// Bounds check for value array access
	triIdx := triangle * 3
	if triIdx+2 >= len(interp.tri.Triangles) {
		LogError("Interpolate triangle index out of bounds: %d >= %d", triIdx+2, len(interp.tri.Triangles))
		return interp.fillValue, interp.fillValue, errors.New("triangle index out of bounds")
	}
	idx0, idx1, idx2 := interp.tri.Triangles[triIdx], interp.tri.Triangles[triIdx+1], interp.tri.Triangles[triIdx+2]
	if idx0 >= len(interp.panValues) || idx1 >= len(interp.panValues) || idx2 >= len(interp.panValues) {
		LogError("Interpolate pan value index out of bounds: indices [%d, %d, %d], panValues len %d", idx0, idx1, idx2, len(interp.panValues))
		return interp.fillValue, interp.fillValue, errors.New("pan value index out of bounds")
	}
	if idx0 >= len(interp.tiltValues) || idx1 >= len(interp.tiltValues) || idx2 >= len(interp.tiltValues) {
		LogError("Interpolate tilt value index out of bounds: indices [%d, %d, %d], tiltValues len %d", idx0, idx1, idx2, len(interp.tiltValues))
		return interp.fillValue, interp.fillValue, errors.New("tilt value index out of bounds")
	}

	pan := baryDist1*interp.panValues[idx0] + baryDist2*interp.panValues[idx1] + baryDist3*interp.panValues[idx2]
	tilt := baryDist1*interp.tiltValues[idx0] + baryDist2*interp.tiltValues[idx1] + baryDist3*interp.tiltValues[idx2]
	return pan, tilt, nil
}
