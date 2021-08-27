package engine

import (
	"encoding/gob"
	"fmt"
	"io/fs"
	"reflect"
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
	db   map[string]Identifier    // Named components by ID
	dex  map[dexKey][]interface{} // Ancestor/behaviour index
}

type dexKey struct {
	ancestor  interface{}
	behaviour reflect.Type
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

func (g *Game) registerComponent(c interface{}, path []interface{}) error {
	// register in g.dex
	ct := reflect.TypeOf(c)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		// TODO: sub-quadratic?
		k := dexKey{c, b}
		g.dex[k] = append(g.dex[k], c)
		for _, p := range path {
			k := dexKey{p, b}
			g.dex[k] = append(g.dex[k], c)
		}
	}

	// register in g.db
	i, ok := c.(Identifier)
	if !ok {
		return nil
	}
	id := i.Ident()
	if id == "" {
		return nil
	}
	if _, exists := g.db[id]; exists {
		return fmt.Errorf("duplicate id %q", id)
	}
	g.db[id] = i
	return nil
}

// Component returns the component with a given ID, or nil if there is none.
// This only returns sensible values after LoadAndPrepare.
func (g *Game) Component(id string) Identifier {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.db[id]
}

// Query looks for components having both a given ancestor and implementing
// a given behaviour (see Behaviors in interface.go). This only returns sensible
// values after LoadAndPrepare. Note that every component is its own ancestor.
func (g *Game) Query(ancestor interface{}, behaviour reflect.Type) []interface{} {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.dex[dexKey{ancestor, behaviour}]
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

// LoadAndPrepare first calls Load on all Loaders. Once loading is complete, it
// builds the component databases and then calls Prepare on every Preparer.
// LoadAndPrepare must be called before any calls to Component or Query.
func (g *Game) LoadAndPrepare(assets fs.FS) error {
	// Load all the Loaders
	if err := Walk(g, func(c interface{}, _ []interface{}) error {
		l, ok := c.(Loader)
		if !ok {
			return nil
		}
		return l.Load(assets)
	}); err != nil {
		return err
	}

	// Build the component databases
	g.dbmu.Lock()
	g.db = make(map[string]Identifier)
	g.dex = make(map[dexKey][]interface{})
	if err := Walk(g, g.registerComponent); err != nil {
		return err
	}
	g.dbmu.Unlock()

	// Prepare all the Preppers
	return Walk(g, func(c interface{}, _ []interface{}) error {
		p, ok := c.(Prepper)
		if !ok {
			return nil
		}
		return p.Prepare(g)
	})
}
