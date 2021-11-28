package release

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

// SuboDotVersion represents the dot version for subo
// it is also the image tag used for builders
var SuboDotVersion = "0.2.1"

// FFIVersion is the FFI version used by this version of subo
var FFIVersion = "0.13.1"

// AtmoVersion is the default version of Atmo that will be used for new projects
var AtmoVersion = "0.4.2"

// SCCTag is the docker tag used for creating new compute core deployments
var SCCTag = "v0.1.0"

// CheckForLatestVersion returns an error if SuboDotVersion does not match the latest GitHub release or if the check fails
func CheckForLatestVersion() error {
	latestRepoRelease, _, err := github.NewClient(nil).Repositories.GetLatestRelease(context.Background(), "suborbital", "subo")
	if err != nil {
		return errors.Wrap(err, "failed to fetch latest subo release")
	}
	latestCmdVersion, err := version.NewVersion(*latestRepoRelease.TagName)
	if err != nil {
		return errors.Wrap(err, "failed to parse latest subo version")
	}
	cmdVersion, err := version.NewVersion(SuboDotVersion)
	if err != nil {
		return errors.Wrap(err, "failed to parse current subo version")
	} else if cmdVersion.LessThan(latestCmdVersion) {
		return errors.New(fmt.Sprintf("upgrade subo %s to the latest release %s\n", cmdVersion, latestCmdVersion))
	}
	return err
}
