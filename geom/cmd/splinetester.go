//go:build example
// +build example

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

package main

import (
	"fmt"
	"log"

	"github.com/DrJosh9000/ichigo/geom"
)

func main() {
	// Put your own points here!
	points := []geom.Float2{
		{X: -7, Y: -2},
		{X: -5, Y: 1},
		{X: -3, Y: 0},
		{X: -2, Y: -3},
		{X: 0, Y: 2},
		{X: 1, Y: -5},
		{X: 3, Y: -2},
		{X: 4, Y: 4},
	}
	cubic := &geom.CubicSpline{
		Points:         points,
		FixedPreslope:  true,
		FixedPostslope: true,
		Preslope:       -5,
		Postslope:      -4,
	}
	if err := cubic.Prepare(); err != nil {
		log.Fatalf("cubic.Prepare() = %v, want nil", err)
	}
	// Produce interpolated points in CSV-like form.
	for x := -8.0; x < 8.0; x += 0.0625 {
		fmt.Printf("%f,%f\n", x, cubic.Interpolate(x))
	}
}
