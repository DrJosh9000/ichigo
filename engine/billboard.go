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
	Scanner
	Transformer
} = &Billboard{}

func init() {
	gob.Register(&Billboard{})
}

// Billboard draws an image at a position.
type Billboard struct {
	ID
	Hidden
	Pos image.Point
	Src ImageRef
	ZOrder
}

// Draw draws the image.
func (b *Billboard) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(b.Src.Image(), opts)
}

// Scan returns a slice containing Src.
func (b *Billboard) Scan() []interface{} { return []interface{}{&b.Src} }

func (b *Billboard) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(cfloat(b.Pos))
	return opts
}
