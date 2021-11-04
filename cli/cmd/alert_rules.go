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
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// alert-rules command is used to manage lacework alert rules
	alertRulesCommand = &cobra.Command{
		Use:     "alert-rule",
		Aliases: []string{"alert-rules", "ar"},
		Short:   "manage alert rules",
		Long: `Manage alert rules to route events to the appropriate people or tools.		
An alert rule has three parts:
  1. Alert channel(s) that should receive the event notification
  2. Event severity and categories to include
  3. Resource group(s) containing the subset of your environment to consider
`,
	}

	// list command is used to list all lacework alert rules
	alertRulesListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all alert rules",
		Long:    "List all alert rules configured in your Lacework account.",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			alertRules, err := cli.LwApi.V2.AlertRules.List()
			if err != nil {
				return errors.Wrap(err, "unable to get alert rules")
			}
			if len(alertRules.Data) == 0 {
				msg := `There are no alert rules configured in your account.

Get started by integrating your alert rules to manage alerting using the command:

    lacework alert-rule create

If you prefer to configure alert rules via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Alert Rules.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}
			if cli.JSONOutput() {
				return cli.OutputJSON(alertRules)
			}

			var rows [][]string
			for _, rule := range alertRules.Data {
				rows = append(rows, []string{rule.Guid, rule.Filter.Name, rule.Filter.Status()})
			}

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "ENABLED"}, rows))
			return nil
		},
	}
	// show command is used to retrieve a lacework alert rule by resource id
	alertRulesShowCommand = &cobra.Command{
		Use:   "show",
		Short: "show an alert rule by id",
		Long:  "Show a single alert rule by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var response api.AlertRuleResponse
			err := cli.LwApi.V2.AlertRules.Get(args[0], &response)
			if err != nil {
				return errors.Wrap(err, "unable to get alert rule")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			alertRule := response.Data
			var headers [][]string
			headers = append(headers, []string{alertRule.Guid, alertRule.Filter.Name, alertRule.Filter.Status()})

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "ENABLED"}, headers))
			cli.OutputHuman("\n")
			cli.OutputHuman(buildAlertRuleDetailsTable(alertRule))

			return nil
		},
	}

	// delete command is used to remove a lacework alert rule by resource id
	alertRulesDeleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "delete a alert rule",
		Long:  "Delete a single alert rule by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			err := cli.LwApi.V2.AlertRules.Delete(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to delete alert rule")
			}
			cli.OutputHuman(fmt.Sprintf("The alert rule with GUID %s was deleted \n", args[0]))
			return nil
		},
	}

	// create command is used to create a new lacework alert rule
	alertRulesCreateCommand = &cobra.Command{
		Use:   "create",
		Short: "create a new alert rule",
		Long:  "Creates a new single alert rule.",
		RunE: func(_ *cobra.Command, args []string) error {
			if !cli.InteractiveMode() {
				return errors.New("interactive mode is disabled")
			}

			response, err := promptCreateAlertRule()
			if err != nil {
				return errors.Wrap(err, "unable to create alert rule")
			}

			cli.OutputHuman(fmt.Sprintf("The alert rule was created with GUID %s \n", response.Data.Guid))
			return nil
		},
	}
)

func init() {
	// add the alert-rule command
	rootCmd.AddCommand(alertRulesCommand)

	// add sub-commands to the alert-rule command
	alertRulesCommand.AddCommand(alertRulesListCommand)
	alertRulesCommand.AddCommand(alertRulesShowCommand)
	alertRulesCommand.AddCommand(alertRulesCreateCommand)
	alertRulesCommand.AddCommand(alertRulesDeleteCommand)
}

func buildAlertRuleDetailsTable(rule api.AlertRule) string {
	var (
		details     [][]string
		updatedTime string
	)
	severities := api.NewAlertRuleSeveritiesFromIntSlice(rule.Filter.Severity).ToStringSlice()

	if nano, err := strconv.ParseInt(rule.Filter.CreatedOrUpdatedTime, 10, 64); err == nil {
		updatedTime = time.Unix(nano/1000, 0).Format(time.RFC3339)
	}
	details = append(details, []string{"SEVERITIES", strings.Join(severities, ", ")})
	details = append(details, []string{"EVENT CATEGORIES", strings.Join(rule.Filter.EventCategories, ", ")})
	details = append(details, []string{"DESCRIPTION", rule.Filter.Description})
	details = append(details, []string{"UPDATED BY", rule.Filter.CreatedOrUpdatedBy})
	details = append(details, []string{"LAST UPDATED", updatedTime})

	detailsTable := &strings.Builder{}
	detailsTable.WriteString(renderOneLineCustomTable("ALERT RULE DETAILS",
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

	if len(rule.Channels) > 0 {
		channels := [][]string{{strings.Join(rule.Channels, "\n")}}
		detailsTable.WriteString(renderCustomTable([]string{"ALERT CHANNELS"}, channels,
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
			}),
		),
		)
		detailsTable.WriteString("\n")
	}

	if len(rule.Filter.ResourceGroups) > 0 {
		resourceGroups := [][]string{{strings.Join(rule.Filter.ResourceGroups, "\n")}}
		detailsTable.WriteString(renderCustomTable([]string{"RESOURCE GROUPS"}, resourceGroups,
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
			}),
		),
		)
	}

	return detailsTable.String()
}

func promptCreateAlertRule() (api.AlertRuleResponse, error) {
	channelList, channelMap := getAlertChannels()

	if len(channelList) < 1 {
		return api.AlertRuleResponse{}, errors.New("no Alert Channels found.")
	}

	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name: "},
			Validate: survey.Required,
		},
		{
			Name:     "description",
			Prompt:   &survey.Input{Message: "Description: "},
			Validate: survey.Required,
		},
		{
			Name: "channels",
			Prompt: &survey.MultiSelect{
				Message: "Select alert channels:",
				Options: channelList,
			},
			Validate: survey.Required,
		},
		{
			Name: "severities",
			Prompt: &survey.MultiSelect{
				Message: "Select severities:",
				Options: []string{"Critical", "High", "Medium", "Low", "Info"},
			},
		},
		{
			Name: "eventCategories",
			Prompt: &survey.MultiSelect{
				Message: "Select event categories:",
				Options: []string{"Compliance", "App", "Cloud", "File", "Machine", "User", "Platform"},
			},
		},
	}

	answers := struct {
		Name            string
		Description     string   `survey:"description"`
		Channels        []string `survey:"channels"`
		Severities      []string `survey:"severities"`
		EventCategories []string `survey:"eventCategories"`
		ResourceGroups  []string `survey:"resourceGroups"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return api.AlertRuleResponse{}, err
	}

	var channels []string
	for _, channel := range answers.Channels {
		channels = append(channels, channelMap[channel])
	}

	resourceGroups, resourceGroupMap := promptAddResourceGroupsToAlertRule()
	var groups []string
	for _, group := range resourceGroups {
		groups = append(groups, resourceGroupMap[group])
	}

	alertRule := api.NewAlertRule(
		answers.Name,
		api.AlertRuleConfig{
			Description:     answers.Description,
			Channels:        channels,
			Severities:      api.NewAlertRuleSeverities(answers.Severities),
			EventCategories: answers.EventCategories,
			ResourceGroups:  groups,
		})

	cli.StartProgress(" Creating alert rule...")
	response, err := cli.LwApi.V2.AlertRules.Create(alertRule)

	cli.StopProgress()
	return response, err
}

