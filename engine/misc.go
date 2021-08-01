package engine

import (
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

// ID implements Identifier directly (as a string value).
type ID string

// Ident returns id as a string.
func (id ID) Ident() string { return string(id) }

// ZPos implements ZPositioner directly (as a float64 value).
type ZPos float64

// Z returns z as a float64.
func (z ZPos) Z() float64 { return float64(z) }

// GeoMDef is a serialisable form of GeoM.
type GeoMDef [6]float64 // Assumption: this has identical memory layout to GeoM

// ToGeoMDef translates a GeoM to a GeoMDef using unsafe.Pointer.
func ToGeoMDef(m *ebiten.GeoM) *GeoMDef {
	return (*GeoMDef)(unsafe.Pointer(m))
}

// GeoM translates a GeoMDef to a GeoM using unsafe.Pointer.
func (d *GeoMDef) GeoM() *ebiten.GeoM {
	return (*ebiten.GeoM)(unsafe.Pointer(d))
}
