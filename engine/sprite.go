package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Sprite combines an Actor with the ability to Draw...
type Sprite struct {
	Actor
	AnimRef
	Hidden bool
	ID
	Src ImageRef
	ZPos

	vx, vy     float64       // TODO: refactor
}

func (s *Sprite) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	if s.Hidden {
		return
	}
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(float64(s.Actor.Position.X), float64(s.Actor.Position.Y))
	op.GeoM.Concat(geom)

	frame := s.Anim().CurrentFrame()
	src := s.Src.Image()
	w, _ := src.Size()
	sx, sy := (frame * s.Actor.Size.X) % w, ((frame * s.Actor.Size.X) / w) * s.Actor.Size.Y

	screen.DrawImage(src.SubImage(image.Rect(sx, sy, sx + s.Actor.Size.X, sy+s.Actor.Size.Y)).(*ebiten.Image), &op)
}

func (s *Sprite) Scan() []interface{}{
	return []interface{}{&s.Actor}
}

func (s *Sprite) Update() error {
	// TODO: delegate updating to something else
	if s.Actor.CollidesAt(s.Actor.Position.Add(image.Pt(0, 1))) {
		// Not falling
		s.vy = 0
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// Jump?
			s.vy = -7
		}
	} else {
		// Falling
		s.vy += gravity
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		s.vx = -3
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		s.vx = 3
	default:
		s.vx = 0
	}
	s.Actor.MoveX(s.vx, func() { s.vx = 0 })
	s.Actor.MoveY(s.vy, func() { s.vy = 0 })
	return s.Anim().Update()
}