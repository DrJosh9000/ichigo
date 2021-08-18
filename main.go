package main

import (
	"embed"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"drjosh.dev/gurgle/engine"
	"drjosh.dev/gurgle/game"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var assets embed.FS

func main() {
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("TODO")

	engine.AssetFS = assets
	// engine.AnimDefs set in game/anims.go

	denseTiles := [][]engine.Tile{
		{engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, engine.StaticTile(9)},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, nil, nil, nil, nil, nil, nil, engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, engine.StaticTile(9), nil, nil, nil},
		{nil, nil, nil, nil, engine.StaticTile(9), engine.StaticTile(9), engine.StaticTile(9), nil, nil, nil, engine.StaticTile(9), nil, nil, nil, nil, nil, nil, nil, nil, nil},
		{engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8), engine.StaticTile(8)},
		{engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "red_tiles"}}, engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), &engine.AnimatedTile{AnimRef: engine.AnimRef{Key: "green_tiles"}}, engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7)},
		{engine.StaticTile(9), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(7), engine.StaticTile(9)},
	}
	tiles := make(map[image.Point]engine.Tile)
	for j, row := range denseTiles {
		for i, tile := range row {
			tiles[image.Pt(i, j)] = tile
		}
	}

	level1 := &engine.Scene{
		ID: "level_1",
		Components: []interface{}{
			&engine.Fill{
				Color: color.Gray{100},
				DrawOrder:  0,
			},
			&engine.Tilemap{
				ID:       "terrain",
				Map:      tiles,
				Src:      engine.ImageRef{Path: "assets/boxes.png"},
				TileSize: 16,
				DrawOrder:     1,
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
			&game.Awakeman{
				CameraID: "level_1_camera",
				Sprite: engine.Sprite{
					ID: "awakeman",
					Actor: engine.Actor{
						CollisionDomain: "level_1",
						Pos:             image.Pt(100, 100),
						Size:            image.Pt(8, 16),
					},
					FrameOffset: image.Pt(-1, 0),
					FrameSize:   image.Pt(10, 16),
					Src:         engine.ImageRef{Path: "assets/aw.png"},
					DrawOrder:        2,
				},
			},
		},
	}

	game := &engine.Game{
		ScreenHeight: 240,
		ScreenWidth:  320,
		Scene: &engine.Scene{
			ID: "root",
			Components: []interface{}{
				&engine.GobDumper{
					KeyCombo: []ebiten.Key{ebiten.KeyControl, ebiten.KeyD},
				},
				&engine.Camera{
					ID:     "level_1_camera",
					Bounds: image.Rect(-32, -32, 320+32, 240+32),
					Scene:  level1,
				},
				engine.PerfDisplay{},
			},
		},
	}
	game.PrepareToRun()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
