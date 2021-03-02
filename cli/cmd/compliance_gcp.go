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
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
)

var (
	// complianceGcpListProjCmd represents the list-projects sub-command inside the gcp command
	complianceGcpListProjCmd = &cobra.Command{
		Use:     "list-projects <organization_id>",
		Aliases: []string{"list-proj", "list"},
		Short:   "list projects from an organization",
		Long: `List all GCP projects from the provided organization ID.

Use the following command to list all GCP integrations in your account:

    $ lacework integrations list --type GCP_CFG

Then, select one GUID from an integration and visualize its details using the command:

    $ lacework integration show <int_guid>
`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Compliance.ListGcpProjects(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to list gcp projects")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Data[0])
			}

			rows := [][]string{}
			for _, gcp := range response.Data {
				for _, proj := range gcp.Projects {
					rows = append(rows, []string{proj})
				}
			}
			cli.OutputHuman(renderSimpleTable([]string{"Projects"}, rows))
			return nil
		},
	}

	// complianceGcpGetReportCmd represents the get-report sub-command inside the gcp command
	complianceGcpGetReportCmd = &cobra.Command{
		Use:     "get-report <organization_id> <project_id>",
		Aliases: []string{"get"},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			switch compCmdState.Type {
			case "CIS", "SOC", "PCI":
				compCmdState.Type = fmt.Sprintf("GCP_%s", compCmdState.Type)
				return nil
			case "GCP_CIS", "GCP_SOC", "GCP_PCI":
				return nil
			default:
				return errors.New("supported report types are: CIS, SOC, or PCI")
			}
		},
		Short: "get the latest GCP compliance report",
		Long: `Get the latest compliance assessment report, these reports run on a regular schedule,
typically once a day. The available report formats are human-readable (default), json and pdf.

To run an ad-hoc compliance assessment use the command:

    $ lacework compliance gcp run-assessment <project_id>
`,
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := api.ComplianceGcpReportConfig{
				OrganizationID: args[0],
				ProjectID:      args[1],
				Type:           compCmdState.Type,
			}

			if compCmdState.Pdf || compCmdState.PdfName != "" {
				pdfName := fmt.Sprintf(
					"%s_Report_%s_%s_%s_%s.pdf",
					config.Type,
					config.OrganizationID,
					config.ProjectID,
					cli.Account, time.Now().Format("20060102150405"),
				)
				if compCmdState.PdfName != "" {
					cli.OutputHuman("(DEPRECATED) This flag has been replaced by '--pdf'\n\n")
					pdfName = compCmdState.PdfName
				}

				cli.StartProgress(" Downloading compliance report...")
				err := cli.LwApi.Compliance.DownloadGcpReportPDF(pdfName, config)
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to get gcp pdf compliance report")
				}

				cli.OutputHuman("The GCP compliance report was downloaded at '%s'\n", pdfName)
				return nil
			}

			if compCmdState.Severity != "" {
				if !array.ContainsStr(api.ValidEventSeverities, compCmdState.Severity) {
					return errors.Errorf("the severity %s is not valid, use one of %s",
						compCmdState.Severity, strings.Join(api.ValidEventSeverities, ", "),
					)
				}
			}
			if compCmdState.Status != "" {
				if !array.ContainsStr(api.ValidComplianceStatuses, compCmdState.Status) {
					return errors.Errorf("the status %s is not valid, use one of %s",
						compCmdState.Status, strings.Join(api.ValidComplianceStatuses, ", "),
					)
				}
			}

			cli.StartProgress(" Getting compliance report...")
			response, err := cli.LwApi.Compliance.GetGcpReport(config)
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get gcp compliance report")
			}

			if len(response.Data) == 0 {
				return errors.New("there is no data found in the report")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Data[0])
			}

			report := response.Data[0]
			cli.OutputHuman("\n")
			recommendations, filteredOutput := complianceReportRecommendationsTable(report.Recommendations)
			cli.OutputHuman(
				buildComplianceReportTable(
					complianceGcpReportDetailsTable(&report),
					complianceReportSummaryTable(report.Summary),
					recommendations,
					filteredOutput,
				),
			)
			return nil
		},
	}

	// complianceGcpRunAssessmentCmd represents the run-assessment sub-command inside the gcp command
	complianceGcpRunAssessmentCmd = &cobra.Command{
		Use:     "run-assessment <org_or_project_id>",
		Aliases: []string{"run"},
		Short:   "run a new GCP compliance assessment",
		Long:    `Run a compliance assessment for the provided GCP organization or project.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := cli.LwApi.Compliance.RunGcpReport(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to run gcp compliance assessment")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			cli.OutputHuman("A new GCP compliance assessment has been initiated.\n")
			cli.OutputHuman("\n")
			cli.OutputHuman(
				renderSimpleTable(
					[]string{"INTEGRATION GUID", "ORG/PROJECT ID"},
					[][]string{[]string{response.IntgGuid, args[0]}},
				),
			)
			return nil
		},
	}
)

func init() {
	// add sub-commands to the gcp command
	complianceGcpCmd.AddCommand(complianceGcpListProjCmd)
	complianceGcpCmd.AddCommand(complianceGcpRunAssessmentCmd)
	complianceGcpCmd.AddCommand(complianceGcpGetReportCmd)

	complianceGcpGetReportCmd.Flags().BoolVar(&compCmdState.Details, "details", false,
		"increase details about the compliance report",
	)
	complianceGcpGetReportCmd.Flags().StringVar(&compCmdState.PdfName, "pdf-file", "",
		"(DEPRECATED) use --pdf",
	)
	complianceGcpGetReportCmd.Flags().BoolVar(&compCmdState.Pdf, "pdf", false,
		"download report in PDF format",
	)

	// GCP report types: GCP_CIS, GCP_SOC, or GCP_PCI.
	complianceGcpGetReportCmd.Flags().StringVar(&compCmdState.Type, "type", "CIS",
		"report type to display, supported types: CIS, SOC, or PCI",
	)

	complianceGcpGetReportCmd.Flags().StringVar(&compCmdState.Category, "category", "",
		"filter the compliance report details view by category (identity-and-access-management, s3, logging...)",
	)

	complianceGcpGetReportCmd.Flags().StringVar(&compCmdState.Service, "service", "",
		"filter the compliance report details view by service (aws:s3, aws:iam, aws:cloudtrail ...)",
	)

	complianceGcpGetReportCmd.Flags().StringVar(&compCmdState.Severity, "severity", "",
		fmt.Sprintf("filter compliance report details view by severity threshold (%s)",
			strings.Join(api.ValidEventSeverities, ", ")),
	)

	complianceGcpGetReportCmd.Flags().StringVar(&compCmdState.Status, "status", "",
		fmt.Sprintf("filter compliance report details view by status (%s)",
			strings.Join(api.ValidComplianceStatuses, ", ")),
	)
}

func complianceGcpReportDetailsTable(report *api.ComplianceGcpReport) [][]string {
	return [][]string{
		[]string{"Report Type", report.ReportType},
		[]string{"Report Title", report.ReportTitle},
		[]string{"Organization ID", report.OrganizationID},
		[]string{"Organization Name", report.OrganizationName},
		[]string{"Project ID", report.ProjectID},
		[]string{"Project Name", report.ProjectName},
		[]string{"Report Time", report.ReportTime.UTC().Format(time.RFC3339)},
	}
}
