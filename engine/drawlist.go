package engine

import (
	"image"
	"math"

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
func edge(u, v Drawer) bool {
	// Common logic for known interfaces (BoundingBoxer, ZPositioner), to
	// simplify DrawOrderer implementations.
	switch x := u.(type) {
	case BoundingBoxer:
		xb := x.BoundingBox()
		switch y := v.(type) {
		case BoundingBoxer:
			yb := y.BoundingBox()
			if xb.Min.Z >= yb.Max.Z { // x is in front of y
				return false
			}
			if xb.Max.Z <= yb.Min.Z { // x is behind y
				return true
			}
			if xb.Max.Y <= yb.Min.Y { // x is above y
				return false
			}
			if xb.Min.Y >= yb.Max.Y { // x is below y
				return true
			}
		case ZPositioner:
			return xb.Max.Z < y.ZPos() // x is before y
		}

	case ZPositioner:
		switch y := v.(type) {
		case BoundingBoxer:
			return x.ZPos() < y.BoundingBox().Min.Z
		case ZPositioner:
			return x.ZPos() < y.ZPos()
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
		ub := image.Rect(math.MinInt, math.MinInt, math.MaxInt, math.MaxInt)
		if x, ok := u.(BoundingBoxer); ok {
			ub = x.BoundingBox().BoundingRect(π)
		}
		// For each possible neighbor...
		for j, v := range d.list {
			if i == j || v == (tombstone{}) {
				continue
			}
			// Does it have a bounding rect? Do overlap test.
			if y, ok := v.(BoundingBoxer); ok {
				if vb := y.BoundingBox().BoundingRect(π); !ub.Overlaps(vb) {
					continue
				}
			}

			// If the edge goes u->v, add it.
			if edge(u, v) {
				edges[i] = append(edges[i], j)
				indegree[j]++
			}
		}
	}

	// Initialise queue with all the zero-indegree vertices
	var queue []int
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
