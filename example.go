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

// This is an example game using the engine. Run it with:
//    go run -tags example example.go
package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/DrJosh9000/ichigo/engine"
	"github.com/DrJosh9000/ichigo/example"
	"github.com/DrJosh9000/ichigo/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	enableCPUProfile  = true
	enableHeapProfile = true
	enableREPL        = true
	hardcodedLevel1   = true
	rewriteLevel1     = false
)

func main() {
	// Change to true to enable cpu profile
	if enableCPUProfile && runtime.GOOS != "js" {
		f, err := os.Create("cpuprofile.pprof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("TODO")

	// Change to true to rewrite level1.gobz
	lev1 := interface{}(&engine.SceneRef{Path: "assets/level1.gobz"})
	if hardcodedLevel1 {
		lev1 = example.Level1()
		if rewriteLevel1 && runtime.GOOS != "js" {
			if err := engine.SaveGobz(lev1, "example/assets/level1.gobz"); err != nil {
				log.Fatalf("Couldn't save level1.gobz: %v", err)
			}
		}
	}

	g := &engine.Game{
		ScreenSize: image.Pt(320, 240), // Window interior is this many pixels.
		// TODO: refactor Projection and VoxelScale into... Scene? Camera?
		// We might want different projections and scales in different levels.
		Projection: geom.SimpleProjection{},
		VoxelScale: geom.Float3{
			// Each voxel counts for this much (Euclidean) space.
			X: 1,
			Y: 1,
			Z: math.Sqrt(3),
		},
		Root: &engine.DrawDFS{
			Child: engine.MakeContainer(
				&engine.Fill{
					ID:    "bg_fill",
					Color: color.Gray{100},
				},
				&engine.LoadingSwitch{
					During: &engine.Billboard{
						ID:  "loading_screen",
						Src: engine.ImageRef{Path: "assets/loading.png"},
					},
					After: &engine.Camera{
						ID:    "game_camera",
						Child: lev1,
					},
				},
				&engine.DebugToast{ID: "toast", Pos: image.Pt(0, 15)},
				engine.PerfDisplay{},
			),
		},
	}
	if err := g.LoadAndPrepare(example.Assets); err != nil {
		log.Fatalf("Loading/preparing error: %v", err)
	}

	if enableREPL && runtime.GOOS != "js" {
		go g.REPL(os.Stdin, os.Stdout, example.Assets)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game error: %v", err)
	}

	if enableHeapProfile && runtime.GOOS != "js" {
		f, err := os.Create("heapprofile.pprof")
		if err != nil {
			log.Fatal("could not create heap profile: ", err)
		}
		defer f.Close()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write heap profile: ", err)
		}
	}
}
