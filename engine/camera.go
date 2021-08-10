package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	ID
	Scene *Scene

	// camera controls
	Zoom   float64
	Centre image.Point

	game *Game
	// TODO: camera constraints
}

func (c *Camera) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	//geom.Concat(*c.Transform.GeoM())
	scx, scy := float64(c.game.ScreenWidth/2), float64(c.game.ScreenHeight/2)
	geom.Translate((scx/c.Zoom - float64(c.Centre.X)), (scy/c.Zoom - float64(c.Centre.Y)))
	geom.Scale(c.Zoom, c.Zoom)
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
