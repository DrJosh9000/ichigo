package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Billboard satisfies interfaces.
var _ interface {
	Identifier
	Drawer
	DrawOrderer
	ParallaxScaler
	Scanner
} = &Billboard{}

func init() {
	gob.Register(&Billboard{})
}

// Billboard draws an image at a position.
type Billboard struct {
	ID
	Hidden
	Parallax
	Pos image.Point
	Src ImageRef
	ZOrder
}

// Draw draws the image.
func (b *Billboard) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if b.Hidden {
		return
	}
	var geom ebiten.GeoM
	geom.Translate(float64(b.Pos.X), float64(b.Pos.Y))
	geom.Concat(opts.GeoM)
	opts.GeoM = geom
	screen.DrawImage(b.Src.Image(), &opts)
}

// Scan returns a slice containing Src.
func (b *Billboard) Scan() []interface{} { return []interface{}{&b.Src} }
