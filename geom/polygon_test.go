package geom

import (
	"image"
	"testing"
)

func TestPolygonContains(t *testing.T) {
	square5 := []image.Point{
		{X: 0, Y: 0},
		{X: 5, Y: 0},
		{X: 5, Y: 5},
		{X: 0, Y: 5},
	}
	tests := []struct {
		polygon []image.Point
		point   image.Point
		want    bool
	}{
		{square5, image.Pt(2, 3), true},
	}

	for _, test := range tests {
		if got, want := PolygonContains(test.polygon, test.point), test.want; got != want {
			t.Errorf("polygonContains(%v, %v) = %v, want %v", test.polygon, test.point, got, want)
		}
	}
}
