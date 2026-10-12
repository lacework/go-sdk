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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/internal/array"
)

var (
	// list of valid API methods
	validApiMethods = []string{"get", "post", "delete", "patch"}

	// data to send for POST/PATCH request
	apiData string

	// file to read the data from, "-" reads it from stdin
	apiDataFile string

	// apiCmd represents the api command
	apiCmd = &cobra.Command{
		Use:   "api <method> <path>",
		Short: "Helper to call Lacework's API",
		Long: `Use this command as a helper to call any available Lacework API v2 endpoint.

### API v2

To list all available Lacework schema types:

    lacework api get /v2/schemas

To receive a json response of all machines within the given time window:

    lacework api post /api/v2/Entities/Machines/search -d "{}"

To receive a json response of all agents within the given time window:

    lacework api post /api/v2/AgentInfo/search -d "{}"

To send the request body from a file (or from stdin with "-f -"):

    lacework api patch /api/v2/Policies -f policies.json
    lacework api patch /api/v2/Policies -d '@policies.json'

Use "--file" in PowerShell 7.2 or earlier. These versions remove the double
quotes inside "--data", so the API receives invalid JSON.

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
		"data to send only for post and patch requests, '@<file>' reads it from a file",
	)
	apiCmd.Flags().StringVarP(&apiDataFile,
		"file", "f", "",
		"file with the data to send, use '-' to read from stdin",
	)
	apiCmd.MarkFlagsMutuallyExclusive("data", "file")
}

func runApiCommand(_ *cobra.Command, args []string) error {
	path := apiDataFile
	if strings.HasPrefix(apiData, "@") {
		path = strings.TrimPrefix(apiData, "@")
		if path == "" {
			return errors.New("missing file name after '@' in '--data'")
		}
	}
	if path != "" {
		data, err := readApiDataFile(path)
		if err != nil {
			return err
		}
		apiData = string(data)
	}

	switch args[0] {
	case "patch":
		if apiData == "" {
			return fmt.Errorf("missing '--data' or '--file' parameter for patch requests")
		}
	case "get":
		if apiData != "" {
			return fmt.Errorf("use '--data' or '--file' only for post, delete and patch requests")
		}
	}

	if apiData != "" && !json.Valid([]byte(apiData)) {
		return errors.New("the request data is not valid JSON. If your shell removes " +
			"the double quotes (PowerShell 7.2 or earlier does), put the JSON in a file and use '--file'")
	}

	response := new(interface{})
	err := cli.LwApi.RequestDecoder(
		strings.ToUpper(args[0]),
		cleanupEndpoint(args[1]),
		strings.NewReader(apiData),
		response,
	)
	if err != nil {
		return errors.Wrap(err, "an error occurred with the request")
	}

	if err := cli.OutputJSON(*response); err != nil {
		return errors.Wrap(err, "unable to format json response")
	}
	return nil
}

func readApiDataFile(path string) ([]byte, error) {
	var (
		data []byte
		err  error
	)
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to read data file")
	}
	// Windows PowerShell 5.1 writes UTF-16 with '>' and Notepad can add a UTF-8 BOM
	if bytes.HasPrefix(data, []byte{0xFF, 0xFE}) || bytes.HasPrefix(data, []byte{0xFE, 0xFF}) {
		return nil, errors.New("the data file is UTF-16 encoded, save it as UTF-8")
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	if len(data) == 0 {
		return nil, errors.New("the data file is empty")
	}
	return data, nil
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
