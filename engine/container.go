package engine

// Container contains many components.
type Container []interface{}

// Scan returns c.
func (c Container) Scan() []interface{} { return c }

// Register records component in the slice, if parent is the container.
func (c *Container) Register(component, parent interface{}) error {
	if parent == c {
		*c = append(*c, component)
	}
	return nil
}

// Unregister searches the slice for the component, and removes it by setting
// to nil. If the number of nil items is greater than half the slice, the slice
// is compacted.
func (c *Container) Unregister(component interface{}) {
	free := 0
	for i, x := range *c {
		switch x {
		case component:
			(*c)[i] = nil
			free++
		case nil:
			free++
		}
	}
	if free > len(*c)/2 {
		c.compact()
	}
}

func (c *Container) compact() {
	i := 0
	for _, x := range *c {
		if x != nil {
			(*c)[i] = x
			i++
		}
	}
	*c = (*c)[:i]
}
