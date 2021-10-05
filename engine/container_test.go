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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeContainer(t *testing.T) {
	c := MakeContainer(69, 420)
	if want := []interface{}{69, 420}; !cmp.Equal(c.items, want) {
		t.Errorf("c.items = %v, want %v", c.items, want)
	}
	if want := map[interface{}]int{69: 0, 420: 1}; !cmp.Equal(c.reverse, want) {
		t.Errorf("c.reverse = %v, want %v", c.reverse, want)
	}
}
