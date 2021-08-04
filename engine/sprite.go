package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const dampen = 0.5
const gravity = 0.2

// Sprite combines an Actor with the ability to Draw...
type Sprite struct {
	Actor
	*Anim // TODO: better
	Hidden bool
	ID
	Src ImageRef
	ZPos

	vx, vy     float64       // TODO: refactor
	animIdle, animWalkLeft, animWalkRight *Anim
}

func (s *Sprite) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	if s.Hidden {
		return
	}
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(float64(s.Actor.Position.X), float64(s.Actor.Position.Y))
	op.GeoM.Concat(geom)

	frame := s.Anim.CurrentFrame()
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
			s.vy = -5
		}
	} else {
		// Falling
		s.vy += gravity
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		s.vx = -2
		s.Anim = s.animWalkLeft
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		s.vx = 2
		s.Anim = s.animWalkRight
	default:
		s.vx = 0
		s.Anim = s.animIdle
	}
	s.Actor.MoveX(s.vx, func() { s.vx = -s.vx * dampen })
	s.Actor.MoveY(s.vy, func() { s.vy = -s.vy * dampen })
	return s.Anim.Update()
}

func (s *Sprite) Build(g *Game) {
	// TODO: better than this
	s.animWalkLeft = &Anim{Def: AnimDefs["aw_walk_left"]}
	s.animWalkRight = &Anim{Def: AnimDefs["aw_walk_right"]}
	s.animIdle = &Anim{Def: AnimDefs["aw_idle"]}
}