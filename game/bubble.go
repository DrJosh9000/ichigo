package game

import (
	"math/rand"

	"drjosh.dev/gurgle/engine"
	"drjosh.dev/gurgle/geom"
)

type Bubble struct {
	Life   int
	Sprite engine.Sprite

	game *engine.Game
}

func (b *Bubble) Scan() []interface{} {
	return []interface{}{&b.Sprite}
}

func (b *Bubble) Prepare(g *engine.Game) error {
	b.game = g
	return nil
}

func (b *Bubble) Update() error {
	b.Life--
	if b.Life <= 0 {
		b.game.Unregister(b)
	}
	// not using MoveX/MoveY/... because collisions are unnecessary -
	// this is an effect particle, if it overlaps a solid, who cares
	b.Sprite.Actor.Pos = b.Sprite.Actor.Pos.Add(geom.Pt3(
		rand.Intn(3)-1, -1, rand.Intn(2)-1,
	))
	return nil
}
