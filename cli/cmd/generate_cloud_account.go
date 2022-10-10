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

	// Uncomment when `cloud-account iac-generate` removed
	//	initGenerateAwsTfCommandFlags()
	//	initGenerateGcpTfCommandFlags()
	//	initGenerateAzureTfCommandFlags()

	generateCloudAccountCommand.AddCommand(generateAwsTfCommand)
	generateCloudAccountCommand.AddCommand(generateGcpTfCommand)
	generateCloudAccountCommand.AddCommand(generateAzureTfCommand)
}
