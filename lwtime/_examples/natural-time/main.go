package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwtime"
)

func main() {
	start, end, err := lwtime.ParseNatural("this month")
	if err != nil {
		fmt.Println("Unable to parse natural time: %s", err)
		os.Exit(1)
	}
	// The time range is 2023-07-01 07:00:00 +0000 UTC to 2023-07-13 01:23:59.921851 +0000 UTC
	fmt.Printf("The time range is %s to %s\n", start.String(), end.String())
}
