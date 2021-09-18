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

var wholePlane = image.Rect(math.MinInt, math.MinInt, math.MaxInt, math.MaxInt)

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
		ubr := wholePlane
		if x, ok := u.(BoundingBoxer); ok {
			ubr = x.BoundingBox().BoundingRect(π)
		}
		// For each possible neighbor...
		for j, v := range d.list {
			if i == j || v == (tombstone{}) {
				continue
			}
			// Does it have a bounding rect? Do overlap test.
			if y, ok := v.(BoundingBoxer); ok {
				if vbr := y.BoundingBox().BoundingRect(π); !ubr.Overlaps(vbr) {
					continue
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
