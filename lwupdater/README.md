# Lacework Updater

A Go library to check for available updates of Lacework projects.

## Usage

Download the library into your `$GOPATH`:

    $ go get github.com/lacework/go-sdk/lwupdater

Import the library into your tool:

```go
import "github.com/lacework/go-sdk/lwupdater"
```

## Examples

This example checks for the latest release of this repository (https://github.com/lacework/go-sdk):
```go
package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwupdater"
)

func main() {
	var (
		project  = "go-sdk"
		sdk, err = lwupdater.Check(project, "v0.1.0")
	)

	if err != nil {
		fmt.Println("Unable to check for updates: %s", err)
		os.Exit(1)
	}

	// Output: The latest release of the go-sdk project is v0.1.7
	fmt.Printf("The latest release of the %s project is %s\n",
		project, sdk.Latest,
	)
}
```

Look at the [examples/](examples/) folder for more examples.
