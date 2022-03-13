package deployer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

const (
	k8sDeployJobType = "kubernetes"
)

// K8sDeployJob represents a deployment job.
type K8sDeployJob struct {
	repo            string
	branch          string
	domain          string
	updateTemplates bool
}

type deploymentData struct {
	Identifier string
	AppVersion string
	ImageName  string
	Domain     string
}

// NewK8sDeployJob creates a new deploy job.
func NewK8sDeployJon(repo, branch, domain string, updateTemplates bool) DeployJob {
	k := &K8sDeployJob{
		repo:            repo,
		branch:          branch,
		domain:          domain,
		updateTemplates: updateTemplates,
	}

	return k
}

// Typw returns the deploy job typw.
func (k *K8sDeployJob) Type() string {
	return k8sDeployJobType
}

// Deploy executes the deployment.
func (k *K8sDeployJob) Deploy(log util.FriendlyLogger, ctx *project.Context) error {
	imageName, err := project.DockerNameFromDirective(ctx.Directive)
	if err != nil {
		return errors.Wrap(err, "failed to DockerNameFromDirective")
	}

	data := deploymentData{
		Identifier: strings.Replace(ctx.Directive.Identifier, ".", "-", -1),
		AppVersion: ctx.Directive.AppVersion,
		ImageName:  imageName,
		Domain:     k.domain,
	}

	if err := os.RemoveAll(filepath.Join(ctx.Cwd, ".deployment")); err != nil {
		if os.IsNotExist(err) {
			//that's fine.
		} else {
			return errors.Wrap(err, "failed to RemoveAll")
		}
	}

	if err := os.MkdirAll(filepath.Join(ctx.Cwd, ".deployment"), 0700); err != nil {
		return errors.Wrap(err, "failed to MkdirAll")
	}

	templatesPath, err := template.TemplateFullPath(k.repo, k.branch)
	if err != nil {
		return errors.Wrap(err, "failed to TemplateDir")
	}

	if k.updateTemplates {
		templatesPath, err = template.UpdateTemplates(k.repo, k.branch)
		if err != nil {
			return errors.Wrap(err, "🚫 failed to UpdateTemplates")
		}
	}

	if err := template.ExecTmplDir(ctx.Cwd, ".deployment", templatesPath, "k8s", data); err != nil {
		// if the templates are missing, try updating them and exec again.
		if err == template.ErrTemplateMissing {
			templatesPath, err = template.UpdateTemplates(k.repo, k.branch)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to UpdateTemplates")
			}

			if err := template.ExecTmplDir(ctx.Cwd, ".deployment", templatesPath, "k8s", data); err != nil {
				return errors.Wrap(err, "🚫 failed to ExecTmplDir")
			}
		} else {
			return errors.Wrap(err, "🚫 failed to ExecTmplDir")
		}
	}

	// if this fails, that's ok (the ns may already exist).
	util.Run("kubectl create ns suborbital")

	if _, err := util.Run("kubectl apply -f .deployment/"); err != nil {
		return errors.Wrap(err, "failed to Run kubectl apply")
	}

	return nil
}
