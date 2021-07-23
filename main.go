package main

import (
	"log"

	"drjosh.dev/gurgle/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

const screenWidth, screenHeight = 320, 240

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("ebiten")
	if err := ebiten.RunGame(&engine.Game{
		ScreenHeight: screenHeight,
		ScreenWidth:  screenWidth,
		Components: []interface{}{
			engine.TPSDisplay{},
		},
	}); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
