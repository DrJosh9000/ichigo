package engine

import (
	"encoding/gob"

	"github.com/DrJosh9000/ichigo/geom"
)

var _ Collider = SolidRect{}

func init() {
	gob.Register(&SolidRect{})
}

type SolidRect struct {
	ID
	geom.Box
}

func (s SolidRect) CollidesWith(r geom.Box) bool {
	return s.Box.Overlaps(r)
}
