package engine

import (
	"encoding/gob"
	"image"
	"strconv"

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

// Point3 is a an element of int^3.
type Point3 struct {
	X, Y, Z int
}

// Pt3(x, y, z) is shorthand for Point3{x, y, z}.
func Pt3(x, y, z int) Point3 {
	return Point3{x, y, z}
}

// String returns a string representation of p like "(3,4,5)".
func (p Point3) String() string {
	return "(" + strconv.Itoa(p.X) + "," + strconv.Itoa(p.Y) + "," + strconv.Itoa(p.Z) + ")"
}

// XY applies the Z-forgetting projection. (It returns just X and Y.)
func (p Point3) XY() image.Point {
	return image.Point{p.X, p.Y}
}

// Add performs vector addition.
func (p Point3) Add(q Point3) Point3 {
	return Point3{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

// Sub performs vector subtraction.
func (p Point3) Sub(q Point3) Point3 {
	return p.Add(q.Neg())
}

// CMul performs componentwise multiplication.
func (p Point3) CMul(q Point3) Point3 {
	return Point3{p.X * q.X, p.Y * q.Y, p.Z * q.Z}
}

// Mul performs scalar multiplication.
func (p Point3) Mul(k int) Point3 {
	return Point3{p.X * k, p.Y * k, p.Z * k}
}

// CDiv performs componentwise division.
func (p Point3) CDiv(q Point3) Point3 {
	return Point3{p.X / q.X, p.Y / q.Y, p.Z / q.Z}
}

// Div performs scalar division by k.
func (p Point3) Div(k int) Point3 {
	return Point3{p.X / k, p.Y / k, p.Z / k}
}

// Neg returns the vector pointing in the opposite direction.
func (p Point3) Neg() Point3 {
	return Point3{-p.X, -p.Y, -p.Z}
}

// Coord returns the components of the vector.
func (p Point3) Coord() (x, y, z int) {
	return p.X, p.Y, p.Z
}

// IsoProject performs isometric projection of a 3D coordinate into 2D.
//
// If π.X = 0, the x returned is p.X; similarly for π.Y and y.
// Otherwise, x projects to x + z/π.X and y projects to y + z/π.Y.
func (p Point3) IsoProject(π image.Point) image.Point {
	/*
		I'm using the π character because I'm a maths wanker.

		Dividing is used because there's little reason for an isometric
		projection in a game to exaggerate the Z position.

		Integers are used to preserve that "pixel perfect" calculation in case
		you are making the next Celeste.
	*/
	q := image.Point{p.X, p.Y}
	if π.X != 0 {
		q.X += p.Z / π.X
	}
	if π.Y != 0 {
		q.Y += p.Z / π.Y
	}
	return q
}

// Box describes an axis-aligned rectangular prism.
type Box struct {
	Min, Max Point3
}

// String returns a string representation of b like "(3,4,5)-(6,5,8)".
func (b Box) String() string {
	return b.Min.String() + "-" + b.Max.String()
}

// Empty reports whether the box contains no points.
func (b Box) Empty() bool {
	return b.Min.X >= b.Max.X || b.Min.Y >= b.Max.Y || b.Min.Z >= b.Max.Z
}

// Eq reports whether b and c contain the same set of points. All empty boxes
// are considered equal.
func (b Box) Eq(c Box) bool {
	return b == c || b.Empty() && c.Empty()
}

// Overlaps reports whether b and c have non-empty intersection.
func (b Box) Overlaps(c Box) bool {
	return !b.Empty() && !c.Empty() &&
		b.Min.X < c.Max.X && c.Min.X < b.Max.X &&
		b.Min.Y < c.Max.Y && c.Min.Y < b.Max.Y &&
		b.Min.Z < c.Max.Z && c.Min.Z < b.Max.Z
}

// Size returns b's width, height, and depth.
func (b Box) Size() Point3 {
	return b.Max.Sub(b.Min)
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
	return z, dot(v.vox.pos.XY(), v.vox.ivm.DrawOrderBias)
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
