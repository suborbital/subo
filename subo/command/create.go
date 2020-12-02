package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/release"
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
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			if bctx.RunnableExists(name) {
				return fmt.Errorf("🚫 runnable %s already exists", name)
			}

			lang, _ := cmd.Flags().GetString("lang")

			fmt.Println(fmt.Sprintf("✨ START: creating runnable %s", name))

			path, err := makeRunnableDir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "failed to makeRunnableDir")
			}

			namespace, _ := cmd.Flags().GetString("namespace")

			dotHive, err := writeDotHive(bctx.Cwd, name, lang, namespace)
			if err != nil {
				return errors.Wrap(err, "failed to writeDotHive")
			}

			tmplPath := filepath.Join(os.TempDir(), "suborbital", "templates")

			if err := copyTmpl(bctx.Cwd, name, tmplPath, dotHive); err != nil {
				return errors.Wrap(err, "failed to copyTmpl")
			}

			fmt.Println(fmt.Sprintf("✨ DONE: %s", path))

			return nil
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "$HOME"
	}

	cmd.Flags().String("dir", cwd, "the directory to put the new runnable in")
	cmd.Flags().String("lang", "rust", "the language of the new runnable")
	cmd.Flags().String("namespace", "default", "the namespace for the new runnable")

	return cmd
}

func makeRunnableDir(cwd, name string) (string, error) {
	path := filepath.Join(cwd, name)

	if err := os.Mkdir(path, 0700); err != nil {
		return "", errors.Wrap(err, "failed to Mkdir")
	}

	return path, nil
}

func writeDotHive(cwd, name, lang, namespace string) (*context.DotHive, error) {
	dotHive := &context.DotHive{
		Name:       name,
		Lang:       lang,
		Namespace:  namespace,
		APIVersion: release.FFIVersion,
	}

	bytes, err := yaml.Marshal(dotHive)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal dotHive")
	}

	path := filepath.Join(cwd, name, ".hive.yml")

	if err := ioutil.WriteFile(path, bytes, 0700); err != nil {
		return nil, errors.Wrap(err, "failed to WriteFile dotHive")
	}

	return dotHive, nil
}

type tmplData struct {
	context.DotHive
	NameCaps  string
	NameCamel string
}

func copyTmpl(cwd, name, templatesPath string, dotHive *context.DotHive) error {
	templatePath := filepath.Join(templatesPath, dotHive.Lang)
	targetPath := filepath.Join(cwd, name)

	if _, err := os.Stat(templatePath); err != nil {
		if err == os.ErrNotExist {
			return fmt.Errorf("template for %s does not exist", dotHive.Lang)
		}

		return errors.Wrap(err, "failed to Stat template directory")
	}

	nameCamel := ""

	nameParts := strings.Split(dotHive.Name, "-")
	for _, part := range nameParts {
		nameCamel += strings.ToUpper(string(part[0]))
		nameCamel += string(part[1:])
	}

	templateData := tmplData{
		DotHive:   *dotHive,
		NameCaps:  strings.ToUpper(strings.Replace(dotHive.Name, "-", "", -1)),
		NameCamel: nameCamel,
	}

	var err error = filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, templatePath, "", 1)
		if relPath == "" {
			return nil
		}

		if info.IsDir() {
			return os.Mkdir(filepath.Join(targetPath, relPath), 0755)
		}

		var data, err1 = ioutil.ReadFile(filepath.Join(templatePath, relPath))
		if err1 != nil {
			return err1
		}

		if strings.HasSuffix(info.Name(), ".tmpl") {
			tmpl, err := template.New("tmpl").Parse(string(data))
			if err != nil {
				return errors.Wrapf(err, "failed to parse template file %s", info.Name())
			}

			builder := &strings.Builder{}
			if err := tmpl.Execute(builder, templateData); err != nil {
				return errors.Wrapf(err, "failed to Execute template for %s", info.Name())
			}

			data = []byte(builder.String())
			relPath = strings.Replace(relPath, ".tmpl", "", 1)
		}

		return ioutil.WriteFile(filepath.Join(targetPath, relPath), data, 0777)
	})

	return err
}
