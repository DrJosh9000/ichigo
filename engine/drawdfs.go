package engine

import "github.com/hajimehoshi/ebiten/v2"

var _ interface {
	Drawer
	DrawManager
	Scanner
} = &DrawDFS{}

// DrawDFS is a DrawLayer that does not add any structure. Components are
// drawn in the order in which they are encountered by a depth-first search
// through the game tree, without any extra sorting based on Z values or
// consideration for DrawOrderer).
type DrawDFS struct {
	Child interface{}
	Hides
}

func (d *DrawDFS) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if d.Hidden() {
		return
	}
	d.draw(d.Child, screen, *opts)
}

// exists to satisfy interface
func (DrawDFS) ManagesDrawingSubcomponents() {}

func (d *DrawDFS) draw(component interface{}, screen *ebiten.Image, opts ebiten.DrawImageOptions) {
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
	// Is it a DrawManager? return early (don't recurse)
	if _, ok := component.(DrawManager); ok {
		return
	}
	// Has subcomponents? recurse
	if sc, ok := component.(Scanner); ok {
		/*for _, ch := range sc.Scan() {
			d.draw(ch, screen, opts)
		}*/
		scv := func(x interface{}) error {
			d.draw(x, screen, opts)
			return nil
		}
		sc.Scan(scv)
	}
}

//func (d *DrawDFS) Scan() []interface{} { return []interface{}{d.Child} }
func (d *DrawDFS) Scan(visit func(interface{}) error) error {
	return visit(d.Child)
}
