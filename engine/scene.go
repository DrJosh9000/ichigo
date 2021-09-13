package engine

import (
	"encoding/gob"
	"io/fs"
	"path/filepath"
)

var (
	_ scener = &Scene{}

	_ interface {
		Loader
		Saver
		scener
	} = &SceneRef{}
)

type scener interface {
	BoundingRecter
	Disabler
	Hider
	Identifier
	Scanner
}

func init() {
	gob.Register(&Scene{})
	gob.Register(&SceneRef{})
}

// Scene contains a bunch of components.
type Scene struct {
	ID
	Bounds     // world coordinates
	Components []interface{}
	Disables
	Hides
}

// Scan returns all immediate subcomponents (including the camera, if not nil).
func (s *Scene) Scan() []interface{} { return s.Components }

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

	*Scene // not gob encoded
}

// GobDecode saves the byte slice as Path.
func (r *SceneRef) GobDecode(b []byte) error {
	r.Path = string(b)
	return nil
}

// GobEncode returns Path as a byte slice.
func (r *SceneRef) GobEncode() ([]byte, error) {
	return []byte(r.Path), nil
}

// Load loads the scene from the file.
func (r *SceneRef) Load(assets fs.FS) error {
	sc := new(Scene)
	if err := LoadGobz(sc, assets, r.Path); err != nil {
		return err
	}
	r.Scene = sc
	return nil
}

// Save saves the scene to a file in the current directory.
func (r *SceneRef) Save() error {
	return SaveGobz(r.Scene, filepath.Base(r.Path))
}
