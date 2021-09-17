package engine

import (
	"image"
	"log"

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

// Slow topological sort
func (d *drawList) topsort() {
	// Produce edge lists - O(|V|^2)
	// Count indegrees - also O(|V|^2)
	edges := make([][]int, len(d.list))
	indegree := make([]int, len(d.list))
	for i, u := range d.list {
		if u == (tombstone{}) {
			continue
		}
		var ub image.Rectangle
		switch x := u.(type) {
		case BoundingBoxer:
			ub = x.BoundingBox().BoundingRect(geom.IntProjection{X: 0, Y: 1})
		default:
			ub = image.Rect(0, 0, 320, 240)
		}
		for j, v := range d.list {
			if i == j {
				continue
			}
			if v == (tombstone{}) {
				continue
			}
			var vb image.Rectangle
			switch y := v.(type) {
			case BoundingBoxer:
				vb = y.BoundingBox().BoundingRect(geom.IntProjection{X: 0, Y: 1})
			default:
				vb = image.Rect(0, 0, 320, 240)
			}
			if !ub.Overlaps(vb) {
				// No overlap, no need to emit an edge
				continue
			}
			if u.DrawBefore(v) || v.DrawAfter(u) {
				edges[i] = append(edges[i], j)
				indegree[j]++
			}
		}
	}

	// Start queue with all zero-indegree vertices
	var queue []int
	for i, n := range indegree {
		if d.list[i] == (tombstone{}) {
			continue
		}
		if n == 0 {
			queue = append(queue, i)
		}
	}

	// Process into new list
	list := make([]Drawer, 0, len(d.list))
	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]
		d.rev[d.list[i]] = len(list)
		list = append(list, d.list[i])
		for _, j := range edges[i] {
			indegree[j]--
			if indegree[j] <= 0 {
				if indegree[j] < 0 {
					log.Printf("indegree[%d] = %d (component %v)", j, indegree[j], d.list[j])
				}
				queue = append(queue, j)
			}
		}
	}

	// Replace list
	d.list = list
	if false {
		// Update rev
		d.rev = make(map[Drawer]int, len(list))
		for i, v := range list {
			d.rev[v] = i
		}
	}
}
