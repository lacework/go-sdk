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
	"strconv"
	"strings"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/lacework/go-sdk/lwseverity"
	"github.com/lacework/go-sdk/lwtime"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// alertListCmd represents the alert list command
	alertListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all alerts",
		Long: `List all alerts.

By default, alerts are shown for the last 24 hours.
Use a custom time range by suppling a range flag...

    lacework alert ls --range "last 7 days"

Or by specifying start and end flags.

    lacework alert ls --start "-7d@d" --end "now"

Start and end times may be specified in one of the following formats:
    A. A relative time specifier
    B. RFC3339 date and time
    C. Epoch time in milliseconds

To list open alerts of type "NewViolations" with high or critical severity.

    lacework alert ls --status Open --severity high --type NewViolations
`,
		Args: cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, args []string) error {
			// validate severity
			if alertCmdState.Severity != "" && !array.ContainsStr(
				api.ValidAlertSeverities, alertCmdState.Severity) {
				return errors.Wrap(
					errors.New(fmt.Sprintf("the severity (%s) is not valid, use one of (%s)",
						alertCmdState.Severity, strings.Join(api.ValidAlertSeverities, ", "))),
					"unable to list alerts",
				)
			}
			// validate status
			if alertCmdState.Status != "" && !array.ContainsStr(
				api.ValidAlertStatuses, alertCmdState.Status) {
				return errors.Wrap(
					errors.New(fmt.Sprintf("the status (%s) is not valid, use one of (%s)",
						alertCmdState.Status, strings.Join(api.ValidAlertStatuses, ", "))),
					"unable to list alerts",
				)
			}
			return nil
		},
		RunE: listAlert,
	}
)

func init() {
	alertCmd.AddCommand(alertListCmd)

	// range time flag
	alertListCmd.Flags().StringVar(
		&alertCmdState.Range,
		"range", "",
		"natural time range for alerts",
	)

	// start time flag
	alertListCmd.Flags().StringVar(
		&alertCmdState.Start,
		"start", "-24h",
		"start time for alerts",
	)
	// end time flag
	alertListCmd.Flags().StringVar(
		&alertCmdState.End,
		"end", "now",
		"end time for alerts",
	)

	// severity flag
	alertListCmd.Flags().StringVar(
		&alertCmdState.Severity,
		"severity", "",
		fmt.Sprintf(
			"filter alerts by severity threshold (%s)",
			strings.Join(api.ValidAlertSeverities, ", "),
		),
	)

	// status flag
	alertListCmd.Flags().StringVar(
		&alertCmdState.Status,
		"status", "",
		fmt.Sprintf(
			"filter alerts by status (%s)",
			strings.Join(api.ValidAlertStatuses, ", "),
		),
	)

	// type flag
	alertListCmd.Flags().StringVar(
		&alertCmdState.Type,
		"type", "",
		"filter alerts by type",
	)
}

func alertListTable(alerts api.Alerts) (out [][]string) {
	alerts.SortByID()
	alerts.SortBySeverity()

	for _, alert := range alerts {
		// filter severity if desired
		if lwseverity.ShouldFilter(alert.Severity, alertCmdState.Severity) {
			continue
		}

		out = append(out, []string{
			strconv.Itoa(alert.ID),
			alert.Type,
			alert.Name,
			alert.Severity,
			alert.StartTime,
			alert.EndTime,
			alert.Status,
		})
	}

	return
}

func renderAlertListTable(alerts api.Alerts) {
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Alert ID", "Type", "Name", "Severity", "Start Time", "End Time", "Status"},
			alertListTable(alerts),
			tableFunc(func(t *tablewriter.Table) {
				t.SetAutoWrapText(false)
				t.SetBorder(false)
			}),
		),
	)
}

func listAlert(_ *cobra.Command, _ []string) error {
	cli.Log.Debugw("listing alerts")

	var (
		err   error
		start time.Time
		end   time.Time
		msg   string = "unable to list alerts"
	)

	// use of if/else intentional here based on logic paths for determining start and end time.Time values
	// if cli user has specified a range we use ParseNatural which gives us start and end time.Time values
	// otherwise we need to convert alertCmdState start and end strings to time.Time values using parseQueryTime
	if alertCmdState.Range != "" {
		cli.Log.Debugw("retrieving natural time range")

		start, end, err = lwtime.ParseNatural(alertCmdState.Range)
		if err != nil {
			return errors.Wrap(err, msg)
		}
	} else {
		// parse start
		start, err = parseQueryTime(alertCmdState.Start)
		if err != nil {
			return errors.Wrap(err, msg)
		}
		// parse end
		end, err = parseQueryTime(alertCmdState.End)
		if err != nil {
			return errors.Wrap(err, msg)
		}
	}

	filters := []api.Filter{}

	if alertCmdState.Status != "" {
		filters = append(filters, api.Filter{
			Expression: "eq",
			Field:      string(api.AlertsFilterFieldStatus),
			Value:      alertCmdState.Status,
		})
	}

	if alertCmdState.Type != "" {
		filters = append(filters, api.Filter{
			Expression: "eq",
			Field:      string(api.AlertsFilterFieldType),
			Value:      alertCmdState.Type,
		})
	}

	cli.StartProgress(
		fmt.Sprintf(
			" Fetching alerts in the time range %s - %s...",
			start.Format("2006-Jan-2 15:04:05 MST"),
			end.Format("2006-Jan-2 15:04:05 MST"),
		),
	)
	var searchRequest = api.SearchFilter{
		TimeFilter: &api.TimeFilter{StartTime: &start, EndTime: &end},
		Filters:    filters,
	}

	listResponse, err := cli.LwApi.V2.Alerts.SearchAll(searchRequest)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, msg)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(listResponse.Data)
	}

	if len(listResponse.Data) == 0 {
		cli.OutputHuman("There are no alerts in your account in the specified time range.\n")
		return nil
	}
	renderAlertListTable(listResponse.Data)

	// breadcrumb
	cli.OutputHuman("\nUse 'lacework alert show <alert_id>' to see details for a specific alert.\n")
	return nil
}
