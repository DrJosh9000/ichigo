package engine

import (
	"encoding/gob"
	"fmt"
	"image"
	"log"
	"math"
	"strings"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ interface {
	Drawer
	DrawManager
	Hider
	Prepper
	Registrar
	Scanner
	Updater
} = &DrawDAG{}

func init() {
	gob.Register(&DrawDAG{})
}

// DrawDAG is a DrawLayer that organises DrawBoxer descendants in a directed
// acyclic graph (DAG), in order to draw them according to ordering constraints.
// It combines a DAG with a spatial index used when updating vertices to reduce
// the number of tests between components.
type DrawDAG struct {
	ChunkSize int
	Child     interface{}
	Hides

	dag
	boxCache  map[DrawBoxer]geom.Box        // used to find components that moved
	chunks    map[image.Point]drawerSet     // chunk coord -> drawers with bounding rects intersecting chunk
	chunksRev map[DrawBoxer]image.Rectangle // comopnent -> rectangle of chunk coords
	game      *Game
}

// Draw draws everything in the DAG in topological order.
func (d *DrawDAG) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if d.Hidden() {
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
		d: {
			hidden: false,
			opts:   *opts,
		},
	}
	// Draw everything in d.dag, where not hidden (itself or any parent)
	d.dag.topWalk(func(x Drawer) {
		// Is d hidden itself?
		if h, ok := x.(Hider); ok && h.Hidden() {
			cache[x] = state{hidden: true}
			return // skip drawing
		}
		// Walk up game tree to find the nearest state in cache.
		var st state
		stack := []interface{}{x}
		for p := d.game.Parent(x); p != nil; p = d.game.Parent(p) {
			if s, found := cache[p]; found {
				st = s
				break
			}
			stack = append(stack, p)
		}
		// Unwind the stack, accumulating state along the way.
		for i := len(stack) - 1; i >= 0; i-- {
			p := stack[i]
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
			return
		}
		x.Draw(screen, &st.opts)
	})
}

// exists to satisfy interface
func (DrawDAG) ManagesDrawingSubcomponents() {}

// Prepare adds all subcomponents to the DAG.
func (d *DrawDAG) Prepare(game *Game) error {
	d.dag = make(dag)
	d.boxCache = make(map[DrawBoxer]geom.Box)
	d.chunks = make(map[image.Point]drawerSet)
	d.chunksRev = make(map[DrawBoxer]image.Rectangle)
	d.game = game

	// Because Game.LoadAndPrepare calls Prepare in a post-order walk, all the
	// descendants should be prepared, meaning BoundingBox (hence Register) is
	// likely to be a safe call.
	return d.Register(d, nil)
}

// Scan visits d.Child.
func (d *DrawDAG) Scan(visit func(interface{}) error) error {
	return visit(d.Child)
}

func (d *DrawDAG) String() string { return "DrawDAG" }

// Update checks for any changes to descendants, and updates its internal
// data structures accordingly.
func (d *DrawDAG) Update() error {
	// Re-evaluate bounding boxes for all descendants. If a box has changed,
	// fix up the edges by removing and re-adding the vertex.
	// Thanks once again to postorder traversal, this happens after all
	// descendant updates.
	var readd []DrawBoxer
	for db, bb := range d.boxCache {
		nbb := db.BoundingBox()
		if bb != nbb {
			d.Unregister(db)
			readd = append(readd, db)
		}
	}
	for _, db := range readd {
		d.Register(db, nil)
	}
	return nil
}

// Register recursively registers compponent and all descendants that are
// DrawBoxers into internal data structures (the DAG, etc) unless they are
// descendants of a different DrawManager.
func (d *DrawDAG) Register(component, _ interface{}) error {
	// *Don't* register the component if it is inside a descendant DrawManager.
	// These queries work because component should be registered in game before
	// this call.
	for dm := range d.game.Query(d, DrawManagerType) {
		if dm == d {
			continue
		}
		dbs := d.game.Query(dm, DrawBoxerType)
		if _, found := dbs[component]; found {
			return nil
		}
	}
	if db, ok := component.(DrawBoxer); ok {
		d.registerOne(db)
	}
	if _, ok := component.(DrawManager); ok && component != d {
		return nil
	}
	if sc, ok := component.(Scanner); ok {
		scv := func(x interface{}) error {
			return d.Register(x, nil)
		}
		return sc.Scan(scv)
	}
	return nil
}

