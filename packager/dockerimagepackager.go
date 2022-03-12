package packager

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

const dockerImageJobType = "docker"

type DockerImageJob struct{}

func NewDockerImageJob() PackageJob {
	b := &DockerImageJob{}

	return b
}

// Type returns the job type
func (b *DockerImageJob) Type() string {
	return dockerImageJobType
}

// Package packages the application
func (b *DockerImageJob) Package(log util.FriendlyLogger, ctx *project.Context) error {
	if err := ctx.HasDockerfile(); err != nil {
		return errors.Wrap(err, "missing Dockerfile")
	}

	if !ctx.Bundle.Exists {
		return errors.New("missing project bundle")
	}

	os.Setenv("DOCKER_BUILDKIT", "0")

	if _, err := util.Run(fmt.Sprintf("docker build . -t=%s:%s", ctx.Directive.Identifier, ctx.Directive.AppVersion)); err != nil {
		return errors.Wrap(err, "🚫 failed to build Docker image")
	}

	util.LogDone(fmt.Sprintf("built Docker image -> %s:%s", ctx.Directive.Identifier, ctx.Directive.AppVersion))

	return nil
}
