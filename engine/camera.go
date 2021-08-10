package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	ID
	Scene *Scene

	// camera controls
	Centre   image.Point
	Rotation float64
	Zoom     float64

	game *Game
	// TODO: camera constraints
}

func (c *Camera) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	// move the c.Centre to the origin
	geom.Translate(-float64(c.Centre.X), -float64(c.Centre.Y))
	// zoom and rotate
	geom.Scale(c.Zoom, c.Zoom)
	geom.Rotate(c.Rotation)
	// move the origin to the centre of screen space
	geom.Translate(float64(c.game.ScreenWidth/2), float64(c.game.ScreenHeight/2))
	c.Scene.Draw(screen, geom)
}

func (c *Camera) Update() error { return c.Scene.Update() }

func (c *Camera) Scan() []interface{} { return []interface{}{c.Scene} }

func (c *Camera) Prepare(game *Game) {
	c.game = game
	if c.Zoom == 0 {
		c.Zoom = 1
	}
}
