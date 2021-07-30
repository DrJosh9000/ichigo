package engine

import (
	"encoding/gob"
	"fmt"
	"math"

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
