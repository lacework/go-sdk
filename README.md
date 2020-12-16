<img src="https://techally-content.s3-us-west-1.amazonaws.com/public-content/lacework_logo_full.png" width="600">

# Lacework Go SDK

[![GitHub release](https://img.shields.io/github/release/lacework/go-sdk.svg)](https://github.com/lacework/go-sdk/releases/)
[![Go version](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/lacework/go-sdk)
[![Go report](https://goreportcard.com/badge/github.com/lacework/go-sdk)](https://goreportcard.com/report/github.com/lacework/go-sdk)
[![CircleCI status](https://circleci.com/gh/lacework/go-sdk.svg?style=shield)](https://circleci.com/gh/lacework/go-sdk)
[![GitHub releases](https://img.shields.io/github/downloads/lacework/go-sdk/total.svg)](https://GitHub.com/lacework/go-sdk/releases/)

This repository provides a set of tools, libraries, relevant documentation, code
samples, processes, and/or guides that allow users and developers to interact with
the Lacework platform.

Find more information about this repository at the following [Wiki page](https://github.com/lacework/go-sdk/wiki).

## Lacework CLI ([`cli`](cli/))

The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. You can use it to manage compliance
reports, external integrations, vulnerability scans, and other operations.

### Install

#### Bash:
```
$ curl https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.sh | bash
```

#### Powershell:
```
C:\> Set-ExecutionPolicy Bypass -Scope Process -Force
C:\> iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.ps1'))
```

Look at the [cli/](cli/) folder for more information.

## Lacework API Client ([`api`](api/))

A Golang API client for interacting with the [Lacework API](https://support.lacework.com/hc/en-us/categories/360002496114-Lacework-API-).

### Basic Usage
```go
package main

import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"
)

func main() {
	lacework, err := api.NewClient("account")
	if err != nil {
		log.Fatal(err)
	}

	tokenRes, err := lacework.GenerateTokenWithKeys("KEY", "SECRET")
	if err != nil {
		log.Fatal(err)
	}

	// Output: YOUR-ACCESS-TOKEN
	fmt.Println(tokenRes.Token())
}
```
Look at the [api/](api/) folder for more information.

## Lacework Logger ([`lwlogger`](lwlogger/))

A Logger wrapper for Lacework based of [zap](https://github.com/uber-go/zap) logger Go package.

### Basic Usage
```go
package main

import "github.com/lacework/go-sdk/lwlogger"

func main() {
	lwL := lwlogger.New("INFO")

	// Output: {"level":"info","ts":"[timestamp]","caller":"main.go:9","msg":"interesting info"}
	lwL.Info("interesting info")
}
```

Look at the [lwlogger/](lwlogger/) folder for more information.

## Lacework Updater ([`lwupdater`](lwupdater/))

A Go library to check for available updates of Lacework projects.

### Basic Usage
```go
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
			project, sdk.LatestVersion,
		)
	}
}
```

Set the environment variable `LW_UPDATES_DISABLE=1` to avoid checking for updates.

## Release Process

The release process of this repository is documented at the following [Wiki page](https://github.com/lacework/go-sdk/wiki/Release-Process).

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
