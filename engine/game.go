package engine

import "github.com/hajimehoshi/ebiten/v2"

// Game implements the ebiten methods using a collection of components.
type Game struct {
	ScreenWidth  int
	ScreenHeight int
	Scene        *Scene
}

// Draw draws the entire thing.
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen, ebiten.GeoM{})
}

// Layout returns the configured screen width/height.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.ScreenWidth, g.ScreenHeight
}

// Update just passes the call onto Layers.
func (g *Game) Update() error {
	return g.Scene.Update()
}
