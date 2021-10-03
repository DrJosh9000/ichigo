package geom

import (
	"errors"
	"fmt"
	"sort"
)

// LinearSpline implements a linear spline.
type LinearSpline struct {
	Points []Float2

	deriv []float64 // slope of segment between points i and i+1
}

// Prepare sorts the points and computes internal information.
func (s *LinearSpline) Prepare() error {
	if len(s.Points) < 1 {
		return errors.New("spline needs at least 1 point")
	}

	// Ensure Points is sorted.
	sort.Slice(s.Points, func(i, j int) bool {
		return s.Points[i].X < s.Points[j].X
	})

	// Check for points with equal X and compute derivatives
	s.deriv = make([]float64, len(s.Points)-1)
	for i := range s.Points[1:] {
		if s.Points[i].X == s.Points[i+1].X {
			return fmt.Errorf("spline value defined twice [%v, %v]", s.Points[i], s.Points[i+1])
		}
		s.deriv[i] = (s.Points[i+1].Y - s.Points[i].Y) / (s.Points[i+1].X - s.Points[i].X)
	}
	return nil
}

// Interpolate, given x, returns y where (x,y) is a point on the spline.
// If x is outside the spline, it extrapolates from either the first or
// last segments of the spline.
func (s *LinearSpline) Interpolate(x float64) float64 {
	N := len(s.Points)
	if N == 1 {
		return s.Points[0].Y
	}
	if x < s.Points[1].X {
		// Comes before the end of the first segment
		return s.Points[0].Y + (x-s.Points[0].X)*s.deriv[0]
	}
	if x > s.Points[N-2].X {
		// Comes after the start of the last segment
		return s.Points[N-1].Y + (x-s.Points[N-1].X)*s.deriv[N-2]
	}
	// Somewhere in the middle
	i := sort.Search(N, func(n int) bool {
		return x <= s.Points[n].X
	})
	if x == s.Points[i].X {
		// Hit a point exactly
		return s.Points[i].Y
	}
	// In the interval between point i and point i+1
	return s.Points[i].Y + (x-s.Points[i].X)*s.deriv[i]
}

// CubicSpline implements a cubic spline.
type CubicSpline struct {
	// Points on the spline
	Points []Float2

	deriv, deriv2 []float64
}

// Prepare
func (s *CubicSpline) Prepare() error {
	if len(s.Points) < 1 {
		return errors.New("spline needs at least 1 point")
	}
	// Ensure Points is sorted.
	sort.Slice(s.Points, func(i, j int) bool {
		return s.Points[i].X < s.Points[j].X
	})
	// Check for points with equal X.
	for i := range s.Points[1:] {
		if s.Points[i].X == s.Points[i+1].X {
			return fmt.Errorf("spline value defined twice [%v, %v]", s.Points[i], s.Points[i+1])
		}
	}
	// TODO: compute deriv and deriv2
	return nil
}

func (s *CubicSpline) Interpolate(x float64) float64 {
	N := len(s.Points)
	if N == 1 {
		return s.Points[0].Y
	}
	// TODO
	return 0
}
