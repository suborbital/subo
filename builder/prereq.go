package builder

import (
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/suborbital/subo/project"
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
				Command: "curl -L https://github.com/suborbital/reactr/archive/v{{ .RunnableDir.Runnable.APIVersion }}.tar.gz -o _lib/_lib.tar.gz",
			},
			Prereq{
				File:    "_lib/suborbital",
				Command: "tar --strip-components=3 -C _lib -xvzf _lib/_lib.tar.gz **/api/grain/suborbital/*",
			},
		},
		"assemblyscript": {
			Prereq{
				File:    "node_modules",
				Command: "{{ .BuildConfig.JsToolchain }} install",
			},
		},
		"tinygo": {},
		"typescript": {
			Prereq{
				File:    "node_modules",
				Command: "{{ .BuildConfig.JsToolchain }} install",
			},
		},
		"javascript": {
			Prereq{
				File:    "node_modules",
				Command: "{{ .BuildConfig.JsToolchain }} install",
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
				Command: "curl -L https://github.com/suborbital/reactr/archive/v{{ .RunnableDir.Runnable.APIVersion }}.tar.gz -o _lib/_lib.tar.gz",
			},
			Prereq{
				File:    "_lib/suborbital",
				Command: "tar --wildcards --strip-components=3 -C _lib -xvzf _lib/_lib.tar.gz **/api/grain/suborbital/*",
			},
		},
		"assemblyscript": {
			Prereq{
				File:    "node_modules",
				Command: "{{ .BuildConfig.JsToolchain }} install",
			},
		},
		"tinygo": {},
		"typescript": {
			Prereq{
				File:    "node_modules",
				Command: "{{ .BuildConfig.JsToolchain }} install",
			},
		},
		"javascript": {
			Prereq{
				File:    "node_modules",
				Command: "{{ .BuildConfig.JsToolchain }} install",
			},
		},
	},
}

// GetCommand takes a RunnableDir, and returns an executed template command string.
func (p Prereq) GetCommand(b BuildConfig, r project.RunnableDir) (string, error) {
	cmdTmpl, err := template.New("cmd").Parse(p.Command)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse prerequisite Command string into template: %s", p.Command)
	}

	type TemplateParams struct {
		RunnableDir project.RunnableDir
		BuildConfig
	}

	data := TemplateParams{
		RunnableDir: r,
		BuildConfig: b,
	}

	var fullCmd strings.Builder
	err = cmdTmpl.Execute(&fullCmd, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute prerequisite Command string with runnableDir")
	}

	return fullCmd.String(), nil
}
