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
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// integrationCmd represents the integration command
	integrationCmd = &cobra.Command{
		Use:     "integration",
		Aliases: []string{"int"},
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

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Integration GUID", "Name", "Type", "Status", "State"})
			table.SetBorder(false)
			table.AppendBulk(integrations.Table())
			table.Render()
			return nil
		},
	}

	// integrationCreateCmd represents the create sub-command inside the integration command
	integrationCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create an external integrations",
		Args:  cobra.NoArgs,
		Long: `This command will prompt an interactive session that will help you create
a new Lacework external integration. If the flag '--noninteractive' is provided,
this command will be disabled.`,
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
		Long: `Delete an external integration by providing its integration GUID, to find the
list of integration configured on your Lacework account use the command:

  $ lacework int list`,
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
				//"Docker V2 Registry",
				//"Amazon Container Registry",
				//"Google Container Registry",
				//"Azure Config",
				//"Azure Activity Log",
				//"GCP Config",
				//"GCP Audit Trail",
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
		//case "Docker V2 Registry":
		//case "Amazon Container Registry":
		//case "Google Container Registry":
		//case "Azure Config":
		//case "Azure Activity Log":
		//case "GCP Config":
		//case "GCP Audit Trail":
		//case "Snowflake Data Share":
	default:
		return errors.New("unknown integration type")
	}
}
