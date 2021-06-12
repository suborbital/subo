package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

type projectData struct {
	Name        string
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

			bctx, err := context.CurrentBuildContext(cwd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			util.LogStart(fmt.Sprintf("creating project %s", name))

			path, err := util.Mkdir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			branch, _ := cmd.Flags().GetString(branchFlag)

			templatesPath, err := util.TemplateFullPath(branch)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to TemplateFullPath")
			}

			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				templatesPath, err = util.UpdateTemplates(bctx, name, branch)
				if err != nil {
					return errors.Wrap(err, "🚫 failed to UpdateTemplates")
				}
			}

			data := projectData{
				Name:        name,
				APIVersion:  release.FFIVersion,
				AtmoVersion: release.AtmoVersion,
			}

			if err := util.ExecTmplDir(bctx.Cwd, name, templatesPath, "project", data); err != nil {
				// if the templates are missing, try updating them and exec again
				if err == util.ErrTemplateMissing {
					templatesPath, err = util.UpdateTemplates(bctx, name, branch)
					if err != nil {
						return errors.Wrap(err, "🚫 failed to UpdateTemplates")
					}

					if err := util.ExecTmplDir(bctx.Cwd, name, templatesPath, "project", data); err != nil {
						return errors.Wrap(err, "🚫 failed to ExecTmplDir")
					}
				} else {
					return errors.Wrap(err, "🚫 failed to ExecTmplDir")
				}
			}

			util.LogDone(path)

			if _, _, err := util.Run(fmt.Sprintf("git init ./%s", name)); err != nil {
				return errors.Wrap(err, "🚫 failed to initialize Run git init")
			}

			return nil
		},
	}

	cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest templates")

	return cmd
}
