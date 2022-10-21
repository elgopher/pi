// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	extension = ".json"
	dir       = "saves"
)

func Load(name string) (string, error) {
	bytes, err := os.ReadFile(filename(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("error reading file %s: %w", filename(name), err)
	}

	return string(bytes), nil
}

func filename(name string) string {
	return fmt.Sprintf("%s%c%s%s", dir, os.PathSeparator, name, extension)
}

func Save(name string, data string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("problem creating %s directory: %w", dir, err)
	}

	tmpName := tmpFilename(name)
	if err := os.WriteFile(tmpName, []byte(data), 0644); err != nil {
		return fmt.Errorf("error writing tmp file %s: %w", tmpName, err)
	}

	if err := os.Rename(tmpName, filename(name)); err != nil {
		return fmt.Errorf("error renaming tmp file %s to %s: %w", tmpName, filename(name), err)
	}

	return nil
}

func tmpFilename(name string) string {
	return filename(name) + ".tmp"
}

func Delete(name string) error {
	if err := os.Remove(filename(name)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("error deleting file %s: %w", filename(name), err)
	}
	return nil
}

func Names() ([]string, error) {
	var names []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(entryName, extension) {
			name := strings.TrimSuffix(entryName, extension)
			names = append(names, name)
		}
	}

	return names, nil
}

func Cleanup() {
	_ = os.RemoveAll(dir)
}
