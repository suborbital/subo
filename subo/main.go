package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := rootCommand()

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	rootCmd.Execute()
}
