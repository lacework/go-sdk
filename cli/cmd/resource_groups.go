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
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// resource-groups command is used to manage lacework resource groups
	resourceGroupsCommand = &cobra.Command{
		Use:     "resource-group",
		Aliases: []string{"resource-groups", "rg"},
		Short:   "Manage resource groups",
		Long:    "Manage Lacework-identifiable assets via the use of resource groups.",
	}

	// list command is used to list all lacework resource groups
	resourceGroupsListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all resource groups",
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

    lacework resource-group create

If you prefer to configure resource groups via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Resource Groups.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}

			groups := make([]resourceGroup, 0)
			for _, g := range resourceGroups.Data {
				var props api.ResourceGroupProps
				var resourceGroupGuid string
				var isDefault int
				var resourceName string

				if g.GetProps() != nil {
					props, _ = parsePropsType(g.Props, g.Type)
					resourceGroupGuid = g.ResourceGuid
					isDefault = g.IsDefault
					resourceName = g.Name
				} else {
					resourceGroupGuid = g.ResourceGroupGuid
					if *g.IsDefaultBoolean {
						isDefault = 1
					} else {
						isDefault = 0
					}
					resourceName = g.NameV2
				}

				groups = append(groups, resourceGroup{
					Id:        resourceGroupGuid,
					ResType:   g.Type,
					Name:      resourceName,
					status:    g.Status(),
					Props:     props,
					Enabled:   g.Enabled,
					IsDefault: isDefault,
					Query:     g.Query,
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
				rows = append(rows, []string{g.Id, g.ResType, g.Name, g.status, IsDefault(g.IsDefault)})
			}

			cli.OutputHuman(renderSimpleTable([]string{"RESOURCE GUID", "TYPE", "NAME", "STATUS", "DEFAULT"}, rows))
			return nil
		},
	}
	// show command is used to retrieve a lacework resource group by resource id
	resourceGroupsShowCommand = &cobra.Command{
		Use:   "show <resource_group_id>",
		Short: "Get resource group by ID",
		Long:  "Get a single resource group by it's resource group ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var response api.ResourceGroupResponse
			err := cli.LwApi.V2.ResourceGroups.Get(args[0], &response)
			if err != nil {
				return errors.Wrap(err, "unable to get resource group")
			}

			var props api.ResourceGroupProps
			var resourceGroupGuid string
			var isDefault int
			var name string

			if response.Data.Props != nil {
				props, _ = parsePropsType(response.Data.Props, response.Data.Type)
				resourceGroupGuid = response.Data.ResourceGuid
				isDefault = response.Data.IsDefault
				name = response.Data.Name
			} else {
				resourceGroupGuid = response.Data.ResourceGroupGuid
				if *response.Data.IsDefaultBoolean {
					isDefault = 1
				} else {
					isDefault = 0
				}
				name = response.Data.NameV2
			}

			group := resourceGroup{
				Id:          resourceGroupGuid,
				ResType:     response.Data.Type,
				Name:        name,
				status:      response.Data.Status(),
				Props:       props,
				Enabled:     response.Data.Enabled,
				IsDefault:   isDefault,
				Query:       response.Data.Query,
				Description: response.Data.Description,
				UpdatedBy:   response.Data.UpdatedBy,
				UpdatedTime: response.Data.UpdatedTime,
			}

			if cli.JSONOutput() {
				jsonOut := struct {
					Group resourceGroup `json:"resource_group"`
				}{Group: group}
				return cli.OutputJSON(jsonOut)
			}

			var groupCommon [][]string

			if group.Props != nil {
				groupCommon = append(groupCommon,
					[]string{group.Id, group.ResType, group.Name, group.status, IsDefault(group.IsDefault)},
				)
				cli.OutputHuman(renderSimpleTable([]string{"RESOURCE ID", "TYPE", "NAME", "STATE", "DEFAULT"}, groupCommon))
				cli.OutputHuman("\n")
				cli.OutputHuman(buildResourceGroupPropsTable(group))
			} else {
				groupCommon = append(groupCommon,
					[]string{group.Id, group.ResType, group.Name, group.Description, group.status,
						IsDefault(group.IsDefault), group.UpdatedBy, group.UpdatedTime.UTC().String()},
				)
				cli.OutputHuman(renderSimpleTable([]string{"RESOURCE ID", "TYPE", "NAME", "DESCRIPTION", "STATE",
					"DEFAULT", "UPDATED BY", "UPDATED_TIME"}, groupCommon))
			}

			return nil
		},
	}

	// delete command is used to remove a lacework resource group by resource id
	resourceGroupsDeleteCommand = &cobra.Command{
		Use:   "delete <resource_group_id>",
		Short: "Delete a resource group",
		Long:  "Delete a single resource group by it's resource group ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			err := cli.LwApi.V2.ResourceGroups.Delete(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to delete resource group")
			}

			cli.OutputHuman("The resource group was deleted.\n")
			return nil
		},
	}

	// create command is used to create a new lacework resource group
	resourceGroupsCreateCommand = &cobra.Command{
		Use:   "create",
		Short: "Create a new resource group",
		Long:  "Creates a new single resource group.",
		RunE: func(_ *cobra.Command, args []string) error {
			if !cli.InteractiveMode() {
				return errors.New("interactive mode is disabled")
			}

			err := promptCreateResourceGroup()
			if err != nil {
				return errors.Wrap(err, "unable to create resource group")
			}

			cli.OutputHuman("The resource group was created.\n")
			return nil
		},
	}
)

