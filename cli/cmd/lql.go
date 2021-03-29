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
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/lacework/go-sdk/api"
)

var (
	lqlCmdState = struct {
		End          string
		Env          bool
		File         string
		Repo         bool
		Start        string
		URL          string
		ValidateOnly bool
	}{}

	// lqlCmd represents the lql parent command
	lqlCmd = &cobra.Command{
		Aliases: []string{"lql"},
		Use:     "query",
		Short:   "run and manage LQL queries",
		Long: `Run and manage LQL queries.

NOTE: This feature is not yet available!`,
	}

	// lqlRunCmd represents the lql run command
	lqlRunCmd = &cobra.Command{
		Use:   "run [query|query_id]",
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
)

func init() {
	// add the lql command
	rootCmd.AddCommand(lqlCmd)

	// add sub-commands to the lql command
	lqlCmd.AddCommand(lqlRunCmd)

	// run specific flags
	setQueryFlags(lqlRunCmd.Flags())

	// env flag to specify a query from disk
	lqlRunCmd.Flags().BoolVarP(
		&lqlCmdState.Env,
		"env", "e", false,
		"run an LQL query by ID (using active profile)",
	)
	// start time flag
	// TODO: come up with reasonable default per UI (1d)
	lqlRunCmd.Flags().StringVarP(
		&lqlCmdState.Start,
		"start", "", "",
		"start time for LQL query",
	)
	// end time flag
	// TODO: come up with reasonable default per UI (1d)
	lqlRunCmd.Flags().StringVarP(
		&lqlCmdState.End,
		"end", "", "",
		"end time for LQL query",
	)
	lqlRunCmd.Flags().BoolVarP(
		&lqlCmdState.ValidateOnly,
		"validate_only", "", false,
		"validate query only (do not run)",
	)
}

func setQueryFlags(cmds ...*flag.FlagSet) {
	for _, cmd := range cmds {
		if cmd != nil {
			// file flag to specify a query from disk
			cmd.StringVarP(
				&lqlCmdState.File,
				"file", "f", "",
				"path to an LQL query to run",
			)
			// repo flag to specify a query from repo
			cmd.BoolVarP(
				&lqlCmdState.Repo,
				"repo", "r", false,
				"id of an LQL query to run via active repo",
			)
			// url flag to specify a query from url
			cmd.StringVarP(
				&lqlCmdState.URL,
				"url", "u", "",
				"url to an LQL query to run",
			)
		}
	}
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
		if lqlCmdState.Env || lqlCmdState.Repo {
			queryID = args[0]
		} else {
			query = args[0]
		}
	}

	if lqlCmdState.Env {
		var queryResponse api.LQLQueryResponse
		queryResponse, err = cli.LwApi.LQL.GetQueryByID(queryID)
		if err == nil && len(queryResponse.Data) != 0 {
			query = queryResponse.Data[0].QueryText
		}
	} else if lqlCmdState.Repo {
		err = errors.New("NotImplementedError")
	} else if lqlCmdState.File != "" {
		var fileData []byte
		fileData, err = ioutil.ReadFile(lqlCmdState.File)
		if err != nil {
			err = errors.Wrap(err, "unable to read file")
			return
		}
		query = string(fileData)
	} else if lqlCmdState.URL != "" {
		msg := "unable to open URL"
		var response *http.Response
		var body []byte

		response, err = http.Get(lqlCmdState.URL)
		if err != nil {
			err = errors.Wrap(err, msg)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			err = errors.Wrap(errors.New(response.Status), msg)
			return
		}

		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			err = errors.Wrap(err, msg)
			return
		}
		query = string(body)
	} else {
		var firstUseWord string
		if lqlCmdState.ValidateOnly {
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

	if lqlCmdState.ValidateOnly {
		// validate_only should compile
		return CompileQueryAndOutput(query)
	} else {
		// !validate_only should should run
		response, err = cli.LwApi.LQL.RunQuery(query, lqlCmdState.Start, lqlCmdState.End)
	}

	return output(response, err, msg)
}
