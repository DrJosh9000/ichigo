//go:build example
// +build example

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
	linear := &geom.LinearSpline{Points: points}
	if err := linear.Prepare(); err != nil {
		log.Fatalf("linear.Prepare() = %v, want nil", err)
	}
	cubic := &geom.CubicSpline{Points: points}
	if err := cubic.Prepare(); err != nil {
		log.Fatalf("cubic.Prepare() = %v, want nil", err)
	}
	// Produce interpolated points in CSV-like form.
	for x := -8.0; x < 8.0; x += 0.125 {
		fmt.Printf("%f,%f,%f\n", x, linear.Interpolate(x), cubic.Interpolate(x))
	}
}
