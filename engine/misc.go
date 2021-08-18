package engine

// ID implements Identifier directly (as a string value).
type ID string

// Ident returns id as a string.
func (id ID) Ident() string { return string(id) }

// Parallax implements ParallaxScaler directly (as a float64 value).
type Parallax float64

// ParallaxFactor returns s as a float64.
func (s Parallax) ParallaxFactor() float64 { return float64(s) }

// ZOrder implements DrawOrderer directly (as a float64 value).
type ZOrder float64

// DrawOrder returns z as a float64.
func (z ZOrder) DrawOrder() float64 { return float64(z) }
