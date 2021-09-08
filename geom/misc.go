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
