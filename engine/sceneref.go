package engine

import (
	"encoding/gob"
	"image"
	"io/fs"
	"path/filepath"
)

// Ensure SceneRef satisfies interfaces.
var _ interface {
	Loader
	Scener
} = &SceneRef{}

func init() {
	gob.Register(&SceneRef{})
}

// SceneRef loads a gzipped, gob-encoded Scene from the asset FS.
// After Load, Scene is usable.
// This is mostly useful for scenes that refer to other scenes, e.g.
//
//    sc := &Scene{
//	    Components: []interface{}{
//			&SceneRef{Path: "assets/foo.gob.gz"}  // inflated at Load time
//		},
//    }
type SceneRef struct {
	Path string

	scene *Scene // not exported for gob reasons
}

// Load loads the scene from the file.
func (r *SceneRef) Load(assets fs.FS) error {
	sc := new(Scene)
	if err := LoadGobz(sc, assets, r.Path); err != nil {
		return err
	}
	r.scene = sc
	return nil
}

// Save saves the scene to a file in the current directory.
func (r *SceneRef) Save() error { return SaveGobz(r.scene, filepath.Base(r.Path)) }

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

// IsHidden returns the value of IsHidden from the scene.
func (r SceneRef) IsHidden() bool { return r.scene.IsHidden() }

// Hide calls Hide on the scene.
func (r SceneRef) Hide() { r.scene.Hide() }

// Show calls Show on the scene.
func (r SceneRef) Show() { r.scene.Show() }

// Ident returns the value of Ident from the scene.
func (r SceneRef) Ident() string { return r.scene.Ident() }

// Scan returns the components in the scene.
func (r SceneRef) Scan() []interface{} { return r.scene.Scan() }
