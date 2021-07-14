package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const screenWidth, screenHeight = 320, 240

type game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	// TODO
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	// TODO
	screen.Fill(color.White)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("ebiten")
	if err := ebiten.RunGame(&game{}); err != nil {
		log.Fatal(err)
	}
}
