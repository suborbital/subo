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
		Long:    `Subo is the full toolchain for using and managing Suborbital Extension Engine (SE2) tools.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmd.SetVersionTemplate("Subo CLI v{{.Version}}\n")

	// create commands.
	create := &cobra.Command{
		Use:   "create",
		Short: "create a plugin, project, or handler",
		Long:  `create a new E2Core project, WebAssembly plugin or handler`,
	}

	if features.EnableReleaseCommands {
		create.AddCommand(command.CreateReleaseCmd())
	}

	create.AddCommand(command.CreateProjectCmd())
	create.AddCommand(command.CreatePluginCmd())
	// TODO: turn into create workflow command
	// Ref: https://github.com/suborbital/subo/issues/347
	// create.AddCommand(command.CreateHandlerCmd()).

	// se2 related commands.
	cmd.AddCommand(se2Command())

	// docs related commands.
	cmd.AddCommand(docsCommand())

	cmd.AddCommand(create)
	cmd.AddCommand(command.BuildCmd())

	// TODO: Re-enable when dev is updated to work with e2core
	// cmd.AddCommand(command.DevCmd())

	cmd.AddCommand(command.CleanCmd())

	if features.EnableRegistryCommands {
		cmd.AddCommand(command.PushCmd())
		cmd.AddCommand(command.DeployCmd())
	}

	return cmd
}

func se2Command() *cobra.Command {
	se2 := &cobra.Command{
		Use:   "se2",
		Short: "SE2 related resources",
		Long:  `manage Suborbital Extension Engine (SE2) resources`,
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create SE2 resources",
		Long:  `create Suborbital Extension Engine (SE2) resources`,
	}

	create.AddCommand(command.SE2CreateTokenCommand())

	se2.AddCommand(create)
	se2.AddCommand(command.SE2DeployCommand())
	se2.AddCommand(command.SE2MigrateStorageCommand())

	return se2
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
