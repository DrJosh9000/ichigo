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
	"fmt"
	"image"
	"io/fs"

	"github.com/DrJosh9000/ichigo/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure Tilemap satisfies interfaces.
var _ interface {
	Identifier
	Collider
	Drawer
	Hider
	Scanner
	Transformer
} = &Tilemap{}

// Ensure StaticTile and AnimatedTile satisfy Tile.
var (
	_ Tile = StaticTile(0)
	_ interface {
		Tile
		Scanner
	} = &AnimatedTile{}
)

func init() {
	gob.Register(&AnimatedTile{})
	gob.Register(StaticTile(0))
	gob.Register(&Tilemap{})
}

// Tilemap renders a grid of rectangular tiles at equal Z position.
type Tilemap struct {
	ID
	Disables
	Hides
	Map    map[image.Point]Tile // tilespace coordinate -> tile
	Ersatz bool                 // disables collisions ("fake wall")
	Offset image.Point          // world coordinates
	Sheet  Sheet
}

// CollidesWith implements Collider.
func (t *Tilemap) CollidesWith(b geom.Box) bool {
	if t.Ersatz {
		return false
	}

	// Probe the map at all tilespace coordinates overlapping the rect.
	r := b.XY().Sub(t.Offset) // TODO: pretend tilemap is a plane in 3D?
	min := geom.CDiv(r.Min, t.Sheet.CellSize)
	max := geom.CDiv(r.Max.Sub(image.Pt(1, 1)), t.Sheet.CellSize) // NB: fencepost

	for j := min.Y; j <= max.Y; j++ {
		for i := min.X; i <= max.X; i++ {
			if t.Map[image.Pt(i, j)] != nil {
				return true
			}
		}
	}
	return false
}

// Draw draws the tilemap.
func (t *Tilemap) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	og := opts.GeoM
	for p, tile := range t.Map {
		if tile == nil {
			continue
		}
		var mat ebiten.GeoM
		mat.Translate(geom.CFloat(geom.CMul(p, t.Sheet.CellSize)))
		mat.Concat(og)
		opts.GeoM = mat

		src := t.Sheet.SubImage(tile.Cell())
		screen.DrawImage(src, opts)
	}
}

// Load instantiates animations for all AnimatedTiles.
func (t *Tilemap) Load(fs.FS) error {
	for _, tile := range t.Map {
		at, ok := tile.(*AnimatedTile)
		if !ok {
			continue
		}
		at.anim = t.Sheet.NewAnim(at.AnimKey)
		if at.anim == nil {
			return fmt.Errorf("missing anim %q", at.AnimKey)
		}
	}
	return nil
}

// Scan visits &t.Sheet and all tiles.
func (t *Tilemap) Scan(visit VisitFunc) error {
	if err := visit(&t.Sheet); err != nil {
		return err
	}
	for _, tile := range t.Map {
		if err := visit(tile); err != nil {
			return err
		}
	}
	return nil
}

// Transform returns a translation by t.Offset.
func (t *Tilemap) Transform() (opts ebiten.DrawImageOptions) {
	opts.GeoM.Translate(geom.CFloat(t.Offset))
	return opts
}

// TileAt returns the tile present at the given world coordinate.
func (t *Tilemap) TileAt(wc image.Point) Tile {
	return t.Map[geom.CDiv(wc.Sub(t.Offset), t.Sheet.CellSize)]
}

// SetTileAt sets the tile at the given world coordinate.
func (t *Tilemap) SetTileAt(wc image.Point, tile Tile) {
	t.Map[geom.CDiv(wc.Sub(t.Offset), t.Sheet.CellSize)] = tile
}

// TileBounds returns a rectangle describing the tile boundary for the tile
// at the given world coordinate.
func (t *Tilemap) TileBounds(wc image.Point) image.Rectangle {
	p := geom.CMul(geom.CDiv(wc.Sub(t.Offset), t.Sheet.CellSize), t.Sheet.CellSize).Add(t.Offset)
	return image.Rectangle{p, p.Add(t.Sheet.CellSize)}
}

// Tile is the interface needed by Tilemap.
type Tile interface {
	Cell() int
}

// StaticTile returns a fixed tile index.
type StaticTile int

// Cell returns s as an int.
func (s StaticTile) Cell() int { return int(s) }

// AnimatedTile uses an Anim to choose a tile index.
type AnimatedTile struct {
	AnimKey string

	anim *Anim
}

// Cell returns the value of Cell provided by the animation.
func (a *AnimatedTile) Cell() int { return a.anim.Cell() }

// Scan visits a.anim.
func (a *AnimatedTile) Scan(visit VisitFunc) error {
	return visit(a.anim)
}
