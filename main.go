package main

import (
	"embed"
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

	game := &engine.Game{
		ScreenHeight: screenHeight,
		ScreenWidth:  screenWidth,
		Scene: &engine.Scene{
			ID: "root",
			Components: []interface{}{
				&engine.GobDumper{
					KeyCombo: []ebiten.Key{ebiten.KeyControl, ebiten.KeyD},
				},
				&engine.Tilemap{
					ID:       "terrain",
					Map:      tiles,
					Src:      engine.ImageRef{Path: "assets/boxes.png"},
					TileSize: 16,
				},
				engine.PerfDisplay{},
			},
		},
	}
	game.Build()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
