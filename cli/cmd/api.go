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

package cmd

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/internal/array"
)

var (
	// list of valid API methods
	validApiMethods = []string{"get", "post", "delete", "patch"}

	// data to send for POST/PATCH request
	apiData string

	// apiCmd represents the api command
	apiCmd = &cobra.Command{
		Use:   "api <method> <path>",
		Short: "Helper to call Lacework's API",
		Long: `Use this command as a helper to call any available Lacework API v1 & v2 endpoint.


### APIv2

To list all available Lacework schema type:

    lacework api get /v2/schemas

Example usage,
	To recieve a json response of all machines within the given time window
	lacework api post /api/v2/Entities/Machines/search -d "{\"timeFilter\":{\"startTime\":\"2023-05-10T00:00:00Z\",\"endTime\":\"2023-05-14T00:00:00Z\"}}"   
	
	To recieve a json response of all agents within the given time window
	lacework api post /api/v2/AgentInfo/search -d "{\"timeFilter\":{\"startTime\":\"2023-05-10T00:00:00Z\",\"endTime\":\"2023-05-14T00:00:00Z\"}}"  

For a complete list of available API v2 endpoints visit:

    https://<ACCOUNT>.lacework.net/api/v2/docs
`,
		Args: argsApiValidator,
		RunE: runApiCommand,
	}
)

func init() {
	// add the api command
	rootCmd.AddCommand(apiCmd)

	apiCmd.Flags().StringVarP(&apiData,
		"data", "d", "",
		"data to send only for post and patch requests",
	)
}

func runApiCommand(_ *cobra.Command, args []string) error {
	switch args[0] {
	case "patch":
		if apiData == "" {
			return fmt.Errorf("missing '--data' parameter patch requests")
		}
	case "get":
		if apiData != "" {
			return fmt.Errorf("use '--data' only for post, delete and patch requests")
		}
	}

	response := new(interface{})
	err := cli.LwApi.RequestDecoder(
		strings.ToUpper(args[0]),
		cleanupEndpoint(args[1]),
		strings.NewReader(apiData),
		response,
	)
	if err != nil {
		return errors.Wrap(err, "unable to send the request")
	}

	if err := cli.OutputJSON(*response); err != nil {
		return errors.Wrap(err, "unable to format json response")
	}
	return nil
}

func argsApiValidator(_ *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("requires 2 argument. (method and path)")
	}
	if !array.ContainsStr(validApiMethods, args[0]) {
		return fmt.Errorf(
			"invalid method specified: '%s' (valid methods are %s)",
			args[0], validApiMethods,
		)
	}
	return nil
}

// cleanupEndpoint will make sure that any provided endpoint is well formatted
// and doesn't contain known fields like /api/v1/foo
func cleanupEndpoint(endpoint string) string {
	splitEndpoint := strings.Split(endpoint, "/")

	if len(splitEndpoint) > 0 && splitEndpoint[0] == "api" {
		return strings.Join(splitEndpoint[1:], "/")
	}

	if len(splitEndpoint) > 1 && splitEndpoint[1] == "api" {
		return strings.Join(splitEndpoint[2:], "/")
	}

	return strings.TrimPrefix(endpoint, "/")
}
