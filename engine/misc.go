package engine

// ID implements Identifier directly (as a string value).
type ID string

// Ident returns id as a string.
func (id ID) Ident() string { return string(id) }

// ZPos implements ZPositioner directly (as a float64 value).
type ZPos float64

// Z returns z as a float64.
func (z ZPos) Z() float64 { return float64(z) }
