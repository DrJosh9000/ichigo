package engine

import (
	"encoding/gob"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(Game{})
}

// Identifier components have a sense of self. This makes it easier for
// components to find and interact with one another.
type Identifier interface {
	Ident() string
}

// Scanner components can be scanned. It is called when the game
// component database is being constructed. It should store the Game reference
// (if needed later on), and return a slice of all subcomponents.
type Scanner interface {
	Scan(game *Game) []interface{}
}

// Game implements the ebiten methods using a collection of components.
type Game struct {
	ScreenWidth  int
	ScreenHeight int
	Scene        *Scene

	componentsByID map[string]interface{}
}

// Draw draws the entire thing.
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen, ebiten.GeoM{})
}

// Layout returns the configured screen width/height.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.ScreenWidth, g.ScreenHeight
}

// Update just passes the call onto Layers.
func (g *Game) Update() error {
	return g.Scene.Update()
}

// RegisterComponent tells the game there is a new component. Currently this is
// only necessary for components with IDs.
func (g *Game) RegisterComponent(c interface{}) {
	i, ok := c.(Identifier)
	if !ok {
		return
	}
	id := i.Ident()
	if id == "" {
		return
	}
	g.componentsByID[id] = c
}

// UnregisterComponent tells the game the component is no more.
func (g *Game) UnregisterComponent(c interface{}) {
	i, ok := c.(Identifier)
	if !ok {
		return
	}
	id := i.Ident()
	if id == "" {
		return
	}
	delete(g.componentsByID, id)
}

// Component returns the component with a given ID, or nil if there is none.
func (g *Game) Component(id string) interface{} { return g.componentsByID[id] }

// Walk calls v with every component, for as long as visit returns true.
func (g *Game) Walk(v func(interface{}) bool) {
	g.walk(g.Scene, v)
}

func (g *Game) walk(c interface{}, v func(interface{}) bool) {
	if !v(c) {
		return
	}
	if sc, ok := c.(Scanner); ok {
		for _, c := range sc.Scan(g) {
			if !v(c) {
				return
			}
			g.walk(c, v)
		}
	}
}

// Build builds the component database.
func (g *Game) Build() {
	g.componentsByID = make(map[string]interface{})
	g.Walk(func(c interface{}) bool {
		g.RegisterComponent(c)
		return true
	})
}
