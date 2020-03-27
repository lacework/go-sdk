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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// integrationCmd represents the integration command
	integrationCmd = &cobra.Command{
		Use:   "integration",
		Short: "Manage external integrations",
	}

	// integrationListCmd represents the list sub-command inside the integration command
	instegrationListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available external integrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			lacework, err := api.NewClient(cli.Account,
				api.WithApiKeys(cli.KeyID, cli.Secret),
			)
			if err != nil {
				return errors.Wrap(err, "unable to generate Lacework API client")
			}

			cli.Log.Debugw("api client generated",
				"version", lacework.ApiVersion(),
				"base_url", lacework.URL(),
			)

			integrations, err := lacework.Integrations.List()
			if err != nil {
				return errors.Wrap(err, "unable to get integrations")
			}

			fmt.Println(integrations.String())
			return nil
		},
	}

	// integrationCreateCmd represents the create sub-command inside the integration command
	instegrationCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an external integrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// integrationDeleteCmd represents the delete sub-command inside the integration command
	instegrationDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an external integrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	// add the integration command
	rootCmd.AddCommand(integrationCmd)

	// add sub-commands to the integration command
	integrationCmd.AddCommand(instegrationListCmd)
	integrationCmd.AddCommand(instegrationCreateCmd)
	integrationCmd.AddCommand(instegrationDeleteCmd)
}
