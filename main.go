package main

import (
	"context"
	"os"
	"time"

	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

const checkVersionTimeout = 500 * time.Millisecond

func main() {
	migrateCache()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	versionChan := checkVersion(ctx)

	rootCmd := rootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}

	select {
	case msg := <-versionChan:
		if msg != "" {
			util.LogInfo(msg)
		}
	case <-time.After(checkVersionTimeout):
		util.LogFail("failed to CheckForLatestVersion due to timeout")
	}
}

func checkVersion(ctx context.Context) chan string {
	versionChan := make(chan string)

	go func() {
		msg := ""
		if version, err := release.CheckForLatestVersion(); err != nil {
			msg = err.Error()
		} else if version != "" {
			msg = version
		}

		select {
		case <-ctx.Done():
		default:
			versionChan <- msg
		}
	}()

	return versionChan
}
