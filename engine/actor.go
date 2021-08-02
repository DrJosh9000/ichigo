package engine

import "math"

// Thorson-style movement:
// https://maddythorson.medium.com/celeste-and-towerfall-physics-d24bd2ae0fc5

const collide = false // TODO: add collision detection

type Actor struct {
	X, Y int

	game       *Game
	xRem, yRem float64
}

func (a *Actor) MoveX(dx float64, onCollide func()) {
	a.xRem += dx
	move := int(math.Round(a.xRem))
	if move == 0 {
		return
	}
	a.xRem -= float64(move)
	sign := sign(move)
	for move != 0 {
		if collide {
			if onCollide != nil {
				onCollide()
			}
			return
		}
		a.X += sign
		move -= sign
	}
}

func (a *Actor) MoveY(dy float64, onCollide func()) {
	a.yRem += dy
	move := int(math.Round(a.yRem))
	if move == 0 {
		return
	}
	a.yRem -= float64(move)
	sign := sign(move)
	for move != 0 {
		if collide {
			if onCollide != nil {
				onCollide()
			}
			return
		}
		a.Y += sign
		move -= sign
	}
}

func (a *Actor) Build(g *Game) {
	a.game = g
}

func sign(m int) int {
	if m < 0 {
		return -1
	}
	return 1
}
