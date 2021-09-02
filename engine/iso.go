package engine

import (
	"encoding/gob"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	_ interface {
		Prepper
		Scanner
	} = &IsoVoxmap{}

	_ interface {
		Prepper
		Scanner
		Transformer
	} = &IsoVoxel{}

	_ interface {
		Drawer
		Transformer
	} = &IsoVoxelSide{}
)

func init() {
	gob.Register(&IsoVoxmap{})
	gob.Register(&IsoVoxel{})
	gob.Register(&IsoVoxelSide{})
}

// IsoVoxmap implements a voxel map, painted using flat images in 2D.
type IsoVoxmap struct {
	ID
	Disabled
	Hidden
	Map           map[Point3]*IsoVoxel
	DrawOrderBias image.Point // so boxes overdraw correctly
	OffsetBack    image.Point // offsets the image drawn for the back
	OffsetFront   image.Point // offsets the image drawn for the front
	Projection    image.Point // IsoProjection parameter
	Sheet         Sheet
	VoxSize       Point3 // size of each voxel
}

// Prepare makes sure all voxels know about the map and where they are, for
// drawing.
func (m *IsoVoxmap) Prepare(*Game) error {
	// Ensure all child units know about wall, which houses common attributes
	for p, u := range m.Map {
		u.pos, u.ivm = p, m
	}
	return nil
}

// Scan returns the Sheet and all voxels in the map.
func (m *IsoVoxmap) Scan() []interface{} {
	c := make([]interface{}, 1, len(m.Map)+1)
	c[0] = &m.Sheet
	for _, voxel := range m.Map {
		c = append(c, voxel)
	}
	return c
}

// IsoVoxel is a voxel in an IsoVoxmap.
type IsoVoxel struct {
	CellBack  int // cell to draw for back side
	CellFront int // cell to draw for front side

	back  IsoVoxelSide
	front IsoVoxelSide
	ivm   *IsoVoxmap
	pos   Point3
}

// Prepare tells the front and back about the voxel.
func (v *IsoVoxel) Prepare(*Game) error {
	v.back.vox = v
	v.front.vox = v
	v.front.front = true
	return nil
}

// Scan returns the back and front of the voxel.
func (v *IsoVoxel) Scan() []interface{} {
	return []interface{}{&v.back, &v.front}
}

// Transform returns a translation of pos.CMul(VoxSize) iso-projected
// (the top-left of the back of the voxel).
func (v *IsoVoxel) Transform() (opts ebiten.DrawImageOptions) {
	p3 := v.pos.CMul(v.ivm.VoxSize)
	p2 := p3.IsoProject(v.ivm.Projection)
	opts.GeoM.Translate(cfloat(p2))
	return opts
}

// IsoVoxelSide is a side of a voxel.
type IsoVoxelSide struct {
	front bool
	vox   *IsoVoxel
}

// Draw draws this side.
func (v *IsoVoxelSide) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	cell := v.vox.CellBack
	if v.front {
		cell = v.vox.CellFront
	}
	screen.DrawImage(v.vox.ivm.Sheet.SubImage(cell), opts)
}

// DrawOrder returns the Z of the nearest or farthest vertex of the voxel,
// with a bias equal to the dot product of the bias vector with pos.XY().
func (v *IsoVoxelSide) DrawOrder() (int, int) {
	z := v.vox.pos.Z * v.vox.ivm.VoxSize.Z
	if v.front {
		z += v.vox.ivm.VoxSize.Z - 1
	}
	bias := dot(v.vox.pos.XY(), v.vox.ivm.DrawOrderBias)
	return z, bias
}

// Transform offsets the image by either OffsetBack or OffsetFront.
func (v *IsoVoxelSide) Transform() (opts ebiten.DrawImageOptions) {
	if v.front {
		opts.GeoM.Translate(cfloat(v.vox.ivm.OffsetFront))
	} else {
		opts.GeoM.Translate(cfloat(v.vox.ivm.OffsetBack))
	}
	return opts
}
