package geom

import (
	"math"
	"testing"
)

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
		{x: -4.5, want: 0.75},
		{x: -4, want: 0.5},
		{x: -3.5, want: 0.25},
		{x: -3, want: 0},
		{x: -2.5, want: -1.5},
		{x: -2, want: -3},
		{x: -1.5, want: -1.75},
		{x: -1, want: -0.5},
		{x: -0.5, want: 0.75},
		{x: 0, want: 2},
		{x: 0.5, want: -1.5},
		{x: 1, want: -5},
		{x: 1.5, want: -4.25},
		{x: 2, want: -3.5},
		{x: 2.5, want: -2.75},
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

func TestCubicSplineNoPoints(t *testing.T) {
	s := &CubicSpline{}
	if err := s.Prepare(); err == nil {
		t.Errorf("s.Prepare() = %v, want error", err)
	}
}

func TestCubicSplineEqualXPoints(t *testing.T) {
	s := &CubicSpline{
		Points: []Float2{{-5, 1}, {-2, 7}, {-2, -3}, {0, 2}, {3, -2}},
	}
	if err := s.Prepare(); err == nil {
		t.Errorf("s.Prepare() = %v, want error", err)
	}
}

func TestCubicSplineOnePoint(t *testing.T) {
	s := &CubicSpline{
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

func TestNaturalCubicSpline(t *testing.T) {
	s := &CubicSpline{
		Points: []Float2{{-7, -2}, {-5, 1}, {-3, 0}, {-2, -3}, {0, 2}, {1, -5}, {3, -2}, {4, 4}},
	}
	if err := s.Prepare(); err != nil {
		t.Errorf("s.Prepare() = %v, want nil", err)
	}
	tests := []struct {
		x, want float64
	}{
		{x: -8, want: -3.648342225609756},
		{x: -7.5, want: -2.824171112804878},
		{x: -7, want: -2},
		{x: -6.5, want: -1.180464581745427},
		{x: -6, want: -0.3887433307926829},
		{x: -5.5, want: 0.3473495855564025},
		{x: -5, want: 1},
		{x: -4.5, want: 1.5067079125381098},
		{x: -4, want: 1.6662299923780488},
		{x: -3.5, want: 1.2426370760289636},
		{x: -3, want: 0},
		{x: -2.5, want: -1.9368449885670733},
		{x: -2, want: -3},
		{x: -1.5, want: -1.855450886051829},
		{x: -1, want: 0.45221989329268286},
		{x: -0.5, want: 2.2837807259908534},
		{x: 0, want: 2},
		{x: 0.5, want: -1.229539824695122},
		{x: 1, want: -5},
		{x: 1.5, want: -6.734946646341463},
		{x: 2, want: -6.406821646341463},
		{x: 2.5, want: -4.6252858231707314},
		{x: 3, want: -2},
		{x: 3.5, want: 0.941477705792683},
		{x: 4, want: 4},
		{x: 4.5, want: 7.078029725609756},
		{x: 5, want: 10.156059451219512},
		{x: 5.5, want: 13.234089176829269},
	}
	for _, test := range tests {
		if got := s.Interpolate(test.x); math.Abs(got-test.want) > 0.0000001 {
			t.Errorf("s.Interpolate(%v) = %v, want %v", test.x, got, test.want)
		}
	}
}

func TestFixedEndSlopesCubicSpline(t *testing.T) {
	s := &CubicSpline{
		Points:         []Float2{{-7, -2}, {-5, 1}, {-3, 0}, {-2, -3}, {0, 2}, {1, -5}, {3, -2}, {4, 4}},
		FixedPreslope:  true,
		FixedPostslope: true,
		Preslope:       -5,
		Postslope:      4,
	}
	if err := s.Prepare(); err != nil {
		t.Errorf("s.Prepare() = %v, want nil", err)
	}
	tests := []struct {
		x, want float64
	}{
		{x: -8, want: 3},
		{x: -7.5, want: 0.5},
		{x: -7, want: -2},
		{x: -6.5, want: -3.213753455247408},
		{x: -6, want: -2.4866758806597526},
		{x: -5.5, want: -0.7662603657422216},
		{x: -5, want: 1},
		{x: -4.5, want: 2.02752418673185},
		{x: -4, want: 2.183379403298762},
		{x: -3.5, want: 1.4975449182162928},
		{x: -3, want: 0},
		{x: -2.5, want: -1.9904880751130807},
		{x: -2, want: -3},
		{x: -1.5, want: -1.8287325238726178},
		{x: -1, want: 0.4715013478320472},
		{x: -0.5, want: 2.2859845456206873},
		{x: 0, want: 2},
		{x: 0.5, want: -1.2096392835930003},
		{x: 1, want: -5},
		{x: 1.5, want: -6.8416283067574355},
		{x: 2, want: -6.6251085119020425},
		{x: 2.5, want: -4.846034461095628},
		{x: 3, want: -2},
		{x: 3.5, want: 1.3035791794215745},
		{x: 4, want: 4},
		{x: 4.5, want: 6},
		{x: 5, want: 8},
		{x: 5.5, want: 10},
	}
	for _, test := range tests {
		if got := s.Interpolate(test.x); math.Abs(got-test.want) > 0.0000001 {
			t.Errorf("s.Interpolate(%v) = %v, want %v", test.x, got, test.want)
		}
	}
}
