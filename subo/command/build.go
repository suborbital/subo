package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/reactr/bundle"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
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
				util.LogInfo("ℹ️  building single Runnable (run from project root to create bundle)")
			} else {
				util.LogStart(fmt.Sprintf("building runnables in %s", bctx.Cwd))
			}

			noBundle, _ := cmd.Flags().GetBool("no-bundle")
			shouldBundle := !noBundle

			useNative, _ := cmd.Flags().GetBool("native")
			shouldDockerBuild, _ := cmd.Flags().GetBool("docker")

			modules := make([]os.File, len(bctx.Runnables))

			for i, r := range bctx.Runnables {
				util.LogStart(fmt.Sprintf("building runnable: %s (%s)", r.Name, r.Runnable.Lang))

				var file *os.File

				if useNative {
					util.LogInfo("🔗 using native toolchain")
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
					util.LogInfo("ℹ️  adding static files to bundle")
				}

				directiveBytes, err := bctx.Directive.Marshal()
				if err != nil {
					return errors.Wrap(err, "failed to Directive.Marshal")
				}

				if err := bundle.Write(directiveBytes, modules, static, bctx.Bundle.Fullpath); err != nil {
					return errors.Wrap(err, "🚫 failed to WriteBundle")
				}

				defer util.LogDone(fmt.Sprintf("bundle was created -> %s", bctx.Bundle.Fullpath))
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
	cmd.Flags().Bool("docker", false, "pass --docker to automatically build a Docker image based on your project's Dockerfile. It will be tagged with the 'identifier' and 'appVersion' from your Directive")

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
