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
	Updater
} = &Sprite{}

func init() {
	gob.Register(&Sprite{})
}

// Sprite combines an Actor with the ability to Draw from a single spritesheet.
type Sprite struct {
	Actor
	FrameOffset image.Point
	Hidden
	Sheet Sheet
	ZOrder

	anim *Anim
}

func (s *Sprite) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if s.Hidden {
		return
	}
	dp := s.Pos.Add(s.FrameOffset)
	var geom ebiten.GeoM
	geom.Translate(float64(dp.X), float64(dp.Y))
	geom.Concat(opts.GeoM)
	opts.GeoM = geom

	src := s.Sheet.SubImage(s.anim.CurrentFrame())
	screen.DrawImage(src, &opts)
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

func (s *Sprite) Update() error { return s.anim.Update() }
