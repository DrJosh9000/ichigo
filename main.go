package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"drjosh.dev/gurgle/engine"
	"drjosh.dev/gurgle/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	enableCPUProfile = false
	hardcodedLevel1  = true
	rewriteLevel1    = false
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
		lev1 = level1()
		if rewriteLevel1 && runtime.GOOS != "js" {
			if err := engine.SaveGobz(lev1, "game/assets/level1.gobz"); err != nil {
				log.Fatalf("Couldn't save level1.gobz: %v", err)
			}
		}
	}

	g := &engine.Game{
		ScreenSize: image.Pt(320, 240),
		Root: &engine.Scene{
			ID: "root",
			Components: []interface{}{
				&engine.Camera{
					ID:    "game_camera",
					Child: lev1,
				},
				&engine.DebugToast{ID: "toast", Pos: image.Pt(0, 15)},
				engine.PerfDisplay{},
			},
		},
	}
	if err := g.LoadAndPrepare(game.Assets); err != nil {
		log.Fatalf("Loading/preparing error: %v", err)
	}

	if runtime.GOOS != "js" {
		// Run a repl on the console.
		go g.REPL(os.Stdin, os.Stdout, game.Assets)
	}

	// ... while the game also runs
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}

func level1() *engine.Scene {
	denseTiles := [][]engine.Tile{
		{engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, engine.StaticTile(9)},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimKey: "red_tile"}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimKey: "red_tile"}, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, &engine.AnimatedTile{AnimKey: "green_tile"}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimKey: "red_tile"}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimKey: "green_tile"}, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, &engine.AnimatedTile{AnimKey: "green_tile"}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimKey: "green_tile"}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, &engine.AnimatedTile{AnimKey: "red_tile"}, nil, nil, nil, nil, nil, nil, engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, engine.StaticTile(9), nil, nil, nil},
		{nil, nil, nil, nil, engine.StaticTile(9), engine.StaticTile(9), engine.StaticTile(9), nil, nil, nil, engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), &engine.AnimatedTile{AnimKey: "red_tile"}, engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8)},
		{engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), &engine.AnimatedTile{AnimKey: "red_tile"}, engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), &engine.AnimatedTile{AnimKey: "green_tile"}, engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7)},
		{engine.StaticTile(9), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(9)},
	}
	tiles := make(map[image.Point]engine.Tile)
	for j, row := range denseTiles {
		for i, tile := range row {
			if tile == nil {
				continue
			}
			tiles[image.Pt(i, j)] = tile
		}
	}

	return &engine.Scene{
		ID:     "level_1",
		Bounds: engine.Bounds(image.Rect(-32, -32, 320+32, 240+32)),
		Components: []interface{}{
			&engine.Fill{
				ID:     "bg_fill",
				Color:  color.Gray{100},
				ZOrder: -1000,
			},
			&engine.Parallax{
				CameraID: "game_camera",
				Child: &engine.Billboard{
					ID:     "bg_image",
					ZOrder: -900,
					Pos:    image.Pt(-160, -120),
					Src:    engine.ImageRef{Path: "assets/space.png"},
				},
				Factor: 0.5,
			},
			&engine.IsoVoxmap{
				ID:            "voxmap",
				DrawOrderBias: image.Pt(1, -1), // left before right, bottom before top
				DrawOffset:    image.Pt(-8, 0),
				Projection:    image.Pt(-2, 2), // each step in Z becomes -1/2 step in X plus 1/2 step in Y.
				VoxSize:       engine.Pt3(16, 16, 16),
				Sheet: engine.Sheet{
					CellSize: image.Pt(24, 24),
					Src:      engine.ImageRef{Path: "assets/box16.png"},
				},
				Map: map[engine.Point3]*engine.IsoVoxel{
					engine.Pt3(0, 2, 0): {CellBack: 0, CellFront: 1},
					engine.Pt3(1, 2, 0): {CellBack: 0, CellFront: 1},
					engine.Pt3(2, 2, 0): {CellBack: 0, CellFront: 1},

					engine.Pt3(5, 0, 0): {CellBack: 0, CellFront: 1},
					engine.Pt3(5, 1, 0): {CellBack: 0, CellFront: 1},
					engine.Pt3(5, 2, 0): {CellBack: 0, CellFront: 1},

					engine.Pt3(8, 2, 0): {CellBack: 0, CellFront: 1},
					engine.Pt3(8, 2, 1): {CellBack: 0, CellFront: 1},
					engine.Pt3(8, 2, 2): {CellBack: 0, CellFront: 1},

					engine.Pt3(11, 2, 0):  {CellBack: 0, CellFront: 1},
					engine.Pt3(11, 2, -1): {CellBack: 0, CellFront: 1},
					engine.Pt3(11, 2, -2): {CellBack: 0, CellFront: 1},
				},
			},
			&engine.Tilemap{
				ID:     "terrain",
				ZOrder: -1,
				Map:    tiles,
				Sheet: engine.Sheet{
					AnimDefs: map[string]*engine.AnimDef{
						"red_tile": {
							Steps: []engine.AnimStep{
								{Cell: 3, Duration: 12},
								{Cell: 4, Duration: 12},
								{Cell: 5, Duration: 12},
								{Cell: 6, Duration: 12},
							}},
						"green_tile": {
							Steps: []engine.AnimStep{
								{Cell: 0, Duration: 16},
								{Cell: 1, Duration: 16},
								{Cell: 2, Duration: 16},
							},
						},
					},
					CellSize: image.Pt(16, 16),
					Src:      engine.ImageRef{Path: "assets/boxes.png"},
				},
			},
			&engine.SolidRect{
				ID:  "ceiling",
				Box: engine.Box{Min: engine.Pt3(0, -1, 0), Max: engine.Pt3(320, 0, 100)},
			},
			&engine.SolidRect{
				ID:  "left_wall",
				Box: engine.Box{Min: engine.Pt3(-1, 0, 0), Max: engine.Pt3(0, 240, 100)},
			},
			&engine.SolidRect{
				ID:  "right_wall",
				Box: engine.Box{Min: engine.Pt3(320, 0, 0), Max: engine.Pt3(321, 240, 100)},
			},
			&game.Awakeman{
				CameraID: "game_camera",
				ToastID:  "toast",
				Sprite: engine.Sprite{
					Actor: engine.Actor{
						CollisionDomain: "level_1",
						Pos:             engine.Pt3(100, 100, 9),
						Size:            engine.Pt3(8, 16, 8),
					},
					FrameOffset: image.Pt(-1, 0),
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
					},
				},
			},
		},
	}
}
