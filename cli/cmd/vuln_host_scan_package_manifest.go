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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// vulHostScanPkgManifestCmd represents the 'lacework vuln host scan-pkg-manifest' command
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

			if err := buildVulnHostScanPkgManifestReports(&response); err != nil {
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
)

// Build the cli output for vuln host scan-package-manifest
func buildVulnHostScanPkgManifestReports(response *api.VulnerabilitySoftwarePackagesResponse) error {
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

func hostScanPackagesVulnToTable(scan *api.VulnerabilitySoftwarePackagesResponse) string {
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
		if !vuln.IsVulnerable() {
			continue
		}

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
		out = append(out, []string{
			vuln.VulnID,
			vuln.Severity,
			vuln.ScoreString(),
			vuln.OsPkgInfo.Pkg,
			vuln.OsPkgInfo.PkgVer,
			vuln.FixInfo.FixedVersion,
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
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
		return api.SeverityOrder(out[i][1]) < api.SeverityOrder(out[j][1])
	})

	return out
}
