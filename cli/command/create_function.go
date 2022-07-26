package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/velo/builder/template"
	"github.com/suborbital/velo/cli/release"
	"github.com/suborbital/velo/cli/util"
	"github.com/suborbital/velo/project"
)

const (
	tmplTypeData = "data"
)

// langAliases are aliases for languages.
var langAliases = map[string]string{
	"as": "assemblyscript",
	"rs": "rust",
	"go": "tinygo",
	"gr": "grain",
	"ts": "typescript",
	"js": "javascript",
}

// CreateRunnableError wraps errors for CreateRunnableCmd() failures.
type CreateRunnableError struct {
	Path  string // The ouput directory for build command CreateRunnableCmd().
	error        // The original error.
}

// NewCreateRunnableError cleans up and returns CreateRunnableError for CreateRunnableCmd() failures.
func NewCreateRunnableError(path string, err error) CreateRunnableError {
	if cleanupErr := os.RemoveAll(path); cleanupErr != nil {
		err = errors.Wrap(err, "failed to clean up runnable outputs")
	}
	return CreateRunnableError{Path: path, error: err}
}

// CreateFunctionCmd returns the build command.
func CreateFunctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "function <name>",
		Aliases: []string{"fn"},
		Short:   "create a new function",
		Long:    `create a new function to be used with Velocity`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			namespace, _ := cmd.Flags().GetString(namespaceFlag)
			lang, _ := cmd.Flags().GetString(langFlag)
			tmplType, _ := cmd.Flags().GetString(typeFlag)
			repo, _ := cmd.Flags().GetString(repoFlag)
			branch, _ := cmd.Flags().GetString(branchFlag)

			actualLang := lang
			if val, exists := langAliases[lang]; exists {
				actualLang = val
			}

			dir, _ := cmd.Flags().GetString(dirFlag)
			bctx, err := project.ForDirectory(dir)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to project.ForDirectory")
			}

			if bctx.RunnableExists(name) {
				return fmt.Errorf("🚫 runnable %s already exists", name)
			}

			util.LogStart(fmt.Sprintf("creating runnable %s", name))

			path, err := util.Mkdir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			runnable, err := writeDotRunnable(bctx.Cwd, name, actualLang, namespace)
			if err != nil {
				return errors.Wrap(NewCreateRunnableError(path, err), "🚫 failed to writeDotRunnable")
			}

			templatesPath, err := template.FullPath(repo, branch)
			if err != nil {
				return errors.Wrap(NewCreateRunnableError(path, err), "failed to template.FullPath")
			}

			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				templatesPath, err = template.UpdateTemplates(repo, branch)
				if err != nil {
					return errors.Wrap(NewCreateRunnableError(path, err), "🚫 failed to UpdateTemplates")
				}
			}

			templateName := tmplNameForLang(actualLang, tmplType)

			if err := template.ExecRunnableTmpl(bctx.Cwd, name, templatesPath, templateName, runnable); err != nil {
				// if the templates are missing, try updating them and exec again.
				if err == template.ErrTemplateMissing {
					templatesPath, err = template.UpdateTemplates(repo, branch)
					if err != nil {
						return errors.Wrap(NewCreateRunnableError(path, err), "🚫 failed to UpdateTemplates")
					}

					if err := template.ExecRunnableTmpl(bctx.Cwd, name, templatesPath, templateName, runnable); err != nil {
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
	cmd.Flags().String(repoFlag, "suborbital/runnable-templates", "git repo to download templates from")
	cmd.Flags().String(branchFlag, "vmain", "git branch to download templates from")
	cmd.Flags().String(typeFlag, "data", "template type - 'data' or 'handler'")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest runnable templates")

	return cmd
}

func writeDotRunnable(cwd, name, lang, namespace string) (*directive.Runnable, error) {
	if valid := project.IsValidLang(lang); !valid {
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

	if err := ioutil.WriteFile(path, bytes, util.PermFilePrivate); err != nil {
		return nil, errors.Wrap(err, "failed to WriteFile runnable")
	}

	return runnable, nil
}

func tmplNameForLang(lang, tmplType string) string {
	if tmplType == tmplTypeData {
		return lang
	}

	return fmt.Sprintf("%s-%s", lang, tmplType)
}
