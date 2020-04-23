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
	"sort"
	"strings"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// eventCmd represents the event command
	eventCmd = &cobra.Command{
		Use:   "event",
		Short: "Inspect Lacework events",
	}

	// eventListCmd represents the list sub-command inside the event command
	eventListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all events from a date range (default last 7 days)",
		Long: `List all events from a data range, by default it displays the last
7 days, but you can specify a different time range.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			lacework, err := api.NewClient(cli.Account,
				api.WithLogLevel(cli.LogLevel),
				api.WithApiKeys(cli.KeyID, cli.Secret),
			)
			if err != nil {
				return errors.Wrap(err, "unable to generate api client")
			}

			response, err := lacework.Events.List()
			if err != nil {
				return errors.Wrap(err, "unable to get events")
			}

			// Sort the events from the response by severity
			sort.Slice(response.Events, func(i, j int) bool {
				return response.Events[i].Severity < response.Events[j].Severity
			})

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Events)
			}

			cli.OutputHuman(eventsToTableReport(response.Events))
			return nil
		},
	}

	// eventShowCmd represents the show sub-command inside the event command
	eventShowCmd = &cobra.Command{
		Use:   "show <event_id>",
		Short: "Create an external integrations",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
)

func init() {
	// add the integration command
	rootCmd.AddCommand(eventCmd)

	// add sub-commands to the integration command
	eventCmd.AddCommand(eventListCmd)
	eventCmd.AddCommand(eventShowCmd)
}

func eventsToTableReport(events []api.Event) string {
	var (
		eventsReport = &strings.Builder{}
		table        = tablewriter.NewWriter(eventsReport)
	)

	table.SetHeader([]string{
		"Event ID",
		"Type",
		"Severity",
		"Start Time",
		"End Time",
	})
	table.SetBorder(false)
	table.AppendBulk(eventsToTable(events))
	table.Render()

	return eventsReport.String()
}

func eventsToTable(events []api.Event) [][]string {
	out := [][]string{}
	for _, event := range events {
		out = append(out, []string{
			event.EventID,
			event.EventType,
			event.SeverityString(),
			event.StartTime.UTC().Format(time.RFC3339),
			event.EndTime.UTC().Format(time.RFC3339),
		})
	}
	return out
}
