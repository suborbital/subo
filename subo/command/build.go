package command

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/atmo/bundle"
	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
	"golang.org/x/mod/semver"
)

// BuildCmd returns the build command
func BuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [dir]",
		Short: "build a WebAssembly runnable",
		Long:  `build a WebAssembly runnable and/or create a Runable Bundle`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			bctx, err := context.CurrentBuildContext(dir)
			if err != nil {
				return errors.Wrap(err, "failed to get CurrentBuildContext")
			}

			if len(bctx.Runnables) == 0 {
				return errors.New("🚫 no runnables found in current directory (no .runnable yaml files found)")
			}

			if bctx.CwdIsRunnable {
				util.LogInfo("building single Runnable (run from project root to create bundle)")
			} else {
				util.LogStart(fmt.Sprintf("building runnables in %s", bctx.Cwd))
			}

			noBundle, _ := cmd.Flags().GetBool("no-bundle")
			shouldBundle := !noBundle && !bctx.CwdIsRunnable

			useNative, _ := cmd.Flags().GetBool("native")
			shouldDockerBuild, _ := cmd.Flags().GetBool("docker")

			modules := make([]os.File, len(bctx.Runnables))

			for i, r := range bctx.Runnables {
				util.LogStart(fmt.Sprintf("building runnable: %s (%s)", r.Name, r.Runnable.Lang))

				var file *os.File

				if useNative {
					util.LogInfo("🔗 using native toolchain")
					if err := checkAndRunPreReqs(r); err != nil {
						return errors.Wrap(err, "🚫 failed to checkAndRunPreReqs")
					}

					file, err = doNativeBuildForRunnable(r)
				} else {
					util.LogInfo("🐳 using Docker toolchain")
					file, err = doBuildForRunnable(r)
				}

				if err != nil {
					return errors.Wrapf(err, "🚫 failed to doBuild for %s", r.Name)
				}

				modules[i] = *file

				fullWasmFilepath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))
				util.LogDone(fmt.Sprintf("%s was built -> %s", r.Name, fullWasmFilepath))
			}

			if shouldBundle {
				if bctx.Directive == nil {
					bctx.Directive = &directive.Directive{
						Identifier: "com.suborbital.app",
						// TODO: insert some git smarts here?
						AppVersion:  "v0.0.1",
						AtmoVersion: fmt.Sprintf("v%s", release.AtmoVersion),
					}
				} else if bctx.Directive.Headless {
					util.LogInfo("updating Directive")

					// bump the appVersion since we're in headless mode
					majorStr := strings.TrimPrefix(semver.Major(bctx.Directive.AppVersion), "v")
					major, _ := strconv.Atoi(majorStr)
					new := fmt.Sprintf("v%d.0.0", major+1)

					bctx.Directive.AppVersion = new

					if err := context.WriteDirective(bctx.Cwd, bctx.Directive); err != nil {
						return errors.Wrap(err, "failed to WriteDirective")
					}
				}

				if err := context.AugmentAndValidateDirectiveFns(bctx.Directive, bctx.Runnables); err != nil {
					return errors.Wrap(err, "🚫 failed to AugmentAndValidateDirectiveFns")
				}

				if err := bctx.Directive.Validate(); err != nil {
					return errors.Wrap(err, "🚫 failed to Validate Directive")
				}

				static, err := context.CollectStaticFiles(bctx.Cwd)
				if err != nil {
					return errors.Wrap(err, "failed to CollectStaticFiles")
				}

				if static != nil {
					util.LogInfo("adding static files to bundle")
				}

				directiveBytes, err := bctx.Directive.Marshal()
				if err != nil {
					return errors.Wrap(err, "failed to Directive.Marshal")
				}

				if err := bundle.Write(directiveBytes, modules, static, bctx.Bundle.Fullpath); err != nil {
					return errors.Wrap(err, "🚫 failed to WriteBundle")
				}

				defer util.LogDone(fmt.Sprintf("bundle was created -> %s @ %s", bctx.Bundle.Fullpath, bctx.Directive.AppVersion))
			}

			if shouldDockerBuild {
				os.Setenv("DOCKER_BUILDKIT", "0")

				if _, _, err := util.Run(fmt.Sprintf("docker build . -t=%s:%s", bctx.Directive.Identifier, bctx.Directive.AppVersion)); err != nil {
					return errors.Wrap(err, "🚫 failed to build Docker image")
				}

				util.LogDone(fmt.Sprintf("built Docker image -> %s:%s", bctx.Directive.Identifier, bctx.Directive.AppVersion))
			}

			return nil
		},
	}

	cmd.Flags().Bool("no-bundle", false, "if passed, a .wasm.zip bundle will not be generated")
	cmd.Flags().Bool("native", false, "if passed, build runnables using native toolchain rather than Docker")
	cmd.Flags().Bool("docker", false, "if passed, build your project's Dockerfile. It will be tagged {identifier}:{appVersion}")

	return cmd
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
