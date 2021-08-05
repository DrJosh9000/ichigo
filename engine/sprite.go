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
	Hidden bool
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
	op.GeoM.Translate(float64(s.Pos.X), float64(s.Pos.Y))
	op.GeoM.Concat(geom)

	frame := s.anim.CurrentFrame()
	src := s.Src.Image()
	w, _ := src.Size()
	sp := image.Pt((frame*s.Size.X)%w, ((frame*s.Size.X)/w)*s.Size.Y)

	screen.DrawImage(src.SubImage(image.Rectangle{sp, sp.Add(s.Size)}).(*ebiten.Image), &op)
}

func (s *Sprite) Scan() []interface{} { return []interface{}{&s.Actor} }

func (s *Sprite) SetAnim(a *Anim) { s.anim = a }

func (s *Sprite) Update() error { return s.anim.Update() }
