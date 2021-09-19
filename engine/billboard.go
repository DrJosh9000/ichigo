package engine

import (
	"encoding/gob"
	"fmt"

	"drjosh.dev/gurgle/geom"
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

// DrawAfter reports if b.Pos.Z >= x.Max.Z.
func (b *Billboard) DrawAfter(x Drawer) bool {
	switch x := x.(type) {
	case BoundingBoxer:
		return b.Pos.Z >= x.BoundingBox().Max.Z
	case ZPositioner:
		return b.Pos.Z > x.ZPos()
	}
	return false
}

// DrawBefore reports if b.Pos.Z < x.Min.Z.
func (b *Billboard) DrawBefore(x Drawer) bool {
	switch x := x.(type) {
	case BoundingBoxer:
		return b.Pos.Z < x.BoundingBox().Min.Z
	case ZPositioner:
		return b.Pos.Z < x.ZPos()
	}
	return false
}

// Prepare saves the reference to Game.
func (b *Billboard) Prepare(g *Game) error {
	b.game = g
	return nil
}

// Scan returns a slice containing Src.
func (b *Billboard) Scan() []interface{} { return []interface{}{&b.Src} }

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
