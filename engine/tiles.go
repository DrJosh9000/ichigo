package engine

import (
	"encoding/gob"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(AnimatedTile{})
	gob.Register(StaticTile(0))
	gob.Register(Tilemap{})
}

type Tilemap struct {
	Map       [][]Tile
	Src       ImageRef // must be a horizontal tile set
	TileSize  int
	Transform ebiten.GeoM
	ZPos
}

func (t *Tilemap) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	geom.Concat(t.Transform)
	for j, row := range t.Map {
		for i, tile := range row {
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(i*t.TileSize), float64(j*t.TileSize))
			op.GeoM.Concat(geom)

			sx := tile.TileIndex() * t.TileSize
			im, err := t.Src.Image()
			if err != nil {
				log.Fatalf("Loading image from reference: %v", err)
			}
			src := im.SubImage(image.Rect(sx, 0, sx+t.TileSize, t.TileSize)).(*ebiten.Image)
			screen.DrawImage(src, &op)
		}
	}
}

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

type Tile interface {
	TileIndex() int
}

type StaticTile int

func (s StaticTile) TileIndex() int { return int(s) }

type AnimatedTile struct {
	Anim
}

func (a *AnimatedTile) TileIndex() int { return a.CurrentFrame() }
