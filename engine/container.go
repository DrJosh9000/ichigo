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
	"bytes"
	"encoding/gob"
	"fmt"
)

var _ interface {
	Scanner
	gob.GobDecoder
	gob.GobEncoder
} = &Container{}

func init() {
	gob.Register(&Container{})
}

// Container is a component that contains many other components, in order.
// It can be used as both a component in its own right, or as a ordered set.
// A nil *Container contains no items and modifications will panic (like a map).
type Container struct {
	items   []interface{}
	reverse map[interface{}]int
}

// MakeContainer puts the items into a new Container.
func MakeContainer(items ...interface{}) *Container {
	c := &Container{items: items}
	c.rebuildReverse()
	return c
}

func (c *Container) rebuildReverse() {
	if c == nil {
		return
	}
	c.reverse = make(map[interface{}]int, len(c.items))
	for i, x := range c.items {
		c.reverse[x] = i
	}
}

// GobDecode decodes a byte slice as though it were a slice of items.
func (c *Container) GobDecode(in []byte) error {
	if err := gob.NewDecoder(bytes.NewReader(in)).Decode(&c.items); err != nil {
		return err
	}
	c.rebuildReverse()
	return nil
}

// GobEncode encodes c as the slice of items.
// When called on a nil *Container, GobEncode returns a nil slice.
func (c *Container) GobEncode() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	c.compact()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(c.items); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Scan visits every non-nil component in the container.
// Scan is safe to call on a nil *Container.
func (c *Container) Scan(visit VisitFunc) error {
	if c == nil {
		return nil
	}
	for _, x := range c.items {
		if x == nil {
			continue
		}
		if err := visit(x); err != nil {
			return err
		}
	}
	return nil
}

// Add adds an item to the end of the container, if not already present.
// Add is _not_ safe to call on a nil *Container.
func (c *Container) Add(component interface{}) {
	if c.Contains(component) {
		return
	}
	c.reverse[component] = len(c.items)
	c.items = append(c.items, component)
}

// Remove replaces an item with nil. If the number of nil items is greater than
// half the slice, the slice is compacted (indexes of items will change).
// Remove is safe to call on a nil *Container.
func (c *Container) Remove(component interface{}) {
	if c == nil {
		return
	}
	i, found := c.reverse[component]
	if !found {
		return
	}
	c.items[i] = nil
	delete(c.reverse, i)
	if len(c.reverse) < len(c.items)/2 {
		c.compact()
	}
}

// Contains reports if an item exists in the container.
// Contains is safe to call on a nil *Container.
func (c *Container) Contains(component interface{}) bool {
	if c == nil {
		return false
	}
	_, found := c.reverse[component]
	return found
}

// IndexOf reports if an item exists in the container and returns the index if
// present.
// IndexOf is safe to call on a nil *Container.
func (c *Container) IndexOf(component interface{}) (int, bool) {
	if c == nil {
		return 0, false
	}
	i, found := c.reverse[component]
	return i, found
}

// ItemCount returns the number of (non-nil) items in the container.
// ItemCount is safe to call on a nil *Container.
func (c *Container) ItemCount() int {
	if c == nil {
		return 0
	}
	return len(c.reverse)
}

// Element returns the item at index i, or nil for a free slot.
// Element is _not_ safe to call on a nil *Container.
func (c *Container) Element(i int) interface{} { return c.items[i] }

// Len returns the number of items plus the number of nil slots in the container.
// Len is safe to call on a nil *Container.
func (c *Container) Len() int {
	if c == nil {
		return 0
	}
	return len(c.items)
}

// Swap swaps any two items, free slots, or a combination.
// Swap is _not_ safe to call on a nil *Container.
func (c *Container) Swap(i, j int) {
	c.items[i], c.items[j] = c.items[j], c.items[i]
	if c.items[i] != nil {
		c.reverse[c.items[i]] = i
	}
	if c.items[j] != nil {
		c.reverse[c.items[j]] = j
	}
}

func (c *Container) String() string {
	if c == nil {
		return "Container(nil)"
	}
	return "Container" + fmt.Sprint(c.items)
}

// compact moves all the items to the front of the items slice, removing any
// free slots, and resets the free counter.
func (c *Container) compact() {
	i := 0
	for _, x := range c.items {
		if x != nil {
			c.items[i] = x
			c.reverse[x] = i
			i++
		}
	}
	c.items = c.items[:i]
}
