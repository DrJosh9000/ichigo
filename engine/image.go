package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Image satisfies interfaces.
var (
	_ Identifier     = &Image{}
	_ Drawer         = &Image{}
	_ DrawOrderer    = &Image{}
	_ ParallaxScaler = &Image{}
)

func init() {
	gob.Register(Image{})
}

// Image draws an image at a position.
type Image struct {
	ID
	Parallax
	ZOrder
	Src ImageRef
	Pos image.Point
}

// Draw draws the image.
func (i *Image) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	var geom ebiten.GeoM
	geom.Translate(float64(i.Pos.X), float64(i.Pos.Y))
	geom.Concat(opts.GeoM)
	opts.GeoM = geom
	screen.DrawImage(i.Src.Image(), &opts)
}
