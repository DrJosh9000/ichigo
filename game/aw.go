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

const awakemanProducesBubbles = true

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
	engine.Disables
	Sprite   engine.Sprite
	CameraID string
	ToastID  string

	game        *engine.Game
	camera      *engine.Camera
	toast       *engine.DebugToast
	vel         geom.Float3
	facingLeft  bool
	coyoteTimer int
	jumpBuffer  int
	noclip      bool
	spawnPoint  geom.Int3
	bubbleTimer int

	anims map[string]*engine.Anim
}

// Ident returns "awakeman". There should be only one!
func (aw *Awakeman) Ident() string { return "awakeman" }

func (aw *Awakeman) Update() error {
	// TODO: better cheat for noclip
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		aw.noclip = !aw.noclip
		aw.vel = geom.Float3{}
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
	aw.camera.PointAt(aw.Sprite.Actor.BoundingBox().Centre(), z)
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
		respawnY       = 1000
		bubblePeriod   = 6
	)

	if awakemanProducesBubbles {
		// Add a bubble?
		aw.bubbleTimer--
		if aw.bubbleTimer <= 0 {
			aw.bubbleTimer = bubblePeriod
			bubble := NewBubble(aw.Sprite.Actor.Pos.Add(geom.Pt3(-3, -20, -1)))
			if err := engine.PreorderWalk(bubble, func(c, _ interface{}) error {
				if p, ok := c.(engine.Loader); ok {
					return p.Load(Assets)
				}
				return nil
			}); err != nil {
				return err
			}
			// Add bubble to same parent as aw
			par := aw.game.Parent(aw)
			aw.game.PathRegister(bubble, par)
			if err := engine.PostorderWalk(bubble, func(c, _ interface{}) error {
				if p, ok := c.(engine.Prepper); ok {
					return p.Prepare(aw.game)
				}
				return nil
			}); err != nil {
				return err
			}
			bubble.Sprite.SetAnim(bubble.Sprite.Sheet.NewAnim("bubble"))
		}
	}

	// Fell below some threshold?
	if aw.Sprite.Actor.Pos.Y > respawnY {
		aw.Sprite.Actor.Pos = aw.spawnPoint
		aw.vel = geom.Float3{}
	}

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
	v0 := aw.vel

	// Has traction?
	if aw.vel.Y >= 0 && aw.Sprite.Actor.CollidesAt(aw.Sprite.Actor.Pos.Add(geom.Pt3(0, 1, 0))) {
		// Not falling.
		// Instantly decelerate (AW absorbs all kinetic E in legs, or something)
		if aw.jumpBuffer > 0 {
			// Tried to jump recently -- so jump
			aw.vel.Y = jumpVelocity
			aw.jumpBuffer = 0
		} else {
			// Can jump now or soon.
			aw.vel.Y = 0
			aw.coyoteTimer = coyoteTime
		}
	} else {
		// Falling. v = v_0 + a, and a = gravity + airResistance(v_0)
		aw.vel.Y += gravity + airResistance*aw.vel.Y
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
			aw.vel.Y = jumpVelocity
		} else {
			// Buffer the jump in case aw hits the ground soon.
			aw.jumpBuffer = jumpBufferTime
		}
	}
	// Left, right, away, toward
	aw.vel.X, aw.vel.Z = 0, 0
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyJ):
		aw.vel.X = -runVelocity
	case ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyL):
		aw.vel.X = runVelocity
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyI):
		aw.vel.Z = -runVelocity
	case ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyK):
		aw.vel.Z = runVelocity
	}

	// Animations and velocity correction
	switch {
	case aw.vel.X != 0 && aw.vel.Z != 0: // Diagonal
		aw.Sprite.SetAnim(aw.anims["run_vert"])
		// Pythagorean theorem; |vx| = |vz|, so the hypotenuse is √2 too big
		// if we want to run at runVelocity always
		aw.vel.X /= sqrt2
		aw.vel.Z /= sqrt2
	case aw.vel.X == 0 && aw.vel.Z != 0: // Vertical
		aw.Sprite.SetAnim(aw.anims["run_vert"])

	// vz == 0 for all remaining cases
	case aw.vel.X < 0: // Left
		aw.Sprite.SetAnim(aw.anims["run_left"])
		aw.facingLeft = true
	case aw.vel.X > 0: // Right
		aw.Sprite.SetAnim(aw.anims["run_right"])
		aw.facingLeft = false
	default: // aw.velocity.X == 0; Idle
		aw.Sprite.SetAnim(aw.anims["idle_right"])
		if aw.facingLeft {
			aw.Sprite.SetAnim(aw.anims["idle_left"])
		}
	}

	// s = (v_0 + v) / 2.
	aw.Sprite.Actor.MoveX((v0.X+aw.vel.X)/2, nil)
	// For Y, on collision from going upwards, bounce a little bit.
	// Does not apply to X because controls override it anyway.
	aw.Sprite.Actor.MoveY((v0.Y+aw.vel.Y)/2, func() {
		if aw.vel.Y > 0 {
			return
		}
		aw.vel.Y *= restitution
		if math.Abs(aw.vel.Y) < ε {
			aw.vel.Y = 0
		}
	})
	aw.Sprite.Actor.MoveZ((v0.Z+aw.vel.Z)/2, nil)
	return nil
}

func (aw *Awakeman) Prepare(game *engine.Game) error {
	aw.game = game
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
	aw.spawnPoint = aw.Sprite.Actor.Pos

	return nil
}

//func (aw *Awakeman) Scan() []interface{} { return []interface{}{&aw.Sprite} }
func (aw *Awakeman) Scan(visit func(interface{}) error) error {
	return visit(&aw.Sprite)
}

func (aw *Awakeman) String() string {
	return fmt.Sprintf("Awakeman@%v", aw.Sprite.Actor.Pos)
}
