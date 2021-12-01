package main

import (
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/command"
	"github.com/suborbital/subo/subo/features"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

func rootCommand() *cobra.Command {
	defer func() {
		version_msg, err := release.CheckForLatestVersion()
		if err != nil {
			util.LogFail(err.Error())
		} else if version_msg != "" {
			util.LogInfo(version_msg)
		}
	}()

	cmd := &cobra.Command{
		Use:     "subo",
		Short:   "Suborbital Development Platform CLI",
		Version: release.SuboDotVersion,
		Long: `Subo is the full toolchain for using and managing Suborbital Development Platform tools,
including building WebAssembly Runnables and Atmo projects.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmd.SetVersionTemplate("Subo CLI v{{.Version}}\n")

	// create commands
	create := &cobra.Command{
		Use:   "create",
		Short: "create a runnable or project",
		Long:  `create a new Atmo project or WebAssembly runnable`,
	}

	if features.EnableReleaseCommands {
		create.AddCommand(command.CreateReleaseCmd())
	}

	create.AddCommand(command.CreateProjectCmd())
	create.AddCommand(command.CreateRunnableCmd())

	// compute network related commands
	cmd.AddCommand(computeCommand())

	// add top-level commands to root
	cmd.AddCommand(create)
	cmd.AddCommand(command.BuildCmd())
	cmd.AddCommand(command.DevCmd())
	cmd.AddCommand(command.CleanCmd())

	return cmd
}

func computeCommand() *cobra.Command {
	compute := &cobra.Command{
		Use:   "compute",
		Short: "compute network related resources",
		Long:  `manage Suborbital Compute Network resources`,
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create compute network resources",
		Long:  `create Suborbital Compute Network resources`,
	}

	create.AddCommand(command.ComputeCreateTokenCommand())
	compute.AddCommand(create)

	deploy := &cobra.Command{
		Use:   "deploy",
		Short: "deploy compute network resources",
		Long:  `deploy Suborbital Compute Network resources`,
	}

	deploy.AddCommand(command.ComputeDeployCoreCommand())
	compute.AddCommand(deploy)

	return compute
}
