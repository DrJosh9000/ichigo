package game

import (
	"encoding/gob"
	"image"

	"drjosh.dev/gurgle/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func init() {
	gob.Register(Awakeman{})
}

type Awakeman struct {
	engine.Sprite

	vx, vy      float64
	facingLeft  bool
	coyoteTimer int

	animIdleLeft, animIdleRight, animRunLeft, animRunRight *engine.Anim
}

func (aw *Awakeman) Update() error {
	const (
		bounceDampen = 0.5
		gravity      = 0.3
		jumpVelocity = -3.5
		runVelocity  = 1.4
	)

	// Standing on something?
	if aw.CollidesAt(aw.Pos.Add(image.Pt(0, 1))) {
		// Not falling
		aw.vy = 0
		aw.coyoteTimer = 5
	} else {
		// Falling
		aw.vy += gravity
		if aw.coyoteTimer > 0 {
			aw.coyoteTimer--
		}
	}
	if aw.coyoteTimer > 0 && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Jump
		aw.vy = jumpVelocity
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA):
		aw.vx = -runVelocity
		aw.SetAnim(aw.animRunLeft)
		aw.facingLeft = true
	case ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD):
		aw.vx = runVelocity
		aw.SetAnim(aw.animRunRight)
		aw.facingLeft = false
	default:
		aw.vx = 0
		aw.SetAnim(aw.animIdleRight)
		if aw.facingLeft {
			aw.SetAnim(aw.animIdleLeft)
		}
	}
	aw.MoveX(aw.vx, func() { aw.vx = -aw.vx * bounceDampen })
	aw.MoveY(aw.vy, func() { aw.vy = -aw.vy * bounceDampen })
	return aw.Sprite.Update()
}

func (aw *Awakeman) Prepare(*engine.Game) {
	aw.animIdleLeft = &engine.Anim{Def: engine.AnimDefs["aw_idle_left"]}
	aw.animIdleRight = &engine.Anim{Def: engine.AnimDefs["aw_idle_right"]}
	aw.animRunLeft = &engine.Anim{Def: engine.AnimDefs["aw_run_left"]}
	aw.animRunRight = &engine.Anim{Def: engine.AnimDefs["aw_run_right"]}
}

func (aw *Awakeman) Scan() []interface{} { return []interface{}{&aw.Sprite} }
