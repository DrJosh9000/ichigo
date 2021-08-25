package main

import (
	"bufio"
	_ "image/png"
	"log"
	"os"
	"strings"

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

	// Run a repl on the console.
	go repl(g)

	// ... while the game also runs
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}

func repl(g *engine.Game) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		tok := strings.Split(sc.Text(), " ")
		if len(tok) == 0 {
			continue
		}
		switch tok[0] {
		case "quit":
			os.Exit(0)
		case "pause":
			g.Disable()
		case "resume", "unpause":
			g.Enable()
		case "save":
			if len(tok) != 2 {
				log.Print("Usage: save ID")
				break
			}
			id := tok[1]
			c := g.Component(id)
			if c == nil {
				log.Printf("Component %q not found", id)
				break
			}
			s, ok := c.(engine.Saver)
			if !ok {
				log.Printf("Component %q not a Saver (type %T)", id, c)
				break
			}
			if err := s.Save(); err != nil {
				log.Printf("Couldn't save %q: %v", id, err)
			}
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("Couldn't scan stdin: %v", err)
	}
}
