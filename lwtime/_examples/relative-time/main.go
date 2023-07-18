package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwtime"
)

func main() {
	start, err := lwtime.ParseRelative("-1d@d")
	if err != nil {
		fmt.Println("Unable to parse start time range: %s", err)
		os.Exit(1)
	}
	end, err := lwtime.ParseRelative("@d")
	if err != nil {
		fmt.Println("Unable to parse end time range: %s", err)
		os.Exit(1)
	}
	// Output: The time range is 2023-07-11 07:00:00 +0000 UTC to 2023-07-12 07:00:00 +0000 UTC
	fmt.Printf("The time range is %s to %s\n", start.String(), end.String())
}
