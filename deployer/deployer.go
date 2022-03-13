package deployer

import (
	"github.com/pkg/errors"

	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

type Deployer struct {
	log util.FriendlyLogger
}

type DeployJob interface {
	Type() string
	Deploy(util.FriendlyLogger, *project.Context) error
}

// New creates a new Deployer.
func New(log util.FriendlyLogger) *Deployer {
	d := &Deployer{
		log: log,
	}

	return d
}

// Deploy executes a DeployJob.
func (d *Deployer) Deploy(ctx *project.Context, job DeployJob) error {
	if err := job.Deploy(d.log, ctx); err != nil {
		return errors.Wrapf(err, "deploy job %s failed", job.Type())
	}

	return nil
}
