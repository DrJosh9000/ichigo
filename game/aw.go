package game

import (
	"encoding/gob"
	"image"
	"math"

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
	jumpBuffer  int
	noclip      bool

	animIdleLeft, animIdleRight, animRunLeft, animRunRight *engine.Anim
}

func (aw *Awakeman) Update() error {
	// TODO: better cheat for noclip
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		aw.noclip = !aw.noclip
	}
	upd := aw.realUpdate
	if aw.noclip {
		upd = aw.noclipUpdate
	}
	if err := upd(); err != nil {
		return err
	}

	// Update the camera
	aw.camera.Zoom = 1
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		aw.camera.Zoom = 2
	}
	// aw.Pos is top-left corner, so add half size to get centre
	aw.camera.Centre = aw.Pos.Add(aw.Size.Div(2))
	return nil
}

func (aw *Awakeman) noclipUpdate() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		aw.Pos.Y--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		aw.Pos.Y++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		aw.Pos.X--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		aw.Pos.X++
	}
	return nil
}

func (aw *Awakeman) realUpdate() error {
	const (
		ε              = 0.2
		restitution    = -0.3
		gravity        = 0.3
		airResistance  = -0.01 // ⇒ terminal velocity = 30
		jumpVelocity   = -4.2
		runVelocity    = 1.4
		coyoteTime     = 5
		jumpBufferTime = 5
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

	// Has traction?
	if aw.CollidesAt(aw.Pos.Add(image.Pt(0, 1))) {
		// Not falling.
		// Instantly decelerate (AW absorbs all kinetic E in legs, or something)
		if aw.jumpBuffer > 0 {
			// Tried to jump recently -- so jump
			aw.vy = jumpVelocity
			aw.jumpBuffer = 0
		} else {
			// Can jump now or soon.
			aw.vy = 0
			aw.coyoteTimer = coyoteTime
		}
	} else {
		// Falling. v = v_0 + a, and a = gravity + airResistance(v_0)
		aw.vy += gravity + airResistance*aw.vy
		if aw.coyoteTimer > 0 {
			aw.coyoteTimer--
		}
		if aw.jumpBuffer > 0 {
			aw.jumpBuffer--
		}
	}

	// Handle controls

	// NB: spacebar sometimes does things on web pages (scrolls down)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		// On ground or recently on ground?
		if aw.coyoteTimer > 0 {
			// Jump. One frame of v = jumpVelocity (ignoring any gravity already applied this tick).
			aw.vy = jumpVelocity
		} else {
			// Buffer the jump in case aw hits the ground soon.
			aw.jumpBuffer = jumpBufferTime
		}
	}
	// Left and right
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

	// s = (v_0 + v) / 2.
	aw.MoveX((ux+aw.vx)/2, nil)
	// For Y, on collision, bounce a little bit.
	// Does not apply to X because controls override it anyway.
	aw.MoveY((uy+aw.vy)/2, func() {
		aw.vy *= restitution
		if math.Abs(aw.vy) < ε {
			aw.vy = 0
		}
	})
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
