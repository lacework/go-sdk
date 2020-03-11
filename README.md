# Lacework Go SDK

This repository provides a Go API client, tools, libraries, relevant documentation, code
samples, processes, and/or guides that allow developers to interact with Lacework.

## API Client ([`api`](api/))

A Golang API client for interacting with the [Lacework API](https://support.lacework.com/hc/en-us/categories/360002496114-Lacework-API-).

### Basic Usage
```go
import "github.com/lacework/go-sdk/api"

lacework, err := api.NewClient("account")
if err == nil {
	log.Fatal(err)
}

tokenRes, err := lacework.GenerateTokenWithKeys("KEY", "SECRET")
if err != nil {
	log.Fatal(err)
}

// Output: YOUR-ACCESS-TOKEN
fmt.Printf("%s\n", tokenRes.Token())
```
Look at the [api/](api/) folder for more documentation.

## Lacework CLI ([`cli`](cli/))

_(work-in-progress)_ The Lacework Command Line Interface.

### Basic Usage

Today, you have to first build the CLI by running `make build-cli`, then you will be
able to execute it directly:
```
$ make build-cli
$ ./bin/lacework-cli
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
