package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tilemap struct {
	Map      [][]int
	Src      *ebiten.Image // must be a horizontal tile set
	TileSize int
	GeoM     ebiten.GeoM

	ZPos
}

func (t *Tilemap) Draw(screen *ebiten.Image) {
	for j, row := range t.Map {
		for i, tile := range row {
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(i*t.TileSize), float64(j*t.TileSize))
			op.GeoM.Concat(t.GeoM)

			sx := tile * t.TileSize
			src := t.Src.SubImage(image.Rect(sx, 0, sx+t.TileSize, t.TileSize)).(*ebiten.Image)
			screen.DrawImage(src, &op)
		}
	}
}
