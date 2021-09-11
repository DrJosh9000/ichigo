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

// ZPosition implements DrawAfter and DrawPosition as a simple Z coordinate.
type ZPosition int

// DrawAfter reports if z >= x.Max.Z.
func (z ZPosition) DrawAfter(x Drawer) bool {
	switch d := x.(type) {
	case BoundingBoxer:
		return int(z) >= d.BoundingBox().Max.Z
	case zpositioner:
		return z.zposition() > d.zposition()
	}
	return false
}

// DrawBefore reports if z < x.Min.Z.
func (z ZPosition) DrawBefore(x Drawer) bool {
	switch d := x.(type) {
	case BoundingBoxer:
		return int(z) < d.BoundingBox().Min.Z
	case zpositioner:
		return z.zposition() < d.zposition()
	}
	return false
}

func (z ZPosition) zposition() int { return int(z) }

type zpositioner interface {
	zposition() int
}
