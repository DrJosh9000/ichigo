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

type PrismMap struct {
	ID
	Disabled
	Hidden

	Map           map[Int3]*Prism // pos -> prism
	DrawOrderBias image.Point     // dot with pos.XY() = bias value
	DrawOffset    image.Point     // offset applies to whole map
	PosToDraw     IntMatrix2x3    // p.pos -> drawspace (before offset and camera and ...)
	PosToWorld    IntMatrix3x4    // p.pos -> worldspace
	PrismSize     Int3            // in worldspace
	Sheet         Sheet
}

func (m *PrismMap) CollidesWith(b Box) bool {
	// Back corner of a prism p is:
	// m.PrismPos.Apply(p.pos)

	return false
}

func (m *PrismMap) Prepare(*Game) error {
	for v, p := range m.Map {
		p.pos = v
		p.pm = m
	}
	return nil
}

func (m *PrismMap) Transform(pt Transform) (tf Transform) {
	tf.Opts.GeoM.Translate(cfloat(m.DrawOffset))
	return tf.Concat(pt)
}

type Prism struct {
	Cell int

	pos Int3
	pm  *PrismMap
}

func (p *Prism) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(p.pm.Sheet.SubImage(p.Cell), opts)
}

func (p *Prism) DrawOrder() (int, int) {
	return p.pm.PosToWorld.Apply(p.pos).Z,
		dot(p.pos.XY(), p.pm.DrawOrderBias)
}

func (p *Prism) Transform(pt Transform) (tf Transform) {
	tf.Opts.GeoM.Translate(cfloat(
		p.pm.PosToDraw.Apply(p.pos),
	))
	return tf.Concat(pt)
}
