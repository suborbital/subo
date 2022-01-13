package util

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func getTokenTmpDir() string {
	tokenPath := filepath.Join(os.TempDir(), "suborbital", "compute", "envtoken")
	return tokenPath
}

func WriteEnvironmentToken(tokenStr string) error {
	tokenPath := getTokenTmpDir()
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(tokenPath), 0755); err != nil {
			return errors.Wrap(err, "failed to Mkdir")
		}
	}

	if err := ioutil.WriteFile(tokenPath, []byte(tokenStr), 0600); err != nil {
		return errors.Wrap(err, "failed to WriteFile for token")
	}
	return nil
}

func ReadEnvironmentToken() (string, error) {
	tokenPath := getTokenTmpDir()
	buf, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to ReadFile for token")
	}
	return string(buf), nil
}
