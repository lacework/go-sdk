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
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// used by integration list to list only a single type of integration
	integrationType string

	// integrationCmd represents the integration command
	integrationCmd = &cobra.Command{
		Use:     "integration",
		Aliases: []string{"integrations", "int"},
		Short:   "Manage external integrations",
		Long:    `Manage external integrations with the Lacework platform`,
	}

	// integrationListCmd represents the list sub-command inside the integration command
	integrationListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available external integrations",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			var (
				integrations api.RawIntegrationsResponse
				err          error
			)

			if integrationType != "" {
				intType, found := api.FindIntegrationType(integrationType)
				if !found {
					return errors.Errorf("unknown integration type '%s'", integrationType)
				}
				integrations, err = cli.LwApi.Integrations.ListByType(intType)
			} else {
				integrations, err = cli.LwApi.Integrations.List()
			}
			if err != nil {
				return errors.Wrap(err, "unable to get integrations")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(integrations.Data)
			}

			if len(integrations.Data) == 0 {
				cli.OutputHuman("There was no integration found.\n")
				return nil
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"Integration GUID", "Name", "Type", "Status", "State"},
					integrationsToTable(integrations.Data),
				),
			)
			return nil
		},
	}

	// integrationShowCmd represents the show sub-command inside the integration command
	integrationShowCmd = &cobra.Command{
		Use:   "show <int_guid>",
		Short: "Show details about a specific external integration",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			integration, err := cli.LwApi.Integrations.Get(args[0])

			if err != nil {
				return errors.Wrap(err, "unable to get integration")
			}

			if len(integration.Data) == 0 {
				msg := "the provided integration GUID was not found\n\n"
				msg += "To list the available integrations in your account run 'lacework integrations list'"
				return errors.New(msg)
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(integration.Data[0])
			}

			integrationType, supported := api.FindIntegrationType(integration.Data[0].Type)
			if supported {
				var resp api.V2CommonIntegration
				err = cli.LwApi.V2.Schemas.GetService(integrationType.Schema()).Get(args[0], &resp)

				if err != nil {
					cli.Log.Debugw("unable to get integration service", "error", err.Error())
				}

				if resp.Data.State != nil {
					integration.Data[0].State.Details = resp.Data.State.Details
				}
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"Integration GUID", "Name", "Type", "Status", "State"},
					integrationsToTable(integration.Data),
				),
			)

			cli.OutputHuman("\n")
			cli.OutputHuman(buildIntDetailsTable(integration.Data))
			return nil
		},
	}

	// integrationCreateCmd represents the create sub-command inside the integration command
	integrationCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an external integrations",
		Args:  cobra.NoArgs,
		Long:  `Creates an external integration in your account through an interactive session.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if !cli.InteractiveMode() {
				return errors.New("interactive mode is disabled")
			}

			err := promptCreateIntegration()
			if err != nil {
				return errors.Wrap(err, "unable to create integration")
			}

			cli.OutputHuman("The integration was created.\n")
			return nil
		},
	}

	// integrationUpdateCmd represents the update sub-command inside the integration command
	integrationUpdateCmd = &cobra.Command{
		Use:    "update",
		Hidden: true,
		Short:  "Update an external integrations",
		Args:   cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	// integrationDeleteCmd represents the delete sub-command inside the integration command
	integrationDeleteCmd = &cobra.Command{
		Use:   "delete <int_guid>",
		Short: "Delete an external integrations",
		Long: `Delete an external integration by providing an integration GUID.

Integration GUIDs can be found by using the 'lacework integration list' command.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cli.Log.Info("deleting integration", "int_guid", args[0])
			cli.StartProgress(" Deleting integration...")
			response, err := cli.LwApi.Integrations.Delete(args[0])
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to delete integration")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response.Data)
			}

			cli.OutputHuman("The integration %s was deleted.\n", args[0])
			return nil
		},
	}
)

func init() {
	// add the integration command
	rootCmd.AddCommand(integrationCmd)

	// add sub-commands to the integration command
	integrationCmd.AddCommand(integrationListCmd)
	integrationCmd.AddCommand(integrationShowCmd)
	integrationCmd.AddCommand(integrationCreateCmd)
	integrationCmd.AddCommand(integrationUpdateCmd)
	integrationCmd.AddCommand(integrationDeleteCmd)

	// add type flag to integration list command
	integrationListCmd.Flags().StringVarP(&integrationType,
		"type", "t", "", "list all integrations of a specific type",
	)
}

