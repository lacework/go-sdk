# Lacework Logger

A Logger wrapper for Lacework based of [zap](https://github.com/uber-go/zap) logger Go package.

## Usage

Download the library into your `$GOPATH`:

    $ go get github.com/lacework/go-sdk/lwlogger

Import the library into your tool:

```go
import "github.com/lacework/go-sdk/lwlogger"
```

## Environment Variables

This package can be controlled via environment variables:

| Environment Variable | Description | Default | Supported Options |
|----------------------|-------------|---------|-------------------|
|`LW_LOG`|Change the verbosity of the logs |`""`| `INFO` or `DEBUG` |
|`LW_LOG_FORMAT`|Controls the format of the logs|`JSON`| `JSON` or `CONSOLE` |
|`LW_LOG_DEV`|Switch the logger instance to development mode (extra verbose)|`false`| `true` or `false` |

## Examples

To create a new logger instance with the log level `INFO`, write an interesting
info message and another debug message. Note that only the info message will be
displayed:
```go
package main

import "github.com/lacework/go-sdk/lwlogger"

func main() {
	lwL := lwlogger.New("INFO")

	// Output: {"level":"info","ts":"[timestamp]","caller":"main.go:9","msg":"interesting info"}
	lwL.Info("interesting info")

	// This message wont be displayed
	lwL.Debug("this is a debug message, too long and only needed when debugging this app")
}
```

Look at the [examples/](examples/) folder for more examples.
