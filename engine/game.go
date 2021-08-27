package engine

import (
	"encoding/gob"
	"io/fs"
	"log"
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

	// Components by ID
	dbmu sync.RWMutex
	db   map[string]interface{}

	// Collision domains - all Collider subcomponents
	cdmu sync.RWMutex
	cd   map[string]map[Collider]struct{}
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
	if _, exists := g.db[id]; exists {
		log.Printf("duplicate id %q", id)
	}
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

// Walk calls v with every path of components reachable from c via Scan, for as
// long as visit returns nil.
func Walk(c interface{}, v func(interface{}, []interface{}) error) error {
	return walk(c, nil, v)
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

	g.cdmu.Lock()
	g.cd = make(map[string]map[Collider]struct{})
	g.cdmu.Unlock()

	g.dbmu.Lock()
	g.db = make(map[string]interface{})
	g.dbmu.Unlock()

	// -> here <- is the moment in time where db is empty.
	Walk(g.Root, func(c interface{}, p []interface{}) error {
		g.RegisterComponent(c)
		return nil
	})
	Walk(g.Root, func(c interface{}, _ []interface{}) error {
		if p, ok := c.(Prepper); ok {
			p.Prepare(g)
		}
		return nil
	})
	return nil
}