func promptCreateIntegration() error {
	var (
		integration = ""
		prompt      = &survey.Select{
			Message: "Choose an integration type to create: ",
			Options: []string{
				"Slack Alert Channel",
				"Email Alert Channel",
				"Amazon S3 Alert Channel",
				"Cisco Webex Alert Channel",
				"Datadog Alert Channel",
				"GCP PubSub Alert Channel",
				"Microsoft Teams Alert Channel",
				"New Relic Insights Alert Channel",
				"Webhook Alert Channel",
				"VictorOps Alert Channel",
				"Splunk Alert Channel",
				"QRadar Alert Channel",
				"Service Now Alert Channel",
				"PagerDuty Alert Channel",
				"Amazon CloudWatch Alert Channel",
				"Jira Cloud Alert Channel",
				"Jira Server Alert Channel",
				"Docker Hub Registry",
				"Docker V2 Registry",
				"Amazon Container Registry (ECR)",
				"Google Container Registry (GCR)",
				"Google Artifact Registry (GAR)",
				"Github Container Registry (GHCR)",
				"AWS Config",
				"AWS CloudTrail",
				"AWS Config (US GovCloud)",
				"AWS CloudTrail (US GovCloud)",
				"GCP Config",
				"GCP Audit Log",
				"Azure Config",
				"Azure Activity Log",
				//"Snowflake Data Share",
			},
		}
		err = survey.AskOne(prompt, &integration)
	)
	if err != nil {
		return err
	}

	switch integration {
	case "Slack Alert Channel":
		return createSlackAlertChannelIntegration()
	case "Email Alert Channel":
		return createEmailAlertChannelIntegration()
	case "GCP PubSub Alert Channel":
		return createGcpPubSubChannelIntegration()
	case "Microsoft Teams Alert Channel":
		return createMicrosoftTeamsChannelIntegration()
	case "New Relic Insights Alert Channel":
		return createNewRelicAlertChannelIntegration()
	case "Amazon S3 Alert Channel":
		return createAwsS3ChannelIntegration()
	case "Cisco Webex Alert Channel":
		return createCiscoWebexChannelIntegration()
	case "Datadog Alert Channel":
		return createDatadogIntegration()
	case "Webhook Alert Channel":
		return createWebhookIntegration()
	case "VictorOps Alert Channel":
		return createVictorOpsChannelIntegration()
	case "Splunk Alert Channel":
		return createSplunkIntegration()
	case "PagerDuty Alert Channel":
		return createPagerDutyAlertChannelIntegration()
	case "QRadar Alert Channel":
		return createQRadarAlertChannelIntegration()
	case "Service Now Alert Channel":
		return createServiceNowAlertChannelIntegration()
	case "Amazon CloudWatch Alert Channel":
		return createAwsCloudWatchAlertChannelIntegration()
	case "Jira Cloud Alert Channel":
		return createJiraCloudAlertChannelIntegration()
	case "Jira Server Alert Channel":
		return createJiraServerAlertChannelIntegration()
	case "Docker Hub Registry":
		return createDockerHubIntegration()
	case "Docker V2 Registry":
		return createDockerV2Integration()
	case "Amazon Container Registry (ECR)":
		return createAwsEcrIntegration()
	case "Google Artifact Registry (GAR)":
		return createGarIntegration()
	case "Github Container Registry (GHCR)":
		return createGhcrIntegration()
	case "Google Container Registry (GCR)":
		return createGcrIntegration()
	case "AWS Config":
		return createAwsConfigIntegration()
	case "AWS CloudTrail":
		return createAwsCloudTrailIntegration()
	case "AWS GovCloud Config":
		return createAwsGovCloudConfigIntegration()
	case "AWS GovCloud CloudTrail":
		return createAwsGovCloudCTIntegration()
	case "GCP Config":
		return createGcpConfigIntegration()
	case "GCP Audit Log":
		return createGcpAuditLogIntegration()
	case "Azure Config":
		return createAzureConfigIntegration()
	case "Azure Activity Log":
		return createAzureActivityLogIntegration()
	//case "Snowflake Data Share":
	default:
		return errors.New("unknown integration type")
	}
}

func integrationsToTable(integrations []api.RawIntegration) [][]string {
	out := [][]string{}
	for _, idata := range integrations {
		out = append(out, []string{
			idata.IntgGuid,
			idata.Name,
			idata.Type,
			idata.Status(),
			idata.StateString(),
		})
	}
	return out
}

func buildIntDetailsTable(integrations []api.RawIntegration) string {
	if len(integrations) == 0 {
		return "ERROR unable to access integration details. No data!\n"
	}

	integration := integrations[0]
	details := reflectIntegrationData(integration)
	details = append(details, []string{"UPDATED AT", integration.CreatedOrUpdatedTime})
	details = append(details, []string{"UPDATED BY", integration.CreatedOrUpdatedBy})
	details = append(details, buildIntegrationState(integration.State)...)

	return renderOneLineCustomTable("INTEGRATION DETAILS",
		renderCustomTable([]string{}, details,
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(false)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
			}),
		),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	)
}

