package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Camera satisfies interfaces.
var _ interface {
	Identifier
	Prepper
	Transformer
} = &Camera{}

func init() {
	gob.Register(&Camera{})
}

// Camera models a camera that is viewing a scene. (Camera is a child of the
// scene it is viewing, for various reasons.) Changes to the fields take effect
// immediately.
type Camera struct {
	ID

	// Camera controls
	Centre   image.Point // world coordinates
	Filter   ebiten.Filter
	Rotation float64 // radians
	Zoom     float64 // unitless

	game *Game
}

// Prepare grabs a copy of game (needed for screen dimensions)
func (c *Camera) Prepare(game *Game) error {
	c.game = game
	return nil
}

// Transform returns the camera transform.
func (c *Camera) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(float2(c.Centre.Mul(-1)))
	opts.GeoM.Scale(c.Zoom, c.Zoom)
	opts.GeoM.Rotate(c.Rotation)
	opts.GeoM.Translate(float64(c.game.ScreenWidth/2), float64(c.game.ScreenHeight/2))
	return opts
}
