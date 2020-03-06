package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lacework/go-sdk/client"
)

// (@afiune) This will become the Lacework CLI at some point in time:
//
//   $ lacework-cli api get integrations
func main() {
	keysFile := flag.String("api-keys", "", "JSON file containing a set of Lacework API keys")
	flag.Parse()

	lacework, err := client.New("customerdemo")
	if err != nil {
		exitWithError("unable to generate api client", err)
	}
	fmt.Printf("Api version: %s\n", lacework.ApiVersion())

	if len(*keysFile) == 0 {
		fmt.Println("\nTry passing '-api-keys [file.json]'")
		os.Exit(0)
	}

	var keys apiKeys
	content, err := ioutil.ReadFile(*keysFile)
	if err != nil {
		exitWithError("unable to generate api client", err)
	}
	json.Unmarshal(content, &keys)

	fmt.Printf("Generation API token: ")
	token, err := lacework.GenerateToken(keys.KeyId, keys.Secret)
	if err != nil {
		exitWithError("unable to generate token", err)
	}

	fmt.Println(token.Message)

	fmt.Printf("List all integrations: ")
	integrations, err := lacework.GetIntegrations()
	if err != nil {
		exitWithError("unable to generate integrations", err)
	}
	fmt.Println(token.Message)
	fmt.Println("---------------------------------")

	fmt.Println(integrations.List())
}

func exitWithError(msg string, err error) {
	fmt.Println("\nERROR: " + msg)
	fmt.Println(err)
	os.Exit(1)
}

// TODO: @afiune should we backport this apiKeys struct to the client package
type apiKeys struct {
	KeyId  string
	Secret string
}
