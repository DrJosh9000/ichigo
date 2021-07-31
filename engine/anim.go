package engine

// Anim is an "instance" of an AnimDef: an animation being displayed,
// together with the current state.
type Anim struct {
	Def   *AnimDef
	Index int
	Ticks int
}

func (a *Anim) CurrentFrame() int { return a.Def.Frames[a.Index].Frame }

// Update increments the tick count and advances the frame if necessary.
func (a *Anim) Update() error {
	a.Ticks++
	if !a.Def.Loop && a.Index == len(a.Def.Frames)-1 {
		// on the last frame of a one shot so remain on final frame
		return nil
	}
	if a.Ticks >= a.Def.Frames[a.Index].Duration {
		a.Ticks = 0
		a.Index++
	}
	if a.Def.Loop && a.Index >= len(a.Def.Frames) {
		a.Index = 0
	}
	return nil
}

// AnimDef describes an animation (sequence of frames and timings).
type AnimDef struct {
	Frames []AnimFrame `json:"frames"`
	Loop   bool        `json:"loop"`
}

// AnimFrame describes a frame in an animation.
type AnimFrame struct {
	Frame    int `json:"frame"`    // show this frame
	Duration int `json:"duration"` // for this long, in ticks
}
