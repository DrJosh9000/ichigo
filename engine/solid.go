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

	"github.com/DrJosh9000/ichigo/geom"
)

var _ Collider = SolidRect{}

func init() {
	gob.Register(&SolidRect{})
}

// SolidRect is a minimal implementation of a Collider defined by a single Box.
type SolidRect struct {
	ID
	geom.Box
}

// CollidesWith reports if r overlaps with s.Box.
func (s SolidRect) CollidesWith(r geom.Box) bool {
	return s.Box.Overlaps(r)
}
