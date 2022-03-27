//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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
	"github.com/spf13/cobra"
)

var (
	// accountCmd represents the account command
	accountCmd = &cobra.Command{
		Use:     "account",
		Aliases: []string{"accounts", "acc"},
		Short:   "Manage accounts in an organization (org admins only)",
		Long: `Manage accounts inside your Lacework organization.

An organization can contain multiple accounts so you can also manage components
such as alerts, resource groups, team members, and audit logs at a more granular
level inside an organization. A team member may have access to multiple accounts
and can easily switch between them.

To enroll your Lacework account in an organization follow the documentation:

  https://docs.lacework.com/organization-overview
    `,
	}

	// accountListCmd represents the list command inside the account command
	accountListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all accounts",
		Long:    `List all accounts in your organization.`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.StartProgress(" Loading account information ...")
			user, err := cli.LwApi.V2.UserProfile.Get()
			cli.StopProgress()
			if err != nil {
				return err
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(user.Data)
			}

			if len(user.Data) == 0 {
				return yikes("unable to load account information.")
			}

			profile := user.Data[0]
			if !profile.OrgAccount {
				cli.OutputHuman("Your account is not enrolled in an organization.\n")
				return nil
			}

			rows := [][]string{{profile.OrgAccountName()}}
			for _, acc := range profile.SubAccountNames() {
				rows = append(rows, []string{acc})
			}

			cli.OutputHuman(renderSimpleTable([]string{"Accounts"}, rows))
			cli.OutputHuman("\nUse '--subaccount <name>' to switch any command to a different account.\n")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountListCmd)
}
