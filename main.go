package main

import (
	"embed"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"drjosh.dev/gurgle/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

const screenWidth, screenHeight = 320, 240

//go:embed assets
var assets embed.FS

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("TODO")

	engine.AssetFS = assets
	engine.AnimDefs = map[string]*engine.AnimDef{
		"green_tiles": {
			Frames: []engine.AnimFrame{
				{Frame: 0, Duration: 16},
				{Frame: 1, Duration: 16},
				{Frame: 2, Duration: 16},
			},
		},
		"red_tiles": {
			Frames: []engine.AnimFrame{
				{Frame: 3, Duration: 12},
				{Frame: 4, Duration: 12},
				{Frame: 5, Duration: 12},
				{Frame: 6, Duration: 12},
			},
		},
		"aw_idle": {
			Frames: []engine.AnimFrame{
				{Frame: 0, Duration: 60},
			},
		},
		"aw_walk_right": {
			Frames: []engine.AnimFrame{
				{Frame: 1, Duration: 6},
				{Frame: 2, Duration: 6},
				{Frame: 3, Duration: 6},
				{Frame: 4, Duration: 6},
			},
		},
		"aw_walk_left": {
			Frames: []engine.AnimFrame{
				{Frame: 5, Duration: 6},
				{Frame: 6, Duration: 6},
				{Frame: 7, Duration: 6},
				{Frame: 8, Duration: 6},
			},
		},
		"aw_run_right": {
			Frames: []engine.AnimFrame{
				{Frame: 9, Duration: 3},
				{Frame: 10, Duration: 3},
				{Frame: 11, Duration: 3},
				{Frame: 12, Duration: 3},
			},
		},
		"aw_run_left": {
			Frames: []engine.AnimFrame{
				{Frame: 13, Duration: 3},
				{Frame: 14, Duration: 3},
				{Frame: 15, Duration: 3},
				{Frame: 16, Duration: 3},
			},
		},
	}

	tiles := [][]engine.Tile{
		{engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, engine.StaticTile(9)},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8)},
		{engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7)},
		{engine.StaticTile(9), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(9)},
	}

	level1 := &engine.Scene{
		ID: "level_1",
		Components: []interface{}{
			engine.Fill{
				Color: color.White,
			},
			&engine.Tilemap{
				ID:       "terrain",
				Map:      tiles,
				Src:      engine.ImageRef{Path: "assets/boxes.png"},
				TileSize: 16,
				ZPos:     0,
			},
			&engine.SolidRect{
				ID:   "ceiling",
				Rect: image.Rect(0, -1, 320, 0),
			},
			&engine.SolidRect{
				ID:   "left_wall",
				Rect: image.Rect(-1, 0, 0, 240),
			},
			&engine.SolidRect{
				ID:   "right_wall",
				Rect: image.Rect(320, 0, 321, 240),
			},
			/*&engine.SolidRect{
				ID:   "ground",
				Rect: image.Rect(0, 192, 320, 240),
			},*/
			&engine.Sprite{
				ID: "protagonist",
				Actor: engine.Actor{
					Position: image.Pt(100, 100),
					Size:     image.Pt(10, 16),
				},
				Src:  engine.ImageRef{Path: "assets/aw.png"},
				ZPos: 1,
			},
		},
	}

	game := &engine.Game{
		ScreenHeight: screenHeight,
		ScreenWidth:  screenWidth,
		Scene: &engine.Scene{
			ID: "root",
			Components: []interface{}{
				&engine.GobDumper{
					KeyCombo: []ebiten.Key{ebiten.KeyControl, ebiten.KeyD},
				},
				level1,
				engine.PerfDisplay{},
			},
		},
	}
	game.Build()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
