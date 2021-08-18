package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Collider components have tangible form.
type Collider interface {
	CollidesWith(image.Rectangle) bool
}

// Drawer components can draw themselves. Draw is called often.
// Each component is responsible for calling Draw on its child components
// (so that hiding the parent can hide the children, etc).
type Drawer interface {
	Draw(screen *ebiten.Image, opts ebiten.DrawImageOptions)
}

// DrawOrderer is used to reorder layers.
type DrawOrderer interface {
	DrawOrder() float64
}

// Identifier components have a sense of self. This makes it easier for
// components to find and interact with one another.
type Identifier interface {
	Ident() string
}

// Prepper components can be prepared. It is called after the component
// database has been populated but before the game is run. The component can
// store the reference to game, if needed, and also query the component database.
type Prepper interface {
	Prepare(game *Game)
}

// Scanner components can be scanned. It is called when the game tree is walked
// (such as when the game component database is constructed).
// Scan should return a slice containing all immediate subcomponents.
type Scanner interface {
	Scan() []interface{}
}

// Updater components can update themselves. Update is called repeatedly.
// Each component is responsible for calling Update on its child components
// (so that disabling the parent prevents updates to the children, etc).
type Updater interface {
	Update() error
}
