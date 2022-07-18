//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	agentListCmdState = struct {
		// The available filters for the agent list command
		AvailableFilters CmdFilters

		// List of filters to apply to the agent list command
		Filters []string
	}{}

	agentListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all hosts with a running agent",
		Long: `List all hosts that have a running agent in your environment.

You can use 'key:value' pairs to filter the list of hosts with the --filter flag.

    lacework agent list --filter 'os:Amazon Linux' --filter 'tags.VpcId:vpc-72225916'

**NOTE:** The value can be a regular expression such as 'hostname:db-server.*'

The available keys for this command are:
` + stringSliceToMarkdownList(
			agentListCmdState.AvailableFilters.GetFiltersFrom(
				api.AgentInfo{},
			),
		),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return validateKeyValuePairs(agentListCmdState.Filters)
		},
		RunE: listAgents,
	}
)

func init() {
	agentListCmd.Flags().StringSliceVar(&agentListCmdState.Filters, "filter", []string{},
		"filter results by key:value pairs (e.g. 'tags.Name:worker-.*')",
	)
}

func listAgents(_ *cobra.Command, _ []string) error {
	var (
		progressMsg = "Fetching list of agents"
		response    = &api.AgentInfoResponse{}
		now         = time.Now().UTC().Add(-1 * time.Minute)
		before      = now.AddDate(0, 0, -7) // 7 days from ago
		filters     = api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
		}
	)

	cleanedF := []string{}
	if len(agentListCmdState.Filters) != 0 {
		filters.Filters = []api.Filter{}
		for _, pair := range agentListCmdState.Filters {

			kv := strings.Split(pair, ":")
			if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
				cli.Log.Warnw("malformed filter, ignoring",
					"pair", pair, "expected_format", "key:value",
				)
				continue
			}

			cleanedF = append(cleanedF, pair)
			cli.Log.Infow("adding filter", "key", kv[0], "value", kv[1])
			filters.Filters = append(filters.Filters, api.Filter{
				Field:      kv[0],
				Expression: cli.lqlOperator, // @afiune we use rlike to allow user to pass regex
				Value:      kv[1],
			})
		}

		if len(cleanedF) != 0 {
			progressMsg = fmt.Sprintf(
				"%s with filters (%s)",
				progressMsg, strings.Join(cleanedF, ", "),
			)
		}

		agentListCmdState.Filters = cleanedF
	}

	cli.StartProgress(fmt.Sprintf("%s...", progressMsg))
	err := cli.LwApi.V2.AgentAccessTokens.Search(response, filters)
	cli.StopProgress()
	if err != nil {
		if strings.Contains(err.Error(), "Invalid field") {
			return errors.Errorf("the provided filter key is invalid.\n\nThe available keys for this command are:\n%s",
				stringSliceToMarkdownList(agentListCmdState.AvailableFilters.Filters))
		}
		return errors.Wrap(err, "unable to list agents via AgentInfo search")
	}

	if response.Paging.Urls.NextPage != "" {
		totalPages := response.Paging.TotalRows / response.Paging.Rows

		all := []api.AgentInfo{}
		page := 0
		for {
			all = append(all, response.Data...)

			cli.StartProgress(fmt.Sprintf("%s [%.0f%%]...", progressMsg, float32(page)/float32(totalPages)*100))
			pageOk, err := cli.LwApi.NextPage(response)
			if err == nil && pageOk {
				page += 1
				continue
			}
			break
		}
		response.Data = all
		response.ResetPaging()
	}

	cli.StartProgress("Crunching agent data...")
	// Sort machine details by last time seen
	sort.Slice(response.Data, func(i, j int) bool {
		return response.Data[i].CreatedTime.After(response.Data[j].CreatedTime)
	})

	// clean duplicate machines
	machines := cleanDuplicateMachine(response.Data)
	cli.StopProgress()

	if cli.JSONOutput() {
		return cli.OutputJSON(machines)
	}

	if len(machines) == 0 {
		if len(agentListCmdState.Filters) != 0 {
			cli.OutputHuman("No agent found with the provided filter(s).\n")
		} else {
			cli.OutputHuman(
				"There are no agents running in your account.\n\nTry installing one with 'lacework agent install <host>%s'\n",
				cli.OutputNonDefaultProfileFlag(),
			)
		}
		return nil
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"MID", "Short Agent Token", "Hostname", "Name", "Status", "IP Address", "External IP", "OS Arch", "Last Checkin"},
			agentInfoToListAgentTable(machines),
		),
	)
	return nil
}

func cleanDuplicateMachine(machines []api.AgentInfo) []api.AgentInfo {
	var cleanedMachines []api.AgentInfo

	for _, m := range machines {
		if machineExist(cleanedMachines, m.Mid) {
			continue
		}
		cleanedMachines = append(cleanedMachines, m)
	}

	return cleanedMachines
}

func machineExist(machines []api.AgentInfo, mid int) bool {
	for _, m := range machines {
		if mid == m.Mid {
			return true
		}
	}
	return false
}

func agentInfoToListAgentTable(machines []api.AgentInfo) [][]string {
	out := [][]string{}
	for _, m := range machines {
		out = append(out, []string{
			fmt.Sprintf("%d", m.Mid),
			m.Tags.LwTokenShort,
			m.Hostname,
			m.Tags.Name,
			m.Status,
			m.Tags.InternalIP,
			m.Tags.ExternalIP,
			fmt.Sprintf("%s/%s", m.Tags.Os, m.Tags.Arch),
			m.CreatedTime.Format(time.RFC3339),
		})
	}
	return out
}
