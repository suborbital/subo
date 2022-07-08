package deployer

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/suborbital/velo/cli/util"
	"github.com/suborbital/velo/project"
)

const (
	flyDeployJobType = "fly.io"
)

// FlyDeployJob represents a deployment job.
type FlyDeployJob struct {
	org    string
	region string
	local  bool
}

// NewFlyDeployJob creates a new deploy job.
func NewFlyDeployJob(org, region string, local bool) DeployJob {
	k := &FlyDeployJob{
		org:    org,
		region: region,
		local:  local,
	}

	return k
}

// Typw returns the deploy job typw.
func (f *FlyDeployJob) Type() string {
	return flyDeployJobType
}

// Deploy executes the deployment.
func (f *FlyDeployJob) Deploy(log util.FriendlyLogger, ctx *project.Context) error {
	nameDashes := strings.Replace(ctx.Directive.Identifier, ".", "-", -1)

	cmd := fmt.Sprintf("flyctl launch --name %s --org %s --region %s --now", nameDashes, f.org, f.region)

	if f.local {
		if _, err := util.Command.Run("flyctl auth docker"); err != nil {
			return errors.Wrap(err, "failed to check Fly.io authentication, run 'flyctl auth login' and try again")
		}

		if _, err := util.Command.Run(fmt.Sprintf("docker buildx build . --platform linux/amd64 --push -t registry.fly.io/%s:%s", nameDashes, ctx.Directive.AppVersion)); err != nil {
			return errors.Wrap(err, "failed to build and push image")
		}

		cmd = fmt.Sprintf("flyctl deploy -i registry.fly.io/%s:%s --name %s --org %s --region %s --now", nameDashes, ctx.Directive.AppVersion, nameDashes, f.org, f.region)
	} else {
		if _, err := util.Command.Run("flyctl auth whoami"); err != nil {
			return errors.Wrap(err, "failed to check Fly.io authentication, run 'flyctl auth login' and try again")
		}

		if _, err := os.Stat("./fly.toml"); err == nil {
			cmd = "flyctl deploy --remote-only --now"
		}
	}

	if _, err := util.Command.Run(cmd); err != nil {
		return errors.Wrap(err, "failed to deploy")
	}

	return nil
}
