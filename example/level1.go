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

package example

import (
	"image"
	"time"

	"github.com/DrJosh9000/ichigo/engine"
	"github.com/DrJosh9000/ichigo/geom"
)

// Level1 creates the level_1 scene.
func Level1() *engine.Scene {
	return &engine.Scene{
		ID:     "level_1",
		Bounds: engine.Bounds(image.Rect(-32, -32, 320+32, 240+32)),
		Child: engine.MakeContainer(
			engine.DummyLoad{
				Duration: 2 * time.Second,
			},
			&engine.Parallax{
				CameraID: "game_camera",
				Child: &engine.Billboard{
					ID:  "bg_image",
					Pos: geom.Pt3(-160, -20, -100),
					Src: engine.ImageRef{Path: "assets/space.png"},
				},
				Factor: 0.5,
			}, // Parallax
			&engine.DrawDAG{
				ChunkSize: 16,
				Child: engine.MakeContainer(
					&engine.PrismMap{
						ID: "hexagons",
						PosToWorld: geom.IntMatrix3x4{
							// For each tile in the X direction, go right by 24 and
							// forward by 8, etc
							0: [4]int{24, 0, 0, 0},
							1: [4]int{0, 16, 0, 0},
							2: [4]int{8, 0, 16, 0},
						},
						PrismSize: geom.Int3{X: 32, Y: 16, Z: 16},
						PrismTop: []image.Point{
							{X: 8, Y: 0},
							{X: 0, Y: 8},
							{X: 8, Y: 16},
							{X: 23, Y: 16},
							{X: 31, Y: 8},
							{X: 23, Y: 0},
						},
						Sheet: engine.Sheet{
							CellSize: image.Pt(32, 32),
							Src:      engine.ImageRef{Path: "assets/hexprism32.png"},
						},
						Map: map[geom.Int3]*engine.Prism{
							geom.Pt3(11, 0, -6): {},
							geom.Pt3(12, 0, -6): {},

							geom.Pt3(9, 0, -5):  {},
							geom.Pt3(10, 0, -5): {},
							geom.Pt3(11, 0, -5): {},
							geom.Pt3(12, 0, -5): {},

							geom.Pt3(7, 0, -4):  {},
							geom.Pt3(8, 0, -4):  {},
							geom.Pt3(9, 0, -4):  {},
							geom.Pt3(10, 0, -4): {},
							geom.Pt3(11, 0, -4): {},
							geom.Pt3(12, 0, -4): {},

							geom.Pt3(5, 0, -3):  {},
							geom.Pt3(6, 0, -3):  {},
							geom.Pt3(7, 0, -3):  {},
							geom.Pt3(8, 0, -3):  {},
							geom.Pt3(9, 0, -3):  {},
							geom.Pt3(10, 0, -3): {},
							geom.Pt3(11, 0, -3): {},
							geom.Pt3(12, 0, -3): {},

							geom.Pt3(3, 0, -2):  {},
							geom.Pt3(4, 0, -2):  {},
							geom.Pt3(5, 0, -2):  {},
							geom.Pt3(6, 0, -2):  {},
							geom.Pt3(7, 0, -2):  {},
							geom.Pt3(8, 0, -2):  {},
							geom.Pt3(9, 0, -2):  {},
							geom.Pt3(10, 0, -2): {},
							geom.Pt3(11, 0, -2): {},
							geom.Pt3(12, 0, -2): {},

							geom.Pt3(1, 0, -1):  {},
							geom.Pt3(2, 0, -1):  {},
							geom.Pt3(3, 0, -1):  {},
							geom.Pt3(4, 0, -1):  {},
							geom.Pt3(5, 0, -1):  {},
							geom.Pt3(6, 0, -1):  {},
							geom.Pt3(7, 0, -1):  {},
							geom.Pt3(8, 0, -1):  {},
							geom.Pt3(9, 0, -1):  {},
							geom.Pt3(10, 0, -1): {},
							geom.Pt3(11, 0, -1): {},
							geom.Pt3(12, 0, -1): {},

							geom.Pt3(0, 0, 0):  {},
							geom.Pt3(1, 0, 0):  {},
							geom.Pt3(2, 0, 0):  {},
							geom.Pt3(3, 0, 0):  {},
							geom.Pt3(4, 0, 0):  {},
							geom.Pt3(5, 0, 0):  {},
							geom.Pt3(6, 0, 0):  {},
							geom.Pt3(7, 0, 0):  {},
							geom.Pt3(8, 0, 0):  {},
							geom.Pt3(9, 0, 0):  {},
							geom.Pt3(10, 0, 0): {},
							geom.Pt3(11, 0, 0): {},
							geom.Pt3(12, 0, 0): {},

							geom.Pt3(0, 0, 1):  {},
							geom.Pt3(1, 0, 1):  {},
							geom.Pt3(2, 0, 1):  {},
							geom.Pt3(3, 0, 1):  {},
							geom.Pt3(4, 0, 1):  {},
							geom.Pt3(5, 0, 1):  {},
							geom.Pt3(6, 0, 1):  {},
							geom.Pt3(7, 0, 1):  {},
							geom.Pt3(8, 0, 1):  {},
							geom.Pt3(9, 0, 1):  {},
							geom.Pt3(10, 0, 1): {},
							geom.Pt3(11, 0, 1): {},
							geom.Pt3(12, 0, 1): {},

							geom.Pt3(0, 0, 2):  {},
							geom.Pt3(1, 0, 2):  {},
							geom.Pt3(2, 0, 2):  {},
							geom.Pt3(3, 0, 2):  {},
							geom.Pt3(4, 0, 2):  {},
							geom.Pt3(5, 0, 2):  {},
							geom.Pt3(6, 0, 2):  {},
							geom.Pt3(7, 0, 2):  {},
							geom.Pt3(8, 0, 2):  {},
							geom.Pt3(9, 0, 2):  {},
							geom.Pt3(10, 0, 2): {},
							geom.Pt3(11, 0, 2): {},
							geom.Pt3(12, 0, 2): {},

							geom.Pt3(0, 0, 3):  {},
							geom.Pt3(1, 0, 3):  {},
							geom.Pt3(2, 0, 3):  {},
							geom.Pt3(3, 0, 3):  {},
							geom.Pt3(4, 0, 3):  {},
							geom.Pt3(5, 0, 3):  {},
							geom.Pt3(6, 0, 3):  {},
							geom.Pt3(7, 0, 3):  {},
							geom.Pt3(8, 0, 3):  {},
							geom.Pt3(9, 0, 3):  {},
							geom.Pt3(10, 0, 3): {},
							geom.Pt3(11, 0, 3): {},
							geom.Pt3(12, 0, 3): {},

							geom.Pt3(0, 0, 4):  {},
							geom.Pt3(1, 0, 4):  {},
							geom.Pt3(2, 0, 4):  {},
							geom.Pt3(3, 0, 4):  {},
							geom.Pt3(4, 0, 4):  {},
							geom.Pt3(5, 0, 4):  {},
							geom.Pt3(6, 0, 4):  {},
							geom.Pt3(7, 0, 4):  {},
							geom.Pt3(8, 0, 4):  {},
							geom.Pt3(9, 0, 4):  {},
							geom.Pt3(10, 0, 4): {},
							geom.Pt3(11, 0, 4): {},
							geom.Pt3(12, 0, 4): {},

							geom.Pt3(0, 0, 5):  {},
							geom.Pt3(1, 0, 5):  {},
							geom.Pt3(2, 0, 5):  {},
							geom.Pt3(3, 0, 5):  {},
							geom.Pt3(4, 0, 5):  {},
							geom.Pt3(5, 0, 5):  {},
							geom.Pt3(6, 0, 5):  {},
							geom.Pt3(6, -1, 5): {Cell: 1},
							geom.Pt3(7, 0, 5):  {},
							geom.Pt3(8, 0, 5):  {},
							geom.Pt3(9, 0, 5):  {},
							geom.Pt3(10, 0, 5): {},
							geom.Pt3(11, 0, 5): {},
							geom.Pt3(12, 0, 5): {},

							geom.Pt3(0, 0, 6):  {},
							geom.Pt3(1, 0, 6):  {},
							geom.Pt3(2, 0, 6):  {},
							geom.Pt3(3, 0, 6):  {},
							geom.Pt3(4, 0, 6):  {},
							geom.Pt3(5, 0, 6):  {},
							geom.Pt3(6, 0, 6):  {},
							geom.Pt3(7, 0, 6):  {},
							geom.Pt3(8, 0, 6):  {},
							geom.Pt3(9, 0, 6):  {},
							geom.Pt3(10, 0, 6): {},
							geom.Pt3(11, 0, 6): {},
							geom.Pt3(12, 0, 6): {},

							geom.Pt3(0, 0, 7):  {},
							geom.Pt3(1, 0, 7):  {},
							geom.Pt3(2, 0, 7):  {},
							geom.Pt3(3, 0, 7):  {},
							geom.Pt3(4, 0, 7):  {},
							geom.Pt3(5, 0, 7):  {},
							geom.Pt3(6, 0, 7):  {},
							geom.Pt3(7, 0, 7):  {},
							geom.Pt3(8, 0, 7):  {},
							geom.Pt3(9, 0, 7):  {},
							geom.Pt3(10, 0, 7): {},
							geom.Pt3(11, 0, 7): {},
							geom.Pt3(12, 0, 7): {},

							geom.Pt3(0, 0, 8):  {},
							geom.Pt3(1, 0, 8):  {},
							geom.Pt3(2, 0, 8):  {},
							geom.Pt3(3, 0, 8):  {},
							geom.Pt3(4, 0, 8):  {},
							geom.Pt3(5, 0, 8):  {},
							geom.Pt3(6, 0, 8):  {},
							geom.Pt3(7, 0, 8):  {},
							geom.Pt3(8, 0, 8):  {},
							geom.Pt3(9, 0, 8):  {},
							geom.Pt3(10, 0, 8): {},

							geom.Pt3(0, 0, 9): {},
							geom.Pt3(1, 0, 9): {},
							geom.Pt3(2, 0, 9): {},
							geom.Pt3(3, 0, 9): {},
							geom.Pt3(4, 0, 9): {},
							geom.Pt3(5, 0, 9): {},
							geom.Pt3(6, 0, 9): {},
							geom.Pt3(7, 0, 9): {},
							geom.Pt3(8, 0, 9): {},

							geom.Pt3(0, 0, 10): {},
							geom.Pt3(1, 0, 10): {},
							geom.Pt3(2, 0, 10): {},
							geom.Pt3(3, 0, 10): {},
							geom.Pt3(4, 0, 10): {},
							geom.Pt3(5, 0, 10): {},
							geom.Pt3(6, 0, 10): {},

							geom.Pt3(0, 0, 11): {},
							geom.Pt3(1, 0, 11): {},
							geom.Pt3(2, 0, 11): {},
							geom.Pt3(3, 0, 11): {},
							geom.Pt3(4, 0, 11): {},

							geom.Pt3(0, 0, 12): {},
							geom.Pt3(1, 0, 12): {},
							geom.Pt3(2, 0, 12): {},

							geom.Pt3(0, 0, 13): {},
						}, // Map
					}, // PrismMap
					&Awakeman{
						CameraID: "game_camera",
						ToastID:  "toast",
						Sprite: engine.Sprite{
							Actor: engine.Actor{
								CollisionDomain: "level_1",
								Pos:             geom.Pt3(100, -64, 100),
								Bounds: geom.Box{
									Min: geom.Pt3(-4, -15, -1),
									Max: geom.Pt3(4, 1, 1),
								},
							},
							DrawOffset: image.Pt(-5, -15),
							Sheet: engine.Sheet{
								AnimDefs: map[string]*engine.AnimDef{
									"idle_left": {Steps: []engine.AnimStep{
										{Cell: 1, Duration: 60},
									}},
									"idle_right": {Steps: []engine.AnimStep{
										{Cell: 0, Duration: 60},
									}},
									"run_left": {Steps: []engine.AnimStep{
										{Cell: 14, Duration: 3},
										{Cell: 15, Duration: 5},
										{Cell: 16, Duration: 3},
										{Cell: 17, Duration: 3},
									}},
									"run_right": {Steps: []engine.AnimStep{
										{Cell: 10, Duration: 3},
										{Cell: 11, Duration: 5},
										{Cell: 12, Duration: 3},
										{Cell: 13, Duration: 3},
									}},
									"run_vert": {Steps: []engine.AnimStep{
										{Cell: 18, Duration: 3},
										{Cell: 19, Duration: 5},
										{Cell: 20, Duration: 3},
										{Cell: 21, Duration: 3},
										{Cell: 22, Duration: 3},
										{Cell: 23, Duration: 5},
										{Cell: 24, Duration: 3},
										{Cell: 25, Duration: 3},
									}},
									"walk_left": {Steps: []engine.AnimStep{
										{Cell: 2, Duration: 6},
										{Cell: 3, Duration: 6},
										{Cell: 4, Duration: 6},
										{Cell: 5, Duration: 6},
									}},
									"walk_right": {Steps: []engine.AnimStep{
										{Cell: 6, Duration: 6},
										{Cell: 7, Duration: 6},
										{Cell: 8, Duration: 6},
										{Cell: 9, Duration: 6},
									}},
								},
								CellSize: image.Pt(10, 16),
								Src:      engine.ImageRef{Path: "assets/aw.png"},
							}, // Sheet
						}, // Sprite
					}, // Awakeman
				), // Container
			}, // DrawDAG
		), // Container
	} // Scene
}
