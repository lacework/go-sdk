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
	// top-level cloud-account command
	suppressionsCommand = &cobra.Command{
		Use:     "suppressions",
		Hidden:  true,
		Aliases: []string{"suppression", "sup", "sups"},
		Short:   "Manage legacy suppressions",
		Long:    "Manage legacy suppressions",
	}

	// suppressionsListCmd represents the list sub-command inside the suppressions command
	suppressionsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all legacy suppressions per CSP",
	}

	// suppressionsMigrateCmd represents the migrate sub-command inside the suppressions command
	suppressionsMigrateCmd = &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"mig"},
		Short:   "Migrate legacy suppressions for selected CSP, to policy exceptions",
	}
)

func init() {
	rootCmd.AddCommand(suppressionsCommand)
	suppressionsCommand.AddCommand(suppressionsListCmd)
	suppressionsListCmd.AddCommand(suppressionsListAwsCmd)

	suppressionsCommand.AddCommand(suppressionsMigrateCmd)
	suppressionsMigrateCmd.AddCommand(suppressionsMigrateAwsCmd)
}
