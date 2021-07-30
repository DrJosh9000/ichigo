package main

import (
	"embed"
	"image"
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

	boxesFile, err := assets.Open("assets/boxes.png")
	if err != nil {
		log.Fatalf("Couldn't open asset: %v", err)
	}
	boxesPNG, _, err := image.Decode(boxesFile)
	if err != nil {
		log.Fatalf("Couldn't decode asset: %v", err)
	}
	boxesFile.Close()

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
	for j := range staticTiles {
		tiles[j] = make([]engine.Tile, len(staticTiles[j]))
		for i := range staticTiles[j] {
			tiles[j][i] = staticTiles[j][i]
		}
	}

	tiles[3][5] = &engine.AnimatedTile{
		AnimDef: []engine.TileAnimFrameDef{
			{Tile: 0, DurationTicks: 16},
			{Tile: 1, DurationTicks: 16},
			{Tile: 2, DurationTicks: 16},
		},
	}
	tiles[2][7] = &engine.AnimatedTile{
		AnimDef: []engine.TileAnimFrameDef{
			{Tile: 3, DurationTicks: 12},
			{Tile: 4, DurationTicks: 12},
			{Tile: 5, DurationTicks: 12},
			{Tile: 6, DurationTicks: 12},
		},
	}

	tilemap := &engine.Tilemap{
		Map:      tiles,
		Src:      ebiten.NewImageFromImage(boxesPNG),
		TileSize: 16,
	}

	game := &engine.Game{
		ScreenHeight: screenHeight,
		ScreenWidth:  screenWidth,
		Scene: &engine.Scene{
			Components: []interface{}{
				tilemap,
				engine.PerfDisplay{},
			},
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
