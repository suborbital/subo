package packager

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

const dockerImageJobType = "docker"

type DockerImagePackageJob struct{}

func NewDockerImagePackageJob() PackageJob {
	b := &DockerImagePackageJob{}

	return b
}

// Type returns the job type.
func (b *DockerImagePackageJob) Type() string {
	return dockerImageJobType
}

// Package packages the application.
func (b *DockerImagePackageJob) Package(log util.FriendlyLogger, ctx *project.Context) error {
	if err := ctx.HasDockerfile(); err != nil {
		return errors.Wrap(err, "missing Dockerfile")
	}

	if !ctx.Bundle.Exists {
		return errors.New("missing project bundle")
	}

	imageName, err := project.DockerNameFromDirective(ctx.Directive)
	if err != nil {
		return errors.Wrap(err, "failed to dockerNameFromDirective")
	}

	os.Setenv("DOCKER_BUILDKIT", "0")

	if _, err := util.Run(fmt.Sprintf("docker build . -t=%s", imageName)); err != nil {
		return errors.Wrap(err, "🚫 failed to build Docker image")
	}

	util.LogDone(fmt.Sprintf("built Docker image -> %s", imageName))

	return nil
}
