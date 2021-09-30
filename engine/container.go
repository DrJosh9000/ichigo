package engine

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

var _ interface {
	Prepper
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
	c.Prepare(nil)
	return c
}

// GobDecode decodes a byte slice as though it were a slice of items.
func (c *Container) GobDecode(in []byte) error {
	if err := gob.NewDecoder(bytes.NewReader(in)).Decode(&c.items); err != nil {
		return err
	}
	c.reverse = nil
	return c.Prepare(nil)
}

// GobEncode encodes c as the slice of items.
func (c *Container) GobEncode() ([]byte, error) {
	c.compact()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(c.items); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Prepare ensures the helper data structures are present and valid.
func (c *Container) Prepare(*Game) error {
	c.reverse = make(map[interface{}]int, len(c.items))
	for i, x := range c.items {
		c.reverse[x] = i
	}
	return nil
}

// Scan visits every non-nil component in the container.
func (c *Container) Scan(visit VisitFunc) error {
	if c == nil {
		return nil
	}
	for _, x := range c.items {
		if x != nil {
			if err := visit(x); err != nil {
				return err
			}
		}
	}
	return nil
}

// Add adds an item to the end of the container, if not already present.
func (c *Container) Add(component interface{}) {
	if c.Contains(component) {
		return
	}
	c.reverse[component] = len(c.items)
	c.items = append(c.items, component)
}

// Remove replaces an item with nil. If the number of nil items is greater than
// half the slice, the slice is compacted (indexes of items will change).
func (c *Container) Remove(component interface{}) {
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
func (c *Container) Contains(component interface{}) bool {
	if c == nil {
		return false
	}
	_, found := c.reverse[component]
	return found
}

// IndexOf reports if an item exists in the container and returns the index if
// present.
func (c *Container) IndexOf(component interface{}) (int, bool) {
	if c == nil {
		return 0, false
	}
	i, found := c.reverse[component]
	return i, found
}

func (c *Container) ItemCount() int {
	if c == nil {
		return 0
	}
	return len(c.reverse)
}

// Element returns the item at index i, or nil for a free slot.
func (c *Container) Element(i int) interface{} { return c.items[i] }

// Len returns the number of items plus the number of free slots in the container.
func (c *Container) Len() int {
	if c == nil {
		return 0
	}
	return len(c.items)
}

// Swap swaps any two items, free slots, or a combination.
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
