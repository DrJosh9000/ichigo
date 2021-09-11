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

// DrawAfter reports if the sprite should be drawn after x.
func (s *Sprite) DrawAfter(x Drawer) bool {
	sb := s.BoundingBox()
	switch d := x.(type) {
	case BoundingBoxer:
		xb := d.BoundingBox()
		if sb.Max.Z <= xb.Min.Z { // s is behind x
			/*if xb.Min.Y < 0 {
				log.Print("s.DrawAfter: sprite is behind prism")
			}*/
			return false
		}
		if sb.Min.Z >= xb.Max.Z { // s is in front of x
			/*if xb.Min.Y < 0 {
				log.Print("s.DrawAfter: sprite is in front of prism")
			}*/
			return true
		}
		if sb.Min.Y >= xb.Max.Y { // s is below x
			/*if xb.Min.Y < 0 {
				log.Print("s.DrawAfter: sprite is below prism")
			}*/
			return false
		}
		if sb.Max.Y <= xb.Min.Y { // s is above x
			/*if xb.Min.Y < 0 {
				log.Print("s.DrawAfter: sprite is above prism")
			}*/
			return true
		}
	case zpositioner:
		return sb.Min.Z > int(d.zposition()) // s is after
	}
	return false
}

// DrawBefore reports if the sprite should be drawn before x.
func (s *Sprite) DrawBefore(x Drawer) bool {
	sb := s.BoundingBox()
	switch d := x.(type) {
	case BoundingBoxer:
		xb := d.BoundingBox()
		if sb.Min.Z >= xb.Max.Z { // s is in front of x
			/*if xb.Min.Y < 0 {
				log.Print("s.DrawBefore: sprite is in front of prism")
			}*/
			return false
		}
		if sb.Max.Z <= xb.Min.Z { // s is behind x
			/*if xb.Min.Y < 0 {
				log.Print("s.DrawBefore: sprite is behind prism")
			}*/
			return true
		}
		if sb.Max.Y <= xb.Min.Y { // s is above x
			/*if xb.Min.Y < 0 {
				log.Print("s.DrawBefore: sprite is above prism")
			}*/
			return false
		}
		if sb.Min.Y >= xb.Max.Y { // s is below x
			/*if xb.Min.Y < 0 {
				log.Print("s.DrawBefore: sprite is below prism")
			}*/
			return true
		}
	case zpositioner:
		return sb.Max.Z < int(d.zposition()) // s is before
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
