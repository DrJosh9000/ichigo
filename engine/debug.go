package engine

import (
	"compress/gzip"
	"encoding/gob"
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func init() {
	gob.Register(GobDumper{})
	gob.Register(PerfDisplay{})
}

// PerfDisplay debugprints CurrentTPS and CurrentFPS in the top left.
type PerfDisplay struct {
	Hidden bool
}

func (p PerfDisplay) Draw(screen *ebiten.Image, _ ebiten.GeoM) {
	if p.Hidden {
		return
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f  FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (PerfDisplay) Z() float64 {
	// Always draw on top
	return math.MaxFloat64
}

// GobDumper waits for a given key combo, then dumps the game into a gob file
// in the current directory.
type GobDumper struct {
	KeyCombo []ebiten.Key

	game *Game
}

// Prepare simply stores the reference to the Game.
func (d *GobDumper) Prepare(g *Game) { d.game = g }

// Update waits for the key combo, then dumps the game state into a gzipped gob.
func (d *GobDumper) Update() error {
	for _, key := range d.KeyCombo {
		if !ebiten.IsKeyPressed(key) {
			return nil
		}
	}
	if d.game == nil {
		return errors.New("nil d.game in GobDumper.Update")
	}
	f, err := os.Create(time.Now().Format("20060102030405.gob.gz"))
	if err != nil {
		return err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	defer gz.Close()
	if err := gob.NewEncoder(gz).Encode(d.game); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	return f.Close()
}
