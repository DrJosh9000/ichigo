package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Updater is a component that can update (called repeatedly).
type Updater interface {
	Update() error
}

// Drawer is a component that can draw itself (called repeatedly).
type Drawer interface {
	Draw(*ebiten.Image)
}

// Game implements the ebiten methods using a collection of components.
type Game struct {
	ScreenWidth  int
	ScreenHeight int
	Components   []interface{}
}

// Update calls Update on all Updater components.
func (g *Game) Update() error {
	for _, c := range g.Components {
		if u, ok := c.(Updater); ok {
			if err := u.Update(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Draw calls Draw on all Drawer components.
func (g *Game) Draw(screen *ebiten.Image) {
	for _, c := range g.Components {
		if d, ok := c.(Drawer); ok {
			d.Draw(screen)
		}
	}
}

// Layout returns the configured screen width/height.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.ScreenWidth, g.ScreenHeight
}
