package engine

import (
	"image"
	"math"

	"drjosh.dev/gurgle/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

const commonDrawerComparisons = false

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

func (d drawList) Less(i, j int) bool {
	// Deal with tombstones first, in case anything else thinks it
	// needs to go last.
	if d.list[i] == (tombstone{}) {
		return false
	}
	if d.list[j] == (tombstone{}) {
		return true
	}

	if commonDrawerComparisons {
		// Common logic for known interfaces (BoundingBoxer, ZPositioner), to
		// simplify Draw{Before,After} implementations.
		switch x := d.list[i].(type) {
		case BoundingBoxer:
			xb := x.BoundingBox()
			switch y := d.list[j].(type) {
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
			switch y := d.list[j].(type) {
			case BoundingBoxer:
				return x.ZPos() < y.BoundingBox().Min.Z
			case ZPositioner:
				return x.ZPos() < y.ZPos()
			}
		}
	}

	// Fallback case: ask the components themselves
	return d.list[i].DrawBefore(d.list[j]) || d.list[j].DrawAfter(d.list[i])
}

func (d drawList) Len() int { return len(d.list) }

func (d drawList) Swap(i, j int) {
	d.rev[d.list[i]], d.rev[d.list[j]] = j, i
	d.list[i], d.list[j] = d.list[j], d.list[i]
}

// Slow topological sort. Uses a projection π to flatten bounding boxes for
// overlap tests, so that the graph is reduced.
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
			if u.DrawBefore(v) || v.DrawAfter(u) {
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
