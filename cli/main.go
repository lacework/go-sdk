//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lacework/go-sdk/api"
)

// (@afiune) This will become the Lacework CLI at some point in time:
//
//   $ lacework-cli api get integrations
func main() {
	keysFile := flag.String("api-keys", "", "JSON file containing a set of Lacework API keys")
	flag.Parse()

	lacework, err := api.NewClient("customerdemo")
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
		exitWithError("unable to read keys file", err)
	}

	err = json.Unmarshal(content, &keys)
	if err != nil {
		exitWithError("unable to parse keys file", err)
	}

	fmt.Printf("Generation API token: ")
	token, err := lacework.GenerateTokenWithKeys(keys.KeyId, keys.Secret)
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

// TODO: @afiune should we backport this apiKeys struct to the api package
type apiKeys struct {
	KeyId  string
	Secret string
}
