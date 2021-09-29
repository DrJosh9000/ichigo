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
}

func (d *DrawDFS) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if d.Hidden() {
		return
	}
	d.drawRecursive(d.Child, screen, *opts)
}

// ManagesDrawingSubcomponents is present so DrawDFS is recognised as a
// DrawManager.
func (DrawDFS) ManagesDrawingSubcomponents() {}

func (d *DrawDFS) drawRecursive(component interface{}, screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	// Hidden? stop drawing
	if h, ok := component.(Hider); ok && h.Hidden() {
		return
	}
	// Has a transform? apply to opts
	if tf, ok := component.(Transformer); ok {
		opts = concatOpts(tf.Transform(), opts)
	}
	// Does it draw itself? Draw
	if dr, ok := component.(Drawer); ok {
		dr.Draw(screen, &opts)
	}
	// Is it a DrawManager? It manages drawing all its subcomponents.
	if _, ok := component.(DrawManager); ok {
		return
	}
	// Has subcomponents? recurse
	// TODO: use g.Children or g.Query - but need to go in Scan order...
	if sc, ok := component.(Scanner); ok {
		sc.Scan(func(x interface{}) error {
			d.drawRecursive(x, screen, opts)
			return nil
		})
	}
}

// Scan visits d.Child.
func (d *DrawDFS) Scan(visit VisitFunc) error {
	return visit(d.Child)
}

func (d *DrawDFS) String() string { return "DrawDFS" }
