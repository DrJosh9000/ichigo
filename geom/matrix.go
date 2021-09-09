package geom

import (
	"errors"
	"image"
)

// IntMatrix3 implements a 3x3 integer matrix.
type IntMatrix3 [3][3]int

// Apply applies the matrix to a vector to obtain a transformed vector.
func (a IntMatrix3) Apply(v Int3) Int3 {
	return Int3{
		X: v.X*a[0][0] + v.Y*a[0][1] + v.Z*a[0][2],
		Y: v.X*a[1][0] + v.Y*a[1][1] + v.Z*a[1][2],
		Z: v.X*a[2][0] + v.Y*a[2][1] + v.Z*a[2][2],
	}
}

// Concat returns the matrix equivalent to applying matrix a and then b.
func (a IntMatrix3) Concat(b IntMatrix3) IntMatrix3 {
	return IntMatrix3{
		[3]int{
			a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0],
			a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1],
			a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2],
		},
		[3]int{
			a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0],
			a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1],
			a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2],
		},
		[3]int{
			a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0],
			a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1],
			a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2],
		},
	}
}

// IntMatrix3x4 implements a 3 row, 4 column integer matrix, capable of
// describing any integer affine transformation.
type IntMatrix3x4 [3][4]int

// Apply applies the matrix to a vector to obtain a transformed vector.
func (a IntMatrix3x4) Apply(v Int3) Int3 {
	return Int3{
		X: v.X*a[0][0] + v.Y*a[0][1] + v.Z*a[0][2] + a[0][3],
		Y: v.X*a[1][0] + v.Y*a[1][1] + v.Z*a[1][2] + a[1][3],
		Z: v.X*a[2][0] + v.Y*a[2][1] + v.Z*a[2][2] + a[2][3],
	}
}

// ToRatMatrix3 returns the 3x3 submatrix as a rational matrix equivalent.
func (a IntMatrix3x4) ToRatMatrix3() RatMatrix3 {
	return RatMatrix3{
		0: [3]Rat{{a[0][0], 1}, {a[0][1], 1}, {a[0][2], 1}},
		1: [3]Rat{{a[1][0], 1}, {a[1][1], 1}, {a[1][2], 1}},
		2: [3]Rat{{a[2][0], 1}, {a[2][1], 1}, {a[2][2], 1}},
	}
}

// Translation returns the translation component of the matrix (last column)
// i.e. what you get if you Apply the matrix to the zero vector.
func (a IntMatrix3x4) Translation() Int3 {
	return Int3{X: a[0][3], Y: a[1][3], Z: a[2][3]}
}

// IntMatrix2x3 implements a 2 row, 3 column matrix (as two row vectors).
type IntMatrix2x3 struct{ X, Y Int3 }

// Apply applies the matrix to a vector to obtain a transformed vector.
func (a IntMatrix2x3) Apply(v Int3) image.Point {
	return image.Point{
		X: v.Dot(a.X),
		Y: v.Dot(a.Y),
	}
}

// RatMatrix3 implements a 3x3 matrix with rational number entries.
type RatMatrix3 [3][3]Rat

// IdentityRatMatrix3x4 is the identity matrix for RatMatrix3x4.
var IdentityRatMatrix3 = RatMatrix3{
	0: [3]Rat{0: {1, 1}},
	1: [3]Rat{1: {1, 1}},
	2: [3]Rat{2: {1, 1}},
}

// IntApply applies the matrix to the integer vector v. Any remainder is lost.
func (a RatMatrix3) IntApply(v Int3) Int3 {
	x, y, z := IntRat(v.X), IntRat(v.Y), IntRat(v.Z)
	return Int3{
		X: x.Mul(a[0][0]).Add(y.Mul(a[0][1])).Add(z.Mul(a[0][2])).Int(),
		Y: x.Mul(a[1][0]).Add(y.Mul(a[1][1])).Add(z.Mul(a[1][2])).Int(),
		Z: x.Mul(a[2][0]).Add(y.Mul(a[2][1])).Add(z.Mul(a[2][2])).Int(),
	}
}

// Mul multiplies the whole matrix by a scalar.
func (a RatMatrix3) Mul(r Rat) RatMatrix3 {
	// "A little repetition..."
	a[0][0] = a[0][0].Mul(r)
	a[0][1] = a[0][1].Mul(r)
	a[0][2] = a[0][2].Mul(r)
	a[1][0] = a[1][0].Mul(r)
	a[1][1] = a[1][1].Mul(r)
	a[1][2] = a[1][2].Mul(r)
	a[2][0] = a[2][0].Mul(r)
	a[2][1] = a[2][1].Mul(r)
	a[2][2] = a[2][2].Mul(r)
	return a
}

// Adjugate returns the adjugate of the matrix.
func (a RatMatrix3) Adjugate() RatMatrix3 {
	return RatMatrix3{
		0: [3]Rat{
			0: a[1][1].Mul(a[2][2]).Sub(a[1][2].Mul(a[2][1])),
			1: a[0][1].Mul(a[2][2]).Sub(a[0][2].Mul(a[2][1])).Neg(),
			2: a[0][1].Mul(a[1][2]).Sub(a[0][2].Mul(a[1][1])),
		},
		1: [3]Rat{
			0: a[1][0].Mul(a[2][2]).Sub(a[1][2].Mul(a[2][0])).Neg(),
			1: a[0][0].Mul(a[2][2]).Sub(a[0][2].Mul(a[2][0])),
			2: a[0][0].Mul(a[1][2]).Sub(a[0][2].Mul(a[1][0])).Neg(),
		},
		2: [3]Rat{
			0: a[1][0].Mul(a[2][1]).Sub(a[1][1].Mul(a[2][0])),
			1: a[0][0].Mul(a[2][1]).Sub(a[0][1].Mul(a[2][0])).Neg(),
			2: a[0][0].Mul(a[1][1]).Sub(a[0][1].Mul(a[1][0])),
		},
	}
}

// Inverse returns the inverse of the matrix.
func (a RatMatrix3) Inverse() (RatMatrix3, error) {
	adj := a.Adjugate()
	det := a[0][0].Mul(adj[0][0]).Add(a[0][1].Mul(adj[1][0])).Add(a[0][2].Mul(adj[2][0]))
	if det.N == 0 {
		return RatMatrix3{}, errors.New("matrix is singular")
	}
	return adj.Mul(det.Invert()), nil
}
