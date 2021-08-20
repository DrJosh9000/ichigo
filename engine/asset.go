package engine

import (
	"compress/gzip"
	"encoding/gob"
	"image"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// Assets (usually embed.FS)
	AssetFS fs.FS

	// AnimDefs are easier to write as Go expressions -
	// so just set this.
	AnimDefs map[string]*AnimDef

	imageCache = make(map[string]*ebiten.Image)
)

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

// SceneRef loads a gzipped, gob-encoded Scene from the asset FS.
type SceneRef struct {
	Path string

	scene *Scene
}

func (r *SceneRef) Scene() *Scene {
	if r.scene != nil {
		return r.scene
	}
	f, err := AssetFS.Open(r.Path)
	if err != nil {
		log.Fatalf("Couldn't open asset: %v", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		log.Fatalf("Couldn't gunzip asset: %v", err)
	}
	r.scene = new(Scene)
	if err := gob.NewDecoder(gz).Decode(r.scene); err != nil {
		log.Fatalf("Couldn't decode asset: %v", err)
	}
	return r.scene
}
