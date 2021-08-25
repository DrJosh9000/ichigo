//go:build js
// +build js

package game

// TODO: make an io/fs wrapper for ebitenutil.OpenFile ?

import "embed"

//go:embed assets
var Assets embed.FS
