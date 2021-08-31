package engine

import (
	"encoding/gob"
	"image"
)

var _ Collider = SolidRect{}

func init() {
	gob.Register(&SolidRect{})
}

type SolidRect struct {
	ID
	Bounds
}

func (s SolidRect) CollidesWith(r image.Rectangle) bool {
	return s.BoundingRect().Overlaps(r)
}
