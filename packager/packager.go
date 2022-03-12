package packager

import (
	"github.com/pkg/errors"
	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

// Packager is responsible for packaging and publishing projects
type Packager struct {
	log util.FriendlyLogger
}

// PackageJob represents a specific type of packaging,
// for example modules into bundle, bundle into container image, etc
type PackageJob interface {
	Type() string
	Package(util.FriendlyLogger, *project.Context) error
}

// PublishJob represents an attempt to publish a packaged application
type PublishJob interface {
	Type() string
	Publish(util.FriendlyLogger, *project.Context) error
}

// New creates a new Packager
func New(log util.FriendlyLogger) *Packager {
	p := &Packager{
		log: log,
	}

	return p
}

// Package executes the given set of PackageJobs, returning an error if any fail
func (p *Packager) Package(ctx *project.Context, jobs ...PackageJob) error {
	for _, j := range jobs {
		if err := j.Package(p.log, ctx); err != nil {
			return errors.Wrapf(err, "package job %s failed", j.Type())
		}
	}

	return nil
}

// Publish executes a PublishJob
func (p *Packager) Publish(ctx *project.Context, job PublishJob) error {
	if err := job.Publish(p.log, ctx); err != nil {
		return errors.Wrapf(err, "publish job %s failed", job.Type())
	}

	return nil
}
