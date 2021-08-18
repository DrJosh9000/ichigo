package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Sprite satisfies interfaces.
var (
	_ Identifier  = &Sprite{}
	_ Drawer      = &Sprite{}
	_ DrawOrderer = &Sprite{}
	_ Scanner     = &Sprite{}
	_ Updater     = &Sprite{}
)

func init() {
	gob.Register(Sprite{})
}

// Sprite combines an Actor with the ability to Draw from a single spritesheet.
type Sprite struct {
	Actor
	ZOrder
	FrameSize   image.Point
	FrameOffset image.Point
	Hidden      bool
	ID
	Src ImageRef

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

	frame := s.anim.CurrentFrame()
	src := s.Src.Image()
	w, _ := src.Size()
	sp := image.Pt((frame*s.FrameSize.X)%w, ((frame*s.FrameSize.X)/w)*s.FrameSize.Y)

	screen.DrawImage(src.SubImage(image.Rectangle{sp, sp.Add(s.FrameSize)}).(*ebiten.Image), &opts)
}

func (s *Sprite) Scan() []interface{} { return []interface{}{&s.Actor} }

func (s *Sprite) SetAnim(a *Anim) {
	if s.anim != a {
		a.Reset()
	}
	s.anim = a
}

func (s *Sprite) Update() error { return s.anim.Update() }
