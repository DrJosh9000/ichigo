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

// Parallax is a container that changes its transform based on the position of a
// camera, intended to produce a "parallax" like effect.
type Parallax struct {
	CameraID string
	Factor   float64 // how much to change in response to the camera
	Child    interface{}

	camera *Camera
}

func (p *Parallax) Prepare(game *Game) error {
	c, ok := game.Component(p.CameraID).(*Camera)
	if !ok {
		return fmt.Errorf("component %q type != *Camera", p.CameraID)
	}
	p.camera = c
	return nil
}

func (p *Parallax) Scan() []interface{} { return []interface{}{p.Child} }

func (p *Parallax) Transform() (opts ebiten.DrawImageOptions) {
	x, y := float2(p.camera.Centre)
	opts.GeoM.Translate(x*p.Factor, y*p.Factor)
	return opts
}
