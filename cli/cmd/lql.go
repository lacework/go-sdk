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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwtime"
)

var (
	lqlCmdState = struct {
		End          string
		File         string
		Repo         bool
		Range        string
		Start        string
		URL          string
		ValidateOnly bool
		FailOnCount  string
	}{}

	// lqlCmd represents the lql parent command
	lqlCmd = &cobra.Command{
		Hidden:  true,
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

Run a query via editor:

	$ lacework query run --range today

Run a query via ID (uses active profile):

	$ lacework query run MyQuery --range

Start and End times are required to run a query:

1.  Start and End times must be specified in one of the following formats:

	A. A relative time specifier
	B. RFC3339 Date and Time
	C. Epoch time in milliseconds

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
	setQuerySourceFlags(lqlRunCmd)

	// since time flag
	lqlRunCmd.Flags().StringVarP(
		&lqlCmdState.Range,
		"range", "", "",
		"natural time range for LQL query",
	)

	// start time flag
	lqlRunCmd.Flags().StringVarP(
		&lqlCmdState.Start,
		"start", "", "@d",
		"start time for LQL query",
	)
	// end time flag
	lqlRunCmd.Flags().StringVarP(
		&lqlCmdState.End,
		"end", "", "now",
		"end time for LQL query",
	)
	lqlRunCmd.Flags().BoolVarP(
		&lqlCmdState.ValidateOnly,
		"validate_only", "", false,
		"validate query only (do not run)",
	)
	// fail on count
	lqlRunCmd.Flags().StringVarP(
		&lqlCmdState.FailOnCount,
		"fail_on_count", "", "",
		"fail if the query matches the fail_on_count expression",
	)
}

func setQuerySourceFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		if cmd != nil {
			action := strings.Split(cmd.Use, " ")[0]

			// file flag to specify a query from disk
			cmd.Flags().StringVarP(
				&lqlCmdState.File,
				"file", "f", "",
				fmt.Sprintf("path to an LQL query to %s", action),
			)
			/* repo flag to specify a query from repo
			cmd.Flags().BoolVarP(
				&lqlCmdState.Repo,
				"repo", "r", false,
				fmt.Sprintf("id of an LQL query to %s via active repo", action),
			)*/
			// url flag to specify a query from url
			cmd.Flags().StringVarP(
				&lqlCmdState.URL,
				"url", "u", "",
				fmt.Sprintf("url to an LQL query to %s", action),
			)
		}
	}
}

// for commands that take a query as input
func inputQuery(cmd *cobra.Command, args []string) (string, error) {
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

func inputQueryFromEnv(queryID string) (query string, err error) {
	var queryResponse api.LQLQueryResponse

	queryResponse, err = cli.LwApi.LQL.GetQueryByID(queryID)
	if err == nil && len(queryResponse.Data) != 0 {
		query = queryResponse.Data[0].QueryText
	}
	return
}

func inputQueryFromRepo() (query string, err error) {
	err = errors.New("NotImplementedError")
	return
}

func inputQueryFromFile(filePath string) (query string, err error) {
	fileData, err := ioutil.ReadFile(filePath)

	if err != nil {
		err = errors.Wrap(err, "unable to read file")
		return
	}

	query = string(fileData)
	return
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

func queryErrorCrumbs(query string, err error) error {
	// not the error we're looking for
	if !strings.Contains(fmt.Sprintf("%s", err), "unable to translate query blob") {
		return err
	}
	// smells like json
	query = strings.TrimLeft(query, " ")
	if strings.HasPrefix(query, "{") || strings.HasPrefix(query, "[") {
		return errors.New(`invalid query

It looks like you attempted to submit an LQL query in JSON format.
Please validate that the JSON is formatted properly and adheres to the following schema:

{
	"QUERY_TEXT": "MyLQL(CloudTrailRawEvents e) { SELECT INSERT_ID }"
}`)
	}
	// smells like plain text
	return errors.New(`invalid query
	
It looks like you attempted to submit an LQL query in plain text format.
Please validate that the text adheres to the following schema:

MyLQL(CloudTrailRawEvents e) { 
	SELECT INSERT_ID 
}
`)
}

func runQuery(cmd *cobra.Command, args []string) error {
	msg := "unable to run LQL query"

	query, err := inputQuery(cmd, args)
	if err != nil {
		return errors.Wrap(err, msg)
	}

	if lqlCmdState.Range != "" {
		cli.Log.Debugw("retrieving natural time range")

		var start, end time.Time
		start, end, err = lwtime.ParseNatural(lqlCmdState.Range)
		if err != nil {
			return errors.Wrap(err, msg)
		}
		lqlCmdState.Start = start.UTC().Format(time.RFC3339)
		lqlCmdState.End = end.UTC().Format(time.RFC3339)
	}
	// fail_on_count pre
	var co countOperation
	if lqlCmdState.FailOnCount != "" {
		err = co.parse(lqlCmdState.FailOnCount)
		if err != nil {
			return err
		}

		_, err = co.isFail(0)
		if err != nil {
			return err
		}
	}

	cli.Log.Debugw("running LQL query", "query", query)

	// validate_only should compile
	if lqlCmdState.ValidateOnly {
		return compileQueryAndOutput(query)
	}
	// !validate_only should should run
	response, err := cli.LwApi.LQL.RunQuery(query, lqlCmdState.Start, lqlCmdState.End)

	if err != nil {
		err = queryErrorCrumbs(query, err)
		return errors.Wrap(err, msg)
	}
	// output
	if err = cli.OutputJSON(response.Data); err != nil {
		return err
	}
	// fail_on_count post
	if lqlCmdState.FailOnCount != "" {
		isFail, err := co.isFail(len(response.Data))
		if err != nil {
			return err
		}
		if isFail {
			os.Exit(9)
		}
	}
	return nil
}

const operationRE = `^(>|>=|<|<=|={1,2}|!=)\s*(\d+)$`

type countOperation struct {
	operator string
	num      int
}

func (co *countOperation) parse(s string) error {
	re := regexp.MustCompile(operationRE)

	s = strings.TrimSpace(s)

	var op_parts []string
	if op_parts = re.FindStringSubmatch(s); s == "" || op_parts == nil {
		return errors.New(
			fmt.Sprintf("count operation (%s) is invalid", s))
	}
	co.num, _ = strconv.Atoi(op_parts[2])
	co.operator = op_parts[1]
	return nil
}

func (co countOperation) isFail(count int) (bool, error) {
	switch co.operator {
	case ">":
		return count > co.num, nil
	case ">=":
		return count >= co.num, nil
	case "<":
		return count < co.num, nil
	case "<=":
		return count <= co.num, nil
	case "=", "==":
		return count == co.num, nil
	case "!=":
		return count != co.num, nil
	}
	return true, errors.New(fmt.Sprintf("count operation (%s) is invalid", co.operator))
}
