package engine

import (
	"encoding/gob"
	"image"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AssetFS fs.FS

	animDefCache = make(map[string]*AnimDef)
	imageCache   = make(map[string]*ebiten.Image)
)

// AnimRef loads AnimDef from an asset and manages an Anim using it.
type AnimRef struct {
	Path string

	anim *Anim
}

func (r *AnimRef) Anim() *Anim {
	if r.anim != nil {
		return r.anim
	}
	if ad := animDefCache[r.Path]; ad != nil {
		r.anim = &Anim{Def: ad}
		return r.anim
	}
	f, err := AssetFS.Open(r.Path)
	if err != nil {
		log.Fatalf("Couldn't open asset: %v", err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	ad := &AnimDef{}
	if err := dec.Decode(ad); err != nil {
		log.Fatalf("Couldn't decode asset: %v", err)
	}
	animDefCache[r.Path] = ad
	r.anim = &Anim{Def: ad}
	return r.anim
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
func (r *ImageRef) Image() *ebiten.Image {
	if r.image != nil {
		return r.image
	}
	r.image = imageCache[r.Path]
	if r.image != nil {
		return r.image
	}
	f, err := AssetFS.Open(r.Path)
	if err != nil {
		log.Fatalf("Couldn't open asset: %v", err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("Couldn't decode asset: %v", err)
	}
	r.image = ebiten.NewImageFromImage(i)
	imageCache[r.Path] = r.image
	return r.image
}
