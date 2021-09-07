package engine

import (
	"encoding/gob"
)

// Ensure Actor satisfies interfaces.
var _ Prepper = &Actor{}

func init() {
	gob.Register(&Actor{})
}

// Thorson-style movement:
// https://maddythorson.medium.com/celeste-and-towerfall-physics-d24bd2ae0fc5

// Actor handles basic movement.
type Actor struct {
	CollisionDomain  string // id of component to look for colliders inside of
	Pos, Size        Int3
	xRem, yRem, zRem float64

	game *Game
}

func (a *Actor) CollidesAt(p Int3) bool {
	bounds := Box{Min: p, Max: p.Add(a.Size)}
	for c := range a.game.Query(a.CollisionDomain, ColliderType) {
		if c.(Collider).CollidesWith(bounds) {
			return true
		}
	}
	return false
}

func (a *Actor) MoveX(x float64, onCollide func()) {
	a.xRem += x
	move := int(a.xRem + 0.5) // Note: math.Round can lead to vibration
	if move == 0 {
		return
	}
	a.xRem -= float64(move)
	sign := sign(move)
	for move != 0 {
		a.Pos.X += sign
		move -= sign
		if !a.CollidesAt(a.Pos) {
			continue
		}
		if onCollide != nil {
			onCollide()
		}
		a.Pos.X -= sign
		a.xRem = 0
		return
	}
}

func (a *Actor) MoveY(y float64, onCollide func()) {
	a.yRem += y
	move := int(a.yRem + 0.5)
	if move == 0 {
		return
	}
	a.yRem -= float64(move)
	sign := sign(move)
	for move != 0 {
		a.Pos.Y += sign
		move -= sign
		if !a.CollidesAt(a.Pos) {
			continue
		}
		if onCollide != nil {
			onCollide()
		}
		a.Pos.Y -= sign
		a.yRem = 0
		return
	}
}

func (a *Actor) MoveZ(z float64, onCollide func()) {
	a.zRem += z
	move := int(a.zRem + 0.5)
	if move == 0 {
		return
	}
	a.zRem -= float64(move)
	sign := sign(move)
	for move != 0 {
		a.Pos.Z += sign
		move -= sign
		if !a.CollidesAt(a.Pos) {
			continue
		}
		if onCollide != nil {
			onCollide()
		}
		a.Pos.Z -= sign
		a.zRem = 0
		return
	}
}

func (a *Actor) Prepare(g *Game) error {
	a.game = g
	return nil
}
