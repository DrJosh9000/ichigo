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

	CameraID string

	camera      *engine.Camera
	vx, vy      float64
	facingLeft  bool
	coyoteTimer int

	animIdleLeft, animIdleRight, animRunLeft, animRunRight *engine.Anim
}

func (aw *Awakeman) Update() error {
	const (
		bounceDampen = 0.5
		gravity      = 0.3
		jumpVelocity = -3.6
		runVelocity  = 1.4
	)

	// High-school physics time! Under constant acceleration:
	//   v = v_0 + a*t
	// and
	//   s = t * (v_0 + v) / 2
	// (note t is in ticks and s is in world units)
	// and since we get one Update per tick (t = 1),
	//   v = v_0 + a,
	// and
	//   s = (v_0 + v) / 2.
	// Capture current v_0 to use later.
	ux, uy := aw.vx, aw.vy

	// Standing on something?
	if aw.CollidesAt(aw.Pos.Add(image.Pt(0, 1))) {
		// Not falling. Let's assume aw always lands safely.
		// Setting a = -v_0 gives v = v_0 - v_0 = 0.
		aw.vy = 0
		aw.coyoteTimer = 5
	} else {
		// Falling. v = v_0 + a, and a is gravity.
		aw.vy += gravity
		if aw.coyoteTimer > 0 {
			aw.coyoteTimer--
		}
	}

	// Handle controls

	// NB: spacebar sometimes does things on web pages (scrolls down)
	if aw.coyoteTimer > 0 && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyZ)) {
		// Jump. One frame of a = jumpVelocity (ignoring any gravity already applied this tick).
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
	aw.camera.Zoom = 1
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		aw.camera.Zoom = 2
	}
	/*
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			aw.camera.Rotation += math.Pi / 6
		}
	*/

	// s = (v_0 + v) / 2.
	// On collision, bounce a little bit.
	aw.MoveX((ux+aw.vx)/2, func() { aw.vx = -aw.vx * bounceDampen })
	aw.MoveY((uy+aw.vy)/2, func() { aw.vy = -aw.vy * bounceDampen })
	// aw.Pos is top-left corner, so add half size to get centre
	aw.camera.Centre = aw.Pos.Add(aw.Size.Div(2))
	return aw.Sprite.Update()
}

func (aw *Awakeman) Prepare(game *engine.Game) {
	aw.camera = game.Component(aw.CameraID).(*engine.Camera)

	aw.animIdleLeft = &engine.Anim{Def: engine.AnimDefs["aw_idle_left"]}
	aw.animIdleRight = &engine.Anim{Def: engine.AnimDefs["aw_idle_right"]}
	aw.animRunLeft = &engine.Anim{Def: engine.AnimDefs["aw_run_left"]}
	aw.animRunRight = &engine.Anim{Def: engine.AnimDefs["aw_run_right"]}
}

func (aw *Awakeman) Scan() []interface{} { return []interface{}{&aw.Sprite} }
