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
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
		Short: "Generates a package-manifest from the local host",
		Long: `Generates a package-manifest formatted for usage with the Lacework
scan package-manifest API.

Additionally, you can automatically generate a package-manifest from
the local host and send it directly to the Lacework API with the command:

    lacework vulnerability host scan-pkg-manifest --local`,
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
		Short: "Request an on-demand host vulnerability assessment from a package-manifest",
		Long: `Request an on-demand host vulnerability assessment of your software packages to
determine if the packages contain any common vulnerabilities and exposures.

Simple usage:

    lacework vulnerability host scan-pkg-manifest '{
        "osPkgInfoList": [
            {
                "os":"Ubuntu",
                "osVer":"18.04",
                "pkg": "openssl",
                "pkgVer": "1.1.1-1ubuntu2.1~18.04.5"
            }
        ]
    }'

To generate a package-manifest from the local host and scan it automatically:

    lacework vulnerability host scan-pkg-manifest --local

**NOTE:**
 - Only packages managed by a package manager for supported OS's are reported.
 - Calls to this operation are rate limited to 10 calls per hour, per access key.
 - This operation is limited to 10k packages per command execution.`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}

			var (
				pkgManifest      = new(api.VulnerabilitiesPackageManifest)
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
			var response api.VulnerabilitySoftwarePackagesResponse
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
				response, err = cli.LwApi.V2.Vulnerabilities.SoftwarePackages.Scan(*pkgManifest)
			}
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to request an on-demand host vulnerability scan")
			}

			if err := buildVulnHostScanPkgManifestReports(response); err != nil {
				return err
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

	vulHostShowAssessmentCmd = &cobra.Command{
		Use:     "show-assessment <machine_id>",
		Aliases: []string{"show"},
		Args:    cobra.ExactArgs(1),
		Short:   "Show results of a host vulnerability assessment",
		Long: `Show results of a host vulnerability assessment.

To find the machine id from hosts in your environment, use the command:

    lacework vulnerability host list-cves

Grab a CVE id and feed it to the command:

    lacework vulnerability host list-hosts my_cve_id`,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if vulCmdState.Csv {
				cli.EnableCSVOutput()

				// when rendering csv output, default to details since there is no output with less verbosity
				if !vulCmdState.Details && !vulCmdState.Packages {
					vulCmdState.Details = true
				}
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}

			var (
				assessment api.VulnerabilitiesHostResponse
				cacheKey   = fmt.Sprintf("host/assessment/v2/%s", args[0])
			)

			expired := cli.ReadCachedAsset(cacheKey, &assessment)
			if expired {
				// check machine exists
				var machineDetailsResponse api.MachineDetailsEntityResponse
				filter := api.SearchFilter{Filters: []api.Filter{{
					Expression: "eq",
					Field:      "mid",
					Value:      args[0],
				}}}

				cli.StartProgress(fmt.Sprintf("Searching for machine with id %s...", args[0]))
				err := cli.LwApi.V2.Entities.Search(&machineDetailsResponse, filter)
				cli.StopProgress()

				if err != nil {
					return errors.Wrapf(err, "unable to get machine details id %s", args[0])
				}

				if len(machineDetailsResponse.Data) == 0 {
					cli.OutputHuman("no hosts found with id %s\n", args[0])
					return nil
				}

				cli.StartProgress("Fetching host vulnerabilities...")
				response, err := cli.LwApi.V2.Vulnerabilities.Hosts.Search(filter)
				if err != nil {
					return errors.Wrapf(err, "unable to get host assessment with id %s", args[0])
				}
				cli.StopProgress()

				assessment = response
				assessment.Data = filterUniqueCves(response.Data)

				cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Hour*1), assessment)
			}

			if err := buildVulnHostReports(assessment); err != nil {
				return err
			}

			if vulFailureFlagsEnabled() {
				cli.Log.Infow("failure flags enabled",
					"fail_on_severity", vulCmdState.FailOnSeverity,
					"fail_on_fixable", vulCmdState.FailOnFixable,
				)
				assessmentCounts := assessment.VulnerabilityCounts()
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
		Short: "List host vulnerability assessments from a time range",
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

	setCsvFlag(
		vulHostListCvesCmd.Flags(),
		vulHostListHostsCmd.Flags(),
		vulHostShowAssessmentCmd.Flags(),
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

func filterUniqueCves(hosts []api.VulnerabilityHost) []api.VulnerabilityHost {
	var uniqueCves []api.VulnerabilityHost
	var cves []string

	for _, host := range hosts {
		if host.VulnID == "" || array.ContainsStr(cves, host.VulnID) {
			continue
		}
		uniqueCves = append(uniqueCves, host)
		cves = append(cves, host.VulnID)

	}
	return uniqueCves
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

func cvesSummary(hosts []api.VulnerabilityHost) map[string]api.VulnCveSummary {
	uniqueCves := make(map[string]api.VulnCveSummary)
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
		sum := api.VulnCveSummary{Host: host, Count: 1}
		sum.Hostnames = append(sum.Hostnames, sum.Host.EvalCtx.Hostname)
		uniqueCves[host.VulnID] = sum
	}
	return uniqueCves
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
	for k, _ := range sevSummaries {
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

func aggregatePackagesWithHosts(slice []packageTable, s packageTable, host bool, hasFix bool) []packageTable {
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
			if hasFix {
				slice[i].fixes++
			}
			return slice
		}
	}
	return append(slice, s)
}

func sortVulnHosts(slice []api.VulnerabilityHost) {
	for range slice {
		sort.Slice(slice[:], func(i, j int) bool {
			switch strings.Compare(slice[i].VulnID, slice[j].VulnID) {
			case -1:
				return true
			case 1:
				return false
			default:
				return false
			}
		})
	}
}

func hostVulnPackagesTable(cves map[string]api.VulnCveSummary, withHosts bool) ([][]string, string) {
	var (
		out                [][]string
		filteredPackages   []packageTable
		aggregatedPackages []packageTable
		cveSlice           []api.VulnerabilityHost
	)

	for _, cve := range cves {
		cveSlice = append(cveSlice, cve.Host)
	}
	sortVulnHosts(cveSlice)

	for _, host := range cveSlice {
		pack := packageTable{
			cveCount:       1,
			severity:       cases.Title(language.English).String(host.Severity),
			packageName:    host.FeatureKey.Namespace,
			currentVersion: host.FeatureKey.VersionInstalled,
			fixVersion:     host.FixInfo.FixedVersion,
			packageStatus:  host.PackageActive(),
		}
		if withHosts {
			pack.hostCount = 1
		}

		if vulCmdState.Active && host.PackageActive() == "" {
			filteredPackages = aggregatePackagesWithHosts(filteredPackages, pack, withHosts, false)
			continue
		}

		if vulCmdState.Fixable && host.FixInfo.FixedVersion == "" {
			filteredPackages = aggregatePackagesWithHosts(filteredPackages, pack, withHosts, false)
			continue
		}

		if vulCmdState.Severity != "" {
			if filterSeverity(host.Severity, vulCmdState.Severity) {
				filteredPackages = aggregatePackagesWithHosts(filteredPackages, pack, withHosts, false)
				continue
			}
		}
		aggregatedPackages = aggregatePackagesWithHosts(aggregatedPackages, pack, withHosts, false)
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

	if len(filteredPackages) > 0 {
		filteredOutput := fmt.Sprintf("%d of %d package(s) showing\n", len(out), len(aggregatedPackages)+len(filteredPackages))
		return out, filteredOutput
	}

	return out, ""
}

func filterHostCVEsTable(cves map[string]api.VulnCveSummary) (map[string]api.VulnCveSummary, string) {
	var out map[string]api.VulnCveSummary
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

func hostVulnCVEsTable(hostSummary map[string]api.VulnCveSummary) [][]string {
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

func filterHostVulnCVEs(cves map[string]api.VulnCveSummary) (map[string]api.VulnCveSummary, int, int) {
	var (
		filtered = 0
		total    = 0
		out      = make(map[string]api.VulnCveSummary)
	)

	for _, c := range cves {
		cve := c.Host
		total++
		// if the user wants to show only vulnerabilities of active packages
		if vulCmdState.Active && cve.Status != "Active" {
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

func hostVulnHostDetailsMainReportTable(assessment api.VulnerabilitiesHostResponse) string {
	host := assessment.Data[0]
	mainBldr := &strings.Builder{}
	mainBldr.WriteString(
		renderCustomTable([]string{"Host Details", "Vulnerabilities"},
			[][]string{[]string{
				renderCustomTable([]string{},
					[][]string{
						{"Machine ID", strconv.Itoa(host.Mid)},
						{"Hostname", host.EvalCtx.Hostname},
						{"External IP", host.MachineTags.ExternalIP},
						{"Internal IP", host.MachineTags.InternalIP},
						{"Os", host.MachineTags.Os},
						{"Arch", host.MachineTags.Arch},
						{"Namespace", host.FeatureKey.Namespace},
						{"Provider", host.MachineTags.VMProvider},
						{"Instance ID", host.MachineTags.InstanceID},
						{"AMI", host.MachineTags.AmiID},
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

	return mainBldr.String()
}

func showPackages() bool {
	return vulCmdState.Details || vulCmdState.Fixable || vulCmdState.Packages || vulCmdState.Active || vulCmdState.Severity != ""
}

func buildVulnHostsDetailsTableCSV(filteredCves map[string]api.VulnCveSummary) ([]string, [][]string) {
	if !showPackages() {
		return nil, nil
	}

	if vulCmdState.Packages {
		packages, _ := hostVulnPackagesTable(filteredCves, false)
		sort.Slice(packages, func(i, j int) bool {
			return severityOrder(packages[i][1]) < severityOrder(packages[j][1])
		})
		// order by cve count
		return []string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version", "Pkg Status"}, packages
	}

	rows := hostVulnCVEsTableForHostViewCSV(filteredCves)
	sort.Slice(rows, func(i, j int) bool {
		return severityOrder(rows[i][1]) < severityOrder(rows[j][1])
	})
	return []string{"CVE ID", "Severity", "Score", "Package", "Package Namespace", "Current Version",
		"Fix Version", "Pkg Status", "First Seen", "Last Status Update", "Vuln Status"}, rows
}

func buildVulnHostsDetailsTable(filteredCves map[string]api.VulnCveSummary) string {
	mainBldr := &strings.Builder{}

	if showPackages() {
		if vulCmdState.Packages {
			packages, filtered := hostVulnPackagesTable(filteredCves, false)
			sort.Slice(packages, func(i, j int) bool {
				return severityOrder(packages[i][1]) < severityOrder(packages[j][1])
			})
			// if the user wants to show only vulnerabilities of active packages
			// and we don't have any, show a friendly message
			if len(packages) == 0 {
				mainBldr.WriteString(buildHostVulnCVEsToTableError())
			} else {
				mainBldr.WriteString(
					renderSimpleTable(
						[]string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version", "Pkg Status"},
						packages,
					),
				)
				if filtered != "" {
					mainBldr.WriteString(filtered)
				}
			}
		} else {
			rows := hostVulnCVEsTableForHostView(filteredCves)
			sort.Slice(rows, func(i, j int) bool {
				return severityOrder(rows[i][1]) < severityOrder(rows[j][1])
			})
			// if the user wants to show only vulnerabilities of active packages
			// and we don't have any, show a friendly message
			if len(rows) == 0 {
				mainBldr.WriteString(buildHostVulnCVEsToTableError())
			} else {
				mainBldr.WriteString(renderSimpleTable([]string{
					"CVE ID", "Severity", "CvssV2", "CvssV3", "Package", "Current Version",
					"Fix Version", "Pkg Status", "Vuln Status"},
					rows,
				))
			}
		}
	}

	if !vulCmdState.Details && !vulCmdState.Active && !vulCmdState.Fixable && !vulCmdState.Packages && vulCmdState.Severity == "" {
		mainBldr.WriteString(
			"\nTry adding '--details' to increase details shown about the vulnerability assessment.\n",
		)
	} else if !vulCmdState.Active {
		mainBldr.WriteString(
			"\nTry adding '--active' to only show vulnerabilities of packages actively running.\n",
		)
	} else if !vulCmdState.Fixable {
		mainBldr.WriteString(
			"\nTry adding '--fixable' to only show fixable vulnerabilities.\n",
		)
	}

	return mainBldr.String()
}

func hostVulnCVEsTableForHostViewCSV(cves map[string]api.VulnCveSummary) [][]string {
	var (
		out      [][]string
		cveSlice []api.VulnerabilityHost
	)

	for _, cve := range cves {
		cveSlice = append(cveSlice, cve.Host)
	}
	sortVulnHosts(cveSlice)

	for _, host := range cveSlice {
		var (
			firstSeen   string
			lastUpdated string
		)

		if host.Props.FirstTimeSeen != nil {
			firstSeen = host.Props.FirstTimeSeen.UTC().String()
		}

		if host.Props.FirstTimeSeen != nil {
			lastUpdated = host.Props.LastUpdatedTime.UTC().String()
		}

		out = append(out, []string{
			host.VulnID,
			host.Severity,
			host.CvssV2(),
			host.CvssV3(),
			host.FeatureKey.Name,
			host.FeatureKey.Namespace,
			host.FeatureKey.VersionInstalled,
			host.FixInfo.FixedVersion,
			host.PackageActive(),
			firstSeen,
			lastUpdated,
			host.PackageActive(),
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
}

func hostVulnCVEsTableForHostView(summary map[string]api.VulnCveSummary) [][]string {
	var (
		out      [][]string
		cveSlice []api.VulnerabilityHost
	)

	for _, cve := range summary {
		cveSlice = append(cveSlice, cve.Host)
	}
	sortVulnHosts(cveSlice)

	for _, host := range cveSlice {
		out = append(out, []string{
			host.VulnID,
			host.Severity,
			host.CvssV2(),
			host.CvssV3(),
			host.FeatureKey.Name,
			host.FeatureKey.VersionInstalled,
			host.FixInfo.FixedVersion,
			host.PackageActive(),
			host.Status,
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
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

func hostScanPackagesVulnToTable(scan api.VulnerabilitySoftwarePackagesResponse) string {
	var (
		mainBldr = &strings.Builder{}
		rows     [][]string
		headers  []string
	)

	if vulCmdState.Packages {
		rows = hostScanPackagesVulnPackagesTable(filterHostScanPackagesVulnPackages(scan.Data))
		headers = []string{
			"CVE Count",
			"Severity",
			"Package",
			"Version",
			"Fixes Available",
		}
	} else {
		rows = hostScanPackagesVulnDetailsTable(scan.Data)
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

func filterHostScanPackagesVulnDetails(vulns []api.VulnerabilitySoftwarePackage) []api.VulnerabilitySoftwarePackage {
	out := make([]api.VulnerabilitySoftwarePackage, 0)

	for _, vuln := range vulns {
		if vulCmdState.Fixable && !vuln.HasFix() {
			continue
		}

		out = append(out, vuln)
	}

	return out
}

func hostScanPackagesVulnDetailsTable(vulns []api.VulnerabilitySoftwarePackage) [][]string {
	var out [][]string
	for _, vuln := range vulns {

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

func filterHostScanPackagesVulnPackages(vulns []api.VulnerabilitySoftwarePackage) filteredPackageTable {
	var (
		filteredPackages   []packageTable
		aggregatedPackages []packageTable
	)

	for _, vuln := range vulns {
		pack := packageTable{
			cveCount:       1,
			severity:       cases.Title(language.English).String(vuln.Severity),
			packageName:    vuln.OsPkgInfo.Pkg,
			currentVersion: vuln.OsPkgInfo.PkgVer,
		}

		if vulCmdState.Fixable && !vuln.HasFix() {
			filteredPackages = aggregatePackagesWithHosts(aggregatedPackages, pack, false, false)
			continue
		}

		aggregatedPackages = aggregatePackagesWithHosts(aggregatedPackages, pack, false, vuln.HasFix())
	}

	return filteredPackageTable{
		packages:        aggregatedPackages,
		totalPackages:   len(aggregatedPackages),
		totalUnfiltered: len(filteredPackages) + len(aggregatedPackages),
	}
}

func hostScanPackagesVulnPackagesTable(pkgs filteredPackageTable) [][]string {
	var out [][]string
	for _, pkg := range pkgs.packages {
		out = append(out, []string{
			"1",
			pkg.severity,
			pkg.packageName,
			pkg.currentVersion,
			strconv.Itoa(pkg.fixes),
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return severityOrder(out[i][1]) < severityOrder(out[j][1])
	})

	return out
}

func removeFixedVulnerabilitiesFromAssessment(assessment []api.VulnerabilityHost) []api.VulnerabilityHost {
	var filteredCves []api.VulnerabilityHost
	for _, cve := range assessment {
		if cve.Status != "Fixed" {
			filteredCves = append(filteredCves, cve)
		}
	}
	return filteredCves
}

// Build the cli output for vuln host show-assessment
func buildVulnHostReports(response api.VulnerabilitiesHostResponse) error {
	// @afiune the UI today doesn't display any vulnerability that has been fixed
	// but the APIs return them, this is causing confusion, to fix this issue we
	// are filtering all of those "Fixed" vulnerabilities here
	response.Data = removeFixedVulnerabilitiesFromAssessment(response.Data)

	hostVulnCounts := response.VulnerabilityCounts()
	if hostVulnCounts.Total == 0 {
		if cli.JSONOutput() {
			return cli.OutputJSON(response.Data)
		}
		cli.OutputHuman("Great news! This host has no vulnerabilities... (time for %s)\n", randomEmoji())
		return nil
	}

	var (
		mainReport                  = hostVulnHostDetailsMainReportTable(response)
		filteredCves, filtered      = filterHostCVEsTable(cvesSummary(response.Data))
		detailsReport               = buildVulnHostsDetailsTable(filteredCves)
		csvHeader, csvDetailsReport = buildVulnHostsDetailsTableCSV(filteredCves)
	)

	switch {
	case cli.JSONOutput():
		return cli.OutputJSON(response.Data)
	case cli.CSVOutput():
		return cli.OutputCSV(csvHeader, csvDetailsReport)
	default:
		cli.OutputHuman(mainReport)
		cli.OutputHuman(detailsReport)
		if filtered != "" {
			cli.OutputHuman(filtered)
		}
		return nil
	}
}

// Build the cli output for vuln host list-cves
func buildListCVEReports(cves []api.VulnerabilityHost) error {
	uniqueCves := cvesSummary(cves)
	uniqueCves, filtered := filterHostCVEsTable(uniqueCves)

	if cli.JSONOutput() {
		if uniqueCves == nil {
			if err := cli.OutputJSON(buildHostVulnCVEsToTableError()); err != nil {
				return err
			}
		} else {
			if err := cli.OutputJSON(uniqueCves); err != nil {
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

	if vulCmdState.Packages {
		packages, filtered := hostVulnPackagesTable(uniqueCves, true)

		if cli.CSVOutput() {
			return cli.OutputCSV(
				[]string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version", "Pkg Status", "Hosts"},
				packages,
			)
		}

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

	rows := hostVulnCVEsTable(uniqueCves)
	// if the user wants to show only online or
	// offline hosts, show a friendly message
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

// Build the cli output for vuln host scan-package-manifest
func buildVulnHostScanPkgManifestReports(response api.VulnerabilitySoftwarePackagesResponse) error {
	response.Data = filterHostScanPackagesVulnDetails(response.Data)

	if cli.JSONOutput() {
		return cli.OutputJSON(response)
	}

	if len(response.Data) == 0 {
		cli.OutputHuman(fmt.Sprintf("There are no vulnerabilities found! Time for %s\n", randomEmoji()))
	} else {
		cli.OutputHuman(hostScanPackagesVulnToTable(response))
	}

	return nil
}

type vulnSummary struct {
	host     api.VulnerabilityHost
	severity []string
	fixable  int
	count    int
}
