# API Client

A Golang API client for interacting with [Lacework APIs](https://docs.lacework.net/api/about-the-lacework-api).

## Usage

Download the library into your `$GOPATH`:

    $ go get github.com/lacework/go-sdk/api

Import the library into your tool:

```go
import "github.com/lacework/go-sdk/api"
```

## Requirements

To interact with Lacework's API you need to have:

1. A Lacework account
2. Either API access keys or token for authentication

## Examples

Create a new Lacework client that will automatically generate a new access token
from the provided set of API keys, then hit the `/api/v2/AlertChannels` endpoint
to list all available alert channels in your account:
```go
package main

import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient("account",
		api.WithTokenFromKeys("KEY", "SECRET"),
	)
	if err != nil {
		log.Fatal(err)
	}

	alertChannels, err := lacework.V2.AlertChannels.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, channel := range alertChannels.Data {
		fmt.Printf("Alert channel: %s\n", channel.Name)
	}
	// Output:
	//
	// Alert channel: DEFAULT EMAIL
}
```

Look at the [_examples/](_examples/) folder for more examples.
