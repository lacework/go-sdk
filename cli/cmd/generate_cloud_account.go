package cmd

import "github.com/spf13/cobra"

var (
	generateCloudAccountCommand = &cobra.Command{
		Use:     "cloud-account",
		Aliases: []string{"cloud", "ca"},
		Short:   "Create IaC code",
		Long:    "Create IaC content for various different cloud environments",
	}
)

func init() {
	generateTfCommand.AddCommand(generateCloudAccountCommand)

	initGenerateAwsTfCommandFlags()
	initGenerateGcpTfCommandFlags()
	initGenerateAzureTfCommandFlags()

	generateCloudAccountCommand.AddCommand(generateAwsTfCommand)
	generateCloudAccountCommand.AddCommand(generateGcpTfCommand)
	generateCloudAccountCommand.AddCommand(generateAzureTfCommand)
}
