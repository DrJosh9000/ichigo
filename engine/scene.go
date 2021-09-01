package engine

import (
	"encoding/gob"
)

// Ensure Scene satisfies Scener.
var _ Scener = &Scene{}

func init() {
	gob.Register(&Scene{})
}

// Scene just contains a bunch of components.
type Scene struct {
	ID
	Bounds     // world coordinates
	Components []interface{}
	Disabled
	Hidden
}

// Scan returns all immediate subcomponents (including the camera, if not nil).
func (s *Scene) Scan() []interface{} { return s.Components }
