package engine

import (
	"encoding/gob"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var _ interface {
	Drawer
	DrawOrderer
	Hider
} = &PerfDisplay{}

func init() {
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
