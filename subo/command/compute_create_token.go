package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/input"
	"github.com/suborbital/subo/subo/scn"
)

// ComputeCreateTokenCommand returns the dev command
func ComputeCreateTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [email]",
		Short: "create a Compute Network token",
		Long:  `create a Compute Network token`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			email := args[0]

			SCN := scn.New()

			verifier, err := SCN.CreateEmailVerifier(email)
			if err != nil {
				return errors.Wrap(err, "failed to CreateEmailVerifier")
			}

			fmt.Print("A verification code was sent to your email address. Enter the code to continue: ")
			code, err := input.ReadStdinString()
			if err != nil {
				return errors.Wrap(err, "failed to ReadStdinString")
			}

			if len(code) != 6 {
				return errors.New("code must be 6 characters in length")
			}

			reqVerifier := &scn.RequestVerifier{
				UUID: verifier.UUID,
				Code: code,
			}

			token, err := SCN.CreateEnvironmentToken(reqVerifier)
			if err != nil {
				return errors.Wrap(err, "failed to CreateEnvironmentToken")
			}

			fmt.Println(token.Token)

			return nil
		},
	}

	return cmd
}
