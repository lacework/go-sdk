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
