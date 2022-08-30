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
		Use:   "list",
		Short: "List all available alert channel integrations",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
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
		},
	}
)

func init() {
	// add the alert-channel command
	rootCmd.AddCommand(alertChannelCommand)
	alertChannelCommand.AddCommand(alertChannelListCmd)
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
