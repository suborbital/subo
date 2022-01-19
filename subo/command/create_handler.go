package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/builder/context"
	"github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

const (
	defaultRepo = "suborbital/subo"
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

os. package 

// OVERWITE IT

// CreateHandlerCmd returns the build command
func CreateHandlerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler <name>", //this <name> has to pull a default runnable which is in form: /name
		Short: "create a new handler",
		Long:  `create a new handler for Subo CLI`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]


			// type, _ := cmd.Flags().GetString(typeFlag)
			// resource, _ := cmd.Flags().GetString(resourceFlag)
			method, _ := cmd.Flags().GetString(methodFlag)
			stream, _ := cmd.Flags().GetString(streamFlag)
			// steps, _ := cmd.Flags().GetString(stepsFlag)

			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			bctx, err := context.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			util.LogStart(fmt.Sprintf("creating handler with function name %s", name))

			// path, err := util.Mkdir(bctx.Cwd, name)
			// if err != nil {
			// 	return errors.Wrap(err, "🚫 failed to Mkdir")
			// }

			// branch, _ := cmd.Flags().GetString(branchFlag)
			// environment, _ := cmd.Flags().GetString(environmentFlag)
			// headless, _ := cmd.Flags().GetBool(headlessFlag)

			// templatesPath, err := template.TemplateFullPath(defaultRepo, branch)
			// if err != nil {
			// 	return errors.Wrap(err, "🚫 failed to TemplateFullPath")
			// }

			// if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
			// 	templatesPath, err = template.UpdateTemplates(defaultRepo, branch)
			// 	if err != nil {
			// 		return errors.Wrap(err, "🚫 failed to UpdateTemplates")
			// 	}
			// }

			data := handlerData{
				Name:        name,
				Environment: environment,
				APIVersion:  release.FFIVersion,
				AtmoVersion: release.AtmoVersion,
				Headless:    headless,
			}
			//Need to set up code to take from templates
			templatesPath, err := template.TemplateFullPath(repo, branch)
			if err != nil {
				return errors.Wrap(err, "failed to TemplateDir")
			}
					//maybe? 
			if update, _ := cmd.Flags().GetBool(updateTemplatesFlag); update {
				templatesPath, err = template.UpdateTemplates(repo, branch)
				if err != nil {
					return errors.Wrap(err, "🚫 failed to UpdateTemplates")
				}
			}


			//Need to grab lib.ts (from templates) and be able to pass that into the handler
			fmt.Sprintf("foo", input) //this is going to be a template literal because we have to pass in the input- 

			if err := template.ExecTmplDir(bctx.Cwd, name, templatesPath, fmt.Sprintf("%s", input), data); err != nil {
				// if the templates are missing, try updating them and exec again
				if err == template.ErrTemplateMissing {
					templatesPath, err = template.UpdateTemplates(defaultRepo, branch)
					if err != nil {
						return errors.Wrap(err, "🚫 failed to UpdateTemplates")
					}

					if err := template.ExecTmplDir(bctx.Cwd, name, templatesPath, fmt.Sprintf("%s", input), data); err != nil {
						return errors.Wrap(err, "🚫 failed to ExecTmplDir")
					}
				} else {
					return errors.Wrap(err, "🚫 failed to ExecTmplDir")
				}
			}
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

	cmd.Flags().String(methodFlag, "POST", "the method for which you want ")
	cmd.Flags().String(streamFlag, "main", "git branch to download templates from") //stream is a stype of handler 
	cmd.Flags().String(stepsFlag, "fn", "Runnable functions to be composed when handling requests to the resource.")

	// cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	// cmd.Flags().String(environmentFlag, "com.suborbital", "project environment name (your company's reverse domain")
	// cmd.Flags().Bool(headlessFlag, false, "use 'headless mode' for this project")
	// cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest templates")

	return cmd
}

func writeHandler(cwd, method, stream, steps string) (*directive.Runnable, error) { //notate optional params      , method, lang, namespace string
	
	handler := &directive.Handler{
		Method:     method,
		Stream:     stream,
		// Steps:  	steps,
	}

	bytes, err := yaml.Marshal(handler)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal handler")
	}

	path := filepath.Join(cwd, name, ".handler.yaml") //How to get it to Directive.yaml?

	if err := ioutil.WriteFile(path, bytes, 0700); err != nil { //Don't need to write entire file
		return nil, errors.Wrap(err, "failed to WriteFile handler")
	}

	return handler, nil
}

// func appendHandler(cwd, method, stream, steps string) (*directive.Runnable, error) { 

// }

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