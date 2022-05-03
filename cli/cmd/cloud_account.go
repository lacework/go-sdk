package cmd

import (
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// top-level cloud-account command
	cloudAccountCommand = &cobra.Command{
		Use:     "cloud-account",
		Aliases: []string{"cloud-accounts", "cloud", "ca"},
		Short:   "Manage cloud accounts",
		Long:    "Manage cloud account integrations with Lacework",
	}

	// used by cloud account list to list only a single type of cloud account
	cloudAccountType string

	// cloudAccountsListCmd represents the list sub-command inside the cloud accounts command
	cloudAccountListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available cloud account integrations",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			var (
				cloudAccounts api.CloudAccountsResponse
				err           error
			)

			if cloudAccountType != "" {
				caType, found := api.FindCloudAccountType(cloudAccountType)
				if !found {
					return errors.Errorf("unknown cloud account type '%s'", cloudAccountType)
				}
				cloudAccounts, err = cli.LwApi.V2.CloudAccounts.ListByType(caType)
			} else {
				cloudAccounts, err = cli.LwApi.V2.CloudAccounts.List()
			}
			if err != nil {
				return errors.Wrap(err, "unable to get cloud accounts")
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(cloudAccounts.Data)
			}

			if len(cloudAccounts.Data) == 0 {
				cli.OutputHuman("There was no cloud account found.\n")
				return nil
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"Cloud Account Integration GUID", "Name", "Type", "Status", "State"},
					cloudAccountsToTable(cloudAccounts.Data),
				),
			)
			return nil
		},
	}
)

func init() {
	// add the cloud-account command
	rootCmd.AddCommand(cloudAccountCommand)
	cloudAccountCommand.AddCommand(cloudAccountListCmd)

	// add type flag to cloud accounts list command
	cloudAccountListCmd.Flags().StringVarP(&cloudAccountType,
		"type", "t", "", "list all cloud accounts of a specific type",
	)
}

func cloudAccountsToTable(integrations []api.CloudAccountRaw) [][]string {
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
