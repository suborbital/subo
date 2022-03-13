package deployer

import "github.com/suborbital/subo/project"

type Deployer struct{}

type DeployJob interface {
	Type() string
	Deploy(project.Context) error
}
