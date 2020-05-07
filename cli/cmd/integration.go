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
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// integrationCmd represents the integration command
	integrationCmd = &cobra.Command{
		Use:     "integration",
		Aliases: []string{"integrations", "int"},
		Short:   "manage external integrations",
		Long:    `Manage external integrations with the Lacework platform`,
	}

	// integrationListCmd represents the list sub-command inside the integration command
	integrationListCmd = &cobra.Command{
		Use:   "list",
		Short: "list all available external integrations",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			lacework, err := api.NewClient(cli.Account,
				api.WithLogLevel(cli.LogLevel),
				api.WithApiKeys(cli.KeyID, cli.Secret),
			)
			if err != nil {
				return errors.Wrap(err, "unable to generate api client")
			}

			integrations, err := lacework.Integrations.List()
			if err != nil {
				return errors.Wrap(err, "unable to get integrations")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(integrations.Data)
			}

			cli.OutputHuman(buildIntegrationsTable(integrations.Data))
			return nil
		},
	}

	// integrationShowCmd represents the show sub-command inside the integration command
	integrationShowCmd = &cobra.Command{
		Use:   "show <int_guid>",
		Short: "Show details about a specific external integration",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			lacework, err := api.NewClient(cli.Account,
				api.WithLogLevel(cli.LogLevel),
				api.WithApiKeys(cli.KeyID, cli.Secret),
			)
			if err != nil {
				return errors.Wrap(err, "unable to generate api client")
			}

			integration, err := lacework.Integrations.Get(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to get integration")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(integration.Data)
			}

			cli.OutputHuman(buildIntegrationsTable(integration.Data))
			cli.OutputHuman("\n")
			cli.OutputHuman(buildIntDetailsTable(integration.Data))
			return nil
		},
	}

	// integrationCreateCmd represents the create sub-command inside the integration command
	integrationCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create an external integrations",
		Args:  cobra.NoArgs,
		Long:  `Creates an external integration in your account through an interactive session.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if !cli.InteractiveMode() {
				return errors.New("interactive mode is disabled")
			}
			lacework, err := api.NewClient(cli.Account,
				api.WithLogLevel(cli.LogLevel),
				api.WithApiKeys(cli.KeyID, cli.Secret),
			)
			if err != nil {
				return errors.Wrap(err, "unable to generate api client")
			}

			err = promptCreateIntegration(lacework)
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
		Short:  "update an external integrations",
		Args:   cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	// integrationDeleteCmd represents the delete sub-command inside the integration command
	integrationDeleteCmd = &cobra.Command{
		Use:   "delete <int_guid>",
		Short: "delete an external integrations",
		Long: `Delete an external integration by providing its integration GUID. Integration
GUIDs can be found by using the 'lacework integration list' command.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			lacework, err := api.NewClient(cli.Account,
				api.WithLogLevel(cli.LogLevel),
				api.WithApiKeys(cli.KeyID, cli.Secret),
			)
			if err != nil {
				return errors.Wrap(err, "unable to generate api client")
			}

			cli.Log.Info("deleting integration", "int_guid", args[0])
			cli.StartProgress(" Deleting integration...")
			response, err := lacework.Integrations.Delete(args[0])
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
}

func promptCreateIntegration(lacework *api.Client) error {
	var (
		integration = ""
		prompt      = &survey.Select{
			Message: "Choose an integration type to create: ",
			Options: []string{
				"Docker Hub",
				"AWS Config",
				"AWS CloudTrail",
				"GCP Config",
				"GCP Audit Log",
				"Azure Config",
				"Azure Activity Log",
				//"Docker V2 Registry",
				//"Amazon Container Registry",
				//"Google Container Registry",
				//"Snowflake Data Share",
			},
		}
		err = survey.AskOne(prompt, &integration)
	)
	if err != nil {
		return err
	}

	switch integration {
	case "Docker Hub":
		return createDockerHubIntegration(lacework)
	case "AWS Config":
		return createAwsConfigIntegration(lacework)
	case "AWS CloudTrail":
		return createAwsCloudTrailIntegration(lacework)
	case "GCP Config":
		return createGcpConfigIntegration(lacework)
	case "GCP Audit Log":
		return createGcpAuditLogIntegration(lacework)
	case "Azure Config":
		return createAzureConfigIntegration(lacework)
	case "Azure Activity Log":
		return createAzureActivityLogIntegration(lacework)
	//case "Docker V2 Registry":
	//case "Amazon Container Registry":
	//case "Google Container Registry":
	//case "Snowflake Data Share":
	default:
		return errors.New("unknown integration type")
	}
}

func integrationsTable(integrations []api.RawIntegration) [][]string {
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

func buildIntegrationsTable(integrations []api.RawIntegration) string {
	var (
		tableBuilder = &strings.Builder{}
		t            = tablewriter.NewWriter(tableBuilder)
	)

	t.SetHeader([]string{
		"Integration GUID",
		"Name",
		"Type",
		"Status",
		"State",
	})
	t.SetBorder(false)
	t.AppendBulk(integrationsTable(integrations))
	t.Render()

	return tableBuilder.String()
}

func buildIntDetailsTable(integrations []api.RawIntegration) string {
	var (
		main    = &strings.Builder{}
		details = &strings.Builder{}
		t       = tablewriter.NewWriter(details)
	)

	t.SetBorder(false)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	if len(integrations) != 0 {
		integration := integrations[0]
		t.AppendBulk(reflectIntegrationData(integration))
		t.AppendBulk(buildIntegrationState(integration.State))
	}
	t.Render()

	t = tablewriter.NewWriter(main)
	t.SetBorder(false)
	t.SetAutoWrapText(false)
	t.SetHeader([]string{"INTEGRATION DETAILS"})
	t.Append([]string{details.String()})
	t.Render()

	return main.String()
}

func buildIntegrationState(state *api.IntegrationState) [][]string {
	if state != nil {
		return [][]string{
			[]string{"LAST UPDATED TIME", state.LastUpdatedTime},
			[]string{"LAST SUCCESSFUL TIME", state.LastSuccessfulTime},
		}
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
			[]string{"LEVEL", iData.IDType},
			[]string{"ORG/PROJECT ID", iData.ID},
			[]string{"CLIENT ID", iData.Credentials.ClientID},
			[]string{"CLIENT EMAIL", iData.Credentials.ClientEmail},
			[]string{"PRIVATE KEY ID", iData.Credentials.PrivateKeyID},
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
			[]string{"ROLE ARN", iData.Credentials.RoleArn},
			[]string{"EXTERNAL ID", iData.Credentials.ExternalID},
		}
		if iData.QueueUrl != "" {
			return append(out, []string{"QUEUE URL", iData.QueueUrl})
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
			[]string{"CLIENT ID", iData.Credentials.ClientID},
			[]string{"CLIENT SECRET", iData.Credentials.ClientSecret},
			[]string{"TENANT ID", iData.TenantID},
		}
		if iData.QueueUrl != "" {
			return append(out, []string{"QUEUE URL", iData.QueueUrl})
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
		return out
	}

	for key, value := range m {
		if s, ok := value.(string); ok {
			out[key] = s
		} else {
			deepMap := deepKeyValueExtract(value)
			for deepK, deepV := range deepMap {
				out[deepK] = deepV
			}
		}
	}

	return out
}
