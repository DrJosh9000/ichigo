package engine

import (
	"encoding/gob"
	"image"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Scene satisfies Scener.
var _ Scener = &Scene{}

func init() {
	gob.Register(Scene{})
}

// Scene manages drawing and updating a bunch of components.
type Scene struct {
	Bounds     image.Rectangle // world coordinates
	Components []interface{}
	Disabled   bool
	Hidden     bool
	ID
	ZOrder
}

// Draw draws all components in order.
func (s *Scene) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if s.Hidden {
		return
	}
	// Draw everything.
	for _, i := range s.Components {
		if d, ok := i.(Drawer); ok {
			d.Draw(screen, opts)
		}
	}
}

// Prepare does an initial Z-order sort.
func (s *Scene) Prepare(game *Game) { s.sortByDrawOrder() }

// sortByDrawOrder sorts the components by Z position.
// Everything without a Z sorts first. Stable sort is used to avoid Z-fighting
// (among layers without a Z, or those with equal Z).
func (s *Scene) sortByDrawOrder() {
	sort.SliceStable(s.Components, func(i, j int) bool {
		a, aok := s.Components[i].(DrawOrderer)
		b, bok := s.Components[j].(DrawOrderer)
		if aok && bok {
			return a.DrawOrder() < b.DrawOrder()
		}
		return !aok && bok
	})
}

// Scan returns all immediate subcomponents.
func (s *Scene) Scan() []interface{} { return s.Components }

// Scene returns itself.
func (s *Scene) Scene() *Scene { return s }

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
	cz := -math.MaxFloat64 // fun fact: this is min float64
	for _, c := range s.Components {
		z, ok := c.(DrawOrderer)
		if !ok {
			continue
		}
		if t := z.DrawOrder(); t > cz {
			cz = t
			continue
		}
		s.sortByDrawOrder()
		return nil
	}
	return nil
}
