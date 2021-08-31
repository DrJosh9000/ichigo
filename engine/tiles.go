package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Tilemap satisfies interfaces.
var _ interface {
	Identifier
	Collider
	Drawer
	Hider
	Scanner
} = &Tilemap{}

// Ensure StaticTile and AnimatedTile satisfy Tile.
var (
	_ Tile = StaticTile(0)
	_ interface {
		Tile
		Scanner
	} = AnimatedTile{}
)

func init() {
	gob.Register(&AnimatedTile{})
	gob.Register(StaticTile(0))
	gob.Register(&Tilemap{})
}

// Tilemap renders a grid of rectangular tiles at equal Z position.
type Tilemap struct {
	ID
	Disabled
	Hidden
	Map    map[image.Point]Tile // tilespace coordinate -> tile
	Ersatz bool                 // disables collisions ("fake wall")
	Offset image.Point          // world coordinates
	Sheet  Sheet
	ZOrder
}

// CollidesWith implements Collider.
func (t *Tilemap) CollidesWith(r image.Rectangle) bool {
	if t.Ersatz {
		return false
	}

	// Probe the map at all tilespace coordinates overlapping the rect.
	r = r.Sub(t.Offset)
	min := div2(r.Min, t.Sheet.CellSize)
	max := div2(r.Max.Sub(image.Pt(1, 1)), t.Sheet.CellSize) // NB: fencepost

	for j := min.Y; j <= max.Y; j++ {
		for i := min.X; i <= max.X; i++ {
			if t.Map[image.Pt(i, j)] != nil {
				return true
			}
		}
	}
	return false
}

// Draw draws the tilemap.
func (t *Tilemap) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	og := opts.GeoM
	var geom ebiten.GeoM
	for p, tile := range t.Map {
		if tile == nil {
			continue
		}
		geom.Reset()
		geom.Translate(float2(mul2(p, t.Sheet.CellSize).Add(t.Offset)))
		geom.Concat(og)
		opts.GeoM = geom

		src := t.Sheet.SubImage(tile.CellIndex())
		screen.DrawImage(src, opts)
	}
}

// Scan returns a slice containing Src and all non-nil tiles.
func (t *Tilemap) Scan() []interface{} {
	c := make([]interface{}, 1, len(t.Map)+1)
	c[0] = &t.Sheet
	for _, tile := range t.Map {
		c = append(c, tile)
	}
	return c
}

// TileAt returns the tile present at the given world coordinate.
func (t *Tilemap) TileAt(wc image.Point) Tile {
	return t.Map[div2(wc.Sub(t.Offset), t.Sheet.CellSize)]
}

// SetTileAt sets the tile at the given world coordinate.
func (t *Tilemap) SetTileAt(wc image.Point, tile Tile) {
	t.Map[div2(wc.Sub(t.Offset), t.Sheet.CellSize)] = tile
}

// TileBounds returns a rectangle describing the tile boundary for the tile
// at the given world coordinate.
func (t *Tilemap) TileBounds(wc image.Point) image.Rectangle {
	p := mul2(div2(wc.Sub(t.Offset), t.Sheet.CellSize), t.Sheet.CellSize).Add(t.Offset)
	return image.Rectangle{p, p.Add(t.Sheet.CellSize)}
}

// Tile is the interface needed by Tilemap.
type Tile interface {
	CellIndex() int
}

// StaticTile returns a fixed tile index.
type StaticTile int

func (s StaticTile) CellIndex() int { return int(s) }

// AnimatedTile uses an Anim to choose a tile index.
type AnimatedTile struct {
	Animer
}

func (a AnimatedTile) CellIndex() int { return a.CurrentFrame() }

// Scan returns a.Animer. (It could be a Loader.)
func (a AnimatedTile) Scan() []interface{} { return []interface{}{a.Animer} }
