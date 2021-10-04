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

	"github.com/hajimehoshi/ebiten/v2"
)

var _ interface {
	Drawer
	DrawManager
	Scanner
} = &DrawDFS{}

func init() {
	gob.Register(&DrawDFS{})
}

// DrawDFS is a DrawManager that does not add any structure. Components are
// drawn in the order in which they are encountered by a depth-first search
// through the game tree using Query, without any extra sorting based on Z
// values or consideration for DrawOrderer.
type DrawDFS struct {
	Child interface{}
	Hides

	game *Game
}

func (d *DrawDFS) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	stack := []ebiten.DrawImageOptions{*opts}
	d.game.Query(d, DrawerType,
		// visitPre
		func(x interface{}) error {
			if h, ok := x.(Hider); ok && h.Hidden() {
				return Skip
			}
			opts := stack[len(stack)-1]
			if tf, ok := x.(Transformer); ok {
				opts = concatOpts(tf.Transform(), opts)
				stack = append(stack, opts)
			}
			if x == d { // neither draw nor skip d itself
				return nil
			}
			if dr, ok := x.(Drawer); ok {
				dr.Draw(screen, &opts)
			}
			if _, isDM := x.(DrawManager); isDM {
				return Skip
			}
			return nil
		},
		// visitPost
		func(x interface{}) error {
			if _, ok := x.(Transformer); ok {
				stack = stack[:len(stack)-1]
			}
			return nil
		},
	)
}

// ManagesDrawingSubcomponents is present so DrawDFS is recognised as a
// DrawManager.
func (DrawDFS) ManagesDrawingSubcomponents() {}

func (d *DrawDFS) Prepare(g *Game) error {
	d.game = g
	return nil
}

// Scan visits d.Child.
func (d *DrawDFS) Scan(visit VisitFunc) error {
	return visit(d.Child)
}

func (d *DrawDFS) String() string { return "DrawDFS" }
