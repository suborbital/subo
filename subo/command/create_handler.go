package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/builder/context"
	"github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
	"gopkg.in/yaml.v2"
)

type handlerData struct {
	Name        string
	Environment string
	Headless    bool
	APIVersion  string
	AtmoVersion string
}

//GOAL
// Would love to be able to run subo create handler /foo and have it added to the directive. 
// The basic command can just create a handler with a placeholder function. 
// Some potential flags: --method, --stream, --steps

//How:  Construct the new handler using the values provided by the user with good defaults otherwise
//How2: Append the new handler to the currently list of handlers


//Psudocode:
//Find how create_project created a handler AND the directive file

//Questions: type and resource always going to be the same or no, would those be flags too? Because type and resource are already set in the example project.
// cwd, err := os.Getwd() //current working directory
// bctx, err := context.ForDirectory(cwd) //BuildContext


//READ ENTIRE FILE
//TAKE FROM DIRECTIVE.go

// OVERWITE IT

// CreateHandlerError wraps errors for CreateHandlerCmd() failures
type CreateHandlerError struct {
	Path  string // The ouput directory for build command CreateHandlerCmd().
	error        // The original error.
}

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
			resource, _ := cmd.Flags().GetString(resourceFlag)
			method, _ := cmd.Flags().GetString(methodFlag)
			// stream, _ := cmd.Flags().GetString(streamFlag)
			// steps, _ := cmd.Flags().GetString(stepsFlag)

			//STEPS: In order to create a handler what information do you need: 
			//1. What directory you are pulling the runnable from- actually NOT PULLING any runnable, 
			//we are just creating handler document 

			//So, if they put /foo we need to put in this handler: "resource: /foo"
			//Handler will look like: 
			// handlers:
			//   - type: request
			//     resource: /foo
			//     method: GET
			//  
			util.LogStart(fmt.Sprintf("creating handler with function name %s", name))

			//Do I need to set up infrastructure where I need to check that the placeholder fn actually exisst in our runnables list?

			//Need to grab lib.ts (from templates) and be able to pass that into the handler
			// fmt.Sprintf("foo", input) //this is going to be a template literal because we have to pass in the input- 

			handler, err := writeHandler(bctx.Cwd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to writeHandler")
			}

			util.LogDone(path)

			if _, err := util.Run(fmt.Sprintf("git init ./%s", name)); err != nil {
				return errors.Wrap(err, "🚫 failed to initialize Run git init")
			}

			return nilq
		},
	}


	cmd.Flags().String(typeFlag, "POST", "the method for which you want ")
	cmd.Flags().String(resourceFlag, "main", "git branch to download templates from") //
	cmd.Flags().String(methodFlag, "POST", "the method for which you want ")
	// cmd.Flags().String(streamFlag, "main", "git branch to download templates from") //stream is a stype of handler 
	// cmd.Flags().String(stepsFlag, "fn", "Runnable functions to be composed when handling requests to the resource.")

	// cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	// cmd.Flags().String(environmentFlag, "com.suborbital", "project environment name (your company's reverse domain")
	// cmd.Flags().Bool(headlessFlag, false, "use 'headless mode' for this project")
	// cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest templates")

	return cmd
}

func writeHandler(cwd, handlerType, resource, method, string) (*directive.Runnable, error) { //notate optional params      , method, lang, namespace string
	
	handler := &directive.Handler{
		HandlerType:       handlerType,
		Resource:   resource,
		Method:     method
	}
		// Stream:     stream,
		// Steps:  	steps,
	

	bytes, err := yaml.Marshal(handler)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal handler")
	}

	path := filepath.Join(cwd, name, ".handler.yaml") //How to get it to Directive.yaml? Is this being appended to the directive.yaml.tmpl file?

	if err := ioutil.WriteFile(path, bytes, 0700); err != nil { //Don't need to write entire file
		return nil, errors.Wrap(err, "failed to WriteFile handler")
	}

	return handler, nil
}

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