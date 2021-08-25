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

var (
	_ Drawer      = PerfDisplay{}
	_ DrawOrderer = PerfDisplay{}
	_ Hider       = &PerfDisplay{}

	_ Disabler = &GobDumper{}
	_ Prepper  = &GobDumper{}
	_ Updater  = &GobDumper{}
)

func init() {
	gob.Register(&GobDumper{})
	gob.Register(&PerfDisplay{})
}

// DebugToast debugprints a string for a while, then disappears.
type DebugToast struct {
	ID
	Hidden
	Timer int // ticks
	Text  string
}

func (d *DebugToast) Draw(screen *ebiten.Image, _ ebiten.DrawImageOptions) {
	if d.Hidden {
		return
	}
	ebitenutil.DebugPrintAt(screen, d.Text, 0, 20)
}

func (d *DebugToast) DrawOrder() float64 {
	return math.MaxFloat64
}

func (d *DebugToast) Toast(text string) {
	d.Text = text
	d.Timer = 120
	d.Hidden = false
}

func (d *DebugToast) Update() error {
	if d.Hidden = d.Timer <= 0; !d.Hidden {
		d.Timer--
	}
	return nil
}

// PerfDisplay debugprints CurrentTPS and CurrentFPS in the top left.
type PerfDisplay struct {
	Hidden
}

func (p PerfDisplay) Draw(screen *ebiten.Image, _ ebiten.DrawImageOptions) {
	if p.Hidden {
		return
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f  FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (PerfDisplay) DrawOrder() float64 {
	// Always draw on top
	return math.MaxFloat64
}

// GobDumper waits for a given key combo, then dumps the game into a gob file
// in the current directory.
type GobDumper struct {
	Disabled
	KeyCombo []ebiten.Key

	game *Game
}

// Prepare simply stores the reference to the Game.
func (d *GobDumper) Prepare(g *Game) { d.game = g }

// Update waits for the key combo, then dumps the game state into a gzipped gob.
func (d *GobDumper) Update() error {
	if d.Disabled {
		return nil
	}
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
