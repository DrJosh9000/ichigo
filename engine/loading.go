package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type LoadingScreen struct{}

func (LoadingScreen) Draw(screen *ebiten.Image, _ *ebiten.DrawImageOptions) {
	ebitenutil.DebugPrint(screen, "Loading...")
}
