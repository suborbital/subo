package util

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// TemplateDir gets the template directory for subo and ensures it exists
func TemplateDir() (string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get UserConfigDir")
	}

	tmplPath := filepath.Join(config, "suborbital", "templates")

	if os.Stat(tmplPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(tmplPath, os.ModePerm); err != nil {
				return "", errors.Wrap(err, "failed to MkdirAll template directory")
			}
		} else {
			return "", errors.Wrap(err, "failed to Stat template directory")
		}
	}

	return tmplPath, nil
}
