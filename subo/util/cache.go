package util

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetCacheDir returns and creates the cache directory
func GetCacheDir() (string, error) {
	cachePath, err := os.UserCacheDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get UserCacheDir")
	}

	targetPath := filepath.Join(cachePath, "suborbital")

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			return "", errors.Wrap(err, "failed to MkdirAll")
		}
	}
	return targetPath, nil
}
