package context

import (
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

// Prereq is a pre-requisite file paired with the native command needed to acquire that file (if it's missing).
type Prereq struct {
	File    string
	Command string
}

// PreRequisiteCommands is a map of OS : language : preReq.
var PreRequisiteCommands = map[string]map[string][]Prereq{
	"darwin": {
		"rust":  {},
		"swift": {},
		"grain": {
			Prereq{
				File:    "_lib",
				Command: "mkdir _lib",
			},
			Prereq{
				File:    "_lib/_lib.tar.gz",
				Command: "curl -L https://github.com/suborbital/reactr/archive/v{{ .Runnable.APIVersion }}.tar.gz -o _lib/_lib.tar.gz",
			},
			Prereq{
				File:    "_lib/suborbital",
				Command: "tar --strip-components=3 -C _lib -xvzf _lib/_lib.tar.gz **/api/grain/suborbital/*",
			},
		},
		"assemblyscript": {
			Prereq{
				File:    "node_modules",
				Command: "npm install --include=dev",
			},
		},
		"tinygo": {},
		"typescript": {
			Prereq{
				File:    "node_modules",
				Command: "npm install --include=dev",
			},
		},
	},
	"linux": {
		"rust":  {},
		"swift": {},
		"grain": {
			Prereq{
				File:    "_lib",
				Command: "mkdir _lib",
			},
			Prereq{
				File:    "_lib/_lib.tar.gz",
				Command: "curl -L https://github.com/suborbital/reactr/archive/v{{ .Runnable.APIVersion }}.tar.gz -o _lib/_lib.tar.gz",
			},
			Prereq{
				File:    "_lib/suborbital",
				Command: "tar --wildcards --strip-components=3 -C _lib -xvzf _lib/_lib.tar.gz **/api/grain/suborbital/*",
			},
		},
		"assemblyscript": {
			Prereq{
				File:    "node_modules",
				Command: "npm install --include=dev",
			},
		},
		"tinygo": {},
		"typescript": {
			Prereq{
				File:    "node_modules",
				Command: "npm install --include=dev",
			},
		},
	},
}

// GetCommand takes a RunnableDir, and returns an executed template command string.
func (p Prereq) GetCommand(r RunnableDir) (string, error) {
	cmdTmpl, err := template.New("cmd").Parse(p.Command)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse prerequisite Command string into template: %s", p.Command)
	}

	var fullCmd strings.Builder
	err = cmdTmpl.Execute(&fullCmd, r)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute prerequisite Command string with runnableDir")
	}

	return fullCmd.String(), nil
}
