//go:build !docker

package main

import (
	"context"
	"time"

	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
)

const checkVersionTimeout = 500 * time.Millisecond

func checkForUpdates() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	util.LogInfo("hhelloo?")
	select {
	case msg := <-versionChan:
		if msg != "" {
			util.LogInfo(msg)
		}
	case <-time.After(checkVersionTimeout):
		util.LogFail("failed to CheckForLatestVersion due to timeout")
	}
}
