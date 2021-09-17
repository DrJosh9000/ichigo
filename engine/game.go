package engine

import (
	"encoding/gob"
	"errors"
	"fmt"
	"image"
	"io/fs"
	"log"
	"reflect"
	"sort"
	"sync"
	"time"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const topologicalDrawSort = true

var _ interface {
	Disabler
	Hider
	Identifier
	Updater
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
	Root       interface{} // typically a *Scene or SceneRef though
	Projection geom.IntProjection
	VoxelScale geom.Float3

	dbmu     sync.RWMutex
	byID     map[string]Identifier              // Named components by ID
	byAB     map[abKey]map[interface{}]struct{} // Ancestor/behaviour index
	drawList drawList                           // draw list :|
	par      map[interface{}]interface{}        // par[x] is parent of x
}

// Draw draws everything.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.Hidden() {
		return
	}

	// Hiding a parent component should hide the child objects, and the
	// transform applied to a child should be the cumulative transform of all
	// parents as well.
	// cache memoises the results for each component.
	type state struct {
		hidden bool
		opts   ebiten.DrawImageOptions
	}
	cache := map[interface{}]state{
		g: {hidden: false},
	}
	// Draw everything in g.drawList, where not hidden (itself or any parent)
	for _, d := range g.drawList.list {
		// Is d hidden itself?
		if h, ok := d.(Hider); ok && h.Hidden() {
			cache[d] = state{hidden: true}
			continue // skip drawing
		}
		// Walk up g.par to find the nearest state in accum.
		var st state
		stack := []interface{}{d}
		for p := g.par[d]; ; p = g.par[p] {
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
			if h, ok := p.(Hider); ok {
				st.hidden = st.hidden || h.Hidden()
			}
			if st.hidden {
				cache[p] = state{hidden: true}
				continue
			}
			// p is not hidden, so compute its cumulative opts.
			if tf, ok := p.(Transformer); ok {
				st.opts = concatOpts(tf.Transform(), st.opts)
			}
			cache[p] = st
		}

		// Skip drawing if hidden.
		if st.hidden {
			continue
		}
		d.Draw(screen, &st.opts)
	}

	if true {
		// Infodump about draw list
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("len(drawList.list) = %d", len(g.drawList.list)), 0, 30)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("len(drawList.rev) = %d", len(g.drawList.list)), 0, 45)
	}
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
	// TODO: do it in a fixed order? map essentially randomises iteration order
	for u := range g.Query(g.Ident(), UpdaterType) {
		// Skip g (note g satisfies Updater, so this would infinitely recurse)
		if u == g {
			continue
		}

		// Is u disabled itself?
		if d, ok := u.(Disabler); ok && d.Disabled() {
			cache[u] = true
			continue
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
			continue
		}

		if err := u.(Updater).Update(); err != nil {
			return err
		}
	}

	// Sort the draw list (on every frame - this isn't as bad as it sounds)
	if topologicalDrawSort {
		g.drawList.topsort()
	} else {
		sort.Stable(g.drawList)

		// Truncate tombstones from the end.
		for i := g.drawList.Len() - 1; i >= 0; i-- {
			if g.drawList.list[i] != (tombstone{}) {
				break
			}
			g.drawList.list = g.drawList.list[:i]
		}
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
	return g.par[c]
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
	return preorderWalk(component, nil, visit)
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
		l, ok := c.(Loader)
		if !ok {
			return nil
		}
		return l.Load(assets)
	}); err != nil {
		return err
	}
	log.Printf("finished loading in %v", time.Since(startLoad))

	// Build the component databases
	startBuild := time.Now()
	g.dbmu.Lock()
	g.byID = make(map[string]Identifier)
	g.byAB = make(map[abKey]map[interface{}]struct{})
	g.drawList.list = nil
	g.drawList.rev = make(map[Drawer]int)
	g.par = make(map[interface{}]interface{})
	if err := PreorderWalk(g, g.register); err != nil {
		return err
	}
	g.dbmu.Unlock()
	log.Printf("finished building db in %v", time.Since(startBuild))

	// Prepare all the Preppers
	startPrep := time.Now()
	for p := range g.Query(g.Ident(), PrepperType) {
		if err := p.(Prepper).Prepare(g); err != nil {
			return err
		}
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
	// walk goes in the right order for registering.
	return preorderWalk(component, parent, g.register)
}

func (g *Game) register(component, parent interface{}) error {
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

	// register in g.drawList
	if d, ok := component.(Drawer); ok {
		if _, exists := g.drawList.rev[d]; exists {
			// already registered
			return fmt.Errorf("double registration of %v", d)
		}
		g.drawList.rev[d] = len(g.drawList.list)
		g.drawList.list = append(g.drawList.list, d)
	}

	// register in g.byAB
	ct := reflect.TypeOf(component)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		// TODO: better than O(len(path)^2) time and memory?
		for p := component; p != nil; p = g.par[p] {
			id, ok := p.(Identifier)
			if !ok || id.Ident() == "" {
				continue
			}
			k := abKey{id.Ident(), b}
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
	postorderWalk(component, nil, func(c, _ interface{}) error {
		g.unregister(c)
		return nil
	})
	g.dbmu.Unlock()
}

func (g *Game) unregister(component interface{}) {
	// unregister from g.byAB, using g.par to trace the path
	ct := reflect.TypeOf(component)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		for p := component; p != nil; p = g.par[p] {
			id, ok := p.(Identifier)
			if !ok || id.Ident() == "" {
				continue
			}
			k := abKey{id.Ident(), b}
			if g.byAB[k] == nil {
				continue
			}
			delete(g.byAB[k], component)
		}
	}

	// unregister from g.par
	delete(g.par, component)

	// unregister from g.drawList
	if d, ok := component.(Drawer); ok {
		if i, found := g.drawList.rev[d]; found {
			g.drawList.list[i] = tombstone{}
			delete(g.drawList.rev, d)
		}
	}

	// unregister from g.byID if needed
	if id, ok := component.(Identifier); ok && id.Ident() != "" {
		delete(g.byID, id.Ident())
	}
}

// --------- Helper stuff ---------

type abKey struct {
	ancestor  string
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
