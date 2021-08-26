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
