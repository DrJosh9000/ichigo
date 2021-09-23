package engine

var _ interface {
	Registrar
	Scanner
} = &Container{}

// Container contains many components, in order.
type Container struct {
	Items []interface{}

	free    map[int]struct{}
	reverse map[interface{}]int
}

func MakeContainer(items ...interface{}) *Container {
	c := &Container{Items: items}
	c.Prepare(nil)
	return c
}

func (c *Container) Prepare(*Game) error {
	if c.reverse == nil {
		c.reverse = make(map[interface{}]int, len(c.Items))
		for i, x := range c.Items {
			c.reverse[x] = i
		}
	}
	if c.free == nil {
		c.free = make(map[int]struct{})
	}
	return nil
}

// Scan visits every component in the container.
func (c *Container) Scan(visit func(interface{}) error) error {
	for _, x := range c.Items {
		if err := visit(x); err != nil {
			return err
		}
	}
	return nil
}

// Len returns the number of items in the container.
func (c *Container) Len() int { return len(c.Items) }

// Swap swaps two items in the container.
func (c *Container) Swap(i, j int) {
	if i == j {
		return
	}
	ifree := c.Items[i] == nil
	jfree := c.Items[j] == nil
	switch {
	case ifree && jfree:
		return
	case ifree:
		c.Items[i] = c.Items[j]
		c.reverse[c.Items[i]] = i
		c.free[j] = struct{}{}
		delete(c.free, i)
	case jfree:
		c.Items[j] = c.Items[i]
		c.reverse[c.Items[j]] = j
		c.free[i] = struct{}{}
		delete(c.free, j)
	default:
		c.Items[i], c.Items[j] = c.Items[j], c.Items[i]
		c.reverse[c.Items[i]] = i
		c.reverse[c.Items[j]] = j
	}
}

func (c Container) String() string { return "Container" }

// Register records component in the slice, if parent is this container. It
// writes the component to an arbitrary free index in the slice, or appends if
// there are none free.
func (c *Container) Register(component, parent interface{}) error {
	if parent != c {
		return nil
	}
	if len(c.free) == 0 {
		c.reverse[component] = len(c.Items)
		c.Items = append(c.Items, component)
		return nil
	}
	for i := range c.free {
		c.reverse[component] = i
		c.Items[i] = component
		delete(c.free, i)
		return nil
	}
	return nil
}

// Unregister searches the slice for the component, and removes it by setting
// to nil. If the number of nil items is greater than half the slice, the slice
// is compacted.
func (c *Container) Unregister(component interface{}) {
	if i, found := c.reverse[component]; found {
		c.Items[i] = nil
		c.free[i] = struct{}{}
		delete(c.reverse, i)
	}
	if len(c.free) > len(c.Items)/2 {
		c.compact()
	}
}

func (c *Container) compact() {
	i := 0
	for _, x := range c.Items {
		if x != nil {
			c.Items[i] = x
			c.reverse[x] = i
			i++
		}
	}
	c.Items = c.Items[:i]
	c.free = make(map[int]struct{})
}
