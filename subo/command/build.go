package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/hive-wasm/wasm"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/util"
)

// BuildCmd returns the build command
func BuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build a Wasm runnable",
		Long:  `build a Wasm runnable from local source files`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			bctx, err := context.CurrentBuildContext(dir)
			if err != nil {
				return errors.Wrap(err, "failed to get CurrentBuildContext")
			}

			if len(bctx.Runnables) == 0 {
				return errors.New("🚫 no runnables found in current directory (no .hive yaml files found)")
			}

			fmt.Println("✨ START: building runnables in", bctx.Cwd)

			results := make([]os.File, len(bctx.Runnables))

			for i, r := range bctx.Runnables {
				fmt.Println(fmt.Sprintf("✨ START: building runnable: %s (%s)", r.Name, r.DotHive.Lang))

				file, err := doBuildForRunnable(r)
				if err != nil {
					fmt.Println("🚫 FAIL:", errors.Wrapf(err, "failed to doBuild for %s", r.Name))
				} else {
					results[i] = *file
					fullWasmFilepath := filepath.Join(r.Fullpath, fmt.Sprintf("%s.wasm", r.Name))
					fmt.Println(fmt.Sprintf("✨ DONE: %s was built -> %s", r.Name, fullWasmFilepath))
				}

			}

			shouldBundle, err := cmd.Flags().GetBool("bundle")
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get bundle flag")
			} else if shouldBundle {
				if err := wasm.WriteBundle(results, bctx.Bundle.Fullpath); err != nil {
					return errors.Wrap(err, "🚫 failed to WriteBundle")
				}

				fmt.Println(fmt.Sprintf("✨ DONE: bundle was created -> %s", bctx.Bundle.Fullpath))
			}

			return nil
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "$HOME"
	}

	cmd.Flags().String("dir", cwd, "the directory to run the build from")
	cmd.Flags().Bool("bundle", false, "if true, bundle all resulting runnables into a deployable .wasm.zip bundle")

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
