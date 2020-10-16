package main

import (
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/command"
	"github.com/suborbital/subo/subo/context"
	"github.com/suborbital/subo/subo/release"
)

func rootCommand(bctx *context.BuildContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subo",
		Short:   "Suborbital Development Platform CLI",
		Long:    `subo includes a full toolchain for using and managing Suborbital Development Platform tools, including building Wasm Runnables and serving Wasm bundles.`,
		Version: release.SuboDotVersion,
	}

	cmd.SetVersionTemplate("{{.Version}}\n")

	cmd.AddCommand(command.BuildCmd(bctx))
	cmd.AddCommand(command.CreateCmd(bctx))

	return cmd
}
