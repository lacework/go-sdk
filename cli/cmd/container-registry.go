package cmd

import (
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	containerRegistryCommand = &cobra.Command{
		Use:     "container-registry",
		Aliases: []string{"container-registries", "cr"},
		Short:   "Manage container registries",
		Long:    "Manage container registries integrations with Lacework",
	}

	// containerRegistriesListCmd represents the list sub-command inside the container registries command
	containerRegistryListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available container registry integrations",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			containerRegistries, err := cli.LwApi.V2.ContainerRegistries.List()

			if err != nil {
				return errors.Wrap(err, "unable to get container registries")
			}

			if len(containerRegistries.Data) == 0 {
				cli.OutputHuman("No container registries found.\n")
				return nil
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(containerRegistries.Data)
			}

			cli.OutputHuman(
				renderSimpleTable(
					[]string{"container registry GUID", "Name", "Type", "Status", "State"},
					containerRegistriesToTable(containerRegistries.Data),
				),
			)
			return nil
		},
	}
)

func init() {
	// add the container-registry command
	rootCmd.AddCommand(containerRegistryCommand)
	containerRegistryCommand.AddCommand(containerRegistryListCmd)
}

func containerRegistriesToTable(integrations []api.ContainerRegistryRaw) [][]string {
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
