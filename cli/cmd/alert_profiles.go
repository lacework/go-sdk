//
// Author:: Darren Murray(<darren.murray@lacework.net>)
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
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// alert-profiles command is used to manage lacework alert profiles
	alertProfilesCommand = &cobra.Command{
		Use:     "alert-profile",
		Aliases: []string{"alert-profiles", "ap"},
		Short:   "Manage alert profiles",
		Long: `Manage alert profiles to define how your LQL queries get consumed into alerts.		

An alert profile consists of the name of the new profile, the name of existing alert profile 
from which this new profile extends and a list of alert templates.
`,
	}

	// list command is used to list all lacework alert profiles
	alertProfilesListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all alert profiles",
		Long:    "List all alert profiles configured in your Lacework account.",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			alertProfiles, err := cli.LwApi.V2.Alert.Profiles.List()
			if err != nil {
				return errors.Wrap(err, "unable to get alert profiles")
			}
			if len(alertProfiles.Data) == 0 {
				msg := `There are no alert profiles configured in your account.

Get started by integrating your alert profiles to manage alerting using the command:

    lacework alert-profile create

If you prefer to configure alert profiles via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Alert Profiles.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}
			if cli.JSONOutput() {
				return cli.OutputJSON(alertProfiles)
			}

			var rows [][]string
			for _, profile := range alertProfiles.Data {
				rows = append(rows, []string{profile.Guid, profile.Extends})
			}

			cli.OutputHuman(renderSimpleTable([]string{"ID", "EXTENDS"}, rows))
			return nil
		},
	}
	// show command is used to retrieve a lacework alert profile by id
	alertProfilesShowCommand = &cobra.Command{
		Use:     "show <alert_profile_id>",
		Short:   "Show an alert profile by ID",
		Aliases: []string{"get"},
		Long:    "Show a single alert profile by it's ID.",
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var response api.AlertProfileResponse
			err := cli.LwApi.V2.Alert.Profiles.Get(args[0], &response)
			if err != nil {
				return errors.Wrap(err, "unable to get alert profile")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			alertProfile := response.Data
			var headers [][]string
			headers = append(headers, []string{alertProfile.Guid, alertProfile.Extends})
			cli.OutputHuman(renderSimpleTable([]string{"ID", "EXTENDS"}, headers))
			cli.OutputHuman("\n")
			cli.OutputHuman(buildAlertProfileDetailsTable(alertProfile))

			return nil
		},
	}

	// delete command is used to remove a lacework alert profile by resource id
	alertProfilesDeleteCommand = &cobra.Command{
		Use:   "delete <alert_profile_id>",
		Short: "Delete a alert profile",
		Long:  "Delete a single alert profile by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			err := cli.LwApi.V2.Alert.Profiles.Delete(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to delete alert profile")
			}
			cli.OutputHuman(fmt.Sprintf("The alert profile with GUID %s was deleted \n", args[0]))
			return nil
		},
	}

	// create command is used to create a new lacework alert profile
	alertProfilesCreateCommand = &cobra.Command{
		Use:   "create",
		Short: "Create a new alert profile",
		RunE: func(_ *cobra.Command, args []string) error {
			if !cli.InteractiveMode() {
				return errors.New("interactive mode is disabled")
			}

			response, err := promptCreateAlertProfile()
			if err != nil {
				return errors.Wrap(err, "unable to create alert profile")
			}

			cli.OutputHuman(fmt.Sprintf("The alert profile was created with GUID %s \n", response.Data.Guid))
			return nil
		},
	}
)

func init() {
	// add the alert-profile command
	rootCmd.AddCommand(alertProfilesCommand)

	// add sub-commands to the alert-profile command
	alertProfilesCommand.AddCommand(alertProfilesListCommand)
	alertProfilesCommand.AddCommand(alertProfilesShowCommand)
	alertProfilesCommand.AddCommand(alertProfilesCreateCommand)
	alertProfilesCommand.AddCommand(alertProfilesDeleteCommand)
}

func buildAlertProfileDetailsTable(profile api.AlertProfile) string {
	var details [][]string

	detailsTable := &strings.Builder{}

	for _, alert := range profile.Alerts {
		details = append(details, []string{alert.Name, alert.EventName, alert.Description, alert.Subject})
	}

	detailsTable.WriteString(renderOneLineCustomTable("ALERT TEMPLATES",
		renderSimpleTable([]string{"NAME", "EVENT NAME", "DESCRIPTION", "SUBJECT"}, details),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	),
	)

	if len(profile.DescriptionKeys) > 0 {
		var keys [][]string
		for _, key := range profile.DescriptionKeys {
			keys = append(keys, []string{key.Name, key.Spec})
		}
		detailsTable.WriteString(renderOneLineCustomTable("DESCRIPTION KEYS",
			renderSimpleTable([]string{"NAME", "SPEC"}, keys),
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetAutoWrapText(false)
			}),
		),
		)
		detailsTable.WriteString("\n")
	}

	if len(profile.Fields) > 0 {
		var fields [][]string
		for _, f := range profile.Fields {
			fields = append(fields, []string{f.Name})
		}
		detailsTable.WriteString(renderOneLineCustomTable("FIELDS",
			renderSimpleTable([]string{"NAME"}, fields),
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetAutoWrapText(false)
			}),
		),
		)
		detailsTable.WriteString("\n")
	}

	return detailsTable.String()
}

func promptCreateAlertProfile() (api.AlertProfileResponse, error) {
	profileList := getAlertProfiles()

	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Profile Name: "},
			Validate: survey.Required,
		},
		{
			Name: "extends",
			Prompt: &survey.Select{
				Message: "Select an alert profile to extend from:",
				Options: profileList,
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name    string `json:"name"`
		Extends string `json:"extends"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return api.AlertProfileResponse{}, err
	}

	if strings.HasPrefix(answers.Name, "LW_") {
		return api.AlertProfileResponse{}, errors.New("profile name prefix 'LW_' is reserved for Lacework defined profiles")
	}

	var templates []api.AlertTemplate
	templates = append(templates, promptAddAlertTemplate())
	addTemplates := false
	for {
		if err := survey.AskOne(&survey.Confirm{
			Message: "Add another alert template?",
		}, &addTemplates); err != nil {
			return api.AlertProfileResponse{}, err
		}

		if addTemplates {
			templates = append(templates, promptAddAlertTemplate())
		} else {
			break
		}
	}
	alertProfile := api.NewAlertProfile(answers.Name, answers.Extends, templates)

	cli.StartProgress(" Creating alert profile...")
	response, err := cli.LwApi.V2.Alert.Profiles.Create(alertProfile)

	cli.StopProgress()
	return response, err
}

func promptAddAlertTemplate() api.AlertTemplate {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Alert Template Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "eventName",
			Prompt:   &survey.Input{Message: "Alert Template Event Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "description",
			Prompt:   &survey.Input{Message: "Alert Template Description: "},
			Validate: survey.Required,
		},
		{
			Name:     "subject",
			Prompt:   &survey.Input{Message: "Alert Template Subject: "},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name        string `json:"name"`
		EventName   string `json:"eventName"`
		Description string `json:"description"`
		Subject     string `json:"subject"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return api.AlertTemplate{}
	}

	return api.AlertTemplate{
		Name:        answers.Name,
		EventName:   answers.EventName,
		Description: answers.Description,
		Subject:     answers.Subject,
	}
}

func getAlertProfiles() []string {
	response, err := cli.LwApi.V2.Alert.Profiles.List()

	if err != nil {
		return nil
	}
	var profiles = make([]string, 0)
	for _, p := range response.Data {
		// profiles can only extend from 'LW_' profiles with >= 1 alert template
		if strings.HasPrefix(p.Guid, "LW_") && len(p.Alerts) >= 1 {
			profiles = append(profiles, p.Guid)
		}
	}

	return profiles
}
