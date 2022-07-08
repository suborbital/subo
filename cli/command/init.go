package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/suborbital/velo/builder/template"
	"github.com/suborbital/velo/cli/release"
	"github.com/suborbital/velo/cli/util"
	"github.com/suborbital/velo/project"
)

const (
	defaultRepo   = "suborbital/runnable-templates"
	defaultBranch = "vmain"
)

type projectData struct {
	Name        string
	Environment string
	Headless    bool
	APIVersion  string
	AtmoVersion string
}

// InitCmd returns the init command.
func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <name>",
		Short: "initialize a new project",
		Long:  `initialize a new project for Velocity`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			bctx, err := project.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to project.ForDirectory")
			}

			name := args[0]
			path := name

			util.LogStart(fmt.Sprintf("creating project %s", name))

			if name != "." {
				if _, err = util.Mkdir(bctx.Cwd, name); err != nil {
					return errors.Wrap(err, "🚫 failed to Mkdir")
				}
			} else {
				// if the provided name was '.', then set the project name to be the current dir name
				pathElements := strings.Split(cwd, string(filepath.Separator))
				name = pathElements[len(pathElements)-1]
			}

			branch, _ := cmd.Flags().GetString(branchFlag)
			environment, _ := cmd.Flags().GetString(environmentFlag)
			projectType, _ := cmd.Flags().GetString(typeFlag)

			templatesPath, err := template.FullPath(defaultRepo, branch)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to template.FullPath")
			}

			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				templatesPath, err = template.UpdateTemplates(defaultRepo, branch)
				if err != nil {
					return errors.Wrap(err, "🚫 failed to UpdateTemplates")
				}
			}

			data := projectData{
				Name:        name,
				Environment: environment,
				APIVersion:  release.FFIVersion,
				AtmoVersion: release.AtmoVersion,
			}

			templateName := "project"
			if projectType != "" {
				templateName = fmt.Sprintf("project-%s", projectType)
			}

			if err := template.ExecTmplDir(bctx.Cwd, path, templatesPath, templateName, data); err != nil {
				// if the templates are missing, try updating them and exec again.
				if err == template.ErrTemplateMissing {
					templatesPath, err = template.UpdateTemplates(defaultRepo, branch)
					if err != nil {
						return errors.Wrap(err, "🚫 failed to UpdateTemplates")
					}

					if err := template.ExecTmplDir(bctx.Cwd, path, templatesPath, templateName, data); err != nil {
						return errors.Wrap(err, "🚫 failed to ExecTmplDir")
					}
				} else {
					return errors.Wrap(err, "🚫 failed to ExecTmplDir")
				}
			}

			util.LogDone(name)

			if _, err := util.Command.Run(fmt.Sprintf("git init ./%s", name)); err != nil {
				return errors.Wrap(err, "🚫 failed to initialize Run git init")
			}

			return nil
		},
	}

	cmd.Flags().String(branchFlag, "vmain", "git branch to download templates from")
	cmd.Flags().String(typeFlag, "", "type of project to create, such as 'js'")
	cmd.Flags().String(environmentFlag, "com.suborbital", "project environment name (your company's reverse domain")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest templates")

	return cmd
}
