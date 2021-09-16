package engine

import "github.com/hajimehoshi/ebiten/v2"

const commonDrawerComparisons = false

var _ Drawer = tombstone{}

type tombstone struct{}

func (tombstone) Draw(*ebiten.Image, *ebiten.DrawImageOptions) {}

func (tombstone) DrawAfter(x Drawer) bool { return x != tombstone{} }
func (tombstone) DrawBefore(Drawer) bool  { return false }

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
