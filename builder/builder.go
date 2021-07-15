package builder

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/bundle"
	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/subo/builder/context"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
	"golang.org/x/mod/semver"
)

// Builder is capable of building Wasm modules from source
type Builder struct {
	Context *context.BuildContext

	results []BuildResult

	log util.FriendlyLogger
}

// BuildResult is the results of a build including the built module and logs
type BuildResult struct {
	Succeeded bool
	OutputLog string
	Module    *os.File
}

type Toolchain string

const (
	ToolchainNative = Toolchain("native")
	ToolchainDocker = Toolchain("docker")
)

// ForDirectory creates a Builder bound to a particular directory
func ForDirectory(logger util.FriendlyLogger, dir string) (*Builder, error) {
	ctx, err := context.ForDirectory(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to context.FirDirectory")
	}

	b := &Builder{
		Context: ctx,
		results: []BuildResult{},
		log:     logger,
	}

	return b, nil
}

func (b *Builder) BuildWithToolchain(tcn Toolchain) error {
	var err error

	b.results = make([]BuildResult, len(b.Context.Runnables))

	for i, r := range b.Context.Runnables {
		b.log.LogStart(fmt.Sprintf("building runnable: %s (%s)", r.Name, r.Runnable.Lang))

		result := &BuildResult{}

		if tcn == ToolchainNative {
			if err := b.checkAndRunPreReqs(r, result); err != nil {
				return errors.Wrap(err, "🚫 failed to checkAndRunPreReqs")
			}

			err = b.doNativeBuildForRunnable(r, result)
		} else {
			err = b.doBuildForRunnable(r, result)
		}

		// even if there was a failure, load the result into the builder
		// since the logs of the failed build are useful
		b.results[i] = *result

		if err != nil {
			return errors.Wrapf(err, "🚫 failed to build %s", r.Name)
		}

		fullWasmFilepath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))
		b.log.LogDone(fmt.Sprintf("%s was built -> %s", r.Name, fullWasmFilepath))
	}

	return nil
}

// Results returns build results for all of the modules built by this builder
// returns os.ErrNotExists if none have been built yet.
func (b *Builder) Results() ([]BuildResult, error) {
	if b.results == nil || len(b.results) == 0 {
		return nil, os.ErrNotExist
	}

	return b.results, nil
}

func (b *Builder) Bundle() error {
	if b.results == nil || len(b.results) == 0 {
		return errors.New("must build before calling Bundle")
	}

	if b.Context.Directive == nil {
		b.Context.Directive = &directive.Directive{
			Identifier: "com.suborbital.app",
			// TODO: insert some git smarts here?
			AppVersion:  "v0.0.1",
			AtmoVersion: fmt.Sprintf("v%s", release.AtmoVersion),
		}
	} else if b.Context.Directive.Headless {
		b.log.LogInfo("updating Directive")

		// bump the appVersion since we're in headless mode
		majorStr := strings.TrimPrefix(semver.Major(b.Context.Directive.AppVersion), "v")
		major, _ := strconv.Atoi(majorStr)
		new := fmt.Sprintf("v%d.0.0", major+1)

		b.Context.Directive.AppVersion = new

		if err := context.WriteDirectiveFile(b.Context.Cwd, b.Context.Directive); err != nil {
			return errors.Wrap(err, "failed to WriteDirectiveFile")
		}
	}

	if err := context.AugmentAndValidateDirectiveFns(b.Context.Directive, b.Context.Runnables); err != nil {
		return errors.Wrap(err, "🚫 failed to AugmentAndValidateDirectiveFns")
	}

	if err := b.Context.Directive.Validate(); err != nil {
		return errors.Wrap(err, "🚫 failed to Validate Directive")
	}

	static, err := context.CollectStaticFiles(b.Context.Cwd)
	if err != nil {
		return errors.Wrap(err, "failed to CollectStaticFiles")
	}

	if static != nil {
		b.log.LogInfo("adding static files to bundle")
	}

	directiveBytes, err := b.Context.Directive.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to Directive.Marshal")
	}

	modules := make([]os.File, len(b.results))
	for i := range b.results {
		modules[i] = *b.results[i].Module
	}

	if err := bundle.Write(directiveBytes, modules, static, b.Context.Bundle.Fullpath); err != nil {
		return errors.Wrap(err, "🚫 failed to WriteBundle")
	}

	return nil
}

func (b *Builder) doBuildForRunnable(r context.RunnableDir, result *BuildResult) error {
	img := r.BuildImage
	if img == "" {
		return fmt.Errorf("%q is not a supported language", r.Runnable.Lang)
	}

	outputLog, err := util.Run(fmt.Sprintf("docker run --rm --mount type=bind,source=%s,target=/root/runnable %s", r.Fullpath, img))

	result.OutputLog = outputLog

	if err != nil {
		result.Succeeded = false
		return errors.Wrap(err, "failed to Run docker command")
	}

	result.Succeeded = true

	targetPath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))

	file, err := os.Open(targetPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open resulting built file %s", targetPath)
	}

	result.Module = file

	return nil
}

// results and resulting file are loaded into the BuildResult pointer
func (b *Builder) doNativeBuildForRunnable(r context.RunnableDir, result *BuildResult) error {
	cmds, err := context.NativeBuildCommands(r.Runnable.Lang)
	if err != nil {
		return errors.Wrap(err, "failed to NativeBuildCommands")
	}

	for _, cmd := range cmds {
		cmdTmpl, err := template.New("cmd").Parse(cmd)
		if err != nil {
			return errors.Wrap(err, "failed to Parse command template")
		}

		fullCmd := &strings.Builder{}
		if err := cmdTmpl.Execute(fullCmd, r); err != nil {
			return errors.Wrap(err, "failed to Execute command template")
		}

		// Even if the command fails, still load the output into the result object
		outputLog, err := util.RunInDir(fullCmd.String(), r.Fullpath)

		result.OutputLog += outputLog + "\n"

		if err != nil {
			result.Succeeded = false
			return errors.Wrap(err, "failed to RunInDir")
		}

		result.Succeeded = true
	}

	targetPath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))

	file, err := os.Open(targetPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open resulting built file %s", targetPath)
	}

	result.Module = file

	return nil
}

func (b *Builder) checkAndRunPreReqs(runnable context.RunnableDir, result *BuildResult) error {
	preReqLangs, ok := context.PreRequisiteCommands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	preReqs, ok := preReqLangs[runnable.Runnable.Lang]
	if !ok {
		return fmt.Errorf("unsupported language: %s", runnable.Runnable.Lang)
	}

	for _, p := range preReqs {
		filepath := filepath.Join(runnable.Fullpath, p.File)

		if _, err := os.Stat(filepath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				b.log.LogStart(fmt.Sprintf("missing %s, fixing...", p.File))

				outputLog, err := util.RunInDir(p.Command, runnable.Fullpath)

				result.OutputLog += outputLog + "\n"

				if err != nil {
					return errors.Wrapf(err, "failed to Run prerequisite: %s", p.Command)
				}

				b.log.LogDone("fixed!")
			}
		}
	}

	return nil
}
