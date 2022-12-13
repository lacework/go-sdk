package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// top-level cloud-account command
	suppressionsCommand = &cobra.Command{
		Use:     "suppressions",
		Aliases: []string{"suppression", "sup", "sups"},
		Short:   "Manage legacy suppressions",
		Long:    "Manage legacy suppressions",
	}

	// suppressionsListCmd represents the list sub-command inside the suppressions command
	suppressionsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all legacy suppressions per CSP",
	}

	// suppressionsMigrateCmd represents the migrate sub-command inside the suppressions command
	suppressionsMigrateCmd = &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"mig"},
		Short:   "Migrate legacy suppressions for selected CSP, to policy exceptions",
	}
)

func init() {
	// add the suppressions command
	rootCmd.AddCommand(suppressionsCommand)

	suppressionsCommand.AddCommand(suppressionsListCmd)
	suppressionsListCmd.AddCommand(suppressionsListAwsCmd)

	suppressionsCommand.AddCommand(suppressionsMigrateCmd)
	suppressionsMigrateCmd.AddCommand(suppressionsMigrateAwsCmd)
}
