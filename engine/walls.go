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
		DrawOrderer
		Disabler
		Hider
		Prepper
		Updater
	} = &WallUnit{}
)

// Wall is a more flexible kind of tilemap. WallUnits can be added at the same
// level as other components and are responsible for their own drawing, so that
// Scene can do draw ordering, e.g. hide the player character behind a wall.
// But Wall is still responsible for collisions.
type Wall struct {
	ID
	Ersatz     bool        // fake wall - disables collisions
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

// WallUnit
type WallUnit struct {
	Disabled
	Hidden
	Pos    image.Point // tilespace coordinates
	Tile   Tile        // chooses which cell in wall.Sheet to draw
	WallID string
	ZOrder

	wall *Wall
}

func (u *WallUnit) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if u.Hidden {
		return
	}
	var geom ebiten.GeoM
	geom.Translate(float2(mul2(u.Pos, u.wall.UnitSize).Add(u.wall.UnitOffset).Add(u.wall.Offset)))
	geom.Concat(opts.GeoM)
	opts.GeoM = geom

	src := u.wall.Sheet.SubImage(u.Tile.CellIndex())
	screen.DrawImage(src, &opts)
}

func (u *WallUnit) Prepare(g *Game) {
	u.wall = g.Component(u.WallID).(*Wall)
	u.wall.regUnit(u)
}

func (u *WallUnit) Update() error {
	if u.Disabled {
		return nil
	}
	if up, ok := u.Tile.(Updater); ok {
		return up.Update()
	}
	return nil
}
