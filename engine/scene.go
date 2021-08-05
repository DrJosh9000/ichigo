package engine

import (
	"encoding/gob"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(Scene{})
}

// Scene manages drawing and updating a bunch of components.
type Scene struct {
	Components []interface{}
	Disabled   bool
	Hidden     bool
	ID
	Transform GeoMDef
	ZPos
}

// Draw draws all components in order.
func (s *Scene) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	if s.Hidden {
		return
	}
	geom.Concat(*s.Transform.GeoM())

	for _, i := range s.Components {
		if d, ok := i.(Drawer); ok {
			d.Draw(screen, geom)
		}
	}
}

// Prepare does an initial Z-order sort.
func (s *Scene) Prepare(*Game) { s.sortByZ() }

// sortByZ sorts the components by Z position.
// Everything without a Z sorts first. Stable sort is used to avoid Z-fighting
// (among layers without a Z, or those with equal Z).
func (s *Scene) sortByZ() {
	sort.SliceStable(s.Components, func(i, j int) bool {
		a, aok := s.Components[i].(ZPositioner)
		b, bok := s.Components[j].(ZPositioner)
		if aok && bok {
			return a.Z() < b.Z()
		}
		return !aok && bok
	})
}

// Scan returns all immediate subcomponents.
func (s *Scene) Scan() []interface{} { return s.Components }

// Update calls Update on all Updater components.
func (s *Scene) Update() error {
	if s.Disabled {
		return nil
	}

	for _, c := range s.Components {
		// Update each updater in turn
		if u, ok := c.(Updater); ok {
			if err := u.Update(); err != nil {
				return err
			}
		}
	}
	// Check if the updates put the components out of order; if so, sort
	curZ := -math.MaxFloat64 // fun fact: this is min float64
	for _, c := range s.Components {
		z, ok := c.(ZPositioner)
		if !ok {
			continue
		}
		t := z.Z()
		if t < curZ {
			s.sortByZ()
			return nil
		}
		curZ = t
	}
	return nil
}
