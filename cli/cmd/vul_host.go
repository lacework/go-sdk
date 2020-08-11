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

import "github.com/spf13/cobra"

var (
	vulHostScanPkgManifestCmd = &cobra.Command{
		Use:     "scan-pkg-manifest",
		Aliases: []string{"scan", "pkg"},
		Short:   "request an on-demand host vulnerability assessment from a package-manifest",
		Long:    "Request an on-demand host vulnerability assessment from a package-manifest",
		RunE: func(_ *cobra.Command, args []string) error {
			return nil
		},
	}
	vulHostListAssessmentsCmd = &cobra.Command{
		Use:     "list-assessments",
		Aliases: []string{"list", "ls"},
		Short:   "list host vulnerability assessments from a time range",
		Long:    "List host vulnerability assessments from a time range.",
		RunE: func(_ *cobra.Command, args []string) error {
			return nil
		},
	}
	vulHostShowAssessmentCmd = &cobra.Command{
		Use:     "show-assessments",
		Aliases: []string{"show"},
		Short:   "show results of a host vulnerability assessment",
		Long:    "Sshow results of a host vulnerability assessment.",
		RunE: func(_ *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	// add sub-commands to the 'vulnerability container' command
	vulHostCmd.AddCommand(vulHostScanPkgManifestCmd)
	vulHostCmd.AddCommand(vulHostListAssessmentsCmd)
	vulHostCmd.AddCommand(vulHostShowAssessmentCmd)
}
