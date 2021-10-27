package util

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getTokenTmpDir() string {
	tokenPath := filepath.Join(os.TempDir(), "suborbital", "compute", "envtoken")
	return tokenPath
}

func WriteEnvironmentToken(tokenStr []byte) error {
	tokenPath := getTokenTmpDir()
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		if _, err := Mkdir(filepath.Dir(tokenPath), ""); err != nil {
            if err != nil {
                return errors.Wrap(err, "failed to write token when create dir")
            }
        }
	}

	if err := ioutil.WriteFile(tokenPath, tokenStr, 0700); err != nil {
		return errors.Wrap(err, "failed to WriteFile for token")
	}
	return nil
}

func ReadEnvironmentToken() ([]byte, error) {
	tokenPath := getTokenTmpDir()
	buf, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read token")
	}
	return buf, nil
}
