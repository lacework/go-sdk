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
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// team-members command is used to manage lacework team members
	teamMembersCommand = &cobra.Command{
		Use:     "team-member",
		Aliases: []string{"team-members", "tm"},
		Short:   "Manage team members",
		Long: `Manage Team Members to grant or restrict access to multiple Lacework Accounts. 
			  Team members can also be granted organization-level roles.
`,
	}

	// list command is used to list all lacework team members
	teamMembersListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all team members",
		Long:    "List all team members configured in your Lacework account.",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.StartProgress(" Fetching team members...")
			teamMembers, err := cli.LwApi.V2.TeamMembers.List()
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get team members")
			}
			if len(teamMembers.Data) == 0 {
				msg := `There are no team members configured in your account.

Get started by configuring your team members using the command:

    lacework team-member create

If you prefer to configure team members via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Team Members.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}
			if cli.JSONOutput() {
				return cli.OutputJSON(teamMembers)
			}

			var rows [][]string
			for _, tm := range teamMembers.Data {
				rows = append(rows, []string{tm.UserGuid, tm.UserName, enabled(tm.UserEnabled)})
			}

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "STATUS"}, rows))
			return nil
		},
	}
	// show command is used to retrieve a lacework team member by guid
	teamMembersShowCommand = &cobra.Command{
		Use:   "show <team_member_id>",
		Short: "Show a team member by id",
		Long:  "Show a single team member by it's id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var response api.TeamMemberResponse
			cli.StartProgress(" Fetching team member...")

			err := cli.LwApi.V2.TeamMembers.Get(args[0], &response)
			if err != nil {
				cli.StopProgress()
				return errors.Wrap(err, "unable to get team member")
			}
			cli.StopProgress()

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			teamMember := response.Data
			headers := [][]string{{teamMember.UserGuid, teamMember.UserName, enabled(teamMember.UserEnabled)}}
			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "STATUS"}, headers))
			cli.OutputHuman("\n")
			cli.OutputHuman(buildTeamMemberDetailsTable(teamMember))

			return nil
		},
	}

	// delete command is used to remove a lacework team member by id
	teamMembersDeleteCommand = &cobra.Command{
		Use:   "delete <team_member_id>",
		Short: "Delete a team member",
		Long:  "Delete a single team member by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cli.StartProgress(" Deleting team member...")
			err := cli.LwApi.V2.TeamMembers.Delete(args[0])
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to delete team member")
			}
			cli.OutputHuman("The team member with GUID %s was deleted\n", args[0])
			return nil
		},
	}

	// create command is used to create a new lacework team member
	teamMembersCreateCommand = &cobra.Command{
		Use:   "create",
		Short: "Create a new team member",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if !cli.InteractiveMode() {
				return errors.New("interactive mode is disabled")
			}

			userID, err := promptCreateTeamMember()
			if err != nil {
				return errors.Wrap(err, "unable to create team member")
			}

			cli.OutputHuman("The team member was created with GUID %s\n", userID)
			return nil
		},
	}
)

func init() {
	// add the team-member command
	rootCmd.AddCommand(teamMembersCommand)

	// add sub-commands to the team-member command
	teamMembersCommand.AddCommand(teamMembersListCommand)
	teamMembersCommand.AddCommand(teamMembersShowCommand)
	teamMembersCommand.AddCommand(teamMembersCreateCommand)
	teamMembersCommand.AddCommand(teamMembersDeleteCommand)
}

func buildTeamMemberDetailsTable(tm api.TeamMember) string {
	var (
		details     [][]string
		updatedTime string
	)

	if tm.Props.UpdatedTime != nil {
		updatedTime = fmt.Sprintf("%v", tm.Props.UpdatedTime)
	}

	details = append(details, []string{"FIRST NAME", tm.Props.FirstName})
	details = append(details, []string{"LAST NAME", tm.Props.LastName})
	details = append(details, []string{"COMPANY", tm.Props.Company})
	details = append(details, []string{"ACCOUNT ADMIN", strconv.FormatBool(tm.Props.AccountAdmin)})
	details = append(details, []string{"ORG ADMIN", strconv.FormatBool(tm.Props.OrgAdmin)})
	details = append(details, []string{"CREATED AT", tm.Props.CreatedTime})
	details = append(details, []string{"JIT CREATED", strconv.FormatBool(tm.Props.JitCreated)})
	details = append(details, []string{"UPDATED BY", tm.Props.UpdatedBy})
	details = append(details, []string{"UPDATED AT", updatedTime})

	detailsTable := &strings.Builder{}
	detailsTable.WriteString(renderOneLineCustomTable("TEAM MEMBER DETAILS",
		renderCustomTable([]string{}, details,
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
	),
	)

	detailsTable.WriteString("\n")
	return detailsTable.String()
}

func promptCreateTeamMember() (string, error) {
	questions := []*survey.Question{
		{
			Name:     "email",
			Prompt:   &survey.Input{Message: "Email: "},
			Validate: validateEmail(),
		},
		{
			Name:     "firstName",
			Prompt:   &survey.Input{Message: "First Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "lastName",
			Prompt:   &survey.Input{Message: "Last Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "company",
			Prompt:   &survey.Input{Message: "Company: "},
			Validate: survey.Required,
		},
		{
			Name:     "orgTeamMemberPrompt",
			Prompt:   &survey.Confirm{Message: "Create at Organization Level?"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Email               string `survey:"email"`
		Description         string `survey:"description"`
		Company             string `survey:"company"`
		FirstName           string `survey:"firstName"`
		LastName            string `survey:"lastName"`
		OrgTeamMemberPrompt bool   `survey:"orgTeamMemberPrompt"`
		OrgAdminRole        bool
		UserRole            bool
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return "", err
	}

	if answers.OrgTeamMemberPrompt {
		orgConfig, err := askConfigureOrgTeamMemberPrompt()
		if err != nil {
			return "", err
		}

		var (
			adminRoleAccounts []string
			userRoleAccounts  []string
		)

		switch orgConfig {
		case "User Role":
			answers.UserRole = true
		case "Admin Role":
			answers.OrgAdminRole = true
		case "Manage Roles for Accounts":
			adminRoleAccounts, userRoleAccounts, err = askManageOrgTeamMemberRolesPrompt()
			if err != nil {
				return "", err
			}
		}

		teamMember := api.NewTeamMemberOrg(answers.Email, api.TeamMemberProps{
			Company:   answers.Company,
			FirstName: answers.FirstName,
			LastName:  answers.LastName,
			OrgAdmin:  answers.OrgAdminRole,
			OrgUser:   answers.UserRole,
		})

		if len(adminRoleAccounts) > 0 {
			teamMember.AdminRoleAccounts = sliceToUpper(adminRoleAccounts)
		}
		if len(userRoleAccounts) > 0 {
			teamMember.UserRoleAccounts = sliceToUpper(userRoleAccounts)
		}

		cli.StartProgress(" Creating team member...")
		defer cli.StopProgress()

		res, err := cli.LwApi.V2.TeamMembers.CreateOrg(teamMember)
		if err != nil {
			return "", err
		}
		userID := res.Data.Accounts[0].UserGuid
		return userID, nil
	}
	var accountAdmin bool
	err = survey.AskOne(&survey.Confirm{Message: "Account Admin?"}, &accountAdmin)
	if err != nil {
		return "", err
	}

	teamMember := api.NewTeamMember(answers.Email, api.TeamMemberProps{
		Company:      answers.Company,
		FirstName:    answers.FirstName,
		LastName:     answers.LastName,
		AccountAdmin: accountAdmin,
	})

	cli.StartProgress(" Creating team member...")
	defer cli.StopProgress()

	res, err := cli.LwApi.V2.TeamMembers.Create(teamMember)
	if err != nil {
		return "", err
	}
	userID := res.Data.UserGuid
	return userID, nil
}

func askManageOrgTeamMemberRolesPrompt() ([]string, []string, error) {
	res, err := cli.LwApi.V2.UserProfile.Get()
	if err != nil {
		return nil, nil, err
	}
	accountsList := res.Data[0].SubAccountNames()

	accounts := struct {
		UserRoleAccounts  []string `survey:"userRoleAccounts"`
		AdminRoleAccounts []string `survey:"adminRoleAccounts"`
	}{}

	questions := []*survey.Question{
		{
			Name: "userRoleAccounts",
			Prompt: &survey.MultiSelect{
				Message: "Select Accounts Team Member will have User role: ",
				Options: accountsList},
		},
		{
			Name: "adminRoleAccounts",
			Prompt: &survey.MultiSelect{
				Message: "Select Accounts Team Member will have Admin role: ",
				Options: accountsList},
		},
	}

	err = survey.Ask(questions, &accounts,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return nil, nil, err
	}
	return accounts.AdminRoleAccounts, accounts.UserRoleAccounts, nil
}

func askConfigureOrgTeamMemberPrompt() (string, error) {
	var configureOrgMember string

	err := survey.AskOne(&survey.Select{
		Message: "Select a role for all accounts: ",
		Options: []string{"User Role", "Admin Role", "Manage Roles for Accounts"}},
		&configureOrgMember,
		survey.WithIcons(promptIconsFunc))

	if err != nil {
		return "", err
	}
	return configureOrgMember, nil
}

func enabled(status int) string {
	if status == 1 {
		return "Enabled"
	}
	return "Disabled"
}

func sliceToUpper(list []string) (upper []string) {
	for _, item := range list {
		upper = append(upper, strings.ToUpper(item))
	}
	return
}

func validateEmail() survey.Validator {
	return func(val interface{}) error {
		emailRegex, _ := regexp.Compile(
			"[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?", //nolint
		)
		if !emailRegex.MatchString(val.(string)) {
			return fmt.Errorf("not a valid email %s", val.(string))
		}
		return nil
	}
}
