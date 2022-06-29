package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/suborbital/velo/builder"
	"github.com/suborbital/velo/cli/util"
	"github.com/suborbital/velo/packager"
)

// BuildCmd returns the build command.
func BuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build a project or single function",
		Long:  `build a build a Velocity project and bundle the results, or build a single function into a WebAssembly module`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")

			bdr, err := builder.ForDirectory(&util.PrintLogger{}, &builder.DefaultBuildConfig, dir)
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
			runCmd, _ := cmd.Flags().GetString("run")

			// Determine if a custom Docker mountpath and relpath were set.
			mountPath, _ := cmd.Flags().GetString("mountpath")
			relPath, _ := cmd.Flags().GetString("relpath")

			if mountPath != "" {
				if relPath == "" {
					// Fallback to the dir arg as that's usually a sane default.
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
				_, err = util.Command.Run(fmt.Sprintf("make %s", makeTarget))
				if err != nil {
					return errors.Wrapf(err, "🚫 failed to make %s", makeTarget)
				}
			}

			if runCmd != "" {
				util.LogStart(fmt.Sprintf("running %s", runCmd))
				_, err = util.Command.Run(runCmd)
				if err != nil {
					return errors.Wrapf(err, "🚫 failed to %s", runCmd)
				}
			}

			var toolchain builder.Toolchain
			if useNative {
				toolchain = builder.ToolchainNative
			} else {
				util.LogInfo("🐳 using Docker toolchain")
				toolchain = builder.ToolchainDocker
			}

			// The builder does the majority of the work.
			if err := bdr.BuildWithToolchain(toolchain); err != nil {
				return errors.Wrap(err, "failed to BuildWithToolchain")
			}

			pkgr := packager.New(&util.PrintLogger{})
			pkgJobs := []packager.PackageJob{}

			if shouldBundle {
				pkgJobs = append(pkgJobs, packager.NewBundlePackageJob())
			}

			if shouldDockerBuild && !bdr.Context.CwdIsRunnable {
				pkgJobs = append(pkgJobs, packager.NewDockerImagePackageJob())
			}

			if err := pkgr.Package(bdr.Context, pkgJobs...); err != nil {
				return errors.Wrap(err, "failed to Package")
			}

			return nil
		},
	}

	cmd.Flags().String("dir", ".", "the directory in which to run the build")
	cmd.Flags().Bool("no-bundle", false, "if passed, a .wasm.zip bundle will not be generated")
	cmd.Flags().Bool("native", false, "use native (locally installed) toolchain rather than Docker")
	cmd.Flags().String("make", "", "execute the provided Make target before building the project bundle")
	cmd.Flags().String(runPartnerFlag, "", "execute the provided command before building the project bundle")
	cmd.Flags().Bool("docker", false, "build your project's Dockerfile. It will be tagged {identifier}:{appVersion}")
	cmd.Flags().StringSlice("langs", []string{}, "build only Runnables for the listed languages (comma-seperated)")
	cmd.Flags().String("mountpath", "", "if passed, the Docker builders will mount their volumes at the provided path")
	cmd.Flags().String("relpath", "", "if passed, the Docker builders will run `subo build` using the provided path, relative to '--mountpath'")
	cmd.Flags().String("builder-tag", "", "use the provided tag for builder images")

	return cmd
}
