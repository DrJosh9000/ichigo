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

// CubicSpline implements a cubic spline. A cubic spline interpolates
// the given Points with cubic polynomials, ensuring first and second
// derivatives along the whole spline are continuous.
type CubicSpline struct {
	Points []Float2

	// If false, CubicSpline defines a natural cubic spline (the slopes at the
	// endpoints are "free" and the moments at the ends are zero.)
	// If true, Preslope (or Postslope, or both) is used to set the slope.
	FixedPreslope, FixedPostslope bool

	// Slope of line before and after spline, for extrapolation.
	// If a natural cubic spline is being used, these are set by Prepare.
	// If instead FixedPreslope or FixedPostslope are true, these are read by
	// Prepare to determine the moments.
	Preslope, Postslope float64

	// moments (second derivative at 1/6 scale) and intervals
	m, h []float64
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
	// Compute moments. "moments" is a term from drafting that basically means
	// the second derivative of the function. Points are known as "knots".
	//
	// Let's start with the natural cubic spline case, where m[0] and m[N-1] are
	// chosen to be 0. I'll skip over the derivation of the equations below, but
	// it follows from putting a cubic in each interval and having each one meet
	// the knots and match second derivatives with its neighbors.
	//
	// Note: these "moments" aren't the true values of the second derivatives
	// at the knots - they are calculated at 1/6th scale to avoid a multiply
	// and divide by 6.
	//
	// Given:
	//
	//    ɣ(i) = 2.0 * (h[i-1] + h[i])
	//    b(i) = ((Y[i+1]-Y[i])/h[i] - (Y[i]-Y[i-1])/h[i-1])
	//
	// we solve for m[i] in the equations:
	//
	//    h[i-1]*m[i-1] + ɣ(i)*m[i] + h[i]*m[i+1] = b(i)
	//
	// for i = 1...N-2.
	//
	// Written as a diagonally dominant tridiagonal matrix equation:
	//
	// [ ɣ(1)  h[1]     0     0     ...       0 ] [   m[1] ]   [   b(1) ]
	// [ h[1]  ɣ(2)  h[2]     0     ...       0 ] [   m[2] ]   [   b(2) ]
	// [    0  h[2]  ɣ(3)  h[3]     ...       0 ] [   m[3] ] = [   b(3) ]
	// [    0     0  h[3]  ɣ(4)     ...     ... ] [    ... ]   [    ... ]
	// [  ...   ...   ...   ...     ...  h[N-3] ] [  .  .. ]   [    ... ]
	// [    0     0   ...     0  h[N-3]  ɣ(N-2) ] [ m[N-2] ]   [ b(N-2) ]
	//
	// This is solvable in O(N) using simplified Gaussian elimination
	// ("Thomas algorithm").
	//
	// For the fixed end-slopes case, we also need to derive m[0] and m[N-1]
	// from the given end-slopes. Given:
	//
	//     b(0) = (Y[1] - Y[0]) / h[0] - Preslope
	//   b(N-1) = Postslope - (Y[N-1] - Y[N-2]) / h[N-2]
	//
	// We solve two additional equations for the new unknowns m[0] and m[N-1]:
	//
	//       2*m[0]*h[0] + m[1]*h[0]       = b(0)
	//     m[N-2]*h[N-2] + 2*m[N-1]*h[N-2] = b(N-1).
	//
	// Fortunately this is still a tridiagonal:
	//
	// [ 2h[0]  h[0]     0     0      0     ...        0 ] [   m[0] ]   [   b(0) ]
	// [  h[0]  ɣ(1)  h[1]     0      0     ...        0 ] [   m[1] ]   [   b(1) ]
	// [     0  h[1]  ɣ(2)  h[2]      0     ...        0 ] [   m[2] ]   [   b(2) ]
	// [     0     0  h[2]  ɣ(3)   h[3]     ...      ... ] [   m[3] ] = [   b(3) ]
	// [     0     0  h[3]  ɣ(4)    ...  h[N-3]        0 ] [    ... ]   [    ... ]
	// [   ...   ...   ...   ... h[N-3]  ɣ(N-2)   h[N-2] ] [ m[N-2] ]   [ b(N-2) ]
	// [     0     0     0   ...     0   h[N-2]  2h[N-2] ] [ m[N-1] ]   [ b(N-1) ].
	//
	// Fixing one end but leaving the other free leads to a mix of the two.

	// Setup:
	diag, upper, B := make([]float64, N), make([]float64, N), make([]float64, N)
	if s.FixedPreslope {
		diag[0] = 2.0 * s.h[0]
		upper[0] = s.h[0]
		B[0] = (s.Points[1].Y-s.Points[0].Y)/s.h[0] - s.Preslope
	}
	for i := 1; i < N-1; i++ {
		diag[i] = 2.0 * (s.h[i-1] + s.h[i])
		upper[i] = s.h[i]
		B[i] = (s.Points[i+1].Y-s.Points[i].Y)/s.h[i] - (s.Points[i].Y-s.Points[i-1].Y)/s.h[i-1]
	}
	if s.FixedPostslope {
		diag[N-1] = 2.0 * s.h[N-2]
		upper[N-1] = s.h[N-2]
		B[N-1] = s.Postslope - (s.Points[N-1].Y-s.Points[N-2].Y)/s.h[N-2]
	}
	// Forward elimination:
	if s.FixedPreslope {
		// Use row 0 to eliminate lower h[0] from row 1.
		// lower[1] = h[0]; diag[0] = 2h[0]; therefore lower[1]/diag[0] = 0.5.
		diag[1] -= 0.5 * upper[0]
		B[1] -= 0.5 * B[0]
	}
	for i := 2; i < N-1; i++ {
		// Use row i-1 to eliminate lower h[i-1] from row i
		t := s.h[i-1] / diag[i-1] // lower[i] / diag[i-1]
		diag[i] -= t * upper[i-1]
		B[i] -= t * B[i-1]
	}
	if s.FixedPostslope {
		// Use row N-2 to eliminate lower h[N-2] from row N-1.
		t := s.h[N-2] / diag[N-2]
		diag[N-1] -= t * upper[N-2]
		B[N-1] -= t * B[N-2]
	}
	// Back substitution:
	if s.FixedPostslope {
		s.m[N-1] = B[N-1] / diag[N-1]
	}
	for i := N - 2; i > 0; i-- {
		s.m[i] = (B[i] - s.h[i]*s.m[i+1]) / diag[i]
	}
	if s.FixedPreslope {
		s.m[0] = (B[0] - s.h[0]*s.m[1]) / diag[0]
	}
	// Derive pre- and post-slope, if not fixed:
	if !s.FixedPreslope {
		s.Preslope = -s.m[1]*s.h[0] + (s.Points[1].Y-s.Points[0].Y)/s.h[0]
	}
	if !s.FixedPostslope {
		s.Postslope = s.m[N-2]*s.h[N-2] + (s.Points[N-1].Y-s.Points[N-2].Y)/s.h[N-2]
	}
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
		return s.Points[0].Y + (x-s.Points[0].X)*s.Preslope
	}
	if x > s.Points[N-1].X {
		// Comes after the end of the spline, extrapolate
		return s.Points[N-1].Y + (x-s.Points[N-1].X)*s.Postslope
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
