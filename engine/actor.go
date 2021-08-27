package engine

import (
	"encoding/gob"
	"fmt"
	"image"
)

// Ensure Actor satisfies interfaces.
var _ interface {
	Bounder
	Prepper
} = &Actor{}

func init() {
	gob.Register(&Actor{})
}

// Thorson-style movement:
// https://maddythorson.medium.com/celeste-and-towerfall-physics-d24bd2ae0fc5

// Actor handles basic movement.
type Actor struct {
	CollisionDomain string
	Pos             image.Point
	Size            image.Point

	collisionDomain []Collider
	xRem, yRem      float64
}

func (a *Actor) BoundingRect() image.Rectangle { return image.Rectangle{a.Pos, a.Pos.Add(a.Size)} }

func (a *Actor) CollidesAt(p image.Point) bool {
	bounds := image.Rectangle{Min: p, Max: p.Add(a.Size)}
	/*return nil != Walk(a.collisionDomain, func(c interface{}, _ []interface{}) error {
		coll, ok := c.(Collider)
		if !ok {
			return nil
		}
		if coll.CollidesWith(bounds) {
			return Collision{With: coll}
		}
		return nil
	})*/
	for _, c := range a.collisionDomain {
		if c.CollidesWith(bounds) {
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

func (a *Actor) Prepare(g *Game) error {
	//a.collisionDomain = g.Component(a.CollisionDomain)

	cs := g.Query(g.Component(a.CollisionDomain), ColliderType)
	a.collisionDomain = make([]Collider, 0, len(cs))
	for _, c := range cs {
		a.collisionDomain = append(a.collisionDomain, c.(Collider))
	}
	return nil
}

func sign(m int) int {
	if m < 0 {
		return -1
	}
	return 1
}

// Collision reports a collision occurred.
type Collision struct {
	With Collider
}

// Error is really only to implement the error interface.
func (c Collision) Error() string {
	return fmt.Sprintf("collision with %v", c.With)
}
