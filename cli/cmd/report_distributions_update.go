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

// update command is used to update an existing lacework report distribution
var reportDistributionsUpdateCommand = &cobra.Command{
	Use:   "update <report_distribution_id>",
	Short: "Update an existing report distribution",
	Long: `Update an existing report distribution.

To update a report distribution:

    lacework report-distribution update <report_distribution_id>
`,
	Args: cobra.ExactArgs(1),
	RunE: updateReportDistribution,
}

func updateReportDistribution(_ *cobra.Command, args []string) error {
	var (
		reportDistribution api.ReportDistributionUpdate
		err                error
	)

	existing, err := cli.LwApi.V2.ReportDistributions.Get(args[0])
	if err != nil {
		return err
	}

	reportDistribution, err = promptUpdateReportDistribution(existing.Data)
	if err != nil {
		return err
	}

	cli.StartProgress("Updated report distribution...")
	resp, err := cli.LwApi.V2.ReportDistributions.Update(args[0], reportDistribution)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to update report distribution")
	}

	cli.OutputHuman("Report distribution updated. To view the report run:\n\n"+
		"lacework report-distribution show %s \n", resp.Data.ReportDistributionGuid)
	return nil
}

func promptUpdateReportDistribution(existing api.ReportDistribution) (reportDistribution api.ReportDistributionUpdate, err error) {
	cli.StartProgress("Fetching list of alert channels...")
	definition, err := cli.LwApi.V2.ReportDefinitions.Get(existing.ReportDefinitionGuid)
	cli.StopProgress()
	if err != nil {
		return api.ReportDistributionUpdate{}, err
	}

	cli.StartProgress("Fetching list of alert channels...")
	alertChannels, err := cli.LwApi.V2.AlertChannels.List()
	cli.StopProgress()
	if err != nil {
		return api.ReportDistributionUpdate{}, err
	}

	channelMap := make(map[string]string, len(alertChannels.Data))
	channelIDMap := make(map[string]string, len(alertChannels.Data))
	var channelOptions []string
	var channelDefaults []string

	for _, def := range alertChannels.Data {
		channelMap[def.Name] = def.IntgGuid
		channelIDMap[def.IntgGuid] = def.Name
		channelOptions = append(channelOptions, def.Name)
	}

	for _, def := range existing.AlertChannels {
		channelDefaults = append(channelDefaults, channelIDMap[def])
	}

	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: UpdateReportDistributionReportNameQuestion, Default: existing.DistributionName},
			Validate: survey.Required,
		},
		{
			Name:     "frequency",
			Prompt:   &survey.Select{Message: UpdateReportDistributionFrequencyQuestion, Options: api.ReportDistributionFrequencies(), Default: existing.Frequency},
			Validate: survey.Required,
		},
		{
			Name:     "channels",
			Prompt:   &survey.MultiSelect{Message: UpdateReportDistributionAlertChannelsQuestion, Options: channelOptions, Default: channelDefaults},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name         string   `survey:"name"`
		Frequency    string   `survey:"frequency"`
		Definition   string   `survey:"definition"`
		Channels     []string `survey:"channels"`
		Integrations []string `survey:"integrations"`
	}{}

	if err = survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc)); err != nil {
		return
	}

	var channelAnswers []string
	for _, c := range answers.Channels {
		channelAnswers = append(channelAnswers, channelMap[c])
	}

	reportDistribution = api.ReportDistributionUpdate{
		DistributionName: answers.Name,
		AlertChannels:    channelAnswers,
		Frequency:        answers.Frequency,
	}

	// scope can only be either resource group or integration
	if err = promptUpdateReportDistributionScope(&reportDistribution,
		definition.Data.SubReportType, existing); err != nil {
		return api.ReportDistributionUpdate{}, err
	}

	// prompt optional fields
	var violations []string
	if err = promptUpdateReportDistributionViolations(violations, existing); err != nil {
		return api.ReportDistributionUpdate{}, err
	}
	reportDistribution.Data.Violations = violations

	var severities []string
	if err = promptUpdateReportDistributionSeverities(severities, existing); err != nil {
		return api.ReportDistributionUpdate{}, err
	}
	reportDistribution.Data.Severities = severities

	return
}

