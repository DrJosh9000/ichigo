package engine

import (
	"encoding/gob"
	"image"

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

// PrismMap
type PrismMap struct {
	ID
	Disabled
	Hidden

	Map           map[Int3]*Prism // pos -> prism
	DrawOrderBias image.Point     // dot with pos.XY() = bias value
	DrawOffset    image.Point     // offset applies to whole map
	PosToWorld    IntMatrix3x4    // p.pos -> voxelspace
	PrismSize     Int3            // in voxelspace
	Sheet         Sheet

	game *Game
}

func (m *PrismMap) CollidesWith(b Box) bool {
	// Back corner of a prism p is:
	// m.PrismPos.Apply(p.pos)

	return false
}

func (m *PrismMap) Prepare(g *Game) error {
	m.game = g
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
