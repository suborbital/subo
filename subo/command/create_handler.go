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

//GOAL
// Would love to be able to run subo create handler /foo and have it added to the directive.
// The basic command can just create a handler with a placeholder function.
// Some potential flags: --method, --stream, --steps

//READ ENTIRE FILE
// Write CreateHandler Function
// CreateHandlerCmd returns the build command
func CreateHandlerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler <name>", //this <name> has to pull a default runnable which is in form: /name
		Short: "create a new handler",
		Long:  `create a new handler for Subo CLI`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			handlerType, _ := cmd.Flags().GetString(typeFlag)
			resource, _ := cmd.Flags().GetString(resourceFlag) //name and the resource are the same
			method, _ := cmd.Flags().GetString(methodFlag)
			// stream, _ := cmd.Flags().GetString(streamFlag)
			// steps, _ := cmd.Flags().GetString(stepsFlag)
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
			// //create new handler object

			handler := directive.Handler{
				Input: directive.Input{
					Type:     handlerType,
					Resource: resource,
					Method:   method,
				},
			}

			// Stream:     stream,
			// Steps:  	steps,

			//add handler object to the directive
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
	cmd.Flags().String(resourceFlag, "/foo", "git branch to download templates from") //
	cmd.Flags().String(methodFlag, "GET", "the method for which you want ")
	// cmd.Flags().String(streamFlag, "main", "git branch to download templates from") //stream is a stype of handler
	// cmd.Flags().String(stepsFlag, "fn", "Runnable functions to be composed when handling requests to the resource.")

	return cmd
}

// // WriteDirectiveFile writes a Directive to disk
// func WriteDirectiveFile(cwd string, directive *directive.Directive) error {
// 	filePath := filepath.Join(cwd, "Directive.yaml")

// 	directiveBytes, err := yaml.Marshal(directive)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to Marshal")
// 	}

// 	if err := ioutil.WriteFile(filePath, directiveBytes, os.FileMode(os.O_WRONLY)); err != nil {
// 		return errors.Wrap(err, "failed to WriteFile")
// 	}

// 	return nil
// }

//Writes the handler
// func writeHandler(cwd string, handlerType, resource, method, string) (*directive.Runnable, error) { //notate optional params      , method, lang, namespace string

// 	handler := &directive.Handler{
// 		HandlerType:       handlerType,
// 		Resource:   resource,
// 		Method:     method,
// 		// Stream:     stream,
// 		// Steps:  	steps,
// 	}

// 	//Not sure if I need this at all
// 	bytes, err := yaml.Marshal(handler)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to Marshal handler")
// 	}

// 	path := filepath.Join(cwd, name, ".handler.yaml") //How to get it to Directive.yaml? Is this being appended to the directive.yaml.tmpl file?

// 	if err := ioutil.WriteFile(path, bytes, 0700); err != nil { //Don't need to write entire file
// 		return nil, errors.Wrap(err, "failed to WriteFile handler")
// 	}

// 	return handler, nil
// }

// func appendHandler(cwd, method, stream, steps string) (*directive.Runnable, error) {

// }

//NOTES

//Handler Structure
// handlers:
//   - type: request
//     resource: /hello
//     method: POST
//     steps:
//       - group:
//         - fn: modify-url
//         - fn: helloworld-rs
//           as: hello
//       - fn: fetch-test
//         with:
//           url: modify-url
//           logme: hello

//
//
//					IN ATMO
//
// // Handlers returns the handlers for the app
// func (h *HeadlessBundleSource) Handlers() []directive.Handler {
// 	if h.bundlesource.bundle == nil {
// 		return []directive.Handler{}
// 	}

// 	handlers := []directive.Handler{}

// 	// for each Runnable, construct a handler that executes it
// 	// based on a POST request to its FQFN URL /identifier/namespace/fn/version
// 	for _, runnable := range h.bundlesource.Runnables() {
// 		handler := directive.Handler{
// 			Input: directive.Input{
// 				Type:     directive.InputTypeRequest,
// 				Method:   http.MethodPost,
// 				Resource: fqfn.Parse(runnable.FQFN).HeadlessURLPath(),
// 			},
// 			Steps: []executable.Executable{
// 				{
// 					CallableFn: executable.CallableFn{
// 						Fn:   runnable.Name,
// 						With: map[string]string{},
// 						FQFN: runnable.FQFN,
// 					},
// 				},
// 			},
// 		}

// 		handlers = append(handlers, handler)
// 	}

// 	return handlers
// }
