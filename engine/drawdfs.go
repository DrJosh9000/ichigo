package engine

import "github.com/hajimehoshi/ebiten/v2"

// DrawDFS is a DrawLayer that does not add any structure. Components are
// drawn in the order in which they are encountered by a depth-first search
// through the game tree, without any extra sorting based on Z values or
// consideration for DrawOrderer).
type DrawDFS struct {
	Components []interface{}
	Hides
}

func (d *DrawDFS) DrawAll(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if d.Hidden() {
		return
	}
	for _, component := range d.Components {
		d.draw(component, screen, *opts)
	}
}

func (d *DrawDFS) Scan() []interface{} { return d.Components }

func (d *DrawDFS) draw(component interface{}, screen *ebiten.Image, opts ebiten.DrawImageOptions) {
	// Hidden? stop drawing
	if h, ok := component.(Hider); ok && h.Hidden() {
		return
	}
	// Has a transform? apply to opts
	if tf, ok := component.(Transformer); ok {
		opts = concatOpts(tf.Transform(), opts)
	}
	// Is it a DrawLayer? draw all and return
	if dl, ok := component.(DrawLayer); ok {
		dl.DrawAll(screen, &opts)
		return
	}
	// Not a draw layer.
	// Does it draw itself? Draw
	if dr, ok := component.(Drawer); ok {
		dr.Draw(screen, &opts)
	}
	// Has subcomponents? recurse
	if sc, ok := component.(Scanner); ok {
		for _, ch := range sc.Scan() {
			d.draw(ch, screen, opts)
		}
	}
}
