package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/builder"
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

			bdr, err := builder.ForDirectory(dir)
			if err != nil {
				return errors.Wrap(err, "failed to builder.ForDirectory")
			}

			if len(bdr.Context.Runnables) == 0 {
				return errors.New("🚫 no runnables found in current directory (no .runnable yaml files found)")
			}

			if bdr.Context.CwdIsRunnable {
				util.LogInfo("building single Runnable (run from project root to create bundle)")
			} else {
				util.LogStart(fmt.Sprintf("building runnables in %s", bdr.Context.Cwd))
			}

			noBundle, _ := cmd.Flags().GetBool("no-bundle")
			shouldBundle := !noBundle && !bdr.Context.CwdIsRunnable
			shouldDockerBuild, _ := cmd.Flags().GetBool("docker")

			useNative, _ := cmd.Flags().GetBool("native")
			makeTarget, _ := cmd.Flags().GetString("make")

			if makeTarget != "" {
				util.LogStart(fmt.Sprintf("make %s", makeTarget))
				_, _, err = util.Run(fmt.Sprintf("make %s", makeTarget))
				if err != nil {
					return errors.Wrapf(err, "🚫 failed to make %s", makeTarget)
				}
			}

			var toolchain builder.Toolchain
			if useNative {
				util.LogInfo("🔗 using native toolchain")
				toolchain = builder.ToolchainNative
			} else {
				util.LogInfo("🐳 using Docker toolchain")
				toolchain = builder.ToolchainDocker
			}

			// the builder does the majority of the work:

			if err := bdr.BuildWithToolchain(toolchain); err != nil {
				return errors.Wrap(err, "failed to BuildWithToolchain")
			}

			if shouldBundle {
				if err := bdr.Bundle(); err != nil {
					return errors.Wrap(err, "failed to Bundle")
				}

				defer util.LogDone(fmt.Sprintf("bundle was created -> %s @ %s", bdr.Context.Bundle.Fullpath, bdr.Context.Directive.AppVersion))
			}

			if shouldDockerBuild {
				os.Setenv("DOCKER_BUILDKIT", "0")

				if _, _, err := util.Run(fmt.Sprintf("docker build . -t=%s:%s", bdr.Context.Directive.Identifier, bdr.Context.Directive.AppVersion)); err != nil {
					return errors.Wrap(err, "🚫 failed to build Docker image")
				}

				util.LogDone(fmt.Sprintf("built Docker image -> %s:%s", bdr.Context.Directive.Identifier, bdr.Context.Directive.AppVersion))
			}

			return nil
		},
	}

	cmd.Flags().Bool("no-bundle", false, "if passed, a .wasm.zip bundle will not be generated")
	cmd.Flags().Bool("native", false, "if passed, build runnables using native toolchain rather than Docker")
	cmd.Flags().String("make", "", "if passed, execute the provided Make target before building the project bundle")
	cmd.Flags().Bool("docker", false, "if passed, build your project's Dockerfile. It will be tagged {identifier}:{appVersion}")

	return cmd
}
