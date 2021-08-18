package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera models a camera that is viewing a scene.
// Changes to the configuration take effect immediately.
// Camera ignores Scene.Draw and calls Scene's children's Draw.
type Camera struct {
	ID
	Scene *Scene

	// Camera controls
	Bounds image.Rectangle // world coordinates
	Centre image.Point     // world coordinates
	Filter ebiten.Filter
	Zoom   float64 // unitless

	game *Game
}

// Draw applies transformations to opts, then calls c.Scene.Draw with it.
func (c *Camera) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if c.Scene.Hidden {
		return
	}

	// Compute the geometry matrix for the camera controls.

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

	// Apply other options
	opts.Filter = c.Filter

	// Draw everything.
	og := opts.GeoM
	for _, i := range c.Scene.Components {
		if d, ok := i.(Drawer); ok {
			cs := 1.0
			if s, ok := i.(CoordScaler); ok {
				cs = s.CoordScale()
			}
			var geom ebiten.GeoM
			// 1. Move centre to the origin, subject to CoordScale
			geom.Translate(-float64(centre.X)*cs, -float64(centre.Y)*cs)
			// 2. Zoom (this is also where rotation would be)
			geom.Scale(zoom, zoom)
			// 3. Move the origin to the centre of screen space.
			geom.Translate(sw2, sh2)
			// 4. Apply transforms from the caller.
			geom.Concat(og)
			opts.GeoM = geom
			d.Draw(screen, opts)
		}
	}
}

// Update passes the call to c.Scene.
func (c *Camera) Update() error { return c.Scene.Update() }

// Scan returns the only child (c.Scene).
func (c *Camera) Scan() []interface{} { return []interface{}{c.Scene} }

// Prepare grabs a copy of game.
func (c *Camera) Prepare(game *Game) { c.game = game }
