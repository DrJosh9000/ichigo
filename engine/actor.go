package engine

import (
	"encoding/gob"
	"image"
)

// Ensure Actor satisfies interfaces.
var _ Prepper = &Actor{}

func init() {
	gob.Register(Actor{})
}

// Thorson-style movement:
// https://maddythorson.medium.com/celeste-and-towerfall-physics-d24bd2ae0fc5

// Actor handles basic movement.
type Actor struct {
	CollisionDomain string
	Pos             image.Point
	Size            image.Point

	collisionDomain interface{}
	xRem, yRem      float64
}

func (a *Actor) CollidesAt(p image.Point) bool {
	hit := false
	Walk(a.collisionDomain, func(c interface{}) bool {
		if coll, ok := c.(Collider); ok {
			if coll.CollidesWith(image.Rectangle{Min: p, Max: p.Add(a.Size)}) {
				hit = true
				return false
			}
		}
		return true
	})
	return hit
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
		if a.CollidesAt(a.Pos.Add(image.Pt(sign, 0))) {
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
		if a.CollidesAt(a.Pos.Add(image.Pt(0, sign))) {
			if onCollide != nil {
				onCollide()
			}
			return
		}
		a.Pos.Y += sign
		move -= sign
	}
}

func (a *Actor) Prepare(g *Game) {
	a.collisionDomain = g.Component(a.CollisionDomain)
}

func sign(m int) int {
	if m < 0 {
		return -1
	}
	return 1
}
