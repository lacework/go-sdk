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
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/failon"
	"github.com/lacework/go-sdk/lwtime"
)

var (
	queryCmdState = struct {
		End           string
		File          string
		Limit         int
		Range         string
		Start         string
		URL           string
		ValidateOnly  bool
		FailOnCount   string
		EmptyTemplate bool
		// create, update validate from library
		CURVFromLibrary string
	}{}

	// queryCmd represents the lql parent command
	queryCmd = &cobra.Command{
		Use:     "query",
		Aliases: []string{"lql", "queries"},
		Short:   "Run and manage queries",
		Long: `Run and manage Lacework Query Language (LQL) queries.

LQL is a SQL-like query language for specifying the selection, filtering, and 
manipulation of data. Queries let you interactively request information from 
specified curated datasources.

Lacework ships a set of default LQL queries that are available in your account.

For more information about LQL, visit:

  https://docs.lacework.com/lql-overview

To view all LQL queries in your Lacework account.

    lacework query ls

To show a query.

    lacework query show <query_id>

To execute a query.

    lacework query run <query_id>

**Note: LQL syntax may change.**
`,
	}

	// queryRunCmd represents the lql run command
	queryRunCmd = &cobra.Command{
		Aliases: []string{"execute"},
		Use:     "run [query_id]",
		Short:   "Run a query",
		Long: `Run an LQL query via editor:

    lacework query run --range today

Run a query via ID (uses active profile):

    lacework query run MyQuery --start "-1w@w" --end "@w"

Start and end times are required to run a query:

1.  Specify start and end times in one of the following formats:

    A. A relative time specifier
    B. RFC3339 date and time
    C. Epoch time in milliseconds

2. Specify start and end times in one of the following ways:

    A. As StartTimeRange and EndTimeRange in the ParamInfo block within the query
    B. As start_time_range and end_time_range if specifying JSON
    C. As --start and --end CLI flags

3. Start and End time precedence:

    A. CLI flags take precedence over JSON specifications`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			// default is 0 hence the '< 0' comparison
			if queryCmdState.Limit < 0 {
				return errors.New("limit must be at least 1")
			}

			if queryCmdState.FailOnCount != "" {
				var co failon.CountOperation
				if err := co.Parse(queryCmdState.FailOnCount); err != nil {
					return err
				}

				if _, err := co.IsFail(0); err != nil {
					return err
				}
			}
			return nil
		},
		RunE: runQuery,
	}
)

func init() {
	// add the lql command
	rootCmd.AddCommand(queryCmd)

	// add sub-commands to the lql command
	queryCmd.AddCommand(queryRunCmd)

	if cli.isLCLInstalled() {
		queryRunCmd.Flags().StringVarP(
			&queryCmdState.CURVFromLibrary,
			"library", "l", "",
			"run query from Lacework Content Library",
		)
	}

	// run specific flags
	setQuerySourceFlags(queryRunCmd)

	// limit flag
	queryRunCmd.Flags().IntVar(
		&queryCmdState.Limit,
		"limit", 0,
		"result limit for query (default 0)",
	)

	// range time flag
	queryRunCmd.Flags().StringVar(
		&queryCmdState.Range,
		"range", "",
		"natural time range for query",
	)

	// start time flag
	queryRunCmd.Flags().StringVar(
		&queryCmdState.Start,
		"start", "-24h",
		"start time for query",
	)
	// end time flag
	queryRunCmd.Flags().StringVar(
		&queryCmdState.End,
		"end", "now",
		"end time for query",
	)
	queryRunCmd.Flags().BoolVar(
		&queryCmdState.ValidateOnly,
		"validate_only", false,
		"validate query only (do not run)",
	)
	// fail on count
	queryRunCmd.Flags().StringVar(
		&queryCmdState.FailOnCount,
		"fail_on_count", "",
		"fail if the results from a query match the provided expression (e.g. '>0')",
	)
	// empty template flag
	queryRunCmd.Flags().BoolVar(
		&queryCmdState.EmptyTemplate,
		"empty", false,
		"start $EDITOR with empty file",
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
	// if running via library (CUV)
	if queryCmdState.CURVFromLibrary != "" {
		return inputQueryFromLibrary(queryCmdState.CURVFromLibrary)
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
		bytes, err := io.ReadAll(os.Stdin)
		return string(bytes), err
	}
	// if running via editor
	action := "validate"
	if !queryCmdState.ValidateOnly {
		action = strings.Split(cmd.Use, " ")[0]
	}
	return inputQueryFromEditor(action)
}

func inputQueryFromLibrary(id string) (string, error) {
	var (
		lcl *LaceworkContentLibrary
		err error
	)
	if lcl, err = cli.LoadLCL(); err != nil {
		return "", err
	}
	return lcl.GetQuery(id)
}

func inputQueryFromFile(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)

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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		err = errors.Wrap(err, msg)
		return
	}
	query = string(body)
	return
}

