package engine

// Ensure Anim satisfies Animer.
var _ Animer = &Anim{}

// AnimFrame describes a frame in an animation.
type AnimFrame struct {
	Frame    int // show this frame
	Duration int // for this long, in ticks
}

// Anim is n animation being displayed, together with the current state.
type Anim struct {
	Frames  []AnimFrame
	OneShot bool
	Index   int
	Ticks   int
}

// Copy makes a shallow copy of the anim.
func (a *Anim) Copy() *Anim {
	a2 := *a
	return &a2
}

// CurrentFrame returns the frame number for the current index.
func (a *Anim) CurrentFrame() int { return a.Frames[a.Index].Frame }

// Reset resets both Index and Ticks to 0.
func (a *Anim) Reset() { a.Index, a.Ticks = 0, 0 }

// Update increments the tick count and advances the frame if necessary.
func (a *Anim) Update() error {
	a.Ticks++
	if a.OneShot && a.Index == len(a.Frames)-1 {
		// on the last frame of a one shot so remain on final frame
		return nil
	}
	if a.Ticks >= a.Frames[a.Index].Duration {
		a.Ticks = 0
		a.Index++
	}
	if !a.OneShot && a.Index >= len(a.Frames) {
		a.Index = 0
	}
	return nil
}
