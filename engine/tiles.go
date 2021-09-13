package engine

import (
	"encoding/gob"
	"fmt"
	"image"
	"io/fs"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Tilemap satisfies interfaces.
var _ interface {
	Identifier
	Collider
	Drawer
	Hider
	Scanner
	Transformer
} = &Tilemap{}

// Ensure StaticTile and AnimatedTile satisfy Tile.
var (
	_ Tile = StaticTile(0)
	_ interface {
		Tile
		Scanner
	} = &AnimatedTile{}
)

func init() {
	gob.Register(&AnimatedTile{})
	gob.Register(StaticTile(0))
	gob.Register(&Tilemap{})
}

// Tilemap renders a grid of rectangular tiles at equal Z position.
type Tilemap struct {
	ID
	Disables
	Hides
	Map    map[image.Point]Tile // tilespace coordinate -> tile
	Ersatz bool                 // disables collisions ("fake wall")
	Offset image.Point          // world coordinates
	Sheet  Sheet
	ZPosition
}

// CollidesWith implements Collider.
func (t *Tilemap) CollidesWith(b geom.Box) bool {
	if t.Ersatz {
		return false
	}

	// Probe the map at all tilespace coordinates overlapping the rect.
	r := b.XY().Sub(t.Offset) // TODO: pretend tilemap is a plane in 3D?
	min := geom.CDiv(r.Min, t.Sheet.CellSize)
	max := geom.CDiv(r.Max.Sub(image.Pt(1, 1)), t.Sheet.CellSize) // NB: fencepost

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
	for p, tile := range t.Map {
		if tile == nil {
			continue
		}
		var mat ebiten.GeoM
		mat.Translate(geom.CFloat(geom.CMul(p, t.Sheet.CellSize)))
		mat.Concat(og)
		opts.GeoM = mat

		src := t.Sheet.SubImage(tile.Cell())
		screen.DrawImage(src, opts)
	}
}

// Load instantiates animations for all AnimatedTiles.
func (t *Tilemap) Load(fs.FS) error {
	for _, tile := range t.Map {
		at, ok := tile.(*AnimatedTile)
		if !ok {
			continue
		}
		at.anim = t.Sheet.NewAnim(at.AnimKey)
		if at.anim == nil {
			return fmt.Errorf("missing anim %q", at.AnimKey)
		}
	}
	return nil
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

// Transform returns a translation by t.Offset.
func (t *Tilemap) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(geom.CFloat(t.Offset))
	return opts
}

// TileAt returns the tile present at the given world coordinate.
func (t *Tilemap) TileAt(wc image.Point) Tile {
	return t.Map[geom.CDiv(wc.Sub(t.Offset), t.Sheet.CellSize)]
}

// SetTileAt sets the tile at the given world coordinate.
func (t *Tilemap) SetTileAt(wc image.Point, tile Tile) {
	t.Map[geom.CDiv(wc.Sub(t.Offset), t.Sheet.CellSize)] = tile
}

// TileBounds returns a rectangle describing the tile boundary for the tile
// at the given world coordinate.
func (t *Tilemap) TileBounds(wc image.Point) image.Rectangle {
	p := geom.CMul(geom.CDiv(wc.Sub(t.Offset), t.Sheet.CellSize), t.Sheet.CellSize).Add(t.Offset)
	return image.Rectangle{p, p.Add(t.Sheet.CellSize)}
}

// Tile is the interface needed by Tilemap.
type Tile interface {
	Cell() int
}

// StaticTile returns a fixed tile index.
type StaticTile int

func (s StaticTile) Cell() int { return int(s) }

// AnimatedTile uses an Anim to choose a tile index.
type AnimatedTile struct {
	AnimKey string

	anim *Anim
}

func (a *AnimatedTile) Cell() int { return a.anim.Cell() }

// Scan returns a.anim.
func (a *AnimatedTile) Scan() []interface{} { return []interface{}{a.anim} }
