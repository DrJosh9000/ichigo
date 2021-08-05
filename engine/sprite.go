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
	*Anim  // TODO: better
	Hidden bool
	ID
	Src ImageRef
	ZPos

	vx, vy     float64 // TODO: refactor
	facingLeft bool

	animIdleLeft, animIdleRight, animRunLeft, animRunRight *Anim
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
	sx, sy := (frame*s.Actor.Size.X)%w, ((frame*s.Actor.Size.X)/w)*s.Actor.Size.Y

	screen.DrawImage(src.SubImage(image.Rect(sx, sy, sx+s.Actor.Size.X, sy+s.Actor.Size.Y)).(*ebiten.Image), &op)
}

func (s *Sprite) Scan() []interface{} {
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
		s.Anim = s.animRunLeft
		s.facingLeft = true
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		s.vx = 2
		s.Anim = s.animRunRight
		s.facingLeft = false
	default:
		s.vx = 0
		s.Anim = s.animIdleRight
		if s.facingLeft {
			s.Anim = s.animIdleLeft
		}
	}
	s.Actor.MoveX(s.vx, func() { s.vx = -s.vx * dampen })
	s.Actor.MoveY(s.vy, func() { s.vy = -s.vy * dampen })
	return s.Anim.Update()
}

func (s *Sprite) Build(g *Game) {
	// TODO: better than this
	s.animRunLeft = &Anim{Def: AnimDefs["aw_run_left"]}
	s.animRunRight = &Anim{Def: AnimDefs["aw_run_right"]}
	s.animIdleLeft = &Anim{Def: AnimDefs["aw_idle_left"]}
	s.animIdleRight = &Anim{Def: AnimDefs["aw_idle_right"]}
}
