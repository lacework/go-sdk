//
// Author:: Ross Moles (<ross.moles@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"github.com/spf13/cobra"
)

var (
	// top-level suppressions command
	suppressionsCommand = &cobra.Command{
		Use:     "suppressions",
		Hidden:  true,
		Aliases: []string{"suppression", "sup", "sups"},
		Short:   "Manage legacy suppressions",
		Long:    "Manage legacy suppressions",
	}

	// suppressionsAwsCmd represents the aws sub-command inside the suppressions command
	suppressionsAwsCmd = &cobra.Command{
		Use:   "aws",
		Short: "Manage legacy suppressions for aws",
	}
)

func init() {
	rootCmd.AddCommand(suppressionsCommand)
	suppressionsCommand.AddCommand(suppressionsAwsCmd)
	suppressionsAwsCmd.AddCommand(suppressionsListAwsCmd)
	suppressionsAwsCmd.AddCommand(suppressionsMigrateAwsCmd)
}
