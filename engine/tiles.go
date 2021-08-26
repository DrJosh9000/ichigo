package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Tilemap satisfies interfaces.
var (
	_ Identifier  = &Tilemap{}
	_ Collider    = &Tilemap{}
	_ Drawer      = &Tilemap{}
	_ DrawOrderer = &Tilemap{}
	_ Hider       = &Tilemap{}
	_ Scanner     = &Tilemap{}
	_ Updater     = &Tilemap{}
)

func init() {
	gob.Register(&AnimatedTile{})
	gob.Register(StaticTile(0))
	gob.Register(&Tilemap{})
}

// Tilemap renders a grid of square tiles.
type Tilemap struct {
	ID
	Disabled
	Hidden
	Map      map[image.Point]Tile // tilespace coordinate -> tile
	Ersatz   bool                 // "fake wall"
	Offset   image.Point          // world coordinates
	Src      ImageRef
	TileSize int
	ZOrder
}

// CollidesWith implements Collider.
func (t *Tilemap) CollidesWith(r image.Rectangle) bool {
	if t.Ersatz {
		return false
	}

	// Probe the map at all tilespace coordinates overlapping the rect.
	r = r.Sub(t.Offset)
	min := r.Min.Div(t.TileSize)
	max := r.Max.Sub(image.Pt(1, 1)).Div(t.TileSize) // NB: fencepost

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
func (t *Tilemap) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	if t.Hidden {
		return
	}
	src := t.Src.Image()
	w, _ := src.Size()
	og := opts.GeoM
	var geom ebiten.GeoM
	for p, tile := range t.Map {
		if tile == nil {
			continue
		}
		geom.Reset()
		geom.Translate(float64(p.X*t.TileSize+t.Offset.X), float64(p.Y*t.TileSize+t.Offset.Y))
		geom.Concat(og)
		opts.GeoM = geom

		s := tile.TileIndex() * t.TileSize
		sx, sy := s%w, (s/w)*t.TileSize
		src := src.SubImage(image.Rect(sx, sy, sx+t.TileSize, sy+t.TileSize)).(*ebiten.Image)
		screen.DrawImage(src, &opts)
	}
}

// Scan returns a slice containing Src and all non-nil tiles.
func (t *Tilemap) Scan() []interface{} {
	c := make([]interface{}, 1, len(t.Map)+1)
	c[0] = &t.Src
	for _, tile := range t.Map {
		c = append(c, tile)
	}
	return c
}

// Update calls Update on any tiles that are Updaters, e.g. AnimatedTile.
func (t *Tilemap) Update() error {
	if t.Disabled {
		return nil
	}
	for _, tile := range t.Map {
		if u, ok := tile.(Updater); ok {
			if err := u.Update(); err != nil {
				return err
			}
		}
	}
	return nil
}

// TileAt returns the tile present at the given world coordinate.
func (t *Tilemap) TileAt(wc image.Point) Tile {
	return t.Map[wc.Sub(t.Offset).Div(t.TileSize)]
}

// SetTileAt sets the tile at the given world coordinate.
func (t *Tilemap) SetTileAt(wc image.Point, tile Tile) {
	t.Map[wc.Sub(t.Offset).Div(t.TileSize)] = tile
}

// TileBounds returns a rectangle describing the tile boundary for the tile
// at the given world coordinate.
func (t *Tilemap) TileBounds(wc image.Point) image.Rectangle {
	p := wc.Sub(t.Offset).Div(t.TileSize).Mul(t.TileSize).Add(t.Offset)
	return image.Rectangle{p, p.Add(image.Pt(t.TileSize, t.TileSize))}
}

// Tile is the interface needed by Tilemap.
type Tile interface {
	TileIndex() int
}

// Ensure StaticTile and AnimatedTile satisfy Tile.
var (
	_ Tile    = StaticTile(0)
	_ Tile    = AnimatedTile{}
	_ Scanner = AnimatedTile{}
)

// StaticTile returns a fixed tile index.
type StaticTile int

func (s StaticTile) TileIndex() int { return int(s) }

// AnimatedTile uses an Anim to choose a tile index.
type AnimatedTile struct {
	Animer
}

func (a AnimatedTile) TileIndex() int { return a.CurrentFrame() }

// Scan returns a.Animer. (It could be a Loader.)
func (a AnimatedTile) Scan() []interface{} { return []interface{}{a.Animer} }