func promptUpdateReportDistributionScope(distribution *api.ReportDistributionUpdate, subReportType string, existing api.ReportDistribution) error {
	distributionScope := ""
	if err := survey.AskOne(&survey.Select{
		Message: CreateReportDistributionScopeQuestion,
		Options: api.ReportDistributionScopes(),
	}, &distributionScope); err != nil {
		return err
	}

	if distributionScope == api.ReportDistributionScopeCloudIntegration.String() {
		if err := promptUpdateReportDistributionIntegration(distribution, subReportType, existing); err != nil {
			return err
		}
		// clear resource group to avoid conflict
		distribution.Data.ResourceGroups = []string{}
		return nil
	} else {
		if err := promptUpdateReportDistributionResourceGroup(distribution, existing); err != nil {
			return err
		}
		// clear integrations to avoid conflict
		distribution.Data.Integrations = []api.ReportDistributionIntegration{}
		return nil
	}
}

func promptUpdateReportDistributionResourceGroup(distribution *api.ReportDistributionUpdate, existing api.ReportDistribution) error {
	cli.StartProgress("Fetching list of resource groups...")
	resourceGroups, err := cli.LwApi.V2.ResourceGroups.List()
	cli.StopProgress()
	if err != nil {
		return err
	}

	groupMap := make(map[string]string, len(resourceGroups.Data))
	groupIDMap := make(map[string]string, len(resourceGroups.Data))
	var (
		groupOptions  []string
		groupDefaults []string
		groupAnswers  []string
	)

	for _, group := range resourceGroups.Data {
		groupMap[group.Name] = group.ID()
		groupIDMap[group.ID()] = group.Name
		groupOptions = append(groupOptions, group.Name)
	}

	for _, group := range existing.Data.ResourceGroups {
		groupDefaults = append(groupDefaults, groupIDMap[group])
	}

	var selectedGroups []string
	if err = survey.AskOne(&survey.MultiSelect{
		Message: CreateReportDistributionResourceGroupsQuestion,
		Options: groupOptions,
		Default: groupDefaults,
	}, &selectedGroups); err != nil {
		return err
	}

	for _, c := range selectedGroups {
		groupAnswers = append(groupAnswers, groupMap[c])
	}

	distribution.Data.ResourceGroups = groupAnswers
	return nil
}

func promptUpdateReportDistributionIntegration(distribution *api.ReportDistributionUpdate, reportType string, existing api.ReportDistribution) error {
	var integrations []api.ReportDistributionIntegration

	switch reportType {
	case api.ReportDefinitionSubTypeAws.String():
		if err := promptUpdateReportDistributionIntegrationsAws(&integrations, existing); err != nil {
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
			reportType)
	}

	distribution.Data.Integrations = integrations
	return nil
}

func promptUpdateReportDistributionViolations(integrations []string, existing api.ReportDistribution) error {
	addViolations := false
	if err := survey.AskOne(&survey.Confirm{
		Message: UpdateReportDistributionAddViolationsQuestion,
	}, &addViolations); err != nil {
		return err
	}

	if addViolations {
		if err := survey.AskOne(&survey.MultiSelect{
			Renderer: survey.Renderer{},
			Message:  CreateReportDistributionViolationsQuestion,
			Options:  api.ReportDistributionViolations(),
			Default:  existing.Data.Violations,
		}, &integrations); err != nil {
			return err
		}
	}

	return nil
}

func promptUpdateReportDistributionSeverities(integrations []string, existing api.ReportDistribution) error {
	addSevs := false
	if err := survey.AskOne(&survey.Confirm{
		Message: UpdateReportDistributionAddSeveritiesQuestion,
	}, &addSevs); err != nil {
		return err
	}

	if addSevs {
		if err := survey.AskOne(&survey.MultiSelect{
			Renderer: survey.Renderer{},
			Message:  CreateReportDistributionSeveritiesQuestion,
			Options:  lwseverity.ValidSeveritiesStrings,
			Default:  existing.Data.Severities,
		}, &integrations); err != nil {
			return err
		}
	}

	return nil
}

func promptUpdateReportDistributionIntegrationsAws(integrations *[]api.ReportDistributionIntegration, existing api.ReportDistribution) error {
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

	var existingAccounts []string
	for _, integration := range existing.Data.Integrations {
		existingAccounts = append(existingAccounts, integration.AccountID)
	}

	var integrationAnswers []string
	if err = survey.AskOne(&survey.MultiSelect{
		Renderer: survey.Renderer{},
		Message:  CreateReportDistributionIntegrationAwsQuestion,
		Options:  integrationOptions,
		Default:  existingAccounts,
	}, &integrationAnswers); err != nil {
		return err
	}

	for _, integration := range integrationAnswers {
		*integrations = append(*integrations, api.ReportDistributionIntegration{AccountID: integration})
	}

	return nil
}
