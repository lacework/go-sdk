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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// the package manifest file
	pkgManifestFile string

	// automatically generate the package manifest from the local host
	pkgManifestLocal bool

	vulHostGenPkgManifestCmd = &cobra.Command{
		Use:   "generate-pkg-manifest",
		Args:  cobra.NoArgs,
		Short: "generates a package-manifest from the local host",
		Long: `Generates a package-manifest formatted for usage with the Lacework
scan package-manifest API.

Additionally, you can automatically generate a package-manifest from
the local host and send it directly to the Lacework API with the command:

    $ lacework vulnerability host scan-pkg-manifest --local`,
		RunE: func(_ *cobra.Command, _ []string) error {
			manifest, err := cli.GeneratePackageManifest()
			if err != nil {
				return errors.Wrap(err, "unable to generate package manifest")
			}

			return cli.OutputJSON(manifest)
		},
	}

	vulHostScanPkgManifestCmd = &cobra.Command{
		Use:   "scan-pkg-manifest <manifest>",
		Args:  cobra.MaximumNArgs(1),
		Short: "request an on-demand host vulnerability assessment from a package-manifest",
		Long: `Request an on-demand host vulnerability assessment of your software packages to
determine if the packages contain any common vulnerabilities and exposures.

Simple usage:

    $ lacework vulnerability host scan-pkg-manifest '{
        "os_pkg_info_list": [
            {
                "os":"Ubuntu",
                "os_ver":"18.04",
                "pkg": "openssl",
                "pkg_ver": "1.1.1-1ubuntu2.1~18.04.5"
            }
        ]
    }'

To generate a package-manifest from the local host and scan it automatically:

    $ lacework vulnerability host scan-pkg-manifest --local

(*) NOTE:
 - Only packages managed by a package manager for supported OS's are reported.
 - Calls to this operation are rate limited to 10 calls per hour, per access key.
 - This operation is limited to 1k of packages per payload. If you require a payload
   larger than 1k, you must make multiple requests.`,
		RunE: func(_ *cobra.Command, args []string) error {
			var pkgManifest = ""
			if len(args) != 0 && args[0] != "" {
				pkgManifest = args[0]
			} else if pkgManifestFile != "" {
				pkgManifestBytes, err := ioutil.ReadFile(pkgManifestFile)
				if err != nil {
					return errors.Wrap(err, "unable to read file")
				}
				pkgManifest = string(pkgManifestBytes)
			} else if pkgManifestLocal {
				manifest, err := cli.GeneratePackageManifest()
				if err != nil {
					return errors.Wrap(err, "unable to generate package manifest")
				}
				manifestString, err := json.Marshal(&manifest)
				if err != nil {
					panic(err)
				}

				pkgManifest = string(manifestString)
			} else {
				// avoid asking for a confirmation before launching the editor
				prompt := &survey.Editor{
					Message:  "Provide a package manifest to scan",
					FileName: "pkg-manifest*.json",
				}
				err := survey.AskOne(prompt, &pkgManifest)
				if err != nil {
					return err
				}
			}

			response, err := cli.LwApi.Vulnerabilities.Host.Scan(pkgManifest)
			if err != nil {
				return errors.Wrap(err, "unable to request an on-demand host vulnerability scan")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			if len(response.Vulns) == 0 {
				// @afiune add a helpful message, possible things are:
				cli.OutputHuman("There are no vulnerabilities found.\n")
				return nil
			}

			cli.OutputHuman(hostScanPackagesVulnToTable(&response))
			return nil
		},
	}

	vulHostListCvesCmd = &cobra.Command{
		Use:   "list-cves",
		Args:  cobra.NoArgs,
		Short: "list the CVEs found in the hosts in your environment",
		Long: `List the CVEs found in the hosts in your environment.

Filter results to only show vulnerabilities actively running in your environment
with fixes:

    $ lacework vulnerability host list-cves --active --fixable`,
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Vulnerabilities.Host.ListCves()
			if err != nil {
				return errors.Wrap(err, "unable to get CVEs from hosts")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.CVEs)
			}

			if len(response.CVEs) == 0 {
				// @afiune add a helpful message, possible things are:
				// 1) host vuln feature is not enabled on the account
				// 2) user doesn't have agents deployed
				// 3) there are actually NO vulnerabilities on any host
				cli.OutputHuman("There are no vulnerabilities on any host in your environment.\n")
				return nil
			}

			if vulCmdState.Packages {
				cli.OutputHuman(hostVulnCVEsPackagesSummary(response.CVEs, true))
			} else {
				cli.OutputHuman(hostVulnCVEsToTable(response.CVEs))
			}

			return nil
		},
	}

	vulHostListHostsCmd = &cobra.Command{
		Use:   "list-hosts <cve_id>",
		Args:  cobra.ExactArgs(1),
		Short: "list the hosts that contain a specified CVE id in your environment",
		Long: `List the hosts that contain a specified CVE id in your environment.

To list the CVEs found in the hosts of your environment run:

    $ lacework vulnerability host list-cves`,
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Vulnerabilities.Host.ListHostsWithCVE(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get hosts with CVE "+args[0])
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Hosts)
			}

			if len(response.Hosts) == 0 {
				// @afiune add a helpful message, possible things are:
				// 1) host vuln feature is not enabled on the account
				// 2) user doesn't have agents deployed
				// 3) there are actually NO vulnerabilities on any host
				cli.OutputHuman("There are no hosts in your environment with the CVE id '%s'\n", args[0])
				return nil
			}

			cli.OutputHuman(hostVulnHostsToTable(response.Hosts))
			return nil
		},
	}

	vulHostShowAssessmentCmd = &cobra.Command{
		Use:     "show-assessment <machine_id>",
		Aliases: []string{"show"},
		Args:    cobra.ExactArgs(1),
		Short:   "show results of a host vulnerability assessment",
		Long: `Show results of a host vulnerability assessment.

To find the machine id from hosts in your environment, use the command:

    $ lacework vulnerability host list-cves

Grab a CVE id and feed it to the command:

    $ lacework vulnerability host list-hosts my_cve_id`,
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Vulnerabilities.Host.GetHostAssessment(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get host assessment with id "+args[0])
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Assessment)
			}

			cli.OutputHuman(hostVulnHostDetailsToTable(response.Assessment))
			return nil
		},
	}

	// @afiune this is not yet supported since there is no external API available
	vulHostListAssessmentsCmd = &cobra.Command{
		Use:    "list-assessments",
		Hidden: true,
		//Aliases: []string{"list", "ls"},
		Short: "list host vulnerability assessments from a time range",
		Long:  "List host vulnerability assessments from a time range.",
		RunE: func(_ *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	// add sub-commands to the 'vulnerability host' command
	vulHostCmd.AddCommand(vulHostScanPkgManifestCmd)
	vulHostCmd.AddCommand(vulHostGenPkgManifestCmd)
	vulHostCmd.AddCommand(vulHostListAssessmentsCmd)
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

	setActiveFlag(
		vulHostShowAssessmentCmd.Flags(),
		vulHostListCvesCmd.Flags(),
	)

	// add online flag to host list-hosts command
	vulHostListHostsCmd.Flags().BoolVar(&vulCmdState.Online,
		"online", false, "only show hosts that are online",
	)
	// add offline flag to host list-hosts command
	vulHostListHostsCmd.Flags().BoolVar(&vulCmdState.Offline,
		"offline", false, "only show hosts that are offline",
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
}

func hostVulnHostsToTable(hosts []api.HostVulnDetail) string {
	var (
		tableBuilder = &strings.Builder{}
		t            = tablewriter.NewWriter(tableBuilder)
		rows         = hostVulnHostsTable(hosts)
	)

	// if the user wants to show only online or
	// offline hosts, show a friendly message
	if len(rows) == 0 {
		if vulCmdState.Online {
			return "There are no online hosts.\n"
		}
		if vulCmdState.Offline {
			return "There are no offline hosts.\n"
		}
	}

	t.SetHeader([]string{
		"Machine ID",
		"Hostname",
		"External IP",
		"Internal IP",
		"Os/Arch",
		"Provider",
		"Instance ID",
		"Vulnerabilities",
		"Status",
	})
	t.SetBorder(false)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.AppendBulk(rows)
	t.Render()

	return tableBuilder.String()
}

func hostVulnHostsTable(hosts []api.HostVulnDetail) [][]string {
	out := [][]string{}
	for _, host := range hosts {

		// filter by machine status: Online / Offline
		if vulCmdState.Online && host.Details.MachineStatus != "Online" {
			continue
		}
		if vulCmdState.Offline && host.Details.MachineStatus != "Offline" {
			continue
		}

		hostVulnSummary, _ := hostVulnSummaryFromHostDetail(&host.Summary)

		out = append(out, []string{
			host.Details.MachineID,
			host.Details.Hostname,
			host.Details.Tags.ExternalIP,
			host.Details.Tags.InternalIP,
			fmt.Sprintf("%s/%s", host.Details.Tags.Os, host.Details.Tags.Arch),
			host.Details.Tags.VmProvider,
			host.Details.Tags.InstanceID,
			hostVulnSummary,
			host.Details.MachineStatus,
		})
	}

	return out
}

func hostVulnSummaryFromHostDetail(hostVulnSummary *api.HostVulnCveSummary) (string, bool) {
	summary := []string{}
	hostVulnCounts := hostVulnSummary.Severity.VulnerabilityCounts()

	summary = addToHostSummary(summary, hostVulnCounts.Critical, "Critical")
	summary = addToHostSummary(summary, hostVulnCounts.High, "High")
	summary = addToHostSummary(summary, hostVulnCounts.Medium, "Medium")
	summary = addToHostSummary(summary, hostVulnCounts.Low, "Low")
	summary = addToHostSummary(summary, hostVulnCounts.Negligible, "Negligible")

	if len(summary) == 0 {
		return fmt.Sprintf("None! Time for %s", randomEmoji()), false
	}

	if hostVulnCounts.TotalFixable != 0 {
		summary = append(summary, fmt.Sprintf("%d Fixable", hostVulnCounts.TotalFixable))
	}

	return strings.Join(summary, " "), true
}

func hostVulnCVEsPackagesSummary(cves []api.HostVulnCVE, withHosts bool) string {
	var (
		tableBuilder = &strings.Builder{}
		t            = tablewriter.NewWriter(tableBuilder)
	)

	headers := []string{
		"CVE Count",
		"Severity",
		"Package",
		"Current Version",
		"Fix Version",
		"Pkg Status",
	}
	if withHosts {
		headers = append(headers, "Hosts")
	}
	t.SetHeader(headers)
	t.SetBorder(false)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.AppendBulk(hostVulnPackagesTable(cves, withHosts))
	t.Render()

	return tableBuilder.String()
}

func hostVulnPackagesTable(cves []api.HostVulnCVE, withHosts bool) [][]string {
	out := [][]string{}
	for _, cve := range cves {
		for _, pkg := range cve.Packages {
			// if the user wants to show only vulnerabilities of acive packages
			if vulCmdState.Active && pkg.PackageStatus == "" {
				continue
			}
			if vulCmdState.Fixable && pkg.FixedVersion == "" {
				continue
			}

			added := false
			for i := range out {
				if out[i][1] == strings.Title(pkg.Severity) &&
					out[i][2] == pkg.Name &&
					out[i][3] == pkg.Version &&
					out[i][4] == pkg.FixedVersion &&
					out[i][5] == pkg.PackageStatus {

					if countCVEs, err := strconv.Atoi(out[i][0]); err == nil {
						out[i][0] = fmt.Sprintf("%d", (countCVEs + 1))
						added = true
					}

					if withHosts {
						if countHosts, err := strconv.Atoi(out[i][6]); err == nil {
							prevCount := stringToInt(pkg.HostCount)
							out[i][6] = fmt.Sprintf("%d", (countHosts + prevCount))
							added = true
						}
					}

				}
			}

			if added {
				continue
			}

			row := []string{
				"1",
				strings.Title(pkg.Severity),
				pkg.Name,
				pkg.Version,
				pkg.FixedVersion,
				pkg.PackageStatus,
			}
			if withHosts {
				row = append(row, pkg.HostCount)
			}
			out = append(out, row)
		}
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
}

func hostVulnCVEsToTable(cves []api.HostVulnCVE) string {
	var (
		tableBuilder = &strings.Builder{}
		t            = tablewriter.NewWriter(tableBuilder)
		rows         = hostVulnCVEsTable(cves)
	)

	// if the user wants to show only online or
	// offline hosts, show a friendly message
	if len(rows) == 0 {
		return buildHostVulnCVEsToTableError()
	}

	t.SetHeader([]string{
		"CVE",
		"Severity",
		"Score",
		"Package",
		"Current Version",
		"Fix Version",
		"OS Version",
		"Hosts",
		"Pkg Status",
		"Vuln Status",
	})
	t.SetBorder(false)
	t.AppendBulk(rows)
	t.Render()

	if !vulCmdState.Active {
		tableBuilder.WriteString(
			"\nTry adding '--active' to only show vulnerabilities of packages actively running.\n",
		)
	} else if !vulCmdState.Fixable {
		tableBuilder.WriteString(
			"\nTry adding '--fixable' to only show fixable vulnerabilities.\n",
		)
	}

	return tableBuilder.String()
}

func hostVulnCVEsTable(cves []api.HostVulnCVE) [][]string {
	out := [][]string{}
	out = append(out, hostVulnCVEsTableForSeverity(cves, "Critical")...)
	out = append(out, hostVulnCVEsTableForSeverity(cves, "High")...)
	out = append(out, hostVulnCVEsTableForSeverity(cves, "Medium")...)
	out = append(out, hostVulnCVEsTableForSeverity(cves, "Low")...)
	//out = append(out, hostVulnCVEsTableForSeverity(cves, "Info")...)
	out = append(out, hostVulnCVEsTableForSeverity(cves, "Negligible")...)
	return out
}

func hostVulnCVEsTableForSeverity(cves []api.HostVulnCVE, severity string) [][]string {
	out := [][]string{}
	for _, cve := range cves {
		for _, pkg := range cve.Packages {
			// if the user wants to show only vulnerabilities of acive packages
			if vulCmdState.Active && pkg.PackageStatus == "" {
				continue
			}
			if vulCmdState.Fixable && pkg.FixedVersion == "" {
				continue
			}

			if pkg.Severity == severity {
				out = append(out, []string{
					cve.ID,
					pkg.Severity,
					pkg.CvssScore,
					pkg.Name,
					pkg.Version,
					pkg.FixedVersion,
					pkg.Namespace,
					pkg.HostCount,
					pkg.PackageStatus,
					pkg.Status,
				})
			}
		}
	}

	// order by total number of host
	sort.Slice(out, func(i, j int) bool {
		return stringToInt(out[i][7]) > stringToInt(out[j][7])
	})

	return out
}

func hostVulnHostDetailsToTable(assessment api.HostVulnHostAssessment) string {
	var (
		tableBuilder        = &strings.Builder{}
		hostDetailsTable    = &strings.Builder{}
		hostVulnCountsTable = &strings.Builder{}
		t                   = tablewriter.NewWriter(hostDetailsTable)
	)

	t.SetColumnSeparator("")
	t.SetBorder(false)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.AppendBulk(
		[][]string{
			[]string{"Machine ID", assessment.Host.MachineID},
			[]string{"Hostname", assessment.Host.Hostname},
			[]string{"External IP", assessment.Host.Tags.ExternalIP},
			[]string{"Internal IP", assessment.Host.Tags.InternalIP},
			[]string{"Os", assessment.Host.Tags.Os},
			[]string{"Arch", assessment.Host.Tags.Arch},
			[]string{"Namespace", getNamespaceFromHostVuln(assessment.CVEs)},
			[]string{"Provider", assessment.Host.Tags.VmProvider},
			[]string{"Instance ID", assessment.Host.Tags.InstanceID},
			[]string{"AMI", assessment.Host.Tags.AmiID},
		},
	)
	t.Render()

	t = tablewriter.NewWriter(hostVulnCountsTable)
	t.SetBorder(false)
	t.SetColumnSeparator(" ")
	t.SetHeader([]string{
		"Severity", "Count", "Fixable",
	})
	t.AppendBulk(hostVulnAssessmentToCountsTable(assessment.VulnerabilityCounts()))
	t.Render()

	t = tablewriter.NewWriter(tableBuilder)
	t.SetBorder(false)
	t.SetAutoWrapText(false)
	t.SetHeader([]string{
		"Host Details",
		"Vulnerabilities",
	})
	t.Append([]string{
		hostDetailsTable.String(),
		hostVulnCountsTable.String(),
	})
	t.Render()

	if vulCmdState.Details || vulCmdState.Fixable || vulCmdState.Packages || vulCmdState.Active {
		if vulCmdState.Packages {
			tableBuilder.WriteString(hostVulnCVEsPackagesSummary(assessment.CVEs, false))
		} else {
			tableBuilder.WriteString(hostVulnHostAssessmentCVEsToTable(assessment))
		}
		tableBuilder.WriteString("\n")
	}

	if !vulCmdState.Details && !vulCmdState.Active && !vulCmdState.Fixable && !vulCmdState.Packages {
		tableBuilder.WriteString(
			"Try adding '--details' to increase details shown about the vulnerability assessment.\n",
		)
	} else if !vulCmdState.Active {
		tableBuilder.WriteString(
			"Try adding '--active' to only show vulnerabilities of packages actively running.\n",
		)
	} else if !vulCmdState.Fixable {
		tableBuilder.WriteString(
			"Try adding '--fixable' to only show fixable vulnerabilities.\n",
		)
	}

	return tableBuilder.String()
}

func hostVulnHostAssessmentCVEsToTable(assessment api.HostVulnHostAssessment) string {
	var (
		tableBuilder = &strings.Builder{}
		t            = tablewriter.NewWriter(tableBuilder)
		rows         = hostVulnCVEsTableForHostView(assessment.CVEs)
	)

	// if the user wants to show only vulnerabilities of active packages
	// and we don't have any, show a friendly message
	if len(rows) == 0 {
		if vulCmdState.Active && vulCmdState.Fixable {
			return "There are no fixable vulnerabilities with packages actively running in this host.\n"
		}
		if vulCmdState.Active {
			return "There are no vulnerabilities with packages actively running in this host.\n"
		}
		if vulCmdState.Active {
			return "There are no fixable vulnerabilities in this host.\n"
		}
	}

	t.SetHeader([]string{
		"CVE",
		"Severity",
		"Score",
		"Package",
		"Current Version",
		"Fix Version",
		"Pgk Status",
		"Vuln Status",
	})
	t.SetBorder(false)
	t.AppendBulk(rows)
	t.Render()

	return tableBuilder.String()
}

func hostVulnCVEsTableForHostView(cves []api.HostVulnCVE) [][]string {
	out := [][]string{}
	for _, cve := range cves {
		for _, pkg := range cve.Packages {
			// if the user wants to show only vulnerabilities of acive packages
			if vulCmdState.Active && pkg.PackageStatus == "" {
				continue
			}

			if vulCmdState.Fixable && pkg.FixedVersion == "" {
				continue
			}

			out = append(out, []string{
				cve.ID,
				pkg.Severity,
				pkg.CvssScore,
				pkg.Name,
				pkg.Version,
				pkg.FixedVersion,
				pkg.PackageStatus,
				pkg.VulnerabilityStatus,
			})
		}
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
}

func getNamespaceFromHostVuln(cves []api.HostVulnCVE) string {
	namespace := ""
	for _, cve := range cves {
		for _, pkg := range cve.Packages {
			if namespace != pkg.Namespace {
				namespace = pkg.Namespace
			}
		}
	}
	return namespace
}

func hostVulnAssessmentToCountsTable(counts api.HostVulnCounts) [][]string {
	return [][]string{
		[]string{"Critical", fmt.Sprint(counts.Critical),
			fmt.Sprint(counts.CritFixable)},
		[]string{"High", fmt.Sprint(counts.High),
			fmt.Sprint(counts.HighFixable)},
		[]string{"Medium", fmt.Sprint(counts.Medium),
			fmt.Sprint(counts.MedFixable)},
		[]string{"Low", fmt.Sprint(counts.Low),
			fmt.Sprint(counts.LowFixable)},
		[]string{"Negligible", fmt.Sprint(counts.Negligible),
			fmt.Sprint(counts.NegFixable)},
	}
}

func buildHostVulnCVEsToTableError() string {
	msg := "There are no"
	if vulCmdState.Fixable {
		msg = fmt.Sprintf("%s fixable", msg)
	}
	msg = fmt.Sprintf("%s vulnerabilities", msg)

	if vulCmdState.Active {
		msg = fmt.Sprintf("%s of packages actively running", msg)
	}
	return fmt.Sprintf("%s in your environment.\n", msg)
}

func addToHostSummary(text []string, num int32, severity string) []string {
	if len(text) == 0 {
		if num != 0 {
			return append(text, fmt.Sprintf("%d %s", num, severity))
		}
	}
	return text
}

func hostScanPackagesVulnToTable(scan *api.HostVulnScanPkgManifestResponse) string {
	var (
		tableBuilder   = &strings.Builder{}
		summaryBuilder = &strings.Builder{}
		t              *tablewriter.Table
		rows           [][]string
		headers        []string
	)

	if vulCmdState.Packages {
		rows = hostScanPackagesVulnPackagesView(scan.Vulns)
		headers = []string{
			"CVE Count",
			"Severity",
			"Package",
			"Version",
			"Fixes Available",
		}
	} else {
		rows = hostScanPackagesVulnDetailsTable(scan.Vulns)
		headers = []string{
			"CVE",
			"Severity",
			"Score",
			"Package",
			"Version",
			"Fix Version",
		}
	}

	if len(rows) == 0 {
		if vulCmdState.Fixable {
			return "There are no fixable vulnerabilities.\n"
		}
		scannedVia := "package manifest"
		if pkgManifestLocal {
			scannedVia = "localhost"
		}
		return fmt.Sprintf("Great news! The %s has no vulnerabilities... (time for %s)\n",
			scannedVia, randomEmoji())
	}

	t = tablewriter.NewWriter(summaryBuilder)
	t.SetBorder(false)
	t.SetColumnSeparator(" ")
	t.SetHeader([]string{
		"Severity", "Count", "Fixable",
	})
	t.AppendBulk(hostVulnAssessmentToCountsTable(scan.VulnerabilityCounts()))
	t.Render()

	t = tablewriter.NewWriter(tableBuilder)
	t.SetBorder(false)
	t.SetAutoWrapText(false)
	t.SetHeader([]string{"Vulnerabilities"})
	t.Append([]string{summaryBuilder.String()})
	t.Render()

	t = tablewriter.NewWriter(tableBuilder)
	t.SetHeader(headers)
	t.SetBorder(false)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.AppendBulk(rows)
	t.Render()

	return tableBuilder.String()
}

func hostScanPackagesVulnDetailsTable(vulns []api.HostScanPackageVulnDetails) [][]string {
	out := [][]string{}
	for _, vuln := range vulns {
		if vuln.Summary.EvalStatus != "MATCH_VULN" {
			continue
		}

		if vulCmdState.Fixable && vuln.FixInfo.EvalStatus != "GOOD" {
			continue
		}

		fixedVersion := ""
		if vuln.FixInfo.EvalStatus == "GOOD" {
			fixedVersion = vuln.FixInfo.FixedVersion
		}

		out = append(out, []string{
			vuln.VulnID,
			vuln.Severity,
			vuln.ScoreString(),
			vuln.OsPkgInfo.Pkg,
			vuln.OsPkgInfo.PkgVer,
			fixedVersion,
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
}

func hostScanPackagesVulnPackagesView(vulns []api.HostScanPackageVulnDetails) [][]string {
	out := [][]string{}
	for _, vuln := range vulns {
		if vuln.Summary.EvalStatus != "MATCH_VULN" {
			continue
		}

		if vulCmdState.Fixable && vuln.FixInfo.EvalStatus != "GOOD" {
			continue
		}

		added := false
		for i := range out {
			if out[i][1] == vuln.Severity &&
				out[i][2] == vuln.OsPkgInfo.Pkg &&
				out[i][3] == vuln.OsPkgInfo.PkgVer {

				if countCVEs, err := strconv.Atoi(out[i][0]); err == nil {
					out[i][0] = fmt.Sprintf("%d", (countCVEs + 1))
					added = true
				}

				if vuln.FixInfo.EvalStatus == "GOOD" {
					if fixes, err := strconv.Atoi(out[i][4]); err == nil {
						out[i][4] = fmt.Sprintf("%d", (fixes + 1))
					}
				}
			}
		}

		if added {
			continue
		}

		fixes := "0"
		if vuln.FixInfo.EvalStatus == "GOOD" {
			fixes = "1"
		}

		out = append(out, []string{
			"1",
			vuln.Severity,
			vuln.OsPkgInfo.Pkg,
			vuln.OsPkgInfo.PkgVer,
			fixes,
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
}
