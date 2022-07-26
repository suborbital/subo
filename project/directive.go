package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/velo/cli/util"
)

// WriteDirectiveFile writes a Directive to disk.
func WriteDirectiveFile(cwd string, localDirective *directive.Directive) error {
	filePath := filepath.Join(cwd, "Directive.yaml")

	directiveBytes, err := yaml.Marshal(localDirective)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal")
	}

	if err := ioutil.WriteFile(filePath, directiveBytes, util.PermFilePrivate); err != nil {
		return errors.Wrap(err, "failed to WriteFile")
	}

	return nil
}

// readDirectiveFile finds a Directive from disk but does not validate it.
func readDirectiveFile(cwd string) (*directive.Directive, error) {
	filePath := filepath.Join(cwd, "Directive.yaml")

	_, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Stat Directive")
	}

	directiveBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadFile for Directive")
	}

	localDirective := &directive.Directive{}
	if err := localDirective.Unmarshal(directiveBytes); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal Directive")
	}

	return localDirective, nil
}

// readQueriesFile finds a queries.yaml from disk.
func readQueriesFile(cwd string) ([]directive.DBQuery, error) {
	filePath := filepath.Join(cwd, "Queries.yaml")

	directiveBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to ReadFile for Queries.yaml")
	}

	localDirective := &directive.Directive{}
	if err := localDirective.Unmarshal(directiveBytes); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal Directive")
	}

	return localDirective.Queries, nil
}

// AugmentAndValidateDirectiveFns ensures that all functions referenced in a handler exist
// in the project and then adds the function list to the provided directive.
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

	// validate augmented directive.
	return dxe.Validate()
}

// getHandlerFnList gets a full list of all functions used in the directive's handlers.
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
			}
		}
	}

	fns := []string{}
	for fn := range fnMap {
		fns = append(fns, fn)
	}

	return fns
}

func DockerNameFromDirective(d *directive.Directive) (string, error) {
	identParts := strings.Split(d.Identifier, ".")
	if len(identParts) != 3 {
		return "", errors.New("ident has incorrect number of parts")
	}

	org := identParts[1]
	repo := identParts[2]

	name := fmt.Sprintf("%s/%s:%s", org, repo, d.AppVersion)

	return name, nil
}
