package engine

import (
	"encoding/gob"
	"errors"
	"fmt"
	"image"
	"io/fs"
	"log"
	"reflect"
	"sync"
	"time"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ interface {
	Disabler
	Hider
	Identifier
	Updater
	Registrar
	Scanner
} = &Game{}

var (
	errNilComponent = errors.New("nil component")
	errNilParent    = errors.New("nil parent")
)

func init() {
	gob.Register(&Game{})
}

// ComponentSet is a set of components.
type ComponentSet map[interface{}]struct{}

// Game implements the ebiten methods using a collection of components. One
// component must be the designated root component.
type Game struct {
	Disables
	Hides
	ScreenSize image.Point
	Root       Drawer
	Projection geom.Projector
	VoxelScale geom.Float3

	dbmu     sync.RWMutex
	byID     map[string]Identifier        // Named components by ID
	byAB     map[abKey]ComponentSet       // Ancestor/behaviour index
	parent   map[interface{}]interface{}  // parent[x] is parent of x
	children map[interface{}]ComponentSet // children[x] are chilren of x
}

// Draw draws everything.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.Hidden() {
		return
	}
	g.Root.Draw(screen, &ebiten.DrawImageOptions{})
}

// Layout returns the configured screen width/height.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.ScreenSize.X, g.ScreenSize.Y
}

// Update updates everything.
func (g *Game) Update() error {
	return g.updateRecursive(g)
}

// updateRecursive updates everything in a post-order traversal. It terminates
// the recursion early if the component reports it is Disabled.
func (g *Game) updateRecursive(c interface{}) error {
	if d, ok := c.(Disabler); ok && d.Disabled() {
		return nil
	}
	if sc, ok := c.(Scanner); ok {
		if err := sc.Scan(g.updateRecursive); err != nil {
			return err
		}
	}
	if u, ok := c.(Updater); ok && c != g {
		return u.Update()
	}
	return nil
}

// Ident returns "__GAME__".
func (g *Game) Ident() string { return "__GAME__" }

// Component returns the component with a given ID, or nil if there is none.
// This only returns sensible values for registered components (e.g. after
// LoadAndPrepare).
func (g *Game) Component(id string) Identifier {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.byID[id]
}

// Parent returns the parent of a given component, or nil if there is none.
// This only returns sensible values for registered components (e.g. after
// LoadAndPrepare).
func (g *Game) Parent(c interface{}) interface{} {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.parent[c]
}

func (g *Game) Children(c interface{}) ComponentSet {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.children[c]
}

// PathRegister calls Register on every Registrar in the path between g and
// parent (top-to-bottom, i.e. g first)
func (g *Game) PathRegister(component, parent interface{}) error {
	rp := g.ReversePath(parent)
	for i := len(rp) - 1; i >= 0; i-- {
		if r, ok := rp[i].(Registrar); ok {
			if err := r.Register(component, parent); err != nil {
				return err
			}
		}
	}
	return nil
}

// PathUnregister calls Unregister on every Registrar in the path between g and
// parent (bottom-to-top, i.e. parent first).
func (g *Game) PathUnregister(component interface{}) {
	for _, p := range g.ReversePath(component) {
		if r, ok := p.(Registrar); ok {
			r.Unregister(component)
		}
	}
}

// Path returns a slice with the path of components to reach component from g
// (including g and component).
func (g *Game) Path(component interface{}) []interface{} {
	stack := g.ReversePath(component)
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}
	return stack
}

// ReversePath returns the same slice as Path, but reversed. (ReversePath is
// faster than Path).
func (g *Game) ReversePath(component interface{}) []interface{} {
	var stack []interface{}
	g.dbmu.RLock()
	for p := component; p != nil; p = g.Parent(p) {
		stack = append(stack, p)
	}
	g.dbmu.RUnlock()
	return stack
}

// Query looks for components having both a given ancestor and implementing
// a given behaviour (see Behaviors in interface.go). This only returns sensible
// values after LoadAndPrepare. Note that every component is its own ancestor.
func (g *Game) Query(ancestor interface{}, behaviour reflect.Type) ComponentSet {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.byAB[abKey{ancestor, behaviour}]
}

// Scan visits g.Root.
func (g *Game) Scan(visit func(interface{}) error) error {
	return visit(g.Root)
}

// Load loads a component and all subcomponents recursively.
// Note that this method does not implement Loader.
func (g *Game) Load(component interface{}, assets fs.FS) error {
	if l, ok := component.(Loader); ok {
		if err := l.Load(assets); err != nil {
			return err
		}
	}
	if sc, ok := component.(Scanner); ok {
		return sc.Scan(func(x interface{}) error {
			return g.Load(x, assets)
		})
	}
	return nil
}

// Prepare prepares a component and all subcomponents recursively.
// Note that this method does not implement Prepper.
func (g *Game) Prepare(component interface{}) error {
	// Postorder traversal, in case ancestors depend on descendants being
	// ready to answer queries.
	if sc, ok := component.(Scanner); ok {
		if err := sc.Scan(g.Prepare); err != nil {
			return err
		}
	}
	if p, ok := component.(Prepper); ok {
		return p.Prepare(g)
	}
	return nil
}

