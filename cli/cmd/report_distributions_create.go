//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// create command is used to create a new lacework report distribution
var reportDistributionsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a report distribution",
	Long: `Create a new report distribution to view the evaluation of a set of policies in a report.

To create a new report distribution:

    lacework report-distribution create

To create a new report distribution from an existing file:

    lacework report-distribution create --file custom-report.json
`,
	Args: cobra.NoArgs,
	RunE: createReportDistribution,
}

func createReportDistribution(_ *cobra.Command, args []string) error {
	var (
		reportDistribution api.ReportDistribution
		err                error
	)
	reportDistribution, err = promptCreateReportDistributionFromNew()
	if err != nil {
		return err
	}

	cli.StartProgress("Creating report distribution...")
	resp, err := cli.LwApi.V2.ReportDistributions.Create(reportDistribution)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to create report distribution")
	}

	cli.OutputHuman("New report distribution created. To view the report run:\n\n"+
		"lacework report-distribution show %s \n", resp.Data.ReportDistributionGuid)
	return nil
}

func promptCreateReportDistributionFromNew() (reportDistribution api.ReportDistribution, err error) {
	cli.StartProgress("Fetching list of report definitions...")
	reportDefinitions, err := cli.LwApi.V2.ReportDefinitions.List()
	cli.StopProgress()
	if err != nil {
		return api.ReportDistribution{}, err
	}
	definitionMap := make(map[string]api.ReportDefinition, len(reportDefinitions.Data))
	var definitionDisplayOptions []string

	for _, definition := range reportDefinitions.Data {
		definitionMap[definition.DisplayName] = definition
		definitionDisplayOptions = append(definitionDisplayOptions, definition.DisplayName)
	}

	cli.StartProgress("Fetching list of alert channels...")
	alertChannels, err := cli.LwApi.V2.AlertChannels.List()
	cli.StopProgress()
	channelMap := make(map[string]string, len(alertChannels.Data))
	var channelOptions []string

	for _, definition := range alertChannels.Data {
		channelMap[definition.Name] = definition.IntgGuid
		channelOptions = append(channelOptions, definition.Name)
	}

	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: CreateReportDistributionReportNameQuestion},
			Validate: survey.Required,
		},
		{
			Name:     "frequency",
			Prompt:   &survey.Select{Message: CreateReportDistributionFrequencyQuestion, Options: api.ReportDistributionFrequencies()},
			Validate: survey.Required,
		},
		{
			Name:     "definition",
			Prompt:   &survey.Select{Message: CreateReportDistributionDefinitionQuestion, Options: definitionDisplayOptions},
			Validate: survey.Required,
		},
		{
			Name:     "channels",
			Prompt:   &survey.MultiSelect{Message: CreateReportDistributionDefinitionQuestion, Options: channelOptions},
			Validate: survey.Required,
		},
		{
			Name:     "groups",
			Prompt:   &survey.MultiSelect{Message: CreateReportDistributionDefinitionQuestion, Options: channelOptions},
			Validate: survey.Required,
		},
		{
			Name:     "integrations",
			Prompt:   &survey.MultiSelect{Message: CreateReportDistributionDefinitionQuestion, Options: channelOptions},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name           string   `survey:"name"`
		Frequency      string   `survey:"frequency"`
		Definition     string   `survey:"definition"`
		Channels       []string `survey:"channels"`
		ResourceGroups []string `survey:"groups"`
		Integrations   []string `survey:"integrations"`
	}{}

	if err = survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc)); err != nil {
		return
	}

	var channelAnswers []string
	for _, c := range answers.Channels {
		channelAnswers = append(channelAnswers, channelMap[c])
	}
	var groupAnswers []string
	for _, c := range answers.ResourceGroups {
		channelAnswers = append(channelAnswers, channelMap[c])
	}

	var integrations []api.ReportDistributionIntegration
	switch definitionMap[answers.Definition].SubReportType {
	case api.ReportDefinitionSubTypeAws.String():
		promptReportDistributionIntegrationsAws(integrations)
	case api.ReportDefinitionSubTypeGcp.String():
		// do gcp integration
		promptReportDistributionIntegrationsAws(integrations)
		integrations = append(integrations, api.ReportDistributionIntegrationGcp{
			OrganizationID: "",
			ProjectID:      "",
		})
	case api.ReportDefinitionSubTypeAzure.String():
		// do azure integration
		integrations = append(integrations, api.ReportDistributionIntegrationAzure{
			TenantID:       "",
			SubscriptionID: "",
		})
	default:
		return api.ReportDistribution{}, errors.Errorf("unsupported report definition type '%s'",
			definitionMap[answers.Definition].SubReportType)
	}

	reportDistribution = api.ReportDistribution{
		ReportDefinitionGuid: definitionMap[answers.Definition].ReportDefinitionGuid,
		DistributionName:     answers.Name,
		AlertChannels:        channelAnswers,
		Frequency:            answers.Frequency,
		Data: api.ReportDistributionData{
			ResourceGroups: groupAnswers,
			Integrations:   integrations,
		},
	}

	// prompt optional fields
	var violations []string
	if err = promptAddReportDistributionViolations(&violations); err != nil {
		return api.ReportDistribution{}, err
	}
	reportDistribution.Data.Violations = violations

	var severities []string
	if err = promptAddReportDistributionSeverities(&severities); err != nil {
		return api.ReportDistribution{}, err
	}
	reportDistribution.Data.Severities = severities

	return
}

func promptReportDistributionIntegrationsAws(integrations []api.ReportDistributionIntegration) error {
	custom := false
	if err := survey.AskOne(&survey.Confirm{
		Message: CreateReportDistributionIntegrationAwsQuestion,
	}, &custom); err != nil {
		return err
	}

	integrations = append(integrations, api.ReportDistributionIntegrationAws{AccountID: ""})
	return nil
}
