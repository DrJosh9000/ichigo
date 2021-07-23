package engine

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TPSDisplay debugprints CurrentTPS in the top left.
type TPSDisplay struct{}

func (TPSDisplay) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (TPSDisplay) Z() float64 {
	// Always draw on top
	return math.MaxFloat64
}
