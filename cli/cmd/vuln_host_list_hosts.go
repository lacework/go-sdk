//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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
	"strconv"
	"strings"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// vulHostListHostsCmd represents the 'lacework vuln host list-hosts <cve-id>' command
	vulHostListHostsCmd = &cobra.Command{
		Use:  "list-hosts <cve_id>",
		Args: cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if vulCmdState.Csv {
				cli.EnableCSVOutput()
			}
			return nil
		},
		Short: "List the hosts that contain a specified CVE ID in your environment",
		Long: `List the hosts that contain a specified CVE ID in your environment.

To list the CVEs found in the hosts of your environment run:

    lacework vulnerability host list-cves`,
		RunE: func(_ *cobra.Command, args []string) error {
			filter := api.SearchFilter{Filters: []api.Filter{{
				Expression: "eq",
				Field:      "vulnId",
				Value:      args[0],
			}}}

			cli.StartProgress("Fetching Hosts...")
			response, err := cli.LwApi.V2.Vulnerabilities.Hosts.SearchAllPages(filter)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get hosts with CVE "+args[0])
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Data)
			}

			if len(response.Data) == 0 {
				cli.OutputHuman("There are no hosts in your environment with the CVE id '%s'\n", args[0])
				return nil
			}

			rows := hostVulnHostsTable(response.Data)
			if cli.CSVOutput() {
				return cli.OutputCSV(
					[]string{"Machine ID", "Hostname", "External IP", "Internal IP",
						"Os/Arch", "Provider", "Instance ID", "Vulnerabilities", "Status"},
					rows,
				)
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"Machine ID", "Hostname", "External IP", "Internal IP",
						"Os/Arch", "Provider", "Instance ID", "Vulnerabilities", "Status"},
					rows,
				),
			)
			return nil
		},
	}
)

func hostVulnHostsTable(hosts []api.VulnerabilityHost) [][]string {
	var out [][]string
	hostSummary := hostsSummary(hosts)
	for _, sum := range hostSummary {
		host := sum.host
		summary := severitySummary(sum.severity, sum.fixable)
		out = append(out, []string{
			strconv.Itoa(host.Mid),
			host.EvalCtx.Hostname,
			host.MachineTags.ExternalIP,
			host.MachineTags.InternalIP,
			fmt.Sprintf("%s/%s", host.MachineTags.Os, host.MachineTags.Arch),
			host.MachineTags.VMProvider,
			host.MachineTags.InstanceID,
			summary,
			host.Status,
		})
	}

	return out
}

func severitySummary(severities []string, fixable int) string {
	summary := &strings.Builder{}
	sevSummaries := make(map[string]int)
	for _, s := range severities {
		switch s {
		case "Critical":
			if v, ok := sevSummaries["Critical"]; ok {
				sevSummaries["Critical"] = v + 1
			}
			sevSummaries["Critical"] = 1
		case "High":
			if v, ok := sevSummaries["High"]; ok {
				sevSummaries["High"] = v + 1
			}
			sevSummaries["High"] = 1
		case "Medium":
			if v, ok := sevSummaries["Medium"]; ok {
				sevSummaries["Medium"] = v + 1
			}
			sevSummaries["Medium"] = 1
		case "Low":
			if v, ok := sevSummaries["Low"]; ok {
				sevSummaries["Low"] = v + 1
			}
			sevSummaries["Low"] = 1
		case "Info":
			if v, ok := sevSummaries["Info"]; ok {
				sevSummaries["Info"] = v + 1
			}
			sevSummaries["Info"] = 1
		}
	}

	var keys []string
	for k := range sevSummaries {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return severityOrder(keys[i]) < severityOrder(keys[j])
	})

	for _, k := range keys {
		v := sevSummaries[k]
		summary.WriteString(fmt.Sprintf(" %d %s", v, k))
	}

	if fixable != 0 {
		summary.WriteString(fmt.Sprintf(" %d Fixable", fixable))
	}
	return summary.String()
}

func hostsSummary(hosts []api.VulnerabilityHost) map[int]vulnSummary {
	uniqueHosts := make(map[int]vulnSummary)
	for _, host := range hosts {
		if v, ok := uniqueHosts[host.Mid]; ok {
			v.severity = append(v.severity, host.Severity)
			if host.FixInfo.FixAvailable != "" && host.FixInfo.FixAvailable != "0" {
				v.fixable++
			}
			uniqueHosts[host.Mid] = v
			continue
		}

		sum := vulnSummary{host: host}
		sum.severity = append(sum.severity, host.Severity)
		if host.FixInfo.FixAvailable != "" && host.FixInfo.FixAvailable != "0" {
			sum.fixable++
		}
		uniqueHosts[host.Mid] = sum
	}
	return uniqueHosts
}

type vulnSummary struct {
	host     api.VulnerabilityHost
	severity []string
	fixable  int
}
