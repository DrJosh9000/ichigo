package engine

import (
	"encoding/gob"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const gravity = 0.5

func init() {
	gob.Register(Actor{})
}

// Thorson-style movement:
// https://maddythorson.medium.com/celeste-and-towerfall-physics-d24bd2ae0fc5

// Collider components have tangible form.
type Collider interface {
	CollidesWith(image.Rectangle) bool
}

type Actor struct {
	ID
	Position image.Point
	Size     image.Point
	ZPos     // TODO: refactor

	game       *Game
	xRem, yRem float64
	src        *ebiten.Image // TODO: refactor
	vx, vy     float64       // TODO: refactor
}

func (a *Actor) collidesAt(p image.Point) bool {
	// TODO: more efficient test?
	hit := false
	a.game.Walk(func(c interface{}) bool {
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

func (a *Actor) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	// TODO: delegate drawing to something else
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(float64(a.Position.X), float64(a.Position.Y))
	op.GeoM.Concat(geom)
	screen.DrawImage(a.src, &op)
}

func (a *Actor) Update() error {
	// TODO: delegate updating to something else
	if a.collidesAt(a.Position.Add(image.Pt(0, 1))) {
		// Not falling
		a.vy = 0
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// Jump?
			a.vy = -7
		}
	} else {
		// Falling
		a.vy += gravity
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		a.vx = -3
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		a.vx = 3
	default:
		a.vx = 0
	}
	a.MoveX(a.vx, func() {
		a.vx = 0
	})
	a.MoveY(a.vy, func() {
		a.vy = 0
	})
	return nil
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
		if a.collidesAt(a.Position.Add(image.Pt(sign, 0))) {
			if onCollide != nil {
				onCollide()
			}
			return
		}
		a.Position.X += sign
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
		if a.collidesAt(a.Position.Add(image.Pt(0, sign))) {
			if onCollide != nil {
				onCollide()
			}
			return
		}
		a.Position.Y += sign
		move -= sign
	}
}

func (a *Actor) Build(g *Game) {
	a.game = g

	// TODO: remove hack temporary image
	a.src = ebiten.NewImage(16, 16)
	a.src.Fill(color.White)
}

func sign(m int) int {
	if m < 0 {
		return -1
	}
	return 1
}
