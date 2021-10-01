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
	"compress/gzip"
	"encoding/gob"
	"io/fs"
	"os"
	"path/filepath"
)

type assetKey struct {
	assets fs.FS
	path   string
}

// LoadGobz gunzips and gob-decodes a component from a file from a FS.
func LoadGobz(dst interface{}, assets fs.FS, path string) error {
	f, err := assets.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	return gob.NewDecoder(gz).Decode(dst)
}

// SaveGobz takes an object, gob-encodes it, gzips it, and writes to disk.
// This requires running on something with a disk to write to (not JS)
func SaveGobz(src interface{}, name string) error {
	f, err := os.CreateTemp(".", filepath.Base(name))
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	gz := gzip.NewWriter(f)
	if err := gob.NewEncoder(gz).Encode(src); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), name)
}
