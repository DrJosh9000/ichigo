package engine

// ID implements Identifier directly (as a string value).
type ID string

// Ident returns id as a string.
func (id ID) Ident() string { return string(id) }

// DrawOrder implements DrawOrderer directly (as a float64 value).
type DrawOrder float64

// DrawOrder returns z as a float64.
func (z DrawOrder) DrawOrder() float64 { return float64(z) }
