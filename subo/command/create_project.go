package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

			logStart(fmt.Sprintf("creating project %s", name))

			path, err := util.Mkdir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			templateRootPath, err := util.TemplateDir()
			if err != nil {
				return errors.Wrap(err, "failed to TemplateDir")
			}

			data := projectData{
				Name:        name,
				APIVersion:  release.FFIVersion,
				AtmoVersion: release.AtmoVersion,
			}

			branch, _ := cmd.Flags().GetString(branchFlag)
			branchDirName := fmt.Sprintf("subo-%s", strings.ReplaceAll(branch, "/", "-"))
			tmplPath := filepath.Join(templateRootPath, branchDirName, "templates")

			// encapsulate this in a function so it can be called if updates are requested or if the first attempt to copy fails
			updateAndCopy := func() error {
				logStart("downloading templates")

				filepath, err := util.DownloadZip(branch, templateRootPath)
				if err != nil {
					return errors.Wrap(err, "🚫 failed to downloadZip for templates")
				}

				// tmplPath may be different than the default if a custom URL was provided
				tmplPath, err = util.ExtractZip(filepath, templateRootPath, branchDirName)
				if err != nil {
					return errors.Wrap(err, "🚫 failed to extractZip for templates")
				}

				logDone("templates downloaded")

				if err = util.ExecTmplDir(bctx.Cwd, name, tmplPath, "project", data); err != nil {
					return errors.Wrap(err, "🚫 failed to copyTmpl")
				}

				return nil
			}

			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				if err := updateAndCopy(); err != nil {
					return errors.Wrap(err, "failed to updateAndCopy")
				}
			} else {
				if err := util.ExecTmplDir(bctx.Cwd, name, tmplPath, "project", data); err != nil {
					if err == util.ErrTemplateMissing {
						if err := updateAndCopy(); err != nil {
							return errors.Wrap(err, "failed to updateAndCopy")
						}
					} else {
						return errors.Wrap(err, "🚫 failed to copyTmpl")
					}
				}
			}

			logDone(path)

			return nil
		},
	}

	cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest templates")

	return cmd
}
