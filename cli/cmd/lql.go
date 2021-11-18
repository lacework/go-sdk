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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwtime"
)

var (
	queryCmdState = struct {
		End          string
		File         string
		Repo         bool
		Range        string
		Start        string
		URL          string
		ValidateOnly bool
	}{}

	// queryCmd represents the lql parent command
	queryCmd = &cobra.Command{
		Use:     "query",
		Aliases: []string{"lql", "queries"},
		Short:   "run and manage queries",
		Long: `Run and manage Lacework Query Language (LQL) queries.

To provide customizable specification of datasets, Lacework provides the Lacework
Query Language (LQL). LQL is a human-readable text syntax for specifying selection,
filtering, and manipulation of data.

Currently, Lacework has introduced LQL for configuration of AWS CloudTrail policies
and queries. This means you can use LQL to customize AWS CloudTrail policies only.
For all other policies, use the previous existing methods.

Lacework ships a set of default LQL queries that are available in your account.

For more information about LQL, visit:

  https://docs.lacework.com/lql-overview

To view all LQL queries in your Lacework account.

    lacework query ls

To show a query.

    lacework query show <query_id>

To execute a query.

    lacework query run <query_id>

** NOTE: LQL syntax may change. **
`,
	}

	// queryRunCmd represents the lql run command
	queryRunCmd = &cobra.Command{
		Aliases: []string{"execute"},
		Use:     "run [query_id]",
		Short:   "run a query",
		Long: `Run an LQL query via editor:

    lacework query run --range today

Run a query via ID (uses active profile):

    lacework query run MyQuery --start "-1w@w" --end "@w"

Start and End times are required to run a query:

1.  Start and End times must be specified in one of the following formats:

    A. A relative time specifier
    B. RFC3339 Date and Time
    C. Epoch time in milliseconds

2. Start and End times must be specified in one of the following ways:

    A.  As StartTimeRange and EndTimeRange in the ParamInfo block within the query
    B.  As start_time_range and end_time_range if specifying JSON
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
	rootCmd.AddCommand(queryCmd)

	// add sub-commands to the lql command
	queryCmd.AddCommand(queryRunCmd)

	// run specific flags
	setQuerySourceFlags(queryRunCmd)

	// since time flag
	queryRunCmd.Flags().StringVarP(
		&queryCmdState.Range,
		"range", "", "",
		"natural time range for query",
	)

	// start time flag
	queryRunCmd.Flags().StringVarP(
		&queryCmdState.Start,
		"start", "", "@d",
		"start time for query",
	)
	// end time flag
	queryRunCmd.Flags().StringVarP(
		&queryCmdState.End,
		"end", "", "now",
		"end time for query",
	)
	queryRunCmd.Flags().BoolVarP(
		&queryCmdState.ValidateOnly,
		"validate_only", "", false,
		"validate query only (do not run)",
	)
}

func setQuerySourceFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		if cmd != nil {
			action := strings.Split(cmd.Use, " ")[0]

			// file flag to specify a query from disk
			cmd.Flags().StringVarP(
				&queryCmdState.File,
				"file", "f", "",
				fmt.Sprintf("path to a query to %s", action),
			)
			/* repo flag to specify a query from repo
			cmd.Flags().BoolVarP(
				&queryCmdState.Repo,
				"repo", "r", false,
				fmt.Sprintf("id of a query to %s via active repo", action),
			)*/
			// url flag to specify a query from url
			cmd.Flags().StringVarP(
				&queryCmdState.URL,
				"url", "u", "",
				fmt.Sprintf("url to a query to %s", action),
			)
		}
	}
}

// for commands that take a query as input
func inputQuery(cmd *cobra.Command) (string, error) {
	// if running via repo
	if queryCmdState.Repo {
		return inputQueryFromRepo()
	}
	// if running via file
	if queryCmdState.File != "" {
		return inputQueryFromFile(queryCmdState.File)
	}
	// if running via URL
	if queryCmdState.URL != "" {
		return inputQueryFromURL(queryCmdState.URL)
	}
	// if running via stdin
	stat, err := os.Stdin.Stat()
	if err != nil {
		cli.Log.Debugw("error retrieving stdin mode", "error", err.Error())
	} else if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		return string(bytes), err
	}
	// if running via editor
	action := "validate"
	if !queryCmdState.ValidateOnly {
		action = strings.Split(cmd.Use, " ")[0]
	}
	return inputQueryFromEditor(action)
}

func inputQueryFromRepo() (query string, err error) {
	err = errors.New("NotImplementedError")
	return
}

func inputQueryFromFile(filePath string) (string, error) {
	fileData, err := ioutil.ReadFile(filePath)

	if err != nil {
		return "", errors.Wrap(err, "unable to read file")
	}

	return string(fileData), nil
}

func inputQueryFromURL(url string) (query string, err error) {
	msg := "unable to access URL"

	response, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, msg)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.Wrap(errors.New(response.Status), msg)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.Wrap(err, msg)
		return
	}
	query = string(body)
	return
}

func inputQueryFromEditor(action string) (query string, err error) {
	prompt := &survey.Editor{
		Message:  fmt.Sprintf("Type a query to %s", action),
		FileName: "query*.lql",
	}
	err = survey.AskOne(prompt, &query)

	return
}

func parseQuery(s string) (api.NewQuery, error) {
	var query api.NewQuery
	var err error

	// valid json
	if err = json.Unmarshal([]byte(s), &query); err == nil {
		return query, err
	}
	// valid yaml
	query = api.NewQuery{}
	err = yaml.Unmarshal([]byte(s), &query)
	if err == nil && !reflect.DeepEqual(query, api.NewQuery{}) { // empty string unmarshals w/o error
		return query, nil
	}
	// invalid policy
	return query, queryErrorCrumbs(s)
}

func parseQueryTime(s string) (time.Time, error) {
	// empty
	if s == "" {
		return time.Time{}, errors.New(fmt.Sprintf("unable to parse time (%s)", s))
	}
	// parse time as relative
	if t, err := lwtime.ParseRelative(s); err == nil {
		return t, err
	}
	// parse time as RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, err
	}
	// parse time as millis
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(0, i*int64(time.Millisecond)), err
	}
	return time.Time{}, errors.New(fmt.Sprintf("unable to parse time (%s)", s))
}

func queryErrorCrumbs(q string) error {
	// smells like json
	q = strings.TrimSpace(q)
	if strings.HasPrefix(q, "[") ||
		strings.HasPrefix(q, "{") {

		return errors.New(`invalid query

It looks like you attempted to submit a query in JSON format.
Please validate that the JSON is formatted properly and adheres to the following schema:

{
    "evaluatorId": "Cloudtrail",
    "queryId": "MyLQL",
    "queryText": "MyLQL { source { CloudTrailRawEvents } filter { EVENT_SOURCE = 's3.amazonaws.com' } return { INSERT_ID } }"
}
`)
	}
	// smells like plain text
	return errors.New(`invalid query
	
It looks like you attempted to submit a query in YAML format.
Please validate that the text adheres to the following schema:

evaluatorId: Cloudtrail
queryId: MyLQL
queryText: |-
  MyLQL {
      source {
          CloudTrailRawEvents
      }
      filter {
          EVENT_SOURCE = 's3.amazonaws.com'
      }
      return {
          INSERT_ID
      }
  }
`)
}

func runQuery(cmd *cobra.Command, args []string) error {
	var (
		err        error
		start      time.Time
		end        time.Time
		msg        string = "unable to run query"
		hasCmdArgs bool   = len(args) != 0 && args[0] != ""
	)

	// validate_only w/ query_id
	if queryCmdState.ValidateOnly && hasCmdArgs {
		return errors.New("flag --validate_only unavailable when specifying query_id argument")
	}
	// validate_only
	if queryCmdState.ValidateOnly {
		return validateQuery(cmd, args)
	}

	// use of if/else intentional here based on logic paths for determining start and end time.Time values
	// if cli user has specified a range we use ParseNatural which gives us start and end time.Time values
	// otherwise we need to convert queryCmdState start and end strings to time.Time values using parseQueryTime
	if queryCmdState.Range != "" {
		cli.Log.Debugw("retrieving natural time range")

		start, end, err = lwtime.ParseNatural(queryCmdState.Range)
		if err != nil {
			return errors.Wrap(err, msg)
		}
	} else {
		// parse start
		start, err = parseQueryTime(queryCmdState.Start)
		if err != nil {
			return errors.Wrap(err, msg)
		}
		// parse end
		end, err = parseQueryTime(queryCmdState.End)
		if err != nil {
			return errors.Wrap(err, msg)
		}
	}

	queryArgs := []api.ExecuteQueryArgument{
		api.ExecuteQueryArgument{
			Name:  "StartTimeRange",
			Value: start.UTC().Format(lwtime.RFC3339Milli),
		},
		api.ExecuteQueryArgument{
			Name:  "EndTimeRange",
			Value: end.UTC().Format(lwtime.RFC3339Milli),
		},
	}

	// query by id
	if hasCmdArgs {
		return outputQueryRunResponse(
			runQueryByID(args[0], queryArgs),
		)
	}
	// adhoc query
	return outputQueryRunResponse(
		runAdhocQuery(cmd, queryArgs),
	)
}

func runQueryByID(id string, args []api.ExecuteQueryArgument) (
	map[string]interface{},
	error,
) {
	cli.Log.Debugw("running query", "query", id)

	cli.StartProgress(" Executing query...")
	defer cli.StopProgress()

	request := api.ExecuteQueryByIDRequest{
		QueryID:   id,
		Arguments: args,
	}
	return cli.LwApi.V2.Query.ExecuteByID(request)
}

func runAdhocQuery(cmd *cobra.Command, args []api.ExecuteQueryArgument) (
	response map[string]interface{},
	err error,
) {
	// input query
	queryString, err := inputQuery(cmd)
	if err != nil {
		return
	}
	// parse query
	newQuery, err := parseQuery(queryString)
	if err != nil {
		return
	}

	cli.StartProgress(" Executing query...")
	defer cli.StopProgress()

	// execute query
	executeQuery := api.ExecuteQueryRequest{
		Query: api.ExecuteQuery{
			QueryText:   newQuery.QueryText,
			EvaluatorID: newQuery.EvaluatorID,
		},
		Arguments: args,
	}

	cli.Log.Debugw("running query", "query", queryString)
	response, err = cli.LwApi.V2.Query.Execute(executeQuery)
	return
}

func outputQueryRunResponse(response map[string]interface{}, err error) error {
	msg := "unable to run query"

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
