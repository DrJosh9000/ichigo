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
	// Currently the centre of the screen c is A^-1.c in world coordinates
	// So it is off by (p - A^-1.c)
	t := c.Transform.GeoM()
	t.Invert()
	wcx, wcy := t.Apply(float64(c.game.ScreenWidth/2), float64(c.game.ScreenHeight/2))
	t.Translate(float64(p.X)-wcx, float64(p.Y)-wcy)
	t.Invert()
}

func (c *Camera) Zoom(f float64) {
	t := c.Transform.GeoM()
	t.Scale(f, f)
}
