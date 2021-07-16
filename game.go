package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const screenWidth, screenHeight = 320, 240

type updater interface {
	Update() error
}

type drawer interface {
	Draw(*ebiten.Image)
}

type game struct {
	components []interface{}
}

func (g *game) Update() error {
	for _, c := range g.components {
		if u, ok := c.(updater); ok {
			if err := u.Update(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	for _, c := range g.components {
		if d, ok := c.(drawer); ok {
			d.Draw(screen)
		}
	}
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
