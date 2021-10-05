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
	"image"

	"github.com/DrJosh9000/ichigo/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Camera satisfies interfaces.
var _ interface {
	Identifier
	Prepper
	Scanner
	Transformer
} = &Camera{}

func init() {
	gob.Register(&Camera{})
}

// Camera models a camera that is viewing something.
type Camera struct {
	ID
	Child interface{}
	Disables
	Hides

	// Camera controls
	// These directly manipulate the camera. If you want to restrict the camera
	// view area to the child's bounding rectangle, use PointAt.
	Centre   image.Point // voxel coordinates
	Rotation float64     // radians
	Zoom     float64     // unitless

	game *Game
}

// PointAt points the camera at a particular centre point and zoom, but adjusts
// for the bounds of the child component (if available).
func (c *Camera) PointAt(centre geom.Int3, zoom float64) {
	// Special sauce: if Child has a BoundingRect, make some adjustments
	bnd, ok := c.Child.(BoundingRecter)
	if !ok {
		c.Centre = geom.Project(c.game.Projection, centre)
		c.Zoom = zoom
		return
	}

	// The child has boundaries; respect them.
	br := bnd.BoundingRect()

	// The lower bound on zoom is the larger of
	// { (ScreenWidth / BoundsWidth), (ScreenHeight / BoundsHeight) }
	sz := br.Size()
	if z := float64(c.game.ScreenSize.X) / float64(sz.X); zoom < z {
		zoom = z
	}
	if z := float64(c.game.ScreenSize.Y) / float64(sz.Y); zoom < z {
		zoom = z
	}

	// If the configured centre puts the camera out of bounds, move it.
	// Camera frame currently Rectangle{ centre ± (screen/(2*zoom)) }.
	sw2, sh2 := geom.CFloat(c.game.ScreenSize.Div(2))
	swz, shz := int(sw2/zoom), int(sh2/zoom)
	cent := geom.Project(c.game.Projection, centre)
	if cent.X-swz < br.Min.X {
		cent.X = br.Min.X + swz
	}
	if cent.Y-shz < br.Min.Y {
		cent.Y = br.Min.Y + shz
	}
	if cent.X+swz > br.Max.X {
		cent.X = br.Max.X - swz
	}
	if cent.Y+shz > br.Max.Y {
		cent.Y = br.Max.Y - shz
	}
	c.Centre, c.Zoom = cent, zoom
}

// Prepare grabs a copy of game (needed for screen dimensions)
func (c *Camera) Prepare(game *Game) error {
	c.game = game
	return nil
}

// Scan visits c.Child.
func (c *Camera) Scan(visit VisitFunc) error {
	return visit(c.Child)
}

func (c *Camera) String() string { return "Camera@" + c.Centre.String() }

// Transform returns the camera transform.
func (c *Camera) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(geom.CFloat(c.Centre.Mul(-1)))
	opts.GeoM.Scale(c.Zoom, c.Zoom)
	opts.GeoM.Rotate(c.Rotation)
	opts.GeoM.Translate(geom.CFloat(c.game.ScreenSize.Div(2)))
	return opts
}
