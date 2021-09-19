package engine

import (
	"image"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ Drawer = tombstone{}

type tombstone struct{}

func (tombstone) Draw(*ebiten.Image, *ebiten.DrawImageOptions) {}

func (tombstone) DrawAfter(x Drawer) bool { return x != tombstone{} }
func (tombstone) DrawBefore(Drawer) bool  { return false }

func (tombstone) String() string { return "tombstone" }

type drawList struct {
	list []Drawer
	rev  map[Drawer]int
}

// edge reports if there is a draw ordering constraint between u and v (where
// u draws before v).
func edge(u, v Drawer, πsign image.Point) bool {
	// Common logic for known interfaces (BoundingBoxer, ZPositioner), to
	// simplify DrawOrderer implementations.
	switch u := u.(type) {
	case BoundingBoxer:
		ub := u.BoundingBox()
		switch v := v.(type) {
		case BoundingBoxer:
			vb := v.BoundingBox()
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
		case ZPositioner:
			return ub.Max.Z < v.ZPos() // u is before v
		}

	case ZPositioner:
		switch y := v.(type) {
		case BoundingBoxer:
			return u.ZPos() < y.BoundingBox().Min.Z
		case ZPositioner:
			return u.ZPos() < y.ZPos()
		}
	}

	// Fallback case: ask the components themselves if they have an opinion
	if do, ok := u.(DrawOrderer); ok && do.DrawBefore(v) {
		return true
	}
	if do, ok := v.(DrawOrderer); ok && do.DrawAfter(u) {
		return true
	}

	// No relation
	return false
}

// Topological sort. Uses a projection π to flatten bounding boxes for
// overlap tests, in order to reduce edge count.
func (d *drawList) topsort(π geom.Projector) {
	// Produce edge lists and count indegrees - O(|V|^2)
	// TODO: optimise this
	edges := make([][]int, len(d.list))
	indegree := make([]int, len(d.list))
	for i, u := range d.list {
		if u == (tombstone{}) {
			// Prevents processing this vertex later on
			indegree[i] = -1
			continue
		}
		// If we can't get a more specific bounding rect, assume entire plane.
		var ubr image.Rectangle
		ub, brCheck := u.(BoundingBoxer)
		if brCheck {
			ubr = ub.BoundingBox().BoundingRect(π)
		}
		// For each possible neighbor...
		for j, v := range d.list {
			if i == j || v == (tombstone{}) {
				continue
			}
			// Does it have a bounding rect? Do overlap test.
			if brCheck {
				if vb, ok := v.(BoundingBoxer); ok {
					if vbr := vb.BoundingBox().BoundingRect(π); !ubr.Overlaps(vbr) {
						continue
					}
				}
			}

			// If the edge goes u->v, add it.
			if edge(u, v, π.Sign()) {
				edges[i] = append(edges[i], j)
				indegree[j]++
			}
		}
	}

	// Initialise queue with all the zero-indegree vertices
	queue := make([]int, 0, len(d.list))
	for i, n := range indegree {
		if n == 0 {
			queue = append(queue, i)
		}
	}

	// Process into new list. O(|V| + |E|)
	list := make([]Drawer, 0, len(d.list))
	for len(queue) > 0 {
		// Get front of queue.
		i := queue[0]
		queue = queue[1:]
		// Add to output list.
		d.rev[d.list[i]] = len(list)
		list = append(list, d.list[i])
		// Reduce indegree for all outgoing edges, enqueue if indegree now 0.
		for _, j := range edges[i] {
			indegree[j]--
			if indegree[j] == 0 {
				queue = append(queue, j)
			}
		}
	}
	// Job done!
	d.list = list
}

type drawDAG struct {
	*dag
	planes    set
	chunks    map[image.Point]set
	chunksRev map[Drawer]image.Rectangle
	chunkSize int
	proj      geom.Projector
}

func newDrawDAG(chunkSize int) *drawDAG {
	return &drawDAG{
		dag:       newDAG(),
		planes:    make(set),                        // drawers that take up whole plane
		chunks:    make(map[image.Point]set),        // chunk coord -> drawers with bounding rects intersecting chunk
		chunksRev: make(map[Drawer]image.Rectangle), // drawer -> rectangle of chunk coords
		chunkSize: chunkSize,
	}
}

// add adds a Drawer and any needed edges to the DAG and chunk map.
func (d *drawDAG) add(x Drawer) {
	switch x := x.(type) {
	case BoundingBoxer:
		br := x.BoundingBox().BoundingRect(d.proj)
		min := br.Min.Div(d.chunkSize)
		max := br.Max.Sub(image.Pt(1, 1)).Div(d.chunkSize)
		cand := make(set)
		for j := min.Y; j <= max.Y; j++ {
			for i := min.X; i <= max.X; i++ {
				for c := range d.chunks[image.Pt(i, j)] {
					cand[c] = struct{}{}
				}
			}
		}
		for c := range cand {
			// TODO: x before or after c?
			d.dag.addEdge(c, x)
			d.dag.addEdge(x, c)
		}

	case ZPositioner:
		// TODO: Flat plane
		d.planes[x] = struct{}{}
	}
}

type set map[interface{}]struct{}

type dag struct {
	in, out map[interface{}]set
}

func newDAG() *dag {
	return &dag{
		in:  make(map[interface{}]set),
		out: make(map[interface{}]set),
	}
}

// addEdge adds the edge u-v in O(1).
func (d *dag) addEdge(u, v interface{}) {
	if d.in[v] == nil {
		d.in[v] = make(set)
	}
	if d.out[u] == nil {
		d.out[u] = make(set)
	}
	d.in[v][u] = struct{}{}
	d.out[u][v] = struct{}{}
}

// removeEdge removes the edge u-v in O(1).
func (d *dag) removeEdge(u, v interface{}) {
	delete(d.in[v], u)
	delete(d.out[u], v)
}

// removeVertex removes all in and out edges associated with v in O(degree(v)).
func (d *dag) removeVertex(v interface{}) {
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
}

// topIterate visits each vertex in topological order, in time O(|V| + |E|) and
// O(|V|) temporary memory.
func (d *dag) topIterate(visit func(interface{})) {
	// Count indegrees - indegree(v) = len(d.in[v]) for each v.
	// If indegree(v) = 0, enqueue. Total: O(|V|).
	queue := make([]interface{}, 0, len(d.in))
	indegree := make(map[interface{}]int)
	for u, e := range d.in {
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
