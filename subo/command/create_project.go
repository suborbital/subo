package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/builder/context"
	"github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

const (
	defaultRepo = "suborbital/subo"
)

type projectData struct {
	Name        string
	Environment string
	Headless    bool
	APIVersion  string
	AtmoVersion string
}

// CreateProjectCmd returns the build command
func CreateProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <name>",
		Short: "create a new project",
		Long:  `create a new project for Atmo or Hive`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			bctx, err := context.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			util.LogStart(fmt.Sprintf("creating project %s", name))

			path, err := util.Mkdir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			branch, _ := cmd.Flags().GetString(branchFlag)
			environment, _ := cmd.Flags().GetString(environmentFlag)
			headless, _ := cmd.Flags().GetBool(headlessFlag)

			templatesPath, err := template.TemplateFullPath(defaultRepo, branch)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to TemplateFullPath")
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
				Headless:    headless,
			}

			if err := template.ExecTmplDir(bctx.Cwd, name, templatesPath, "project", data); err != nil {
				// if the templates are missing, try updating them and exec again
				if err == template.ErrTemplateMissing {
					templatesPath, err = template.UpdateTemplates(defaultRepo, branch)
					if err != nil {
						return errors.Wrap(err, "🚫 failed to UpdateTemplates")
					}

					if err := template.ExecTmplDir(bctx.Cwd, name, templatesPath, "project", data); err != nil {
						return errors.Wrap(err, "🚫 failed to ExecTmplDir")
					}
				} else {
					return errors.Wrap(err, "🚫 failed to ExecTmplDir")
				}
			}

			util.LogDone(path)

			if _, err := util.Run(fmt.Sprintf("git init ./%s", name)); err != nil {
				return errors.Wrap(err, "🚫 failed to initialize Run git init")
			}

			return nil
		},
	}

	cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	cmd.Flags().String(environmentFlag, "com.suborbital", "project environment name (your company's reverse domain")
	cmd.Flags().Bool(headlessFlag, false, "use 'headless mode' for this project")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest templates")

	return cmd
}
