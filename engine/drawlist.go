package engine

import "github.com/hajimehoshi/ebiten/v2"

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
	if d.list[i] == (tombstone{}) {
		return false
	}
	if d.list[j] == (tombstone{}) {
		return true
	}
	return d.list[i].DrawBefore(d.list[j]) || d.list[j].DrawAfter(d.list[i])
}

func (d drawList) Len() int { return len(d.list) }

func (d drawList) Swap(i, j int) {
	d.rev[d.list[i]], d.rev[d.list[j]] = j, i
	d.list[i], d.list[j] = d.list[j], d.list[i]
}
