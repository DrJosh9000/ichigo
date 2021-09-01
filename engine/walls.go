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

// Prepare makes sure all WallUnits know about Wall.
func (w *Wall) Prepare(*Game) error {
	// Ensure all child units know about wall, which houses common attributes
	for _, u := range w.Units {
		u.wall = w
	}
	return nil
}

// Transform returns a GeoM translation by Offset.
func (w *Wall) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(float2(w.Offset))
	return opts
}

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

// Draw draws this wall unit.
func (u *WallUnit) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(u.wall.Sheet.SubImage(u.Tile.CellIndex()), opts)
}

// Scan returns the Tile.
func (u *WallUnit) Scan() []interface{} { return []interface{}{u.Tile} }

func (u *WallUnit) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(float2(mul2(u.Pos, u.wall.UnitSize).Add(u.wall.UnitOffset)))
	return opts
}
