//go:build !js
// +build !js

package game

import "os"

var Assets = os.DirFS("game/")