func getAlertChannels() ([]string, map[string]string) {
	response, err := cli.LwApi.V2.AlertChannels.List()

	if err != nil {
		return nil, nil
	}
	var items = make(map[string]string)
	var channels = make([]string, 0)
	for _, i := range response.Data {
		displayName := fmt.Sprintf("%s - %s", i.ID(), i.Name)
		channels = append(channels, displayName)
		items[displayName] = i.ID()
	}

	return channels, items
}

func getResourceGroups() ([]string, map[string]string) {
	response, err := cli.LwApi.V2.ResourceGroups.List()

	if err != nil {
		return nil, nil
	}
	var items = make(map[string]string)
	var groups = make([]string, 0)

	for _, i := range response.Data {
		displayName := fmt.Sprintf("%s - %s", i.ID(), i.Name)
		groups = append(groups, displayName)
		items[displayName] = i.ID()
	}

	return groups, items
}

func promptAddResourceGroupsToAlertRule() ([]string, map[string]string) {
	addResourceGroups := false
	err := survey.AskOne(&survey.Confirm{
		Message: "Add Resource Groups to Alert Rule?",
	}, &addResourceGroups)

	if err != nil {
		return nil, nil
	}

	if addResourceGroups {
		var groups []string
		groupList, groupMap := getResourceGroups()

		err = survey.AskOne(&survey.MultiSelect{
			Message: "Select Resource Groups:",
			Options: groupList,
		}, &groups)

		if err != nil {
			return nil, nil
		}
		return groups, groupMap
	}
	return nil, nil
}