// parsePropsType converts props json string to interface of resource group props type
func parsePropsType(props interface{}, resourceType string) (api.ResourceGroupProps, error) {
	propsString := props.(string)

	switch resourceType {
	case api.AwsResourceGroup.String():
		return unmarshallAwsPropString([]byte(propsString))
	case api.AzureResourceGroup.String():
		return unmarshallAzurePropString([]byte(propsString))
	case api.ContainerResourceGroup.String():
		return unmarshallContainerPropString([]byte(propsString))
	case api.GcpResourceGroup.String():
		return unmarshallGcpPropString([]byte(propsString))
	case api.LwAccountResourceGroup.String():
		return unmarshallLwAccountPropString([]byte(propsString))
	case api.MachineResourceGroup.String():
		return unmarshallMachinePropString([]byte(propsString))
	}
	return nil, errors.New("Unable to determine resource group props type")
}

func promptCreateResourceGroup() error {

	resourceGroupOptions := []string{
		"AWS",
		"AZURE",
		"CONTAINER",
		"GCP",
		"LW_ACCOUNT",
		"MACHINE",
	}

	isRGv2Enabled := false
	ffResponse, _ := cli.LwApi.V2.FeatureFlags.GetFeatureFlagsMatchingPrefix(api.ApiV2CliFeatureFlag)
	if len(ffResponse.Data.Flags) >= 1 {
		isRGv2Enabled = true
	}

	if isRGv2Enabled {
		resourceGroupOptions = append(resourceGroupOptions,
			"AWS(v2)",
			"AZURE(v2)",
			"GCP(v2)",
			"CONTAINER(v2)",
			"MACHINE(v2)",
		)
	}

	var (
		group  = ""
		prompt = &survey.Select{
			Message: "Choose a resource group type to create: ",
			Options: resourceGroupOptions,
		}
		err = survey.AskOne(prompt, &group)
	)
	if err != nil {
		return err
	}

	switch group {
	case "AWS":
		return createAwsResourceGroup()
	case "AZURE":
		return createAzureResourceGroup()
	case "CONTAINER":
		return createContainerResourceGroup()
	case "GCP":
		return createGcpResourceGroup()
	case "LW_ACCOUNT":
		return createLwAccountResourceGroup()
	case "MACHINE":
		return createMachineResourceGroup()
	case "AWS(v2)":
		return createResourceGroupV2("AWS")
	case "AZURE(v2)":
		return createResourceGroupV2("AZURE")
	case "GCP(v2)":
		return createResourceGroupV2("GCP")
	case "CONTAINER(v2)":
		return createResourceGroupV2("CONTAINER")
	case "MACHINE(v2)":
		return createResourceGroupV2("MACHINE")
	default:
		return errors.New("unknown resource group type")
	}
}

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

func determineResourceGroupProps(resType string, props api.ResourceGroupProps) [][]string {
	propsString, err := json.Marshal(props)
	if err != nil {
		return [][]string{}
	}
	details := setBaseProps(props)

	switch resType {
	case api.AwsResourceGroup.String():
		details = append(details, setAwsProps(propsString))
	case api.AzureResourceGroup.String():
		details = append(details, setAzureProps(propsString)...)
	case api.ContainerResourceGroup.String():
		details = append(details, setContainerProps(propsString)...)
	case api.GcpResourceGroup.String():
		details = append(details, setGcpProps(propsString)...)
	case api.LwAccountResourceGroup.String():
		details = append(details, setLwAccountProps(propsString))
	case api.MachineResourceGroup.String():
		details = append(details, setMachineProps(propsString))
	}

	return details
}

func setBaseProps(props api.ResourceGroupProps) [][]string {
	var (
		details [][]string
	)
	lastUpdated := props.GetBaseProps().LastUpdated
	details = append(details, []string{"DESCRIPTION", props.GetBaseProps().Description})
	details = append(details, []string{"UPDATED BY", props.GetBaseProps().UpdatedBy})
	details = append(details, []string{"LAST UPDATED", lastUpdated.String()})
	return details
}

func init() {
	// add the resource-group command
	rootCmd.AddCommand(resourceGroupsCommand)

	// add sub-commands to the resource-group command
	resourceGroupsCommand.AddCommand(resourceGroupsListCommand)
	resourceGroupsCommand.AddCommand(resourceGroupsShowCommand)
	resourceGroupsCommand.AddCommand(resourceGroupsCreateCommand)
	resourceGroupsCommand.AddCommand(resourceGroupsDeleteCommand)
}

func IsDefault(isDefault int) string {
	if isDefault == 1 {
		return "True"
	}
	return "False"
}

type resourceGroup struct {
	Id          string                 `json:"resource_guid"`
	ResType     string                 `json:"type"`
	Name        string                 `json:"name"`
	Props       api.ResourceGroupProps `json:"props"`
	Enabled     int                    `json:"enabled"`
	IsDefault   int                    `json:"isDefault"`
	Query       *api.RGQuery           `json:"query"`
	Description string                 `json:"description,omitempty"`
	UpdatedTime *time.Time             `json:"updatedTime,omitempty"`
	UpdatedBy   string                 `json:"updatedBy,omitempty"`
	status      string
}
