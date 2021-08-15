package engine

import (
	"encoding/gob"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(Game{})
}

type GameMode int

const (
	GameModeMenu = GameMode(iota)
	GameModePlay
	GameModePause
	GameModeEdit
)

// Game implements the ebiten methods using a collection of components.
type Game struct {
	Mode         GameMode
	ScreenWidth  int
	ScreenHeight int
	*Scene

	componentsByID map[string]interface{}
}

// Draw draws the entire thing, with default draw options.
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen, ebiten.DrawImageOptions{})
}

// Layout returns the configured screen width/height.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.ScreenWidth, g.ScreenHeight
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
// Note this does not remove any references held by other components.
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

// Scan implements Scanner.
func (g *Game) Scan() []interface{} { return []interface{}{g.Scene} }

// Walk calls v with every component reachable from c via Scan, recursively,
// for as long as visit returns true.
func Walk(c interface{}, v func(interface{}) bool) {
	if !v(c) {
		return
	}
	if sc, ok := c.(Scanner); ok {
		for _, c := range sc.Scan() {
			if !v(c) {
				return
			}
			Walk(c, v)
		}
	}
}

// PrepareToRun builds the component database (using Walk) and then calls
// Prepare on every Preparer. You must call PrepareToRun before passing to
// ebiten.RunGame.
func (g *Game) PrepareToRun() {
	g.componentsByID = make(map[string]interface{})
	Walk(g, func(c interface{}) bool {
		g.RegisterComponent(c)
		return true
	})
	Walk(g, func(c interface{}) bool {
		if p, ok := c.(Prepper); ok {
			p.Prepare(g)
		}
		return true
	})
}
