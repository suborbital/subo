package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/context"
)

func main() {
	bctx, err := context.CurrentBuildContext()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to get CurrentBuildContext"))
	}

	rootCmd := rootCommand(bctx)

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {

		return nil
	}

	rootCmd.Execute()
}
