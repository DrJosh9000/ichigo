package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(Sprite{})
}

// Sprite combines an Actor with the ability to Draw from a single spritesheet.
type Sprite struct {
	Actor
	FrameSize   image.Point
	FrameOffset image.Point
	Hidden      bool
	ID
	Src ImageRef
	ZPos

	anim *Anim
}

func (s *Sprite) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	if s.Hidden {
		return
	}
	var op ebiten.DrawImageOptions
	dp := s.Pos.Add(s.FrameOffset)
	op.GeoM.Translate(float64(dp.X), float64(dp.Y))
	op.GeoM.Concat(geom)

	frame := s.anim.CurrentFrame()
	src := s.Src.Image()
	w, _ := src.Size()
	sp := image.Pt((frame*s.FrameSize.X)%w, ((frame*s.FrameSize.X)/w)*s.FrameSize.Y)

	screen.DrawImage(src.SubImage(image.Rectangle{sp, sp.Add(s.FrameSize)}).(*ebiten.Image), &op)
}

func (s *Sprite) Scan() []interface{} { return []interface{}{&s.Actor} }

func (s *Sprite) SetAnim(a *Anim) {
	if s.anim != a {
		a.Reset()
	}
	s.anim = a
}

func (s *Sprite) Update() error { return s.anim.Update() }
