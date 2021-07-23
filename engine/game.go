package engine

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// Updater is a component that can update. Update is called repeatedly.
type Updater interface {
	Update() error
}

// Drawer is a component that can draw itself. Draw is called often.
// DrawAfter is used to reorder components.
type Drawer interface {
	Draw(*ebiten.Image)
	DrawAfter(Drawer) bool
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

// Sort sorts the components by draw order.
// Non-Drawers are sorted first.
func (g *Game) Sort() {
	sort.Slice(g.Components, func(i, j int) bool {
		a, aok := g.Components[i].(Drawer)
		b, bok := g.Components[j].(Drawer)
		if aok && bok {
			return b.DrawAfter(a)
		}
		return !aok && bok
	})
}
