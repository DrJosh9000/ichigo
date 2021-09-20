package engine

import (
	"fmt"
	"image"
	"strings"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// DrawDAG is a DrawLayer that organises DrawBoxer descendants in a directed
// acyclic graph (DAG), in order to draw them according to ordering constraints.
// It combines a DAG with a spatial index used when adding new vertices
// in order to reduce the number of tests between components.
type DrawDAG struct {
	ChunkSize  int
	Components []interface{}
	Hides

	*dag
	boxCache  map[DrawBoxer]geom.Box
	chunks    map[image.Point]drawerSet     // chunk coord -> drawers with bounding rects intersecting chunk
	chunksRev map[DrawBoxer]image.Rectangle // comopnent -> rectangle of chunk coords
	parent    func(x interface{}) interface{}
	proj      geom.Projector
}

// Draw draws everything in the DAG in topological order.
func (d *DrawDAG) DrawAll(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
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
	d.dag.topIterate(func(x Drawer) {
		// Is d hidden itself?
		if h, ok := x.(Hider); ok && h.Hidden() {
			cache[x] = state{hidden: true}
			return // skip drawing
		}
		// Walk up game tree to find the nearest state in cache.
		var st state
		stack := []interface{}{x}
		for p := d.parent(x); ; p = d.parent(p) {
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
			return
		}
		x.Draw(screen, &st.opts)
	})
}

func (d *DrawDAG) Prepare(game *Game) error {
	d.dag = newDAG()
	d.boxCache = make(map[DrawBoxer]geom.Box)
	d.chunks = make(map[image.Point]drawerSet)
	d.chunksRev = make(map[DrawBoxer]image.Rectangle)
	d.parent = game.Parent
	d.proj = game.Projection

	// Descendants might not be prepared yet, so fill the cache with zero boxes
	// and fill remaining data structures during update
	// TODO: work out a system for dependent prepares........ sync.Once?
	return PreorderWalk(d, func(c, _ interface{}) error {
		if db, ok := c.(DrawBoxer); ok {
			d.boxCache[db] = geom.Box{}
		}
		return nil
	})
}

func (d *DrawDAG) Scan() []interface{} { return d.Components }

func (d *DrawDAG) Update() error {
	// Re-evaluate bounding boxes for all descendants. If a box has changed,
	// fix up the edges by removing and re-adding the vertex.
	var readd []DrawBoxer
	for db, bb := range d.boxCache {
		nbb := db.BoundingBox()
		if bb != nbb {
			d.Remove(db)
			readd = append(readd, db)
		}
	}
	for _, db := range readd {
		d.Add(db)
	}
	return nil
}

// Add adds a Drawer and any needed edges to the DAG and chunk map.
func (d *DrawDAG) Add(x DrawBoxer) {
	πsign := d.proj.Sign()

	// Update the box cache
	bb := x.BoundingBox()
	d.boxCache[x] = bb

	// Update the reverse chunk map
	br := bb.BoundingRect(d.proj)
	revr := image.Rectangle{
		Min: br.Min.Div(d.ChunkSize),
		Max: br.Max.Sub(image.Pt(1, 1)).Div(d.ChunkSize),
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
	for c := range cand {
		y := c.(DrawBoxer)
		// Bounding rectangle overlap test
		// No overlap, no edge.
		if ybr := y.BoundingBox().BoundingRect(d.proj); !br.Overlaps(ybr) {
			continue
		}
		switch {
		case drawOrderConstraint(y, x, πsign):
			d.dag.addEdge(y, x)
		case drawOrderConstraint(x, y, πsign):
			d.dag.addEdge(x, y)
		}
	}
}

// Remove removes a Drawer and all associated edges and metadata.
func (d *DrawDAG) Remove(x DrawBoxer) {
	// Remove from chunk map
	revr := d.chunksRev[x]
	for j := revr.Min.Y; j <= revr.Max.Y; j++ {
		for i := revr.Min.X; i <= revr.Max.X; i++ {
			delete(d.chunks[image.Pt(i, j)], x)
		}
	}
	// Remove from reverse chunk map
	delete(d.chunksRev, x)
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

type dag struct {
	all     drawerSet
	in, out map[Drawer]drawerSet
}

func newDAG() *dag {
	return &dag{
		all: make(drawerSet),
		in:  make(map[Drawer]drawerSet),
		out: make(map[Drawer]drawerSet),
	}
}

func (d *dag) String() string {
	var sb strings.Builder
	sb.WriteString("digraph {\n")
	for v, e := range d.out {
		fmt.Fprintf(&sb, "%v -> %v\n", v, e)
	}
	sb.WriteString(" }\n")
	return sb.String()
}

// addEdge adds the edge u-v in O(1).
func (d *dag) addEdge(u, v Drawer) {
	d.all[u], d.all[v] = struct{}{}, struct{}{}
	if d.in[v] == nil {
		d.in[v] = make(drawerSet)
	}
	if d.out[u] == nil {
		d.out[u] = make(drawerSet)
	}
	d.in[v][u] = struct{}{}
	d.out[u][v] = struct{}{}
}

/*
// removeEdge removes the edge u-v in O(1).
func (d *dag) removeEdge(u, v Drawer) {
	delete(d.in[v], u)
	delete(d.out[u], v)
}
*/

// removeVertex removes all in and out edges associated with v in O(degree(v)).
func (d *dag) removeVertex(v Drawer) {
	for u := range d.in[v] {
		// u-v is no longer an edge
		delete(d.out[u], v)
	}
	for w := range d.out[v] {
		// v-w is no longer an edge
		delete(d.in[w], v)
	}
	delete(d.in, v)
	delete(d.out, v)
	delete(d.all, v)
}

// topIterate visits each vertex in topological order, in time O(|V| + |E|) and
// O(|V|) temporary memory.
func (d *dag) topIterate(visit func(Drawer)) {
	// Count indegrees - indegree(v) = len(d.in[v]) for each v.
	// If indegree(v) = 0, enqueue. Total: O(|V|).
	queue := make([]Drawer, 0, len(d.in))
	indegree := make(map[Drawer]int)
	for u := range d.all {
		// NB: zero indegree vertices may be missing from d.in
		e := d.in[u]
		if len(e) == 0 {
			queue = append(queue, u)
		} else {
			indegree[u] = len(e)
		}
	}

	// Visit every vertex (O(|V|)) and decrement indegrees for every out edge
	// of each vertex visited (O(|E|)). Total: O(|V|+|E|).
	for len(queue) > 0 {
		u := queue[0]
		visit(u)
		queue = queue[1:]

		// Decrement indegree for all out edges, and enqueue target if its
		// indegree is now 0.
		for v := range d.out[u] {
			indegree[v]--
			if indegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
}
