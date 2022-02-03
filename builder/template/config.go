package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/suborbital/subo/subo/util"
)

func TemplateFullPath(repo, branch string) (string, error) {
	repoParts := strings.Split(repo, "/")
	if len(repoParts) != 2 {
		return "", fmt.Errorf("repo is invalid, contains %d parts", len(repoParts))
	}

	repoName := repoParts[1]

	root, err := TemplateRootDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to TemplateRootDir")
	}

	return filepath.Join(root, fmt.Sprintf("%s-%s", repoName, strings.ReplaceAll(branch, "/", "-")), "templates"), nil
}

// TemplateRootDir gets the template directory for subo and ensures it exists
func TemplateRootDir() (string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get UserConfigDir")
	}

	tmplPath := filepath.Join(config, "suborbital", "templates")

	if os.Stat(tmplPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(tmplPath, util.PermDirectory); err != nil {
				return "", errors.Wrap(err, "failed to MkdirAll template directory")
			}
		} else {
			return "", errors.Wrap(err, "failed to Stat template directory")
		}
	}

	return tmplPath, nil
}
