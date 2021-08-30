package engine

import (
	"encoding/gob"
	"fmt"
	"io/fs"
	"reflect"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ interface {
	Disabler
	Hider
	Identifier
	Updater
	Scanner
} = &Game{}

func init() {
	gob.Register(&Game{})
}

// Game implements the ebiten methods using a collection of components. One
// component must be the designated root component - usually a scene of some
// kind.
type Game struct {
	Disabled
	Hidden
	ScreenWidth  int
	ScreenHeight int
	Root         DrawUpdater // typically a *Scene or SceneRef though

	dbmu sync.RWMutex
	byID map[string]Identifier              // Named components by ID
	byAB map[abKey]map[interface{}]struct{} // Ancestor/behaviour index
}

type abKey struct {
	ancestor  string
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

// Ident returns "__GAME__".
func (g *Game) Ident() string { return "__GAME__" }

// Component returns the component with a given ID, or nil if there is none.
// This only returns sensible values after LoadAndPrepare.
func (g *Game) Component(id string) Identifier {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.byID[id]
}

// Query looks for components having both a given ancestor and implementing
// a given behaviour (see Behaviors in interface.go). This only returns sensible
// values after LoadAndPrepare. Note that every component is its own ancestor.
func (g *Game) Query(ancestorID string, behaviour reflect.Type) map[interface{}]struct{} {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.byAB[abKey{ancestorID, behaviour}]
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
	// Load all the Loaders.
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
	g.byID = make(map[string]Identifier)
	g.byAB = make(map[abKey]map[interface{}]struct{})
	if err := Walk(g, g.registerComponent); err != nil {
		return err
	}
	g.dbmu.Unlock()

	// Prepare all the Preppers
	for p := range g.Query(g.Ident(), PrepperType) {
		if err := p.(Prepper).Prepare(g); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) registerComponent(c interface{}, path []interface{}) error {
	// register in g.dex
	ct := reflect.TypeOf(c)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		// TODO: sub-quadratic?
		for _, p := range append(path, c) {
			i, ok := p.(Identifier)
			if !ok || i.Ident() == "" {
				continue
			}
			k := abKey{i.Ident(), b}
			if g.byAB[k] == nil {
				g.byAB[k] = make(map[interface{}]struct{})
			}
			g.byAB[k][c] = struct{}{}
		}
	}

	// register in g.db
	i, ok := c.(Identifier)
	if !ok || i.Ident() == "" {
		return nil
	}
	id := i.Ident()
	if _, exists := g.byID[id]; exists {
		return fmt.Errorf("duplicate id %q", id)
	}
	g.byID[id] = i
	return nil
}

func (g *Game) unregisterComponent(c interface{}, path []interface{}) {
	// unregister from g.dex
	ct := reflect.TypeOf(c)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		for _, p := range append(path, c) {
			i, ok := p.(Identifier)
			if !ok || i.Ident() == "" {
				continue
			}
			k := abKey{i.Ident(), b}
			if g.byAB[k] == nil {
				continue
			}
			delete(g.byAB[k], c)
		}
	}

	// unregister from g.db
	i, ok := c.(Identifier)
	if !ok || i.Ident() == "" {
		return
	}
	delete(g.byID, i.Ident())
}
