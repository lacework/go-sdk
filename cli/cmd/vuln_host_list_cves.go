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

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/lacework/go-sdk/lwseverity"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// vulHostListCvesCmd represents the 'lacework vuln host list-cves' command
	vulHostListCvesCmd = &cobra.Command{
		Use:  "list-cves",
		Args: cobra.NoArgs,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if vulCmdState.Csv {
				cli.EnableCSVOutput()
			}
			return nil
		},
		Short: "List the CVEs found in the hosts in your environment",
		Long: `List the CVEs found in the hosts in your environment.

Filter results to only show vulnerabilities actively running in your environment
with fixes:

    lacework vulnerability host list-cves --active --fixable`,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}

			cli.StartProgress("Fetching CVEs in your environment...")
			response, err := cli.LwApi.V2.Vulnerabilities.Hosts.SearchAllPages(api.SearchFilter{})
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get CVEs from hosts")
			}

			if err := buildListCVEReports(response.Data); err != nil {
				return err
			}
			return nil
		},
	}
)

// Build the cli output for vuln host list-cves
func buildListCVEReports(cves []api.VulnerabilityHost) error {
	uniqueCves := cvesSummary(cves)
	filteredCves, filtered := filterHostCVEsTable(uniqueCves)

	if cli.JSONOutput() {
		if filteredCves == nil {
			if err := cli.OutputJSON(buildHostVulnCVEsToTableError()); err != nil {
				return err
			}
		} else {
			// fix here too
			if err := cli.OutputJSON(summaryToHostList(filteredCves)); err != nil {
				return err
			}
		}
		return nil
	}

	if len(cves) == 0 {
		// @afiune add a helpful message, possible things are:
		// 1) host vuln feature is not enabled on the account
		// 2) user doesn't have agents deployed
		// 3) there are actually NO vulnerabilities on any host
		cli.OutputHuman("There are no vulnerabilities on any host in your environment.\n")
		return nil
	}

	// packages output
	if vulCmdState.Packages {
		packages, filteredPackages := hostVulnListCvesPackagesTable(cves)
		if cli.CSVOutput() {

			// order by cve count

			return cli.OutputCSV(
				[]string{"CVE Count", "Highest", "Package", "Current Version", "Fix Version", "Pkg Status", "Hosts"},
				packages,
			)
		}
		vulnListCvesPackagesOutput(packages, filteredPackages)
		return nil
	}

	rows := hostVulnCVEsTable(filteredCves)
	if len(rows) == 0 {
		cli.OutputHuman(buildHostVulnCVEsToTableError())
		return nil
	}

	if cli.CSVOutput() {
		return cli.OutputCSV(
			[]string{"CVE ID", "Severity", "CvssV2", "CvssV3", "Package", "Current Version",
				"Fix Version", "OS Version", "Hosts", "Pkg Status", "Vuln Status"},
			rows,
		)
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"CVE ID", "Severity", "CvssV2", "CvssV3", "Package", "Current Version",
				"Fix Version", "OS Version", "Hosts", "Pkg Status", "Vuln Status"},
			rows,
		),
	)

	if filtered != "" {
		cli.OutputHuman(filtered)
	}

	if !vulCmdState.Active {
		cli.OutputHuman(
			"\nTry adding '--active' to only show vulnerabilities of packages actively running.\n",
		)
	} else if !vulCmdState.Fixable {
		cli.OutputHuman(
			"\nTry adding '--fixable' to only show fixable vulnerabilities.\n",
		)
	}
	return nil
}

func vulnListCvesPackagesOutput(packages [][]string, filteredPackagesMsg string) {
	// sort by highest cve count
	sort.Slice(packages, func(i, j int) bool {
		return stringToInt(packages[i][0]) > stringToInt(packages[j][0])
	})

	cli.OutputHuman(
		renderSimpleTable(
			[]string{
				"CVE Count",
				"Highest Severity",
				"Package",
				"Current Version",
				"Fix Version",
				"Pkg Status",
				"Hosts Impacted",
			},
			packages,
		),
	)
	if filteredPackagesMsg != "" {
		cli.OutputHuman(filteredPackagesMsg)
	}
}

