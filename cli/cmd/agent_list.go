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

    lacework agent list --filter 'os:Linux' --filter 'tags.VpcId:vpc-72225916'

**NOTE:** The value can be a regular expression such as 'hostname:db-server.*'

To filter hosts with a running agent version '5.8.0'.

    lacework agent list --filter 'agentVersion:5.8.0.*' --filter 'status:ACTIVE'

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
		"filter results by key:value pairs (e.g. 'hostname:db-server.*')",
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
	err := cli.LwApi.V2.AgentInfo.Search(response, filters)
	cli.StopProgress()
	if err != nil {
		if strings.Contains(err.Error(), "Invalid field") {
			return errors.Errorf("the provided filter key is invalid.\n\nThe available keys for this command are:\n%s",
				stringSliceToMarkdownList(agentListCmdState.AvailableFilters.Filters))
		}
		return errors.Wrap(err, "unable to list agents via AgentInfo search")
	}

	agents := response.Data

	if response.Paging.Urls.NextPage != "" {
		totalPages := response.Paging.TotalRows / response.Paging.Rows

		agents = []api.AgentInfo{}
		page := 0
		for {
			agents = append(agents, response.Data...)

			cli.StartProgress(fmt.Sprintf("%s [%.0f%%]...", progressMsg, float32(page)/float32(totalPages)*100))
			pageOk, err := cli.LwApi.NextPage(response)
			if err == nil && pageOk {
				page += 1
				continue
			}
			break
		}
		response.ResetPaging()
		response.Data = agents
	}

	cli.StartProgress("Crunching agent data...")
	// Sort agents by last updated time (last time they check-in)
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].LastUpdate.After(agents[j].LastUpdate)
	})
	cli.StopProgress()

	if cli.JSONOutput() {
		return cli.OutputJSON(agents)
	}

	if len(agents) == 0 {
		if len(agentListCmdState.Filters) != 0 {
			cli.OutputHuman("No agents found with the provided filter(s).\n")
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
			[]string{
				"MID", "Status", "Agent Version", "Hostname", "Name",
				"Internal IP", "External IP", "OS Arch", "Last Check-in", "Created Time",
			},
			agentInfoToListAgentTable(agents),
		),
	)

	// breadcrumbs
	if len(agentListCmdState.Filters) == 0 {
		cli.OutputHuman("\nTry adding '--filter status:ACTIVE' to show only active agents.\n")
	} else if hasWindowsAgents(agents) && !agentListFiltersContains("os") {
		cli.OutputHuman("\nTry adding '--filter os:Windows' to show only windows agents.\n")
	} else if !hasWindowsAgents(agents) && !agentListFiltersContains("agentVersion") {
		cli.OutputHuman("\nTry adding '--filter \"agentVersion:5.8.0.*\"' to show agents with version '5.8.0'.\n")
	}
	return nil
}

func agentInfoToListAgentTable(agents []api.AgentInfo) [][]string {
	out := [][]string{}
	for _, a := range agents {
		out = append(out, []string{
			fmt.Sprintf("%d", a.Mid),
			a.Status,
			a.AgentVersion,
			a.Hostname,
			a.Tags.Name,
			a.Tags.InternalIP,
			a.Tags.ExternalIP,
			fmt.Sprintf("%s/%s", a.Tags.Os, a.Tags.Arch),
			a.LastUpdate.Format(time.RFC3339),
			a.CreatedTime.Format(time.RFC3339),
		})
	}

	return out
}

// agentListFiltersContains returns true if one of filters passed to this function
// matches the filters provided to the 'agent list' command
func agentListFiltersContains(filters ...string) bool {
	for _, cmdFilter := range agentListCmdState.Filters {
		for _, expectedFilter := range filters {
			if strings.Contains(cmdFilter, expectedFilter) {
				return true
			}
		}
	}
	return false
}

// hasWindowsAgents returns true if there the user has windows agents
func hasWindowsAgents(agents []api.AgentInfo) bool {
	for _, a := range agents {
		if a.Os == "Windows" {
			return true
		}
	}
	return false
}
