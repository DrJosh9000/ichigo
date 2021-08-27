package engine

import (
	"encoding/gob"
	"fmt"
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
	db   map[string]interface{} // Components by ID
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

func (g *Game) registerComponent(c interface{}, p []interface{}) error {
	i, ok := c.(Identifier)
	if !ok {
		return nil
	}
	id := i.Ident()
	if id == "" {
		return nil
	}
	g.dbmu.Lock()
	if _, exists := g.db[id]; exists {
		return fmt.Errorf("duplicate id %q", id)
	}
	g.db[id] = c
	g.dbmu.Unlock()
	return nil
}

// Component returns the component with a given ID, or nil if there is none.
func (g *Game) Component(id string) interface{} {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.db[id]
}

// Scan implements Scanner.
func (g *Game) Scan() []interface{} { return []interface{}{g.Root} }

// Walk calls v with every path of components reachable from c via Scan, for as
// long as visit returns nil.
func Walk(c interface{}, v func(interface{}, []interface{}) error) error {
	return walk(c, make([]interface{}, 0, 16), v)
}

func walk(c interface{}, p []interface{}, v func(interface{}, []interface{}) error) error {
	if err := v(c, p); err != nil {
		return err
	}
	sc, ok := c.(Scanner)
	if !ok {
		return nil
	}
	p = append(p, c)
	for _, c := range sc.Scan() {
		if err := walk(c, p, v); err != nil {
			return err
		}
	}
	return nil
}

// LoadAndPrepare first calls Load on all Loaders. Once loading is complete,
// it builds the component database and then calls Prepare on every Preparer.
//  You must call Prepare before any calls
// to Component. You may call Prepare again (e.g. as an alternative to
// fastidiously calling RegisterComponent/UnregisterComponent).
func (g *Game) LoadAndPrepare(assets fs.FS) error {
	if err := Walk(g.Root, func(c interface{}, _ []interface{}) error {
		l, ok := c.(Loader)
		if !ok {
			return nil
		}
		return l.Load(assets)
	}); err != nil {
		return err
	}

	g.dbmu.Lock()
	g.db = make(map[string]interface{})
	g.dbmu.Unlock()

	// -> here <- is the moment in time where db is empty.
	Walk(g.Root, g.registerComponent)
	Walk(g.Root, func(c interface{}, _ []interface{}) error {
		if p, ok := c.(Prepper); ok {
			p.Prepare(g)
		}
		return nil
	})
	return nil
}
