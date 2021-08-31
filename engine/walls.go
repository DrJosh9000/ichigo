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
	} = &Wall{}

	_ interface {
		Drawer
		Disabler
		Hider
		Prepper
	} = &WallUnit{}
)

// Wall is a more flexible kind of tilemap. WallUnits can be added at the same
// level as other components and are responsible for their own drawing, so that
// Scene can do draw ordering, e.g. hide the player character behind a wall.
// But Wall is still responsible for collisions.
type Wall struct {
	ID
	Ersatz     bool        // disables collisions ("fake wall")
	Offset     image.Point //  offset the whole wall
	Sheet      Sheet
	UnitOffset image.Point // drawing offset
	UnitSize   image.Point // tile size

	units map[image.Point]*WallUnit
}

func (w *Wall) regUnit(u *WallUnit) {
	if w.units == nil {
		w.units = make(map[image.Point]*WallUnit)
	}
	w.units[u.Pos] = u
}

func (w *Wall) CollidesWith(r image.Rectangle) bool {
	if w.Ersatz {
		return false
	}

	// Probe the map at all tilespace coordinates overlapping the rect.
	r = r.Sub(w.Offset)
	min := div2(r.Min, w.UnitSize)
	max := div2(r.Max.Sub(image.Pt(1, 1)), w.UnitSize) // NB: fencepost

	for j := min.Y; j <= max.Y; j++ {
		for i := min.X; i <= max.X; i++ {
			if w.units[image.Pt(i, j)] != nil {
				return true
			}
		}
	}
	return false
}

func (w *Wall) Scan() []interface{} { return []interface{}{&w.Sheet} }

// WallUnit is a unit in a wall. Unlike a tile in a tilemap, WallUnit is
// responsible for drawing itself.
type WallUnit struct {
	Disabled
	Hidden
	Pos    image.Point // tilespace coordinates
	Tile   Tile        // chooses which cell in wall.Sheet to draw
	WallID string
	ZOrder

	wall *Wall
}

func (u *WallUnit) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(u.wall.Sheet.SubImage(u.Tile.CellIndex()), opts)
}

func (u *WallUnit) Prepare(g *Game) error {
	u.wall = g.Component(u.WallID).(*Wall)
	u.wall.regUnit(u)
	return nil
}

func (u *WallUnit) Scan() []interface{} { return []interface{}{u.Tile} }

func (u *WallUnit) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(float2(mul2(u.Pos, u.wall.UnitSize).Add(u.wall.UnitOffset).Add(u.wall.Offset)))
	return opts
}
