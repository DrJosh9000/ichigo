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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Fill satisfies Drawer.
var _ interface {
	Drawer
	Hider
	Identifier
} = &Fill{}

func init() {
	gob.Register(&Fill{})
	gob.Register(color.Gray{})
	gob.Register(color.RGBA{})
}

// Fill fills the screen with a colour.
type Fill struct {
	ID
	Color color.Color
	Hides
}

func (f *Fill) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.Fill(opts.ColorM.Apply(f.Color))
}

func (f *Fill) String() string { return "Fill" }
