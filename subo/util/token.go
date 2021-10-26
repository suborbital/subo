package util

import (
    "github.com/pkg/errors"
    "io/ioutil"
    "os"
    "path/filepath"
)

const tokenDir = "suborbital/compute/envtoken"

type TokenData struct {}

func (_ TokenData) GetTokenTmpDir() string {
	tokenPath := filepath.Join(os.TempDir(), tokenDir)
	return tokenPath
}

func (token TokenData) WriteToken(tokenStr []byte) error {
	tokenPath := token.GetTokenTmpDir()
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		_, err := Mkdir(filepath.Dir(tokenPath), "")
		if err != nil {
			return errors.Wrap(err, "failed to write token when create dir")
		}
	}

	if err := ioutil.WriteFile(tokenPath, tokenStr, 0700); err != nil {
		return errors.Wrap(err, "failed to write token")
	}
	return nil
}

func (token TokenData) ReadToken() ([]byte, error) {
	tokenPath := token.GetTokenTmpDir()
	buf, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read token")
	}
	return buf, nil
}
