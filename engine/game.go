/*
Copyright 2021 Josh Deprez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	"github.com/DrJosh9000/ichigo/geom"
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
// component must be the designated root component.
type Game struct {
	Disables
	Hides
	Projection geom.Projector
	Root       Drawer
	ScreenSize image.Point
	VoxelScale geom.Float3

	dbmu     sync.RWMutex
	byID     map[string]Identifier  // Named components by ID
	byAB     map[abKey]*Container   // paths matching interface
	parent   map[any]any 	        // parent[x] is parent of x
	children map[any]*Container     // children[x] are children of x
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

// Update updates everything. Subcomponents are updated before parent
// components. Disabled components, and components with a disabled ancestor, are
// not updated.
func (g *Game) Update() error {
	return g.Query(g.Root, UpdaterType,
		func(c any) error {
			if d, ok := c.(Disabler); ok && d.Disabled() {
				// Do not update this component or descendants.
				return Skip
			}
			return nil
		},
		func(c any) error {
			if u, ok := c.(Updater); ok {
				return u.Update()
			}
			return nil
		},
	)
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
// This only returns sensible values for registered components.
func (g *Game) Parent(c any) any {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.parent[c]
}

// Children returns the direct subcomponents of the given component, or nil if
// there are none. This only returns sensible values for registered components.
func (g *Game) Children(c any) *Container {
	g.dbmu.RLock()
	defer g.dbmu.RUnlock()
	return g.children[c]
}

// PathRegister calls Register on every Registrar in the path between g and
// parent (top-to-bottom, i.e. game first, component last).
func (g *Game) PathRegister(component, parent any) error {
	for _, p := range g.Path(parent) {
		if r, ok := p.(Registrar); ok {
			if err := r.Register(component, parent); err != nil {
				return err
			}
		}
	}
	return nil
}

// PathUnregister calls Unregister on every Registrar in the path between g and
// parent (bottom-to-top, i.e. component first, game last).
func (g *Game) PathUnregister(component any) {
	for _, p := range g.ReversePath(component) {
		if r, ok := p.(Registrar); ok {
			r.Unregister(component)
		}
	}
}

// Path returns a slice with the path of components to reach component from g
// (including g and component).
func (g *Game) Path(component any) []any {
	stack := g.ReversePath(component)
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}
	return stack
}

// ReversePath returns the same slice as Path, but reversed. (ReversePath is
// faster than Path).
func (g *Game) ReversePath(component any) []any {
	var stack []any
	g.dbmu.RLock()
	for p := component; p != nil; p = g.parent[p] {
		stack = append(stack, p)
	}
	g.dbmu.RUnlock()
	return stack
}

// Query recursively searches for components having both a given ancestor and
// implementing a given behaviour (see Behaviors in interface.go).
// visitPre is called before descendants are visited, while visitPost is called
// after descendants are visited. nil visitPre/visitPost are ignored.
//
// It is up to the visitPre and visitPost callbacks to handle components that
// do not themselves implement the behaviour - more specifically, every ancestor
// (up to the given one) of each component with the behaviour will be visited.
// Visiting components in the tree that *don't* implement the behaviour is
// important when behaviours of the parent need to influence the behaviours of
// the children (e.g. a component can be a Hider and hiding all descendants, but
// not necessarily be a Drawer itself).
//
// Query only visits components that are registered.
//
// Note that every component is an ancestor of itself.
//
// Query returns the first error returned from either visitor callback, except
// Skip when it is returned from a recursive call. Returning Skip from visitPre
// will cause visitPost and the descendants of the component to be skipped (see
// the implementation of Update for an example of how to use this).
func (g *Game) Query(ancestor any, behaviour reflect.Type, visitPre, visitPost VisitFunc) error {
	if visitPre != nil {
		if err := visitPre(ancestor); err != nil {
			return err
		}
	}
	// * Update uses Query.
	// * Updaters can Register new components.
	// * Register acquires g.dbmu.Lock.
	// ==> Wrapping the whole thing in RLock would deadlock.
	// Make the read lock as tight as possible.
	g.dbmu.RLock()
	q := g.byAB[abKey{ancestor, behaviour}]
	g.dbmu.RUnlock()
	if err := q.Scan(func(x any) error {
		if err := g.Query(x, behaviour, visitPre, visitPost); err != nil {
			if errors.Is(err, Skip) {
				return nil
			}
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if visitPost != nil {
		return visitPost(ancestor)
	}
	return nil
}

// Scan visits g.Root.
func (g *Game) Scan(visit VisitFunc) error {
	return visit(g.Root)
}

// Load loads a component and all subcomponents recursively.
// Note that this method does not implement Loader itself.
func (g *Game) Load(component any, assets fs.FS) error {
	// Query cannot be used for this method because Load might cause
	// subcomponents to spring into existence.
	if l, ok := component.(Loader); ok {
		if err := l.Load(assets); err != nil {
			return err
		}
	}
	if sc, ok := component.(Scanner); ok {
		return sc.Scan(func(x any) error {
			return g.Load(x, assets)
		})
	}
	return nil
}

// Prepare prepares a component and all subcomponents recursively.
// Note that this method does not implement Prepper itself.
func (g *Game) Prepare(component any) error {
	// Postorder traversal, in case ancestors depend on descendants being
	// ready to answer queries.
	return g.Query(component, PrepperType, nil, func(c any) error {
		if p, ok := c.(Prepper); ok {
			return p.Prepare(g)
		}
		return nil
	})
}

// LoadAndPrepare first calls Load on all Loaders. Once loading is complete, it
// builds the component databases and then calls Prepare on every Preparer.
// LoadAndPrepare must be called before any calls to Component or Query.
func (g *Game) LoadAndPrepare(assets fs.FS) error {
	if g.Projection == nil {
		g.Projection = geom.ElevationProjection{}
	}
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
	g.byAB = make(map[abKey]*Container)
	g.parent = make(map[any]any)
	g.children = make(map[any]*Container)
	return g.registerRecursive(g, nil)
}

// Register registers a component into the component database (as the
// child of a given parent). Passing a nil component or parent is an error.
// Registering multiple components with the same ID is also an error.
// Registering a component will recursively register all children found via
// Scan.
func (g *Game) Register(component, parent any) error {
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

func (g *Game) registerRecursive(component, parent any) error {
	if err := g.registerOne(component, parent); err != nil {
		return err
	}
	if sc, ok := component.(Scanner); ok {
		return sc.Scan(func(x any) error {
			return g.registerRecursive(x, component)
		})
	}
	return nil
}

func (g *Game) registerOne(component, parent any) error {
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
		g.children[parent] = MakeContainer(component)
	} else {
		g.children[parent].Add(component)
	}

	// register in g.byAB
	ct := reflect.TypeOf(component)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		for c, p := component, g.parent[component]; p != nil; c, p = p, g.parent[p] {
			k := abKey{p, b}
			if g.byAB[k] == nil {
				g.byAB[k] = MakeContainer(c)
				continue
			}
			if g.byAB[k].Contains(c) {
				break
			}
			g.byAB[k].Add(c)
		}
	}
	return nil
}

// Unregister removes the component from the component database.
// Passing a nil component has no effect. Unregistering a component will
// recursively unregister child components found via Scan.
func (g *Game) Unregister(component any) {
	if component == nil {
		return
	}
	g.dbmu.Lock()
	g.unregisterRecursive(component)
	g.dbmu.Unlock()
}

func (g *Game) unregisterRecursive(component any) {
	g.children[component].Scan(func(x any) error {
		g.unregisterRecursive(x)
		return nil
	})
	g.unregisterOne(component)
}

func (g *Game) unregisterOne(component any) {
	parent := g.parent[component]

	// unregister from g.byAB
	ct := reflect.TypeOf(component)
	for _, b := range Behaviours {
		if !ct.Implements(b) {
			continue
		}
		for c, p := component, parent; p != nil; c, p = p, g.parent[p] {
			k := abKey{p, b}
			g.byAB[k].Remove(c)
			if g.byAB[k].ItemCount() > 0 {
				break
			}
		}
	}

	// unregister from g.parent and g.children
	g.children[parent].Remove(component)
	delete(g.parent, component)

	// unregister from g.byID if needed
	if id, ok := component.(Identifier); ok && id.Ident() != "" {
		delete(g.byID, id.Ident())
	}
}

func (g *Game) String() string { return "Game" }

// --------- Helper stuff ---------

// abKey is the key type for game.byAB.
type abKey struct {
	parent    any
	behaviour reflect.Type
}

func (a abKey) String() string {
	return fmt.Sprintf("(%v %s)", a.parent, a.behaviour.Name())
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

// VisitFunc callbacks are either provided or called by various Game functions.
// For example, Query takes two VisitFuncs that are called for each result, and
// Scan is given a VisitFunc that should be called with each component. For
// recursive operations, return Skip for components that should be skipped.
type VisitFunc func(any) error

// Many calls a VisitFunc for multiple args, and returns on first non-nil error.
func (v VisitFunc) Many(x ...any) error {
	for _, c := range x {
		if err := v(c); err != nil {
			return err
		}
	}
	return nil
}

// Skip is an "error" value that can be returned from visitor callbacks. It
// tells recursive methods of Game to skip processing the current item and its
// descendants, but will otherwise continue processing.
const Skip = skip("skip")

type skip string

func (s skip) Error() string { return string(s) }
