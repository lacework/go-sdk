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

	"github.com/lacework/go-sdk/api"
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
		Short: "Helper to call Lacework's RestfulAPI",
		Long: `Use this command as a helper to call any available Lacework API endpoint.

For example, to list all integrations configured in your account run:

    lacework api get /external/integrations

For a complete list of available API endpoints visit:

    https://<ACCOUNT>.lacework.net/api/v1/external/docs
`,
		Args: argsApiValidator,
		RunE: runApiCommand,
	}
)

func init() {
	rootCmd.AddCommand(apiCmd)

	apiCmd.Flags().StringVarP(&apiData,
		"data", "d", "",
		"data to send only for post and patch requests",
	)
}

func runApiCommand(_ *cobra.Command, args []string) error {
	switch args[0] {
	case "post", "patch":
		if apiData == "" {
			return fmt.Errorf("missing '--data' parameter for post or patch requests")
		}
	case "delete", "get":
		if apiData != "" {
			return fmt.Errorf("use '--data' only for post and patch requests")
		}
	}

	lacework, err := api.NewClient(cli.Account,
		api.WithLogLevel(cli.LogLevel),
		api.WithApiKeys(cli.KeyID, cli.Secret),
	)
	if err != nil {
		return errors.Wrap(err, "unable to generate Lacework api client")
	}

	response := new(map[string]interface{})
	err = lacework.RequestDecoder(
		strings.ToUpper(args[0]),
		strings.TrimPrefix(args[1], "/"),
		strings.NewReader(apiData),
		response,
	)
	if err != nil {
		return errors.Wrap(err, "unable to send the request")
	}

	pretty, err := cli.JsonF.Marshal(*response)
	if err != nil {
		cli.Log.Debugw("unable to pretty print JSON response", "raw", response)
		return errors.Wrap(err, "unable to format json response")
	}

	fmt.Println(string(pretty))
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
