package cmd

import (
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
)

func init() {
	// add the cloud-account command
	rootCmd.AddCommand(cloudAccountCommand)
}
