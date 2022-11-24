package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// vulContainerListRegistriesCmd represents the list-registries sub-command inside the container
	// vulnerability command
	vulContainerListRegistriesCmd = &cobra.Command{
		Use:     "list-registries",
		Aliases: []string{"list-reg", "registries"},
		Short:   "List all container registries configured",
		Long:    `List all container registries configured in your account.`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			registries, err := getContainerRegistries()
			if err != nil {
				return err
			}
			if len(registries) == 0 {
				msg := `There are no container registries configured in your account.

Get started by integrating your container registry using the command:

    lacework integration create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Integrations > Container Registry.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(registries)
			}

			var rows [][]string
			for _, acc := range registries {
				rows = append(rows, []string{acc})
			}

			cli.OutputHuman(renderSimpleTable([]string{"Container Registries"}, rows))
			return nil
		},
	}
)