func buildIntegrationState(state *api.IntegrationState) [][]string {
	if state != nil {
		details := [][]string{
			{"STATE UPDATED AT", state.LastUpdatedTime},
			{"LAST SUCCESSFUL STATE", state.LastSuccessfulTime},
		}

		if len(state.Details) != 0 {
			detailsStr, err := json.Marshal(state.Details)
			if err != nil {
				cli.Log.Debugw("unable to marshall state details", "error", err.Error())
				return details
			}

			detailsJSON, err := cli.FormatJSONString(string(detailsStr))
			if err != nil {
				cli.Log.Debugw("unable to json format state details", "error", err.Error())
				return details
			}
			details = append(details, []string{"STATE DETAILS", detailsJSON})
		}
		return details
	}

	return [][]string{}
}

func reflectIntegrationData(raw api.RawIntegration) [][]string {
	switch raw.Type {

	case api.GcpCfgIntegration.String(),
		api.GcpAuditLogIntegration.String():

		var iData api.GcpIntegrationData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"LEVEL", iData.IDType},
			{"ORG/PROJECT ID", iData.ID},
			{"CLIENT ID", iData.Credentials.ClientID},
			{"CLIENT EMAIL", iData.Credentials.ClientEmail},
			{"PRIVATE KEY ID", iData.Credentials.PrivateKeyID},
		}
		if iData.SubscriptionName != "" {
			return append(out, []string{"SUBSCRIPTION NAME", iData.SubscriptionName})
		}
		return out

	case api.AwsCfgIntegration.String(),
		api.AwsCloudTrailIntegration.String():

		var iData api.AwsIntegrationData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"ROLE ARN", iData.Credentials.RoleArn},
			{"EXTERNAL ID", iData.Credentials.ExternalID},
		}
		if iData.QueueUrl != "" {
			out = append(out, []string{"QUEUE URL", iData.QueueUrl})
		}

		accountMapping, err := iData.DecodeAccountMappingFile()
		if err != nil {
			cli.Log.Debugw("unable to decode account mapping file",
				"integration_type", raw.Type,
				"raw_data", iData.AccountMappingFile,
				"error", err,
			)
		}

		if len(accountMapping) != 0 {
			// @afiune should we disable the colors here?
			accountMappingJSON, err := cli.FormatJSONString(string(accountMapping))
			if err != nil {
				accountMappingJSON = string(accountMapping)
			}
			out = append(out, []string{"ACCOUNT MAPPING FILE", accountMappingJSON})
		}
		return out

	case api.AzureCfgIntegration.String(),
		api.AzureActivityLogIntegration.String():

		var iData api.AzureIntegrationData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"CLIENT ID", iData.Credentials.ClientID},
			{"CLIENT SECRET", iData.Credentials.ClientSecret},
			{"TENANT ID", iData.TenantID},
		}
		if iData.QueueUrl != "" {
			return append(out, []string{"QUEUE URL", iData.QueueUrl})
		}
		return out

	case api.SlackChannelIntegration.String():

		var iData api.SlackChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"SLACK URL", iData.SlackUrl},
		}

		return out

	case api.WebhookIntegration.String():

		var iData api.WebhookChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{{"WEBHOOK URL", iData.WebhookUrl}}

		return out

	case api.VictorOpsChannelIntegration.String():

		var iData api.VictorOpsChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"WEBHOOK URL", iData.WebhookURL},
		}

		return out

	case api.SplunkIntegration.String():

		var iData api.SplunkChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"CHANNEL", iData.Channel},
			{"HEC TOKEN", iData.HecToken},
			{"HOST", iData.Host},
			{"PORT", fmt.Sprintf("%d", iData.Port)},
			{"INDEX", iData.EventData.Index},
			{"SOURCE", iData.EventData.Source},
		}
		if iData.Ssl {
			out = append(out, []string{"SSL", "ENABLE"})
		} else {
			out = append(out, []string{"SSL", "DISABLE"})
		}

		return out

	case api.ServiceNowChannelIntegration.String():

		var iData api.ServiceNowChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}

		templateString, err := iData.DecodeCustomTemplateFile()
		if err != nil {
			cli.Log.Debugw("unable to decode custom template file",
				"integration_type", raw.Type,
				"raw_data", iData.CustomTemplateFile,
				"error", err,
			)
		}

		tmplStrPretty, err := cli.FormatJSONString(templateString)
		if err != nil {
			tmplStrPretty = templateString
		}
		out := [][]string{
			{"INSTANCE URL", iData.InstanceURL},
			{"USERNAME", iData.Username},
			{"PASSWORD", iData.Password},
			{"CUSTOM TEMPLATE FILE", tmplStrPretty},
			{"ISSUE GROUPING", iData.IssueGrouping},
		}

		return out

	case api.AwsS3ChannelIntegration.String():

		var iData api.AwsS3ChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"ROLE ARN", iData.Credentials.RoleArn},
			{"BUCKET ARN", iData.Credentials.BucketArn},
			{"EXTERNAL ID", iData.Credentials.ExternalID},
		}

		return out

	case api.QRadarChannelIntegration.String():

		var iData api.QRadarChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"HOST PORT", fmt.Sprint(iData.HostPort)},
			{"HOST URL", iData.HostURL},
			{"COMMUNICATION TYPE", string(iData.CommunicationType)},
		}

		return out

	case api.CiscoWebexChannelIntegration.String():

		var iData api.CiscoWebexChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"WEBHOOK URL", iData.WebhookURL},
		}

		return out

	case api.NewRelicChannelIntegration.String():

		var iData api.NewRelicChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"ACCOUNT ID", fmt.Sprint(iData.AccountID)},
			{"INSERT API KEY", iData.InsertKey},
		}

		return out

	case api.DatadogChannelIntegration.String():

		var iData api.DatadogChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"DATADOG SITE", string(iData.DatadogSite)},
			{"DATADOG SERVICE", string(iData.DatadogService)},
			{"API KEY", iData.ApiKey},
		}

		return out

	case api.GcpPubSubChannelIntegration.String():

		var iData api.GcpPubSubChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"PROJECT ID", iData.ProjectID},
			{"TOPIC ID", iData.TopicID},
			{"CLIENT ID", iData.Credentials.ClientID},
			{"CLIENT EMAIL", iData.Credentials.ClientEmail},
			{"PRIVATE_KEY_ID", iData.Credentials.PrivateKeyID},
			{"ISSUE GROUPING", iData.IssueGrouping},
		}

		return out

	case api.MicrosoftTeamsChannelIntegration.String():

		var iData api.MicrosoftTeamsChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{{"WEBHOOK URL", iData.WebhookURL}}

		return out

	case api.AwsCloudWatchIntegration.String():

		var iData api.AwsCloudWatchData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"EVENT BUS ARN", iData.EventBusArn},
			{"ISSUE GROUPING", iData.IssueGrouping},
		}

		return out

	case api.ContainerRegistryIntegration.String():

		var iData api.ContainerRegData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"REGISTRY TYPE", iData.RegistryType},
			{"REGISTRY DOMAIN", iData.RegistryDomain},
			{"LIMIT BY TAG", iData.LimitByTag},
			{"LIMIT BY LABEL", iData.LimitByLabel},
			{"LIMIT BY REPOSITORY", iData.LimitByRep},
			{"LIMIT NUM IMAGES PER REPO", fmt.Sprintf("%d", iData.LimitNumImg)},
		}

		switch iData.RegistryType {
		case api.DockerHubRegistry.String():
			out = append(out, []string{"USERNAME", iData.Credentials.Username})
			out = append(out, []string{"PASSWORD", iData.Credentials.Password})
		case api.GhcrContainerRegistry.String(),
			api.DockerV2Registry.String():
			out = append(out, []string{"USERNAME", iData.Credentials.Username})
			out = append(out, []string{"PASSWORD", iData.Credentials.Password})
			if iData.Credentials.SSL {
				out = append(out, []string{"SSL", "ENABLE"})
			} else {
				out = append(out, []string{"SSL", "DISABLE"})
			}
		case api.GcpGarContainerRegistry.String(),
			api.GcrRegistry.String():
			out = append(out, []string{"CLIENT ID", iData.Credentials.ClientID})
			out = append(out, []string{"CLIENT EMAIL", iData.Credentials.ClientEmail})
			out = append(out, []string{"PRIVATE KEY ID", iData.Credentials.PrivateKeyID})
		case api.EcrRegistry.String():
			var ecrData api.AwsEcrCommonData
			err := mapstructure.Decode(raw.Data, &ecrData)
			if err != nil {
				cli.Log.Debugw("unable to decode integration data",
					"integration_type", raw.Type,
					"registry_type", iData.RegistryType,
					"raw_data", raw.Data,
					"error", err,
				)
				break
			}

			out = append(out, []string{"AWS AUTH TYPE", ecrData.AwsAuthType})

			switch ecrData.AwsAuthType {
			case api.AwsEcrAccessKey.String():
				var ecrIAMData api.AwsEcrDataWithAccessKeyCreds
				err := mapstructure.Decode(raw.Data, &ecrIAMData)
				if err != nil {
					cli.Log.Debugw("unable to decode ECR integration data",
						"integration_type", raw.Type,
						"registry_type", iData.RegistryType,
						"raw_data", raw.Data,
						"error", err,
					)
					break
				}
				out = append(out, []string{"ACCESS KEY ID", ecrIAMData.Credentials.AccessKeyID})
			case api.AwsEcrIAM.String():
				var ecrCrossAccountData api.AwsEcrDataWithCrossAccountCreds
				err := mapstructure.Decode(raw.Data, &ecrCrossAccountData)
				if err != nil {
					cli.Log.Debugw("unable to decode ECR integration data",
						"integration_type", raw.Type,
						"registry_type", iData.RegistryType,
						"raw_data", raw.Data,
						"error", err,
					)
					break
				}
				out = append(out, []string{"ROLE ARN", ecrCrossAccountData.Credentials.RoleArn})
				out = append(out, []string{"EXTERNAL ID", ecrCrossAccountData.Credentials.ExternalID})
			}
		}

		return out

	case api.PagerDutyIntegration.String():

		var iData api.PagerDutyData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}
		out := [][]string{
			{"INTEGRATION KEY", iData.IntegrationKey},
			{"ISSUE GROUPING", iData.IssueGrouping},
		}

		return out

	case api.EmailIntegration.String():
		// Use v2 endpoint for Email Alert Channel
		emailAlertChan, err := cli.LwApi.V2.AlertChannels.GetEmailUser(raw.IntgGuid)
		if err != nil {
			cli.Log.Debugw("unable to get EmailUser Alert Channel (v2)",
				"error", err.Error(),
			)
			break
		}
		return [][]string{
			{"RECIPIENTS",
				strings.Join(emailAlertChan.Data.Data.ChannelProps.Recipients, "\n")},
		}

	case api.JiraIntegration.String():

		var iData api.JiraAlertChannelData
		err := mapstructure.Decode(raw.Data, &iData)
		if err != nil {
			cli.Log.Debugw("unable to decode integration data",
				"integration_type", raw.Type,
				"raw_data", raw.Data,
				"error", err,
			)
			break
		}

		templateString, err := iData.DecodeCustomTemplateFile()
		if err != nil {
			cli.Log.Debugw("unable to decode custom template file",
				"integration_type", raw.Type,
				"raw_data", iData.CustomTemplateFile,
				"error", err,
			)
		}

		// @afiune should we disable the colors here?
		tmplStrPretty, err := cli.FormatJSONString(templateString)
		if err != nil {
			tmplStrPretty = templateString
		}
		out := [][]string{
			{"JIRA INTEGRATION TYPE", iData.JiraType},
			{"JIRA URL", iData.JiraUrl},
			{"PROJECT KEY", iData.ProjectID},
			{"USERNAME", iData.Username},
			{"ISSUE TYPE", iData.IssueType},
			{"ISSUE GROUPING", iData.IssueGrouping},
			{"CUSTOM TEMPLATE FILE", tmplStrPretty},
		}

		return out

	default:
		out := [][]string{}
		for key, value := range deepKeyValueExtract(raw.Data) {
			out = append(out, []string{key, value})
		}
		return out
	}

	return [][]string{}
}

func deepKeyValueExtract(v interface{}) map[string]string {
	out := map[string]string{}

	m, ok := v.(map[string]interface{})
	if !ok {
		cli.Log.Warnw("unable to parse raw data field", "type", fmt.Sprintf("%T", v))
		return out
	}

	for key, value := range m {
		if s, ok := value.(string); ok {
			out[key] = s
		} else if i, ok := value.(int); ok {
			out[key] = fmt.Sprintf("%d", i)
		} else if i, ok := value.(int32); ok {
			out[key] = fmt.Sprintf("%d", i)
		} else if i, ok := value.(float32); ok {
			out[key] = fmt.Sprintf("%.0f", i)
		} else if i, ok := value.(float64); ok {
			out[key] = fmt.Sprintf("%.0f", i)
		} else if b, ok := value.(bool); ok {
			if b {
				out[key] = "ENABLE"
			} else {
				out[key] = "DISABLE"
			}
		} else {
			deepMap := deepKeyValueExtract(value)
			for deepK, deepV := range deepMap {
				out[deepK] = deepV
			}
		}
	}

	return out
}
