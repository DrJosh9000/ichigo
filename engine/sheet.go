package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	_ Prepper = &Sheet{}
	_ Scanner = &Sheet{}
)

// Sheet handles images that consist of a grid of equally sized regions
// (cells) and can produce subimages for the cell at an index. This is useful
// for various applications such as sprite animation and tile maps.
type Sheet struct {
	CellSize image.Point
	Src      ImageRef

	w int // width as measured in number of cells
}

func (s *Sheet) Prepare(*Game) {
	s.w, _ = s.Src.Image().Size()
	s.w /= s.CellSize.X
}

func (s *Sheet) Scan() []interface{} { return []interface{}{&s.Src} }

// SubImage returns an *ebiten.Image corresponding to the cell at the given
// index.
func (s *Sheet) SubImage(i int) *ebiten.Image {
	p := mul2(image.Pt(i%s.w, i/s.w), s.CellSize)
	r := image.Rectangle{p, p.Add(s.CellSize)}
	return s.Src.Image().SubImage(r).(*ebiten.Image)
}
