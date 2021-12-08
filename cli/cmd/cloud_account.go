package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// top-level cloud-account command
	cloudAccountCommand = &cobra.Command{
		Use:     "cloud-account",
		Aliases: []string{"cloud", "ca"},
		Short:   "Manage cloud account(s)",
		Long:    "Manage cloud account integration with Lacework",
	}
)

func init() {
	// add the cloud-account command
	rootCmd.AddCommand(cloudAccountCommand)
}
