package release

import "fmt"

// These variables are set at buildtime. See the Makefile.
var CommitHash = ""
var BuildTime = ""

func Version() string {
	if CommitHash != "" && BuildTime != "" {
		return fmt.Sprintf(`%s %s (Built at %s)`, SuboVersion, CommitHash, BuildTime)
	}
	return SuboVersion
}
