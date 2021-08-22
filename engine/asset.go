package engine

import (
	"compress/gzip"
	"encoding/gob"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// AnimDefs are easier to write as Go expressions -
	// so just set this.
	AnimDefs map[string]*AnimDef

	imageCache = make(map[string]*ebiten.Image)

	// Ensure ref types satisfy interfaces.
	_ Loader = &ImageRef{}
	_ Loader = &SceneRef{}
	_ Scener = &SceneRef{}
)

func init() {
	gob.Register(AnimRef{})
	gob.Register(ImageRef{})
	gob.Register(SceneRef{})
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
func (r *ImageRef) Load(g *Game) error {
	// Fast path load from cache
	r.image = imageCache[r.Path]
	if r.image != nil {
		return nil
	}
	// Slow path
	f, err := g.AssetFS.Open(r.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	r.image = ebiten.NewImageFromImage(i)
	imageCache[r.Path] = r.image
	return nil
}

// SceneRef loads a gzipped, gob-encoded Scene from the asset FS.
// After Load, Scene is usable.
// This is mostly useful for scenes that refer to other scenes, e.g.
//
//    var sc = &Scene{
//	    Components: []interface{}{
//			SceneRef{Path: "assets/foo.gob.gz"}
//		},
//    }
type SceneRef struct {
	Path string

	scene *Scene // not exported for gob reasons
}

// Load loads the scene from the file.
func (r *SceneRef) Load(g *Game) error {
	f, err := g.AssetFS.Open(r.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	sc := new(Scene)
	if err := gob.NewDecoder(gz).Decode(sc); err != nil {
		return err
	}
	r.scene = sc
	return nil
}

// Scene returns the loaded scene, or nil if not yet loaded.
func (r SceneRef) Scene() *Scene { return r.scene }

// The rest of the methods forward to r.scene, as such they will
// panic if the scene isn't loaded.

// Draw draws the scene.
func (r SceneRef) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	r.scene.Draw(screen, opts)
}

// DrawOrder returns the value of DrawOrder from the scene.
func (r SceneRef) DrawOrder() float64 { return r.scene.DrawOrder() }

// Ident returns the value of Ident from the scene.
func (r SceneRef) Ident() string { return r.scene.Ident() }

// Prepare prepares the scene.
func (r SceneRef) Prepare(g *Game) { r.scene.Prepare(g) }

// Scan returns the components in the scene.
func (r SceneRef) Scan() []interface{} { return r.scene.Scan() }

// Update updates the scene.
func (r SceneRef) Update() error { return r.scene.Update() }
