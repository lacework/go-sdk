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
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
)

var (
	eventsCmdState = struct {
		// start time for listing events
		Start string

		// end time for listing events
		End string

		// list events from an specific number of days
		Days int

		// list events with a specific severity
		Severity string
	}{}

	// eventCmd represents the event command
	eventCmd = &cobra.Command{
		Use:     "event",
		Aliases: []string{"events"},
		Short:   "inspect Lacework events",
		Long:    `Inspect events reported by the Lacework platform`,
	}

	// eventListCmd represents the list sub-command inside the event command
	eventListCmd = &cobra.Command{
		Use:   "list",
		Short: "list all events (default last 7 days)",
		Long: `List all events for the last 7 days by default, or pass --start and --end to
specify a custom time period. You can also pass --serverity to filter by a
severity threshold.

Additionally, pass --days to list events for a specified number of days.

For example, to list all events from the last day with severity medium and above
(Critical, High and Medium) run:

    $ lacework events list --severity medium --days 1`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			var (
				response api.EventsResponse
				err      error
			)

			if eventsCmdState.Severity != "" {
				if !array.ContainsStr(api.ValidEventSeverities, eventsCmdState.Severity) {
					return errors.Errorf("the severity %s is not valid, use one of %s",
						eventsCmdState.Severity, strings.Join(api.ValidEventSeverities, ", "),
					)
				}
			}

			if eventsCmdState.Start != "" || eventsCmdState.End != "" {
				start, end, errT := parseStartAndEndTime(eventsCmdState.Start, eventsCmdState.End)
				if errT != nil {
					return errors.Wrap(errT, "unable to parse time range")
				}

				cli.Log.Infow("requesting list of events from custom time range",
					"start_time", start, "end_time", end,
				)
				response, err = cli.LwApi.Events.ListDateRange(start, end)
			} else if eventsCmdState.Days != 0 {
				end := time.Now()
				start := end.Add(time.Hour * 24 * time.Duration(eventsCmdState.Days) * -1)

				cli.Log.Infow("requesting list of events from specific days",
					"days", eventsCmdState.Days, "start_time", start, "end_time", end,
				)
				response, err = cli.LwApi.Events.ListDateRange(start, end)
			} else {
				cli.Log.Info("requesting list of events from the last 7 days")
				response, err = cli.LwApi.Events.List()
			}

			if err != nil {
				return errors.Wrap(err, "unable to get events")
			}

			cli.Log.Debugw("events", "raw", response)

			// filter events by severity, if the user didn't specify a severity
			// the funtion will return it back without modifications
			events := filterEventsWithSeverity(response.Events)

			// Sort the events by severity
			sort.Slice(events, func(i, j int) bool {
				return events[i].Severity < events[j].Severity
			})

			if cli.JSONOutput() {
				return cli.OutputJSON(events)
			}

			if len(events) == 0 {
				if eventsCmdState.Severity != "" {
					cli.OutputHuman("There are no events with the specified severity.\n")
				} else {
					cli.OutputHuman("There are no events in your account in the specified time range.\n")
				}
				return nil
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"Event ID", "Type", "Severity", "Start Time", "End Time"},
					eventsToTable(events),
				),
			)
			return nil
		},
	}

	// eventShowCmd represents the show sub-command inside the event command
	eventShowCmd = &cobra.Command{
		Use:   "show <event_id>",
		Short: "show details about a specific event",
		Long:  "Show details about a specific event.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cli.Log.Infow("requesting event details", "event_id", args[0])
			response, err := cli.LwApi.Events.Details(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get event details")
			}

			cli.Log.Debugw("event details",
				"event_id", args[0],
				"raw", response,
			)
			if len(response.Events) == 0 {
				return errors.Errorf("there are no details about the event '%s'", args[0])
			}

			// @afiune why do we have an array of events when we ask for details
			// about a single event? Let us use the first one for now
			if cli.JSONOutput() {
				return cli.OutputJSON(response.Events[0])
			}

			cli.OutputHuman(eventDetailsSummaryReport(response.Events[0]))
			for _, entityTable := range eventEntityMapTables(response.Events[0].EntityMap) {
				cli.OutputHuman("\n")
				cli.OutputHuman(entityTable)
			}

			cli.OutputHuman(
				"\nFor further investigation of this event navigate to %s\n",
				eventLinkBuilder(args[0]),
			)
			return nil
		},
	}

	// eventOpenCmd represents the open sub-command inside the event command
	eventOpenCmd = &cobra.Command{
		Use:   "open <event_id>",
		Short: "open a specified event in a web browser",
		Long:  "Open a specified event in a web browser.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// Event IDs should be only numeric values
			if _, err := strconv.Atoi(args[0]); err != nil {
				return errors.Errorf("invalid event id %s. Event id should be a numeric value", args[0])
			}

			var (
				err error
				url = eventLinkBuilder(args[0])
			)

			switch runtime.GOOS {
			case "linux":
				err = exec.Command("xdg-open", url).Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
			case "darwin":
				err = exec.Command("open", url).Start()
			default:
				err = fmt.Errorf("unsupported platform\n\nNavigate to %s", url)
			}
			if err != nil {
				return errors.Wrap(err, "unable to open web browser")
			}

			return nil
		},
	}
)

