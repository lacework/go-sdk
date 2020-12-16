package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/lacework/go-sdk/lwlogger"
)

func main() {
	var fileFlag = flag.String("f", "/tmp/out.log", "log file")
	flag.Parse()

	logWriter, err := os.OpenFile(*fileFlag, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
	if err != nil {
		fmt.Println("ERROR unable to open file to initialize logger: %s", err)
		os.Exit(1)
	}

	// This function allows you to pass any io.Writer
	// as an example we show how to write to a file
	lwL := lwlogger.NewWithWriter("DEBUG", logWriter)

	// Content of /tmp/out.log
	// {"level":"info","ts":"2020-12-16T10:47:16+01:00","caller":"io-writer/main.go:25","msg":"interesting info"}
	lwL.Info("interesting info")
}
