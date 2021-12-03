/*
Copyright 2021 Josh Deprez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	return visit.Many(&s.Actor, &s.Sheet)
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
