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
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	lqlEnd      string
	lqlEnv      bool
	lqlFile     string
	lqlRepo     bool
	lqlStart    string
	lqlURL      string
	lqlValidate bool

	// lqlCmd represents the lql parent command
	lqlCmd = &cobra.Command{
		Aliases: []string{"lql"},
		Use:     "query",
		Short:   "Run and manage LQL queries",
		Long: `Run and manage LQL queries.

NOTE: This feature is not yet available!`,
	}

	// lqlRunCmd represents the lql run command
	lqlRunCmd = &cobra.Command{
		Use:   "run <query|queryID>",
		Short: "run an LQL query",
		Long: `Run an LQL query.

Run a query via text:

	$ lacework lql run 'SimpleLQL_3(CloudTrailRawEvents e) {SELECT INSERT_ID}' --start <start> --end <end>

Run a query via ID:

	$ lacework lql run MyQuery -e

Start and End times are required to run a query:

1.  Start and End times must be specified in one of the following formats:

	A. ISO 8601 Date and Time
	B. Epoch time in milliseconds

2. Start and End times must be specified in one of the following ways:

	A.  As StartTimeRange and EndTimeRange in the ParamInfo block within the LQL query
	B.  As START_TIME_RANGE and END_TIME_RANGE if specifying JSON
	C.  As --start and --end CLI flags
	
3. Start and End time precedence:

	A.  CLI flags take precedence over JSON specifications
	B.  JSON specifications take precedence over ParamInfo specifications`,
		Args: cobra.MaximumNArgs(1),
		RunE: runQuery,
	}

	// create a slice of the cobra.Command pointers
	// which need query "text" as input
	lqlQueryCommands []*cobra.Command = []*cobra.Command{
		lqlCreateCmd,
		lqlRunCmd,
		lqlUpdateCmd,
		lqlValidateCmd,
	}
)

func init() {
	// add the lql command
	rootCmd.AddCommand(lqlCmd)

	// add sub-commands to the lql command
	lqlCmd.AddCommand(lqlRunCmd)

	// for commands that take query "text" as input
	for _, cmd := range lqlQueryCommands {
		// file flag to specify a query from disk
		cmd.Flags().StringVarP(
			&lqlFile,
			"file", "f", "",
			"path to an LQL query to run",
		)
		// repo flag to specify a query from repo
		cmd.Flags().BoolVarP(
			&lqlRepo,
			"repo", "r", false,
			"id of an LQL query to run via active repo",
		)
		// url flag to specify a query from url
		cmd.Flags().StringVarP(
			&lqlURL,
			"url", "u", "",
			"url to an LQL query to run",
		)
	}

	// run specific flags
	// env flag to specify a query from disk
	lqlRunCmd.Flags().BoolVarP(
		&lqlEnv,
		"env", "e", false,
		"run an LQL query by ID (using active profile)",
	)
	// start time flag
	// TODO: come up with reasonable default per UI (1d)
	lqlRunCmd.Flags().StringVarP(
		&lqlStart,
		"start", "", "",
		"start time for LQL query",
	)
	// end time flag
	// TODO: come up with reasonable default per UI (1d)
	lqlRunCmd.Flags().StringVarP(
		&lqlEnd,
		"end", "", "",
		"end time for LQL query",
	)
	lqlRunCmd.Flags().BoolVarP(
		&lqlValidate,
		"validate_only", "", false,
		"validate query only (do not run)",
	)
}

// for commands that take a query as input
func inputQuery(cmd *cobra.Command, args []string) (
	query string,
	err error,
) {
	var queryID string

	// if an inline argument was provided
	// determine if it's a query or a query identifier
	if len(args) != 0 && args[0] != "" {
		if lqlEnv || lqlRepo {
			queryID = args[0]
		} else {
			query = args[0]
		}
	}

	if lqlEnv {
		var queryResponse api.LQLQueryResponse
		queryResponse, err = cli.LwApi.LQL.GetQueryByID(queryID)
		if err == nil && len(queryResponse.Data) != 0 {
			query = queryResponse.Data[0].QueryText
		}
	} else if lqlRepo {
		err = errors.New("NotImplementedError")
	} else if lqlFile != "" {
		var fileData []byte
		fileData, err = ioutil.ReadFile(lqlFile)
		if err != nil {
			err = errors.Wrap(err, "unable to read file")
			return
		}
		query = string(fileData)
	} else if lqlURL != "" {
		msg := "unable to open URL"
		var response *http.Response
		var body []byte

		response, err = http.Get(lqlURL)
		if err != nil {
			err = errors.Wrap(err, msg)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			err = errors.Wrap(errors.New(response.Status), msg)
			return
		}

		body, err = io.ReadAll(response.Body)
		if err != nil {
			err = errors.Wrap(err, msg)
			return
		}
		query = string(body)
	} else {
		var firstUseWord string
		if lqlValidate {
			firstUseWord = "validate"
		} else {
			firstUseWord = strings.Split(cmd.Use, " ")[0]
		}
		prompt := &survey.Editor{
			Message:  fmt.Sprintf("Type a query to %s", firstUseWord),
			FileName: "query*.sh",
		}
		err = survey.AskOne(prompt, &query)
	}

	return
}

// standardized cli/error output
func output(response map[string]interface{}, err error, msg string) error {
	if err != nil {
		return errors.Wrap(err, msg)
	}

	if data, ok := response["data"]; ok {
		err := cli.OutputJSON(data)
		return err
	}

	if err := cli.OutputJSON(response); err != nil {
		return errors.Wrap(err, "unable to format json response")
	}

	return nil
}

func runQuery(cmd *cobra.Command, args []string) error {
	msg := "unable to run LQL query"
	var response map[string]interface{}

	query, err := inputQuery(cmd, args)
	if err != nil {
		return output(response, err, msg)
	}

	cli.Log.Debugw("running LQL query", "query", query)

	if lqlValidate {
		// validate_only should compile
		return CompileQueryAndOutput(query)
	} else {
		// !validate_only should should run
		response, err = cli.LwApi.LQL.RunQuery(query, lqlStart, lqlEnd)
	}

	return output(response, err, msg)
}
