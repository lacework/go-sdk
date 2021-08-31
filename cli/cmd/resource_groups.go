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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// resource_groups command is used to manage lacework resource groups
	resourceGroupsCommand = &cobra.Command{
		Use:     "resource-group",
		Aliases: []string{"group", "rg"},
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
	// show command is used to retrieve a lacework resource group by resource id
	resourceGroupsShowCommand = &cobra.Command{
		Use:   "show",
		Short: "get resource group by id",
		Long:  "Get a single resource group by it's Resource ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var response api.ResourceGroupResponse
			err := cli.LwApi.V2.ResourceGroups.Get(args[0], &response)
			if err != nil {
				return errors.Wrap(err, "unable to get resource groups")
			}

			group := resourceGroup{
				Id:      response.Data.ResourceGuid,
				ResType: response.Data.Type,
				Name:    response.Data.Name,
				State:   response.Data.Status(),
				Props:   response.Data.Props,
			}

			if cli.JSONOutput() {
				jsonOut := struct {
					Group resourceGroup `json:"resource_group"`
				}{Group: group}
				return cli.OutputJSON(jsonOut)
			}

			groupCommon := [][]string{}
			groupCommon = append(groupCommon, []string{group.Id, group.ResType, group.Name, group.State})

			cli.OutputHuman(renderSimpleTable([]string{"RESOURCE ID", "TYPE", "NAME", "STATE"}, groupCommon))
			cli.OutputHuman("\n")
			cli.OutputHuman(buildResourceGroupPropsTable(group))

			return nil
		},
	}
)

func buildResourceGroupPropsTable(group resourceGroup) string {
	props := determineResourceGroupProps(group.ResType, group.Props)

	return renderOneLineCustomTable("RESOURCE GROUP PROPS",
		renderCustomTable([]string{}, props,
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

func determineResourceGroupProps(resType string, props interface{}) [][]string {
	details := setBaseProps(props)

	switch resType {
	case api.AwsResourceGroup.String():
		details = append(details, setAwsProps(props))
	case api.AzureResourceGroup.String():
		details = append(details, setAzureProps(props))
	case api.ContainerResourceGroup.String():
		details = append(details, setContainerProps(props)...)
	case api.GcpResourceGroup.String():
		details = append(details, setGcpProps(props)...)
	case api.LwAccountResourceGroup.String():
		details = append(details, setLwAccountProps(props))
	case api.MachineResourceGroup.String():
		details = append(details, setMachineProps(props))
	}

	return details
}

func setBaseProps(props interface{}) [][]string {
	var (
		baseProps resourceGroupPropsBase
		details   [][]string
	)

	err := json.Unmarshal([]byte(props.(string)), &baseProps)
	if err != nil {
		return [][]string{}
	}

	details = append(details, []string{"DESCRIPTION", baseProps.Description})
	details = append(details, []string{"UPDATED BY", baseProps.UpdatedBy})
	details = append(details, []string{"LAST UPDATED", strconv.Itoa(baseProps.LastUpdated)})
	return details
}

func setAwsProps(group interface{}) []string {
	var awsProps api.AwsResourceGroupProps
	err := json.Unmarshal([]byte(group.(string)), &awsProps)
	if err != nil {
		return []string{}
	}

	return []string{"ACCOUNT IDS", strings.Join(awsProps.AccountIDs, ",")}
}

func setAzureProps(group interface{}) []string {
	var azProps api.AzureResourceGroupProps
	err := json.Unmarshal([]byte(group.(string)), &azProps)
	if err != nil {
		return []string{}
	}

	return []string{"ACCOUNT IDS", strings.Join(azProps.Subscriptions, ",")}
}

func setContainerProps(group interface{}) [][]string {
	var (
		ctrProps api.ContainerResourceGroupProps
		labels   []string
		details  [][]string
	)
	err := json.Unmarshal([]byte(group.(string)), &ctrProps)
	if err != nil {
		return [][]string{}
	}

	for _, labelMap := range ctrProps.ContainerLabels {
		for key, val := range labelMap {
			labels = append(labels, fmt.Sprintf("%s: %v", key, val))

		}
	}
	details = append(details, []string{"CONTAINER LABELS", strings.Join(labels, ",")})
	details = append(details, []string{"CONTAINER TAGS", strings.Join(ctrProps.ContainerTags, ",")})
	return details
}

func setGcpProps(group interface{}) [][]string {
	var (
		gcpProps api.GcpResourceGroupProps
		details  [][]string
	)
	err := json.Unmarshal([]byte(group.(string)), &gcpProps)
	if err != nil {
		return [][]string{}
	}

	details = append(details, []string{"ORGANIZATION", gcpProps.Organization})
	details = append(details, []string{"PROJECTS", strings.Join(gcpProps.Projects, ",")})
	return details
}

func setLwAccountProps(group interface{}) []string {
	var lwProps api.LwAccountResourceGroupProps
	err := json.Unmarshal([]byte(group.(string)), &lwProps)
	if err != nil {
		return []string{}
	}

	return []string{"LW ACCOUNTS", strings.Join(lwProps.LwAccounts, ",")}
}

func setMachineProps(group interface{}) []string {
	var machineProps api.MachineResourceGroupProps

	err := json.Unmarshal([]byte(group.(string)), &machineProps)
	if err != nil {
		return []string{}
	}

	var tags []string
	for _, tagMap := range machineProps.MachineTags {
		for key, val := range tagMap {
			tags = append(tags, fmt.Sprintf("%s: %v", key, val))

		}
	}
	return []string{"MACHINE TAGS", strings.Join(tags, ",")}
}

func init() {
	// add the resource_group command
	rootCmd.AddCommand(resourceGroupsCommand)

	// add sub-commands to the resource_group command
	resourceGroupsCommand.AddCommand(resourceGroupsListCommand)
	resourceGroupsCommand.AddCommand(resourceGroupsShowCommand)
}

type resourceGroup struct {
	Id      string      `json:"resource_guid"`
	ResType string      `json:"type"`
	Name    string      `json:"name"`
	State   string      `json:"state"`
	Props   interface{} `json:"props"`
}

type resourceGroupPropsBase struct {
	Description string `json:"description"`
	UpdatedBy   string `json:"UPDATED_BY,omitempty"`
	LastUpdated int    `json:"LAST_UPDATED,omitempty"`
}