func init() {
	// add the event command
	rootCmd.AddCommand(eventCmd)

	// add sub-commands to the event command
	eventCmd.AddCommand(eventListCmd)

	// add start flag to events list command
	eventListCmd.Flags().StringVar(&eventsCmdState.Start,
		"start", "", "start of the time range in UTC (format: yyyy-MM-ddTHH:mm:ssZ)",
	)
	// add end flag to events list command
	eventListCmd.Flags().StringVar(&eventsCmdState.End,
		"end", "", "end of the time range in UTC (format: yyyy-MM-ddTHH:mm:ssZ)",
	)
	// add days flag to events list command
	eventListCmd.Flags().IntVar(&eventsCmdState.Days,
		"days", 0, "list events for specified number of days (max: 7 days)",
	)
	// add severity flag to events list command
	eventListCmd.Flags().StringVar(&eventsCmdState.Severity,
		"severity", "",
		fmt.Sprintf(
			"filter events by severity threshold (%s)",
			strings.Join(api.ValidEventSeverities, ", "),
		),
	)

	eventCmd.AddCommand(eventShowCmd)
	eventCmd.AddCommand(eventOpenCmd)
}

// Generates a URL similar to:
//   => https://account.lacework.net/ui/investigate/recents/EventDossier-123
func eventLinkBuilder(id string) string {
	return fmt.Sprintf("https://%s.lacework.net/ui/investigation/recents/EventDossier-%s", cli.Account, id)
}