func inputQueryFromEditor(action string) (query string, err error) {
	language := "LQL"
	if action == "create" && !queryCmdState.EmptyTemplate && cli.LwApi.V2.Query.RegoQueryEnabled() {
		languageSelect := &survey.Select{
			Message: "Choose query language: ",
			Options: []string{
				"LQL",
				"Rego",
			},
		}
		err = survey.AskOne(languageSelect, &language)
		if err != nil {
			return
		}
	}

	var queryTextTemplate string
	switch language {
	case "LQL":
		queryTextTemplate = `queryId: YourQueryID
queryText: |-
    {
        source {
            --- Select a datasource. To list all available datasources, use 'lacework query sources'.
        }
        filter {
            --- Add query filter(s), if any. If not, remove this block.
        }
        return {
            --- List fields to return from the selected source. Use 'lacework query describe <datasource>'.
        }
    }`
	case "Rego":
		queryTextTemplate = `queryId: YourQueryID
queryLanguage: Rego
queryText: |-
  package your.package
  import your.dependency
  import data.lacework
  source := lacework.your.data.source.provider.function("servce", "apiKey")
  assess := your.assess.rule`
	default:
		err = errors.New("Invalid query language: " + language)
		return
	}

	prompt := &survey.Editor{
		Message:  fmt.Sprintf("Type a query to %s", action),
		FileName: "query*.yaml",
	}

	if (action == "create" || action == "run") && !queryCmdState.EmptyTemplate {
		prompt.Default = queryTextTemplate
		prompt.HideDefault = true
		prompt.AppendDefault = true
	} else if (action == "create" || action == "run") && queryCmdState.EmptyTemplate {
		prompt.Default = ``
		prompt.HideDefault = true
		prompt.AppendDefault = true
	}

	err = survey.AskOne(prompt, &query)
	return
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
Verify that the JSON is formatted properly and adheres to the following schema:

{
    "queryId": "MyLQL",
    "queryText": "{ source { CloudTrailRawEvents } filter { EVENT_SOURCE = 's3.amazonaws.com' } return { INSERT_ID } }"
}
`)
	}
	// smells like plain text
	return errors.New(`invalid query
	
It looks like you attempted to submit a query in YAML format.
Verify that the text adheres to the following schema:

queryId: MyLQL
queryText: |-
  {
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
		response   api.ExecuteQueryResponse
		msg        string = "unable to run query"
		hasCmdArgs bool   = len(args) != 0 && args[0] != ""
	)

	// check use of <query_id> with other flags
	if hasCmdArgs {
		var naFlag string

		if queryCmdState.File != "" {
			naFlag = "file"
		}
		if queryCmdState.CURVFromLibrary != "" {
			naFlag = "library"
		}
		if queryCmdState.URL != "" {
			naFlag = "url"
		}
		if queryCmdState.ValidateOnly {
			naFlag = "validate_only"
		}
		if queryCmdState.EmptyTemplate {
			naFlag = "empty"
		}
		if naFlag != "" {
			return errors.New(
				fmt.Sprintf(
					"flag --%s not applicable when specifying query_id argument",
					naFlag,
				),
			)
		}
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
		{
			Name:  api.QueryStartTimeRange,
			Value: start.UTC().Format(lwtime.RFC3339Milli),
		},
		{
			Name:  api.QueryEndTimeRange,
			Value: end.UTC().Format(lwtime.RFC3339Milli),
		},
	}

	if hasCmdArgs {
		// query by id
		response, err = runQueryByID(args[0], queryArgs)
	} else {
		// adhoc query
		response, err = runAdhocQuery(cmd, queryArgs)
	}

	if err != nil {
		return errors.Wrap(err, "unable to run query")
	}

	// output
	if err = cli.OutputJSON(response.Data); err != nil {
		return err
	}

	// fail_on_count post
	if queryCmdState.FailOnCount != "" {
		cli.Log.Infow("enforce failure flag(s)",
			"fail_on_count", queryCmdState.FailOnCount,
		)

		queryFailonError := NewQueryFailonError(
			queryCmdState.FailOnCount,
			len(response.Data),
		)
		if queryFailonError.NonCompliant() {
			cmd.SilenceUsage = true
			return queryFailonError
		}
	}
	return nil
}

