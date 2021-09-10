package engine

import (
	"encoding/gob"
	"image"

	"drjosh.dev/gurgle/geom"
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
	Hidden
	Sheet Sheet

	anim *Anim
}

// BoundingBox forwards the call to Actor.
func (s *Sprite) BoundingBox() geom.Box { return s.Actor.BoundingBox() }

// Draw draws the current cell to the screen.
func (s *Sprite) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(s.Sheet.SubImage(s.anim.Cell()), opts)
}

// DrawAfter reports if the sprite must be drawn after x.
func (s *Sprite) DrawAfter(x Drawer) bool {
	sb := s.BoundingBox()
	switch d := x.(type) {
	case BoundingBoxer:
		xb := d.BoundingBox()
		// Z ?
		if sb.Min.Z >= xb.Max.Z { // s is unambiguously in front
			return true
		}
		if sb.Max.Z <= xb.Min.Z { // s is unambiguously behind
			return false
		}
		// Y ? (NB: up is negative)
		if sb.Max.Y <= xb.Min.Y { // s is unambiguously above
			return true
		}
		if sb.Min.Y >= xb.Max.Y { // s is unambiguously below
			return false
		}
		// Hexagon special
		if sb.Min.Z > xb.Min.Z+8 {
			return true
		}
		if sb.Max.Z < sb.Min.Z+8 {
			return false
		}
	case zpositioner:
		return sb.Min.Z > int(d.zposition())
	}
	return false
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

// Transform returns a translation by the DrawOffset and Actor.Pos projected
func (s *Sprite) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(geom.CFloat(
		// Reaching into Actor for a reference to Game so I don't have to
		// implement Prepare in this file, but writing this long comment
		// providing exposition...
		s.Actor.game.Projection.Project(s.Actor.Pos).Add(s.DrawOffset),
	))
	return opts
}

// Update updates the Sprite's anim. anim can change a bit so we don't tell Game
// about it, but that means it must be updated manually.
func (s *Sprite) Update() error { return s.anim.Update() }
