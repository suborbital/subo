package context

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/suborbital/hive-wasm/directive"
	"gopkg.in/yaml.v2"
)

var dockerImageForLang = map[string]string{
	"rust":  "suborbital/hive-rs:1.46.0-3",
	"swift": "UNSUPPORTED",
}

// BuildContext describes the context under which the tool is being run
type BuildContext struct {
	Cwd           string
	CwdIsRunnable bool
	Runnables     []RunnableDir
	Bundle        RunnableBundle
	Directive     *directive.Directive
}

// RunnableDir represents a directory containing a Runnable
type RunnableDir struct {
	Name           string
	UnderscoreName string
	Fullpath       string
	DotHive        DotHive
	BuildImage     string
}

// RunnableBundle contains information about a bundle in the current context
type RunnableBundle struct {
	Exists   bool
	Fullpath string
}

// DotHive represents a .hive.yanl file
type DotHive struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Lang      string `yaml:"lang"`
}

// CurrentBuildContext returns the build context for the provided working directory
func CurrentBuildContext(cwd string) (*BuildContext, error) {
	runnables, cwdIsRunnable, err := getRunnableDirs(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to getRunnableDirs")
	}

	bundle, err := bundleTargetPath(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to bundleIfExists")
	}

	directive, err := readDirectiveFile(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to readDirectiveFile")
	}

	bctx := &BuildContext{
		Cwd:           cwd,
		CwdIsRunnable: cwdIsRunnable,
		Runnables:     runnables,
		Bundle:        *bundle,
		Directive:     directive,
	}

	return bctx, nil
}

// RunnableExists returns true if the context contains a runnable with name <name>
func (b *BuildContext) RunnableExists(name string) bool {
	for _, r := range b.Runnables {
		if r.Name == name {
			return true
		}
	}

	return false
}

func getRunnableDirs(cwd string) ([]RunnableDir, bool, error) {
	runnables := []RunnableDir{}

	// go through all of the dirs in the current dir
	topLvlFiles, err := ioutil.ReadDir(cwd)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to list directory")
	}

	// check to see if we're running from within a Runnable directory
	// and return true if so.
	runnableDir, err := getRunnableFromFiles(cwd, topLvlFiles)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to getRunnableFromFiles")
	} else if runnableDir != nil {
		runnables = append(runnables, *runnableDir)
		return runnables, true, nil
	}

	for _, tf := range topLvlFiles {
		if !tf.IsDir() {
			continue
		}

		dirPath := filepath.Join(cwd, tf.Name())

		// determine if a .hive.yaml exists in that dir
		innerFiles, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return nil, false, errors.Wrapf(err, "failed to list files in %s", tf.Name())
		}

		runnableDir, err := getRunnableFromFiles(dirPath, innerFiles)
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to getRunnableFromFiles")
		} else if runnableDir == nil {
			continue
		}

		runnables = append(runnables, *runnableDir)
	}

	return runnables, false, nil
}

func containsDotHiveYaml(files []os.FileInfo) (string, bool) {
	for _, f := range files {
		if f.Name() == ".hive.yaml" {
			return ".hive.yaml", true
		} else if f.Name() == ".hive.yml" {
			return ".hive.yml", true
		}
	}

	return "", false
}

func getRunnableFromFiles(wd string, files []os.FileInfo) (*RunnableDir, error) {
	filename, exists := containsDotHiveYaml(files)
	if !exists {
		return nil, nil
	}

	dotHiveBytes, err := ioutil.ReadFile(filepath.Join(wd, filename))
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadFile .hive yaml")
	}

	dotHive := DotHive{}
	if err := yaml.Unmarshal(dotHiveBytes, &dotHive); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal .hive yaml")
	}

	if dotHive.Name == "" {
		dotHive.Name = filepath.Base(wd)
	}

	if dotHive.Namespace == "" {
		dotHive.Namespace = "default"
	}

	img := imageForLang(dotHive.Lang)
	if img == "" {
		return nil, fmt.Errorf("(%s) %s is not a valid lang", dotHive.Name, dotHive.Lang)
	}

	absolutePath, err := filepath.Abs(wd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Abs filepath")
	}

	runnable := &RunnableDir{
		Name:           filepath.Base(wd),
		UnderscoreName: strings.Replace(filepath.Base(wd), "-", "_", -1),
		Fullpath:       absolutePath,
		DotHive:        dotHive,
		BuildImage:     img,
	}

	return runnable, nil
}

func imageForLang(lang string) string {
	img, ok := dockerImageForLang[lang]
	if !ok {
		return ""
	}

	return img
}

func bundleTargetPath(cwd string) (*RunnableBundle, error) {
	path := filepath.Join(cwd, "runnables.wasm.zip")

	b := &RunnableBundle{
		Fullpath: path,
		Exists:   false,
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// do nothing
		} else {
			return nil, err
		}
	}

	b.Exists = true

	return b, nil
}
