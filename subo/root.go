package main

import (
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/command"
	"github.com/suborbital/subo/subo/release"
)

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subo",
		Short:   "Suborbital Development Platform CLI",
		Long:    `subo includes a full toolchain for using and managing Suborbital Development Platform tools, including building Wasm Runnables and Atmo projects.`,
		Version: release.SuboDotVersion,
	}

	cmd.SetVersionTemplate("{{.Version}}\n")

	create := &cobra.Command{
		Use:     "create",
		Short:   "create an element",
		Long:    `create an element such as a new project or runnable`,
		Version: release.SuboDotVersion,
	}

	create.AddCommand(command.CreateProjectCmd())
	create.AddCommand(command.CreateRunnableCmd())
	cmd.AddCommand(create)

	cmd.AddCommand(command.BuildCmd())

	return cmd
}
