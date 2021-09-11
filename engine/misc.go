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

// ZPosition implements DrawAfter and DrawBefore using only Z information.
type ZPosition int

// DrawAfter reports if z >= x.Max.Z.
func (z ZPosition) DrawAfter(x Drawer) bool {
	switch x := x.(type) {
	case BoundingBoxer:
		return int(z) >= x.BoundingBox().Max.Z
	case ZPositioner:
		return z.ZPos() > x.ZPos()
	}
	return false
}

// DrawBefore reports if z < x.Min.Z.
func (z ZPosition) DrawBefore(x Drawer) bool {
	switch x := x.(type) {
	case BoundingBoxer:
		return int(z) < x.BoundingBox().Min.Z
	case ZPositioner:
		return z.ZPos() < x.ZPos()
	}
	return false
}

// ZPos returns itself.
func (z ZPosition) ZPos() int { return int(z) }
