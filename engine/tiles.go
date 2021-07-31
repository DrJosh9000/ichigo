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
	Map       [][]Tile
	Src       ImageRef // must be a horizontal tile set
	TileSize  int
	Transform ebiten.GeoM
	ZPos
}

// Draw draws the tilemap.
func (t *Tilemap) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	geom.Concat(t.Transform)
	for j, row := range t.Map {
		for i, tile := range row {
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(i*t.TileSize), float64(j*t.TileSize))
			op.GeoM.Concat(geom)

			sx := tile.TileIndex() * t.TileSize
			src := t.Src.Image().SubImage(image.Rect(sx, 0, sx+t.TileSize, t.TileSize)).(*ebiten.Image)
			screen.DrawImage(src, &op)
		}
	}
}

// Update calls Update on any tiles that are Updaters, e.g. AnimatedTile.
func (t *Tilemap) Update() error {
	for j := range t.Map {
		for i := range t.Map[j] {
			if tile, ok := t.Map[j][i].(Updater); ok {
				if err := tile.Update(); err != nil {
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
