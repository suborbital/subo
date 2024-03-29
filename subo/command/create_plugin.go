package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
	"github.com/suborbital/systemspec/tenant"
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

// CreatePluginError wraps errors for CreatePluginCmd() failures.
type CreatePluginError struct {
	Path  string // The ouput directory for build command CreatePluginCmd().
	error        // The original error.
}

// NewCreatePluginError cleans up and returns CreatePluginError for CreatePluginCmd() failures.
func NewCreatePluginError(path string, err error) CreatePluginError {
	if cleanupErr := os.RemoveAll(path); cleanupErr != nil {
		err = errors.Wrap(err, "failed to clean up plugin outputs")
	}
	return CreatePluginError{Path: path, error: err}
}

// CreatePluginCmd returns the build command.
func CreatePluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin <name>",
		Short: "create a new plugin",
		Long:  `create a new plugin to be used with E2Core or Suborbital Extension Engine (SE2)`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			namespace, _ := cmd.Flags().GetString(namespaceFlag)
			lang, _ := cmd.Flags().GetString(langFlag)
			repo, _ := cmd.Flags().GetString(repoFlag)
			branch, _ := cmd.Flags().GetString(branchFlag)

			dir, _ := cmd.Flags().GetString(dirFlag)
			bctx, err := project.ForDirectory(dir)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to project.ForDirectory")
			}

			if bctx.ModuleExists(name) {
				return fmt.Errorf("🚫 plugin %s already exists", name)
			}

			util.LogStart(fmt.Sprintf("creating plugin %s", name))

			path, err := util.Mkdir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			module, err := writeDotModule(bctx.Cwd, name, lang, namespace)
			if err != nil {
				return errors.Wrap(NewCreatePluginError(path, err), "🚫 failed to writeDotModule")
			}

			templatesPath, err := template.FullPath(repo, branch)
			if err != nil {
				return errors.Wrap(NewCreatePluginError(path, err), "failed to template.FullPath")
			}

			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				templatesPath, err = template.UpdateTemplates(repo, branch)
				if err != nil {
					return errors.Wrap(NewCreatePluginError(path, err), "🚫 failed to UpdateTemplates")
				}
			}

			if err := template.ExecModuleTmpl(bctx.Cwd, name, templatesPath, module); err != nil {
				// if the templates are missing, try updating them and exec again.
				if err == template.ErrTemplateMissing {
					templatesPath, err = template.UpdateTemplates(repo, branch)
					if err != nil {
						return errors.Wrap(NewCreatePluginError(path, err), "🚫 failed to UpdateTemplates")
					}

					if err := template.ExecModuleTmpl(bctx.Cwd, name, templatesPath, module); err != nil {
						return errors.Wrap(NewCreatePluginError(path, err), "🚫 failed to ExecTmplDir")
					}
				} else {
					return errors.Wrap(NewCreatePluginError(path, err), "🚫 failed to ExecTmplDir")
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

	cmd.Flags().String(dirFlag, cwd, "the directory to put the new module in")
	cmd.Flags().String(langFlag, "rust", "the language of the new module")
	cmd.Flags().String(namespaceFlag, "default", "the namespace for the new module")
	cmd.Flags().String(repoFlag, defaultRepo, "git repo to download templates from")
	cmd.Flags().String(branchFlag, defaultBranch, "git branch to download templates from")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest module templates")

	return cmd
}

func writeDotModule(cwd, name, lang, namespace string) (*tenant.Module, error) {
	if actual, exists := langAliases[lang]; exists {
		lang = actual
	}

	if valid := project.IsValidLang(lang); !valid {
		return nil, fmt.Errorf("%s is not an available language", lang)
	}

	module := &tenant.Module{
		Name:       name,
		Lang:       lang,
		Namespace:  namespace,
		APIVersion: release.SDKVersion,
	}

	bytes, err := yaml.Marshal(module)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal module")
	}

	path := filepath.Join(cwd, name, ".module.yaml")

	if err := ioutil.WriteFile(path, bytes, util.PermFilePrivate); err != nil {
		return nil, errors.Wrap(err, "failed to WriteFile module")
	}

	return module, nil
}
