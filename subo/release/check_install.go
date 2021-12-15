package release

import (
	"os/exec"
)

// InstallInfo checks if subo was installed using Homebrew or via some other method and returns update instructions
func InstallInfo() string {
	if _, err := exec.Command("brew", "list", "subo").Output(); err != nil {
		return "Check out our install/upgrade instructions at https://github.com/suborbital/subo#installing"
	}
	return "Run: 'brew upgrade subo' to update subo"
}
