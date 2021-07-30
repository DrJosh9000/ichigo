package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tilemap struct {
	Map       [][]Tile
	Src       *ebiten.Image // must be a horizontal tile set
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
			src := t.Src.SubImage(image.Rect(sx, 0, sx+t.TileSize, t.TileSize)).(*ebiten.Image)
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
	Frame         int // index into AnimDef
	DurationTicks int // time spent showing current frame
	AnimDef       []TileAnimFrameDef
}

func (a *AnimatedTile) TileIndex() int { return a.AnimDef[a.Frame].Tile }

func (a *AnimatedTile) Update() error {
	a.DurationTicks++
	if a.DurationTicks >= a.AnimDef[a.Frame].DurationTicks {
		a.DurationTicks = 0
		a.Frame++
	}
	if a.Frame >= len(a.AnimDef) {
		a.Frame = 0
	}
	return nil
}

type TileAnimFrameDef struct {
	Tile          int // show this tile
	DurationTicks int // show it for this long
}
