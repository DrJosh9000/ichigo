package engine

// ZPos implements ZPositioner directly (as a float64 value).
type ZPos float64

// Z returns z as a float64.
func (z ZPos) Z() float64 { return float64(z) }