func runQueryByID(id string, args []api.ExecuteQueryArgument) (
	api.ExecuteQueryResponse,
	error,
) {
	cli.Log.Debugw("running query", "query", id)

	cli.StartProgress(getRunStartProgressMessage(args))
	defer cli.StopProgress()

	opts := api.ExecuteQueryOptions{}
	// only add limit if > 0
	if queryCmdState.Limit > 0 {
		opts.Limit = &queryCmdState.Limit
	}

	request := api.ExecuteQueryByIDRequest{
		QueryID:   id,
		Options:   opts,
		Arguments: args,
	}
	return cli.LwApi.V2.Query.ExecuteByID(request)
}

func runAdhocQuery(cmd *cobra.Command, args []api.ExecuteQueryArgument) (
	response api.ExecuteQueryResponse,
	err error,
) {
	// input query
	queryString, err := inputQuery(cmd)
	if err != nil {
		return
	}
	// parse query
	newQuery, err := api.ParseNewQuery(queryString)
	if err != nil {
		err = queryErrorCrumbs(queryString)
		return
	}

	opts := api.ExecuteQueryOptions{}
	// only add limit if > 0
	if queryCmdState.Limit > 0 {
		opts.Limit = &queryCmdState.Limit
	}

	cli.StartProgress(getRunStartProgressMessage(args))
	defer cli.StopProgress()

	// execute query
	executeQuery := api.ExecuteQueryRequest{
		Query: api.ExecuteQuery{
			QueryText: newQuery.QueryText,
		},
		Options:   opts,
		Arguments: args,
	}

	cli.Log.Debugw("running query", "query", queryString)
	response, err = cli.LwApi.V2.Query.Execute(executeQuery)
	return
}

func getRunStartProgressMessage(args []api.ExecuteQueryArgument) string {
	var (
		startTime, endTime time.Time
		startErr           error = errors.New("StartTimeRange not present in ExecuteQueryArgument list")
		endErr             error = errors.New("EndTimeRange not present in ExecuteQueryArgument list")
	)
	for _, arg := range args {
		switch arg.Name {
		case api.QueryStartTimeRange:
			startTime, startErr = time.Parse(time.RFC3339, arg.Value)
		case api.QueryEndTimeRange:
			endTime, endErr = time.Parse(time.RFC3339, arg.Value)
		}
	}

	msg := "Executing query"
	if startErr == nil && endErr == nil {
		msg = fmt.Sprintf(
			"%s in the time range %s - %s",
			msg,
			startTime.Format("2006-Jan-2 15:04:05 MST"),
			endTime.Format("2006-Jan-2 15:04:05 MST"),
		)
	}
	return msg
}
