# Go API Client

A Golang API client for interacting with the [Lacework API](https://support.lacework.com/hc/en-us/categories/360002496114-Lacework-API-).

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
from the provided set of API keys, then hit the `/external/integrations` endpoint
to list all available integrations from your account:
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

	integrations, err := lacework.Integrations.List()
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// CUSTOMER_123456B DATADOG
	// CUSTOMER_123456A CONT_VULN_CFG
	// CUSTOMER_123456C PAGER_DUTY_API
	fmt.Println(integrations.String())
}
```

Look at the [_examples/](_examples/) folder for more examples.
