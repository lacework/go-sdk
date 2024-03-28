package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/lacework/go-sdk/internal/format"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
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

	// used by cloud account list to list only a single type of cloud account
	cloudAccountType string

	// cloudAccountsListCmd represents the list sub-command inside the cloud accounts command
	cloudAccountListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available cloud account integrations",
		Args:    cobra.NoArgs,
		RunE:    cloudAccountList,
	}

	// cloudAccountShowCmd represents the show sub-command inside the cloud accounts command
	cloudAccountShowCmd = &cobra.Command{
		Use:     "show",
		Aliases: []string{"get"},
		Short:   "Show a single cloud account integration",
		Args:    cobra.ExactArgs(1),
		RunE:    cloudAccountShow,
	}

	// cloudAccountCreateCmd represents the show sub-command inside the cloud accounts command
	cloudAccountCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new cloud account integration",
		Args:  cobra.NoArgs,
		RunE:  cloudAccountCreate,
	}

	// cloudAccountDeleteCmd represents the delete sub-command inside the cloud accounts command
	cloudAccountDeleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "Delete a cloud account integration",
		Args:    cobra.ExactArgs(1),
		RunE:    cloudAccountDelete,
	}

	// cloudAccountMigrateCmd represents the migrate sub-command inside the cloud accounts command
	cloudAccountMigrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Mark a cloud account integration for migration",
		Args:  cobra.ExactArgs(1),
		RunE:  cloudAccountMigrate,
	}
)

func init() {
	// add the cloud-account command
	rootCmd.AddCommand(cloudAccountCommand)
	cloudAccountCommand.AddCommand(cloudAccountListCmd)
	cloudAccountCommand.AddCommand(cloudAccountShowCmd)
	cloudAccountCommand.AddCommand(cloudAccountDeleteCmd)
	cloudAccountCommand.AddCommand(cloudAccountCreateCmd)
	cloudAccountCommand.AddCommand(cloudAccountMigrateCmd)

	// add type flag to cloud accounts list command
	cloudAccountListCmd.Flags().StringVarP(&cloudAccountType,
		"type", "t", "", "list all cloud accounts of a specific type",
	)
}

func cloudAccountsToTable(cloudAccounts []api.CloudAccountRaw) [][]string {
	var out [][]string
	for _, cadata := range cloudAccounts {
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

func cloudAccountList(_ *cobra.Command, _ []string) error {
	var (
		cloudAccounts api.CloudAccountsResponse
		err           error
	)

	if cloudAccountType != "" {
		caType, found := api.FindCloudAccountType(cloudAccountType)
		if !found {
			return errors.Errorf("unknown cloud account type '%s'", cloudAccountType)
		}
		cloudAccounts, err = cli.LwApi.V2.CloudAccounts.ListByType(caType)
	} else {
		cloudAccounts, err = cli.LwApi.V2.CloudAccounts.List()
	}
	if err != nil {
		return errors.Wrap(err, "unable to get cloud accounts")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(cloudAccounts.Data)
	}

	if len(cloudAccounts.Data) == 0 {
		cli.OutputHuman("No cloud accounts found.\n")
		return nil
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Cloud Account GUID", "Name", "Type", "Status", "State"},
			cloudAccountsToTable(cloudAccounts.Data),
		),
	)
	return nil
}

func cloudAccountDelete(_ *cobra.Command, args []string) error {
	cli.StartProgress(" Deleting cloud account...")
	err := cli.LwApi.V2.CloudAccounts.Delete(args[0])
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to delete cloud account")
	}
	cli.OutputHuman("The cloud account %s was deleted.\n", args[0])
	return nil
}

func cloudAccountMigrate(_ *cobra.Command, args []string) error {
	cli.StartProgress(" Initiating migration for cloud account...")
	err := cli.LwApi.V2.CloudAccounts.Migrate(args[0])
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to initiate migration for cloud-account.")
	}
	cli.OutputHuman("The cloud account %s was marked for migration.\n", args[0])
	return nil
}

func cloudAccountShow(_ *cobra.Command, args []string) error {
	var (
		cloudAccount api.CloudAccountResponse
		out          [][]string
	)
	cli.StartProgress(" Fetching cloud account...")
	err := cli.LwApi.V2.CloudAccounts.Get(args[0], &cloudAccount)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve cloud account")
	}

	out = append(out, []string{cloudAccount.Data.IntgGuid,
		cloudAccount.Data.Name,
		cloudAccount.Data.Type,
		cloudAccount.Data.Status(),
		cloudAccount.Data.StateString()})

	if cli.JSONOutput() {
		return cli.OutputJSON(cloudAccount.Data)
	}

	cli.OutputHuman(renderSimpleTable([]string{"Cloud Account GUID", "Name", "Type", "Status", "State"}, out))
	cli.OutputHuman("\n")
	cli.OutputHuman(buildDetailsTable(cloudAccount.Data))
	return nil
}

