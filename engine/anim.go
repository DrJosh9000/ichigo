package engine

import "encoding/gob"

// Ensure Anim satisfies Animer.
var _ interface {
	Cell() int
	Reset()
	Updater
} = &Anim{}

func init() {
	gob.Register(&Anim{})
}

// AnimDef defines an animation, as a sequence of steps and other information.
type AnimDef struct {
	Steps   []AnimStep
	OneShot bool
}

// NewAnim spawns a new anim using this def, or nil if d is nil.
func (d *AnimDef) NewAnim() *Anim {
	if d == nil {
		return nil
	}
	return &Anim{Def: d}
}

// AnimStep describes a step in an animation.
type AnimStep struct {
	Cell     int // show this cell
	Duration int // for this long, in ticks
}

// Anim is the current state of an animation being played (think of it as an
// instance of an AnimDef). nil *Anim can be used, but always returns 0 for the
// current frame.
type Anim struct {
	Def   *AnimDef
	Index int // current step index
	Ticks int // ticks spent at this step
}

// Cell returns the cell index for the current step.
func (a *Anim) Cell() int {
	if a == nil {
		return 0
	}
	return a.Def.Steps[a.Index].Cell
}

// Reset resets both Index and Ticks to 0.
func (a *Anim) Reset() {
	if a == nil {
		return
	}
	a.Index, a.Ticks = 0, 0
}

// Update increments the tick count and advances the frame if necessary.
func (a *Anim) Update() error {
	if a == nil {
		return nil
	}
	a.Ticks++
	if a.Def.OneShot && a.Index == len(a.Def.Steps)-1 {
		// on the last frame of a one shot so remain on final frame
		return nil
	}
	if a.Ticks >= a.Def.Steps[a.Index].Duration {
		a.Ticks = 0
		a.Index++
	}
	if !a.Def.OneShot && a.Index >= len(a.Def.Steps) {
		a.Index = 0
	}
	return nil
}
