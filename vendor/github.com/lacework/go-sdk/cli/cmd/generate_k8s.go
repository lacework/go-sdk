package cmd

import "github.com/spf13/cobra"

var (
	generateK8sCommand = &cobra.Command{
		Use:   "k8s",
		Short: "Generate Kubernetes integration IaC",
		Long: `Generate IaC to deploy Lacework into a Kubernetes platform.

This command creates Infrastructure as Code (IaC) in the form of Terraform HCL, with the option of running
Terraform and deploying Lacework into GKE.
`,
	}
)

func init() {
	generateTfCommand.AddCommand(generateK8sCommand)

	initGenerateGkeTfCommandFlags()
	initGenerateAwsEksAuditTfCommandFlags()
	generateK8sCommand.AddCommand(generateGkeTfCommand)
	generateK8sCommand.AddCommand(generateAwsEksAuditTfCommand)
}
