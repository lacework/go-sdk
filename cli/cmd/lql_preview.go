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
	"github.com/lacework/go-sdk/lwtime"
)

var (
	queryPreviewSourceCmd = &cobra.Command{
		Use:   "preview-source <datasource_id>",
		Short: "Preview Lacework query data source",
		Long:  `Preview Lacework query data source.`,
		Args:  cobra.ExactArgs(1),
		RunE:  previewQuerySource,
	}
	queryPreviewSourceTemplate = `CLIAdhocPreview { source { %s } return distinct { %s } }`
)

func init() {
	queryCmd.AddCommand(queryPreviewSourceCmd)
}

func previewQuerySource(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("retrieving datasource", "id", args[0])

	cli.StartProgress(" Retrieving datasource...")
	datasourceResponse, err := cli.LwApi.V2.Datasources.Get(args[0])
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "unable to retrieve data source")
	}

	// build returns list from datasource fields
	var returns []string
	for _, ret := range datasourceResponse.Data.ResultSchema {
		returns = append(returns, ret.Name)
	}
	if len(returns) == 0 {
		return errors.New("unable to parse data source schema")
	}

	// initialize query
	executeQuery := api.ExecuteQueryRequest{
		Query: api.ExecuteQuery{
			QueryText: fmt.Sprintf(
				queryPreviewSourceTemplate, args[0], strings.Join(returns, ",")),
		},
	}
	if args[0] == "CloudTrailRawEvents" {
		executeQuery.Query.EvaluatorID = "Cloudtrail"
	}

	// initialize time attempts
	timeAttempts := []map[string]string{
		{"start": "-24h", "end": "now"},
		{"start": "-7d", "end": "-24h"},
		{"start": "-30d", "end": "-7d"},
	}

	for _, timeAttempt := range timeAttempts {
		start, _ := lwtime.ParseRelative(timeAttempt["start"])
		end, _ := lwtime.ParseRelative(timeAttempt["end"])

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

		executeQuery.Arguments = queryArgs

		// execute query
		cli.Log.Debugw("running query", "query", executeQuery.Query.QueryText)
		cli.StartProgress(" Executing preview query...")
		response, err := cli.LwApi.V2.Query.Execute(executeQuery)
		cli.StopProgress()
		if err != nil {
			return errors.Wrap(err, "unable to preview data source")
		}

		// check and output
		data, ok := response["data"]
		if !ok {
			err = errors.New("preview results missing data")
			return errors.Wrap(err, "unable to preview data source")
		}
		if slice, ok := data.([]interface{}); ok {
			if len(slice) == 0 {
				continue
			}
			return cli.OutputJSON(slice[0])
		} else {
			err = errors.New("preview results data is not a slice")
			return errors.Wrap(err, "unable to preview data source")
		}
	}
	cli.OutputHuman("No results found for data source")
	return nil
}
