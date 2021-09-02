package game

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math"

	"drjosh.dev/gurgle/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var _ interface {
	engine.Identifier
	engine.Disabler
	engine.Prepper
	engine.Scanner
	engine.Updater
} = &Awakeman{}

func init() {
	gob.Register(&Awakeman{})
}

// Awakeman is a bit of a god object for now...
type Awakeman struct {
	engine.Disabled
	Sprite   engine.Sprite
	CameraID string
	ToastID  string

	camera      *engine.Camera
	toast       *engine.DebugToast
	vx, vy, vz  float64
	facingLeft  bool
	coyoteTimer int
	jumpBuffer  int
	noclip      bool

	animIdleLeft, animIdleRight, animRunLeft, animRunRight, animWalkLeft, animWalkRight *engine.Anim
}

// Ident returns "awakeman". There should be only one!
func (aw *Awakeman) Ident() string { return "awakeman" }

func (aw *Awakeman) Update() error {
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
	// aw.Pos is top-left corner, so add half size to get centre
	z := 1.0
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		z = 2.0
	}
	pos := aw.Sprite.Actor.Pos.XY()
	size := aw.Sprite.Actor.Size.XY()
	aw.camera.PointAt(pos.Add(size.Div(2)), z)
	return nil
}

func (aw *Awakeman) noclipUpdate() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		aw.Sprite.Actor.Pos.Y--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		aw.Sprite.Actor.Pos.Y++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		aw.Sprite.Actor.Pos.X--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		aw.Sprite.Actor.Pos.X++
	}
	return nil
}

func (aw *Awakeman) realUpdate() error {
	const (
		ε              = 0.2
		restitution    = -0.3
		gravity        = 0.2
		airResistance  = -0.02 // ⇒ terminal velocity = 10
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
	ux, uy, uz := aw.vx, aw.vy, aw.vz

	// Has traction?
	if aw.Sprite.Actor.CollidesAt(aw.Sprite.Actor.Pos.Add(engine.Pt3(0, 1, 0))) {
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
		aw.Sprite.SetAnim(aw.animRunLeft)
		aw.facingLeft = true
	case ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD):
		aw.vx = runVelocity
		aw.Sprite.SetAnim(aw.animRunRight)
		aw.facingLeft = false
	default:
		aw.vx = 0
		aw.Sprite.SetAnim(aw.animIdleRight)
		if aw.facingLeft {
			aw.Sprite.SetAnim(aw.animIdleLeft)
		}
	}
	// Up and down (away and closer)
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW):
		aw.vz = -runVelocity
	case ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS):
		aw.vz = runVelocity
	default:
		aw.vz = 0
	}

	// s = (v_0 + v) / 2.
	aw.Sprite.Actor.MoveX((ux+aw.vx)/2, nil)
	// For Y, on collision, bounce a little bit.
	// Does not apply to X because controls override it anyway.
	aw.Sprite.Actor.MoveY((uy+aw.vy)/2, func() {
		aw.vy *= restitution
		if math.Abs(aw.vy) < ε {
			aw.vy = 0
		}
	})
	aw.Sprite.Actor.MoveZ((uz+aw.vz)/2, nil)
	return nil
}

func (aw *Awakeman) Prepare(game *engine.Game) error {
	cam, ok := game.Component(aw.CameraID).(*engine.Camera)
	if !ok {
		return fmt.Errorf("component %q not *engine.Camera", aw.CameraID)
	}
	aw.camera = cam
	tst, ok := game.Component(aw.ToastID).(*engine.DebugToast)
	if !ok {
		return fmt.Errorf("component %q not *engine.DebugToast", aw.ToastID)
	}
	aw.toast = tst

	aw.animIdleLeft = aw.Sprite.Sheet.NewAnim("idle_left")
	if aw.animIdleLeft == nil {
		return errors.New("missing anim idle_left")
	}
	aw.animIdleRight = aw.Sprite.Sheet.NewAnim("idle_right")
	if aw.animIdleRight == nil {
		return errors.New("missing anim idle_right")
	}
	aw.animRunLeft = aw.Sprite.Sheet.NewAnim("run_left")
	if aw.animRunLeft == nil {
		return errors.New("missing anim run_left")
	}
	aw.animRunRight = aw.Sprite.Sheet.NewAnim("run_right")
	if aw.animRunRight == nil {
		return errors.New("missing anim run_right")
	}
	aw.animWalkRight = aw.Sprite.Sheet.NewAnim("walk_left")
	if aw.animWalkRight == nil {
		return errors.New("missing anim walk_left")
	}
	aw.animWalkLeft = aw.Sprite.Sheet.NewAnim("walk_right")
	if aw.animWalkLeft == nil {
		return errors.New("missing anim walk_right")
	}
	return nil
}

func (aw *Awakeman) Scan() []interface{} { return []interface{}{&aw.Sprite} }
