package engine

import "image"

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
//describing any integer affine transformation.
type IntMatrix3x4 [3][4]int

// Apply applies the matrix to a vector to obtain a transformed vector.
func (a IntMatrix3x4) Apply(v Int3) Int3 {
	return Int3{
		X: v.X*a[0][0] + v.Y*a[0][1] + v.Z*a[0][2] + a[0][3],
		Y: v.X*a[1][0] + v.Y*a[1][1] + v.Z*a[1][2] + a[1][3],
		Z: v.X*a[2][0] + v.Y*a[2][1] + v.Z*a[2][2] + a[2][3],
	}
}

// IntMatrix2x3 implements a 2 row, 3 column matrix (as two row vectors).
type IntMatrix2x3 [2]Int3

// Apply applies the matrix to a vector to obtain a transformed vector.
func (a IntMatrix2x3) Apply(v Int3) image.Point {
	return image.Point{
		X: v.Dot(a[0]),
		Y: v.Dot(a[1]),
	}
}
