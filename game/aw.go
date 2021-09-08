package game

import (
	"encoding/gob"
	"fmt"
	"math"

	"drjosh.dev/gurgle/engine"
	"drjosh.dev/gurgle/geom"
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

	anims map[string]*engine.Anim
}

// Ident returns "awakeman". There should be only one!
func (aw *Awakeman) Ident() string { return "awakeman" }

func (aw *Awakeman) Update() error {
	// TODO: better cheat for noclip
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		aw.noclip = !aw.noclip
		aw.vx, aw.vy, aw.vz = 0, 0, 0
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
	pos := aw.Sprite.Actor.Pos
	size := aw.Sprite.Actor.Size
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
		gravity        = 0.25
		airResistance  = -0.005 // ⇒ terminal velocity = 10
		jumpVelocity   = -3.3
		sqrt2          = 1.414213562373095
		runVelocity    = sqrt2
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
	if aw.vy >= 0 && aw.Sprite.Actor.CollidesAt(aw.Sprite.Actor.Pos.Add(geom.Pt3(0, 1, 0))) {
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
	// Left, right, away, toward
	aw.vx, aw.vz = 0, 0
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		aw.vx = -runVelocity
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		aw.vx = runVelocity
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		aw.vz = -runVelocity
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		aw.vz = runVelocity
	}

	// Animations and velocity correction
	switch {
	case aw.vx != 0 && aw.vz != 0: // Diagonal
		aw.Sprite.SetAnim(aw.anims["run_vert"])
		// Pythagorean theorem; |vx| = |vz|, so the hypotenuse is √2 too big
		// if we want to run at runVelocity always
		aw.vx /= sqrt2
		aw.vz /= sqrt2
	case aw.vx == 0 && aw.vz != 0: // Vertical
		aw.Sprite.SetAnim(aw.anims["run_vert"])

	// vz == 0 for all remaining cases
	case aw.vx < 0: // Left
		aw.Sprite.SetAnim(aw.anims["run_left"])
		aw.facingLeft = true
	case aw.vx > 0: // Right
		aw.Sprite.SetAnim(aw.anims["run_right"])
		aw.facingLeft = false
	default: // aw.vx == 0; Idle
		aw.Sprite.SetAnim(aw.anims["idle_right"])
		if aw.facingLeft {
			aw.Sprite.SetAnim(aw.anims["idle_left"])
		}
	}

	// s = (v_0 + v) / 2.
	aw.Sprite.Actor.MoveX((ux+aw.vx)/2, nil)
	// For Y, on collision from going upwards, bounce a little bit.
	// Does not apply to X because controls override it anyway.
	aw.Sprite.Actor.MoveY((uy+aw.vy)/2, func() {
		if aw.vy > 0 {
			return
		}
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
	aw.anims = aw.Sprite.Sheet.NewAnims()

	/*
		idle_left
		idle_right
		run_left
		run_right
		run_vert
	*/

	return nil
}

func (aw *Awakeman) Scan() []interface{} { return []interface{}{&aw.Sprite} }
