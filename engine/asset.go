package engine

import "github.com/hajimehoshi/ebiten/v2"

type ImageRef struct {
	Path string

	image *ebiten.Image
}

func (r *ImageRef) Image() *ebiten.Image {
	if r.image == nil {
		// TODO
	}
	return r.image
}
