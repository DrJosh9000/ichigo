package engine

// Anim is an "instance" of an AnimDef: an animation being displayed,
// together with the current state.
type Anim struct {
	Def          *AnimDef
	CurrentIndex int
	CurrentTicks int
}

func (a *Anim) CurrentFrame() int { return a.Def.Frames[a.CurrentIndex].Frame }

// Update increments the tick count and advances the frame if necessary.
func (a *Anim) Update() error {
	a.CurrentTicks++
	if !a.Def.Loop && a.CurrentIndex == len(a.Def.Frames)-1 {
		// on the last frame of a one shot so remain on final frame
		return nil
	}
	if a.CurrentTicks >= a.Def.Frames[a.CurrentIndex].DurationTicks {
		a.CurrentTicks = 0
		a.CurrentIndex++
	}
	if a.Def.Loop && a.CurrentIndex >= len(a.Def.Frames) {
		a.CurrentIndex = 0
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
	Frame         int `json:"frame"`    // show this frame
	DurationTicks int `json:"duration"` // for this long
}
