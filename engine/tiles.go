package engine

import "github.com/hajimehoshi/ebiten/v2"

type Tilemap struct {
	Map      [][]int
	Src      *ebiten.Image
	TileSize int
	GeoM     *ebiten.GeoM

	ZPos
}

func (t *Tilemap) Draw(screen *ebiten.Image) {

}
