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
