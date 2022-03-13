package packager

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

const (
	dockerPublishJobType = "docker"
)

type DockerPublishJob struct{}

// NewDockerPublishJob returns a new PublishJob for Bindle.
func NewDockerPublishJob() PublishJob {
	d := &DockerPublishJob{}

	return d
}

// Type returns the publish job's type.
func (b *DockerPublishJob) Type() string {
	return dockerPublishJobType
}

// Publish publishes the application.
func (b *DockerPublishJob) Publish(log util.FriendlyLogger, ctx *project.Context) error {
	if ctx.Directive == nil {
		return errors.New("cannot publish without Directive.yaml")
	}

	if !ctx.Bundle.Exists {
		return errors.New("cannot publish without runnables.wasm.zip, run `subo build` first")
	}

	imagesCmd := "docker images --format '{{ .Repository }}:{{ .Tag }}'"

	imagesOutput, err := util.RunSilent(imagesCmd)
	if err != nil {
		return errors.Wrap(err, "failed to Run images command")
	}

	imageName, err := project.DockerNameFromDirective(ctx.Directive)
	if err != nil {
		return errors.Wrap(err, "failed to dockerNameFromDirective")
	}

	if !strings.Contains(imagesOutput, imageName) {
		return fmt.Errorf("image %s not found, run `subo build --docker` first", imageName)
	}

	if _, err := util.Run(fmt.Sprintf("docker push %s", imageName)); err != nil {
		return errors.Wrap(err, "failed to Run docker push")
	}

	util.LogDone(fmt.Sprintf("pushed Docker image -> %s", imageName))

	return nil
}
