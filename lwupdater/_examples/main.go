package main

import (
	"fmt"

	"github.com/lacework/go-sdk/lwupdater"
)

func main() {
	var (
		project  = "go-sdk"
		sdk, err = lwupdater.Check(project, "v0.1.0")
	)

	if err != nil {
		fmt.Println("Unable to check for updates: %s", err)
	} else {
		// Output: The latest release of the go-sdk project is v0.1.7
		fmt.Printf("The latest release of the %s project is %s\n",
			project, sdk.Latest,
		)
	}
}
