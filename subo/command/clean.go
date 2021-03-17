package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/util"
)

//CleanCmd  removes all of the target/.build folders for Runnables and deletes the .wasm files.
func CleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "remove build folders and .wasm files",
		Long:  "remove all of target/.build folders and deletes .wasm files",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				cwd = "$HOME"
			}

			bctx, err := context.CurrentBuildContext(cwd)
			if err != nil {
				return errors.Wrap(err, "failed to get CurrentBuildContext")
			}

			if len(bctx.Runnables) == 0 {
				return errors.New("🚫 no runnables found in current directory (no .runnable yaml files found)")
			}
			logStart(fmt.Sprintf("cleaning in %s", bctx.Cwd))

			//remove  target folder

			if outStr, _, _ := util.Run(fmt.Sprintf("find . -type d -name target")); outStr != "" {
				if _, _, err := util.Run(fmt.Sprintf("rm -r ./target")); err != nil {
					return errors.Wrap(err, "🚫 failed to remove target folder")
				}
				logDone("removed target folder")
			}

			//remove  .build folder
			if outStr, _, _ := util.Run(fmt.Sprintf("find . -type d -name .build")); outStr != "" {
				if _, _, err := util.Run(fmt.Sprintf("rm -r ./.build")); err != nil {
					return errors.Wrap(err, "🚫 failed to remove target folder")
				}
				logDone("removed .build folder")
			}

			//find all .wasm files
			if outStr, _, _ := util.Run(fmt.Sprintf("find . -type f -name *.wasm")); outStr != "" {
				if _, _, err := util.Run(fmt.Sprintf("find . -type f -name *.wasm -delete")); err != nil {
					return errors.Wrap(err, "🚫 failed to delete .wasm files")
				}
				logDone("removed all .wasm")
			}

			logDone("cleaned")
			return nil
		},
	}

	return cmd
}
