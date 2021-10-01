package example

import (
	"fmt"
	"image"
	"math/rand"

	"github.com/DrJosh9000/ichigo/engine"
	"github.com/DrJosh9000/ichigo/geom"
)

var _ interface {
	engine.Scanner
	engine.Prepper
	engine.Updater
} = &Bubble{}

type Bubble struct {
	Life   int
	Sprite engine.Sprite

	game *engine.Game
}

func NewBubble(pos geom.Int3) *Bubble {
	return &Bubble{
		Life: 60,
		Sprite: engine.Sprite{
			Actor: engine.Actor{
				CollisionDomain: "level_1",
				Pos:             pos,
				Bounds: geom.Box{
					Min: geom.Pt3(-4, -4, -4),
					Max: geom.Pt3(4, 4, 4),
				},
			},
			DrawOffset: image.Pt(-4, -4),
			Sheet: engine.Sheet{
				AnimDefs: map[string]*engine.AnimDef{
					"bubble": {
						Steps: []engine.AnimStep{
							{Cell: 0, Duration: 5},
							{Cell: 1, Duration: 15},
							{Cell: 2, Duration: 20},
							{Cell: 3, Duration: 15},
							{Cell: 4, Duration: 3},
							{Cell: 5, Duration: 2},
						},
						OneShot: true,
					},
				},
				CellSize: image.Pt(8, 8),
				Src:      engine.ImageRef{Path: "assets/bubble.png"},
			},
		},
	}
}

func (b *Bubble) Scan(visit engine.VisitFunc) error {
	return visit(&b.Sprite)
}

func (b *Bubble) String() string {
	return fmt.Sprintf("Bubble@%v", b.Sprite.Actor.Pos)
}

func (b *Bubble) Prepare(g *engine.Game) error {
	b.game = g
	return nil
}

func (b *Bubble) Update() error {
	b.Life--
	if b.Life <= 0 {
		b.game.PathUnregister(b)
	}
	die := func() { b.Life = 0 }
	b.Sprite.Actor.MoveX(float64(rand.Intn(3)-1), die)
	b.Sprite.Actor.MoveY(-1, die)
	//lint:ignore SA4000 one random minus another is not always zero...
	b.Sprite.Actor.MoveZ(float64(rand.Intn(2)-rand.Intn(2)), die)
	return nil
}
