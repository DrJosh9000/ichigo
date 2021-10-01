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
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	_ interface {
		Drawer
		Hider
	} = &PerfDisplay{}

	_ interface {
		Drawer
		Hider
		Updater
	} = &DebugToast{}
)

func init() {
	gob.Register(&DebugToast{})
	gob.Register(&PerfDisplay{})
}

// DebugToast debugprints a string for a while, then disappears.
type DebugToast struct {
	ID
	Hides
	Pos   image.Point
	Timer int // ticks
	Text  string
}

func (d *DebugToast) Draw(screen *ebiten.Image, _ *ebiten.DrawImageOptions) {
	ebitenutil.DebugPrintAt(screen, d.Text, d.Pos.X, d.Pos.Y)
}

func (d *DebugToast) String() string {
	return fmt.Sprintf("DebugToast@%v", d.Pos)
}

func (d *DebugToast) Toast(text string) {
	d.Text = text
	d.Timer = 120
	d.Hides = false
}

func (d *DebugToast) Update() error {
	if d.Hides = d.Timer <= 0; !d.Hides {
		d.Timer--
	}
	return nil
}

// PerfDisplay debugprints CurrentTPS and CurrentFPS in the top left.
type PerfDisplay struct {
	Hides
}

func (p PerfDisplay) Draw(screen *ebiten.Image, _ *ebiten.DrawImageOptions) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f  FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (PerfDisplay) String() string { return "PerfDisplay" }
