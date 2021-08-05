package engine

import "github.com/hajimehoshi/ebiten/v2"

// Drawer components can draw themselves. Draw is called often.
// Each component is responsible for calling Draw on its child components.
type Drawer interface {
	Draw(screen *ebiten.Image, geom ebiten.GeoM)
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
// Each component is responsible for calling Update on its child components.
type Updater interface {
	Update() error
}

// ZPositioner is used to reorder layers.
type ZPositioner interface {
	Z() float64
}
