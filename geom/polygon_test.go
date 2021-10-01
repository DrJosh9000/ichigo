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
	"testing"
)

func TestPolygonContains(t *testing.T) {
	square5 := []image.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 5},
		{X: 5, Y: 5},
		{X: 5, Y: 0},
	}
	pentagon := []image.Point{
		{X: -5, Y: -4},
		{X: -7, Y: 2},
		{X: 0, Y: 7},
		{X: 7, Y: 2},
		{X: 5, Y: -4},
	}
	hexagon := []image.Point{
		{X: 8, Y: 0},
		{X: 0, Y: 8},
		{X: 8, Y: 16},
		{X: 23, Y: 16},
		{X: 31, Y: 8},
		{X: 23, Y: 0},
	}

	tests := []struct {
		polygon []image.Point
		point   image.Point
		want    bool
	}{
		{square5, image.Pt(2, 3), true},
		{square5, image.Pt(0, 0), true},
		{square5, image.Pt(5, 5), true},
		{square5, image.Pt(-5, 0), false},
		{square5, image.Pt(0, -5), false},
		{square5, image.Pt(6, 6), false},

		{pentagon, image.Pt(0, 0), true},
		{pentagon, image.Pt(8, 0), false},
		{pentagon, image.Pt(1, 1), true},
		{pentagon, image.Pt(-1, -1), true},
		{pentagon, image.Pt(-10000, 10000), false},

		{hexagon, image.Pt(0, 0), false},
		{hexagon, image.Pt(16, 8), true},
		{hexagon, image.Pt(0, 8), true},
		{hexagon, image.Pt(-1, 8), false},
		{hexagon, image.Pt(31, 8), true},
		{hexagon, image.Pt(32, 8), false},
		{hexagon, image.Pt(10000, 10000), false},
	}

	for _, test := range tests {
		if got, want := PolygonContains(test.polygon, test.point), test.want; got != want {
			t.Errorf("PolygonContains(%v, %v) = %v, want %v", test.polygon, test.point, got, want)
		}
	}
}
