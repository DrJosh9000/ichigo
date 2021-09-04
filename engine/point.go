package engine

import (
	"image"
	"strconv"
)

// Point3 is a an element of int^3.
type Point3 struct {
	X, Y, Z int
}

// Pt3(x, y, z) is shorthand for Point3{x, y, z}.
func Pt3(x, y, z int) Point3 {
	return Point3{x, y, z}
}

// String returns a string representation of p like "(3,4,5)".
func (p Point3) String() string {
	return "(" + strconv.Itoa(p.X) + "," + strconv.Itoa(p.Y) + "," + strconv.Itoa(p.Z) + ")"
}

// XY applies the Z-forgetting projection. (It returns just X and Y.)
func (p Point3) XY() image.Point {
	return image.Point{p.X, p.Y}
}

// Add performs vector addition.
func (p Point3) Add(q Point3) Point3 {
	return Point3{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

// Sub performs vector subtraction.
func (p Point3) Sub(q Point3) Point3 {
	return p.Add(q.Neg())
}

// CMul performs componentwise multiplication.
func (p Point3) CMul(q Point3) Point3 {
	return Point3{p.X * q.X, p.Y * q.Y, p.Z * q.Z}
}

// Mul performs scalar multiplication.
func (p Point3) Mul(k int) Point3 {
	return Point3{p.X * k, p.Y * k, p.Z * k}
}

// CDiv performs componentwise division.
func (p Point3) CDiv(q Point3) Point3 {
	return Point3{p.X / q.X, p.Y / q.Y, p.Z / q.Z}
}

// Div performs scalar division by k.
func (p Point3) Div(k int) Point3 {
	return Point3{p.X / k, p.Y / k, p.Z / k}
}

// Neg returns the vector pointing in the opposite direction.
func (p Point3) Neg() Point3 {
	return Point3{-p.X, -p.Y, -p.Z}
}

// Coord returns the components of the vector.
func (p Point3) Coord() (x, y, z int) {
	return p.X, p.Y, p.Z
}

// Sign returns a sign vector.
func (p Point3) Sign() Point3 {
	return Point3{sign(p.X), sign(p.Y), sign(p.Z)}
}

func sign(m int) int {
	if m == 0 {
		return 0
	}
	if m < 0 {
		return -1
	}
	return 1
}

// IsoProject performs isometric projection of a 3D coordinate into 2D.
//
// If π.X = 0, the x returned is p.X; similarly for π.Y and y.
// Otherwise, x projects to x + z/π.X and y projects to y + z/π.Y.
func (p Point3) IsoProject(π image.Point) image.Point {
	/*
		I'm using the π character because I'm a maths wanker.

		Dividing is used because there's little reason for an isometric
		projection in a game to exaggerate the Z position.

		Integers are used to preserve that "pixel perfect" calculation in case
		you are making the next Celeste.
	*/
	q := image.Point{p.X, p.Y}
	if π.X != 0 {
		q.X += p.Z / π.X
	}
	if π.Y != 0 {
		q.Y += p.Z / π.Y
	}
	return q
}
