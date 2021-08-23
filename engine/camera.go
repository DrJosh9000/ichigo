package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Camera satisfies interfaces.
var (
	_ Identifier = &Camera{}
	_ Drawer     = &Camera{}
	_ Prepper    = &Camera{}
	_ Scanner    = &Camera{}
	_ Updater    = &Camera{}
)

func init() {
	gob.Register(Camera{})
}

// Camera models a camera that is viewing a scene.
// Changes to the configuration take effect immediately.
// Camera ignores Scene.Draw and calls Scene's children's Draw.
type Camera struct {
	ID
	Scene Scener

	// Camera controls
	Centre image.Point // world coordinates
	Filter ebiten.Filter
	Zoom   float64 // unitless

	game *Game
}

// Draw applies transformations to opts, then calls c.Scene.Draw with it.
func (c *Camera) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if c.Scene.IsHidden() {
		return
	}

	br := c.Scene.BoundingRect()

	// The lower bound on zoom is the larger of
	// { (ScreenWidth / BoundsWidth), (ScreenHeight / BoundsHeight) }
	zoom := c.Zoom
	sz := br.Size()
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
	if centre.X-swz < br.Min.X {
		centre.X = br.Min.X + swz
	}
	if centre.Y-shz < br.Min.Y {
		centre.Y = br.Min.Y + shz
	}
	if centre.X+swz > br.Max.X {
		centre.X = br.Max.X - swz
	}
	if centre.Y+shz > br.Max.Y {
		centre.Y = br.Max.Y - shz
	}

	// Apply other options
	opts.Filter = c.Filter

	// Compute common matrix (parts independent of parallax).
	// Moving centre to the origin happens per component.
	var comm ebiten.GeoM
	// 2. Zoom (this is also where rotation would be)
	comm.Scale(zoom, zoom)
	// 3. Move the origin to the centre of screen space.
	comm.Translate(sw2, sh2)
	// 4. Apply transforms from the caller.
	comm.Concat(opts.GeoM)

	// Draw everything.
	for _, i := range c.Scene.Scan() {
		if d, ok := i.(Drawer); ok {
			pf := 1.0
			if s, ok := i.(ParallaxScaler); ok {
				pf = s.ParallaxFactor()
			}
			var geom ebiten.GeoM
			// 1. Move centre to the origin, subject to parallax factor
			geom.Translate(-float64(centre.X)*pf, -float64(centre.Y)*pf)
			geom.Concat(comm)
			opts.GeoM = geom
			d.Draw(screen, opts)
		}
	}
}

// Update passes the call to c.Scene.
func (c *Camera) Update() error { return c.Scene.Update() }

// Scan returns the only child (c.Scene).
func (c *Camera) Scan() []interface{} { return []interface{}{c.Scene} }

// Prepare grabs a copy of game (needed for screen dimensions)
func (c *Camera) Prepare(game *Game) { c.game = game }
