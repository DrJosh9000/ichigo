/*
Copyright 2021 Josh Deprez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package engine

import (
	"encoding/gob"
	"fmt"

	"github.com/DrJosh9000/ichigo/geom"
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

// Scan visits p.Child.
func (p *Parallax) Scan(visit VisitFunc) error {
	return visit(p.Child)
}

func (p *Parallax) String() string { return "Parallax" }

// Transform returns a GeoM translation of Factor * camera.Centre.
func (p *Parallax) Transform() (opts ebiten.DrawImageOptions) {
	x, y := geom.CFloat(p.camera.Centre)
	opts.GeoM.Translate(x*p.Factor, y*p.Factor)
	return opts
}
