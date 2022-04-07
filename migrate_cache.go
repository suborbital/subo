package main

import (
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/pkg/errors"

	"github.com/suborbital/subo/subo/util"
)

// migrateCache migrates subo's cached files from the old temp directory to the
// new location if necessary.
func migrateCache() {
	userCachePath, err := os.UserCacheDir()
	if err != nil {
		// migration not possible.
		return
	}

	tmpPath := filepath.Join(os.TempDir(), util.CacheBaseDir)
	if _, err = os.Stat(tmpPath); os.IsNotExist(err) {
		return
	}

	newPath := filepath.Join(userCachePath, util.CacheBaseDir)

	if _, err = os.Stat(newPath); os.IsNotExist(err) {
		if cpErr := cp.Copy(tmpPath, newPath); cpErr != nil {
			util.LogWarn(errors.Wrap(cpErr, "failed to migrate cache directory").Error())
			return
		} else {
			util.LogInfo("successfully migrated cache directory")
		}
	}

	if err = os.RemoveAll(tmpPath); err != nil {
		util.LogWarn(errors.Wrap(err, "failed to remove old cache directory").Error())
	}
}
