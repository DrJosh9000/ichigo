package engine

import (
	"bytes"
	"encoding/gob"
)

var _ interface {
	Registrar
	Scanner
	gob.GobDecoder
	gob.GobEncoder
} = &Container{}

func init() {
	gob.Register(&Container{})
}

// Container contains many components, in order.
type Container struct {
	items   []interface{}
	free    map[int]struct{}
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
	c.free, c.reverse = nil, nil
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

// Prepare ensures the helper data structures are present.
func (c *Container) Prepare(*Game) error {
	if c.reverse == nil {
		c.reverse = make(map[interface{}]int, len(c.items))
		for i, x := range c.items {
			c.reverse[x] = i
		}
	}
	if c.free == nil {
		c.free = make(map[int]struct{})
	}
	return nil
}

// Scan visits every non-nil component in the container.
func (c *Container) Scan(visit func(interface{}) error) error {
	for _, x := range c.items {
		if x != nil {
			if err := visit(x); err != nil {
				return err
			}
		}
	}
	return nil
}

// Element returns the item at index i, or nil for a free slot.
func (c *Container) Element(i int) interface{} { return c.items[i] }

// Len returns the number of items plus the number of free slots in the container.
func (c *Container) Len() int { return len(c.items) }

// Swap swaps any two items, free slots, or a combination.
func (c *Container) Swap(i, j int) {
	if i == j {
		return
	}
	ifree := c.items[i] == nil
	jfree := c.items[j] == nil
	switch {
	case ifree && jfree:
		return
	case ifree:
		c.items[i] = c.items[j]
		c.reverse[c.items[i]] = i
		c.free[j] = struct{}{}
		delete(c.free, i)
	case jfree:
		c.items[j] = c.items[i]
		c.reverse[c.items[j]] = j
		c.free[i] = struct{}{}
		delete(c.free, j)
	default:
		c.items[i], c.items[j] = c.items[j], c.items[i]
		c.reverse[c.items[i]] = i
		c.reverse[c.items[j]] = j
	}
}

func (c Container) String() string { return "Container" }

// Register records component into the slice, if parent is this container. It
// writes the component to an arbitrary free index in the slice, or appends if
// there are none free.
func (c *Container) Register(component, parent interface{}) error {
	if parent != c {
		return nil
	}
	if len(c.free) == 0 {
		c.reverse[component] = len(c.items)
		c.items = append(c.items, component)
		return nil
	}
	for i := range c.free {
		c.reverse[component] = i
		c.items[i] = component
		delete(c.free, i)
		return nil
	}
	return nil
}

// Unregister searches the slice for the component, and removes it by setting
// to nil. If the number of nil items is greater than half the slice, the slice
// is compacted.
func (c *Container) Unregister(component interface{}) {
	i, found := c.reverse[component]
	if !found {
		return
	}
	c.items[i] = nil
	c.free[i] = struct{}{}
	delete(c.reverse, i)
	if len(c.free) > len(c.items)/2 {
		c.compact()
	}
}

// compact moves all the items to the front of the items slice, removing any
// free slots, and empties the free map.
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
	c.free = make(map[int]struct{})
}
