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

// Game implements the ebiten methods using a collection of components. One
// component must be the designated root component - usually a scene of some
// kind.
type Game struct {
	Disables
	Hides
	ScreenSize image.Point
	Root       Drawer // usually a DrawManager
	Projection geom.Projector
	VoxelScale geom.Float3

	dbmu sync.RWMutex
	byID map[string]Identifier              // Named components by ID
	byAB map[abKey]map[interface{}]struct{} // Ancestor/behaviour index
	par  map[interface{}]interface{}        // par[x] is parent of x
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
	if g.Disabled() {
		return nil
	}

	// Need to do a similar trick for Draw: disabling a parent object should
	// disable the child objects.
	// cache memoises the disabled state for each component.
	cache := map[interface{}]bool{
		g: false,
	}

	// Update everything that is not disabled.
	return PostorderWalk(g, func(c, _ interface{}) error {
		// Skip g (note g satisfies Updater, so this would infinitely recurse)
		if c == g {
			return nil
		}
		u, ok := c.(Updater)
		if !ok {
			return nil
		}

		// Is u disabled itself?
		if d, ok := u.(Disabler); ok && d.Disabled() {
			cache[u] = true
			return nil
		}

		// Walk up g.par to find the nearest state in accum.
		var st bool
		stack := []interface{}{u}
		for p := g.par[u]; ; p = g.par[p] {
			if s, found := cache[p]; found {
				st = s
				break
			}
			stack = append(stack, p)
		}
		// Unwind the stack, accumulating state along the way.
		for len(stack) > 0 {
			l1 := len(stack) - 1
			p := stack[l1]
			stack = stack[:l1]
			if d, ok := p.(Disabler); ok {
				st = st || d.Disabled()
			}
			cache[p] = st
		}

		// Skip updating if disabled.
		if st {
			return nil
		}

		// Update
		return u.Update()
	})
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
	return g.par[c]
}

// WalkUp visits the component, its parent, its parent, ..., and then g.
func (g *Game) WalkUp(component interface{}, visit func(interface{}) error) error {
	for p := component; p != nil; p = g.Parent(p) {
		if err := visit(p); err != nil {
			return err
		}
	}
	return nil
}

// WalkDown visits g, the subcomponent of g, ..., and then the component.
func (g *Game) WalkDown(component interface{}, visit func(interface{}) error) error {
	var stack []interface{}
	g.dbmu.RLock()
	for p := component; p != nil; p = g.Parent(p) {
		stack = append(stack, p)
	}
	g.dbmu.RUnlock()
	for _, p := range stack {
		if err := visit(p); err != nil {
			return err
		}
	}
	return nil
}

// Query looks for components having both a given ancestor and implementing
// a given behaviour (see Behaviors in interface.go). This only returns sensible
// values after LoadAndPrepare. Note that every component is its own ancestor.
func (g *Game) Query(ancestor interface{}, behaviour reflect.Type) map[interface{}]struct{} {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.byAB[abKey{ancestor, behaviour}]
}

// Scan returns g.Root.
func (g *Game) Scan() Components { return Components{g.Root} }

// PreorderWalk calls visit with every component and its parent, reachable from
// the  given component via Scan, for as long as visit returns nil. The parent
// value passed to visit when visiting component will be nil. The parent will be
// visited before the children.
func PreorderWalk(component interface{}, visit func(component, parent interface{}) error) error {
	return preorderWalk(component, nil, visit)
}

func preorderWalk(component, parent interface{}, visit func(component, parent interface{}) error) error {
	if err := visit(component, parent); err != nil {
		return err
	}
	sc, ok := component.(Scanner)
	if !ok {
		return nil
	}
	for _, c := range sc.Scan() {
		if err := preorderWalk(c, component, visit); err != nil {
			return err
		}
	}
	return nil
}

// PostorderWalk calls visit with every component and its parent, reachable from
// the  given component via Scan, for as long as visit returns nil. The parent
// value passed to visit when visiting component will be nil. The children will
// be visited before the parent.
func PostorderWalk(component interface{}, visit func(component, parent interface{}) error) error {
	return postorderWalk(component, nil, visit)
}

func postorderWalk(component, parent interface{}, visit func(component, parent interface{}) error) error {
	if sc, ok := component.(Scanner); ok {
		for _, c := range sc.Scan() {
			if err := postorderWalk(c, component, visit); err != nil {
				return err
			}
		}
	}
	return visit(component, parent)
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
	if err := PreorderWalk(g, func(c, _ interface{}) error {
		if l, ok := c.(Loader); ok {
			return l.Load(assets)
		}
		return nil
	}); err != nil {
		return err
	}
	log.Printf("finished loading in %v", time.Since(startLoad))

	// Build the component databases
	startBuild := time.Now()
	g.dbmu.Lock()
	g.byID = make(map[string]Identifier)
	g.byAB = make(map[abKey]map[interface{}]struct{})
	g.par = make(map[interface{}]interface{})
	if err := PreorderWalk(g, g.registerOne); err != nil {
		return err
	}
	g.dbmu.Unlock()
	log.Printf("finished building db in %v", time.Since(startBuild))

	// Prepare all the Preppers
	startPrep := time.Now()
	if err := PostorderWalk(g, func(c, _ interface{}) error {
		if p, ok := c.(Prepper); ok {
			return p.Prepare(g)
		}
		return nil
	}); err != nil {
		return err
	}

	log.Printf("finished preparing in %v", time.Since(startPrep))
	return nil
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
	// preorderWalk goes in the right order for registering.
	return preorderWalk(component, parent, g.registerOne)
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

	// register in g.par
	if parent != nil {
		g.par[component] = parent
	}

	// register in g.byAB
	ct := reflect.TypeOf(component)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		// TODO: better than O(len(path)^2) time and memory?
		for p := component; p != nil; p = g.par[p] {
			k := abKey{p, b}
			if g.byAB[k] == nil {
				g.byAB[k] = make(map[interface{}]struct{})
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
	// postorderWalk goes in the right order for unregistering.
	postorderWalk(component, nil, g.unregisterOne)
	g.dbmu.Unlock()
}

func (g *Game) unregisterOne(component, _ interface{}) error {
	// unregister from g.byAB, using g.par to trace the path
	ct := reflect.TypeOf(component)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		for p := component; p != nil; p = g.par[p] {
			if k := (abKey{p, b}); g.byAB[k] != nil {
				delete(g.byAB[k], component)
			}
		}
	}

	// unregister from g.par
	delete(g.par, component)

	// unregister from g.byID if needed
	if id, ok := component.(Identifier); ok && id.Ident() != "" {
		delete(g.byID, id.Ident())
	}
	return nil
}

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
