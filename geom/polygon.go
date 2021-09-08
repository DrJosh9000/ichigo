package geom

import "image"

// PolygonContains reports if a polygon contains a point
func PolygonContains(polygon []image.Point, p image.Point) bool {
	for i, q := range polygon {
		r := polygon[(i+1)%len(polygon)]
		// ∆(p q r) should have positive signed area
		q, r = q.Sub(p), r.Sub(p)
		if q.X*r.Y > r.X*q.Y {
			return false
		}
	}
	return true
}

// PolygonRectOverlap reports if a polygon overlaps a rectangle.
func PolygonRectOverlap(polygon []image.Point, rect image.Rectangle) bool {
	if polygon[0].In(rect) {
		return true
	}

	// Check if any vertex of the rect is inside the polygon.
	if PolygonContains(polygon, rect.Min) {
		return true
	}

	// Reduced Max (to the inclusive bound).
	rmax := rect.Max.Sub(image.Pt(1, 1))

	// Only remaining cases involve line intersection between the rect and poly
	// having eliminated the possibility that one is entirely within another.

	// Since rect is an axis-aligned rectangle, we only need vertical and
	// horizontal line intersection tests.
	// Walk each edge of polygon.
	for i, p := range polygon {
		q := polygon[(i+1)%len(polygon)]
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
