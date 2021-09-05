package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type PrismMap struct {
	Map           map[Point3]*Prism
	DrawOrderBias image.Point // dot with (X,Y) = bias
	DrawOffset    image.Point // offset to apply to whole map
	DrawZStride   image.Point // (X,Y) draw translation for each unit in Z
	PrismSize     Point3
	Sheet         Sheet
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

func (p *Prism) Draw(screen *ebiten.Image, tf *Transform) {
	screen.DrawImage(p.pm.Sheet.SubImage(p.Cell), &tf.Opts)
}

func (p *Prism) DrawOrder() (int, int) {
	return p.pos.Z * p.pm.PrismSize.Z,
		dot(p.pos.XY(), p.pm.DrawOrderBias)
}

func (p *Prism) Transform(pt Transform) (tf Transform) {
	v := p.pos.CMul(p.pm.PrismSize)
	tf.Opts.GeoM.Translate(cfloat(
		v.XY().Add(p.pm.DrawZStride.Mul(v.Z)),
	))
	return tf.Concat(pt)
}
