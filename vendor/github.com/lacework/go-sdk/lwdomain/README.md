# Lacework Domain

Use this package to disseminate a domain URL into account, cluster and whether or not
it is an internal account.

## Usage

Download the library into your `$GOPATH`:

    $ go get github.com/lacework/go-sdk/lwdomain

Import the library into your tool:

```go
import "github.com/lacework/go-sdk/lwdomain"
```

## Examples

The following URL `https://account.lacework.net` would be disseminated into:
* `account` as the account name

```go
package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/lwdomain"
)

func main() {
	domain, err := lwdomain.New("https://account.lacework.net")
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}

	// Output: Lacework Account Name: account
	fmt.Println("Lacework Account Name: %s", domain.Account)
}
```
