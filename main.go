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

	staticTiles := [][]engine.StaticTile{
		{0, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 1},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		{2, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 3},
	}

	tiles := make([][]engine.Tile, len(staticTiles))
	for j, row := range staticTiles {
		tiles[j] = make([]engine.Tile, len(row))
		for i, t := range row {
			tiles[j][i] = t
		}
	}

	tiles[4][5] = &engine.AnimatedTile{
		AnimRef: engine.AnimRef{Key: "green_tiles"},
	}
	tiles[6][7] = &engine.AnimatedTile{
		AnimRef: engine.AnimRef{Key: "red_tiles"},
	}

	game := &engine.Game{
		ScreenHeight: screenHeight,
		ScreenWidth:  screenWidth,
		Scene: &engine.Scene{
			Components: []interface{}{
				&engine.Tilemap{
					Map:      tiles,
					Src:      engine.ImageRef{Path: "assets/boxes.png"},
					TileSize: 16,
				},
				engine.PerfDisplay{},
			},
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
