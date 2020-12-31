package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/hive-wasm/directive"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
	"gopkg.in/yaml.v2"
)

var errTemplateMissing = errors.New("template missing")

const (
	langFlag            = "lang"
	dirFlag             = "dir"
	namespaceFlag       = "namespace"
	branchFlag          = "branch"
	updateTemplatesFlag = "update-templates"
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

			path, err := makeRunnableDir(bctx.Cwd, name)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to makeRunnableDir")
			}

			namespace, _ := cmd.Flags().GetString(namespaceFlag)

			runnable, err := writeDotHive(bctx.Cwd, name, lang, namespace)
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

				filepath, err := downloadZip(branch, templateRootPath)
				if err != nil {
					return errors.Wrap(err, "🚫 failed to downloadZip for templates")
				}

				// tmplPath may be different than the default if a custom URL was provided
				tmplPath, err = extractZip(filepath, templateRootPath, branchDirName)
				if err != nil {
					return errors.Wrap(err, "🚫 failed to extractZip for templates")
				}

				logDone("templates downloaded")

				if err = copyTmpl(bctx.Cwd, name, tmplPath, runnable); err != nil {
					return errors.Wrap(err, "🚫 failed to copyTmpl")
				}

				return nil
			}

			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				if err := updateAndCopy(); err != nil {
					return errors.Wrap(err, "failed to updateAndCopy")
				}
			} else {
				if err := copyTmpl(bctx.Cwd, name, tmplPath, runnable); err != nil {
					if err == errTemplateMissing {
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

func makeRunnableDir(cwd, name string) (string, error) {
	path := filepath.Join(cwd, name)

	if err := os.Mkdir(path, 0700); err != nil {
		return "", errors.Wrap(err, "failed to Mkdir")
	}

	return path, nil
}

func writeDotHive(cwd, name, lang, namespace string) (*directive.Runnable, error) {
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

type tmplData struct {
	directive.Runnable
	NameCaps  string
	NameCamel string
}

func copyTmpl(cwd, name, templatesPath string, runnable *directive.Runnable) error {
	templatePath := filepath.Join(templatesPath, runnable.Lang)
	targetPath := filepath.Join(cwd, name)

	if _, err := os.Stat(templatePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errTemplateMissing
		}

		return errors.Wrap(err, "failed to Stat template directory")
	}

	nameCamel := ""

	nameParts := strings.Split(runnable.Name, "-")
	for _, part := range nameParts {
		nameCamel += strings.ToUpper(string(part[0]))
		nameCamel += string(part[1:])
	}

	templateData := tmplData{
		Runnable:  *runnable,
		NameCaps:  strings.ToUpper(strings.Replace(runnable.Name, "-", "", -1)),
		NameCamel: nameCamel,
	}

	var err error = filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, templatePath, "", 1)
		if relPath == "" {
			return nil
		}

		targetRelPath := relPath
		if strings.Contains(relPath, ".tmpl") {
			tmpl, err := template.New("tmpl").Parse(strings.Replace(relPath, ".tmpl", "", -1))
			if err != nil {
				return errors.Wrapf(err, "failed to parse template directory name %s", info.Name())
			}

			builder := &strings.Builder{}
			if err := tmpl.Execute(builder, templateData); err != nil {
				return errors.Wrapf(err, "failed to Execute template for %s", info.Name())
			}

			targetRelPath = builder.String()
		}

		if info.IsDir() {
			return os.Mkdir(filepath.Join(targetPath, targetRelPath), 0755)
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
		}

		return ioutil.WriteFile(filepath.Join(targetPath, targetRelPath), data, 0777)
	})

	return err
}

func downloadZip(branch, targetPath string) (string, error) {
	url := fmt.Sprintf("https://github.com/suborbital/subo/archive/%s.zip", branch)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to NewRequest")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to Do request")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response was non-200: %d", resp.StatusCode)
	}

	filePath := filepath.Join(targetPath, "subo.zip")

	// check if the zip already exists, and delete it if it does
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return "", errors.Wrap(err, "failed to delete exising templates zip")
		}
	}

	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to MkdirAll")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to Open file")
	}

	defer resp.Body.Close()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", errors.Wrap(err, "failed to Copy data to file")
	}

	return filePath, nil
}

func extractZip(filePath, destPath, branchDirName string) (string, error) {
	escapedFilepath := strings.ReplaceAll(filePath, " ", "\\ ")
	escapedDestpath := strings.ReplaceAll(destPath, " ", "\\ ") + string(filepath.Separator)

	if _, _, err := util.Run(fmt.Sprintf("unzip -q %s -d %s", escapedFilepath, escapedDestpath)); err != nil {
		return "", errors.Wrap(err, "failed to Run unzip")
	}

	files, err := ioutil.ReadDir(destPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to ReadDir")
	}

	for _, f := range files {
		if f.IsDir() && f.Name() == branchDirName {
			return filepath.Join(destPath, f.Name(), "templates"), nil
		}
	}

	return "", errors.New("templates not availale")
}
