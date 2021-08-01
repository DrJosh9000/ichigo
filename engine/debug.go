package engine

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func init() {
	gob.Register(PerfDisplay{})
}

// PerfDisplay debugprints CurrentTPS and CurrentFPS in the top left.
type PerfDisplay struct{}

func (PerfDisplay) Draw(screen *ebiten.Image, _ ebiten.GeoM) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f  FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (PerfDisplay) Z() float64 {
	// Always draw on top
	return math.MaxFloat64
}

// GobDumper waits for a given key combo, then dumps the game into a gob file
// in the current directory.
type GobDumper struct {
	combo []ebiten.Key
	game  ebiten.Game
}

func (d GobDumper) Update() error {
	for _, key := range d.combo {
		if !ebiten.IsKeyPressed(key) {
			return nil
		}
	}
	f, err := os.Create(time.Now().Format("20060102030405.gob"))
	if err != nil {
		return err
	}
	return gob.NewEncoder(f).Encode(d.game)
}
