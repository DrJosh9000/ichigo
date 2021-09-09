package geom

import "image"

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

// DrawOrder computes a draw-order value for a point under this projection.
// Each projection has an implied camera angle - Z alone is insufficient to
// order things properly.
func (π IntProjection) DrawOrder(p Int3) float64 {
	z := float64(p.Z)
	if π.X != 0 {
		z -= float64(p.X) / float64(π.X)
	}
	if π.Y != 0 {
		z -= float64(p.Y) / float64(π.Y)
	}
	return z
}
