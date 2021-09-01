package engine

import (
	"encoding/gob"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ interface {
	Prepper
	Scanner
	Transformer
} = &Parallax{}

func init() {
	gob.Register(&Parallax{})
}

// Parallax is a container that translates based on the position of a
// camera, intended to produce a "parallax" like effect.
type Parallax struct {
	CameraID string
	Factor   float64 // how much to translate in response to the camera
	Child    interface{}

	camera *Camera
}

// Prepare obtains a reference to the camera.
func (p *Parallax) Prepare(game *Game) error {
	c, ok := game.Component(p.CameraID).(*Camera)
	if !ok {
		return fmt.Errorf("component %q type != *Camera", p.CameraID)
	}
	p.camera = c
	return nil
}

// Scan returns the child component.
func (p *Parallax) Scan() []interface{} { return []interface{}{p.Child} }

// Transform returns a GeoM translation of Factor * camera.Centre.
func (p *Parallax) Transform() (opts ebiten.DrawImageOptions) {
	x, y := float2(p.camera.Centre)
	opts.GeoM.Translate(x*p.Factor, y*p.Factor)
	return opts
}
