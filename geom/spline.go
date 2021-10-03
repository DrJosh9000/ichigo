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
		// Hit the point i exactly
		return s.Points[i].Y
	}
	// In the interval between point i-1 and point i
	return s.Points[i-1].Y + (x-s.Points[i-1].X)*s.deriv[i-1]
}

// CubicSpline implements a natural cubic spline. A cubic spline interpolates
// the given Points while ensuring first and second derivatives are continuous.
type CubicSpline struct {
	Points []Float2

	// moments and intervals
	m, h []float64

	// slope of line before and after spline, for extrapolation
	preslope, postslope float64
}

// Prepare sorts the points and computes internal information.
func (s *CubicSpline) Prepare() error {
	if len(s.Points) < 1 {
		return errors.New("spline needs at least 1 point")
	}
	// Ensure Points is sorted.
	sort.Slice(s.Points, func(i, j int) bool {
		return s.Points[i].X < s.Points[j].X
	})
	// Check for points with equal X, and compute intervals.
	N := len(s.Points)
	if N == 1 {
		return nil
	}
	s.m = make([]float64, N)
	s.h = make([]float64, N-1)
	for i := range s.Points[1:] {
		if s.Points[i].X == s.Points[i+1].X {
			return fmt.Errorf("spline value defined twice [%v, %v]", s.Points[i], s.Points[i+1])
		}
		s.h[i] = s.Points[i+1].X - s.Points[i].X
	}
	// Compute moments. m[0] and m[N-1] are chosen to be 0 (natural cubic spline).
	// Given:
	//    ɣ(i) = 2.0 * (h[i-1] + h[i])
	//    b(i) = 6.0 * ((Points[i+1].Y-Points[i].Y)/h[i] - (Points[i].Y-Points[i-1].Y)/h[i-1])
	// we solve for m[i] in the equations:
	//    h[i-1]*m[i-1] + ɣ(i)*m[i] + h[i]*m[i+1] = b(i)
	// for i = 1...N-2.
	//
	// Written as a diagonally dominant tridiagonal matrix equation:
	//
	// [ɣ(1)  h[1]  0     0    ...    0      ] [  m[1]  ]   [  b(1)  ]
	// [h[1]  ɣ(2)  h[2]  0    ...    0      ] [  m[2]  ]   [  b(2)  ]
	// [0     h[2]  ɣ(3)  h[3] ...    0      ] [  m[3]  ] = [  b(3)  ]
	// [0     0     h[3]  ɣ(4) ...    ...    ] [  ...   ]   [  ...   ]
	// [...................... ...    h[N-3] ] [  ...   ]   [  ...   ]
	// [0     0     ...   0    h[N-3] ɣ(N-2) ] [ m[N-2] ]   [ b(N-2) ]
	//
	// This is solvable in O(N) using simplified Gaussian elimination
	// ("Thomas algorithm").

	// Setup:
	diag, upper, B := make([]float64, N-1), make([]float64, N-1), make([]float64, N-1)
	for i := 1; i < N-1; i++ {
		diag[i] = 2.0 * (s.h[i-1] + s.h[i])
		upper[i] = s.h[i]
		B[i] = 6.0 * ((s.Points[i+1].Y-s.Points[i].Y)/s.h[i] - (s.Points[i].Y-s.Points[i-1].Y)/s.h[i-1])
	}
	// Forward elimination:
	for i := 2; i < N-1; i++ {
		t := s.h[i-1] / diag[i-1] // lower[i] / diag[i-1]
		diag[i] -= t * upper[i-1]
		B[i] -= t * B[i-1]
	}
	// Back substitution:
	for i := N - 2; i > 0; i-- {
		s.m[i] = (B[i] - s.h[i]*s.m[i+1]) / diag[i]
	}
	// Divide all the moments by 6, since all the terms with moments in them
	// from this point onwards are divided by six.
	for i := range s.m {
		s.m[i] /= 6.0
	}
	// Pre- and post-slope:
	s.preslope = -s.m[1]*s.h[0] + (s.Points[1].Y-s.Points[0].Y)/s.h[0]
	s.postslope = s.m[N-2]*s.h[N-2] + (s.Points[N-1].Y-s.Points[N-2].Y)/s.h[N-2]
	return nil
}

// Interpolate, given x, returns y where (x,y) is a point on the spline.
// If x is outside the spline, it extrapolates from either the first or
// last segments of the spline.
func (s *CubicSpline) Interpolate(x float64) float64 {
	N := len(s.Points)
	if N == 1 {
		return s.Points[0].Y
	}
	if x < s.Points[0].X {
		// Comes before the start of the spline, extrapolate
		return s.Points[0].Y + (x-s.Points[0].X)*s.preslope
	}
	if x > s.Points[N-1].X {
		// Comes after the end of the spline, extrapolate
		return s.Points[N-1].Y + (x-s.Points[N-1].X)*s.postslope
	}
	// Somewhere in the middle
	i := sort.Search(N, func(n int) bool {
		return x <= s.Points[n].X
	})
	if x == s.Points[i].X {
		// Hit the point i exactly
		return s.Points[i].Y
	}
	// In the interval between point i-1 and point i
	x0, x1 := x-s.Points[i-1].X, s.Points[i].X-x
	return (s.m[i-1]*(x1*x1*x1)+s.m[i]*(x0*x0*x0))/s.h[i-1] -
		(s.m[i-1]*x1+s.m[i]*x0)*s.h[i-1] +
		(s.Points[i-1].Y*x1+s.Points[i].Y*x0)/s.h[i-1]
}