func hostVulnListCvesPackagesTable(cves []api.VulnerabilityHost) ([][]string, string) {
	var (
		out                [][]string
		filteredPackages   []string
		aggregatedPackages []packageTable
	)

	// Get all unique package names
	var packageNames []string
	for _, c := range cves {
		if c.VulnID != "" {
			packageNames = append(packageNames, c.FeatureKey.Name)
		}
	}

	var uniquePackageNames []string = array.Unique(packageNames)
	var added []string

	for _, u := range uniquePackageNames {
		var (
			pack              packageTable
			cveIDs            []string
			hosts             []string
			severities        []lwseverity.Severity
			active            string
			packageIdentifier string
		)
		for _, host := range cves {
			packageIdentifier = fmt.Sprintf("%s-%s", host.VulnID, host.FeatureKey.VersionInstalled)
			if host.FeatureKey.Name == u {
				if host.PackageActive() == "ACTIVE" {
					active = "ACTIVE"
				}

				if array.ContainsStr(added, host.FeatureKey.Name) {
					if host.Severity != "" {
						cveIDs = append(cveIDs, packageIdentifier)
						severities = append(severities, lwseverity.NewSeverity(host.Severity))
						hosts = append(hosts, host.EvalCtx.Hostname)
					}
					continue
				}

				pack = packageTable{
					severity:       cases.Title(language.English).String(host.Severity),
					packageName:    host.FeatureKey.Name,
					currentVersion: host.FeatureKey.VersionInstalled,
					fixVersion:     host.FixInfo.FixedVersion,
				}

				added = append(added, host.FeatureKey.Name)
				cveIDs = append(cveIDs, fmt.Sprintf("%s-%s", host.VulnID, host.FeatureKey.VersionInstalled))
				severities = append(severities, lwseverity.NewSeverity(host.Severity))
				hosts = append(hosts, host.EvalCtx.Hostname)
			}
		}

		// Exclude filtered packages and packages without vulns
		if !array.ContainsStr(filteredPackages, packageIdentifier) && len(cveIDs) > 0 {
			var unqCves []string = array.Unique(cveIDs)
			var unqHosts []string = array.Unique(hosts)
			pack.packageStatus = active
			pack.cveCount = len(unqCves)
			pack.hostCount = len(unqHosts)

			// set highest known severity of the package
			if len(severities) > 0 {
				lwseverity.SortSlice(severities)
				pack.severity = severities[0].GetSeverity()
			}
			aggregatedPackages = append(aggregatedPackages, pack)
		}
	}

	for _, p := range aggregatedPackages {
		// apply package filters
		if vulCmdState.Active && p.packageStatus == "" {
			filteredPackages = append(filteredPackages, p.packageName)
			continue
		}

		if vulCmdState.Fixable && p.fixVersion == "" {
			filteredPackages = append(filteredPackages, p.packageName)
			continue
		}

		if vulCmdState.Severity != "" {
			if p.severity == "Unknown" {
				continue
			}
			if lwseverity.ShouldFilter(p.severity, vulCmdState.Severity) {
				filteredPackages = append(filteredPackages, p.packageName)
				continue
			}
		}

		output := []string{
			strconv.Itoa(p.cveCount),
			p.severity,
			p.packageName,
			p.currentVersion,
			p.fixVersion,
			p.packageStatus}
		if p.hostCount > 0 {
			output = append(output, strconv.Itoa(p.hostCount))
		}
		out = append(out, output)
	}

	filteredOutput := fmt.Sprintf("%d of %d package(s) showing\n", len(out), len(uniquePackageNames))
	return out, filteredOutput
}

func hostVulnCVEsTable(hostSummary map[string]VulnCveSummary) [][]string {
	var out [][]string
	for _, sum := range hostSummary {
		host := sum.Host
		out = append(out, []string{
			host.VulnID,
			host.Severity,
			host.CvssV2(),
			host.CvssV3(),
			host.FeatureKey.Name,
			host.FeatureKey.VersionInstalled,
			host.FixInfo.FixedVersion,
			host.FeatureKey.Namespace,
			strconv.Itoa(sum.Count),
			host.PackageActive(),
			host.Status,
		})
	}

	// order by the total number of host
	sort.Slice(out, func(i, j int) bool {
		return stringToInt(out[i][8]) > stringToInt(out[j][8])
	})

	return out
}
