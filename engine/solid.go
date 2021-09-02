package engine

import (
	"encoding/gob"
)

var _ Collider = SolidRect{}

func init() {
	gob.Register(&SolidRect{})
}

type SolidRect struct {
	ID
	Box
}

func (s SolidRect) CollidesWith(r Box) bool {
	return s.Box.Overlaps(r)
}
