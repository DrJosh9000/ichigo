package geom

import "testing"

func TestLinearSplineNoPoints(t *testing.T) {
	s := &LinearSpline{}
	if err := s.Prepare(); err == nil {
		t.Errorf("s.Prepare() = %v, want error", err)
	}
}

func TestLinearSplineEqualXPoints(t *testing.T) {
	s := &LinearSpline{
		Points: []Float2{{-5, 1}, {-2, 7}, {-2, -3}, {0, 2}, {3, -2}},
	}
	if err := s.Prepare(); err == nil {
		t.Errorf("s.Prepare() = %v, want error", err)
	}
}

func TestLinearSplineOnePoint(t *testing.T) {
	s := &LinearSpline{
		Points: []Float2{{-2, -3}},
	}
	if err := s.Prepare(); err != nil {
		t.Errorf("s.Prepare() = %v, want nil", err)
	}
	for _, x := range []float64{-5, -4, -2, 0, 1, 7} {
		if got, want := s.Interpolate(x), float64(-3); got != want {
			t.Errorf("s.Interpolate(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestLinearSpline(t *testing.T) {
	s := &LinearSpline{
		Points: []Float2{{-7, -2}, {-5, 1}, {-3, 0}, {-2, -3}, {0, 2}, {1, -5}, {3, -2}, {4, 4}},
	}
	if err := s.Prepare(); err != nil {
		t.Errorf("s.Prepare() = %v, want nil", err)
	}
	tests := []struct {
		x, want float64
	}{
		{x: -8, want: -3.5},
		{x: -7.5, want: -2.75},
		{x: -7, want: -2},
		{x: -6.5, want: -1.25},
		{x: -6, want: -0.5},
		{x: -5.5, want: 0.25},
		{x: -5, want: 1},
		{x: -4.5, want: 4.5},
		{x: -4, want: 3},
		{x: -3.5, want: 1.5},
		{x: -3, want: 0},
		{x: -2.5, want: -4.25},
		{x: -2, want: -3},
		{x: -1.5, want: 12.5},
		{x: -1, want: 9},
		{x: -0.5, want: 5.5},
		{x: 0, want: 2},
		{x: 0.5, want: -5.75},
		{x: 1, want: -5},
		{x: 1.5, want: -11},
		{x: 2, want: -8},
		{x: 2.5, want: -5},
		{x: 3, want: -2},
		{x: 3.5, want: 1},
		{x: 4, want: 4},
		{x: 4.5, want: 7},
		{x: 5, want: 10},
		{x: 5.5, want: 13},
	}
	for _, test := range tests {
		if got := s.Interpolate(test.x); got != test.want {
			t.Errorf("s.Interpolate(%v) = %v, want %v", test.x, got, test.want)
		}
	}
}
