package geom

import "image"

// Projection uses floats to define a projection.
type Projection struct{ X, Y float64 }

// Project performs a parallel projection of a 3D coordiante into 2D.
// x projects to (x + z*π.X), and y to (y + z*π.Y)
func (π Projection) Project(p Int3) image.Point {
	return image.Pt(
		p.X+int(π.X*float64(p.Z)),
		p.Y+int(π.Y*float64(p.Z)),
	)
}

// IntProjection holds an integer projection definition.
// It is designed for projecting Z onto X and Y with integer fractions as would
// be used in e.g. a diametric projection (IntProjection{X:0, Y:-2}).
type IntProjection image.Point

// Project performs an integer parallel projection of a 3D coordinate into 2D.
// If π.X = 0, the x returned is p.X; similarly for π.Y and y.
// Otherwise, x projects to x + z/π.X and y projects to y + z/π.Y.
func (π IntProjection) Project(p Int3) image.Point {
	/*
		I'm using the π character because I'm a maths wanker.

		Dividing is used because there's little reason for an isometric
		projection in a game to exaggerate the Z position.

		Integers are used to preserve "pixel perfect" calculation in case you
		are making the next Celeste.
	*/
	q := p.XY()
	if π.X != 0 {
		q.X += p.Z / π.X
	}
	if π.Y != 0 {
		q.Y += p.Z / π.Y
	}
	return q
}
