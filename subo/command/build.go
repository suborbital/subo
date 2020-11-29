package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/hive-wasm/directive"
	"github.com/suborbital/hive-wasm/wasm"
	"github.com/suborbital/subo/subo/context"
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
				return errors.New("🚫 no runnables found in current directory (no .hive yaml files found)")
			}

			fmt.Println("✨ START: building runnables in", bctx.Cwd)

			shouldBundle, _ := cmd.Flags().GetBool("bundle")
			useNative, _ := cmd.Flags().GetBool("native")

			results := make([]os.File, len(bctx.Runnables))

			for i, r := range bctx.Runnables {
				fmt.Println(fmt.Sprintf("✨ START: building runnable: %s (%s)", r.Name, r.DotHive.Lang))

				var file *os.File

				if useNative {
					fmt.Println("🔗 using native toolchain")
					file, err = doNativeBuildForRunnable(r)
				} else {
					fmt.Println("🐳 using Docker toolchain")
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
					fmt.Println(fmt.Sprintf("✨ DONE: %s was built -> %s", r.Name, fullWasmFilepath))
				}

			}

			if shouldBundle {
				if bctx.Directive == nil {
					bctx.Directive = &directive.Directive{
						Identifier: "com.suborbital.app",
						// TODO: insert some git smarts here?
						Version: "v0.0.1",
					}
				}

				if err := context.AugmentAndValidateDirectiveFns(bctx.Directive, bctx.Runnables); err != nil {
					return errors.Wrap(err, "🚫 failed to AugmentAndValidateDirectiveFns")
				}

				if err := bctx.Directive.Validate(); err != nil {
					return errors.Wrap(err, "🚫 failed to Validate Directive")
				}

				if err := wasm.WriteBundle(bctx.Directive, results, bctx.Bundle.Fullpath); err != nil {
					return errors.Wrap(err, "🚫 failed to WriteBundle")
				}

				fmt.Println(fmt.Sprintf("✨ DONE: bundle was created -> %s", bctx.Bundle.Fullpath))
			}

			return nil
		},
	}

	cmd.Flags().Bool("bundle", false, "if passed, bundle all resulting runnables into a deployable .wasm.zip bundle")
	cmd.Flags().Bool("native", false, "if passed, build runnables using native toolchain rather than Docker")

	return cmd
}

func doBuildForRunnable(r context.RunnableDir) (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CWD")
	}

	img := r.BuildImage
	if img == "" {
		return nil, fmt.Errorf("%q is not a supported language", r.DotHive.Lang)
	}

	_, _, err = util.Run(fmt.Sprintf("docker run --rm --mount type=bind,source=%s,target=/root/rs-wasm %s", r.Fullpath, img))
	if err != nil {
		return nil, errors.Wrap(err, "failed to Run docker command")
	}

	targetPath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))
	os.Rename(filepath.Join(cwd, r.Name, "wasm_runner_bg.wasm"), targetPath)

	file, err := os.Open(targetPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open resulting built file %s", targetPath)
	}

	return file, nil
}

func doNativeBuildForRunnable(r context.RunnableDir) (*os.File, error) {
	cmds, err := context.NativeBuildCommands(r.DotHive.Lang)
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
