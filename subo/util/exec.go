package util

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Run runs a command, outputting to terminal and returning the full output and/or error
func Run(cmd string) (string, string, error) {
	return RunInDir(cmd, "")
}

// RunInDir runs a command in the specified directory and returns the full output or error
func RunInDir(cmd, dir string) (string, string, error) {
	// you can uncomment this below if you want to see exactly the commands being run
	// fmt.Println("▶️", cmd)

	argL := strings.Split(cmd, " ")
	command := exec.Command(argL[0], argL[1:]...)

	command.Dir = dir

	var stdoutBuf, stderrBuf bytes.Buffer
	command.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	command.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := command.Run()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to Run command")
	}

	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())

	return outStr, errStr, nil
}
