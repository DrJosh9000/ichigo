package geom

import (
	"image"
	"math"
)

const (
	East = iota
	North
	West
	South
)

// PolygonExtrema returns the most easterly, northerly, westerly, and southerly
// points (north is in the -Y direction, east is in the +X direction, etc). If
// there are multiple points furthest in any direction, the first one is used.
func PolygonExtrema(polygon []image.Point) [4]image.Point {
	var e, n, w, s image.Point
	e.X = math.MinInt
	n.Y = math.MaxInt
	w.X = math.MaxInt
	s.Y = math.MinInt
	for _, p := range polygon {
		if p.X > e.X {
			e = p
		}
		if p.X < w.X {
			w = p
		}
		if p.Y > s.Y {
			s = p
		}
		if p.Y < n.Y {
			n = p
		}
	}
	return [4]image.Point{East: e, North: n, West: w, South: s}
}

// PolygonContains reports if a convex polygon contains a point. The polygon
// must be in clockwise order if +Y is pointing upwards, or anticlockwise if +Y
// is pointing downwards.
func PolygonContains(convex []image.Point, p image.Point) bool {
	for i, q := range convex {
		r := convex[(i+1)%len(convex)]
		// ∆(p q r) should have positive signed area
		q, r = q.Sub(p), r.Sub(p)
		if q.X*r.Y > r.X*q.Y {
			return false
		}
	}
	return true
}

// PolygonRectOverlap reports if a convex polygon overlaps a rectangle.
func PolygonRectOverlap(convex []image.Point, rect image.Rectangle) bool {
	if convex[0].In(rect) {
		return true
	}

	// Check if any vertex of the rect is inside the polygon.
	if PolygonContains(convex, rect.Min) {
		return true
	}
	// Reduced Max (to the inclusive bound).
	rmax := rect.Max.Sub(image.Pt(1, 1))
	// Since we went to the trouble of computing another point...
	// TODO: this shouldn't be necessary
	if PolygonContains(convex, rmax) {
		return true
	}

	// Only remaining cases involve line intersection between the rect and poly
	// having eliminated the possibility that one is entirely within another.

	// Since rect is an axis-aligned rectangle, we only need vertical and
	// horizontal line intersection tests.
	// Walk each edge of polygon.
	for i, p := range convex {
		q := convex[(i+1)%len(convex)]
		// Pretend the edge is a rectangle. Exclude those that don't overlap.
		if !rect.Overlaps(image.Rectangle{p, q}.Canon()) {
			continue
		}

		d := q.Sub(p)
		// If the polygon edge is not vertical, test left and right sides
		if d.X != 0 {
			if d.X < 0 {
				d = d.Mul(-1)
			}
			min := (rect.Min.Y - p.Y) * d.X
			max := (rect.Max.Y - p.Y) * d.X
			// Test left side of rect
			if t := (rect.Min.X - p.X) * d.Y; min <= t && t < max {
				return true
			}
			// Test right side of rect
			if t := (rmax.X - p.X) * d.Y; min <= t && t < max {
				return true
			}
		}
		// If the polygon edge is not horizontal, test the top and bottom sides
		if d.Y != 0 {
			if d.Y < 0 {
				d = d.Mul(-1)
			}
			min := (rect.Min.X - p.X) * d.Y
			max := (rect.Max.X - p.X) * d.Y
			// Test top side of rect
			if t := (rect.Min.Y - p.Y) * d.X; min <= t && t < max {
				return true
			}
			// Test bottom side of rect
			if t := (rmax.Y - p.Y) * d.X; min <= t && t < max {
				return true
			}
		}
	}
	return false
}
