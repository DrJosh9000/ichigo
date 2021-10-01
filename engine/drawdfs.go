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
// through the game tree using Scan, without any extra sorting based on Z values
// or consideration for DrawOrderer. Also, children registered in Game but not
// registered by subcomponents (such that they are visited with Scan) won't be
// drawn.
type DrawDFS struct {
	Child interface{}
	Hides

	game *Game
}

func (d *DrawDFS) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	d.drawRecursive(screen, *opts, d)
}

// ManagesDrawingSubcomponents is present so DrawDFS is recognised as a
// DrawManager.
func (DrawDFS) ManagesDrawingSubcomponents() {}

func (d *DrawDFS) drawRecursive(screen *ebiten.Image, opts ebiten.DrawImageOptions, component interface{}) {
	// Hidden? stop drawing
	if h, ok := component.(Hider); ok && h.Hidden() {
		return
	}
	// Has a transform? apply to opts
	if tf, ok := component.(Transformer); ok {
		opts = concatOpts(tf.Transform(), opts)
	}
	if component != d {
		// Does it draw itself? Draw
		if dr, ok := component.(Drawer); ok {
			dr.Draw(screen, &opts)
		}
		// Is it a DrawManager? It manages drawing all its subcomponents.
		if _, ok := component.(DrawManager); ok {
			return
		}
	}
	// Has subcomponents? recurse
	d.game.Children(component).Scan(func(x interface{}) error {
		d.drawRecursive(screen, opts, x)
		return nil
	})
}

func (d *DrawDFS) Prepare(g *Game) error {
	d.game = g
	return nil
}

// Scan visits d.Child.
func (d *DrawDFS) Scan(visit VisitFunc) error {
	return visit(d.Child)
}

func (d *DrawDFS) String() string { return "DrawDFS" }
