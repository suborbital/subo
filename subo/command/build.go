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

			bdr, err := builder.ForDirectory(&util.PrintLogger{}, dir)
			if err != nil {
				return errors.Wrap(err, "failed to builder.ForDirectory")
			}

			if len(bdr.Context.Runnables) == 0 {
				return errors.New("🚫 no runnables found in current directory (no .runnable.yaml files found)")
			}

			if bdr.Context.CwdIsRunnable {
				util.LogInfo("building single Runnable (run from project root to create bundle)")
			}

			langs, _ := cmd.Flags().GetStringSlice("langs")
			bdr.Context.Langs = langs

			noBundle, _ := cmd.Flags().GetBool("no-bundle")
			shouldBundle := !noBundle && !bdr.Context.CwdIsRunnable && len(langs) == 0
			shouldDockerBuild, _ := cmd.Flags().GetBool("docker")

			if bdr.Context.CwdIsRunnable && shouldDockerBuild {
				return errors.New("🚫 cannot build Docker image for a single Runnable (must be a project)")
			}

			useNative, _ := cmd.Flags().GetBool("native")
			makeTarget, _ := cmd.Flags().GetString("make")

			// determine if a custom Docker mountpath and relpath were set
			mountPath, _ := cmd.Flags().GetString("mountpath")
			relPath, _ := cmd.Flags().GetString("relpath")

			if mountPath != "" {
				if relPath == "" {
					// fallback to the dir arg as that's usually a sane default
					relPath = dir
				}

				bdr.Context.MountPath = mountPath
				bdr.Context.RelDockerPath = relPath
			}

			builderTag, _ := cmd.Flags().GetString("builder-tag")
			if builderTag != "" {
				bdr.Context.BuilderTag = builderTag
			}

			if makeTarget != "" {
				util.LogStart(fmt.Sprintf("make %s", makeTarget))
				_, err = util.Run(fmt.Sprintf("make %s", makeTarget))
				if err != nil {
					return errors.Wrapf(err, "🚫 failed to make %s", makeTarget)
				}
			}

			var toolchain builder.Toolchain
			if useNative {
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

			if shouldDockerBuild && !bdr.Context.CwdIsRunnable {
				if bdr.Context.Directive == nil {
					return errors.New("🚫 failed to open Directive.yaml")
				}

				os.Setenv("DOCKER_BUILDKIT", "0")

				if _, err := util.Run(fmt.Sprintf("docker build . -t=%s:%s", bdr.Context.Directive.Identifier, bdr.Context.Directive.AppVersion)); err != nil {
					return errors.Wrap(err, "🚫 failed to build Docker image")
				}

				util.LogDone(fmt.Sprintf("built Docker image -> %s:%s", bdr.Context.Directive.Identifier, bdr.Context.Directive.AppVersion))
			}

			return nil
		},
	}

	cmd.Flags().Bool("no-bundle", false, "if passed, a .wasm.zip bundle will not be generated")
	cmd.Flags().Bool("native", false, "use native (locally installed) toolchain rather than Docker")
	cmd.Flags().String("make", "", "execute the provided Make target before building the project bundle")
	cmd.Flags().Bool("docker", false, "build your project's Dockerfile. It will be tagged {identifier}:{appVersion}")
	cmd.Flags().StringSlice("langs", []string{}, "build only Runnables for the listed languages (comma-seperated)")
	cmd.Flags().String("mountpath", "", "if passed, the Docker builders will mount their volumes at the provided path")
	cmd.Flags().String("relpath", "", "if passed, the Docker builders will run `subo build` using the provided path, relative to '--mountpath'")
	cmd.Flags().String("builder-tag", "", "use the provided tag for builder images")

	return cmd
}
