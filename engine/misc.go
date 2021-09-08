package engine

import (
	"image"
)

// ID implements Identifier directly (as a string value).
type ID string

// Ident returns id as a string.
func (id ID) Ident() string { return string(id) }

// Bounds implements Bounder directly (as an image.Rectangle value).
type Bounds image.Rectangle

// BoundingRect returns b as an image.Rectangle.
func (b Bounds) BoundingRect() image.Rectangle { return image.Rectangle(b) }

// Disabled implements Disabler directly (as a bool).
type Disabled bool

// IsHidden returns h as a bool.
func (d Disabled) IsDisabled() bool { return bool(d) }

// Hide sets h to true.
func (d *Disabled) Disable() { *d = true }

// Show sets h to false.
func (d *Disabled) Enable() { *d = false }

// Hidden implements Hider directly (as a bool).
type Hidden bool

// IsHidden returns h as a bool.
func (h Hidden) IsHidden() bool { return bool(h) }

// Hide sets h to true.
func (h *Hidden) Hide() { *h = true }

// Show sets h to false.
func (h *Hidden) Show() { *h = false }

// ZOrder implements DrawOrder (in Drawer) directly (as an int value).
type ZOrder int

// DrawOrder returns z as a int with 0 bias.
func (z ZOrder) DrawOrder() (int, int) { return int(z), 0 }

// ---------- Some helpers for image.Point ----------

// cmul performs componentwise multiplication of two image.Points.
func cmul(p, q image.Point) image.Point {
	return image.Point{p.X * q.X, p.Y * q.Y}
}

// cdiv performs componentwise division of two image.Points.
func cdiv(p, q image.Point) image.Point {
	return image.Point{p.X / q.X, p.Y / q.Y}
}

// cfloat returns the components of an image.Point as two floats.
func cfloat(p image.Point) (x, y float64) {
	return float64(p.X), float64(p.Y)
}

// dot returns the dot product of two image.Points.
func dot(p, q image.Point) int {
	return p.X*q.X + p.Y*q.Y
}

// polygonContains reports if a polygon contains a point
func polygonContains(polygon []image.Point, p image.Point) bool {
	for i, p1 := range polygon {
		p2 := polygon[(i+1)%len(polygon)]
		// ∆(p p1 p2) should have positive signed area
		p1, p2 = p1.Sub(p), p2.Sub(p)
		if p2.X*p1.Y-p1.X*p2.Y < 0 {
			return false
		}
	}
	return true
}

// polygonRectOverlap reports if a polygon overlaps a rectangle.
func polygonRectOverlap(polygon []image.Point, rect image.Rectangle) bool {
	// There's a good chance a vertex from one thing is inside the other...

	// Check if any vertex of the polygon is inside the rect.
	for _, p := range polygon {
		if p.In(rect) {
			return true
		}
	}

	// Reduced Max (to the inclusive bound).
	rmax := rect.Max.Sub(image.Pt(1, 1))

	// Check if any vertex of the rect is inside the polygon.
	if polygonContains(polygon, rect.Min) {
		return true
	}
	if polygonContains(polygon, rmax) {
		return true
	}
	if polygonContains(polygon, image.Pt(rect.Min.X, rmax.Y)) {
		return true
	}
	if polygonContains(polygon, image.Pt(rmax.X, rect.Min.Y)) {
		return true
	}

	// Only remaining cases involve line intersection between the rect and poly.

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
