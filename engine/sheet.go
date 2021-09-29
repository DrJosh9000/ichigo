package engine

import (
	"image"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ interface {
	Prepper
	Scanner
} = &Sheet{}

// Sheet handles images that consist of a grid of equally sized regions
// (cells) and can produce subimages for the cell at an index. This is useful
// for various applications such as sprite animation and tile maps. Additionally
// each sheet carries a collection of animations that use the sheet.
type Sheet struct {
	AnimDefs map[string]*AnimDef
	CellSize image.Point
	Src      ImageRef

	w int // width as measured in number of cells
}

// NewAnim returns a new Anim for the given key, or nil if not found in
// AnimDefs.
func (s *Sheet) NewAnim(key string) *Anim {
	return s.AnimDefs[key].NewAnim()
}

// NewAnims returns a new Anim for every AnimDef in the AnimDefs map.
func (s *Sheet) NewAnims() map[string]*Anim {
	m := make(map[string]*Anim, len(s.AnimDefs))
	for k, d := range s.AnimDefs {
		m[k] = d.NewAnim()
	}
	return m
}

// Prepare computes the width of the image (in cells).
func (s *Sheet) Prepare(*Game) error {
	s.w, _ = s.Src.Image().Size()
	s.w /= s.CellSize.X
	return nil
}

// Scan visits &s.Src.
func (s *Sheet) Scan(visit VisitFunc) error {
	return visit(&s.Src)
}

// SubImage returns an *ebiten.Image corresponding to the given cell index.
func (s *Sheet) SubImage(i int) *ebiten.Image {
	p := geom.CMul(image.Pt(i%s.w, i/s.w), s.CellSize)
	r := image.Rectangle{p, p.Add(s.CellSize)}
	return s.Src.Image().SubImage(r).(*ebiten.Image)
}

func (s *Sheet) String() string { return "Sheet" }
