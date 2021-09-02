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
	Pos, Size        Point3
	xRem, yRem, zRem float64

	game *Game
}

func (a *Actor) CollidesAt(p Point3) bool {
	bounds := Box{Min: p, Max: p.Add(a.Size)}.XY() // TODO: 3D collision
	for c := range a.game.Query(a.CollisionDomain, ColliderType) {
		if c.(Collider).CollidesWith(bounds) {
			return true
		}
	}
	return false
}

func (a *Actor) MoveX(dx float64, onCollide func()) {
	a.xRem += dx
	move := int(a.xRem + 0.5) // Note: math.Round can lead to vibration
	if move == 0 {
		return
	}
	a.xRem -= float64(move)
	sign := sign(move)
	for move != 0 {
		if a.CollidesAt(a.Pos.Add(Pt3(sign, 0, 0))) {
			if onCollide != nil {
				onCollide()
			}
			return
		}
		a.Pos.X += sign
		move -= sign
	}
}

func (a *Actor) MoveY(dy float64, onCollide func()) {
	a.yRem += dy
	move := int(a.yRem + 0.5)
	if move == 0 {
		return
	}
	a.yRem -= float64(move)
	sign := sign(move)
	for move != 0 {
		if a.CollidesAt(a.Pos.Add(Pt3(0, sign, 0))) {
			if onCollide != nil {
				onCollide()
			}
			return
		}
		a.Pos.Y += sign
		move -= sign
	}
}

func (a *Actor) MoveZ(dz float64, onCollide func()) {
	a.zRem += dz
	move := int(a.zRem + 0.5)
	if move == 0 {
		return
	}
	a.zRem -= float64(move)
	sign := sign(move)
	for move != 0 {
		if a.CollidesAt(a.Pos.Add(Pt3(0, 0, sign))) {
			if onCollide != nil {
				onCollide()
			}
			return
		}
		a.Pos.Z += sign
		move -= sign
	}
}

func (a *Actor) Prepare(g *Game) error {
	a.game = g
	return nil
}

func sign(m int) int {
	if m < 0 {
		return -1
	}
	return 1
}
