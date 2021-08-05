package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(AnimatedTile{})
	gob.Register(StaticTile(0))
	gob.Register(Tilemap{})
}

// Tilemap renders a grid of tiles.
type Tilemap struct {
	Disabled bool
	Hidden   bool
	ID
	Map       [][]Tile
	Ersatz    bool // "fake wall"
	Src       ImageRef
	TileSize  int
	Transform GeoMDef
	ZPos
}

// CollidesWith implements Collider.
func (t *Tilemap) CollidesWith(r image.Rectangle) bool {
	if t.Ersatz {
		return false
	}
	// TODO: optimise?
	for j, row := range t.Map {
		for i, tile := range row {
			if tile == nil {
				continue
			}
			if r.Overlaps(image.Rect(i*t.TileSize, j*t.TileSize, (i+1)*t.TileSize, (j+1)*t.TileSize)) {
				return true
			}
		}
	}
	return false
}

// Draw draws the tilemap.
func (t *Tilemap) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	if t.Hidden {
		return
	}
	src := t.Src.Image()
	w, _ := src.Size()
	geom.Concat(*t.Transform.GeoM())
	for j, row := range t.Map {
		for i, tile := range row {
			if tile == nil {
				continue
			}
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(i*t.TileSize), float64(j*t.TileSize))
			op.GeoM.Concat(geom)

			s := tile.TileIndex() * t.TileSize
			sx, sy := s%w, (s/w)*t.TileSize
			src := src.SubImage(image.Rect(sx, sy, sx+t.TileSize, sy+t.TileSize)).(*ebiten.Image)
			screen.DrawImage(src, &op)
		}
	}
}

// Update calls Update on any tiles that are Updaters, e.g. AnimatedTile.
func (t *Tilemap) Update() error {
	if t.Disabled {
		return nil
	}
	for _, row := range t.Map {
		for _, tile := range row {
			if u, ok := tile.(Updater); ok {
				if err := u.Update(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Tile is the interface needed by Tilemap.
type Tile interface {
	TileIndex() int
}

// StaticTile returns a fixed tile index.
type StaticTile int

func (s StaticTile) TileIndex() int { return int(s) }

// AnimatedTile uses an Anim to choose a tile index.
type AnimatedTile struct {
	AnimRef
}

func (a *AnimatedTile) TileIndex() int { return a.Anim().CurrentFrame() }

func (a *AnimatedTile) Update() error { return a.Anim().Update() }
