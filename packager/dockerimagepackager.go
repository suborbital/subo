package packager

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

const dockerImagePackageJobType = "docker"

type DockerImagePackageJob struct{}

func NewDockerImagePackageJob() PackageJob {
	b := &DockerImagePackageJob{}

	return b
}

// Type returns the job type.
func (b *DockerImagePackageJob) Type() string {
	return dockerImagePackageJobType
}

// Package packages the application.
func (b *DockerImagePackageJob) Package(log util.FriendlyLogger, ctx *project.Context) error {
	if err := ctx.HasDockerfile(); err != nil {
		return errors.Wrap(err, "missing Dockerfile")
	}

	if !ctx.Bundle.Exists {
		return errors.New("missing project bundle")
	}

	if err := os.Setenv("DOCKER_BUILDKIT", "0"); err != nil {
		util.LogWarn("DOCKER_BUILDKIT=0 could not be set, Docker build may be problematic on M1 Macs.")
	}

	imageName, err := dockerNameFromDirective(ctx.Directive)
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

func dockerNameFromDirective(d *directive.Directive) (string, error) {
	identParts := strings.Split(d.Identifier, ".")
	if len(identParts) != 3 {
		return "", errors.New("ident has incorrect number of parts")
	}

	org := identParts[1]
	repo := identParts[2]

	name := fmt.Sprintf("%s/%s:%s", org, repo, d.AppVersion)

	return name, nil
}
