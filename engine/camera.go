package engine

import (
	"encoding/gob"
	"image"
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

	// Camera controls
	// These directly manipulate the camera. If you want to restrict the camera
	// view area to the child's bounding rectangle, use PointAt.
	Centre     image.Point // world coordinates
	Rotation   float64     // radians
	Zoom       float64     // unitless
	Projection IntProjection

	game *Game
}

// PointAt points the camera at a particular centre point and zoom, but adjusts
// for the bounds of the child component (if available).
func (c *Camera) PointAt(centre Point3, zoom float64) {
	// Special sauce: if Child has a BoundingRect, make some adjustments
	bnd, ok := c.Child.(Bounder)
	if !ok {
		c.Centre = c.Projection.Project(centre)
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
	sw2, sh2 := cfloat(c.game.ScreenSize.Div(2))
	swz, shz := int(sw2/zoom), int(sh2/zoom)
	cent := c.Projection.Project(centre)
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

// Scan returns s.Child.
func (c *Camera) Scan() []interface{} { return []interface{}{c.Child} }

// Transform returns the camera transform.
func (c *Camera) Transform(pt Transform) (tf Transform) {
	tf.Projection = c.Projection
	tf.Opts.GeoM.Translate(cfloat(c.Centre.Mul(-1)))
	tf.Opts.GeoM.Scale(c.Zoom, c.Zoom)
	tf.Opts.GeoM.Rotate(c.Rotation)
	tf.Opts.GeoM.Translate(cfloat(c.game.ScreenSize.Div(2)))
	return tf.Concat(pt)
}
