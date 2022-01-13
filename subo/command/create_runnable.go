package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/subo/builder/context"
	"github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
	"gopkg.in/yaml.v2"
)

// validLangs are the available languages
var validLangs = map[string]bool{
	"rust":           true,
	"swift":          true,
	"assemblyscript": true,
	"tinygo":         true,
	"grain":          true,
}

// langAliases are aliases for languages
var langAliases = map[string]string{
	"typescript": "assemblyscript",
	"rs":         "rust",
	"go":         "tinygo",
	"gr":         "grain",
}

// CreateRunnableError wraps errors for CreateRunnableCmd() failures
type CreateRunnableError struct {
	Path  string // The ouput directory for build command CreateRunnableCmd().
	error        // The original error.
}

// NewCreateRunnableError cleans up and returns CreateRunnableError for CreateRunnableCmd() failures
func NewCreateRunnableError(path string, err error) CreateRunnableError {
	if cleanupErr := os.RemoveAll(path); cleanupErr != nil {
		err = errors.Wrap(err, "failed to clean up runnable outputs")
	}
	return CreateRunnableError{Path: path, error: err}
}

// CreateRunnableCmd returns the build command
func CreateRunnableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runnable <name>",
		Short: "create a new Runnable",
		Long:  `create a new Runnable to be used with Atmo or Reactr`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			namespace, _ := cmd.Flags().GetString(namespaceFlag)
			lang, _ := cmd.Flags().GetString(langFlag)
			repo, _ := cmd.Flags().GetString(repoFlag)
			branch, _ := cmd.Flags().GetString(branchFlag)

			dir, _ := cmd.Flags().GetString(dirFlag)
			bctx, err := context.ForDirectory(dir)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			if bctx.RunnableExists(name) {
				return fmt.Errorf("🚫 runnable %s already exists", name)
			}

			util.LogStart(fmt.Sprintf("creating runnable %s", name))

			path, err := util.Mkdir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			runnable, err := writeDotRunnable(bctx.Cwd, name, lang, namespace)
			if err != nil {
				return errors.Wrap(NewCreateRunnableError(path, err), "🚫 failed to writeDotRunnable")
			}

			templatesPath, err := template.TemplateFullPath(repo, branch)
			if err != nil {
				return errors.Wrap(NewCreateRunnableError(path, err), "failed to TemplateDir")
			}

			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				templatesPath, err = template.UpdateTemplates(repo, branch)
				if err != nil {
					return errors.Wrap(NewCreateRunnableError(path, err), "🚫 failed to UpdateTemplates")
				}
			}

			if err := template.ExecRunnableTmpl(bctx.Cwd, name, templatesPath, runnable); err != nil {
				// if the templates are missing, try updating them and exec again
				if err == template.ErrTemplateMissing {
					templatesPath, err = template.UpdateTemplates(repo, branch)
					if err != nil {
						return errors.Wrap(NewCreateRunnableError(path, err), "🚫 failed to UpdateTemplates")
					}

					if err := template.ExecRunnableTmpl(bctx.Cwd, name, templatesPath, runnable); err != nil {
						return errors.Wrap(NewCreateRunnableError(path, err), "🚫 failed to ExecTmplDir")
					}
				} else {
					return errors.Wrap(NewCreateRunnableError(path, err), "🚫 failed to ExecTmplDir")
				}
			}

			util.LogDone(path)

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
	cmd.Flags().String(repoFlag, "suborbital/subo", "git repo to download templates from")
	cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest runnable templates")

	return cmd
}

func writeDotRunnable(cwd, name, lang, namespace string) (*directive.Runnable, error) {
	if actual, exists := langAliases[lang]; exists {
		lang = actual
	}

	if _, valid := validLangs[lang]; !valid {
		return nil, fmt.Errorf("%s is not an available language", lang)
	}

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

	if err := ioutil.WriteFile(path, bytes, 0600); err != nil {
		return nil, errors.Wrap(err, "failed to WriteFile runnable")
	}

	return runnable, nil
}
