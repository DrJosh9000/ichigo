package engine

import (
	"encoding/gob"
	"image"
)

func init() {
	gob.Register(SolidRect{})
}

type SolidRect struct {
	Rect image.Rectangle
}

func (s SolidRect) CollidesWith(r image.Rectangle) bool { return s.Rect.Overlaps(r) }
