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
	modules []os.File
}

type Toolchain string

const (
	ToolchainNative = Toolchain("native")
	ToolchainDocker = Toolchain("docker")
)

// ForDirectory creates a Builder bound to a particular directory
func ForDirectory(dir string) (*Builder, error) {
	ctx, err := context.ForDirectory(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to context.FirDirectory")
	}

	b := &Builder{
		Context: ctx,
	}

	return b, nil
}

func (b *Builder) BuildWithToolchain(tcn Toolchain) error {
	var err error

	modules := make([]os.File, len(b.Context.Runnables))

	for i, r := range b.Context.Runnables {
		util.LogStart(fmt.Sprintf("building runnable: %s (%s)", r.Name, r.Runnable.Lang))

		var file *os.File

		if tcn == ToolchainNative {
			if err := checkAndRunPreReqs(r); err != nil {
				return errors.Wrap(err, "🚫 failed to checkAndRunPreReqs")
			}

			file, err = doNativeBuildForRunnable(r)
		} else {
			file, err = doBuildForRunnable(r)
		}

		if err != nil {
			return errors.Wrapf(err, "🚫 failed to build %s", r.Name)
		}

		modules[i] = *file

		fullWasmFilepath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))
		util.LogDone(fmt.Sprintf("%s was built -> %s", r.Name, fullWasmFilepath))
	}

	b.modules = modules

	return nil
}

func (b *Builder) Bundle() error {
	if b.modules == nil || len(b.modules) == 0 {
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
		util.LogInfo("updating Directive")

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
		util.LogInfo("adding static files to bundle")
	}

	directiveBytes, err := b.Context.Directive.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to Directive.Marshal")
	}

	if err := bundle.Write(directiveBytes, b.modules, static, b.Context.Bundle.Fullpath); err != nil {
		return errors.Wrap(err, "🚫 failed to WriteBundle")
	}

	return nil
}

func doBuildForRunnable(r context.RunnableDir) (*os.File, error) {
	img := r.BuildImage
	if img == "" {
		return nil, fmt.Errorf("%q is not a supported language", r.Runnable.Lang)
	}

	_, _, err := util.Run(fmt.Sprintf("docker run --rm --mount type=bind,source=%s,target=/root/runnable %s", r.Fullpath, img))
	if err != nil {
		return nil, errors.Wrap(err, "failed to Run docker command")
	}

	targetPath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))

	file, err := os.Open(targetPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open resulting built file %s", targetPath)
	}

	return file, nil
}

func doNativeBuildForRunnable(r context.RunnableDir) (*os.File, error) {
	cmds, err := context.NativeBuildCommands(r.Runnable.Lang)
	if err != nil {
		return nil, errors.Wrap(err, "failed to NativeBuildCommands")
	}

	for _, cmd := range cmds {
		cmdTmpl, err := template.New("cmd").Parse(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to Parse command template")
		}

		fullCmd := &strings.Builder{}
		if err := cmdTmpl.Execute(fullCmd, r); err != nil {
			return nil, errors.Wrap(err, "failed to Execute command template")
		}

		_, _, err = util.RunInDir(fullCmd.String(), r.Fullpath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to RunInDir")
		}
	}

	targetPath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))

	file, err := os.Open(targetPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open resulting built file %s", targetPath)
	}

	return file, nil
}

func checkAndRunPreReqs(runnable context.RunnableDir) error {
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
				util.LogStart(fmt.Sprintf("missing %s, fixing...", p.File))

				_, _, err := util.RunInDir(p.Command, runnable.Fullpath)
				if err != nil {
					return errors.Wrapf(err, "failed to Run prerequisite: %s", p.Command)
				}

				util.LogDone("fixed!")
			}
		}
	}

	return nil
}
