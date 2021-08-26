package main

import (
	_ "image/png"
	"log"
	"os"
	"runtime"

	"drjosh.dev/gurgle/engine"
	"drjosh.dev/gurgle/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("TODO")

	g := &engine.Game{
		ScreenHeight: 240,
		ScreenWidth:  320,
		Root: &engine.Scene{
			ID: "root",
			Components: []interface{}{
				&engine.Camera{
					ID:    "game_camera",
					Scene: &engine.SceneRef{Path: "assets/level1.gobz"},
				},
				&engine.DebugToast{ID: "toast"},
				engine.PerfDisplay{},
			},
		},
	}
	if err := g.Load(game.Assets); err != nil {
		log.Fatalf("Loading error: %v", err)
	}
	g.Prepare()

	if runtime.GOOS != "js" {
		// Run a repl on the console.
		go g.REPL(os.Stdin, os.Stdout, game.Assets)
	}

	// ... while the game also runs
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
