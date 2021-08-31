package engine

// Box describes an axis-aligned rectangular prism.
type Box struct {
	X, Y, Z int // coordinate of the left-top-farthest corner
	W, H, D int // width, height, depth
}

// IsoProjection translates an integer 3D coordinate into an integer 2D
// coordinate.
type IsoProjection struct {
	ZX, ZY int
}

// Project projects a 3D coordinate into 2D.
// If ZX = 0, x is unchanged; similarly for ZY and y.
// Otherwise, x becomes x + z/ZX and y becomes y + z/ZY.
// This is because there's little reason for an isometric projection in a game
// to exaggerate the Z position.
func (π IsoProjection) Project(x, y, z int) (xp, yp int) {
	xp, yp = x, y
	if π.ZX != 0 {
		xp += z / π.ZX
	}
	if π.ZY != 0 {
		yp += z / π.ZY
	}
	return xp, yp
}
