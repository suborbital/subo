package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/velo/cli/release"
	"github.com/suborbital/velo/cli/util"
)

// validLangs are the available languages.
var validLangs = map[string]struct{}{
	"rust":           {},
	"swift":          {},
	"assemblyscript": {},
	"tinygo":         {},
	"grain":          {},
	"typescript":     {},
	"javascript":     {},
	"wat":            {},
}

// Context describes the context under which the tool is being run.
type Context struct {
	Cwd           string
	CwdIsRunnable bool
	Runnables     []RunnableDir
	Bundle        BundleRef
	Directive     *directive.Directive
	AtmoVersion   string
	Langs         []string
	MountPath     string
	RelDockerPath string
	BuilderTag    string
}

// RunnableDir represents a directory containing a Runnable.
type RunnableDir struct {
	Name           string
	UnderscoreName string
	Fullpath       string
	Runnable       *directive.Runnable
	CompilerFlags  string
}

// BundleRef contains information about a bundle in the current context.
type BundleRef struct {
	Exists   bool
	Fullpath string
}

// ForDirectory returns the build context for the provided working directory.
func ForDirectory(dir string) (*Context, error) {
	fullDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Abs path")
	}

	runnables, cwdIsRunnable, err := getRunnableDirs(fullDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to getRunnableDirs")
	}

	bundle, err := bundleTargetPath(fullDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to bundleIfExists")
	}

	directive, err := readDirectiveFile(fullDir)
	if err != nil {
		if !os.IsNotExist(errors.Cause(err)) {
			return nil, errors.Wrap(err, "failed to readDirectiveFile")
		}
	}

	queries, err := readQueriesFile(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to readQueriesFile")
	} else if len(queries) > 0 {
		directive.Queries = queries
	}

	bctx := &Context{
		Cwd:           fullDir,
		CwdIsRunnable: cwdIsRunnable,
		Runnables:     runnables,
		Bundle:        *bundle,
		Directive:     directive,
		Langs:         []string{},
		MountPath:     fullDir,
		RelDockerPath: ".",
		BuilderTag:    fmt.Sprintf("v%s", release.SuboDotVersion),
	}

	if directive != nil {
		bctx.AtmoVersion = directive.AtmoVersion
	}

	return bctx, nil
}

// RunnableExists returns true if the context contains a runnable with name <name>.
func (b *Context) RunnableExists(name string) bool {
	for _, r := range b.Runnables {
		if r.Name == name {
			return true
		}
	}

	return false
}

// ShouldBuildLang returns true if the provided language is safe-listed for building.
func (b *Context) ShouldBuildLang(lang string) bool {
	if len(b.Langs) == 0 {
		return true
	}

	for _, l := range b.Langs {
		if l == lang {
			return true
		}
	}

	return false
}

func (b *Context) Modules() ([]os.File, error) {
	modules := []os.File{}

	for _, r := range b.Runnables {
		wasmPath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))

		file, err := os.Open(wasmPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to Open module file %s", wasmPath)
		}

		modules = append(modules, *file)
	}

	return modules, nil
}

// HasDockerfile returns a nil error if the project's Dockerfile exists.
func (b *Context) HasDockerfile() error {
	dockerfilePath := filepath.Join(b.Cwd, "Dockerfile")

	if _, err := os.Stat(dockerfilePath); err != nil {
		return errors.Wrap(err, "failed to Stat Dockerfile")
	}

	return nil
}

// HasModule returns a nil error if the Runnable's .wasm file exists.
func (r *RunnableDir) HasModule() error {
	runnablePath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))

	if _, err := os.Stat(runnablePath); err != nil {
		return errors.Wrapf(err, "failed to Stat %s", runnablePath)
	}

	return nil
}

func getRunnableDirs(cwd string) ([]RunnableDir, bool, error) {
	runnables := []RunnableDir{}

	// Go through all of the dirs in the current dir.
	topLvlFiles, err := ioutil.ReadDir(cwd)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to list directory")
	}

	// Check to see if we're running from within a Runnable directory
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

		// Determine if a .runnable file exists in that dir.
		innerFiles, err := ioutil.ReadDir(dirPath)
		if err != nil {
			util.LogWarn(fmt.Sprintf("couldn't read files in %v", dirPath))
			continue
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

// ContainsRunnableYaml finds any .runnable file in a list of files.
func ContainsRunnableYaml(files []os.FileInfo) (string, bool) {
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".runnable.") {
			return f.Name(), true
		}
	}

	return "", false
}

// IsValidLang returns true if a language is valid.
func IsValidLang(lang string) bool {
	_, exists := validLangs[lang]

	return exists
}

func getRunnableFromFiles(wd string, files []os.FileInfo) (*RunnableDir, error) {
	filename, exists := ContainsRunnableYaml(files)
	if !exists {
		return nil, nil
	}

	runnableBytes, err := ioutil.ReadFile(filepath.Join(wd, filename))
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadFile .runnable yaml")
	}

	runnable := directive.Runnable{}
	if err := yaml.Unmarshal(runnableBytes, &runnable); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal .runnable yaml")
	}

	if runnable.Name == "" {
		runnable.Name = filepath.Base(wd)
	}

	if runnable.Namespace == "" {
		runnable.Namespace = "default"
	}

	if ok := IsValidLang(runnable.Lang); !ok {
		return nil, fmt.Errorf("(%s) %s is not a valid lang", runnable.Name, runnable.Lang)
	}

	absolutePath, err := filepath.Abs(wd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Abs filepath")
	}

	runnableDir := &RunnableDir{
		Name:           runnable.Name,
		UnderscoreName: strings.Replace(runnable.Name, "-", "_", -1),
		Fullpath:       absolutePath,
		Runnable:       &runnable,
	}

	return runnableDir, nil
}

func bundleTargetPath(cwd string) (*BundleRef, error) {
	path := filepath.Join(cwd, "runnables.wasm.zip")

	b := &BundleRef{
		Fullpath: path,
		Exists:   false,
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return b, nil
		} else {
			return nil, err
		}
	}

	b.Exists = true

	return b, nil
}
