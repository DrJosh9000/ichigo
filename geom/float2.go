package geom

import "fmt"

// Float2 is an element of float64^2
type Float2 struct {
	X, Y float64
}

// String returns a string representation of p like "(3.0,4.0,5.0)".
func (p Float2) String() string {
	return fmt.Sprintf("(%f,%f)", p.X, p.Y)
}

// Add performs vector addition.
func (p Float2) Add(q Float2) Float2 {
	return Float2{p.X + q.X, p.Y + q.Y}
}

// Sub performs vector subtraction.
func (p Float2) Sub(q Float2) Float2 {
	return p.Add(q.Neg())
}

// CMul performs componentwise multiplication.
func (p Float2) CMul(q Float2) Float2 {
	return Float2{p.X * q.X, p.Y * q.Y}
}

// Mul performs scalar multiplication.
func (p Float2) Mul(k float64) Float2 {
	return Float2{p.X * k, p.Y * k}
}

// CDiv performs componentwise division.
func (p Float2) CDiv(q Float2) Float2 {
	return Float2{p.X / q.X, p.Y / q.Y}
}

// Div performs scalar division by k.
func (p Float2) Div(k float64) Float2 {
	return Float2{p.X / k, p.Y / k}
}

// Neg returns the vector pointing in the opposite direction.
func (p Float2) Neg() Float2 {
	return Float2{-p.X, -p.Y}
}

// Coord returns the components of the vector.
func (p Float2) Coord() (x, y float64) {
	return p.X, p.Y
}

// Sign returns a sign vector.
func (p Float2) Sign() Float2 {
	return Float2{FSign(p.X), FSign(p.Y)}
}

// Dot returns the dot product of the two vectors.
func (p Float2) Dot(q Float2) float64 {
	return p.X*q.X + p.Y*q.Y
}
