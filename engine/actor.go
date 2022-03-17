/*
Copyright 2021 Josh Deprez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package engine

import (
	"encoding/gob"
	"errors"
	"log"

	"github.com/DrJosh9000/ichigo/geom"
)

// Ensure Actor satisfies interfaces.
var _ interface {
	BoundingBoxer
	Prepper
} = &Actor{}

var errCollision = errors.New("collision detected")

func init() {
	gob.Register(&Actor{})
}

// Thorson-style movement:
// https://maddythorson.medium.com/celeste-and-towerfall-physics-d24bd2ae0fc5

// Actor handles basic movement.
type Actor struct {
	CollisionDomain string    // id of component to look for colliders inside of
	Pos             geom.Int3 // in voxels; multiply by game.VoxelScale for regular Euclidean space
	Bounds          geom.Box  // in voxels; relative to Pos

	rem  geom.Float3
	game *Game
}

// BoundingBox returns the box Bounds.Add(Pos).
func (a *Actor) BoundingBox() geom.Box {
	return a.Bounds.Add(a.Pos)
}

// CollidesAt runs a collision test of the actor, supposing the actor is at a
// given position (not necessarily a.Pos).
func (a *Actor) CollidesAt(p geom.Int3) bool {
	bounds := a.Bounds.Add(p)
	cd := a.game.Component(a.CollisionDomain)
	if cd == nil {
		log.Printf("collision domain %q not found", a.CollisionDomain)
		return false
	}
	return errCollision == a.game.Query(cd, ColliderType, nil, func(c any) error {
		if cl, ok := c.(Collider); ok && cl.CollidesWith(bounds) {
			return errCollision
		}
		return nil
	})
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

func (a *Actor) String() string { return "Actor@" + a.Pos.String() }
