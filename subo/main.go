package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := rootCommand()

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
