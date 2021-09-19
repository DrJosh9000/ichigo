package geom

import "image"

// Box describes an axis-aligned rectangular prism.
type Box struct {
	Min, Max Int3
}

// String returns a string representation of b like "(3,4,5)-(6,5,8)".
func (b Box) String() string {
	return b.Min.String() + "-" + b.Max.String()
}

// Empty reports whether the box contains no points.
func (b Box) Empty() bool {
	return b.Min.X >= b.Max.X || b.Min.Y >= b.Max.Y || b.Min.Z >= b.Max.Z
}

// Eq reports whether b and c contain the same set of points. All empty boxes
// are considered equal.
func (b Box) Eq(c Box) bool {
	return b == c || b.Empty() && c.Empty()
}

// Overlaps reports whether b and c have non-empty intersection.
func (b Box) Overlaps(c Box) bool {
	return !b.Empty() && !c.Empty() &&
		b.Min.X < c.Max.X && c.Min.X < b.Max.X &&
		b.Min.Y < c.Max.Y && c.Min.Y < b.Max.Y &&
		b.Min.Z < c.Max.Z && c.Min.Z < b.Max.Z
}

// Size returns b's width, height, and depth.
func (b Box) Size() Int3 {
	return b.Max.Sub(b.Min)
}

// Centre returns the centre point of the box.
func (b Box) Centre() Int3 {
	return b.Min.Add(b.Max).Div(2)
}

// Add offsets the box by vector p.
func (b Box) Add(p Int3) Box {
	return Box{
		Min: b.Min.Add(p),
		Max: b.Max.Add(p),
	}
}

// Sub offsets the box by (-p).
func (b Box) Sub(p Int3) Box {
	return Box{
		Min: b.Min.Sub(p),
		Max: b.Max.Sub(p),
	}
}

// Canon returns a copy of b that is well-formed.
func (b Box) Canon() Box {
	if b.Max.X < b.Min.X {
		b.Min.X, b.Max.X = b.Max.X, b.Min.X
	}
	if b.Max.Y < b.Min.Y {
		b.Min.Y, b.Max.Y = b.Max.Y, b.Min.Y
	}
	if b.Max.Z < b.Min.Z {
		b.Min.Z, b.Max.Z = b.Max.Z, b.Min.Z
	}
	return b
}

// BoundingRect returns an image.Rectangle that bounds the box if it were
// projected.
func (b Box) BoundingRect(π Projector) image.Rectangle {
	return b.Back(π).Union(b.Front(π))
}

// Back returns an image.Rectangle representing the back of the box, using
// the given projection π.
func (b Box) Back(π Projector) image.Rectangle {
	p := π.Project(b.Min.Z)
	return image.Rectangle{
		Min: b.Min.XY().Add(p),
		Max: b.Max.XY().Add(p),
	}
}

// Front returns an image.Rectangle representing the front of the box, using
// the given projection π.
func (b Box) Front(π Projector) image.Rectangle {
	p := π.Project(b.Max.Z)
	return image.Rectangle{
		Min: b.Min.XY().Add(p),
		Max: b.Max.XY().Add(p),
	}
}

// XY returns the image.Rectangle representing the box if we forgot about Z
// (a vertical cross-section).
func (b Box) XY() image.Rectangle {
	return image.Rectangle{
		Min: b.Min.XY(),
		Max: b.Max.XY(),
	}
}

// XZ returns the image.Rectangle representing the box if we forgot about Y
// (the horizontal cross-section).
func (b Box) XZ() image.Rectangle {
	return image.Rectangle{
		Min: b.Min.XZ(),
		Max: b.Max.XZ(),
	}
}
