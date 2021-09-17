package engine

import (
	"encoding/gob"
	"fmt"
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

// DrawAfter reports if the sprite should be drawn after x.
func (s *Sprite) DrawAfter(x Drawer) bool {
	if !commonDrawerComparisons {
		sb := s.BoundingBox()
		switch x := x.(type) {
		case BoundingBoxer:
			xb := x.BoundingBox()
			if sb.Max.Z <= xb.Min.Z { // s is behind x
				return false
			}
			if sb.Min.Z >= xb.Max.Z { // s is in front of x
				return true
			}
			if sb.Min.Y >= xb.Max.Y { // s is below x
				return false
			}
			if sb.Max.Y <= xb.Min.Y { // s is above x
				return true
			}
		case ZPositioner:
			return sb.Min.Z > x.ZPos() // s is after
		}
	}
	return false
}

// DrawBefore reports if the sprite should be drawn before x.
func (s *Sprite) DrawBefore(x Drawer) bool {
	if !commonDrawerComparisons {
		sb := s.BoundingBox()
		switch x := x.(type) {
		case BoundingBoxer:
			xb := x.BoundingBox()
			if sb.Min.Z >= xb.Max.Z { // s is in front of x
				return false
			}
			if sb.Max.Z <= xb.Min.Z { // s is behind x
				return true
			}
			if sb.Max.Y <= xb.Min.Y { // s is above x
				return false
			}
			if sb.Min.Y >= xb.Max.Y { // s is below x
				return true
			}
		case ZPositioner:
			return sb.Max.Z < x.ZPos() // s is before
		}
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

func (s *Sprite) String() string {
	return fmt.Sprintf("Sprite@%v", s.Actor.Pos)
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
