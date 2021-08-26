package engine

import "image"

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

// Parallax implements ParallaxScaler directly (as a float64 value).
type Parallax float64

// ParallaxFactor returns s as a float64.
func (s Parallax) ParallaxFactor() float64 { return float64(s) }

// ZOrder implements DrawOrderer directly (as a float64 value).
type ZOrder float64

// DrawOrder returns z as a float64.
func (z ZOrder) DrawOrder() float64 { return float64(z) }

// Some math helpers

func mul2(p, q image.Point) image.Point {
	p.X *= q.X
	p.Y *= q.Y
	return p
}

func div2(p, q image.Point) image.Point {
	p.X /= q.X
	p.Y /= q.Y
	return p
}

func float2(p image.Point) (float64, float64) {
	return float64(p.X), float64(p.Y)
}
