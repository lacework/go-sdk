//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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
	"strings"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// inventoryAwsListCmd represents the list inventory inside the aws command
	inventoryAwsListCmd = &cobra.Command{
		Use:     "list [resource_type]",
		Aliases: []string{"ls"},
		Short:   "List all Aws inventory resources",
		Long: `List all Aws inventory resources collected in your account.

To list only resources of a specific type use:

    lacework inventory aws list ec2:instance

`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var (
				awsInventorySearchResponse api.InventoryAwsResponse
				cacheKey                   = fmt.Sprintf("inventory/aws")
				now                        = time.Now()
				before                     = now.AddDate(0, 0, -1) // last day
				filter                     = api.InventorySearch{SearchFilter: api.SearchFilter{
					TimeFilter: &api.TimeFilter{
						StartTime: &before,
						EndTime:   &now,
					}},
					Csp: api.AwsInventoryType,
				}
			)

			if len(args) > 0 {
				filter.Filters = []api.Filter{{
					Field:      "resourceType",
					Expression: "eq",
					Value:      args[0],
				}}
				cacheKey = fmt.Sprintf("%s/resource/%s", cacheKey, args[0])
			}

			expired := cli.ReadCachedAsset(cacheKey, &awsInventorySearchResponse)
			if expired {
				cli.StartProgress("Fetching list of Aws inventory resources...")
				err := cli.LwApi.V2.Inventory.Search(&awsInventorySearchResponse, &filter)
				cli.StopProgress()
				if err != nil {
					return errors.Wrap(err, "unable to get Aws inventory resources")
				}
				cli.WriteAssetToCache(cacheKey, time.Now().Add(time.Hour*1), awsInventorySearchResponse)
			}

			cli.OutputHuman(buildAwsInventoryTable(awsInventorySearchResponse))

			return nil
		},
	}
	//// inventoryAwsShowCmd represents the list inventory inside the aws command
	//inventoryAwsShowCmd = &cobra.Command{
	//	Use:     "list-accounts",
	//	Aliases: []string{"list", "ls"},
	//	Short:   "List all Aws inventory resources",
	//	Long:    `List all Aws inventory resources collected in your account.`,
	//	Args:    cobra.NoArgs,
	//	RunE: func(_ *cobra.Command, _ []string) error {
	//		var (
	//			awsInventorySearchResponse api.InventoryAwsResponse
	//			now                        = time.Now()
	//			before                     = now.AddDate(0, 0, -1) // last day
	//			filter                     = api.InventorySearch{SearchFilter: api.SearchFilter{
	//				Filters: []api.Filter{{
	//					Expression: "eq",
	//					Field:      "urn",
	//					Value:      "text",
	//				}},
	//				TimeFilter: &api.TimeFilter{
	//					StartTime: &before,
	//					EndTime:   &now,
	//				}},
	//				Csp: api.AwsInventoryType,
	//			}
	//		)
	//
	//		cli.StartProgress("Fetching list of Aws inventory resources...")
	//		err := cli.LwApi.V2.Inventory.Search(&awsInventorySearchResponse, &filter)
	//		cli.StopProgress()
	//		if err != nil {
	//			return errors.Wrap(err, "unable to get Aws inventory resources")
	//		}
	//
	//		cli.OutputHuman(fmt.Sprint(awsInventorySearchResponse))
	//
	//		return nil
	//	},
	//}
)

func buildAwsInventoryTable(response api.InventoryAwsResponse) string {
	var rows [][]string

	inventoryTable := &strings.Builder{}

	for _, resource := range response.Data {
		rows = append(rows, []string{resource.Urn, resource.ResourceType, resource.CloudDetails.AccountID, resource.EndTime})
	}

	inventoryTable.WriteString(renderOneLineCustomTable("AWS INVENTORY RESOURCES",
		renderSimpleTable([]string{"ARN", "TYPE", "ACCOUNT ID", "COLLECTED TIME"}, rows),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
			t.SetColWidth(50)
		}),
	),
	)

	return inventoryTable.String()
}

func init() {
	// add sub-commands to the aws inventory command
	inventoryAwsCmd.AddCommand(inventoryAwsListCmd)
}
