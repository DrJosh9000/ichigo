package game

import "drjosh.dev/gurgle/engine"

func init() {
	engine.AnimDefs = map[string]*engine.AnimDef{
		"green_tiles": {
			Frames: []engine.AnimFrame{
				{Frame: 0, Duration: 16},
				{Frame: 1, Duration: 16},
				{Frame: 2, Duration: 16},
			},
		},
		"red_tiles": {
			Frames: []engine.AnimFrame{
				{Frame: 3, Duration: 12},
				{Frame: 4, Duration: 12},
				{Frame: 5, Duration: 12},
				{Frame: 6, Duration: 12},
			},
		},
		"aw_idle_right": {
			Frames: []engine.AnimFrame{
				{Frame: 0, Duration: 60},
			},
		},
		"aw_idle_left": {
			Frames: []engine.AnimFrame{
				{Frame: 1, Duration: 60},
			},
		},
		"aw_walk_right": {
			Frames: []engine.AnimFrame{
				{Frame: 2, Duration: 6},
				{Frame: 3, Duration: 6},
				{Frame: 4, Duration: 6},
				{Frame: 5, Duration: 6},
			},
		},
		"aw_walk_left": {
			Frames: []engine.AnimFrame{
				{Frame: 6, Duration: 6},
				{Frame: 7, Duration: 6},
				{Frame: 8, Duration: 6},
				{Frame: 9, Duration: 3},
			},
		},
		"aw_run_right": {
			Frames: []engine.AnimFrame{
				{Frame: 10, Duration: 3},
				{Frame: 11, Duration: 3},
				{Frame: 12, Duration: 3},
				{Frame: 13, Duration: 3},
			},
		},
		"aw_run_left": {
			Frames: []engine.AnimFrame{
				{Frame: 14, Duration: 3},
				{Frame: 15, Duration: 3},
				{Frame: 16, Duration: 3},
				{Frame: 17, Duration: 3},
			},
		},
	}
}
