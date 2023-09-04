package cmd

import "github.com/spf13/cobra"

var (
	generateCloudAccountCommand = &cobra.Command{
		Use:     "cloud-account",
		Aliases: []string{"cloud", "ca"},
		Short:   "Generate cloud integration IaC",
		Long: `Generate cloud-account IaC to deploy Lacework into a cloud environment.

This command creates Infrastructure as Code (IaC) in the form of Terraform HCL, with the option of running
Terraform and deploying Lacework into AWS, Azure, or GCP.
`,
	}
)

func init() {
	generateTfCommand.AddCommand(generateCloudAccountCommand)

	// Uncomment when `cloud-account iac-generate` removed
	//	initGenerateAwsTfCommandFlags()
	//	initGenerateGcpTfCommandFlags()
	//	initGenerateAzureTfCommandFlags()
	initGenerateOciTfCommandFlags()

	generateCloudAccountCommand.AddCommand(generateAwsTfCommand)
	generateCloudAccountCommand.AddCommand(generateGcpTfCommand)
	generateCloudAccountCommand.AddCommand(generateAzureTfCommand)
	generateCloudAccountCommand.AddCommand(generateOciTfCommand)
}
