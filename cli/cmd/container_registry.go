package cmd

import (
	"github.com/AlecAivazis/survey/v2"
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
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available container registry integrations",
		Args:    cobra.NoArgs,
		RunE:    containerRegistryList,
	}

	// containerRegistryShowCmd represents the show sub-command inside the container registry command
	containerRegistryShowCmd = &cobra.Command{
		Use:     "show",
		Aliases: []string{"get"},
		Short:   "Show a single container registry integration",
		Args:    cobra.ExactArgs(1),
		RunE:    containerRegistryShow,
	}

	// containerRegistryCreateCmd represents the show sub-command inside the container registries command
	containerRegistryCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new container registry integration",
		Args:  cobra.NoArgs,
		RunE:  containerRegistryCreate,
	}

	// containerRegistryDeleteCmd represents the delete sub-command inside the container registries command
	containerRegistryDeleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "Delete a container registry integration",
		Args:    cobra.ExactArgs(1),
		RunE:    containerRegistryDelete,
	}
)

func init() {
	// add the container-registry command
	rootCmd.AddCommand(containerRegistryCommand)
	containerRegistryCommand.AddCommand(containerRegistryListCmd)
	containerRegistryCommand.AddCommand(containerRegistryShowCmd)
	containerRegistryCommand.AddCommand(containerRegistryCreateCmd)
	containerRegistryCommand.AddCommand(containerRegistryDeleteCmd)
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

func containerRegistryList(_ *cobra.Command, _ []string) error {
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
}

func containerRegistryShow(_ *cobra.Command, args []string) error {
	var (
		containerRegistry api.ContainerRegistryResponse
		out               [][]string
	)
	cli.StartProgress(" Fetching container registry...")
	err := cli.LwApi.V2.ContainerRegistries.Get(args[0], &containerRegistry)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve container registry")
	}

	out = append(out, []string{containerRegistry.Data.IntgGuid,
		containerRegistry.Data.Name,
		containerRegistry.Data.Type,
		containerRegistry.Data.Status(),
		containerRegistry.Data.StateString()})

	cli.OutputHuman(renderSimpleTable([]string{"Container Registry GUID", "Name", "Type", "Status", "State"}, out))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildDetailsTable(containerRegistry.Data))
	return nil
}

func containerRegistryCreate(_ *cobra.Command, args []string) error {
	if !cli.InteractiveMode() {
		return errors.New("interactive mode is disabled")
	}

	err := promptCreateContainerRegistry()
	if err != nil {
		return errors.Wrap(err, "unable to create integration")
	}

	cli.OutputHuman("The integration was created.\n")
	return nil
}

func containerRegistryDelete(_ *cobra.Command, args []string) error {
	cli.StartProgress(" Deleting container registry...")
	err := cli.LwApi.V2.ContainerRegistries.Delete(args[0])
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to delete container registry")
	}
	cli.OutputHuman("The container registry %s was deleted.\n", args[0])
	return nil
}

func promptCreateContainerRegistry() error {
	var (
		integration = ""
		prompt      = &survey.Select{
			Message: "Choose a container registry type to create: ",
			Options: []string{
				"Docker Hub Registry",
				"Docker V2 Registry",
				"Amazon Container Registry (ECR)",
				"Google Container Registry (GCR)",
				"Google Artifact Registry (GAR)",
				"Github Container Registry (GHCR)",
				"Inline Scanner Container Registry",
			},
		}
		err = survey.AskOne(prompt, &integration)
	)
	if err != nil {
		return err
	}

	switch integration {
	case "Docker Hub Registry":
		return createDockerHubIntegration()
	case "Docker V2 Registry":
		return createDockerV2Integration()
	case "Amazon Container Registry (ECR)":
		return createAwsEcrIntegration()
	case "Google Artifact Registry (GAR)":
		return createGarIntegration()
	case "Inline Scanner Container Registry":
		return createInlineScannerIntegration()
	case "Github Container Registry (GHCR)":
		return createGhcrIntegration()
	case "Google Container Registry (GCR)":
		return createGcrIntegration()
	default:
		return errors.New("unknown container registry type")
	}
}
