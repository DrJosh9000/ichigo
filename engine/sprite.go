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

func (s *Sprite) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(s.Sheet.SubImage(s.anim.CurrentFrame()), opts)
}

func (s *Sprite) Scan() []interface{} {
	return []interface{}{
		&s.Actor,
		&s.Sheet,
	}
}

func (s *Sprite) SetAnim(a *Anim) {
	if s.anim != a {
		a.Reset()
	}
	s.anim = a
}

func (s *Sprite) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(float2(s.Actor.Pos.Add(s.FrameOffset)))
	return opts
}

// anim can change a bit so we don't tell Game about it, but that means it must
// be updated here.
func (s *Sprite) Update() error { return s.anim.Update() }
