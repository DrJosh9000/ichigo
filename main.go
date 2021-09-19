package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"

	"drjosh.dev/gurgle/engine"
	"drjosh.dev/gurgle/game"
	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	enableCPUProfile = true
	enableREPL       = true
	hardcodedLevel1  = true
	rewriteLevel1    = false
)

func main() {
	// Change to true to enable cpu profile
	if enableCPUProfile && runtime.GOOS != "js" {
		f, err := os.Create("cpuprofile.pprof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("TODO")

	// Change to true to rewrite level1.gobz
	lev1 := interface{}(&engine.SceneRef{Path: "assets/level1.gobz"})
	if hardcodedLevel1 {
		lev1 = game.Level1()
		if rewriteLevel1 && runtime.GOOS != "js" {
			if err := engine.SaveGobz(lev1, "game/assets/level1.gobz"); err != nil {
				log.Fatalf("Couldn't save level1.gobz: %v", err)
			}
		}
	}

	g := &engine.Game{
		ScreenSize: image.Pt(320, 240), // Window interior is this many pixels.
		// TODO: refactor Projection and VoxelScale into... Scene? Camera?
		// We might want different projections and scales in different levels.
		Projection: geom.IntProjection{
			// Each 1 voxel step in Z is projected into 1 pixel in Y.
			X: 0,
			Y: 1,
		},
		VoxelScale: geom.Float3{
			// Each voxel counts for this much Eucliden space.
			X: 1,
			Y: 1,
			Z: math.Sqrt(3),
		},
		Root: &engine.Scene{
			ID: "root",
			Components: []interface{}{
				&engine.Camera{
					ID:    "game_camera",
					Child: lev1,
				},
				&engine.DebugToast{ID: "toast", Pos: image.Pt(0, 15)},
				engine.PerfDisplay{},
			},
		},
	}
	if err := g.LoadAndPrepare(game.Assets); err != nil {
		log.Fatalf("Loading/preparing error: %v", err)
	}

	if enableREPL && runtime.GOOS != "js" {
		go g.REPL(os.Stdin, os.Stdout, game.Assets)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}
