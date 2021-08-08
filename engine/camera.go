package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	ID
	Scene     *Scene
	Transform GeoMDef

	game *Game
}

func (c *Camera) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	geom.Concat(*c.Transform.GeoM())
	c.Scene.Draw(screen, geom)
}

func (c *Camera) Update() error { return c.Scene.Update() }

func (c *Camera) Scan() []interface{} { return []interface{}{c.Scene} }

func (c *Camera) Prepare(game *Game) { c.game = game }

func (c *Camera) Centre(p image.Point) {
	t := c.Transform.GeoM()
	t.Reset()
	t.Translate(float64(c.game.ScreenWidth/2-p.X), float64(c.game.ScreenHeight/2-p.Y))
}