// registerOne adds component and any needed edges to the DAG and chunk map.
func (d *DrawDAG) registerOne(x DrawBoxer) {
	// Ensure vertex is present
	d.dag.addVertex(x)

	// Update the box cache
	xb := x.BoundingBox()
	d.boxCache[x] = xb

	// Update the reverse chunk map
	xbr := xb.BoundingRect(d.game.Projection)
	revr := image.Rectangle{
		Min: xbr.Min.Div(d.ChunkSize),
		Max: xbr.Max.Sub(image.Pt(1, 1)).Div(d.ChunkSize),
	}
	d.chunksRev[x] = revr

	// Find possible edges between x and items in the overlapping cells.
	// First, a set of all the items in those cells.
	cand := make(drawerSet)
	var p image.Point
	for p.Y = revr.Min.Y; p.Y <= revr.Max.Y; p.Y++ {
		for p.X = revr.Min.X; p.X <= revr.Max.X; p.X++ {
			cell := d.chunks[p]
			if cell == nil {
				cell = make(drawerSet)
				d.chunks[p] = cell
			}
			// Merge cell contents into cand
			for c := range cell {
				cand[c] = struct{}{}
			}
			// Add x to cell
			cell[x] = struct{}{}
		}
	}
	// Add edges between x and elements of cand
	πsign := d.game.Projection.Sign()
	for c := range cand {
		y := c.(DrawBoxer)
		// Bounding rectangle overlap test
		// No overlap, no edge.
		if ybr := y.BoundingBox().BoundingRect(d.game.Projection); !xbr.Overlaps(ybr) {
			continue
		}
		switch {
		case drawOrderConstraint(x, y, πsign):
			d.dag.addEdge(x, y)
		case drawOrderConstraint(y, x, πsign):
			d.dag.addEdge(y, x)
		}
	}
}

// Unregister unregisters the component and all subcomponents.
func (d *DrawDAG) Unregister(component interface{}) {
	if db, ok := component.(DrawBoxer); ok {
		d.unregisterOne(db)
	}
	if _, ok := component.(DrawManager); ok && component != d {
		return
	}
	if sc, ok := component.(Scanner); ok {
		scv := func(x interface{}) error {
			d.Unregister(x)
			return nil
		}
		sc.Scan(scv)
	}
}

func (d *DrawDAG) unregisterOne(x DrawBoxer) {
	// Remove from chunk map
	revr := d.chunksRev[x]
	for j := revr.Min.Y; j <= revr.Max.Y; j++ {
		for i := revr.Min.X; i <= revr.Max.X; i++ {
			delete(d.chunks[image.Pt(i, j)], x)
		}
	}
	// Remove from reverse chunk map
	delete(d.chunksRev, x)
	// Remove from box cache
	delete(d.boxCache, x)
	// Remove from DAG
	d.dag.removeVertex(x)
}

// drawOrderConstraint reports if there is a draw ordering constraint between u
// and v (where u must draw before v).
func drawOrderConstraint(u, v DrawBoxer, πsign image.Point) bool {
	// Common logic for known interfaces (BoundingBoxer, ZPositioner), to
	// simplify DrawOrderer implementations.
	ub, vb := u.BoundingBox(), v.BoundingBox()
	if ub.Min.Z >= vb.Max.Z { // u is in front of v
		return false
	}
	if ub.Max.Z <= vb.Min.Z { // u is behind v
		return true
	}
	if πsign.X != 0 {
		if ub.Max.X*πsign.X <= vb.Min.X*πsign.X { // u is to the left of v
			return false
		}
		if ub.Min.X*πsign.X >= vb.Max.X*πsign.X { // u is to the right of v
			return true
		}
	}
	if πsign.Y != 0 {
		if ub.Max.Y*πsign.Y <= vb.Min.Y*πsign.Y { // u is above v
			return false
		}
		if ub.Min.Y*πsign.Y >= vb.Max.Y*πsign.Y { // u is below v
			return true
		}
	}

	// Ask the components themselves if they have an opinion
	if do, ok := u.(DrawOrderer); ok && do.DrawBefore(v) {
		return true
	}
	if do, ok := v.(DrawOrderer); ok && do.DrawAfter(u) {
		return true
	}

	// No relation
	return false
}

