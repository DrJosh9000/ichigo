package engine

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestContainerLiteral(t *testing.T) {
	c := &Container{}
	if err := c.Prepare(nil); err != nil {
		t.Errorf("c.Prepare() = %v, want nil", err)
	}
}

func TestMakeContainer(t *testing.T) {
	c := MakeContainer(69, 420)
	if want := []interface{}{69, 420}; !cmp.Equal(c.items, want) {
		t.Errorf("c.items = %v, want %v", c.items, want)
	}
	if want := make(map[int]struct{}); !cmp.Equal(c.free, want) {
		t.Errorf("c.free = %v, want %v", c.free, want)
	}
	if want := map[interface{}]int{69: 0, 420: 1}; !cmp.Equal(c.reverse, want) {
		t.Errorf("c.reverse = %v, want %v", c.reverse, want)
	}
}
