//go:build js
// +build js

package game

// TODO: make io/fs wrapper for ebitenutil.OpenFile ?

import "embed"

//go:embed assets
var Assets embed.FS
