package main

import (
	"os"

	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

func main() {
	rootCmd := rootCommand()

	done := update()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}

	<-done
}

func update() chan bool {
	done := make(chan bool)

	go func() {
		version, err := release.CheckForLatestVersion()
		if err != nil {
			util.LogFail(err.Error())
		} else if version != "" {
			util.LogInfo(version)
		}
		done <- true
	}()

	return done
}
