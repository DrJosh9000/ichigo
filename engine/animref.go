package engine

import (
	"encoding/gob"
	"io/fs"
)

// TODO: tidy this crap up
// Anims probably belong with Sheet

var (
	animCache = make(map[assetKey]Anim)

	_ interface {
		Animer
		Loader
	} = &AnimRef{}
)

func init() {
	gob.Register(&AnimRef{})
}

// AnimRef manages an Anim using a premade AnimDef from the cache.
type AnimRef struct {
	Path string

	anim Anim
}

func (r *AnimRef) Load(assets fs.FS) error {
	// Fast path: set r.anim to a copy
	anim, found := animCache[assetKey{assets, r.Path}]
	if found {
		r.anim = anim
		return nil
	}
	// Slow path: load from gobz file
	if err := LoadGobz(&r.anim, assets, r.Path); err != nil {
		return err
	}
	animCache[assetKey{assets, r.Path}] = r.anim
	return nil
}

// CurrentFrame returns the value of CurrentFrame from r.anim.
func (r *AnimRef) CurrentFrame() int { return r.anim.CurrentFrame() }

// Reset calls Reset on r.anim.
func (r *AnimRef) Reset() { r.anim.Reset() }

// Update calls Update on r.anim.
func (r *AnimRef) Update() error { return r.anim.Update() }
