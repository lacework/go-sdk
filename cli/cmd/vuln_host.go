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
	// the maximum number of packages per scan request
	manifestPkgsCap = 1000

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
 - This operation is limited to 10k packages per command execution.`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}

			var (
				pkgManifest      = new(api.PackageManifest)
				pkgManifestBytes []byte
				err              error
			)

			if len(args) != 0 && args[0] != "" {
				pkgManifestBytes = []byte(args[0])
				cli.Log.Debugw("package manifest loaded from arguments", "raw", args[0])
			} else if pkgManifestFile != "" {
				pkgManifestBytes, err = ioutil.ReadFile(pkgManifestFile)
				if err != nil {
					return errors.Wrap(err, "unable to read file")
				}
				cli.Log.Debugw("package manifest loaded from file", "raw", string(pkgManifestBytes))
			} else if pkgManifestLocal {
				pkgManifest, err = cli.GeneratePackageManifest()
				if err != nil {
					return errors.Wrap(err, "unable to generate package manifest")
				}
				cli.Log.Debugw("package manifest generated from localhost", "raw", pkgManifest)
			} else {
				// avoid asking for a confirmation before launching the editor
				var content string
				prompt := &survey.Editor{
					Message:  "Provide a package manifest to scan",
					FileName: "package-manifest*.json",
				}
				err = survey.AskOne(prompt, &content)
				if err != nil {
					return errors.Wrap(err, "unable to load package manifest from editor")
				}
				pkgManifestBytes = []byte(content)
				cli.Log.Debugw("package manifest loaded via editor", "raw", content)
			}

			if len(pkgManifestBytes) != 0 {
				err = json.Unmarshal(pkgManifestBytes, pkgManifest)
				if err != nil {
					return errors.Wrap(err, "invalid package manifest json file")
				}
			}

			totalPkgs := len(pkgManifest.OsPkgInfoList)
			cli.StartProgress(" Scanning packages...")
			cli.Log.Infow("manifest", "total_packages", totalPkgs)
			var response api.HostVulnScanPkgManifestResponse
			// check if the package manifest has more than the maximum
			// number of packages, if so, make multiple API requests
			if totalPkgs >= manifestPkgsCap {
				cli.Log.Infow("manifest over the limit, splitting up")
				cli.Event.Feature = featSplitPkgManifest
				cli.Event.AddFeatureField("total_packages", totalPkgs)
				response, err = fanOutHostScans(
					splitPackageManifest(pkgManifest, manifestPkgsCap)...,
				)
			} else {
				response, err = cli.LwApi.Vulnerabilities.Host.Scan(pkgManifest)
			}
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to request an on-demand host vulnerability scan")
			}

			if cli.JSONOutput() {
				if err := cli.OutputJSON(response); err != nil {
					return err
				}
			} else {
				if len(response.Vulns) == 0 {
					// @afiune add a helpful message, possible things are:
					cli.OutputHuman("There are no vulnerabilities found.\n")
					return nil
				}
				cli.OutputHuman(hostScanPackagesVulnToTable(&response))
			}

			if vulFailureFlagsEnabled() {
				cli.Log.Infow("failure flags enabled",
					"fail_on_severity", vulCmdState.FailOnSeverity,
					"fail_on_fixable", vulCmdState.FailOnFixable,
				)
				assessmentCounts := response.VulnerabilityCounts()
				vulnPolicy := NewVulnerabilityPolicyError(
					&assessmentCounts,
					vulCmdState.FailOnSeverity,
					vulCmdState.FailOnFixable,
				)
				if vulnPolicy.NonCompliant() {
					c.SilenceUsage = true
					return vulnPolicy
				}
			}
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
			if err := validateSeverityFlags(); err != nil {
				return err
			}

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
				packages, filtered := hostVulnPackagesTable(response.CVEs, true)
				cli.OutputHuman(
					renderSimpleTable(
						[]string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version", "Pkg Status", "Hosts"},
						packages,
					),
				)
				if filtered != "" {
					cli.OutputHuman(filtered)
				}
				return nil
			}

			rows, filtered := hostVulnCVEsTable(response.CVEs)
			// if the user wants to show only online or
			// offline hosts, show a friendly message
			if len(rows) == 0 {
				cli.OutputHuman(buildHostVulnCVEsToTableError())
				return nil
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"CVE ID", "Severity", "Score", "Package", "Current Version",
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

			rows := hostVulnHostsTable(response.Hosts)
			// if the user wants to show only online or
			// offline hosts, show a friendly message
			if len(rows) == 0 {
				if vulCmdState.Online {
					cli.OutputHuman("There are no online hosts.\n")
				}
				if vulCmdState.Offline {
					cli.OutputHuman("There are no offline hosts.\n")
				}
				return nil
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
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}
			response, err := cli.LwApi.Vulnerabilities.Host.GetHostAssessment(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get host assessment with id "+args[0])
			}

			if cli.JSONOutput() {
				if err = cli.OutputJSON(response.Assessment); err != nil {
					return err
				}
			} else {
				cli.OutputHuman(hostVulnHostDetailsToTable(response.Assessment))
			}

			if vulFailureFlagsEnabled() {
				cli.Log.Infow("failure flags enabled",
					"fail_on_severity", vulCmdState.FailOnSeverity,
					"fail_on_fixable", vulCmdState.FailOnFixable,
				)
				assessmentCounts := response.Assessment.VulnerabilityCounts()
				vulnPolicy := NewVulnerabilityPolicyError(
					&assessmentCounts,
					vulCmdState.FailOnSeverity,
					vulCmdState.FailOnFixable,
				)
				if vulnPolicy.NonCompliant() {
					c.SilenceUsage = true
					return vulnPolicy
				}
			}
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

func aggregatePackagesWithHosts(slice []packageTable, s packageTable, host bool) []packageTable {
	for i, item := range slice {
		if item.packageName == s.packageName &&
			item.currentVersion == s.currentVersion &&
			item.severity == s.severity &&
			item.fixVersion == s.fixVersion &&
			item.packageStatus == s.packageStatus {
			slice[i].cveCount++
			if host {
				slice[i].hostCount++
			}
			return slice
		}
	}
	return append(slice, s)
}

func hostVulnPackagesTable(cves []api.HostVulnCVE, withHosts bool) ([][]string, string) {
	var (
		out                [][]string
		filteredPackages   []packageTable
		aggregatedPackages []packageTable
	)

	for _, cve := range cves {
		for _, pkg := range cve.Packages {
			pack := packageTable{
				cveCount:       1,
				severity:       strings.Title(pkg.Severity),
				packageName:    pkg.Name,
				currentVersion: pkg.Version,
				fixVersion:     pkg.FixedVersion,
				packageStatus:  pkg.PackageStatus,
			}
			if withHosts {
				pack.hostCount = 1
			}

			if vulCmdState.Active && pkg.PackageStatus == "" {
				filteredPackages = aggregatePackagesWithHosts(filteredPackages, pack, withHosts)
				continue
			}

			if vulCmdState.Fixable && pkg.FixedVersion == "" {
				filteredPackages = aggregatePackagesWithHosts(filteredPackages, pack, withHosts)
				continue
			}

			if vulCmdState.Severity != "" {
				if filterSeverity(pkg.Severity, vulCmdState.Severity) {
					filteredPackages = aggregatePackagesWithHosts(filteredPackages, pack, withHosts)
					continue
				}
			}
			aggregatedPackages = aggregatePackagesWithHosts(aggregatedPackages, pack, withHosts)
		}
	}

	for _, p := range aggregatedPackages {
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

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	if len(filteredPackages) > 0 {
		filteredOutput := fmt.Sprintf("%d of %d package(s) showing \n", len(out), len(aggregatedPackages)+len(filteredPackages))
		return out, filteredOutput
	}

	return out, ""
}

func hostVulnCVEsTable(cves []api.HostVulnCVE) ([][]string, string) {
	var out [][]string
	var filteredCves = 0
	var totalCves = 0

	for _, severity := range api.ValidEventSeverities {
		_, sev := severityToProperTypes(severity)
		output, filtered, total := hostVulnCVEsTableForSeverity(cves, sev)
		filteredCves += filtered
		totalCves += total
		out = append(out, output...)
	}

	if filteredCves > 0 {
		showing := totalCves - filteredCves
		return out, fmt.Sprintf("\n%d of %d cve(s) showing \n", showing, totalCves)
	}

	return out, ""
}

func hostVulnCVEsTableForSeverity(cves []api.HostVulnCVE, severity string) ([][]string, int, int) {
	var (
		filtered = 0
		total    = 0
		out      = [][]string{}
	)

	for _, cve := range cves {
		for _, pkg := range cve.Packages {
			if pkg.Severity == severity {
				total++

				// if the user wants to show only vulnerabilities of active packages
				if vulCmdState.Active && pkg.PackageStatus == "" {
					filtered++
					continue
				}
				if vulCmdState.Fixable && pkg.FixedVersion == "" {
					filtered++
					continue
				}

				if vulCmdState.Severity != "" {
					if filterSeverity(severity, vulCmdState.Severity) {
						filtered++
						continue
					}
				}

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

	return out, filtered, total
}

func hostVulnHostDetailsToTable(assessment api.HostVulnHostAssessment) string {
	mainBldr := &strings.Builder{}
	mainBldr.WriteString(
		renderCustomTable([]string{"Host Details", "Vulnerabilities"},
			[][]string{[]string{
				renderCustomTable([]string{},
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
					tableFunc(func(t *tablewriter.Table) {
						t.SetColumnSeparator("")
						t.SetBorder(false)
						t.SetAlignment(tablewriter.ALIGN_LEFT)

					}),
				),
				renderCustomTable(
					[]string{"Severity", "Count", "Fixable"},
					hostVulnAssessmentToCountsTable(assessment.VulnerabilityCounts()),
					tableFunc(func(t *tablewriter.Table) {
						t.SetBorder(false)
						t.SetColumnSeparator(" ")
					}),
				),
			}},
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetAutoWrapText(false)
				t.SetColumnSeparator(" ")
			}),
		),
	)

	mainBldr.WriteString("\n")
	if vulCmdState.Details || vulCmdState.Fixable || vulCmdState.Packages || vulCmdState.Active || vulCmdState.Severity != "" {
		if vulCmdState.Packages {
			packages, filtered := hostVulnPackagesTable(assessment.CVEs, false)
			mainBldr.WriteString(
				renderSimpleTable(
					[]string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version", "Pkg Status"},
					packages,
				),
			)
			if filtered != "" {
				mainBldr.WriteString(filtered)
			}
		} else {
			rows, filtered := hostVulnCVEsTableForHostView(assessment.CVEs)

			// if the user wants to show only vulnerabilities of active packages
			// and we don't have any, show a friendly message
			if len(rows) == 0 {
				if vulCmdState.Active && vulCmdState.Fixable {
					mainBldr.WriteString("There are no fixable vulnerabilities with packages actively running in this host.\n")
				}
				if vulCmdState.Active {
					mainBldr.WriteString("There are no vulnerabilities with packages actively running in this host.\n")
				}
				if vulCmdState.Active {
					mainBldr.WriteString("There are no fixable vulnerabilities in this host.\n")
				}
			} else {
				mainBldr.WriteString(renderSimpleTable([]string{
					"CVE ID", "Severity", "Score", "Package", "Current Version",
					"Fix Version", "Pgk Status", "Vuln Status"},
					rows,
				))
			}

			if filtered != "" {
				mainBldr.WriteString(filtered)
			}
		}
		mainBldr.WriteString("\n")
	}

	if !vulCmdState.Details && !vulCmdState.Active && !vulCmdState.Fixable && !vulCmdState.Packages && vulCmdState.Severity == "" {
		mainBldr.WriteString(
			"Try adding '--details' to increase details shown about the vulnerability assessment.\n",
		)
	} else if !vulCmdState.Active {
		mainBldr.WriteString(
			"Try adding '--active' to only show vulnerabilities of packages actively running.\n",
		)
	} else if !vulCmdState.Fixable {
		mainBldr.WriteString(
			"Try adding '--fixable' to only show fixable vulnerabilities.\n",
		)
	}

	return mainBldr.String()
}

func hostVulnCVEsTableForHostView(cves []api.HostVulnCVE) ([][]string, string) {
	var (
		total = 0
		out   [][]string
	)

	for _, cve := range cves {
		for _, pkg := range cve.Packages {
			total++
			// if the user wants to show only vulnerabilities of active packages
			if vulCmdState.Active && pkg.PackageStatus == "" {
				continue
			}

			if vulCmdState.Fixable && pkg.FixedVersion == "" {
				continue
			}

			if vulCmdState.Severity != "" && filterSeverity(pkg.Severity, vulCmdState.Severity) {
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

	if len(out) < total {
		showing := total - len(out)
		return out, fmt.Sprintf("\n%d of %d cve(s) showing \n", showing, total)
	}

	return out, ""
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

	if vulCmdState.Severity != "" {
		msg = fmt.Sprintf("%s %s", msg, vulCmdState.Severity)
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
		mainBldr = &strings.Builder{}
		rows     [][]string
		headers  []string
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
			"CVE ID",
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
		return fmt.Sprintf(
			"Great news! The %s has no vulnerabilities... (time for %s)\n",
			scannedVia, randomEmoji(),
		)
	}

	mainBldr.WriteString(
		renderOneLineCustomTable("Vulnerabilities",
			renderCustomTable(
				[]string{"Severity", "Count", "Fixable"},
				hostVulnAssessmentToCountsTable(scan.VulnerabilityCounts()),
				tableFunc(func(t *tablewriter.Table) {
					t.SetBorder(false)
					t.SetColumnSeparator(" ")
				}),
			),
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetAutoWrapText(false)
			}),
		),
	)

	mainBldr.WriteString(renderSimpleTable(headers, rows))

	return mainBldr.String()
}

func hostScanPackagesVulnDetailsTable(vulns []api.HostScanPackageVulnDetails) [][]string {
	out := [][]string{}
	for _, vuln := range vulns {
		if vulCmdState.Fixable && vuln.HasFix() {
			continue
		}

		fixedVersion := ""
		if vuln.HasFix() {
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
		if vulCmdState.Fixable && vuln.HasFix() {
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

				if vuln.HasFix() {
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
		if vuln.HasFix() {
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
