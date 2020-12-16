package main

import "github.com/lacework/go-sdk/lwlogger"

func main() {
	lwL := lwlogger.New("INFO")

	// Output: {"level":"info","ts":"2020-12-16T10:48:08+01:00","caller":"simple-usage/main.go:9","msg":"interesting info"}
	lwL.Info("interesting info")

	// This message wont be displayed
	lwL.Debug("here is a debug message, too long and only needed when debugging this app")
}
