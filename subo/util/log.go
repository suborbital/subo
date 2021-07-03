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
	LogWarn(string)
}

// PrintLogger is a struct wrapper around the logging functions used by Subo
type PrintLogger struct{}

func (p *PrintLogger) LogInfo(msg string)  { LogInfo(msg) }
func (p *PrintLogger) LogStart(msg string) { LogStart(msg) }
func (p *PrintLogger) LogDone(msg string)  { LogDone(msg) }
func (p *PrintLogger) LogFail(msg string)  { LogFail(msg) }
func (p *PrintLogger) LogWarn(msg string)  { LogWarn(msg) }

// Keeping it DRY
func Log(msg string) {
	if _, exists := os.LookupEnv("SUBO_DOCKER"); !exists {
		fmt.Println(msg)
	}
}

// LogInfo logs information
func LogInfo(msg string) {
	Log(fmt.Sprintf("ℹ️  %s", msg))
}

// LogStart logs the start of something
func LogStart(msg string) {
	Log(fmt.Sprintf("⏩ START: %s", msg))
}

// LogDone logs the success of something
func LogDone(msg string) {
	Log(fmt.Sprintf("✅ DONE: %s", msg))
}

// LogFail logs the failure of something
func LogFail(msg string) {
	Log(fmt.Sprintf("🚫 FAILED: %s", msg))
}

// LogWarn logs a warning from something
func LogWarn(msg string) {
	Log(fmt.Sprintf("⚠️ WARNING: %s", msg))
}
