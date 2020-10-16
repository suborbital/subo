package context

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var dockerImageForLang = map[string]string{
	"rust": "suborbital/hive-rs:1.46.0-3",
}

// BuildContext describes the context under which the tool is being run
type BuildContext struct {
	Cwd       string
	Runnables []RunnableDir
	Bundle    RunnableBundle
}

// RunnableDir represents a directory containing a Runnable
type RunnableDir struct {
	Name       string
	Fullpath   string
	DotHive    DotHive
	BuildImage string
}

// RunnableBundle contains information about a bundle in the current context
type RunnableBundle struct {
	Exists   bool
	Fullpath string
}

// DotHive represents a .hive.yanl file
type DotHive struct {
	Lang string `yaml:"lang"`
}

// CurrentBuildContext returns the current build context
func CurrentBuildContext() (*BuildContext, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CWD")
	}

	runnables, err := getRunnableDirs(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to getRunnableDirs")
	}

	bundle, err := bundleTargetPath(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to bundleIfExists")
	}

	bctx := &BuildContext{
		Cwd:       cwd,
		Runnables: runnables,
		Bundle:    *bundle,
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

func getRunnableDirs(cwd string) ([]RunnableDir, error) {
	runnables := []RunnableDir{}

	// go through all of the dirs in the current dir
	topLvlFiles, err := ioutil.ReadDir(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list directory")
	}

	for _, tf := range topLvlFiles {
		if !tf.IsDir() {
			continue
		}

		// determine if a .hive.yaml exists in that dir
		innerFiles, err := ioutil.ReadDir(filepath.Join(cwd, tf.Name()))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list files in %s", tf.Name())
		}

		if containsDotHiveYaml(innerFiles) {
			dotHive := DotHive{
				Lang: "rust",
			}

			img := imageForLang(dotHive.Lang)

			runnable := RunnableDir{
				Name:       tf.Name(),
				Fullpath:   filepath.Join(cwd, tf.Name()),
				DotHive:    dotHive,
				BuildImage: img,
			}

			runnables = append(runnables, runnable)
		}
	}

	return runnables, nil
}

func containsDotHiveYaml(files []os.FileInfo) bool {
	for _, f := range files {
		if f.Name() == ".hive.yaml" || f.Name() == ".hive.yml" {
			return true
		}
	}

	return false
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
