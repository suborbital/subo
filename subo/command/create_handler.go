package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/subo/builder/context"
	"github.com/suborbital/subo/subo/util"
)

type handlerData struct {
	HandlerType string
	Request     string
	Method      string
}

func CreateHandlerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler <resource>",
		Short: "create a new handler",
		Long:  `create a new handler for Subo CLI`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resource := args[0]

			handlerType, _ := cmd.Flags().GetString(typeFlag)
			method, _ := cmd.Flags().GetString(methodFlag)

			util.LogStart(fmt.Sprintf("creating handler for %s", resource))

			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			bctx, err := context.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			if bctx.Directive == nil {
				return errors.New("cannot create handler, Directive.yaml not found")
			}

			//Create a new handler object
			handler := directive.Handler{
				Input: directive.Input{
					Type:     handlerType,
					Resource: resource,
					Method:   method,
				},
			}

			//Add the handler object to the directive file
			bctx.Directive.Handlers = append(bctx.Directive.Handlers, handler)

			//Write Directive File which overwrites the entire file
			if err := context.WriteDirectiveFile(bctx.Cwd, bctx.Directive); err != nil {
				return errors.Wrap(err, "failed to WriteDirectiveFile")
			}

			util.LogDone(fmt.Sprintf("handler for %s created", resource))

			return nil
		},
	}

	cmd.Flags().String(typeFlag, "request", "the handler's input type")
	cmd.Flags().String(methodFlag, "GET", "the HTTP method for 'request' handlers")

	return cmd
}
