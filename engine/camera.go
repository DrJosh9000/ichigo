package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera models a camera that is viewing a scene. Changes to the
// configuration take effect immediately.
type Camera struct {
	ID
	Scene *Scene

	// Camera controls
	Bounds image.Rectangle // world coordinates
	Centre image.Point     // world coordinates
	//Rotation float64       // radians
	Zoom float64 // unitless

	Filter ebiten.Filter

	game *Game
}

// Draw applies transformations to opts, then calls c.Scene.Draw with it.
func (c *Camera) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	// The lower bound on zoom is the larger of
	// { (ScreenWidth / BoundsWidth), (ScreenHeight / BoundsHeight) }
	zoom := c.Zoom
	sz := c.Bounds.Size()
	if z := float64(c.game.ScreenWidth) / float64(sz.X); zoom < z {
		zoom = z
	}
	if z := float64(c.game.ScreenHeight) / float64(sz.Y); zoom < z {
		zoom = z
	}

	// If the configured centre puts the camera out of bounds, move it.
	centre := c.Centre
	// Camera frame currently Rectangle{ centre ± (screen/(2*zoom)) }.
	sw2, sh2 := float64(c.game.ScreenWidth/2), float64(c.game.ScreenHeight/2)
	swz, shz := int(sw2/zoom), int(sh2/zoom)
	if centre.X-swz < c.Bounds.Min.X {
		centre.X = c.Bounds.Min.X + swz
	}
	if centre.Y-shz < c.Bounds.Min.Y {
		centre.Y = c.Bounds.Min.Y + shz
	}
	if centre.X+swz > c.Bounds.Max.X {
		centre.X = c.Bounds.Max.X - swz
	}
	if centre.Y+shz > c.Bounds.Max.Y {
		centre.Y = c.Bounds.Max.Y - shz
	}

	// Apply camera controls to geom.
	// 1. Move centre to the origin
	opts.GeoM.Translate(-float64(centre.X), -float64(centre.Y))
	// 2. Zoom and rotate
	opts.GeoM.Scale(zoom, zoom)
	//geom.Rotate(c.Rotation)
	// 3. Move the origin to the centre of screen space.
	opts.GeoM.Translate(sw2, sh2)

	// Apply other options
	opts.Filter = c.Filter

	c.Scene.Draw(screen, opts)
}

// Update passes the call to c.Scene.
func (c *Camera) Update() error { return c.Scene.Update() }

// Scan returns the only child (c.Scene).
func (c *Camera) Scan() []interface{} { return []interface{}{c.Scene} }

// Prepare grabs a copy of game.
func (c *Camera) Prepare(game *Game) { c.game = game }
