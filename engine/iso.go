package engine

import "strconv"

// Point3 is a en element of int^3.
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

// Box describes an axis-aligned rectangular prism.
type Box struct {
	Min, Max Point3
}

// String returns a string representation of b like "(3,4,5)-(6,5,8)".
func (b Box) String() string {
	return b.Min.String() + "-" + b.Max.String()
}

// Empty reports whether the box contains no points.
func (b Box) Empty() bool {
	return b.Min.X >= b.Max.X || b.Min.Y >= b.Max.Y || b.Min.Z >= b.Max.Z
}

// Eq reports whether b and c contain the same set of points. All empty boxes
// are considered equal.
func (b Box) Eq(c Box) bool {
	return b == c || b.Empty() && c.Empty()
}

// Overlaps reports whether b and c have non-empty intersection.
func (b Box) Overlaps(c Box) bool {
	return !b.Empty() && !c.Empty() &&
		b.Min.X < c.Max.X && c.Min.X < b.Max.X &&
		b.Min.Y < c.Max.Y && c.Min.Y < b.Max.Y &&
		b.Min.Z < c.Max.Z && c.Min.Z < b.Max.Z
}

// Size returns b's width, height, and depth.
func (b Box) Size() Point3 {
	return b.Max.Sub(b.Min)
}

// IsoProjection translates an integer 3D coordinate into an integer 2D
// coordinate.
type IsoProjection struct {
	ZX, ZY int
}

// Project projects a 3D coordinate into 2D.
// If ZX = 0, x is unchanged; similarly for ZY and y.
// Otherwise, x projects to x + z/ZX and y projects to y + z/ZY.
// Dividing is used because there's little reason for an isometric projection
// in a game to exaggerate the Z position, and integers are used to preserve
// "pixel perfect" calculation in case you are making the next Celeste.
func (π IsoProjection) Project(x, y, z int) (xp, yp int) {
	// I'm using the π character because I'm a maths wanker
	xp, yp = x, y
	if π.ZX != 0 {
		xp += z / π.ZX
	}
	if π.ZY != 0 {
		yp += z / π.ZY
	}
	return xp, yp
}
