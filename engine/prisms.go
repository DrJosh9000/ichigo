package engine

import (
	"encoding/gob"
	"fmt"
	"image"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	_ interface {
		Identifier
		Collider
		Disabler
		Hider
		Prepper
		Scanner
		Transformer
	} = &PrismMap{}

	_ interface {
		BoundingBoxer
		Drawer
		Transformer
	} = &Prism{}
)

func init() {
	gob.Register(&PrismMap{})
	gob.Register(&Prism{})
}

// PrismMap is a generalised 3D tilemap/wallmap/voxelmap etc.
type PrismMap struct {
	ID
	Disables
	Hides
	Ersatz     bool
	Map        map[geom.Int3]*Prism // pos -> prism
	DrawOffset image.Point          // offset applies to whole map
	PosToWorld geom.IntMatrix3x4    // p.pos -> world voxelspace
	PrismSize  geom.Int3            // in world voxelspace units
	PrismTop   []image.Point        // polygon vertices anticlockwise, Y means Z
	Sheet      Sheet

	game      *Game
	pwinverse geom.RatMatrix3
	topext    [4]image.Point
}

// CollidesWith checks if the box collides with any prism.
func (m *PrismMap) CollidesWith(b geom.Box) bool {
	if m.Ersatz {
		return false
	}

	// To find the prisms need to test, we need to invert PosToWorld.

	// Step 1: subtract whatever the translation component of PosToWorld is,
	// reducing the rest of the problem to the 3x3 submatrix.
	rb := b.Sub(m.PosToWorld.Translation())
	// Step 2: invert the rest of the fucking matrix.
	// (Spoilers: I did this already in Prepare)
	rb.Min = m.pwinverse.IntApply(rb.Min)
	rb.Max = m.pwinverse.IntApply(rb.Max) //.Sub(Int3{1, 1, 1}))

	rb = rb.Canon() // inverse might flip the corners around...

	// Check neighboring prisms too because there's a fencepost somewhere here
	rb.Min = rb.Min.Sub(geom.Int3{X: 1, Y: 1, Z: 1})
	rb.Max = rb.Max.Add(geom.Int3{X: 1, Y: 1, Z: 1})

	var pp geom.Int3
	for pp.Z = rb.Min.Z; pp.Z <= rb.Max.Z; pp.Z++ {
		for pp.Y = rb.Min.Y; pp.Y <= rb.Max.Y; pp.Y++ {
			for pp.X = rb.Min.X; pp.X <= rb.Max.X; pp.X++ {
				// Is there a prism here?
				prism, found := m.Map[pp]
				if !found {
					continue
				}
				// Do a cheaper test first against the bounding box.
				if !b.Overlaps(prism.BoundingBox()) {
					continue
				}
				// Exact test that takes into account the prism shape.
				r := b.XZ().Sub(prism.pos.XZ())
				if geom.PolygonRectOverlap(m.PrismTop, r) {
					return true
				}
			}
		}
	}
	return false
}

// Prepare computes an inverse of PosToWorld and prepares all the prisms.
func (m *PrismMap) Prepare(g *Game) error {
	m.game = g
	pwi, err := m.PosToWorld.ToRatMatrix3().Inverse()
	if err != nil {
		return fmt.Errorf("inverting PosToWorld: %w", err)
	}
	m.pwinverse = pwi
	for v, p := range m.Map {
		p.pos = m.PosToWorld.Apply(v)
		p.m = m
	}
	m.topext = geom.PolygonExtrema(m.PrismTop)
	return nil
}

// Scan returns the Sheet and all the Prisms.
func (m *PrismMap) Scan() []interface{} {
	c := make([]interface{}, 1, len(m.Map)+1)
	c[0] = &m.Sheet
	for _, prism := range m.Map {
		c = append(c, prism)
	}
	return c
}

// Transform retrurns a translation by the draw offset.
func (m *PrismMap) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(geom.CFloat(m.DrawOffset))
	return opts
}

// Prism represents a single prism in a PrismMap.
type Prism struct {
	Cell int

	pos geom.Int3 // world coordinates
	m   *PrismMap
}

// BoundingBox returns a bounding box for the prism.
func (p *Prism) BoundingBox() geom.Box {
	return geom.Box{Min: p.pos, Max: p.pos.Add(p.m.PrismSize)}
}

// Draw draws the prism.
func (p *Prism) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(p.m.Sheet.SubImage(p.Cell), opts)
}

// DrawAfter reports if the prism should be drawn after x.
func (p *Prism) DrawAfter(x Drawer) bool {
	pb := p.BoundingBox()
	switch x := x.(type) {
	case *Prism:
		// Fast path for other prisms
		if p.pos.Z == x.pos.Z {
			// TODO: account for projection
			return p.pos.Y < x.pos.Y
		}
		return p.pos.Z > x.pos.Z

	case BoundingBoxer:
		xb := x.BoundingBox()
		// The prism special
		split := p.m.topext[geom.North].X
		threshold := p.m.topext[geom.East].Y
		if xb.Min.X > split {
			threshold = p.m.topext[geom.West].Y
		}
		if pb.Min.Z+threshold <= xb.Min.Z { // x is in front of the front half of p
			return false
		}
		if pb.Min.Z+threshold >= xb.Max.Z { // x is behind the front half of p
			return true
		}
	}
	return false
}

// DrawBefore reports if the prism should be drawn before x.
func (p *Prism) DrawBefore(x Drawer) bool {
	pb := p.BoundingBox()
	switch x := x.(type) {
	case *Prism:
		// Fast path for other prisms
		if p.pos.Z == x.pos.Z {
			// TODO: account for projection
			return p.pos.Y > x.pos.Y
		}
		return p.pos.Z < x.pos.Z

	case BoundingBoxer:
		xb := x.BoundingBox()
		// The prism special
		split := p.m.topext[geom.North].X
		threshold := p.m.topext[geom.East].Y
		if xb.Min.X > split {
			threshold = p.m.topext[geom.West].Y
		}
		if pb.Min.Z+threshold >= xb.Max.Z { // x is behind the front half of p
			return false
		}
		if pb.Min.Z+threshold <= xb.Min.Z { // x is in front of the front half of p
			return true
		}
	}
	return false
}

func (p *Prism) String() string {
	return fmt.Sprintf("Prism(%d)@%v", p.Cell, p.pos)
}

// Transform returns a translation by the projected position.
func (p *Prism) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(geom.CFloat(
		geom.Project(p.m.game.Projection, p.pos),
	))
	return opts
}
