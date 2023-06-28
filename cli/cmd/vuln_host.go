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

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
)

var (
	// the maximum number of packages per scan request
	manifestPkgsCap = 1000

	// the package manifest file
	pkgManifestFile string

	// automatically generate the package manifest from the local host
	pkgManifestLocal bool
)

func init() {
	// add sub-commands to the 'vulnerability host' command
	vulHostCmd.AddCommand(vulHostScanPkgManifestCmd)
	vulHostCmd.AddCommand(vulHostGenPkgManifestCmd)
	vulHostCmd.AddCommand(vulHostListCvesCmd)
	vulHostCmd.AddCommand(vulHostListHostsCmd)
	vulHostCmd.AddCommand(vulHostShowAssessmentCmd)

	setFixableFlag(
		vulHostListCvesCmd.Flags(),
		vulHostShowAssessmentCmd.Flags(),
		vulHostScanPkgManifestCmd.Flags(),
	)

	setPackagesFlag(
		vulHostListCvesCmd.Flags(),
		vulHostShowAssessmentCmd.Flags(),
		vulHostScanPkgManifestCmd.Flags(),
	)

	setDetailsFlag(
		vulHostShowAssessmentCmd.Flags(),
	)

	setSeverityFlag(
		vulHostListCvesCmd.Flags(),
		vulHostShowAssessmentCmd.Flags(),
	)

	setFailOnSeverityFlag(
		vulHostShowAssessmentCmd.Flags(),
		vulHostScanPkgManifestCmd.Flags(),
	)

	setFailOnFixableFlag(
		vulHostShowAssessmentCmd.Flags(),
		vulHostScanPkgManifestCmd.Flags(),
	)

	setActiveFlag(
		vulHostShowAssessmentCmd.Flags(),
		vulHostListCvesCmd.Flags(),
	)

	setCsvFlag(
		vulHostListCvesCmd.Flags(),
		vulHostListHostsCmd.Flags(),
		vulHostShowAssessmentCmd.Flags(),
	)

	setTimeRangeFlags(
		vulHostListHostsCmd.Flags(),
	)

	// the package manifest file
	vulHostScanPkgManifestCmd.Flags().StringVarP(&pkgManifestFile,
		"file", "f", "",
		"path to a package manifest to scan",
	)

	// automatically generate the package manifest from the local host
	vulHostScanPkgManifestCmd.Flags().BoolVarP(&pkgManifestLocal,
		"local", "l", false,
		"automatically generate the package manifest from the local host",
	)

	// the collector_type of the assessment
	vulHostShowAssessmentCmd.Flags().StringVar(&vulCmdState.CollectorType,
		"collector_type", vulnHostCollectorTypeAgentless,
		"filter assessments by collector type(Agent/Agentless)",
	)
}

func cvesSummary(hosts []api.VulnerabilityHost) map[string]VulnCveSummary {
	uniqueCves := make(map[string]VulnCveSummary)
	for _, host := range hosts {
		if host.VulnID == "" {
			continue
		}

		if v, ok := uniqueCves[host.VulnID]; ok {
			if array.ContainsStr(v.Hostnames, host.EvalCtx.Hostname) {
				continue
			}

			v.Count++
			v.Hostnames = append(v.Hostnames, host.EvalCtx.Hostname)
			uniqueCves[host.VulnID] = v
			continue
		}
		sum := VulnCveSummary{Host: host, Count: 1}
		sum.Hostnames = append(sum.Hostnames, sum.Host.EvalCtx.Hostname)
		uniqueCves[host.VulnID] = sum
	}
	return uniqueCves
}

func aggregatePackagesWithHosts(slice []packageTable, s packageTable, host bool, hasFix bool) []packageTable {
	for i, item := range slice {
		// if packages are the same, group together
		if item.equals(s) {
			slice[i].cveCount++
			if host {
				slice[i].hostCount++
			}
			if hasFix {
				slice[i].fixes++
			}
			return slice
		}
	}
	return append(slice, s)
}

func filterHostCVEsTable(cves map[string]VulnCveSummary) (map[string]VulnCveSummary, string) {
	var out map[string]VulnCveSummary
	var filteredCves = 0
	var totalCves = 0

	out, filtered, total := filterHostVulnCVEs(cves)
	filteredCves += filtered
	totalCves += total

	if filteredCves > 0 {
		showing := totalCves - filteredCves
		return out, fmt.Sprintf("\n%d of %d cve(s) showing\n", showing, totalCves)
	}

	return out, ""
}

func filterHostVulnCVEs(cves map[string]VulnCveSummary) (map[string]VulnCveSummary, int, int) {
	var (
		filtered = 0
		total    = 0
		out      = make(map[string]VulnCveSummary)
	)

	for _, c := range cves {
		cve := c.Host
		total++
		// if the user wants to show only vulnerabilities of active packages
		if vulCmdState.Active && cve.PackageActive() == "" {
			filtered++
			continue
		}
		if vulCmdState.Fixable && (cve.FixInfo.FixAvailable == "" || cve.FixInfo.FixAvailable == "0") {
			filtered++
			continue
		}

		if vulCmdState.Severity != "" {
			if filterSeverity(cve.Severity, vulCmdState.Severity) {
				filtered++
				continue
			}
		}
		out[cve.VulnID] = c
	}

	return out, filtered, total
}

func hostVulnAssessmentToCountsTable(counts api.HostVulnCounts) [][]string {
	return [][]string{
		{"Critical", fmt.Sprint(counts.Critical),
			fmt.Sprint(counts.CritFixable)},
		{"High", fmt.Sprint(counts.High),
			fmt.Sprint(counts.HighFixable)},
		{"Medium", fmt.Sprint(counts.Medium),
			fmt.Sprint(counts.MedFixable)},
		{"Low", fmt.Sprint(counts.Low),
			fmt.Sprint(counts.LowFixable)},
		{"Info", fmt.Sprint(counts.Info),
			fmt.Sprint(counts.InfoFixable)},
	}
}

func buildHostVulnCVEsToTableError() string {
	msg := "There are no"
	if vulCmdState.Fixable {
		msg = fmt.Sprintf("%s fixable", msg)
	}

	if vulCmdState.Severity != "" {
		msg = fmt.Sprintf("%s %s", msg, vulCmdState.Severity)
	}

	msg = fmt.Sprintf("%s vulnerabilities", msg)

	if vulCmdState.Active {
		msg = fmt.Sprintf("%s of packages actively running", msg)
	}
	return fmt.Sprintf("%s in your environment.\n", msg)
}

func summaryToHostList(sum map[string]VulnCveSummary) (hosts []api.VulnerabilityHost) {
	for _, v := range sum {
		hosts = append(hosts, v.Host)
	}
	return
}

type VulnCveSummary struct {
	Host      api.VulnerabilityHost
	Count     int
	Hostnames []string
}
