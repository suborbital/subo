package util

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CacheDir returns the cache directory and creates it if it doesn't exist
func CacheDir() (string, error) {
	targetPath := filepath.Join(os.TempDir(), "suborbital", "cache")

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			return "", errors.Wrap(err, "failed to MkdirAll")
		}
	}
	return targetPath, nil
}
