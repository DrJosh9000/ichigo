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
	"image"
)

// ID implements Identifier directly (as a string value).
type ID string

// Ident returns id as a string.
func (id ID) Ident() string { return string(id) }

// Bounds implements Bounder directly (as an image.Rectangle value).
type Bounds image.Rectangle

// BoundingRect returns b as an image.Rectangle.
func (b Bounds) BoundingRect() image.Rectangle { return image.Rectangle(b) }

// Disables implements Disabler directly (as a bool).
type Disables bool

// Disabled returns d as a bool.
func (d Disables) Disabled() bool { return bool(d) }

// Disable sets d to true.
func (d *Disables) Disable() { *d = true }

// Enable sets d to false.
func (d *Disables) Enable() { *d = false }

// Hides implements Hider directly (as a bool).
type Hides bool

// Hidden returns h as a bool.
func (h Hides) Hidden() bool { return bool(h) }

// Hide sets h to true.
func (h *Hides) Hide() { *h = true }

// Show sets h to false.
func (h *Hides) Show() { *h = false }
