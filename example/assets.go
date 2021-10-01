package example

// TODO: make an io/fs wrapper for ebitenutil.OpenFile ?

import "embed"

//go:embed assets
var Assets embed.FS
