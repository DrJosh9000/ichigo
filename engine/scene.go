package engine

import (
	"encoding/gob"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Scene satisfies Scener.
var _ Scener = &Scene{}

func init() {
	gob.Register(&Scene{})
}

// Scene manages drawing and updating a bunch of components.
type Scene struct {
	ID
	Bounds             // world coordinates
	Camera     *Camera // optional; applies a bunch of transforms to draw calls
	Components []interface{}
	Disabled
	Hidden
	ZOrder
}

/*
// Draw draws all components in order.
func (s *Scene) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if s.Hidden {
		return
	}
	if s.Camera == nil {
		// Draw everything, no camera transforms.
		for _, i := range s.Components {
			if d, ok := i.(Drawer); ok {
				d.Draw(screen, opts)
			}
		}
		return
	}

	// There is a camera; apply camera transforms.
	br := s.BoundingRect()

	// The lower bound on zoom is the larger of
	// { (ScreenWidth / BoundsWidth), (ScreenHeight / BoundsHeight) }
	zoom := s.Camera.Zoom
	sz := br.Size()
	if z := float64(s.Camera.game.ScreenWidth) / float64(sz.X); zoom < z {
		zoom = z
	}
	if z := float64(s.Camera.game.ScreenHeight) / float64(sz.Y); zoom < z {
		zoom = z
	}

	// If the configured centre puts the camera out of bounds, move it.
	centre := s.Camera.Centre
	// Camera frame currently Rectangle{ centre ± (screen/(2*zoom)) }.
	sw2, sh2 := float64(s.Camera.game.ScreenWidth/2), float64(s.Camera.game.ScreenHeight/2)
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
	opts.Filter = s.Camera.Filter

	// Compute common matrix (parts independent of parallax, which is step 1).
	// Moving centre to the origin happens per component.
	var comm ebiten.GeoM
	// 2. Zoom and rotate
	comm.Scale(zoom, zoom)
	comm.Rotate(s.Camera.Rotation)
	// 3. Move the origin to the centre of screen space.
	comm.Translate(sw2, sh2)
	// 4. Apply transforms from the caller.
	comm.Concat(opts.GeoM)

	// Draw everything.
	for _, i := range s.Components {
		d, ok := i.(Drawer)
		if !ok {
			continue
		}
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
*/

// Transform returns the camera transform
func (s *Scene) Transform() ebiten.DrawImageOptions {
	if s.Camera == nil {
		return ebiten.DrawImageOptions{}
	}
	return s.Camera.Transform()
}

// Scan returns all immediate subcomponents (including the camera, if not nil).
func (s *Scene) Scan() []interface{} {
	if s.Camera != nil {
		return append(s.Components, s.Camera)
	}
	return s.Components
}
