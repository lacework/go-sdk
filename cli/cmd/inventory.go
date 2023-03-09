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

import "github.com/spf13/cobra"

var (
	// inventoryCmd represents the inventory command
	inventoryCmd = &cobra.Command{
		Use:   "inventory",
		Short: "Manage inventory resources",
		Long: `Manage inventory resources for Google, Azure, or Aws cloud providers.

...
`,
	}

	// inventoryAwsCmd represents the aws sub-command inside the inventory command
	inventoryAwsCmd = &cobra.Command{
		Use:   "aws",
		Short: "Inventory for Aws",
		Long: `Manage inventory resources collected by Lacework for Aws.

To list all inventory resources for Aws:

    lacework inventory aws list
`,
	}
)

func init() {
	// add the inventory command
	rootCmd.AddCommand(inventoryCmd)

	// add sub-commands to the inventory command
	inventoryCmd.AddCommand(inventoryAwsCmd)
}
