package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/suborbital/subo/subo/input"
	"github.com/suborbital/subo/subo/util"
)

const (
	se2EndpointEnvKey = "SUBO_SE2_ENDPOINT"
)

// SE2CreateTokenCommand returns the dev command.
func SE2CreateTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [email]",
		Short: "create an SE2 token",
		Long:  `create a Suborbital Extension Engine (SE2) token`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			email := args[0]

			vapi, err := se2API().ForVerifiedEmail(email, getVerifierCode)
			if err != nil {
				return errors.Wrap(err, "failed to ForVerifiedEmail")
			}

			token, err := vapi.CreateEnvironmentToken()
			if err != nil {
				return errors.Wrap(err, "failed to CreateEnvironmentToken")
			}

			fmt.Println(token.Token)

			if err := util.WriteEnvironmentToken(token.Token); err != nil {
				return errors.Wrap(err, "failed to WriteEnvironmentToken")
			}
			return nil
		},
	}

	return cmd
}

// getVerifierCode gets the 6-character code from stdin.
func getVerifierCode() (string, error) {
	fmt.Print("A verification code was sent to your email address. " +
		"Enter the code to continue, " +
		"and your environment token will print below (keep it safe!): ")
	code, err := input.ReadStdinString()
	if err != nil {
		return "", errors.Wrap(err, "failed to ReadStdinString")
	}

	if len(code) != 6 {
		return "", errors.New("code must be 6 characters in length")
	}

	return code, nil
}
