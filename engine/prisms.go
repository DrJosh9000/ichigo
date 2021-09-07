package engine

import (
	"encoding/gob"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	_ interface {
		Identifier
		Collider
		Disabler
		Hider
		Prepper
		Transformer
	} = &PrismMap{}

	_ interface {
		Drawer
		Transformer
	} = &Prism{}
)

func init() {
	gob.Register(&PrismMap{})
	gob.Register(&Prism{})
}

// PrismMap is a generalised tilemap/wallmap/etc.
type PrismMap struct {
	ID
	Disabled
	Hidden
	Ersatz        bool
	Map           map[Int3]*Prism // pos -> prism
	DrawOrderBias image.Point     // dot with pos.XY() = bias value
	DrawOffset    image.Point     // offset applies to whole map
	PosToWorld    IntMatrix3x4    // p.pos -> world voxelspace
	PrismSize     Int3            // in world voxelspace units
	Sheet         Sheet

	game      *Game
	pwinverse RatMatrix3
}

func (m *PrismMap) CollidesWith(b Box) bool {
	if m.Ersatz {
		return false
	}
	// To find the prisms need to test, we need to invert PosToWorld.

	// Step 1: subtract whatever the translation component of PosToWorld is,
	// reducing the rest of the problem to the 3x3 submatrix.
	b = b.Sub(m.PosToWorld.Translation())
	// Step 2: invert the rest of the fucking matrix.
	// (Spoilers: I did this already in Prepare)
	b.Min = m.pwinverse.IntApply(b.Min)
	b.Max = m.pwinverse.IntApply(b.Max.Sub(Int3{1, 1, 1}))
	for k := b.Min.Z; k < b.Max.Z; k++ {
		for j := b.Min.Y; j < b.Max.Y; j++ {
			for i := b.Min.X; i < b.Max.X; i++ {
				// TODO: take into account the prism shape...
				if _, found := m.Map[Int3{i, j, k}]; found {
					return true
				}
			}
		}
	}
	return false
}

func (m *PrismMap) Prepare(g *Game) error {
	m.game = g
	pwi, err := m.PosToWorld.ToRatMatrix3().Inverse()
	if err != nil {
		return fmt.Errorf("inverting PosToWorld: %w", err)
	}
	m.pwinverse = pwi
	log.Printf("inverted PosToWorld: %v", pwi)
	for v, p := range m.Map {
		p.pos = v
		p.m = m
	}
	return nil
}

func (m *PrismMap) Scan() []interface{} {
	c := make([]interface{}, 1, len(m.Map)+1)
	c[0] = &m.Sheet
	for _, prism := range m.Map {
		c = append(c, prism)
	}
	return c
}

func (m *PrismMap) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(cfloat(m.DrawOffset))
	return opts
}

type Prism struct {
	Cell int

	pos Int3
	m   *PrismMap
}

func (p *Prism) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(p.m.Sheet.SubImage(p.Cell), opts)
}

func (p *Prism) DrawOrder() (int, int) {
	return p.m.PosToWorld.Apply(p.pos).Z,
		dot(p.pos.XY(), p.m.DrawOrderBias)
}

func (p *Prism) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(cfloat(
		p.m.game.Projection.Project(p.m.PosToWorld.Apply(p.pos)),
	))
	return opts
}
