package engine

import "github.com/hajimehoshi/ebiten/v2"

// Transform is a bucket of things that affect drawing.
type Transform struct {
	// Projection is used by isometric 3D components to project their
	// coordinates into 2D. There's usually only one component in the tree that
	// sets this field, but it would apply to all descendants.
	Projection IntProjection

	// Opts contains the 2D geometry matrix, the colour matrix, filter mode, and
	// composition mode.
	Opts ebiten.DrawImageOptions
}

// Concat returns the combined transform (a transform equivalent to applying t
// and then u).
func (t Transform) Concat(u Transform) Transform {
	if u.Projection != (IntProjection{}) {
		t.Projection = u.Projection
	}
	t.Opts.ColorM.Concat(u.Opts.ColorM)
	t.Opts.GeoM.Concat(u.Opts.GeoM)
	if u.Opts.CompositeMode != 0 {
		t.Opts.CompositeMode = u.Opts.CompositeMode
	}
	if u.Opts.Filter != 0 {
		t.Opts.Filter = u.Opts.Filter
	}
	return t
}
