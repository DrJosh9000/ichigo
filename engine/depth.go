package engine

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const screenShaderSrc = `package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return imageSrc0UnsafeAt(texCoord)
}
`

const depthShaderSrc = `package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return imageSrc0UnsafeAt(texCoord)
}
`

var depthShader, screenShader *ebiten.Shader

func init() {
	ds, err := ebiten.NewShader([]byte(depthShaderSrc))
	if err != nil {
		log.Fatalf("Compiling depth shader: %v", err)
	}
	depthShader = ds

	ss, err := ebiten.NewShader([]byte(screenShaderSrc))
	if err != nil {
		log.Fatalf("Compiling screen shader: %v", err)
	}
	screenShader = ss
}

type DepthBuffer struct {
	Buffer *ebiten.Image
}
