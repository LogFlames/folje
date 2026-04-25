package main

import (
	"errors"

	"github.com/fogleman/delaunay"
)

type Linear2DPanTiltInterpolator struct {
	points     []delaunay.Point
	tri        *delaunay.Triangulation
	panValues  []float64
	tiltValues []float64
	fillValue  float64
}

// cross returns the cross product of vectors (b-a) and (p-a).
func cross(a, b, p delaunay.Point) float64 {
	return (b.X-a.X)*(p.Y-a.Y) - (b.Y-a.Y)*(p.X-a.X)
}

// pointInTriangle returns true if p lies inside or on the boundary of triangle abc.
// Orientation-agnostic: accepts the triangle whether vertices are CW or CCW.
func pointInTriangle(a, b, c, p delaunay.Point) bool {
	d1 := cross(a, b, p)
	d2 := cross(b, c, p)
	d3 := cross(c, a, p)

	hasNeg := d1 < 0 || d2 < 0 || d3 < 0
	hasPos := d1 > 0 || d2 > 0 || d3 > 0

	return !(hasNeg && hasPos)
}

func (interp *Linear2DPanTiltInterpolator) LocatePoint(p delaunay.Point) (int, error) {
	numEdges := len(interp.tri.Triangles)
	if numEdges == 0 {
		return -1, errors.New("no triangles in the triangulation")
	}

	for t := 0; t < numEdges; t += 3 {
		a := interp.points[interp.tri.Triangles[t]]
		b := interp.points[interp.tri.Triangles[t+1]]
		c := interp.points[interp.tri.Triangles[t+2]]
		if pointInTriangle(a, b, c, p) {
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
