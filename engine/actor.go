package engine

import (
	"encoding/gob"

	"drjosh.dev/gurgle/geom"
)

// Ensure Actor satisfies interfaces.
var _ Prepper = &Actor{}

func init() {
	gob.Register(&Actor{})
}

// Thorson-style movement:
// https://maddythorson.medium.com/celeste-and-towerfall-physics-d24bd2ae0fc5

// Actor handles basic movement.
type Actor struct {
	CollisionDomain string    // id of component to look for colliders inside of
	Pos, Size       geom.Int3 // in voxels; multiply by game.VoxelScale for regular Euclidean space

	rem  geom.Float3
	game *Game
}

// CollidesAt runs a collision test of the actor, supposing the actor is at a
// given position (not necessarily a.Pos).
func (a *Actor) CollidesAt(p geom.Int3) bool {
	bounds := geom.Box{Min: p, Max: p.Add(a.Size)}
	for c := range a.game.Query(a.CollisionDomain, ColliderType) {
		if c.(Collider).CollidesWith(bounds) {
			return true
		}
	}
	return false
}

// MoveX moves the actor x units in world space. It takes Game.VoxelScale into
// account (so MoveX(x) moves the actor x/VoxelScale.X voxel units). onCollide
// is called if a collision occurs, and the actor wil be in the colliding
// position during the call.
func (a *Actor) MoveX(x float64, onCollide func()) {
	a.rem.X += x / a.game.VoxelScale.X
	move := int(a.rem.X + 0.5) // Note: math.Round can lead to vibration
	if move == 0 {
		return
	}
	a.rem.X -= float64(move)
	sign := geom.Sign(move)
	for move != 0 {
		a.Pos.X += sign
		move -= sign
		if !a.CollidesAt(a.Pos) {
			continue
		}
		if onCollide != nil {
			onCollide()
		}
		a.Pos.X -= sign
		a.rem.X = 0
		return
	}
}

// MoveY is like MoveX but in the Y dimension. See MoveX for more information.
func (a *Actor) MoveY(y float64, onCollide func()) {
	a.rem.Y += y / a.game.VoxelScale.Y
	move := int(a.rem.Y + 0.5)
	if move == 0 {
		return
	}
	a.rem.Y -= float64(move)
	sign := geom.Sign(move)
	for move != 0 {
		a.Pos.Y += sign
		move -= sign
		if !a.CollidesAt(a.Pos) {
			continue
		}
		if onCollide != nil {
			onCollide()
		}
		a.Pos.Y -= sign
		a.rem.Y = 0
		return
	}
}

// MoveZ is like MoveX but in the Y dimension. See MoveX for more information.
func (a *Actor) MoveZ(z float64, onCollide func()) {
	a.rem.Z += z / a.game.VoxelScale.Z
	move := int(a.rem.Z + 0.5)
	if move == 0 {
		return
	}
	a.rem.Z -= float64(move)
	sign := geom.Sign(move)
	for move != 0 {
		a.Pos.Z += sign
		move -= sign
		if !a.CollidesAt(a.Pos) {
			continue
		}
		if onCollide != nil {
			onCollide()
		}
		a.Pos.Z -= sign
		a.rem.Z = 0
		return
	}
}

// Prepare stores a reference to the game.
func (a *Actor) Prepare(g *Game) error {
	a.game = g
	return nil
}
