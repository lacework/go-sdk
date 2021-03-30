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
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var (
	lqlCmdState = struct {
		End          string
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
		Use:   "run [query_id]",
		Short: "run an LQL query",
		Long: `Run an LQL query.

Run a query via text:

	$ lacework query run

Run a query via ID (uses active profile):

	$ lacework query run MyQuery

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
	// if a query_id was specified
	if len(args) != 0 && args[0] != "" {
		return inputQueryFromEnv(args[0])
	}
	// if running via repo
	if lqlCmdState.Repo {
		return inputQueryFromRepo()
	}
	// if running via file
	if lqlCmdState.File != "" {
		return inputQueryFromFile(lqlCmdState.File)
	}
	// if running via URL
	if lqlCmdState.URL != "" {
		return inputQueryFromURL(lqlCmdState.URL)
	}
	// if running via editor
	action := "validate"
	if !lqlCmdState.ValidateOnly {
		action = strings.Split(cmd.Use, " ")[0]
	}
	return inputQueryFromEditor(action)
}

func inputQueryFromEnv(queryID string) (
	query string,
	err error,
) {
	var queryResponse api.LQLQueryResponse

	queryResponse, err = cli.LwApi.LQL.GetQueryByID(queryID)
	if err == nil && len(queryResponse.Data) != 0 {
		query = queryResponse.Data[0].QueryText
	}
	return
}

func inputQueryFromRepo() (
	query string,
	err error,
) {
	err = errors.New("NotImplementedError")
	return
}

func inputQueryFromFile(filePath string) (
	query string,
	err error,
) {
	var fileData []byte
	fileData, err = ioutil.ReadFile(filePath)

	if err != nil {
		err = errors.Wrap(err, "unable to read file")
		return
	}

	query = string(fileData)
	return
}

func inputQueryFromURL(url string) (
	query string,
	err error,
) {
	msg := "unable to open URL"
	var response *http.Response
	var body []byte

	response, err = http.Get(url)
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
	return
}

func inputQueryFromEditor(action string) (
	query string,
	err error,
) {
	prompt := &survey.Editor{
		Message:  fmt.Sprintf("Type a query to %s", action),
		FileName: "query*.sh",
	}
	err = survey.AskOne(prompt, &query)

	return
}

func runQuery(cmd *cobra.Command, args []string) error {
	msg := "unable to run LQL query"
	var response map[string]interface{}

	query, err := inputQuery(cmd, args)
	if err != nil {
		return errors.Wrap(err, msg)
	}

	cli.Log.Debugw("running LQL query", "query", query)

	// validate_only should compile
	if lqlCmdState.ValidateOnly {
		return CompileQueryAndOutput(query)
	}
	// !validate_only should should run
	response, err = cli.LwApi.LQL.RunQuery(query, lqlCmdState.Start, lqlCmdState.End)

	if err != nil {
		return errors.Wrap(err, msg)
	}
	if data, ok := response["data"]; ok {
		return cli.OutputJSON(data)
	}
	if err := cli.OutputJSON(response); err != nil {
		return errors.Wrap(err, "unable to format json response")
	}
	return nil
}
