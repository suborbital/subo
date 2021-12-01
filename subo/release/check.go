package release

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/suborbital/subo/subo/util"
)

func getTimestampCache() (bool, error) {
	cachePath, err := util.GetCacheDir()
	if err != nil {
		return false, errors.Wrap(err, "failed to GetCacheDir")
	}

	var cached_timestamp time.Time
	filePath := filepath.Join(cachePath, "latest_version_check.txt")
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		cached_timestamp = time.Time{}
	} else if err != nil {
		return false, errors.Wrap(err, "failed to Stat")
	} else {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return false, errors.Wrap(err, "failed to ReadFile")
		}

		cached_timestamp, err = time.Parse(time.RFC3339, string(data))
		if err != nil {
			return false, errors.Wrap(err, "failed to parse cached timestamp")
		}
	}

	// check if 1 hour has passed since the last version check, and update the cached timestamp if so
	current_timestamp := time.Now().UTC()
	if cached_timestamp.IsZero() || current_timestamp.After(cached_timestamp.Add(time.Duration(1)*time.Hour)) {
		data := []byte(current_timestamp.Format(time.RFC3339))
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
		return fmt.Sprintf("upgrade subo %s to the latest release %s\n", cmdVersion, latestCmdVersion), nil
	}

	return "", nil
}
