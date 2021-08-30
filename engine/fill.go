package engine

import (
	"encoding/gob"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Fill satisfies Drawer.
var _ interface {
	Drawer
	Hider
	Identifier
} = &Fill{}

func init() {
	gob.Register(&Fill{})
	gob.Register(color.Gray{})
	gob.Register(color.RGBA{})
}

// Fill fills the screen with a colour.
type Fill struct {
	ID
	Color color.Color
	Hidden
	ZOrder
}

func (f *Fill) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if f.Hidden {
		return
	}
	screen.Fill(opts.ColorM.Apply(f.Color))
}
