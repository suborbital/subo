package main

import (
	"github.com/spf13/cobra"

	"github.com/suborbital/velo/cli/command"
	"github.com/suborbital/velo/cli/features"
	"github.com/suborbital/velo/cli/release"
)

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "velo",
		Short:   "The Velocity CLI",
		Version: release.Version(),
		Long:    `Velocity is a function server that adds a backend to any application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmd.SetVersionTemplate("Velo v{{.Version}}\n")

	// velo init
	cmd.AddCommand(command.InitCmd())
	// velo build
	cmd.AddCommand(command.BuildCmd())
	// velo dev
	cmd.AddCommand(command.DevCmd())
	// velo start
	cmd.AddCommand(command.StartCmd())
	// velo clean
	cmd.AddCommand(command.CleanCmd())
	// compute related commands.
	cmd.AddCommand(computeCommand())

	if features.EnableRegistryCommands {
		// velo push
		cmd.AddCommand(command.PushCmd())
		// velo deploy
		cmd.AddCommand(command.DeployCmd())
	}

	if features.EnableDocsCommands {
		// docs related commands.
		cmd.AddCommand(docsCommand())
	}

	// add commands.
	add := &cobra.Command{
		Use:   "add",
		Short: "add a function or handler",
		Long:  `add a new WebAssembly function or route handler`,
	}

	// velo add function
	add.AddCommand(command.CreateFunctionCmd())

	if features.EnableReleaseCommands {
		// velo add release
		add.AddCommand(command.CreateReleaseCmd())
	}

	cmd.AddCommand(add)

	cmd.CompletionOptions.HiddenDefaultCmd = true

	return cmd
}

func computeCommand() *cobra.Command {
	compute := &cobra.Command{
		Use:    "compute",
		Short:  "compute network related resources",
		Long:   `manage Suborbital Compute Network resources`,
		Hidden: true,
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
