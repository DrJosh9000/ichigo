package engine

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// Drawer is a component that can draw itself. Draw is called often.
type Drawer interface {
	Draw(screen *ebiten.Image, geom ebiten.GeoM)
}

// Updater is a component that can update. Update is called repeatedly.
type Updater interface {
	Update() error
}

// ZPositioner is used to reorder layers.
type ZPositioner interface {
	Z() float64
}

// Scene manages drawing and updating a bunch of components.
type Scene struct {
	Components []interface{}
	needsSort  bool
}

// Draw draws all components in order.
func (l *Scene) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	for _, i := range l.Components {
		if d, ok := i.(Drawer); ok {
			d.Draw(screen, geom)
		}
	}
}

// SetNeedsSort informs l that its layers may be out of order.
func (l *Scene) SetNeedsSort() {
	l.needsSort = true
}

// sortByZ sorts the components by Z position.
// Stable sort is used to avoid Z-fighting among layers without a Z, or
// among those with equal Z. All non-ZPositioners are sorted first.
func (l *Scene) sortByZ() {
	l.needsSort = false
	sort.SliceStable(l.Components, func(i, j int) bool {
		a, aok := l.Components[i].(ZPositioner)
		b, bok := l.Components[j].(ZPositioner)
		if aok && bok {
			return a.Z() < b.Z()
		}
		return !aok && bok
	})
}

// Update calls Update on all Updater components.
func (l *Scene) Update() error {
	for _, c := range l.Components {
		if u, ok := c.(Updater); ok {
			if err := u.Update(); err != nil {
				return err
			}
		}
	}
	if l.needsSort {
		l.sortByZ()
	}
	return nil
}
