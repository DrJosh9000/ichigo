/*
Copyright 2021 Josh Deprez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package geom

import "image"

// Projector types can be used to project 3D coordinates into 2D. It only
// supports projecting Z into a 2D offset (i.e. not a general projection).
type Projector interface {
	// Sign returns a {-1, 0, 1}-valued 2D vector pointing in the direction that
	// positive Z values are projected to.
	Sign() image.Point
	// Project converts a Z coordinate to a 2D offset.
	Project(int) image.Point
}

// Project is shorthand for π.Project(p.Z).Add(p.XY()).
func Project(π Projector, p Int3) image.Point {
	return π.Project(p.Z).Add(p.XY())
}

// ElevationProjection throws away Z.
type ElevationProjection struct{}

// Sign returns the zero point.
func (ElevationProjection) Sign() image.Point { return image.Point{} }

// Project returns the zero point.
func (ElevationProjection) Project(int) image.Point { return image.Point{} }

// SimpleProjection projects Z onto Y only.
type SimpleProjection struct{}

// Sign returns (0, 1).
func (SimpleProjection) Sign() image.Point { return image.Pt(0, 1) }

// Project returns (0, z).
func (SimpleProjection) Project(z int) image.Point { return image.Pt(0, z) }

// Projection uses two floats to define a custom projection.
type Projection struct{ X, Y float64 }

// Sign returns the componentwise sign of π.
func (π Projection) Sign() image.Point {
	return image.Pt(int(FSign(π.X)), int(FSign(π.Y)))
}

// Project returns (z*π.X, z*π.Y).
func (π Projection) Project(z int) image.Point {
	return image.Pt(
		int(π.X*float64(z)),
		int(π.Y*float64(z)),
	)
}

// IntProjection uses two integers to define a custom projection.
// It is designed for projecting Z onto X and Y with integer fractions as would
// be used in e.g. a diametric projection (IntProjection{X:0, Y:2}).
type IntProjection image.Point

// Sign returns CSign(π).
func (π IntProjection) Sign() image.Point { return CSign(image.Point(π)) }

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
