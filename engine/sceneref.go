package engine

import (
	"encoding/gob"
	"image"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure SceneRef satisfies interfaces.
var (
	_ Loader = &SceneRef{}
	_ Scener = &SceneRef{}
)

func init() {
	gob.Register(SceneRef{})
}

// SceneRef loads a gzipped, gob-encoded Scene from the asset FS.
// After Load, Scene is usable.
// This is mostly useful for scenes that refer to other scenes, e.g.
//
//    var sc = &Scene{
//	    Components: []interface{}{
//			SceneRef{Path: "assets/foo.gob.gz"}  // inflated at Load time
//		},
//    }
type SceneRef struct {
	Path string

	scene *Scene // not exported for gob reasons
}

// Load loads the scene from the file.
func (r *SceneRef) Load(assets fs.FS) error {
	sc := new(Scene)
	if err := loadGobz(sc, assets, r.Path); err != nil {
		return err
	}
	r.scene = sc
	return nil
}

// The rest of the methods forward to r.scene, as such they will
// panic if the scene isn't loaded.

// BoundingRect returns the Bounds from the scene.
func (r SceneRef) BoundingRect() image.Rectangle { return r.scene.BoundingRect() }

// IsDisabled returns the value of IsDisabled from the scene.
func (r SceneRef) IsDisabled() bool { return r.scene.IsDisabled() }

// Disable calls Disable on the scene.
func (r SceneRef) Disable() { r.scene.Disable() }

// Enable calls Enable on the scene.
func (r SceneRef) Enable() { r.scene.Enable() }

// Draw draws the scene.
func (r SceneRef) Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	r.scene.Draw(screen, opts)
}

// DrawOrder returns the value of DrawOrder from the scene.
func (r SceneRef) DrawOrder() float64 { return r.scene.DrawOrder() }

// IsHidden returns the value of IsHidden from the scene.
func (r SceneRef) IsHidden() bool { return r.scene.IsHidden() }

// Hide calls Hide on the scene.
func (r SceneRef) Hide() { r.scene.Hide() }

// Show calls Show on the scene.
func (r SceneRef) Show() { r.scene.Show() }

// Ident returns the value of Ident from the scene.
func (r SceneRef) Ident() string { return r.scene.Ident() }

// Prepare prepares the scene.
func (r SceneRef) Prepare(g *Game) { r.scene.Prepare(g) }

// Scan returns the components in the scene.
func (r SceneRef) Scan() []interface{} { return r.scene.Scan() }

// Update updates the scene.
func (r SceneRef) Update() error { return r.scene.Update() }
