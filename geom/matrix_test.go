package geom

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRatMatrix3InverseAndConcat(t *testing.T) {
	identity := RatMatrix3{
		0: [3]Rat{{1, 1}, {0, 1}, {0, 1}},
		1: [3]Rat{{0, 1}, {1, 1}, {0, 1}},
		2: [3]Rat{{0, 1}, {0, 1}, {1, 1}},
	}
	tests := []struct {
		in, want RatMatrix3
	}{
		{ // identity
			in:   identity,
			want: identity,
		},
		{ // diagonal
			in: RatMatrix3{
				0: [3]Rat{{3, 1}, {0, 1}, {0, 1}},
				1: [3]Rat{{0, 1}, {1, 2}, {0, 1}},
				2: [3]Rat{{0, 1}, {0, 1}, {7, 13}},
			},
			want: RatMatrix3{
				0: [3]Rat{{1, 3}, {0, 1}, {0, 1}},
				1: [3]Rat{{0, 1}, {2, 1}, {0, 1}},
				2: [3]Rat{{0, 1}, {0, 1}, {13, 7}},
			},
		},
		{ // hexagonal prism layout
			in: RatMatrix3{
				0: [3]Rat{{24, 1}, {0, 1}, {0, 1}},
				1: [3]Rat{{0, 1}, {16, 1}, {0, 1}},
				2: [3]Rat{{8, 1}, {0, 1}, {16, 1}},
			},
			want: RatMatrix3{
				0: [3]Rat{{1, 24}, {0, 1}, {0, 1}},
				1: [3]Rat{{0, 1}, {1, 16}, {0, 1}},
				2: [3]Rat{{-1, 48}, {0, 1}, {1, 16}},
			},
		},
		{ // random small positive integer matrix
			in: RatMatrix3{
				0: [3]Rat{{90, 1}, {40, 1}, {15, 1}},
				1: [3]Rat{{28, 1}, {54, 1}, {77, 1}},
				2: [3]Rat{{35, 1}, {99, 1}, {9, 1}},
			},
			want: RatMatrix3{
				0: [3]Rat{{7137, 531380}, {-225, 106276}, {-227, 53138}},
				1: [3]Rat{{-2443, 531380}, {-57, 106276}, {651, 53138}},
				2: [3]Rat{{-441, 265690}, {751, 53138}, {-187, 26569}},
			},
		},
		{ // random small integer matrix
			in: RatMatrix3{
				0: [3]Rat{{-38, 1}, {-23, 1}, {60, 1}},
				1: [3]Rat{{-75, 1}, {-22, 1}, {33, 1}},
				2: [3]Rat{{-5, 1}, {-6, 1}, {-14, 1}},
			},
			want: RatMatrix3{
				0: [3]Rat{{46, 2647}, {-62, 2647}, {51, 2647}},
				1: [3]Rat{{-1215, 29117}, {832, 29117}, {-3246, 29117}},
				2: [3]Rat{{340, 29117}, {-113, 29117}, {-889, 29117}},
			},
		},
	}

	for _, test := range tests {
		got, err := test.in.Inverse()
		if err != nil {
			t.Errorf("(%v).Inverse() error = %v, want nil", test.in, err)
		}
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("(%v).Inverse() diff:\n%s", test.in, diff)
		}
		got2 := test.in.Concat(got)
		if diff := cmp.Diff(got2, identity); diff != "" {
			t.Errorf("(%v).Concat(%v) diff:\n%s", test.in, got2, diff)
		}
		got3 := got.Concat(test.in)
		if diff := cmp.Diff(got3, identity); diff != "" {
			t.Errorf("(%v).Concat(%v) diff\n%s", got, test.in, diff)
		}
	}
}

func TestRatMatrix3InvertSingular(t *testing.T) {
	tests := []RatMatrix3{
		{ // zero row and column
			0: [3]Rat{{1, 1}, {0, 1}, {0, 1}},
			1: [3]Rat{{0, 1}, {1, 1}, {0, 1}},
			2: [3]Rat{{0, 1}, {0, 1}, {0, 1}},
		},
		{ // zero row
			0: [3]Rat{{1, 1}, {0, 1}, {0, 1}},
			1: [3]Rat{{0, 1}, {1, 1}, {1, 1}},
			2: [3]Rat{{0, 1}, {0, 1}, {0, 1}},
		},
		{ // zero column
			0: [3]Rat{{1, 1}, {0, 1}, {0, 1}},
			1: [3]Rat{{0, 1}, {1, 1}, {0, 1}},
			2: [3]Rat{{0, 1}, {1, 1}, {0, 1}},
		},
		{ // product of random 3x2 and 2x3 integer matrices
			0: [3]Rat{{-4330, 1}, {1283, 1}, {2717, 1}},
			1: [3]Rat{{1978, 1}, {-10171, 1}, {571, 1}},
			2: [3]Rat{{-4962, 1}, {-3689, 1}, {4089, 1}},
		},
	}
	for _, test := range tests {
		if _, err := test.Inverse(); err != errSingularMatrix {
			t.Errorf("(%v).Inverse() error = %v, want 'matrix is singular'", test, err)
		}
	}
}
