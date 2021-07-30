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
	Transform  ebiten.GeoM

	needsSort bool
}

// Draw draws all components in order.
func (s *Scene) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	geom.Concat(s.Transform)
	for _, i := range s.Components {
		if d, ok := i.(Drawer); ok {
			d.Draw(screen, geom)
		}
	}
}

// SetNeedsSort informs l that its layers may be out of order.
func (s *Scene) SetNeedsSort() {
	s.needsSort = true
}

// sortByZ sorts the components by Z position.
// Stable sort is used to avoid Z-fighting among layers without a Z, or
// among those with equal Z. All non-ZPositioners are sorted first.
func (s *Scene) sortByZ() {
	s.needsSort = false
	sort.SliceStable(s.Components, func(i, j int) bool {
		a, aok := s.Components[i].(ZPositioner)
		b, bok := s.Components[j].(ZPositioner)
		if aok && bok {
			return a.Z() < b.Z()
		}
		return !aok && bok
	})
}

// Update calls Update on all Updater components.
func (s *Scene) Update() error {
	for _, c := range s.Components {
		if u, ok := c.(Updater); ok {
			if err := u.Update(); err != nil {
				return err
			}
		}
	}
	if s.needsSort {
		s.sortByZ()
	}
	return nil
}
