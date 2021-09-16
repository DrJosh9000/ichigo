package game

import (
	"math/rand"

	"drjosh.dev/gurgle/engine"
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
	b.Sprite.Actor.MoveX(float64(rand.Intn(3)-1), nil)
	b.Sprite.Actor.MoveY(-1, nil)
	return nil
}
