package release

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

func InstallInfo(v *version.Version) (string, error) {
	check, dir, err := checkInstall()

	if err != nil {
		return "", errors.Wrap(err, "failed to check installation method")
	}

	if check == 0 {
		return "You should consider upgrading via the 'brew upgrade subo' command, then build again using 'make subo'", nil
	} else if check == 1 {
		return fmt.Sprintf("You should consider upgrading via the 'git -C %s pull', then build again using 'make subo'", dir), nil
	} else {
		return fmt.Sprintf("You should consider upgrading, go to https://github.com/suborbital/subo/releases/tag/v%s and download the latest source code, then build again using 'make subo'", v), nil
	}
}

func checkInstall() (int, string, error) {
	//Checks if subo formula found
	_, err := exec.Command("brew", "list", "subo").Output()

	if err != nil {
		mydir, err := os.Getwd()

		if err != nil {
			return -1, "", errors.Wrap(err, "Failed to get working directory")
		}

		var remote = "https://github.com/suborbital/subo.git"

		//gets remote origin of the local repository
		cmd, err := exec.Command("git", "-C", mydir, "remote", "get-url", "--all", "origin").Output()

		if err != nil {
			return 2, "", nil
		}

		if strings.TrimRight(string(cmd), "\n") == remote {
			return 1, mydir, nil
		}

		return 2, "", nil
	}

	return 0, "", nil
}
