# Lacework Time Library

A simple relative and natural time package.

## Usage

Download the library into your `$GOPATH`:

    $ go get github.com/lacework/go-sdk/lwtime

Import the library:

```go
import "github.com/lacework/go-sdk/lwtime"
```

## Relative Time Specifiers

Relative times allow you to represent time values dynamically, using specifiers that represent an offset from the current time. For instance, a relative time of `-24h` produces a date/time that is 24 hours less the current time. Relative times can also snap to a particular time. For instance, a relative time of `@d` would represent the start of the current day.

For example, to generate a time range (using a start and end time) that represents the previous day:
```go
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
```

A relative time has three components:
* A signed (+/-) integer
* A relative time unit
* A relative time snap

Lacework supports the following relative time units:
* y - year
* mon - month
* w - week
* d - day
* h - hour
* m - minute
* s - second

Additional considerations include:
* To represent the current time, you can specify either `now` or `+0s`.
* When specifying an integer and relative time unit, snaps are optional.
* When specifying a snap, the integer and relative time unit are optional. For instance, `@d` is actually interpreted as `+0s@d`.


## Natural Time Ranges

Natural time ranges allow you to represent time range values using natural language. For instance, a natural time range of `yesterday` represents a relative start time of `-1d@d` and a relative end time of `@d`.

For example, to generate a time range of this month:
```go
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
```

A natural time has three components:
* An adjective
* A positive number (only when using the last adjective)
* The full text representation of a relative time unit (i.e., year/years)

Lacework supports the following adjectives (disambiguating previous and last by design):
* this/current
* previous
* last

Additional considerations include:
* `last` implies "in the last". So last week reads as "in the last week" and represents a start time of `-1w` and an end time of `now`.
* `previous` always snaps. So "previous week" represents a start time of `-1w@w` and an end time of `@w`.
* `yesterday` is a valid natural time and is equivalent to previous day.
* `today` is a valid natural time and is equivalent to this day or current day.

## Examples

Look at the [_examples/](_examples/) folder for more examples.