// LoadAndPrepare first calls Load on all Loaders. Once loading is complete, it
// builds the component databases and then calls Prepare on every Preparer.
// LoadAndPrepare must be called before any calls to Component or Query.
func (g *Game) LoadAndPrepare(assets fs.FS) error {
	if g.VoxelScale == (geom.Float3{}) {
		g.VoxelScale = geom.Float3{X: 1, Y: 1, Z: 1}
	}

	// Load all the Loaders.
	startLoad := time.Now()
	if err := g.Load(g.Root, assets); err != nil {
		return err
	}
	log.Printf("finished loading in %v", time.Since(startLoad))

	// Build the component databases
	startBuild := time.Now()
	if err := g.build(); err != nil {
		return err
	}
	log.Printf("finished building db in %v", time.Since(startBuild))

	// Prepare all the Preppers
	startPrep := time.Now()
	if err := g.Prepare(g.Root); err != nil {
		return err
	}
	log.Printf("finished preparing in %v", time.Since(startPrep))
	return nil
}

func (g *Game) build() error {
	g.dbmu.Lock()
	defer g.dbmu.Unlock()
	g.byID = make(map[string]Identifier)
	g.byAB = make(map[abKey]ComponentSet)
	g.parent = make(map[interface{}]interface{})
	g.children = make(map[interface{}]ComponentSet)
	return g.registerRecursive(g, nil)
}

// Register registers a component into the component database (as the
// child of a given parent). Passing a nil component or parent is an error.
// Registering multiple components with the same ID is also an error.
// Registering a component will recursively register all children found via
// Scan.
func (g *Game) Register(component, parent interface{}) error {
	if component == nil {
		return errNilComponent
	}
	if parent == nil && component != g {
		return errNilParent
	}
	g.dbmu.Lock()
	defer g.dbmu.Unlock()
	return g.registerRecursive(component, parent)
}

func (g *Game) registerRecursive(component, parent interface{}) error {
	if err := g.registerOne(component, parent); err != nil {
		return err
	}
	if sc, ok := component.(Scanner); ok {
		return sc.Scan(func(x interface{}) error {
			return g.registerRecursive(x, component)
		})
	}
	return nil
}

func (g *Game) registerOne(component, parent interface{}) error {
	// register in g.byID if needed
	if i, ok := component.(Identifier); ok {
		if id := i.Ident(); id != "" {
			if _, exists := g.byID[id]; exists {
				return fmt.Errorf("duplicate id %q", id)
			}
			g.byID[id] = i
		}
	}

	// register in g.parent and g.children
	g.parent[component] = parent
	if g.children[parent] == nil {
		g.children[parent] = make(ComponentSet)
	}
	g.children[parent][component] = struct{}{}

	// register in g.byAB
	ct := reflect.TypeOf(component)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		// TODO: better than O(len(path)^2) time and memory?
		for p := component; p != nil; p = g.parent[p] {
			k := abKey{p, b}
			if g.byAB[k] == nil {
				g.byAB[k] = make(ComponentSet)
			}
			g.byAB[k][component] = struct{}{}
		}
	}
	return nil
}

// Unregister removes the component from the component database.
// Passing a nil component has no effect. Unregistering a component will
// recursively unregister child components found via Scan.
func (g *Game) Unregister(component interface{}) {
	if component == nil {
		return
	}
	g.dbmu.Lock()
	g.unregisterRecursive(component)
	g.dbmu.Unlock()
}

func (g *Game) unregisterRecursive(component interface{}) {
	if sc, ok := component.(Scanner); ok {
		sc.Scan(func(x interface{}) error {
			g.unregisterRecursive(x)
			return nil
		})
	}
	g.unregisterOne(component)
}

func (g *Game) unregisterOne(component interface{}) {
	// unregister from g.byAB, using g.par to trace the path
	ct := reflect.TypeOf(component)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		for p := component; p != nil; p = g.parent[p] {
			if k := (abKey{p, b}); g.byAB[k] != nil {
				delete(g.byAB[k], component)
			}
		}
	}

	// unregister from g.parent and g.children
	delete(g.children[g.parent[component]], component)
	delete(g.parent, component)

	// unregister from g.byID if needed
	if id, ok := component.(Identifier); ok && id.Ident() != "" {
		delete(g.byID, id.Ident())
	}
}

func (g *Game) String() string { return "Game" }

// --------- Helper stuff ---------

type abKey struct {
	ancestor  interface{}
	behaviour reflect.Type
}

// concatOpts returns the combined options (as though a was applied and then b).
func concatOpts(a, b ebiten.DrawImageOptions) ebiten.DrawImageOptions {
	a.ColorM.Concat(b.ColorM)
	a.GeoM.Concat(b.GeoM)
	if b.CompositeMode != 0 {
		a.CompositeMode = b.CompositeMode
	}
	if b.Filter != 0 {
		a.Filter = b.Filter
	}
	return a
}
