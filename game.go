package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const screenWidth, screenHeight = 320, 240

type game struct{}

func (g *game) Update() error {
	// TODO
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 50, 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

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
