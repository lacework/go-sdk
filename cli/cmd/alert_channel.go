package cmd

import (
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	alertChannelCommand = &cobra.Command{
		Use:     "alert-channel",
		Aliases: []string{"alert-channels", "ac"},
		Short:   "Manage alert channels",
		Long:    "Manage alert channels integrations with Lacework",
	}

	// alertChannelsListCmd represents the list sub-command inside the alert channels command
	alertChannelListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available alert channel integrations",
		Args:    cobra.NoArgs,
		RunE:    alertChannelList,
	}

	// alertChannelShowCmd represents the show sub-command inside the alert channel command
	alertChannelShowCmd = &cobra.Command{
		Use:     "show",
		Aliases: []string{"get"},
		Short:   "Show a single alert channel integration",
		Args:    cobra.ExactArgs(1),
		RunE:    alertChannelShow,
	}

	// alertChannelCreateCmd represents the show sub-command inside the alert channels command
	alertChannelCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new alert channel integration",
		Args:  cobra.NoArgs,
		RunE:  alertChannelCreate,
	}

	// alertChannelDeleteCmd represents the delete sub-command inside the alert channels command
	alertChannelDeleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "Delete a alert channel integration",
		Args:    cobra.ExactArgs(1),
		RunE:    alertChannelDelete,
	}
)

func init() {
	// add the alert-channel command
	rootCmd.AddCommand(alertChannelCommand)
	alertChannelCommand.AddCommand(alertChannelListCmd)
	alertChannelCommand.AddCommand(alertChannelShowCmd)
	alertChannelCommand.AddCommand(alertChannelCreateCmd)
	alertChannelCommand.AddCommand(alertChannelDeleteCmd)
}

func alertChannelsToTable(alertChannels []api.AlertChannelRaw) [][]string {
	var out [][]string
	for _, cadata := range alertChannels {
		out = append(out, []string{
			cadata.IntgGuid,
			cadata.Name,
			cadata.Type,
			cadata.Status(),
			cadata.StateString(),
		})
	}
	return out
}

func alertChannelList(_ *cobra.Command, _ []string) error {
	alertChannels, err := cli.LwApi.V2.AlertChannels.List()

	if err != nil {
		return errors.Wrap(err, "unable to get alert channels")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(alertChannels.Data)
	}

	if len(alertChannels.Data) == 0 {
		cli.OutputHuman("No alert channels found.\n")
		return nil
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"alert channel GUID", "Name", "Type", "Status", "State"},
			alertChannelsToTable(alertChannels.Data),
		),
	)
	return nil
}

func alertChannelShow(_ *cobra.Command, args []string) error {
	var (
		alertChannel api.AlertChannelResponse
		out          [][]string
	)
	cli.StartProgress("Fetching alert channel...")
	time.Sleep(time.Second * 3)
	err := cli.LwApi.V2.AlertChannels.Get(args[0], &alertChannel)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve alert channel")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(alertChannel.Data)
	}

	out = append(out, []string{alertChannel.Data.IntgGuid,
		alertChannel.Data.Name,
		alertChannel.Data.Type,
		alertChannel.Data.Status(),
		alertChannel.Data.StateString()})

	cli.OutputHuman(renderSimpleTable([]string{"Alert Channel GUID", "Name", "Type", "Status", "State"}, out))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildDetailsTable(alertChannel.Data))
	return nil
}

func alertChannelCreate(_ *cobra.Command, _ []string) error {
	if !cli.InteractiveMode() {
		return errors.New("interactive mode is disabled")
	}

	err := promptCreateAlertChannel()
	if err != nil {
		return errors.Wrap(err, "unable to create alert channel")
	}

	cli.OutputHuman("The alert channel was created.\n")
	return nil
}

func alertChannelDelete(_ *cobra.Command, args []string) error {
	cli.StartProgress(" Deleting alert channel...")
	err := cli.LwApi.V2.AlertChannels.Delete(args[0])
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to delete alert channel")
	}
	cli.OutputHuman("The alert channel %s was deleted.\n", args[0])
	return nil
}

func promptCreateAlertChannel() error {
	var (
		alertChannel = ""
		prompt       = &survey.Select{
			Message: "Choose an alert channel type to create: ",
			Options: []string{
				"Slack",
				"Email",
				"Amazon S3",
				"Cisco Webex",
				"Datadog",
				"GCP PubSub",
				"Microsoft Teams",
				"New Relic Insights",
				"Webhook",
				"VictorOps",
				"Splunk",
				"QRadar",
				"Service Now",
				"PagerDuty",
				"Amazon CloudWatch",
				"Jira Cloud",
				"Jira Server",
			},
		}
		err = survey.AskOne(prompt, &alertChannel)
	)
	if err != nil {
		return err
	}

	switch alertChannel {
	case "Slack":
		return createSlackAlertChannelIntegration()
	case "Email":
		return createEmailAlertChannelIntegration()
	case "GCP PubSub":
		return createGcpPubSubChannelIntegration()
	case "Microsoft Teams":
		return createMicrosoftTeamsChannelIntegration()
	case "New Relic Insights":
		return createNewRelicAlertChannelIntegration()
	case "Amazon S3":
		return createAwsS3ChannelIntegration()
	case "Cisco Webex":
		return createCiscoWebexChannelIntegration()
	case "Datadog":
		return createDatadogIntegration()
	case "Webhook":
		return createWebhookIntegration()
	case "VictorOps":
		return createVictorOpsChannelIntegration()
	case "Splunk":
		return createSplunkIntegration()
	case "PagerDuty":
		return createPagerDutyAlertChannelIntegration()
	case "QRadar":
		return createQRadarAlertChannelIntegration()
	case "Service Now":
		return createServiceNowAlertChannelIntegration()
	case "Amazon CloudWatch":
		return createAwsCloudWatchAlertChannelIntegration()
	case "Jira Cloud":
		return createJiraCloudAlertChannelIntegration()
	case "Jira Server":
		return createJiraServerAlertChannelIntegration()
	default:
		return errors.New("unknown alert channel type")
	}
}
