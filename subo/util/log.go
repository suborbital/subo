package util

import (
	"fmt"
	"os"
)

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
