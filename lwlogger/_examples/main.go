package main

import "github.com/lacework/go-sdk/lwlogger"

func main() {
	lwL := lwlogger.New("INFO")

	// Output: {"level":"info","ts":"[timestamp]","caller":"main.go:9","msg":"interesting info"}
	lwL.Info("interesting info")
}
