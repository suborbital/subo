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
				cwd = "$HOME"
			}

			bctx, err := project.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "failed to project.ForDirectory")
			}

			if bctx.Directive == nil {
				return errors.New("current directory is not a project; Directive is missing")
			}

			port, _ := cmd.Flags().GetString("port")

			_, err = util.Run(fmt.Sprintf("docker run -v=%s:/home/atmo -e=ATMO_HTTP_PORT=%s -p=%s:%s suborbital/atmo:%s atmo", bctx.Cwd, port, port, port, bctx.AtmoVersion))
			if err != nil {
				return errors.Wrap(err, "🚫 failed to run dev server")
			}

			return nil
		},
	}

	cmd.Flags().String("port", "8080", "set the port that Atmo serves on")

	return cmd
}
