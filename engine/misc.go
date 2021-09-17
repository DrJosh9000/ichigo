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

// Disables implements Disabler directly (as a bool).
type Disables bool

// Disabled returns d as a bool.
func (d Disables) Disabled() bool { return bool(d) }

// Disable sets d to true.
func (d *Disables) Disable() { *d = true }

// Enable sets d to false.
func (d *Disables) Enable() { *d = false }

// Hides implements Hider directly (as a bool).
type Hides bool

// Hidden returns h as a bool.
func (h Hides) Hidden() bool { return bool(h) }

// Hide sets h to true.
func (h *Hides) Hide() { *h = true }

// Show sets h to false.
func (h *Hides) Show() { *h = false }

// ZPosition implements ZPositioner directly (as an int).
type ZPosition int

// ZPos returns itself.
func (z ZPosition) ZPos() int { return int(z) }
