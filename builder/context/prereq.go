package context

// PreReq is a pre-requisite file paired with the native command needed to acquire that file (if it's missing)
type Prereq struct {
	File    string
	Command string
}

// PreRequisiteCommands is a map of OS : language : preReq
var PreRequisiteCommands = map[string]map[string][]Prereq{
	"darwin": {
		"rust":  {},
		"swift": {},
		"assemblyscript": {
			Prereq{
				File:    "node_modules",
				Command: "npm install --include=dev",
			},
		},
		"tinygo": {},
	},
	"linux": {
		"rust":  {},
		"swift": {},
		"assemblyscript": {
			Prereq{
				File:    "node_modules",
				Command: "npm install --include=dev",
			},
		},
		"tinygo": {},
	},
}
