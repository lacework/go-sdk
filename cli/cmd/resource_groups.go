//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
)

var (
	// resource_groups command is used to manage lacework resource groups
	resourceGroupsCommand = &cobra.Command{
		Use:     "resource_group",
		Aliases: []string{"group"},
		Short:   "manage resource groups",
		Long:    "Manage resource groups.",
	}

	// list command is used to list all lacework resource groups
	resourceGroupsListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all resource groups",
		Long:    "List all resource groups configured in your Lacework account.",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			resourceGroups, err := cli.LwApi.V2.ResourceGroups.List()
			if err != nil {
				return errors.Wrap(err, "unable to get resource groups")
			}
			if len(resourceGroups.Data) == 0 {
				msg := `There are no resource groups configured in your account.

Get started by integrating your resource groups to manage alerting using the command:

    $ lacework resource_group create

If you prefer to configure resource groups via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Resource Groups.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}

			groups := make([]resourceGroup, 0)
			for _, g := range resourceGroups.Data {

				groups = append(groups, resourceGroup{
					Id:      g.ResourceGuid,
					ResType: g.Type,
					Name:    g.Name,
					State:   g.Status(),
				})
			}

			if cli.JSONOutput() {
				jsonOut := struct {
					Groups []resourceGroup `json:"resource_groups"`
				}{Groups: groups}
				return cli.OutputJSON(jsonOut)
			}

			rows := [][]string{}
			for _, g := range groups {
				rows = append(rows, []string{g.Id, g.ResType, g.Name, g.State})
			}

			cli.OutputHuman(renderSimpleTable([]string{"RESOURCE ID", "TYPE", "NAME", "STATE"}, rows))
			return nil
		},
	}
)

func init() {
	// add the resource_group command
	rootCmd.AddCommand(resourceGroupsCommand)

	// add sub-commands to the resource_group command
	resourceGroupsCommand.AddCommand(resourceGroupsListCommand)
}

type resourceGroup struct {
	Id      string `json:"resource_guid"`
	ResType string `json:"type"`
	Name    string `json:"name"`
	State   string `json:"state"`
}
