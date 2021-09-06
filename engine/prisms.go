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
	
	Map           map[Point3]*Prism  // pos -> prism
	DrawOrderBias image.Point // dot with pos.XY() = bias value
	DrawOffset    image.Point // offset applies to whole map
	DrawZStride   image.Point // draw offset for each pos unit in Z
	PrismSize     Point3      // (prismsize cmul pos) = world position
	Sheet         Sheet
}

func (m *PrismMap) CollidesWith(b Box) bool {
	// TODO
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

	pos Point3
	pm  *PrismMap
}

func (p *Prism) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(p.pm.Sheet.SubImage(p.Cell), opts)
}

func (p *Prism) DrawOrder() (int, int) {
	return p.pos.Z * p.pm.PrismSize.Z,
		dot(p.pos.XY(), p.pm.DrawOrderBias)
}

func (p *Prism) Transform(pt Transform) (tf Transform) {
	tf.Opts.GeoM.Translate(cfloat(
		cmul(p.pos.XY(), p.pm.PrismSize.XY()).
			Add(p.pm.DrawZStride.Mul(p.pos.Z)),
	))
	return tf.Concat(pt)
}
