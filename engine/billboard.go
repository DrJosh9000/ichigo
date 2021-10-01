package engine

import (
	"encoding/gob"
	"fmt"

	"github.com/DrJosh9000/ichigo/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Billboard satisfies interfaces.
var _ interface {
	BoundingBoxer
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
	Hides
	Pos geom.Int3
	Src ImageRef

	game *Game
}

// BoundingBox returns a 0-depth box incorporating the image size.
func (b *Billboard) BoundingBox() geom.Box {
	sx, sy := b.Src.Image().Size()
	return geom.Box{
		Min: b.Pos,
		Max: b.Pos.Add(geom.Pt3(sx, sy, 0)),
	}
}

// Draw draws the image.
func (b *Billboard) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(b.Src.Image(), opts)
}

// Prepare saves the reference to Game.
func (b *Billboard) Prepare(g *Game) error {
	b.game = g
	return nil
}

// Scan visits &b.Src.
func (b *Billboard) Scan(visit VisitFunc) error {
	return visit(&b.Src)
}

// String returns "Billboard@(b.Pos)".
func (b *Billboard) String() string {
	return fmt.Sprintf("Billboard@%v", b.Pos)
}

// Transform returns a translation by the projected position.
func (b *Billboard) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(geom.CFloat(
		geom.Project(b.game.Projection, b.Pos),
	))
	return opts
}
