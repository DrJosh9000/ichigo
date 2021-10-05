package engine

import (
	"io/fs"
	"time"
)

// DummyLoad is a loader that just takes up time and doesn't actually load
// anything.
type DummyLoad struct {
	time.Duration
}

// Load sleeps for d.Duration, then returns nil.
func (d DummyLoad) Load(fs.FS) error {
	time.Sleep(d.Duration)
	return nil
}