func buildDetailsTable(integration api.V2RawType) string {
	var details [][]string

	if caMap, ok := integration.GetData().(map[string]interface{}); ok {
		for k, v := range caMap {
			switch val := v.(type) {
			case int:
				details = append(details, []string{strings.ToUpper(format.SpaceUpperCase(k)), strconv.Itoa(val)})
			case float64:
				details = append(details,
					[]string{strings.ToUpper(format.SpaceUpperCase(k)), strconv.FormatFloat(val, 'f', -1, 64)})
			case string:
				details = append(details, []string{strings.ToUpper(format.SpaceUpperCase(k)), val})
			case map[string]any:
				for i, j := range val {
					if v, ok := j.(string); ok {
						details = append(details, []string{strings.ToUpper(format.SpaceUpperCase(i)), v})
					} else {
						cli.Log.Warn("unable to build table details, unknown type", "type", i, "key", k)
					}
				}
			case []any:
				var values []string
				for _, i := range val {
					if _, ok := i.(string); ok {
						values = append(values, i.(string))
					} else if _, ok := i.(map[string]interface{}); ok {
						for m, n := range i.(map[string]interface{}) {
							values = append(values, fmt.Sprintf("%s:%s", m, n))
						}
					} else {
						cli.Log.Warn("unable to build table details, unknown type", "type", i, "key", k)
					}
				}
				details = append(details, []string{strings.ToUpper(format.SpaceUpperCase(k)), strings.Join(values, ",")})
			}
		}
	}

	// get server token for container registry type only
	if c, ok := integration.(api.ContainerRegistryRaw); ok {
		if c.ServerToken != nil {
			details = append(details, []string{"SERVER_TOKEN", c.ServerToken.ServerToken})
			details = append(details, []string{"SERVER_TOKEN_URI", c.ServerToken.Uri})
		}
	}

	//common
	details = append(details, []string{"UPDATED AT", integration.GetCommon().CreatedOrUpdatedTime})
	details = append(details, []string{"UPDATED BY", integration.GetCommon().CreatedOrUpdatedBy})
	if integration.GetCommon().State != nil {
		details = append(details, []string{
			"STATE UPDATED AT", integration.GetCommon().State.LastUpdatedTime.Format(time.RFC3339)})
		details = append(details, []string{
			"LAST SUCCESSFUL STATE", integration.GetCommon().State.LastSuccessfulTime.Format(time.RFC3339)})

		//output state details as a json string
		jsonState, err := json.Marshal(integration.GetCommon().State.Details)
		if err == nil {
			detailsJSON, err := cli.FormatJSONString(string(jsonState))
			if err == nil {
				details = append(details, []string{"STATE DETAILS", detailsJSON})
			}
		}
	}

	array.Sort2D(details)
	return renderOneLineCustomTable("DETAILS",
		renderCustomTable([]string{}, details,
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(false)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
			}),
		),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	)

}

func cloudAccountCreate(_ *cobra.Command, _ []string) error {
	if !cli.InteractiveMode() {
		return errors.New("interactive mode is disabled")
	}

	err := promptCreateCloudAccount()
	if err != nil {
		return errors.Wrap(err, "unable to create cloud account")
	}

	cli.OutputHuman("The cloud account was created.\n")
	return nil
}

func promptCreateCloudAccount() error {
	var (
		cloudAccount = ""
		prompt       = &survey.Select{
			Message: "Choose a cloud account type to create: ",
			Options: []string{
				"AWS Config",
				"AWS CloudTrail",
				"AWS Config (US GovCloud)",
				"AWS CloudTrail (US GovCloud)",
				"GCP Config",
				"GCP Audit Log",
				"GCP Audit Log PubSub",
				"Azure Config",
				"Azure Activity Log",
				"OCI Config",
			},
		}
		err = survey.AskOne(prompt, &cloudAccount)
	)
	if err != nil {
		return err
	}

	switch cloudAccount {
	case "AWS Config":
		return createAwsConfigIntegration()
	case "AWS CloudTrail":
		return createAwsCloudTrailIntegration()
	case "AWS GovCloud Config":
		return createAwsGovCloudConfigIntegration()
	case "AWS GovCloud CloudTrail":
		return createAwsGovCloudCTIntegration()
	case "GCP Config":
		return createGcpConfigIntegration()
	case "GCP Audit Log":
		return createGcpAuditLogIntegration()
	case "GCP Audit Log PubSub":
		return createGcpPubSubAuditLogIntegration()
	case "Azure Config":
		return createAzureConfigIntegration()
	case "Azure Activity Log":
		return createAzureActivityLogIntegration()
	case "Azure Active Directory Activity Log":
		return createAzureAdAlIntegration()
	case "OCI Config":
		return createOciConfigIntegration()
	default:
		return errors.New("unknown cloud account type")
	}
}
