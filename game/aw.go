package game

import (
	"encoding/gob"
	"image"
	"math"

	"drjosh.dev/gurgle/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	_ engine.Identifier  = &Awakeman{}
	_ engine.Drawer      = &Awakeman{} // provided by Sprite
	_ engine.DrawOrderer = &Awakeman{} // provided by Sprite
	_ engine.Disabler    = &Awakeman{}
	_ engine.Hider       = &Awakeman{} // provided by Sprite
	_ engine.Prepper     = &Awakeman{}
	_ engine.Scanner     = &Awakeman{}
	_ engine.Updater     = &Awakeman{}
)

func init() {
	gob.Register(&Awakeman{})
}

type Awakeman struct {
	engine.Disabled
	engine.Sprite

	CameraID string
	ToastID  string

	camera      *engine.Camera
	toast       *engine.DebugToast
	vx, vy      float64
	facingLeft  bool
	coyoteTimer int
	jumpBuffer  int
	noclip      bool

	animIdleLeft, animIdleRight, animRunLeft, animRunRight, animWalkLeft, animWalkRight *engine.Anim
}

// Ident returns "awakeman". There should be only one!
func (aw *Awakeman) Ident() string { return "awakeman" }

func (aw *Awakeman) Update() error {
	if aw.Disabled {
		return nil
	}

	// TODO: better cheat for noclip
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		aw.noclip = !aw.noclip
		if aw.toast != nil {
			if aw.noclip {
				aw.toast.Toast("noclip enabled")
			} else {
				aw.toast.Toast("noclip disabled")
			}
		}
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
	return aw.Sprite.Update()
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
	return nil
}

func (aw *Awakeman) Prepare(game *engine.Game) {
	aw.camera = game.Component(aw.CameraID).(*engine.Camera)
	aw.toast, _ = game.Component(aw.ToastID).(*engine.DebugToast)

	aw.animIdleLeft = &engine.Anim{Frames: []engine.AnimFrame{
		{Frame: 1, Duration: 60},
	}}
	aw.animIdleRight = &engine.Anim{Frames: []engine.AnimFrame{
		{Frame: 0, Duration: 60},
	}}
	aw.animRunLeft = &engine.Anim{Frames: []engine.AnimFrame{
		{Frame: 14, Duration: 3},
		{Frame: 15, Duration: 5},
		{Frame: 16, Duration: 3},
		{Frame: 17, Duration: 3},
	}}
	aw.animRunRight = &engine.Anim{Frames: []engine.AnimFrame{
		{Frame: 10, Duration: 3},
		{Frame: 11, Duration: 5},
		{Frame: 12, Duration: 3},
		{Frame: 13, Duration: 3},
	}}
	aw.animWalkRight = &engine.Anim{Frames: []engine.AnimFrame{
		{Frame: 2, Duration: 6},
		{Frame: 3, Duration: 6},
		{Frame: 4, Duration: 6},
		{Frame: 5, Duration: 6},
	}}
	aw.animWalkLeft = &engine.Anim{Frames: []engine.AnimFrame{
		{Frame: 6, Duration: 6},
		{Frame: 7, Duration: 6},
		{Frame: 8, Duration: 6},
		{Frame: 9, Duration: 6},
	}}
}

func (aw *Awakeman) Scan() []interface{} { return []interface{}{&aw.Sprite} }
