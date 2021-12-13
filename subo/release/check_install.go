package release

import (
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-version"
)

func InstallInfo(v *version.Version) string {

	_, err := exec.Command("brew", "list", "subo").Output()

	if err != nil {
		return fmt.Sprintf("Head over to https://github.com/suborbital/subo/releases/tag/v%s or pull new changes", v)
	}

	return "Run: 'brew upgrade subo' to update subo"

}
