package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
	"gopkg.in/yaml.v2"
)

const (
	langFlag            = "lang"
	dirFlag             = "dir"
	namespaceFlag       = "namespace"
	branchFlag          = "branch"
	updateTemplatesFlag = "update-templates"
)

// CreateRunnableCmd returns the build command
func CreateRunnableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runnable <name>",
		Short: "create a new Runnable",
		Long:  `create a new Runnable to be used with Atmo or Hive`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			dir, _ := cmd.Flags().GetString(dirFlag)
			bctx, err := context.CurrentBuildContext(dir)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			if bctx.RunnableExists(name) {
				return fmt.Errorf("🚫 runnable %s already exists", name)
			}

			lang, _ := cmd.Flags().GetString(langFlag)

			logStart(fmt.Sprintf("creating runnable %s", name))

			path, err := util.Mkdir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			namespace, _ := cmd.Flags().GetString(namespaceFlag)

			runnable, err := writeDotRunnable(bctx.Cwd, name, lang, namespace)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to writeDotHive")
			}

			templateRootPath, err := util.TemplateDir()
			if err != nil {
				return errors.Wrap(err, "failed to TemplateDir")
			}

			branch, _ := cmd.Flags().GetString(branchFlag)
			branchDirName := fmt.Sprintf("subo-%s", strings.ReplaceAll(branch, "/", "-"))
			tmplPath := filepath.Join(templateRootPath, branchDirName, "templates")

			// encapsulate this in a function so it can be called if updates are requested or if the first attempt to copy fails
			updateAndCopy := func() error {
				logStart("downloading runnable templates")

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

				if err = util.CopyRunnableTmpl(bctx.Cwd, name, tmplPath, runnable); err != nil {
					return errors.Wrap(err, "🚫 failed to copyTmpl")
				}

				return nil
			}

			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				if err := updateAndCopy(); err != nil {
					return errors.Wrap(err, "failed to updateAndCopy")
				}
			} else {
				if err := util.CopyRunnableTmpl(bctx.Cwd, name, tmplPath, runnable); err != nil {
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

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "$HOME"
	}

	cmd.Flags().String(dirFlag, cwd, "the directory to put the new runnable in")
	cmd.Flags().String(langFlag, "rust", "the language of the new runnable")
	cmd.Flags().String(namespaceFlag, "default", "the namespace for the new runnable")
	cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest runnable templates")

	return cmd
}

func writeDotRunnable(cwd, name, lang, namespace string) (*directive.Runnable, error) {
	runnable := &directive.Runnable{
		Name:       name,
		Lang:       lang,
		Namespace:  namespace,
		APIVersion: release.FFIVersion,
	}

	bytes, err := yaml.Marshal(runnable)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal runnable")
	}

	path := filepath.Join(cwd, name, ".runnable.yaml")

	if err := ioutil.WriteFile(path, bytes, 0700); err != nil {
		return nil, errors.Wrap(err, "failed to WriteFile runnable")
	}

	return runnable, nil
}
