package engine

import (
	"encoding/gob"
	"fmt"
	"image"

	"github.com/DrJosh9000/ichigo/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Sprite satisfies interfaces.
var _ interface {
	BoundingBoxer
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
	Actor      Actor
	DrawOffset image.Point
	Hides
	Sheet Sheet

	anim *Anim
}

// BoundingBox forwards the call to Actor.
func (s *Sprite) BoundingBox() geom.Box { return s.Actor.BoundingBox() }

// Draw draws the current cell to the screen.
func (s *Sprite) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(s.Sheet.SubImage(s.anim.Cell()), opts)
}

// Scan visits &s.Actor and &s.Sheet.
func (s *Sprite) Scan(visit VisitFunc) error {
	if err := visit(&s.Actor); err != nil {
		return err
	}
	return visit(&s.Sheet)
}

// Anim returns the current Anim.
func (s *Sprite) Anim() *Anim { return s.anim }

// SetAnim sets the Anim to use for the sprite. If it is not the same as the
// one currently set, it resets the new anim.
func (s *Sprite) SetAnim(a *Anim) {
	if s.anim != a {
		a.Reset()
	}
	s.anim = a
}

func (s *Sprite) String() string {
	return fmt.Sprintf("Sprite@%v", s.Actor.Pos)
}

// Transform returns a translation by the DrawOffset and Actor.Pos projected
func (s *Sprite) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(geom.CFloat(
		// Reaching into Actor for a reference to Game so I don't have to
		// implement Prepare in this file, but writing this long comment
		// providing exposition...
		geom.Project(s.Actor.game.Projection, s.Actor.Pos).
			Add(s.DrawOffset),
	))
	return opts
}

// Update updates the Sprite's anim. anim can change a bit so we don't tell Game
// about it, but that means it must be updated manually.
func (s *Sprite) Update() error { return s.anim.Update() }
