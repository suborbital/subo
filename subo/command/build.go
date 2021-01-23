package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/hive-wasm/bundle"
	"github.com/suborbital/hive-wasm/directive"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

// BuildCmd returns the build command
func BuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [dir]",
		Short: "build a Wasm runnable",
		Long:  `build a Wasm runnable from local source files`,
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

			logStart(fmt.Sprintf("building runnables in %s", bctx.Cwd))

			shouldBundle, _ := cmd.Flags().GetBool("bundle")
			useNative, _ := cmd.Flags().GetBool("native")

			results := make([]os.File, len(bctx.Runnables))

			for i, r := range bctx.Runnables {
				logStart(fmt.Sprintf("building runnable: %s (%s)", r.Name, r.Runnable.Lang))

				var file *os.File

				if useNative {
					logInfo("🔗 using native toolchain")
					file, err = doNativeBuildForRunnable(r)
				} else {
					logInfo("🐳 using Docker toolchain")
					file, err = doBuildForRunnable(r)
				}

				if err != nil {
					buildErr := errors.Wrapf(err, "🚫 failed to doBuild for %s", r.Name)

					if shouldBundle {
						return buildErr
					}

					fmt.Println(buildErr)
				} else {
					results[i] = *file

					fullWasmFilepath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))
					logDone(fmt.Sprintf("%s was built -> %s", r.Name, fullWasmFilepath))
				}

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

				if err := bundle.Write(bctx.Directive, results, bctx.Bundle.Fullpath); err != nil {
					return errors.Wrap(err, "🚫 failed to WriteBundle")
				}

				logDone(fmt.Sprintf("bundle was created -> %s", bctx.Bundle.Fullpath))
			}

			return nil
		},
	}

	cmd.Flags().Bool("bundle", false, "if passed, bundle all resulting runnables into a deployable .wasm.zip bundle")
	cmd.Flags().Bool("native", false, "if passed, build runnables using native toolchain rather than Docker")

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

func logInfo(msg string) {
	if _, exists := os.LookupEnv("SUBO_DOCKER"); !exists {
		fmt.Println(msg)
	}
}

func logStart(msg string) {
	if _, exists := os.LookupEnv("SUBO_DOCKER"); !exists {
		fmt.Println(fmt.Sprintf("⏩ START: %s", msg))
	}
}

func logDone(msg string) {
	if _, exists := os.LookupEnv("SUBO_DOCKER"); !exists {
		fmt.Println(fmt.Sprintf("✅ DONE: %s", msg))
	}
}
