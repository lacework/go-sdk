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
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
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

				groups = append(groups, resourceGroup{
					Id:               g.ResourceGroupGuid,
					ResType:          g.Type,
					Name:             g.Name,
					Enabled:          g.Enabled,
					IsDefaultBoolean: g.IsDefaultBoolean,
					Query:            g.Query,
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
				rows = append(rows, []string{g.Id, g.ResType, g.Name, strconv.Itoa(g.Enabled),
					strconv.FormatBool(*g.IsDefaultBoolean)})
			}

			cli.OutputHuman(renderSimpleTable([]string{"RESOURCE GROUP ID", "TYPE", "NAME", "STATUS", "DEFAULT"}, rows))
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

			group := resourceGroup{
				Id:               response.Data.ResourceGroupGuid,
				ResType:          response.Data.Type,
				Name:             response.Data.Name,
				Enabled:          response.Data.Enabled,
				IsDefaultBoolean: response.Data.IsDefaultBoolean,
				Query:            response.Data.Query,
				Description:      response.Data.Description,
				UpdatedBy:        response.Data.UpdatedBy,
				UpdatedTime:      response.Data.UpdatedTime,
				CreatedTime:      response.Data.CreatedTime,
				CreatedBy:        response.Data.CreatedBy,
			}

			if cli.JSONOutput() {
				jsonOut := struct {
					Group resourceGroup `json:"resource_group"`
				}{Group: group}
				return cli.OutputJSON(jsonOut)
			}

			var groupCommon [][]string

			groupCommon = append(groupCommon,
				[]string{group.Id, group.ResType, group.Name, group.Description, strconv.Itoa(group.Enabled),
					strconv.FormatBool(*group.IsDefaultBoolean), group.CreatedBy, group.CreatedTime.UTC().String(),
					group.UpdatedBy, group.UpdatedTime.UTC().String()},
			)
			cli.OutputHuman(renderSimpleTable([]string{"RESOURCE GROUP ID", "TYPE", "NAME", "DESCRIPTION", "STATE",
				"DEFAULT", "CREATED BY", "CREATED TIME", "UPDATED BY", "UPDATED TIME"}, groupCommon))

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

func promptCreateResourceGroup() error {

	resourceGroupOptions := []string{
		"AWS",
		"AZURE",
		"CONTAINER",
		"GCP",
		"MACHINE",
		"OCI",
		"KUBERNETES",
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
		return createResourceGroup("AWS")
	case "AZURE":
		return createResourceGroup("AZURE")
	case "GCP":
		return createResourceGroup("GCP")
	case "CONTAINER":
		return createResourceGroup("CONTAINER")
	case "MACHINE":
		return createResourceGroup("MACHINE")
	case "OCI":
		return createResourceGroup("OCI")
	case "KUBERNETES":
		return createResourceGroup("KUBERNETES")
	default:
		return errors.New("unknown resource group type")
	}
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

type resourceGroup struct {
	Id               string       `json:"resourceGroupGuid"`
	ResType          string       `json:"type"`
	Name             string       `json:"name"`
	Enabled          int          `json:"enabled"`
	IsDefaultBoolean *bool        `json:"isDefaultBoolean"`
	Query            *api.RGQuery `json:"query"`
	Description      string       `json:"description,omitempty"`
	UpdatedTime      *time.Time   `json:"updatedTime,omitempty"`
	UpdatedBy        string       `json:"updatedBy,omitempty"`
	CreatedBy        string       `json:"createdBy,omitempty"`
	CreatedTime      *time.Time   `json:"createdTime,omitempty"`
}
