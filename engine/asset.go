package engine

import (
	"encoding/gob"
	"image"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// AnimDefs are easier to write as Go expressions -
	// so just set this.
	// TODO: put in Game
	AnimDefs map[string]*AnimDef

	imageCache = make(map[assetKey]*ebiten.Image)

	// Ensure ImageRef satisfies interfaces.
	_ Loader = &ImageRef{}
)

func init() {
	gob.Register(AnimRef{})
	gob.Register(ImageRef{})
}

type assetKey struct {
	assets fs.FS
	path   string
}

// AnimRef manages an Anim using a premade AnimDef from the cache.
type AnimRef struct {
	Key string

	anim *Anim
}

func (r *AnimRef) Anim() *Anim {
	if r.anim != nil {
		return r.anim
	}
	ad := AnimDefs[r.Key]
	if ad == nil {
		log.Fatalf("Unknown AnimDef %q", r.Key)
		return nil
	}
	r.anim = &Anim{Def: ad}
	return r.anim
}

// ImageRef loads images from the AssetFS into *ebiten.Image form.
// It is your responsibility to import _ "image/..." for whatever
// format the files are in, and to load it (either return it as a
// subcomponent from Scan so that Game will Load it, or call Load
// yourself).
type ImageRef struct {
	Path string

	image *ebiten.Image
}

// Image returns the image, or nil if not loaded.
// Multiple distinct ImageRefs can use the same path.
func (r *ImageRef) Image() *ebiten.Image {
	return r.image
}

// Load loads the image. Load is required before Image returns.
// Loading the same path multiple times uses a cache to return
// the same image.
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
