package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Sprite satisfies interfaces.
var _ interface {
	Drawer
	Scanner
	Transformer
	Updater
} = &Sprite{}

func init() {
	gob.Register(&Sprite{})
}

// Sprite combines an Actor with the ability to Draw from a single spritesheet.
type Sprite struct {
	Actor       Actor
	FrameOffset image.Point
	Hidden
	Sheet Sheet
	ZOrder

	anim *Anim
}

// Draw draws the current cell to the screen.
func (s *Sprite) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(s.Sheet.SubImage(s.anim.Cell()), opts)
}

// Scan returns the Actor and the Sheet.
func (s *Sprite) Scan() []interface{} {
	return []interface{}{
		&s.Actor,
		&s.Sheet,
	}
}

// SetAnim sets the Anim to use for the sprite. If it is not the same as the
// one currently set, it resets the new anim.
func (s *Sprite) SetAnim(a *Anim) {
	if s.anim != a {
		a.Reset()
	}
	s.anim = a
}

// Transform returns a translation by the FrameOffset.
func (s *Sprite) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(cfloat(s.Actor.Pos.Add(s.FrameOffset)))
	return opts
}

// Update updates the Sprite's anim. anim can change a bit so we don't tell Game
// about it, but that means it must be updated manually.
func (s *Sprite) Update() error { return s.anim.Update() }
