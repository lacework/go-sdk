//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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
	"fmt"
	"sort"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	inventoryCmd = &cobra.Command{
		Use:   "inventory",
		Short: "Search and compare cloud resources",
		Long: `
Lacework provides visibility into cloud resources for AWS, Google and Azure clouds. A resource
can be any entity within the cloud deployment, such as an S3 bucket, security group, or Pub/Sub
topics.

For more information about Lacework Resource Inventory, visit:

	https://docs.lacework.net/console/category/resource-inventory
`,
		Hidden: true,
	}
	inventorySearchCmd = &cobra.Command{
		Use:   "search <regex>",
		Short: "Search resources in all cloud providers",
		Args:  cobra.ExactArgs(1),
		RunE:  runInventorySearchCmd,
	}
	inventoryAwsCmd = &cobra.Command{
		Use:   "aws",
		Short: "Inventory for Amazon Web Services",
	}
	inventoryAwsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all AWS resources",
		RunE:  runInventoryAwsListCmd,
	}
	inventoryGoogleCmd = &cobra.Command{
		Use:   "google",
		Short: "Inventory for Google Cloud",
	}
	inventoryAzureCmd = &cobra.Command{
		Use:   "azure",
		Short: "Inventory for Azure Cloud",
	}
)

func init() {
	// add the inventory commands
	rootCmd.AddCommand(inventoryCmd)
	inventoryCmd.AddCommand(inventorySearchCmd)
	inventoryCmd.AddCommand(inventoryAwsCmd)
	inventoryCmd.AddCommand(inventoryGoogleCmd)
	inventoryCmd.AddCommand(inventoryAzureCmd)
	inventoryAwsCmd.AddCommand(inventoryAwsListCmd)
}

func showInventoryResource(inventoryType api.InventoryType, inventory []api.InventoryCommon) error {
	cli.OutputHuman(
		renderSimpleTable(
			[]string{
				"Snapshot Time", "Status", "Cloud", "Resource Region", "Service", "Resource Type", "Resource ID",
			},
			buildInventoryResourcesTableA(inventoryType, inventory),
		),
	)
	return nil
}

func runInventorySearchCmd(_ *cobra.Command, args []string) error {
	var (
		now                     = time.Now().UTC()
		before                  = now.Add(-25 * time.Hour)
		searchInventoryResponse api.InventoryRawResponse
		filters                 = api.InventorySearch{
			SearchFilter: api.SearchFilter{
				TimeFilter: &api.TimeFilter{
					StartTime: &before,
					EndTime:   &now,
				},
				Filters: []api.Filter{{
					Expression: "rlike",
					Field:      "urn",
					Value:      fmt.Sprintf(".*%s.*", args[0]),
				}},
			},
		}
	)

	for _, cloud := range []api.InventoryType{
		api.AwsInventoryType, api.GcpInventoryType, api.AzureInventoryType,
	} {
		cli.StartProgress(fmt.Sprintf("Searching resource on %s...", cloud))
		filters.Csp = cloud
		err := cli.LwApi.V2.Inventory.Search(&searchInventoryResponse, filters)
		cli.StopProgress()
		if err != nil {
			return errors.Wrap(err, "unable to search inventory")
		}

		if len(searchInventoryResponse.Data) != 0 {
			return showInventoryResource(cloud, searchInventoryResponse.Data)
		}
	}

	return errors.New("Resource not found.")
}

func runInventoryAwsListCmd(_ *cobra.Command, args []string) error {
	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -1) // last day
	)

	var (
		awsInventorySearchResponse api.InventoryAwsResponse
		filters                    = api.InventorySearch{
			Csp: api.AwsInventoryType,
			SearchFilter: api.SearchFilter{
				TimeFilter: &api.TimeFilter{
					StartTime: &before,
					EndTime:   &now,
				},
			},
		}
	)

	cli.StartProgress("Fetching AWS inventory...")
	// TODO @afiune search all pages
	err := cli.LwApi.V2.Inventory.Search(&awsInventorySearchResponse, filters)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to execute search inventory")
	}

	// TODO @afiune only display latest resource (de-dup)

	if cli.JSONOutput() {
		return cli.OutputJSON(awsInventorySearchResponse.Data)
	}

	if len(awsInventorySearchResponse.Data) == 0 {
		cli.OutputHuman("There are no AWS resources in your Lacework account.\n")
		return nil
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{
				"Snapshot Time", "Status", "Resource Region", "AWS Account", "Service", "Resource Type", "Resource ID",
			},
			buildInventoryResourcesTableB(awsInventorySearchResponse.Data),
		),
	)
	return nil
}

func buildInventoryResourcesTableA(inventoryType api.InventoryType, inventory []api.InventoryCommon) (out [][]string) {
	for _, resource := range inventory {
		out = append(out, []string{
			resource.StartTime.Format(time.RFC3339),
			resource.Status.Status,
			string(inventoryType),
			resource.ResourceRegion,
			resource.Service,
			resource.ResourceType,
			resource.ResourceID,
		})
	}

	// order by snapshot
	sort.Slice(out, func(i, j int) bool {
		return out[i][0] < out[j][0]
	})

	return
}

func buildInventoryResourcesTableB(inventory []api.InventoryAws) (out [][]string) {
	for _, resource := range inventory {
		out = append(out, []string{
			resource.StartTime,
			resource.Status.Status,
			resource.ResourceRegion,
			resource.CloudDetails.AccountID,
			resource.Service,
			resource.ResourceType,
			resource.ResourceID,
		})
	}

	// order by snapshot
	sort.Slice(out, func(i, j int) bool {
		return out[i][0] < out[j][0]
	})

	return
}
