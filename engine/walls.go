package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	_ interface {
		Collider
		Identifier
		Scanner
		Prepper
		Transformer
	} = &Wall{}

	_ interface {
		Drawer
		Disabler
		Hider
		Scanner
		Transformer
	} = &WallUnit{}
)

// Wall is a more flexible kind of tilemap. WallUnits can be added at the same
// level as other components and are responsible for their own drawing, so that
// Game can do draw ordering, e.g. hide the player character behind a wall.
// But Wall is still responsible for collisions.
type Wall struct {
	ID
	Ersatz     bool        // disables collisions ("fake wall")
	Offset     image.Point //  offset the whole wall
	Sheet      Sheet
	UnitOffset image.Point // drawing offset
	UnitSize   image.Point // tile size
	Units      map[image.Point]*WallUnit
}

// CollidesWith implements a tilerange collosion check, similar to Tilemap.
func (w *Wall) CollidesWith(b Box) bool {
	if w.Ersatz {
		return false
	}

	// Probe the map at all tilespace coordinates overlapping the rect.
	r := b.XY().Sub(w.Offset)
	min := cdiv(r.Min, w.UnitSize)
	max := cdiv(r.Max.Sub(image.Pt(1, 1)), w.UnitSize) // NB: fencepost

	for j := min.Y; j <= max.Y; j++ {
		for i := min.X; i <= max.X; i++ {
			if w.Units[image.Pt(i, j)] != nil {
				return true
			}
		}
	}
	return false
}

// Scan returns the Sheet and all WallUnits.
func (w *Wall) Scan() []interface{} {
	c := make([]interface{}, 1, len(w.Units)+1)
	c[0] = &w.Sheet
	for _, unit := range w.Units {
		c = append(c, unit)
	}
	return c
}

// Prepare makes sure all WallUnits know about Wall and where they are, for
// drawing.
func (w *Wall) Prepare(*Game) error {
	// Ensure all child units know about wall, which houses common attributes
	for p, u := range w.Units {
		u.pos, u.wall = p, w
	}
	return nil
}

// Transform returns a GeoM translation by Offset.
func (w *Wall) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(cfloat(w.Offset))
	return opts
}

// WallUnit is a unit in a wall. Unlike a tile in a tilemap, WallUnit is
// responsible for drawing itself.
type WallUnit struct {
	Disabled
	Hidden
	Tile Tile // chooses which cell in wall.Sheet to draw
	ZOrder

	pos  image.Point // tilespace coordinates
	wall *Wall
}

// Draw draws this wall unit.
func (u *WallUnit) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(u.wall.Sheet.SubImage(u.Tile.Cell()), opts)
}

// Scan returns the Tile.
func (u *WallUnit) Scan() []interface{} { return []interface{}{u.Tile} }

func (u *WallUnit) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(cfloat(cmul(u.pos, u.wall.UnitSize).Add(u.wall.UnitOffset)))
	return opts
}
