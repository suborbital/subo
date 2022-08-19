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
		Version: release.Version(),
		Long: `Subo is the full toolchain for using and managing Suborbital Development Platform tools,
including building WebAssembly Runnables and Atmo projects.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmd.SetVersionTemplate("Subo CLI v{{.Version}}\n")

	// create commands.
	create := &cobra.Command{
		Use:   "create",
		Short: "create a runnable, project, or handler",
		Long:  `create a new Atmo project, WebAssembly runnable or handler`,
	}

	if features.EnableReleaseCommands {
		create.AddCommand(command.CreateReleaseCmd())
	}

	create.AddCommand(command.CreateProjectCmd())
	create.AddCommand(command.CreateRunnableCmd())
	// TODO: turn into create workflow command
	// Ref: https://github.com/suborbital/subo/issues/347
	// create.AddCommand(command.CreateHandlerCmd()).

	// compute network related commands.
	cmd.AddCommand(computeCommand())

	// docs related commands.
	cmd.AddCommand(docsCommand())

	cmd.AddCommand(create)
	cmd.AddCommand(command.BuildCmd())
	cmd.AddCommand(command.DevCmd())
	cmd.AddCommand(command.CleanCmd())

	if features.EnableRegistryCommands {
		cmd.AddCommand(command.PushCmd())
		cmd.AddCommand(command.DeployCmd())
	}

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

func docsCommand() *cobra.Command {
	docs := &cobra.Command{
		Use:   "docs",
		Short: "documentation generation resources",
		Long:  "test and generate code embedded markdown documentation",
	}
	docs.AddCommand(command.DocsBuildCmd())
	docs.AddCommand(command.DocsTestCmd())

	return docs
}
