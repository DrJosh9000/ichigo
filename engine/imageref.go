package engine

import (
	"encoding/gob"
	"image"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	imageCache = make(map[assetKey]*ebiten.Image)

	// Ensure types satisfy interfaces.
	_ Loader = &ImageRef{}
)

func init() {
	gob.Register(ImageRef{})
}

// ImageRef loads images from the AssetFS into *ebiten.Image form. It is your
// responsibility to import _ "image/..." for whatever format the files are in,
// and to load it (either return it as a subcomponent from Scan so that Game
// will Load it, or call Load yourself).
type ImageRef struct {
	Path string

	image *ebiten.Image
}

// Image returns the image, or nil if not loaded. Multiple distinct ImageRefs
// can use the same path efficiently.
func (r *ImageRef) Image() *ebiten.Image {
	return r.image
}

// Load loads the image. Load is required before Image returns. Loading the same
// path multiple times uses a cache to return the same image.
func (r *ImageRef) Load(assets fs.FS) error {
	// Fast path load from cache
	r.image = imageCache[assetKey{assets, r.Path}]
	if r.image != nil {
		return nil
	}
	// Slow path
	f, err := assets.Open(r.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	r.image = ebiten.NewImageFromImage(i)
	imageCache[assetKey{assets, r.Path}] = r.image
	return nil
}
