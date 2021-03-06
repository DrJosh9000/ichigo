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
	"log"
	"time"
)

// LoadingSwitch switches between two subcomponents. While After is being
// loaded asynchronously, During is shown. Once loading is complete, During
// is hidden and After is shown.
type LoadingSwitch struct {
	During, After interface {
		Disabler
		Hider
	}

	assets fs.FS
}

// Scan only scans s.During. Only s.During is loaded automatically - s.After is
// loaded asynchronously from Prepare below.
func (s *LoadingSwitch) Scan(visit VisitFunc) error {
	return visit(s.During)
}

// Load stores a copy of assets to pass to s.After.Load later.
func (s *LoadingSwitch) Load(assets fs.FS) error {
	s.assets = assets
	return nil
}

// Prepare loads, registers, and prepares.After in a separate goroutine. Once
// ready, LoadingSwitch hides s.During and shows s.After.
func (s *LoadingSwitch) Prepare(game *Game) error {
	go s.loadAfter(game)
	return nil
}

func (s *LoadingSwitch) loadAfter(game *Game) {
	startLoad := time.Now()
	if err := game.Load(s.After, s.assets); err != nil {
		log.Printf("Couldn't load: %v", err)
		return
	}
	log.Printf("LoadingSwitch: finished loading in %v", time.Since(startLoad))

	s.After.Disable()
	s.After.Hide()

	startBuild := time.Now()
	if err := game.Register(s.After, s); err != nil {
		log.Printf("Couldn't register: %v", err)
		return
	}
	log.Printf("LoadingSwitch: finished registering in %v", time.Since(startBuild))
	startPrep := time.Now()
	if err := game.Prepare(s.After); err != nil {
		log.Printf("Couldn't prepare: %v", err)
		return
	}
	log.Printf("LoadingSwitch: finished preparing in %v", time.Since(startPrep))

	// TODO: better scene transitions
	s.During.Disable()
	s.During.Hide()
	s.After.Enable()
	s.After.Show()
}
