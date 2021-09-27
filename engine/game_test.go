package engine

import "testing"

func TestGameLoadAndPrepare(t *testing.T) {
	g := &Game{
		Root: fakeDrawBoxer("fake"),
	}
	if err := g.LoadAndPrepare(nil); err != nil {
		t.Errorf("LoadAndPrepare(nil) = %v, want nil", err)
	}
}
