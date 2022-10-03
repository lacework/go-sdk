package cmd

import "github.com/spf13/cobra"

var (
	generateCloudAccountCommand = &cobra.Command{
		Use:     "cloud-account",
		Aliases: []string{"cloud", "ca"},
		Short:   "Generate Lacework integration IaC",
		Long:    "Generate Lacework integration IaC for cloud environment(s)",
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
