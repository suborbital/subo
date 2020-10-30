package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/command/template"
	"github.com/suborbital/subo/subo/context"
	"gopkg.in/yaml.v2"
)

// CreateCmd returns the build command
func CreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "create an empty Wasm runnable",
		Long:  `create an empty Wasm runnable`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			dir, _ := cmd.Flags().GetString("dir")
			bctx, err := context.CurrentBuildContext(dir)
			if err != nil {
				return errors.Wrap(err, "failed to get CurrentBuildContext")
			}

			if bctx.RunnableExists(name) {
				return fmt.Errorf("🚫 runnable %s already exists", name)
			}

			lang, err := cmd.Flags().GetString("lang")
			if err != nil {
				return errors.Wrap(err, "failed to get lang flag")
			}

			filename, tmpl, err := template.ForLang(lang)
			if err != nil {
				return errors.Wrap(err, "failed to template.ForLang")
			}

			fmt.Println(fmt.Sprintf("✨ START: creating runnable %s", name))

			path, err := makeRunnableDir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "failed to makeRunnableDir")
			}

			if err := writeDotHive(path, name, lang); err != nil {
				return errors.Wrap(err, "failed to writeDotHive")
			}

			if err := writeTmpl(path, tmpl, filename); err != nil {
				return errors.Wrap(err, "failed to writeTmpl")
			}

			fmt.Println(fmt.Sprintf("✨ DONE: %s", path))

			return nil
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "$HOME"
	}

	cmd.Flags().String("dir", cwd, "the directory to run the build from")
	cmd.Flags().String("lang", "rust", "the language used for the runnable")

	return cmd
}

func makeRunnableDir(cwd, name string) (string, error) {
	path := filepath.Join(cwd, name)

	if err := os.Mkdir(path, 0700); err != nil {
		return "", errors.Wrap(err, "failed to Mkdir")
	}

	return path, nil
}

func writeDotHive(dir, name, lang string) error {
	dotHive := context.DotHive{
		Name: name,
		Lang: lang,
	}

	bytes, err := yaml.Marshal(dotHive)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal dotHive")
	}

	path := filepath.Join(dir, ".hive.yml")

	if err := ioutil.WriteFile(path, bytes, 0700); err != nil {
		return errors.Wrap(err, "failed to WriteFile dotHive")
	}

	return nil
}

func writeTmpl(dir, tmpl, name string) error {
	path := filepath.Join(dir, name)

	if err := ioutil.WriteFile(path, []byte(tmpl), 0700); err != nil {
		return errors.Wrap(err, "failed to WriteFile template")
	}

	return nil
}
