package geom

import "image"

// Projector is used by Box and others to accept arbitrary
type Projector interface {
	// Sign returns a {-1, 0, 1}-valued 2D vector pointing in the direction that
	// positive Z values are projected to.
	Sign() image.Point
	// Project projects a Z coordinate into 2D offset.
	Project(int) image.Point
}

// Project is shorthand for π.Project(p.Z).Add(p.XY()).
func Project(π Projector, p Int3) image.Point {
	return π.Project(p.Z).Add(p.XY())
}

// Projection uses floats to define a projection.
type Projection struct{ X, Y float64 }

func (π Projection) Sign() (s image.Point) {
	return image.Pt(int(FSign(π.X)), int(FSign(π.Y)))
}

// Project returns (z*π.X, z*π.Y).
func (π Projection) Project(z int) image.Point {
	return image.Pt(
		int(π.X*float64(z)),
		int(π.Y*float64(z)),
	)
}

// IntProjection holds an integer projection definition.
// It is designed for projecting Z onto X and Y with integer fractions as would
// be used in e.g. a diametric projection (IntProjection{X:0, Y:2}).
type IntProjection image.Point

func (π IntProjection) Sign() image.Point { return image.Point(π) }

// Project returns (z/π.X, z/π.Y), unless π.X or π.Y are 0, in which case that
// component is zero
func (π IntProjection) Project(z int) image.Point {
	/*
		Dividing is used because there's little reason for an isometric
		projection in a game to exaggerate the Z position.

		Integers are used to preserve "pixel perfect" calculation in case you
		are making the next Celeste.
	*/
	var q image.Point
	if π.X != 0 {
		q.X = z / π.X
	}
	if π.Y != 0 {
		q.Y = z / π.Y
	}
	return q
}
