package engine

import (
	"compress/gzip"
	"encoding/gob"
	"image"
	"io/fs"
	"os"

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

type assetKey struct {
	assets fs.FS
	path   string
}

func loadGobz(dst interface{}, assets fs.FS, path string) error {
	f, err := assets.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	return gob.NewDecoder(gz).Decode(dst)
}

// saveGobz takes an object, gob-encodes it, gzips it, and writes to disk.
func saveGobz(src interface{}, name string) error {
	f, err := os.CreateTemp(".", name)
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	gz := gzip.NewWriter(f)
	if err := gob.NewEncoder(gz).Encode(src); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), name)
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
