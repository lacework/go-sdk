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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/structs"
	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// report-rules command is used to manage lacework report rules
	reportRulesCommand = &cobra.Command{
		Use:     "report-rule",
		Aliases: []string{"report-rules", "rr"},
		Short:   "Manage report rules",
		Long: `Manage report rules to route reports to one or more email alert channels.		

A report rule has four parts:

  1. Email alert channel(s) that should receive the report
  2. One or more severities to include
  3. Resource group(s) containing the subset of your environment to consider
  4. Notification types containing which report information to send
`,
	}

	// list command is used to list all lacework report rules
	reportRulesListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all report rules",
		Long:    "List all report rules configured in your Lacework account.",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.StartProgress(" Fetching report rules...")
			reportRules, err := cli.LwApi.V2.ReportRules.List()
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to get report rules")
			}
			if len(reportRules.Data) == 0 {
				msg := `There are no report rules configured in your account.

Get started by configuring your report rules using the command:

    lacework report-rule create

If you prefer to configure report rules via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Report Rules.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}
			if cli.JSONOutput() {
				return cli.OutputJSON(reportRules)
			}

			var rows [][]string
			for _, rule := range reportRules.Data {
				rows = append(rows, []string{rule.Guid, rule.Filter.Name, rule.Filter.Status()})
			}

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "ENABLED"}, rows))
			return nil
		},
	}
	// show command is used to retrieve a lacework report rule by guid
	reportRulesShowCommand = &cobra.Command{
		Use:   "show <report_rule_id>",
		Short: "Show a report rule by ID",
		Long:  "Show a single report rule by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var response api.ReportRuleResponse
			cli.StartProgress(" Fetching report rule...")

			err := cli.LwApi.V2.ReportRules.Get(args[0], &response)
			if err != nil {
				cli.StopProgress()
				return errors.Wrap(err, "unable to get report rule")
			}
			cli.StopProgress()

			if cli.JSONOutput() {
				return cli.OutputJSON(response)
			}

			reportRule := response.Data
			headers := [][]string{
				[]string{reportRule.Guid, reportRule.Filter.Name, reportRule.Filter.Status()},
			}

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "ENABLED"}, headers))
			cli.OutputHuman("\n")
			cli.OutputHuman(buildReportRuleDetailsTable(reportRule))

			return nil
		},
	}

	// delete command is used to remove a lacework report rule by id
	reportRulesDeleteCommand = &cobra.Command{
		Use:   "delete <report_rule_id>",
		Short: "Delete a report rule",
		Long:  "Delete a single report rule by it's ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cli.StartProgress(" Deleting report rule...")
			err := cli.LwApi.V2.ReportRules.Delete(args[0])
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "unable to delete report rule")
			}
			cli.OutputHuman("The report rule with GUID %s was deleted\n", args[0])
			return nil
		},
	}

	// create command is used to create a new lacework report rule
	reportRulesCreateCommand = &cobra.Command{
		Use:   "create",
		Short: "Create a new report rule",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if !cli.InteractiveMode() {
				return errors.New("interactive mode is disabled")
			}

			response, err := promptCreateReportRule()
			if err != nil {
				return errors.Wrap(err, "unable to create report rule")
			}

			cli.OutputHuman("The report rule was created with GUID %s\n", response.Data.Guid)
			return nil
		},
	}
)

func init() {
	// add the report-rule command
	rootCmd.AddCommand(reportRulesCommand)

	// add sub-commands to the report-rule command
	reportRulesCommand.AddCommand(reportRulesListCommand)
	reportRulesCommand.AddCommand(reportRulesShowCommand)
	reportRulesCommand.AddCommand(reportRulesCreateCommand)
	reportRulesCommand.AddCommand(reportRulesDeleteCommand)
}

func buildReportRuleDetailsTable(rule api.ReportRule) string {
	var (
		details       [][]string
		notifications [][]string
		updatedTime   string
	)
	severities := api.NewReportRuleSeveritiesFromIntSlice(rule.Filter.Severity).ToStringSlice()

	if nano, err := strconv.ParseInt(rule.Filter.CreatedOrUpdatedTime, 10, 64); err == nil {
		updatedTime = time.Unix(nano/1000, 0).Format(time.RFC3339)
	}
	details = append(details, []string{"SEVERITIES", strings.Join(severities, ", ")})
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

	notifcationsMap := rule.ReportNotificationTypes.ToMap()
	// sort keys
	keys := make([]string, 0, len(notifcationsMap))
	for k := range notifcationsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		notifications = append(notifications, []string{key, cases.Title(language.English).String(strconv.FormatBool(notifcationsMap[key]))})
	}

	detailsTable.WriteString(renderCustomTable([]string{"NOTIFICATION TYPES", "ENABLED"}, notifications,
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetColumnSeparator(" ")
			t.SetAutoWrapText(false)
		}),
	),
	)
	detailsTable.WriteString("\n")

	if len(rule.EmailAlertChannels) > 0 {
		channels := [][]string{{strings.Join(rule.EmailAlertChannels, "\n")}}
		detailsTable.WriteString(renderCustomTable([]string{"EMAIL ALERT CHANNELS"}, channels,
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

func promptCreateReportRule() (api.ReportRuleResponse, error) {
	channelList, channelMap := getEmailAlertChannels()
	notificationFields := structs.Names(api.ReportRuleNotificationTypes{})
	notificationsMap := make(map[string]bool)

	if len(channelList) < 1 {
		return api.ReportRuleResponse{}, errors.New("no email alert channels found.")
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
				Message: "Select email alert channels:",
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
			Name: "notifications",
			Prompt: &survey.MultiSelect{
				Message: "Select report notification types:",
				Options: notificationFields,
			},
		},
	}

	answers := struct {
		Name           string
		Description    string   `survey:"description"`
		Channels       []string `survey:"channels"`
		Severities     []string `survey:"severities"`
		ResourceGroups []string `survey:"resourceGroups"`
		Notifications  []string `survey:"notifications"`
	}{}

	err := survey.Ask(questions, &answers,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return api.ReportRuleResponse{}, err
	}

	var channels []string
	for _, channel := range answers.Channels {
		channels = append(channels, channelMap[channel])
	}

	resourceGroups, resourceGroupMap := promptAddResourceGroupsToReportRule()
	var groups []string
	for _, group := range resourceGroups {
		groups = append(groups, resourceGroupMap[group])
	}

	for _, n := range answers.Notifications {
		notificationsMap[n] = true
	}

	notifications := api.ReportRuleNotificationTypes{}
	err = api.TransformReportRuleNotification(notificationsMap, &notifications)
	if err != nil {
		return api.ReportRuleResponse{}, err
	}

	reportRule, err := api.NewReportRule(
		answers.Name,
		api.ReportRuleConfig{
			Description:        answers.Description,
			Severities:         api.NewReportRuleSeverities(answers.Severities),
			ResourceGroups:     groups,
			EmailAlertChannels: channels,
			NotificationTypes:  api.ReportRuleNotifications{notifications},
		})

	if err != nil {
		return api.ReportRuleResponse{}, err
	}

	cli.StartProgress(" Creating report rule...")
	defer cli.StopProgress()

	return cli.LwApi.V2.ReportRules.Create(reportRule)
}

func getEmailAlertChannels() ([]string, map[string]string) {
	cli.StartProgress("")
	defer cli.StopProgress()
	response, err := cli.LwApi.V2.AlertChannels.List()

	if err != nil {
		return nil, nil
	}
	var items = make(map[string]string)
	var channels = make([]string, 0)
	for _, i := range response.Data {
		if i.AlertChannelType() == api.EmailUserAlertChannelType {
			displayName := fmt.Sprintf("%s - %s", i.ID(), i.Name)
			channels = append(channels, displayName)
			items[displayName] = i.ID()
		}
	}

	return channels, items
}

func promptAddResourceGroupsToReportRule() ([]string, map[string]string) {
	addResourceGroups := false
	err := survey.AskOne(&survey.Confirm{
		Message: "Add Resource Groups to Report Rule?",
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
