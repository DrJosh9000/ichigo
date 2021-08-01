package engine

import (
	"encoding/gob"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(Game{})
}

// IDer components have a name. This makes it easier for components to
// find and interact with one another.
type IDer interface {
	ID() string
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

	allComponents  []interface{}
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

// Component returns the component with a given ID.
func (g *Game) Component(id string) interface{} { return g.componentsByID[id] }

// Walk calls visit with every component, for as long as visit returns true.
func (g *Game) Walk(visit func(interface{}) bool) {
	for _, c := range g.allComponents {
		if !visit(c) {
			return
		}
	}
}

// Build builds the component database.
func (g *Game) Build() {
	byID := make(map[string]interface{})
	all := []interface{}{g.Scene}
	for offset := 0; offset < len(all); offset++ {
		head := all[offset]
		if id, ok := head.(IDer); ok {
			byID[id.ID()] = head
		}
		if sc, ok := head.(Scanner); ok {
			all = append(all, sc.Scan(g))
		}
	}
	g.allComponents = all
	g.componentsByID = byID
}
