package command

import (
	"github.com/spf13/cobra"
)

// DevCmd returns the dev command.
func DevCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "run a development Atmo server using Docker",
		Long:  `run a development Atmo server using Docker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			build := BuildCmd()
			if err := build.Execute(); err != nil {
				return err
			}

			start := StartCmd()
			if err := start.Execute(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(appNameFlag, "Velocity", "if passed, it'll be used as VELOCITY_APP_NAME, otherwise 'Velocity' will be used")
	cmd.Flags().String(runPartnerFlag, "", "if passed, the provided command will be run as the partner application")
	cmd.Flags().String(domainFlag, "", "if passed, it'll be used as VELOCITY_DOMAIN and HTTPS will be used, otherwise HTTP will be used")
	cmd.Flags().Int(httpPortFlag, 8080, "if passed, it'll be used as VELOCITY_HTTP_PORT, otherwise '8080' will be used")
	cmd.Flags().Int(tlsPortFlag, 443, "if passed, it'll be used as VELOCITY_TLS_PORT, otherwise '443' will be used")

	return cmd
}
