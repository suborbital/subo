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
	}

	cmd.SetVersionTemplate("Subo CLI v{{.Version}}\n")

	create := &cobra.Command{
		Use:     "create",
		Short:   "create a runnable or project",
		Long:    `create a new Atmo project or WebAssembly runnable`,
		Version: release.SuboDotVersion,
	}

	if features.EnableReleaseCommands {
		create.AddCommand(command.CreateReleaseCmd())
	}

	create.AddCommand(command.CreateProjectCmd())
	create.AddCommand(command.CreateRunnableCmd())
	cmd.AddCommand(create)

	cmd.AddCommand(command.BuildCmd())
	cmd.AddCommand(command.DevCmd())

	return cmd
}