type drawerSet map[Drawer]struct{}

func (s drawerSet) String() string {
	var sb strings.Builder
	sb.WriteString("{ ")
	for x := range s {
		fmt.Fprintf(&sb, "%v ", x)
	}
	sb.WriteString("}")
	return sb.String()
}

type edges struct {
	in, out drawerSet
}

type dag map[Drawer]edges

// Dot returns a dot-syntax-like description of the graph.
func (d dag) Dot() string {
	var sb strings.Builder
	sb.WriteString("digraph {\n")
	for v, e := range d {
		fmt.Fprintf(&sb, "%v -> %v\n", v, e.out)
	}
	sb.WriteString(" }\n")
	return sb.String()
}

// addVertex ensures the vertex is present, even if there are no edges.
func (d dag) addVertex(v Drawer) {
	if _, found := d[v]; found {
		return
	}
	d[v] = edges{
		in:  make(drawerSet),
		out: make(drawerSet),
	}
}

// removeVertex removes all in and out edges associated with v in O(degree(v)).
func (d dag) removeVertex(v Drawer) {
	for u := range d[v].in {
		d.removeEdge(u, v)
	}
	for w := range d[v].out {
		d.removeEdge(v, w)
	}
	delete(d, v)
}

// addEdge adds the edge (u -> v) in O(1).
func (d dag) addEdge(u, v Drawer) {
	d.addVertex(u)
	d.addVertex(v)
	d[v].in[u] = struct{}{}
	d[u].out[v] = struct{}{}
}

// removeEdge removes the edge (u -> v) in O(1).
func (d dag) removeEdge(u, v Drawer) {
	delete(d[v].in, u)
	delete(d[u].out, v)
}

// topWalk visits each vertex in topological order, in time O(|V| + |E|) and
// O(|V|) temporary memory (for acyclic graphs) and a bit longer if it has to
// break cycles.
func (d dag) topWalk(visit func(Drawer)) {
	// Count indegrees - indegree(v) = len(d[v].in) for each vertex v.
	// If indegree(v) = 0, enqueue. Total: O(|V|).
	queue := make([]Drawer, 0, len(d))
	indegree := make(map[Drawer]int)
	for v, e := range d {
		if len(e.in) == 0 {
			queue = append(queue, v)
		} else {
			indegree[v] = len(e.in)
		}
	}

	for len(indegree) > 0 || len(queue) > 0 {
		if len(queue) == 0 {
			// Getting here means there are unprocessed vertices, because none
			// have zero indegree, otherwise they would be in the queue.
			// Arbitrarily break a cycle by enqueueing the least-indegree vertex
			mind, minv := math.MaxInt, Drawer(nil)
			for v, d := range indegree {
				if d < mind {
					mind, minv = d, v
				}
			}
			log.Printf("breaking cycle in 'DAG' by enqueueing %v with indegree %d", minv, mind)
			queue = append(queue, minv)
			delete(indegree, minv)
		}

		// Visit every vertex (O(|V|)) and decrement indegrees for every out edge
		// of each vertex visited (O(|E|)). Total: O(|V|+|E|).
		for len(queue) > 0 {
			u := queue[0]
			visit(u)
			queue = queue[1:]

			// Decrement indegree for all out-edges, and enqueue target if its
			// indegree is now 0.
			for v := range d[u].out {
				if _, ready := indegree[v]; !ready {
					// Vertex already drawn. This happens if there was a cycle.
					continue
				}
				indegree[v]--
				if indegree[v] == 0 {
					queue = append(queue, v)
					delete(indegree, v)
				}
			}
		}
	}
}
