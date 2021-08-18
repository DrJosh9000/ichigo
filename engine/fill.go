package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Fill satisfies Drawer.
var _ Drawer = &Fill{}

// Fill fills the screen with a colour.
type Fill struct {
	Color color.Color
	ZOrder
	Hidden bool
	ID
}

func (f *Fill) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if f.Hidden {
		return
	}
	screen.Fill(opts.ColorM.Apply(f.Color))
}
