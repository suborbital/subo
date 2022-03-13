package util

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Run runs a command, outputting to terminal and returning the full output and/or error.
func Run(cmd string) (string, error) {
	return run(cmd, "", false)
}

// RunInDir runs a command in the specified directory and returns the full output or error.
func RunInDir(cmd, dir string) (string, error) {
	return run(cmd, dir, false)
}

// RunSilent runs a command without printing to stdout and returns the full output or error.
func RunSilent(cmd string) (string, error) {
	return run(cmd, "", true)
}

func run(cmd, dir string, silent bool) (string, error) {
	// you can uncomment this below if you want to see exactly the commands being run
	// fmt.Println("▶️", cmd).

	command := exec.Command("sh", "-c", cmd)

	command.Dir = dir

	var outBuf bytes.Buffer

	if silent {
		command.Stdout = &outBuf
		command.Stderr = &outBuf
	} else {
		command.Stdout = io.MultiWriter(os.Stdout, &outBuf)
		command.Stderr = io.MultiWriter(os.Stderr, &outBuf)
	}

	runErr := command.Run()

	outStr := outBuf.String()

	if runErr != nil {
		return outStr, errors.Wrap(runErr, "failed to Run command")
	}

	return outStr, nil
}
