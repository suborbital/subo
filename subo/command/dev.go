package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

// DevCmd returns the dev command.
func DevCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "run a development Atmo server using Docker",
		Long:  `run a development Atmo server using Docker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			bctx, err := project.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "failed to project.ForDirectory")
			}

			if bctx.Directive == nil {
				return errors.New("current directory is not a project; Directive is missing")
			}

			port, _ := cmd.Flags().GetString("port")
			verbose, _ := cmd.Flags().GetBool("verbose")

			envvar := ""

			if verbose {
				envvar = "-e ATMO_LOG_LEVEL=debug"
			}

			//dockerCmd stores the Docker command to be displayed and executed.
			dockerCmd := fmt.Sprintf("docker run -v=%s:/home/atmo -e=ATMO_HTTP_PORT=%s %s -p=%s:%s suborbital/atmo:%s atmo", bctx.Cwd, port, envvar, port, port, bctx.AtmoVersion)
			fmt.Printf("\nRunning Docker command :\n%s\n\n", dockerCmd)

			_, err = util.Command.Run(dockerCmd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to run dev server")
			}

			return nil
		},
	}

	cmd.Flags().String("port", "8080", "set the port that Atmo serves on")
	cmd.Flags().BoolP("verbose", "v", false, "display debug messages and docker commands")
	cmd.Flags().Lookup("verbose").NoOptDefVal = "true"

	return cmd
}
