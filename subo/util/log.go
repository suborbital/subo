package util

import (
	"fmt"
	"os"
)

// FriendlyLogger describes a logger designed to provide friendly output for interactive CLI purposes
type FriendlyLogger interface {
	LogInfo(string)
	LogStart(string)
	LogDone(string)
	LogFail(string)
}

// PrintLogger is a struct wrapper around the logging functions used by Subo
type PrintLogger struct{}

func (p *PrintLogger) LogInfo(msg string)  { LogInfo(msg) }
func (p *PrintLogger) LogStart(msg string) { LogStart(msg) }
func (p *PrintLogger) LogDone(msg string)  { LogDone(msg) }
func (p *PrintLogger) LogFail(msg string)  { LogFail(msg) }

// LogInfo logs information
func LogInfo(msg string) {
	if _, exists := os.LookupEnv("SUBO_DOCKER"); !exists {
		fmt.Println(fmt.Sprintf("ℹ️  %s", msg))
	}
}

// LogStart logs the start of something
func LogStart(msg string) {
	if _, exists := os.LookupEnv("SUBO_DOCKER"); !exists {
		fmt.Println(fmt.Sprintf("⏩ START: %s", msg))
	}
}

// LogDone logs the success of something
func LogDone(msg string) {
	if _, exists := os.LookupEnv("SUBO_DOCKER"); !exists {
		fmt.Println(fmt.Sprintf("✅ DONE: %s", msg))
	}
}

// LogFail logs the failure of something
func LogFail(msg string) {
	if _, exists := os.LookupEnv("SUBO_DOCKER"); !exists {
		fmt.Println(fmt.Sprintf("🚫 FAILED: %s", msg))
	}
}
