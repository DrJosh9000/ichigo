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

import (
	"image"
	"strconv"
)

// Int3 is a an element of int^3.
type Int3 struct {
	X, Y, Z int
}

// Pt3(x, y, z) is shorthand for Int3{x, y, z}.
func Pt3(x, y, z int) Int3 {
	return Int3{x, y, z}
}

// String returns a string representation of p like "(3,4,5)".
func (p Int3) String() string {
	return "(" + strconv.Itoa(p.X) + "," + strconv.Itoa(p.Y) + "," + strconv.Itoa(p.Z) + ")"
}

// XY applies the Z-forgetting projection. (It returns just X and Y.)
func (p Int3) XY() image.Point {
	return image.Point{X: p.X, Y: p.Y}
}

// XZ applies the Y-forgetting projection. (It returns just X and Z (as Y).)
func (p Int3) XZ() image.Point {
	return image.Point{X: p.X, Y: p.Z}
}

// Add performs vector addition.
func (p Int3) Add(q Int3) Int3 {
	return Int3{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

// Sub performs vector subtraction.
func (p Int3) Sub(q Int3) Int3 {
	return p.Add(q.Neg())
}

// CMul performs componentwise multiplication.
func (p Int3) CMul(q Int3) Int3 {
	return Int3{p.X * q.X, p.Y * q.Y, p.Z * q.Z}
}

// Mul performs scalar multiplication.
func (p Int3) Mul(k int) Int3 {
	return Int3{p.X * k, p.Y * k, p.Z * k}
}

// CDiv performs componentwise division.
func (p Int3) CDiv(q Int3) Int3 {
	return Int3{p.X / q.X, p.Y / q.Y, p.Z / q.Z}
}

// Div performs scalar division by k.
func (p Int3) Div(k int) Int3 {
	return Int3{p.X / k, p.Y / k, p.Z / k}
}

// Neg returns the vector pointing in the opposite direction.
func (p Int3) Neg() Int3 {
	return Int3{-p.X, -p.Y, -p.Z}
}

// Coord returns the components of the vector.
func (p Int3) Coord() (x, y, z int) {
	return p.X, p.Y, p.Z
}

// Sign returns a sign vector.
func (p Int3) Sign() Int3 {
	return Int3{Sign(p.X), Sign(p.Y), Sign(p.Z)}
}

// Dot returns the dot product of the two vectors.
func (p Int3) Dot(q Int3) int {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}
