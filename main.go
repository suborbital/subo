package main

import (
	"context"
	"os"
	"time"

	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

const checkVersionTimeout = 2500 * time.Millisecond

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := checkVersion(ctx)

	rootCmd := rootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}

	select {
	case <-done:
	case <-time.After(checkVersionTimeout):
		util.LogFail("failed to CheckForLatestVersion due to timeout")
	}
}

func checkVersion(ctx context.Context) chan bool {
	done := make(chan bool)

	go func() {
		if version, err := release.CheckForLatestVersion(); err != nil {
			util.LogFail(err.Error())
		} else if version != "" {
			util.LogInfo(version)
		}
		select {
		case <-ctx.Done():
		default:
			done <- true
		}
	}()

	return done
}
