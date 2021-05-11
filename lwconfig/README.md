# Lacework Config

A Go library to help you manage the Lacework configuration file (`$HOME/.lacework.toml`)

## Usage

Download the library into your `$GOPATH`:

    $ go get github.com/lacework/go-sdk/lwconfig

Import the library into your tool:

```go
import "github.com/lacework/go-sdk/lwconfig"
```

## Examples

Load the default Lacework configuration file and detect if there is a profile named `test`:
```go
package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwconfig"
)

func main() {
	profiles, err := lwconfig.LoadProfiles()
	if err != nil {
		fmt.Printf("Error trying to load profiles: %s\n", err)
		os.Exit(1)
	}

	config, ok := profiles["test"]
	if !ok {
		fmt.Println("You have a test profile configured!")
	} else {
		fmt.Println("'test' profile not found")
	}
}
```

Look at the [examples/](examples/) folder for more examples.
