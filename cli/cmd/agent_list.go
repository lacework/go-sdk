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
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	agentListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all hosts with a running agent",
		Long:    `List all hosts that have a running agent in your environment.`,
		RunE:    listAgents,
	}
)

func listAgents(_ *cobra.Command, _ []string) error {
	cli.StartProgress("Fetching list of agents...")
	response, err := cli.LwApi.V2.Entities.ListAllMachineDetails()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to list agents via MachineDetails search")
	}

	// Sort machine details by last time seen
	sort.Slice(response.Data, func(i, j int) bool {
		return response.Data[i].CreatedTime.After(response.Data[j].CreatedTime)
	})

	// clean duplicate machines
	machines := cleanDuplicateMachine(response.Data)

	if cli.JSONOutput() {
		return cli.OutputJSON(machines)
	}

	if len(machines) == 0 {
		cli.OutputHuman(
			"There are no agents running in your account.\n\nTry installing one with 'lacework agent install <host>%s'\n",
			cli.OutputNonDefaultProfileFlag(),
		)
		return nil
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"MID", "Short Agent Token", "Hostname", "Name", "IP Address", "External IP", "OS Arch", "Last Checkin"},
			machineDetailsToListAgentTable(machines),
		),
	)
	return nil
}

func cleanDuplicateMachine(machines []api.MachineDetails) []api.MachineDetails {
	var cleanedMachines []api.MachineDetails

	for _, m := range machines {
		if machineExist(cleanedMachines, m.Mid) {
			continue
		}
		cleanedMachines = append(cleanedMachines, m)
	}

	return cleanedMachines
}

func machineExist(machines []api.MachineDetails, mid int) bool {
	for _, m := range machines {
		if mid == m.Mid {
			return true
		}
	}
	return false
}

func machineDetailsToListAgentTable(machines []api.MachineDetails) [][]string {
	out := [][]string{}
	for _, m := range machines {
		out = append(out, []string{
			fmt.Sprintf("%d", m.Mid),
			m.Tags.LwTokenShort,
			m.Hostname,
			m.Tags.Name,
			m.Tags.InternalIP,
			m.Tags.ExternalIP,
			fmt.Sprintf("%s/%s", m.Tags.Os, m.Tags.Arch),
			m.CreatedTime.Format(time.RFC3339),
		})
	}
	return out
}
