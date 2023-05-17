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
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// vulHostListCvesCmd represents the 'lacework vuln host list-cves' command
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

				cli.StartProgress(fmt.Sprintf("Searching for machine with id '%s'...", args[0]))
				err := cli.LwApi.V2.Entities.Search(&machineDetailsResponse, filter)
				cli.StopProgress()

				if err != nil {
					return errors.Wrapf(err, "unable to get machine details id %s", args[0])
				}

				if len(machineDetailsResponse.Data) == 0 {
					return errors.Errorf("no hosts found with id %s\n", args[0])
				}

				machineDetails := machineDetailsResponse.Data[0]

				cli.StartProgress(
					fmt.Sprintf("Searching for latest host evaluation for machine %s (%d)...",
						machineDetails.Hostname, machineDetails.Mid,
					))
				evalGUID, err := searchLastestHostEvaluationGuid(args[0])
				cli.StopProgress()
				if err != nil {
					return errors.Wrapf(err, "unable to find information of host '%s'", args[0])
				}

				cli.Log.Infow("latest assessment found", "eval_guid", evalGUID)

				var (
					now    = time.Now().UTC()
					before = now.AddDate(0, 0, -7) // 7 days from ago
				)

				filter.TimeFilter = &api.TimeFilter{
					StartTime: &before,
					EndTime:   &now,
				}
				filter.Filters = append(filter.Filters, api.Filter{
					Expression: "eq",
					Field:      "evalGuid",
					Value:      "c147082bf2b571841a0a24c4d7efff92",
				})

				cli.StartProgress(
					fmt.Sprintf("Fetching vulnerabilities from host evaluation '%s'...", evalGUID),
				)
				assessment, err = cli.LwApi.V2.Vulnerabilities.Hosts.SearchAllPages(filter)
				if err != nil {
					return errors.Wrapf(err, "unable to get host assessment with id %s", args[0])
				}
				cli.StopProgress()

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
)

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
		return cli.OutputJSON(summaryToHostList(filteredCves))
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

func searchLastestHostEvaluationGuid(mid string) (string, error) {
	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -7) // 7 days from ago
		filter = api.SearchFilter{
			TimeFilter: &api.TimeFilter{
				StartTime: &before,
				EndTime:   &now,
			},
			Filters: []api.Filter{{
				Expression: "eq",
				Field:      "mid",
				Value:      mid,
			}},
			Returns: []string{"evalGuid", "startTime"},
		}
	)

	cli.Log.Infow("retrieve host evaluation information", "mid", mid)
	response, err := cli.LwApi.V2.Vulnerabilities.Hosts.SearchAllPages(filter)
	if err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", errors.New("no data found")
	}

	return getUniqueHostEvalGUID(response), nil
}

func getUniqueHostEvalGUID(host api.VulnerabilitiesHostResponse) string {
	var (
		guid      string
		startTime time.Time
	)
	for _, ctr := range host.Data {
		if ctr.EvalGUID != guid {
			if ctr.StartTime.After(startTime) {
				startTime = ctr.StartTime
				guid = ctr.EvalGUID
			}
		}
	}
	return guid
}

func buildVulnHostsDetailsTableCSV(filteredCves map[string]VulnCveSummary) ([]string, [][]string) {
	if !showPackages() {
		return nil, nil
	}

	if vulCmdState.Packages {
		packages, _ := hostVulnPackagesTable(filteredCves, false)
		sort.Slice(packages, func(i, j int) bool {
			return api.SeverityOrder(packages[i][1]) < api.SeverityOrder(packages[j][1])
		})
		// order by cve count
		return []string{"CVE Count", "Severity", "Package", "Current Version", "Fix Version", "Pkg Status"}, packages
	}

	rows := hostVulnCVEsTableForHostViewCSV(filteredCves)
	sort.Slice(rows, func(i, j int) bool {
		return api.SeverityOrder(rows[i][1]) < api.SeverityOrder(rows[j][1])
	})
	return []string{"CVE ID", "Severity", "Score", "Package", "Package Namespace", "Current Version",
		"Fix Version", "Pkg Status", "First Seen", "Last Status Update", "Vuln Status"}, rows
}

func buildVulnHostsDetailsTable(filteredCves map[string]VulnCveSummary) string {
	mainBldr := &strings.Builder{}

	if showPackages() {
		if vulCmdState.Packages {
			packages, filtered := hostVulnPackagesTable(filteredCves, false)
			sort.Slice(packages, func(i, j int) bool {
				return api.SeverityOrder(packages[i][1]) < api.SeverityOrder(packages[j][1])
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
				return api.SeverityOrder(rows[i][1]) < api.SeverityOrder(rows[j][1])
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

func hostVulnHostDetailsMainReportTable(assessment api.VulnerabilitiesHostResponse) string {
	host := assessment.Data[0]
	machineTags, err := host.GetMachineTags()
	if err != nil {
		cli.Log.Debug("failed to parse machine tags")
	}
	mainBldr := &strings.Builder{}
	mainBldr.WriteString(
		renderCustomTable([]string{"Host Details", "Vulnerabilities"},
			[][]string{[]string{
				renderCustomTable([]string{},
					[][]string{
						{"Machine ID", strconv.Itoa(host.Mid)},
						{"Hostname", host.EvalCtx.Hostname},
						{"External IP", machineTags.ExternalIP},
						{"Internal IP", machineTags.InternalIP},
						{"Os", machineTags.Os},
						{"Arch", machineTags.Arch},
						{"Namespace", host.FeatureKey.Namespace},
						{"Provider", machineTags.VMProvider},
						{"Instance ID", machineTags.InstanceID},
						{"AMI", machineTags.AmiID},
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

func hostVulnCVEsTableForHostViewCSV(cves map[string]VulnCveSummary) [][]string {
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
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
	})

	return out
}

func hostVulnCVEsTableForHostView(summary map[string]VulnCveSummary) [][]string {
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
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
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

func hostVulnPackagesTable(cves map[string]VulnCveSummary, withHosts bool) ([][]string, string) {
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
			packageName:    host.FeatureKey.Name,
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
		// add all packages that have not been filtered
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
