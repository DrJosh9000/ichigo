package engine

import (
	"encoding/gob"
	"io/fs"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	gob.Register(&Game{})
}

// Game implements the ebiten methods using a collection of components.
// One component must be the designated root component - usually a
// scene of some kind.
type Game struct {
	Disabled
	Hidden
	ScreenWidth  int
	ScreenHeight int
	Root         DrawUpdater // typically a *Scene or SceneRef though

	dbmu sync.RWMutex
	db   map[string]interface{}
}

// Draw draws the entire thing, with default draw options.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.Hidden {
		return
	}
	g.Root.Draw(screen, ebiten.DrawImageOptions{})
}

// Layout returns the configured screen width/height.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.ScreenWidth, g.ScreenHeight
}

// Update updates the scene.
func (g *Game) Update() error {
	if g.Disabled {
		return nil
	}
	return g.Root.Update()
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
	g.dbmu.Lock()
	g.db[id] = c
	g.dbmu.Unlock()
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
	g.dbmu.Lock()
	delete(g.db, id)
	g.dbmu.Unlock()
}

// Component returns the component with a given ID, or nil if there is none.
func (g *Game) Component(id string) interface{} {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.db[id]
}

// Scan implements Scanner.
func (g *Game) Scan() []interface{} { return []interface{}{g.Root} }

// Walk calls v with every component reachable from c via Scan, recursively,
// for as long as visit returns nil.
func Walk(c interface{}, v func(interface{}) error) error {
	if err := v(c); err != nil {
		return err
	}
	sc, ok := c.(Scanner)
	if !ok {
		return nil
	}
	for _, c := range sc.Scan() {
		if err := Walk(c, v); err != nil {
			return err
		}
	}
	return nil
}

// Load calls Load on all Loaders reachable via Scan (using Walk).
// It stops on the first error.
func (g *Game) Load(assets fs.FS) error {
	return Walk(g.Root, func(c interface{}) error {
		l, ok := c.(Loader)
		if !ok {
			return nil
		}
		return l.Load(assets)
	})
}

// Prepare builds the component database (using Walk) and then calls
// Prepare on every Preparer. You must call Prepare before any calls
// to Component. You may call Prepare again (e.g. as an alternative to
// fastidiously calling RegisterComponent/UnregisterComponent).
func (g *Game) Prepare() {
	g.dbmu.Lock()
	g.db = make(map[string]interface{})
	g.dbmu.Unlock()
	// Moment in time where db is empty... whatev.
	Walk(g.Root, func(c interface{}) error {
		g.RegisterComponent(c)
		return nil
	})
	Walk(g.Root, func(c interface{}) error {
		if p, ok := c.(Prepper); ok {
			p.Prepare(g)
		}
		return nil
	})
}
