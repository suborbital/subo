package main

import (
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/command"
	"github.com/suborbital/subo/subo/features"
	"github.com/suborbital/subo/subo/release"
)

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subo",
		Short:   "Suborbital Development Platform CLI",
		Version: release.SuboDotVersion,
		Long: `
Subo is the full toolchain for using and managing Suborbital Development Platform tools,
including building WebAssembly Runnables and Atmo projects.

Explore the available commands by running 'subo --help'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}
			return nil
		},
	}

	cmd.SetVersionTemplate("Subo CLI v{{.Version}}\n")

	// create commands
	create := &cobra.Command{
		Use:   "create",
		Short: "create a runnable or project",
		Long:  `create a new Atmo project or WebAssembly runnable`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}
			return nil
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}
			return nil
		},
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create compute network resources",
		Long:  `create Suborbital Compute Network resources`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}
			return nil
		},
	}

	create.AddCommand(command.ComputeCreateTokenCommand())
	compute.AddCommand(create)

	deploy := &cobra.Command{
		Use:   "deploy",
		Short: "deploy compute network resources",
		Long:  `deploy Suborbital Compute Network resources`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}
			return nil
		},
	}

	deploy.AddCommand(command.ComputeDeployCoreCommand())
	compute.AddCommand(deploy)

	return compute
}
