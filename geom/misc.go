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

// ---------- Some helpers for image.Point ----------

// CMul performs componentwise multiplication of two image.Points.
func CMul(p, q image.Point) image.Point {
	return image.Point{p.X * q.X, p.Y * q.Y}
}

// CDiv performs componentwise division of two image.Points.
func CDiv(p, q image.Point) image.Point {
	return image.Point{p.X / q.X, p.Y / q.Y}
}

// CFloat returns the components of an image.Point as two floats.
func CFloat(p image.Point) (x, y float64) {
	return float64(p.X), float64(p.Y)
}

// Dot returns the Dot product of two image.Points.
func Dot(p, q image.Point) int {
	return p.X*q.X + p.Y*q.Y
}

// ---------- Some other helpers ----------

// FSign returns the sign of the float64 (-1, 0, or 1).
func FSign(m float64) float64 {
	if m == 0 {
		return 0
	}
	if m < 0 {
		return -1
	}
	return 1
}

// Sign returns the sign of the int (-1, 0, or 1).
func Sign(m int) int {
	if m == 0 {
		return 0
	}
	if m < 0 {
		return -1
	}
	return 1
}
