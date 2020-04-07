<img src="https://techally-content.s3-us-west-1.amazonaws.com/public-content/lacework_logo_full.png" width="600">

# Lacework Go SDK

This repository provides a Go API client, tools, libraries, relevant documentation, code
samples, processes, and/or guides that allow developers to interact with Lacework.

## API Client ([`api`](api/))

A Golang API client for interacting with the [Lacework API](https://support.lacework.com/hc/en-us/categories/360002496114-Lacework-API-).

### Basic Usage
```go
import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"
)

lacework, err := api.NewClient("account")
if err == nil {
	log.Fatal(err)
}

tokenRes, err := lacework.GenerateTokenWithKeys("KEY", "SECRET")
if err != nil {
	log.Fatal(err)
}

// Output: YOUR-ACCESS-TOKEN
fmt.Println(tokenRes.Token())
```
Look at the [api/](api/) folder for more documentation.

## Lacework CLI ([`cli`](cli/))

The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. You can use it to manage compliance
reports, external integrations, vulnerability scans, and other operations.

### Basic Usage

Build and install the CLI by running `make install-cli`, the automation will
install the tool at `/usr/local/bin/lacework`:
```
$ make install-cli
$ lacework version
lacework v0.1.1 (sha:ca9f95d17f4f2092f89dba7b64eaed6db7493a5a) (time:20200406091143)
```
Look at the [cli/](cli/) folder for more documentation.

## Lacework Logger ([`lwlogger`](lwlogger/))

A Logger wrapper for Lacework based of zap logger Go package.

### Basic Usage
```go
import "github.com/lacework/go-sdk/lwlogger"

func main() {
	lwL := lwlogger.New("INFO")
	lwL.Info("interesting info")
}
```

## License and Copyright

Copyright 2020, Lacework Inc.

```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
