package release

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/suborbital/subo/subo/util"
)

func getTimestampCache() (time.Time, error) {
	cachePath, err := util.CacheDir()
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to CacheDir")
	}

	cachedTimestamp := time.Time{}
	filePath := filepath.Join(cachePath, "subo_last_checked.txt")
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
	} else if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to Stat")
	} else {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return time.Time{}, errors.Wrap(err, "failed to ReadFile")
		}

		cachedTimestamp, err = time.Parse(time.RFC3339, string(data))
		if err != nil {
			return time.Time{}, errors.Wrap(err, "failed to parse cached timestamp")
		}
	}
	return cachedTimestamp, nil
}


	// check if 1 hour has passed since the last version check, and update the cached timestamp if so
	currentTimestamp := time.Now().UTC()
	if cachedTimestamp.IsZero() || currentTimestamp.After(cachedTimestamp.Add(time.Hour)) {
		data := []byte(currentTimestamp.Format(time.RFC3339))
		if err := ioutil.WriteFile(filePath, data, os.ModePerm); err != nil {
			return false, errors.Wrap(err, "failed to WriteFile")
		}
	} else {
		// if 1 hour has not passed, skip version check
		return false, nil
	}

	return true, nil
}

// CheckForLatestVersion returns an error if SuboDotVersion does not match the latest GitHub release or if the check fails
func CheckForLatestVersion() (string, error) {
	getLatestVersion, err := getTimestampCache()
	if err != nil {
		return "", errors.Wrap(err, "failed to getTimestampCache")
	} else if !getLatestVersion {
		return "", nil
	}

	latestRepoRelease, _, err := github.NewClient(nil).Repositories.GetLatestRelease(context.Background(), "suborbital", "subo")
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch latest subo release")
	}

	latestCmdVersion, err := version.NewVersion(*latestRepoRelease.TagName)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse latest subo version")
	}

	cmdVersion, err := version.NewVersion(SuboDotVersion)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse current subo version")
	} else if cmdVersion.LessThan(latestCmdVersion) {
		return fmt.Sprintf("An upgrade for subo is available: %s → %s\n", cmdVersion, latestCmdVersion), nil
	}

	return "", nil
}
