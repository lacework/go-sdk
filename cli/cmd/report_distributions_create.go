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
	"github.com/lacework/go-sdk/lwseverity"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// create command is used to create a new lacework report distribution
var reportDistributionsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a report distribution",
	Long: `Create associates a report definition with a distribution channel for the report. 
A report distribution can refine the scope of the report by filtering its content by incident severity,
resource groups, and integrations.

To create a new report distribution:

    lacework report-distribution create
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
			Prompt:   &survey.MultiSelect{Message: CreateReportDistributionAlertChannelsQuestion, Options: channelOptions},
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

	reportDistribution = api.ReportDistribution{
		ReportDefinitionGuid: definitionMap[answers.Definition].ReportDefinitionGuid,
		DistributionName:     answers.Name,
		AlertChannels:        channelAnswers,
		Frequency:            answers.Frequency,
	}

	// scope can only be either resource group or integration
	if err = promptReportDistributionScope(&reportDistribution,
		definitionMap[answers.Definition].SubReportType); err != nil {
		return api.ReportDistribution{}, err
	}

	// prompt optional fields
	var violations []string
	if err = promptAddReportDistributionViolations(violations); err != nil {
		return api.ReportDistribution{}, err
	}
	reportDistribution.Data.Violations = violations

	var severities []string
	if err = promptAddReportDistributionSeverities(severities); err != nil {
		return api.ReportDistribution{}, err
	}
	reportDistribution.Data.Severities = severities

	return
}

func promptReportDistributionScope(distribution *api.ReportDistribution, subReportType string) error {
	distributionScope := ""
	if err := survey.AskOne(&survey.Select{
		Message: CreateReportDistributionScopeQuestion,
		Options: api.ReportDistributionScopes(),
	}, &distributionScope); err != nil {
		return err
	}

	if distributionScope == api.ReportDistributionScopeCloudIntegration.String() {
		if err := promptReportDistributionIntegration(distribution, subReportType); err != nil {
			return err
		}
		return nil
	} else {
		if err := promptReportDistributionResourceGroup(distribution); err != nil {
			return err
		}
		return nil
	}
}

func promptReportDistributionResourceGroup(distribution *api.ReportDistribution) error {
	cli.StartProgress("Fetching list of resource groups...")
	resourceGroups, err := cli.LwApi.V2.ResourceGroups.List()
	cli.StopProgress()
	groupMap := make(map[string]string, len(resourceGroups.Data))

	var (
		groupOptions []string
		groupAnswers []string
	)

	for _, definition := range resourceGroups.Data {
		groupMap[definition.Name] = definition.ID()
		groupOptions = append(groupOptions, definition.Name)
	}

	var selectedGroups []string
	if err = survey.AskOne(&survey.MultiSelect{
		Message: CreateReportDistributionResourceGroupsQuestion,
		Options: groupOptions,
	}, &selectedGroups); err != nil {
		return err
	}

	for _, c := range selectedGroups {
		groupAnswers = append(groupAnswers, groupMap[c])
	}

	distribution.Data.ResourceGroups = groupAnswers

	return nil
}

func promptReportDistributionIntegration(distribution *api.ReportDistribution, subReportType string) error {
	var integrations []api.ReportDistributionIntegration

	switch subReportType {
	case api.ReportDefinitionSubTypeAws.String():
		if err := promptReportDistributionIntegrationsAws(&integrations); err != nil {
			return err
		}
	case api.ReportDefinitionSubTypeGcp.String():
		if err := promptReportDistributionIntegrationsGcp(&integrations); err != nil {
			return err
		}
	case api.ReportDefinitionSubTypeAzure.String():
		if err := promptReportDistributionIntegrationsAzure(&integrations); err != nil {
			return err
		}
	default:
		return errors.Errorf("unsupported report definition type '%s'",
			subReportType)
	}

	distribution.Data.Integrations = integrations
	return nil
}

func promptReportDistributionIntegrationsAws(integrations *[]api.ReportDistributionIntegration) error {
	cli.StartProgress("Fetching Aws Account IDs...")
	accounts, err := cli.LwApi.V2.CloudAccounts.ListByType(api.AwsCfgCloudAccount)
	cli.StopProgress()

	if err != nil {
		return err
	}

	var integrationOptions []string

	for _, ca := range accounts.Data {
		if caMap, ok := ca.GetData().(map[string]interface{}); ok {
			integrationOptions = append(integrationOptions, caMap["awsAccountId"].(string))
		}
	}

	var integrationAnswers []string
	if err = survey.AskOne(&survey.MultiSelect{
		Renderer: survey.Renderer{},
		Message:  CreateReportDistributionIntegrationAwsQuestion,
		Options:  integrationOptions,
	}, &integrationAnswers); err != nil {
		return err
	}

	for _, integration := range integrationAnswers {
		*integrations = append(*integrations, api.ReportDistributionIntegration{AccountID: integration})
	}

	return nil
}

func promptReportDistributionIntegrationsGcp(integrations *[]api.ReportDistributionIntegration) error {
	err := promptReportDistributionIntegrationsGcpData(integrations)
	if err != nil {
		return err
	}

	addIntegration := false
	for {
		if err := survey.AskOne(&survey.Confirm{
			Message: "Add another Gcp Integration?",
		}, &addIntegration); err != nil {
			return err
		}

		if addIntegration {
			err = promptReportDistributionIntegrationsGcpData(integrations)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

func promptReportDistributionIntegrationsGcpData(integrations *[]api.ReportDistributionIntegration) error {
	questions := []*survey.Question{
		{
			Name:     "org",
			Prompt:   &survey.Input{Message: "Organization ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "project",
			Prompt:   &survey.Input{Message: "Project ID:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Org     string `survey:"org"`
		Project string `survey:"project"`
	}{}

	if err := survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc)); err != nil {
		return err
	}

	*integrations = append(*integrations, api.ReportDistributionIntegration{OrganizationID: answers.Org,
		ProjectID: answers.Project})

	return nil
}

func promptReportDistributionIntegrationsAzure(integrations *[]api.ReportDistributionIntegration) error {
	err := promptReportDistributionIntegrationsAzureData(integrations)
	if err != nil {
		return err
	}

	addIntegration := false
	for {
		if err := survey.AskOne(&survey.Confirm{
			Message: "Add another Azure Integration?",
		}, &addIntegration); err != nil {
			return err
		}

		if addIntegration {
			err = promptReportDistributionIntegrationsAzureData(integrations)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

func promptReportDistributionIntegrationsAzureData(integrations *[]api.ReportDistributionIntegration) error {
	questions := []*survey.Question{
		{
			Name:     "tenant",
			Prompt:   &survey.Input{Message: "Tenant ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "subscription",
			Prompt:   &survey.Input{Message: "Subscription ID:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Tenant       string `survey:"tenant"`
		Subscription string `survey:"subscription"`
	}{}

	if err := survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc)); err != nil {
		return err
	}

	*integrations = append(*integrations, api.ReportDistributionIntegration{TenantID: answers.Tenant,
		SubscriptionID: answers.Subscription})

	return nil
}

func promptAddReportDistributionViolations(integrations []string) error {
	addViolations := false
	if err := survey.AskOne(&survey.Confirm{
		Message: CreateReportDistributionAddViolationsQuestion,
	}, &addViolations); err != nil {
		return err
	}

	if addViolations {
		if err := survey.AskOne(&survey.MultiSelect{
			Renderer: survey.Renderer{},
			Message:  CreateReportDistributionViolationsQuestion,
			Options:  api.ReportDistributionViolations(),
		}, &integrations); err != nil {
			return err
		}
	}

	return nil
}

func promptAddReportDistributionSeverities(integrations []string) error {
	addSevs := false
	if err := survey.AskOne(&survey.Confirm{
		Message: CreateReportDistributionAddSeveritiesQuestion,
	}, &addSevs); err != nil {
		return err
	}

	if addSevs {
		if err := survey.AskOne(&survey.MultiSelect{
			Renderer: survey.Renderer{},
			Message:  CreateReportDistributionSeveritiesQuestion,
			Options:  lwseverity.ValidSeveritiesStrings,
		}, &integrations); err != nil {
			return err
		}
	}

	return nil
}
