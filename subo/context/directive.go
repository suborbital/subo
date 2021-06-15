package context

import (
	"fmt"

	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/directive"
	"gopkg.in/yaml.v2"
)

// readDirectiveFile finds a Directive from disk but does not validate it
func readDirectiveFile(cwd string) (*directive.Directive, error) {
	filePath := filepath.Join(cwd, "Directive.yaml")

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to Stat Directive")
	}

	directiveBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadFile for Directive")
	}

	directive := &directive.Directive{}
	if err := directive.Unmarshal(directiveBytes); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal Directive")
	}

	return directive, nil
}

// WriteDirective writes a Directive to disk
func WriteDirective(cwd string, directive *directive.Directive) error {
	filePath := filepath.Join(cwd, "Directive.yaml")

	directiveBytes, err := yaml.Marshal(directive)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal")
	}

	if err := ioutil.WriteFile(filePath, directiveBytes, fs.FileMode(os.O_WRONLY)); err != nil {
		return errors.Wrap(err, "failed to WriteFile")
	}

	return nil
}

// AugmentAndValidateDirectiveFns ensures that all functions referenced in a handler exist
// in the project and then adds the function list to the provided directive
func AugmentAndValidateDirectiveFns(dxe *directive.Directive, fns []RunnableDir) error {
	fnMap := map[string]bool{}
	for _, fn := range fns {
		fnMap[fn.Name] = true
	}

	handlerFns := getHandlerFnList(dxe)

	for _, fn := range handlerFns {
		if good, exists := fnMap[fn]; !good || !exists {
			return fmt.Errorf("project does not contain function %s listed in Directive", fn)
		}
	}

	dirRunnables := make([]directive.Runnable, len(fns))
	for i := range fns {
		dirRunnables[i] = *fns[i].Runnable
	}

	dxe.Runnables = dirRunnables

	return nil
}

// getHandlerFnList gets a full list of all functions used in the directive's handlers
func getHandlerFnList(dxe *directive.Directive) []string {
	fnMap := map[string]bool{}

	for _, h := range dxe.Handlers {
		for _, step := range h.Steps {
			if step.IsFn() {
				fnMap[step.Fn] = true
			} else if step.IsGroup() {
				for _, fn := range step.Group {
					fnMap[fn.Fn] = true
				}
			} else if step.IsForEach() {
				fnMap[step.ForEach.Fn] = true
			}
		}
	}

	fns := []string{}
	for fn := range fnMap {
		fns = append(fns, fn)
	}

	return fns
}
