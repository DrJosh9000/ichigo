package engine

import (
	"encoding/gob"
	"fmt"
	"image"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AssetFS fs.FS

	animDefCache = make(map[string]*AnimDef)
	imageCache   = make(map[string]*ebiten.Image)
)

// AnimRef
type AnimRef struct {
	Path string

	anim *Anim
}

func (r *AnimRef) Anim() (*Anim, error) {
	if r.anim != nil {
		return r.anim, nil
	}
	if ad := animDefCache[r.Path]; ad != nil {
		r.anim = &Anim{Def: ad}
		return r.anim, nil
	}
	f, err := AssetFS.Open(r.Path)
	if err != nil {
		return nil, fmt.Errorf("open asset: %w", err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	ad := &AnimDef{}
	if err := dec.Decode(ad); err != nil {
		return nil, fmt.Errorf("decode asset: %w", err)
	}
	animDefCache[r.Path] = ad
	r.anim = &Anim{Def: ad}
	return r.anim, nil
}

// ImageRef loads images from the AssetFS into *ebiten.Image form.
// It is your responsibility to import _ "image/..." for whatever
// format the files are in.
type ImageRef struct {
	Path string

	image *ebiten.Image
}

// Image returns the image. If it hasn't been loaded yet, it loads.
// Multiple distinct ImageRefs can use the same path.
func (r *ImageRef) Image() (*ebiten.Image, error) {
	if r.image != nil {
		return r.image, nil
	}
	r.image = imageCache[r.Path]
	if r.image != nil {
		return r.image, nil
	}
	f, err := AssetFS.Open(r.Path)
	if err != nil {
		return nil, fmt.Errorf("open asset: %w", err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode asset: %w", err)
	}
	r.image = ebiten.NewImageFromImage(i)
	imageCache[r.Path] = r.image
	return r.image, nil
}
