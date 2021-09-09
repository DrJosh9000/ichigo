package engine

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const screenShaderSrc = `package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	// Image 0 is the image being drawn.
	// Image 1 is the depth map of the image being drawn.
	// Image 2 is the depth buffer.

	// Convert the colours from the depth map and depth buffer into scalars.
	stairs := vec4(1, 256, 65536, 16777216)
	c1 := imageSrc1UnsafeAt(texCoord)
	c2 := imageSrc2UnsafeAt(texCoord)
	d1 := dot(c1, stairs)
	d2 := dot(c2, stairs)

	// If the depth buffer value is higher, return transparent.
	if d2 > d1 {
		return vec4(0)
	}
	return imageSrc0UnsafeAt(texCoord)
}
`

const depthShaderSrc = `package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	// Image 0 is the depth buffer being updated.
	// Image 1 is the depth map being drawn from.

	// Convert the colour vectors into scalars.
	stairs := vec4(1, 256, 65536, 16777216)
	c0 := imageSrc0UnsafeAt(texCoord)
	c1 := imageSrc1UnsafeAt(texCoord)
	d0 := dot(c0, stairs)
	d1 := dot(c1, stairs)

	// Write back the larger of the two.
	if d0 > d1 {
		return c0
	}
	return c1
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
	buffer *ebiten.Image
}

func NewDepthBuffer(size image.Point) *DepthBuffer {
	return &DepthBuffer{
		buffer: ebiten.NewImage(size.X, size.Y),
	}
}

func (b *DepthBuffer) Draw(dst, src, srcDepth *ebiten.Image) {
	w, h := 0, 0
	dst.DrawRectShader(w, h, screenShader, &ebiten.DrawRectShaderOptions{
		Images: [4]*ebiten.Image{src, srcDepth, b.buffer},
	})
	dst.DrawRectShader(w, h, depthShader, &ebiten.DrawRectShaderOptions{
		Images: [4]*ebiten.Image{b.buffer, srcDepth},
	})
}
