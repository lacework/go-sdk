package cmd

import (
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

func alertChannelsToTable(integrations []api.AlertChannelRaw) [][]string {
	var out [][]string
	for _, cadata := range integrations {
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

	if len(alertChannels.Data) == 0 {
		cli.OutputHuman("No alert channels found.\n")
		return nil
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(alertChannels.Data)
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
	cli.StartProgress(" Fetching alert channel...")
	err := cli.LwApi.V2.AlertChannels.Get(args[0], &alertChannel)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve alert channel")
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

func alertChannelCreate(_ *cobra.Command, args []string) error {
	// Todo: alert channel create
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
