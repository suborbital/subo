package command

import (
	"fmt"

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
		Use:   "handler <name>",
		Short: "create a new handler",
		Long:  `create a new handler for Subo CLI`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			handlerType, _ := cmd.Flags().GetString(typeFlag)
			resource, _ := cmd.Flags().GetString(resourceFlag) //name and the resource are the same
			method, _ := cmd.Flags().GetString(methodFlag)

			dir, _ := cmd.Flags().GetString(dirFlag)

			util.LogStart(fmt.Sprintf("creating handler with function name %s", name))

			bctx, err := context.ForDirectory(dir)
			if err != nil {
				return errors.Wrap(err, "failed to ForDirectory")
			}

			if bctx.Directive == nil {
				util.LogFail("Handlers must be created in a project")
				return errors.New("Directive.yaml not found")
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

			util.LogDone(fmt.Sprintf("handler with resource name %s created", name))

			return nil
		},
	}

	cmd.Flags().String(typeFlag, "request", "the method for which you want ")
	cmd.Flags().String(resourceFlag, "/foo", "git branch to download templates from")
	cmd.Flags().String(methodFlag, "GET", "the method for which you want ")

	return cmd
}