// easily add or remove borders to all event details tables
func eventsTableOptions() tableOption {
	return tableFunc(func(t *tablewriter.Table) {
		t.SetBorders(tablewriter.Border{
			Left:   false,
			Right:  false,
			Top:    true,
			Bottom: false,
		})
		t.SetColumnSeparator(" ")
	})
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

func eventDetailsSummaryReport(details api.EventDetails) string {
	return renderCustomTable(
		[]string{"Event ID", "Type", "Actor", "Model", "Start Time", "End Time"},
		[][]string{[]string{
			details.EventID,
			details.EventType,
			details.EventActor,
			details.EventModel,
			details.StartTime.UTC().Format(time.RFC3339),
			details.EndTime.UTC().Format(time.RFC3339),
		}},
		eventsTableOptions(),
	)
}

func eventEntityMapTables(eventEntities api.EventEntityMap) []string {
	tables := []string{}

	if machineTable := eventMachineEntitiesTable(eventEntities.Machine); machineTable != "" {
		tables = append(tables, machineTable)
	}
	if containerTable := eventContainerEntitiesTable(eventEntities.Container); containerTable != "" {
		tables = append(tables, containerTable)
	}
	if appTable := eventApplicationEntitiesTable(eventEntities.Application); appTable != "" {
		tables = append(tables, appTable)
	}
	if userTable := eventUserEntitiesTable(eventEntities.User); userTable != "" {
		tables = append(tables, userTable)
	}
	if ipaddressTable := eventIpAddressEntitiesTable(eventEntities.IpAddress); ipaddressTable != "" {
		tables = append(tables, ipaddressTable)
	}
	if sourceIpAddrTable := eventSourceIpAddressEntitiesTable(eventEntities.SourceIpAddress); sourceIpAddrTable != "" {
		tables = append(tables, sourceIpAddrTable)
	}
	if dnsTable := eventDnsNameEntitiesTable(eventEntities.DnsName); dnsTable != "" {
		tables = append(tables, dnsTable)
	}
	if apiTable := eventAPIEntitiesTable(eventEntities.API); apiTable != "" {
		tables = append(tables, apiTable)
	}
	if ctUserTable := eventCTUserEntitiesTable(eventEntities.CTUser); ctUserTable != "" {
		tables = append(tables, ctUserTable)
	}
	if regionTable := eventRegionEntitiesTable(eventEntities.Region); regionTable != "" {
		tables = append(tables, regionTable)
	}
	if processTable := eventProcessEntitiesTable(eventEntities.Process); processTable != "" {
		tables = append(tables, processTable)
	}
	if exePathTable := eventFileExePathEntitiesTable(eventEntities.FileExePath); exePathTable != "" {
		tables = append(tables, exePathTable)
	}
	if dataHashTable := eventFileDataHashEntitiesTable(eventEntities.FileDataHash); dataHashTable != "" {
		tables = append(tables, dataHashTable)
	}
	if cRuleTable := eventCustomRuleEntitiesTable(eventEntities.CustomRule); cRuleTable != "" {
		tables = append(tables, cRuleTable)
	}
	if violationTable := eventNewViolationEntitiesTable(eventEntities.NewViolation); violationTable != "" {
		tables = append(tables, violationTable)
	}
	if recordsTable := eventRecIDEntitiesTable(eventEntities.RecID); recordsTable != "" {
		tables = append(tables, recordsTable)
	}
	if vReasonTable := eventViolationReasonEntitiesTable(eventEntities.ViolationReason); vReasonTable != "" {
		tables = append(tables, vReasonTable)
	}
	if resourceTable := eventResourceEntitiesTable(eventEntities.Resource); resourceTable != "" {
		tables = append(tables, resourceTable)
	}

	return tables
}

func eventRegionEntitiesTable(regions []api.EventRegionEntity) string {
	if len(regions) == 0 {
		return ""
	}

	data := [][]string{}
	for _, user := range regions {
		data = append(data, []string{
			user.Region,
			strings.Join(user.AccountList, ", "),
		})
	}

	return renderCustomTable(
		[]string{"Region", "Accounts"}, data,
		eventsTableOptions(),
	)
}

func eventCTUserEntitiesTable(users []api.EventCTUserEntity) string {
	if len(users) == 0 {
		return ""
	}

	data := [][]string{}
	for _, user := range users {
		mfa := "Disabled"
		if user.Mfa != 0 {
			mfa = "Enabled"
		}
		data = append(data, []string{
			user.Username,
			user.AccountID,
			user.PrincipalID,
			mfa,
			strings.Join(user.ApiList, ", "),
			strings.Join(user.RegionList, ", "),
		})
	}

	return renderCustomTable(
		[]string{"Username", "Account ID", "Principal ID", "MFA", "List of APIs", "Regions"}, data,
		eventsTableOptions(),
	)
}

func eventDnsNameEntitiesTable(dnss []api.EventDnsNameEntity) string {
	if len(dnss) == 0 {
		return ""
	}

	data := [][]string{}
	for _, d := range dnss {
		data = append(data, []string{
			d.Hostname,
			array.JoinInt32(d.PortList, ", "),
			fmt.Sprintf("%.3f", d.TotalInBytes),
			fmt.Sprintf("%.3f", d.TotalOutBytes),
		})
	}

	return renderCustomTable(
		[]string{"DNS Hostname", "List of Ports", "Inbound Bytes", "Outboud Bytes"}, data,
		eventsTableOptions(),
	)
}

func eventAPIEntitiesTable(apis []api.EventAPIEntity) string {
	if len(apis) == 0 {
		return ""
	}

	data := [][]string{}
	for _, a := range apis {
		data = append(data, []string{a.Service, a.Api})
	}

	return renderCustomTable(
		[]string{"Service", "API"}, data,
		eventsTableOptions(),
	)
}

func eventSourceIpAddressEntitiesTable(ips []api.EventSourceIpAddressEntity) string {
	if len(ips) == 0 {
		return ""
	}

	data := [][]string{}
	for _, ip := range ips {
		data = append(data, []string{ip.IpAddress, ip.Country, ip.Region})
	}

	return renderCustomTable(
		[]string{"Source IP Address", "Country", "Region"}, data,
		eventsTableOptions(),
	)
}

func eventIpAddressEntitiesTable(ips []api.EventIpAddressEntity) string {
	if len(ips) == 0 {
		return ""
	}

	data := [][]string{}
	for _, ip := range ips {
		data = append(data, []string{
			ip.IpAddress,
			fmt.Sprintf("%.3f", ip.TotalInBytes),
			fmt.Sprintf("%.3f", ip.TotalOutBytes),
			array.JoinInt32(ip.PortList, ", "),
			ip.FirstSeenTime.UTC().Format(time.RFC3339),
			ip.ThreatTags,
			fmt.Sprintf("%v", ip.ThreatSource),
			ip.Country,
			ip.Region,
		})
	}

	return renderCustomTable(
		[]string{"IP Address", "Inbound Bytes", "Outboud Bytes", "List of Ports",
			"First Time Seen", "Threat Tags", "Threat Source", "Country", "Region",
		}, data,
		eventsTableOptions(),
	)
}

func eventFileDataHashEntitiesTable(dataHashes []api.EventFileDataHashEntity) string {
	if len(dataHashes) == 0 {
		return ""
	}

	data := [][]string{}
	for _, dHash := range dataHashes {
		knownBad := "No"
		if dHash.IsKnownBad != 0 {
			knownBad = "Yes"
		}
		data = append(data, []string{
			strings.Join(dHash.ExePathList, ", "),
			dHash.FiledataHash,
			fmt.Sprintf("%d", dHash.MachineCount),
			dHash.FirstSeenTime.UTC().Format(time.RFC3339),
			knownBad,
		})
	}

	return renderCustomTable(
		[]string{"Executable Paths", "File Hash", "Number of Machines", "First Time Seen", "Known Bad"},
		data, eventsTableOptions(),
	)
}

func eventFileExePathEntitiesTable(exePaths []api.EventFileExePathEntity) string {
	if len(exePaths) == 0 {
		return ""
	}

	data := [][]string{}
	for _, exe := range exePaths {
		data = append(data, []string{
			exe.ExePath,
			exe.FirstSeenTime.UTC().Format(time.RFC3339),
			exe.LastFiledataHash,
			exe.LastPackageName,
			exe.LastVersion,
			exe.LastFileOwner,
		})
	}

	return renderCustomTable(
		[]string{
			"Executable Path", "First Time Seen", "Last File Hash",
			"Last Package Name", "Last Version", "Last File Owner",
		},
		data, eventsTableOptions(),
	)
}

func eventProcessEntitiesTable(processes []api.EventProcessEntity) string {
	if len(processes) == 0 {
		return ""
	}

	alignLeft := tableFunc(func(t *tablewriter.Table) {
		t.SetAlignment(tablewriter.ALIGN_LEFT)
	})

	data := [][]string{}
	for _, proc := range processes {
		data = append(data, []string{
			fmt.Sprintf("%d", proc.ProcessID),
			proc.Hostname,
			proc.ProcessStartTime.UTC().Format(time.RFC3339),
			fmt.Sprintf("%.3f", proc.CpuPercentage),
			proc.Cmdline,
		})
	}

	return renderCustomTable(
		[]string{"Process ID", "Hostname", "Start Time", "CPU Percentage", "Command"},
		data, eventsTableOptions(), alignLeft,
	)
}

func eventContainerEntitiesTable(containers []api.EventContainerEntity) string {
	if len(containers) == 0 {
		return ""
	}

	data := [][]string{}
	for _, container := range containers {
		containerType := ""
		if container.IsClient != 0 {
			containerType = "Client"
		}
		if container.IsServer != 0 {
			if containerType != "" {
				containerType = "Server/Client"
			} else {
				containerType = "Server"
			}
		}
		data = append(data, []string{
			container.ImageRepo,
			container.ImageTag,
			fmt.Sprintf("%d", container.HasExternalConns),
			containerType,
			container.FirstSeenTime.UTC().Format(time.RFC3339),
			container.PodNamespace,
			container.PodIpAddr,
		})
	}
	return renderCustomTable(
		[]string{
			"Image Repo", "Image Tag", "External Connections", "Type",
			"First Time Seen", "Pod Namespace", "Pod Ipaddress",
		},
		data, eventsTableOptions(),
	)
}

func eventUserEntitiesTable(users []api.EventUserEntity) string {
	if len(users) == 0 {
		return ""
	}

	data := [][]string{}
	for _, user := range users {
		data = append(data, []string{user.Username, user.MachineHostname})
	}

	return renderCustomTable(
		[]string{"Username", "Hostname"},
		data, eventsTableOptions(),
	)
}

func eventApplicationEntitiesTable(applications []api.EventApplicationEntity) string {
	if len(applications) == 0 {
		return ""
	}

	data := [][]string{}
	for _, app := range applications {
		appType := ""
		if app.IsClient != 0 {
			appType = "Client"
		}
		if app.IsServer != 0 {
			if appType != "" {
				appType = "Server/Client"
			} else {
				appType = "Server"
			}
		}
		data = append(data, []string{
			app.Application,
			fmt.Sprintf("%d", app.HasExternalConns),
			appType,
			app.EarliestKnownTime.UTC().Format(time.RFC3339),
		})
	}

	return renderCustomTable(
		[]string{"Application", "External Connections", "Type", "Earliest Known Time"},
		data, eventsTableOptions(),
	)
}

func eventCustomRuleEntitiesTable(rules []api.EventCustomRuleEntity) string {
	if len(rules) == 0 {
		return ""
	}

	r := &strings.Builder{}
	for _, rule := range rules {
		r.WriteString(eventCustomRuleEntityTable(rule))
		r.WriteString("\n")
		r.WriteString(eventCustomRuleDisplayFilerTable(rule))
	}

	return r.String()
}

func eventCustomRuleEntityTable(rule api.EventCustomRuleEntity) string {
	autoWrapText := tableFunc(func(t *tablewriter.Table) {
		t.SetAutoWrapText(false)
	})

	return renderCustomTable(
		[]string{"Rule GUID", "Last Updated User", "Last Updated Time"},
		[][]string{[]string{
			rule.RuleGuid, rule.LastUpdatedUser, rule.LastUpdatedTime.UTC().Format(time.RFC3339),
		}},
		eventsTableOptions(), autoWrapText,
	)
}

func eventCustomRuleDisplayFilerTable(rule api.EventCustomRuleEntity) string {
	filter, err := cli.FormatJSONString(rule.DisplayFilter)
	if err != nil {
		filter = rule.DisplayFilter
	}

	tblOption := tableFunc(func(t *tablewriter.Table) {
		t.SetBorder(true)
		t.SetAutoWrapText(false)
		t.SetAlignment(tablewriter.ALIGN_LEFT)
	})

	return renderCustomTable(
		[]string{"Display Filter"},
		[][]string{[]string{filter}},
		eventsTableOptions(), tblOption,
	)
}

func eventRecIDEntitiesTable(records []api.EventRecIDEntity) string {
	if len(records) == 0 {
		return ""
	}

	data := [][]string{}
	for _, rec := range records {
		data = append(data, []string{
			rec.RecID,
			rec.AccountID,
			rec.AccountAlias,
			rec.Title,
			rec.Status,
			rec.EvalType,
			rec.EvalGuid,
		})
	}
	return renderCustomTable(
		[]string{"Record ID", "Account ID", "Account Alias",
			"Description", "Status", "Evaluation Type", "Evaluation GUID",
		},
		data, eventsTableOptions(),
	)
}

func eventViolationReasonEntitiesTable(reasons []api.EventViolationReasonEntity) string {
	if len(reasons) == 0 {
		return ""
	}

	data := [][]string{}
	for _, reason := range reasons {
		data = append(data, []string{
			reason.RecID,
			reason.Reason,
		})
	}

	return renderCustomTable(
		[]string{"Violation ID", "Reason"},
		data, eventsTableOptions(),
	)
}

func eventResourceEntitiesTable(resources []api.EventResourceEntity) string {
	if len(resources) == 0 {
		return ""
	}

	data := [][]string{}
	for _, res := range resources {
		data = append(data, []string{
			res.Name,
			fmt.Sprintf("%v", res.Value),
		})
	}

	return renderCustomTable(
		[]string{"Name", "Value"},
		data, eventsTableOptions(),
	)
}

func eventNewViolationEntitiesTable(violations []api.EventNewViolationEntity) string {
	if len(violations) == 0 {
		return ""
	}

	data := [][]string{}
	for _, v := range violations {
		data = append(data, []string{
			v.RecID,
			v.Reason,
			v.Resource,
		})
	}

	return renderCustomTable(
		[]string{"Violation ID", "Reason", "Resource"},
		data, eventsTableOptions(),
	)
}

func eventMachineEntitiesTable(machines []api.EventMachineEntity) string {
	if len(machines) == 0 {
		return ""
	}

	data := [][]string{}
	for _, m := range machines {
		data = append(data, []string{
			m.Hostname,
			m.ExternalIp,
			m.InstanceID,
			m.InstanceName,
			fmt.Sprintf("%.3f", m.CpuPercentage),
			m.InternalIpAddress,
		})
	}

	return renderCustomTable(
		[]string{"Hostname", "External IP", "Instance ID",
			"Instance Name", "CPU Percentage", "Internal Ipaddress"},
		data, eventsTableOptions(),
	)
}

func filterEventsWithSeverity(events []api.Event) []api.Event {
	if eventsCmdState.Severity == "" {
		return events
	}

	sevThreshold, sevString := eventSeverityToProperTypes(eventsCmdState.Severity)
	cli.Log.Debugw("filtering events", "threshold", sevThreshold, "severity", sevString)
	eFiltered := []api.Event{}
	for _, event := range events {
		eventSeverity, _ := eventSeverityToProperTypes(event.Severity)
		if eventSeverity <= sevThreshold {
			eFiltered = append(eFiltered, event)
		}
	}

	cli.Log.Debugw("filtered events", "events", eFiltered)
	return eFiltered
}

func eventSeverityToProperTypes(severity string) (int, string) {
	switch strings.ToLower(severity) {
	case "1", "critical":
		return 1, "Critical"
	case "2", "high":
		return 2, "High"
	case "3", "medium":
		return 3, "Medium"
	case "4", "low":
		return 4, "Low"
	case "5", "info":
		return 5, "Info"
	default:
		return 6, "Unknown"
	}
}
