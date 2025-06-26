package cmd

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/lwgenerate/azure"
	"github.com/pkg/errors"
)

// Question labels
const (
	IconAzureConfig 	= "[Configuration]"
	IconActivityLog 	= "[Activity Log]"
	IconEntraID     	= "[Entra ID Activity Log]"
	IconAD          	= "[Active Directory Application]"
	IconAzureAgentless 	= "[Agentless]"
)

var (
	// Define question text here so they can be reused in testing
	// Core questions
	QuestionAzureEnableAgentless     = "Enable Agentless integration?"
	QuestionAzureEnableConfig        = "Enable Configuration integration?"
	QuestionAzureConfigName          = "Custom Configuration integration name: (optional)"
	QuestionEnableActivityLog        = "Enable Activity Log Integration?"
	QuestionActivityLogName          = "Custom Activity Log integration name: (optional)"
	QuestionEnableEntraIdActivityLog = "Enable Entra ID Activity Log Integration?"
	QuestionEntraIdActivityLogName   = "Custom EntraID Activity Log integration name: (optional)"
	QuestionAzureSubscriptionID      = "Subscription ID to be used to provision Lacework resources:"

	// Advanced options
	QuestionAzureCustomizeOutputLocation = "Provide the location for the output to be written: (optional)"

	// EntraID Activity Log
	QuestionEventHubLocation       = "Region where the event hub for logging will reside:"
	QuestionEventHubPartitionCount = "Number of partitions in the event hub for logging:"

	// Active Directory
	QuestionEnableAdIntegration = "Create Active Directory Application?"
	QuestionADApplicationPass   = "Password of an existing Active Directory application:"
	QuestionADApplicationId     = "ID of an existing Active Directory application:"
	QuestionADServicePrincpleId = "Service Principle ID of an existing Active Directory application:"

	// Storage Account
	QuestionUseExistingStorageAccount   = "Use an existing Storage Account?"
	QuestionAzureRegion                 = "Region to be used by Storage Account logging:"
	QuestionStorageAccountName          = "Existing Storage Account name:"
	QuestionStorageAccountResourceGroup = "Existing Storage Account Resource Group:"
	QuestionStorageLocation             = "Region where Storage Account for logging resides: (optional)"

	// Subscriptions
	QuestionEnableAllSubscriptions = "Enable all subscriptions?"
	QuestionSubscriptionIds        = "List of subscription ids to enable logging:"

	// Management Group
	QuestionEnableManagementGroup = "Enable Management Group level Integration?"
	QuestionManagementGroupId     = "Management Group ID:"

	// Select options
	AzureSubscriptions = "Configure Subscriptions (optional)"

	// Regex patterns for validation
	AzureSubscriptionIDRegex = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

	GenerateAzureCommandState      = &azure.GenerateAzureTfConfigurationArgs{}
	GenerateAzureCommandExtraState = &AzureGenerateCommandExtraState{}
	CachedAzureAssetIacParams      = "iac-azure-generate-params"
	CachedAzureAssetExtraState     = "iac-azure-extra-state"

	// List of valid Azure Storage locations
	validAzureLocations = map[string]bool{
		"East US":                  true,
		"East US 2":                true,
		"South Central US":         true,
		"West US 2":                true,
		"West US 3":                true,
		"Australia East":           true,
		"Southeast Asia":           true,
		"North Europe":             true,
		"Sweden Central":           true,
		"UK South":                 true,
		"West Europe":              true,
		"Central US":               true,
		"North Central US":         true,
		"West US":                  true,
		"South Africa North":       true,
		"Central India":            true,
		"East Asia":                true,
		"Japan East":               true,
		"Jio India West":           true,
		"Korea Central":            true,
		"Canada Central":           true,
		"France Central":           true,
		"Germany West Central":     true,
		"Norway East":              true,
		"Switzerland North":        true,
		"UAE North":                true,
		"Brazil South":             true,
		"Central US (Stage)":       true,
		"East US (Stage)":          true,
		"East US 2 (Stage)":        true,
		"North Central US (Stage)": true,
		"South Central US (Stage)": true,
		"West US (Stage)":          true,
		"West US 2 (Stage)":        true,
		"Asia":                     true,
		"Asia Pacific":             true,
		"Australia":                true,
		"Brazil":                   true,
		"Canada":                   true,
		"Europe":                   true,
		"France":                   true,
		"Germany":                  true,
		"Global":                   true,
		"India":                    true,
		"Japan":                    true,
		"Korea":                    true,
		"Norway":                   true,
		"South Africa":             true,
		"Switzerland":              true,
		"United Arab Emirates":     true,
		"United Kingdom":           true,
		"United States":            true,
		"United States EUAP":       true,
		"East Asia (Stage)":        true,
		"Southeast Asia (Stage)":   true,
		"Central US EUAP":          true,
		"East US 2 EUAP":           true,
		"West Central US":          true,
		"South Africa West":        true,
		"Australia Central":        true,
		"Australia Central 2":      true,
		"Australia Southeast":      true,
		"Japan West":               true,
		"Jio India Central":        true,
		"Korea South":              true,
		"South India":              true,
		"West India":               true,
		"Canada East":              true,
		"France South":             true,
		"Germany North":            true,
		"Norway West":              true,
		"Switzerland West":         true,
		"UK West":                  true,
		"UAE Central":              true,
		"Brazil Southeast":         true,
	}

	// Azure command used to generate TF code for azure
	generateAzureTfCommand = &cobra.Command{
		Use:     "azure",
		Aliases: []string{"az"},
		Short:   "Generate and/or execute Terraform code for Azure integration",
		Long: `Use this command to generate Terraform code for deploying Lacework into new Azure environment.

By default, this command will function interactively, prompting for the required information to setup
the new cloud account. In interactive mode, this command will:

* Prompt for the required information to setup the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version will be confirmed suitable for use
  * If Terraform is not installed, or the version installed is not suitable, a new version will be
    installed into a temporary location
  * Once Terraform is detected or installed, Terraform plan will be executed
  * The command will prompt with the outcome of the plan and allow to view more details or continue
    with Terraform apply
  * If confirmed, Terraform apply will be run, completing the setup of the cloud account
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate TF Code
			cli.StartProgress("Generating Azure Terraform Code...")

			if cli.Profile != "default" {
				GenerateAzureCommandState.LaceworkProfile = cli.Profile
			}

			// Setup modifiers for NewTerraform constructor
			mods := []azure.AzureTerraformModifier{
				azure.WithLaceworkProfile(GenerateAzureCommandState.LaceworkProfile),
				azure.WithSubscriptionID(GenerateAzureCommandState.SubscriptionID),
				azure.WithAllSubscriptions(GenerateAzureCommandState.AllSubscriptions),
				azure.WithManagementGroup(GenerateAzureCommandState.ManagementGroup),
				azure.WithExistingStorageAccount(GenerateAzureCommandState.ExistingStorageAccount),
				azure.WithStorageAccountName(GenerateAzureCommandState.StorageAccountName),
				azure.WithStorageLocation(GenerateAzureCommandState.StorageLocation),
				azure.WithActivityLogIntegrationName(GenerateAzureCommandState.ActivityLogIntegrationName),
				azure.WithConfigIntegrationName(GenerateAzureCommandState.ConfigIntegrationName),
				azure.WithEntraIdActivityLogIntegrationName(GenerateAzureCommandState.EntraIdIntegrationName),
				azure.WithEventHubLocation(GenerateAzureCommandState.EventHubLocation),
				azure.WithEventHubPartitionCount(GenerateAzureCommandState.EventHubPartitionCount),
			}

			// Check if AD Creation is required, need to set values for current integration
			if !GenerateAzureCommandState.CreateAdIntegration {
				mods = append(mods, azure.WithAdApplicationId(GenerateAzureCommandState.AdApplicationId))
				mods = append(mods, azure.WithAdApplicationPassword(GenerateAzureCommandState.AdApplicationPassword))
				mods = append(mods, azure.WithAdServicePrincipalId(GenerateAzureCommandState.AdServicePrincipalId))
			}

			// Check subscriptions
			if !GenerateAzureCommandState.AllSubscriptions {
				if len(GenerateAzureCommandState.SubscriptionIds) > 0 {
					mods = append(mods, azure.WithSubscriptionIds(GenerateAzureCommandState.SubscriptionIds))
				}
			}

			// Check management groups
			if GenerateAzureCommandState.ManagementGroup {
				mods = append(mods, azure.WithManagementGroupId(GenerateAzureCommandState.ManagementGroupId))
			}

			// Check storage account
			if GenerateAzureCommandState.ExistingStorageAccount {
				mods = append(mods,
					azure.WithStorageAccountResourceGroup(GenerateAzureCommandState.StorageAccountResourceGroup))
			}

			// Create new struct
			data := azure.NewTerraform(
				GenerateAzureCommandState.Config,
				GenerateAzureCommandState.ActivityLog,
				GenerateAzureCommandState.Agentless,
				GenerateAzureCommandState.EntraIdActivityLog,
				GenerateAzureCommandState.CreateAdIntegration,
				mods...)

			// Generate HCL for azure deployment
			hcl, err := data.Generate()
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "failed to generate terraform code")
			}

			// Write-out generated code to location specified
			dirname, _, err := writeGeneratedCodeToLocation(cmd, hcl, "azure")
			if err != nil {
				return err
			}

			// Prompt to execute, if the command line flag has not been set
			if !GenerateAzureCommandExtraState.TerraformApply {
				err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
					Prompt:   &survey.Confirm{Default: GenerateAzureCommandExtraState.TerraformApply, Message: QuestionRunTfPlan},
					Response: &GenerateAzureCommandExtraState.TerraformApply,
				})

				if err != nil {
					return errors.Wrap(err, "failed to prompt for terraform execution")
				}
			}

			locationDir, _ := determineOutputDirPath(dirname, "azure")
			if GenerateAzureCommandExtraState.TerraformApply {
				// Execution pre-run check
				err := executionPreRunChecks(dirname, locationDir, "azure")
				if err != nil {
					return err
				}
			}

			// Output where code was generated
			if !GenerateAzureCommandExtraState.TerraformApply {
				cli.OutputHuman(provideGuidanceAfterExit(false, false, locationDir, "terraform"))
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			// Validate output location is OK if supplied
			dirname, err := cmd.Flags().GetString("output")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateOutputLocation(dirname); err != nil {
				return err
			}

			// Validate Storage Location
			storageLocation, err := cmd.Flags().GetString("location")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAzureLocation(storageLocation); storageLocation != "" && err != nil {
				return err
			}

			// Load any cached inputs if interactive
			if cli.InteractiveMode() {
				cachedOptions := &azure.GenerateAzureTfConfigurationArgs{}
				iacParamsExpired := cli.ReadCachedAsset(CachedAzureAssetIacParams, &cachedOptions)
				if iacParamsExpired {
					cli.Log.Debug("loaded previously set values for Azure iac generation")
				}

				extraState := &AzureGenerateCommandExtraState{}
				extraStateParamsExpired := cli.ReadCachedAsset(CachedAzureAssetExtraState, &extraState)
				if extraStateParamsExpired {
					cli.Log.Debug("loaded previously set values for Azure iac generation (extra state)")
				}

				// Determine if previously cached options exists; prompt user if they'd like to continue
				answer := false
				if !iacParamsExpired || !extraStateParamsExpired {
					if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
						Prompt:   &survey.Confirm{Message: QuestionUsePreviousCache, Default: false},
						Response: &answer,
					}); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
				}

				// If the user decides NOT to use the previous values; we won't load them.  However, every time the command runs
				// we are going to write out new cached values, so if they run it - bail out - and run it again they'll get
				// re-prompted.
				if answer {
					// Merge cached inputs to current options (current options win)
					if err := mergo.Merge(GenerateAzureCommandState, cachedOptions); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
					if err := mergo.Merge(GenerateAzureCommandExtraState, extraState); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
				}
			}

			// Collect and/or confirm parameters
			err = promptAzureGenerate(GenerateAzureCommandState, GenerateAzureCommandExtraState)
			if err != nil {
				return errors.Wrap(err, "collecting/confirming parameters")
			}

			return nil
		},
	}
)

type AzureGenerateCommandExtraState struct {
	Output         string
	TerraformApply bool
}

func (a *AzureGenerateCommandExtraState) isEmpty() bool {
	return a.Output == "" && !a.TerraformApply
}

// Flush current state of the struct to disk, provided it's not empty
func (a *AzureGenerateCommandExtraState) writeCache() {
	if !a.isEmpty() {
		cli.WriteAssetToCache(CachedAzureAssetExtraState, time.Now().Add(time.Hour*1), a)
	}
}

func validateAzureLocation(val interface{}) error {
	if str, ok := val.(string); ok {
		if str == "" {
			return nil
		}
		if !validAzureLocations[str] {
			return errors.New("invalid Azure region. Please use a valid Azure region like 'East US', 'West Europe', etc.")
		}
	}
	return nil
}

func validateAzureSubscriptionID(val interface{}) error {
	if str, ok := val.(string); ok {
		if matched, _ := regexp.MatchString(AzureSubscriptionIDRegex, str); !matched {
			return errors.New("invalid Azure subscription ID format")
		}
	}
	return nil
}

func initGenerateAzureTfCommandFlags() {
	// Azure sub-command flags
	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.ActivityLog,
		"activity_log",
		false,
		"enable activity log integration")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.ActivityLogIntegrationName,
		"activity_log_integration_name",
		"",
		"specify a custom activity log integration name")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.Agentless,
		"agentless",
		false,
		"enable agentless integration")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.EntraIdActivityLog,
		"entra_id_activity_log",
		false,
		"enable Entra ID activity log integration")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.EntraIdIntegrationName,
		"entra_id_activity_log_integration_name",
		"",
		"specify a custom Entra ID activity log integration name")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.Config,
		"configuration",
		false,
		"enable configuration integration")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.ConfigIntegrationName,
		"configuration_name",
		"",
		"specify a custom configuration integration name")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.SubscriptionID,
		"subscription_id",
		"",
		"specify the Azure Subscription ID to be used to provision Lacework resources")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.CreateAdIntegration,
		"ad_create",
		true,
		"create new active directory integration")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.ManagementGroup,
		"management_group",
		false,
		"management group level integration")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.ManagementGroupId,
		"management_group_id",
		"",
		"specify management group id. Required if mgmt_group provided")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.ExistingStorageAccount,
		"existing_storage",
		false,
		"use existing storage account")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.EventHubLocation,
		"event_hub_location",
		"",
		"specify the location where the Event Hub for logging will reside")

	generateAzureTfCommand.PersistentFlags().IntVar(
		&GenerateAzureCommandState.EventHubPartitionCount,
		"event_hub_partition_count",
		1,
		"specify the number of partitions for the Event Hub")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.StorageAccountName,
		"storage_account_name",
		"",
		"specify storage account name")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.StorageAccountResourceGroup,
		"storage_resource_group",
		"",
		"specify storage resource group")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.StorageLocation,
		"location",
		"",
		"specify azure region where storage account logging resides")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.AllSubscriptions,
		"all_subscriptions",
		false,
		"grant read access to ALL subscriptions within Tenant (overrides `subscription ids`)")

	generateAzureTfCommand.PersistentFlags().StringSliceVar(
		&GenerateAzureCommandState.SubscriptionIds,
		"subscription_ids",
		[]string{},
		`list of subscriptions to grant read access; format is id1,id2,id3`)

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.AdApplicationPassword,
		"ad_pass",
		"",
		"existing active directory application password")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.AdApplicationId,
		"ad_id",
		"",
		"existing active directory application id")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.AdServicePrincipalId,
		"ad_pid",
		"",
		"existing active directory application service principle id")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandExtraState.TerraformApply,
		"terraform-apply",
		false,
		"run terraform apply for the generated hcl")

	_ = generateAzureTfCommand.PersistentFlags().MarkHidden("terraform-apply")

	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandExtraState.TerraformApply,
		"apply",
		false,
		"run terraform apply for the generated hcl")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandExtraState.Output,
		"output",
		"",
		"location to write generated content (default is ~/lacework/azure)",
	)
}

func promptAzureEntraIdQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {
	// Ask for Entra ID integration name
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconEntraID,
			Prompt:   &survey.Input{Message: QuestionEntraIdActivityLogName, Default: config.EntraIdIntegrationName},
			Response: &config.EntraIdIntegrationName,
		},
		{
			Icon: IconEntraID,
			Prompt: &survey.Input{
				Message: QuestionEventHubLocation,
				Default: config.EventHubLocation,
				Help:    "Enter a valid region (e.g., 'East US', 'West Europe')",
			},
			Required: true,
			Response: &config.EventHubLocation,
			Opts:     []survey.AskOpt{survey.WithValidator(validateAzureLocation)},
		},
		{
			Icon: IconEntraID,
			Prompt: &survey.Input{
				Message: QuestionEventHubPartitionCount,
				Default: strconv.Itoa(config.EventHubPartitionCount),
			},
			Response: &config.EventHubPartitionCount,
		},
	}); err != nil {
		return err
	}

	return nil
}

func promptAzureSubscriptionQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {
	// First ask if user wants to configure subscription options
	configureSubscriptions := false
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: AzureSubscriptions, Default: false},
			Response: &configureSubscriptions,
		},
	}); err != nil {
		return err
	}

	// Only proceed with subscription questions if user wants to configure
	if !configureSubscriptions {
		return nil
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionEnableAllSubscriptions, Default: config.AllSubscriptions},
			Response: &config.AllSubscriptions,
		},
	}); err != nil {
		return err
	}
	var idList string
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionSubscriptionIds, Default: strings.Join(config.SubscriptionIds, ",")},
			Checks:   []*bool{allSubscriptionsDisabled(config)},
			Required: true,
			Response: &idList,
		},
	}); err != nil {
		return err
	}
	config.SubscriptionIds = strings.Split(strings.ReplaceAll(idList, " ", ""), ",")

	return nil
}

func promptAzureManagementGroupQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconAzureConfig,
			Prompt:   &survey.Confirm{Message: QuestionEnableManagementGroup, Default: config.ManagementGroup},
			Response: &config.ManagementGroup,
		},
		{
			Icon:     IconAzureConfig,
			Prompt:   &survey.Input{Message: QuestionManagementGroupId, Default: config.ManagementGroupId},
			Checks:   []*bool{&config.ManagementGroup},
			Required: true,
			Response: &config.ManagementGroupId,
		},
	}); err != nil {
		return err
	}
	return nil
}

func promptAzureAdIntegrationQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconAD,
			Prompt:   &survey.Input{Message: QuestionADApplicationPass, Default: config.AdApplicationPassword},
			Required: true,
			Response: &config.AdApplicationPassword,
		},
		{
			Icon:     IconAD,
			Prompt:   &survey.Input{Message: QuestionADApplicationId, Default: config.AdApplicationId},
			Required: true,
			Response: &config.AdApplicationId,
		},
		{
			Icon:     IconAD,
			Prompt:   &survey.Input{Message: QuestionADServicePrincpleId, Default: config.AdServicePrincipalId},
			Required: true,
			Response: &config.AdServicePrincipalId,
		},
	}); err != nil {
		return err
	}
	return nil
}

func promptCustomizeAzureOutputLocation(extraState *AzureGenerateCommandExtraState) error {

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionAzureCustomizeOutputLocation, Default: extraState.Output},
			Response: &extraState.Output,
		},
	}); err != nil {
		return err
	}

	return nil
}

func askAzureSubscriptionID(config *azure.GenerateAzureTfConfigurationArgs) error {
	// if subscription has been set by --subscription_id flag do not prompt
	if config.SubscriptionID != "" {
		return nil
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionAzureSubscriptionID},
			Response: &config.SubscriptionID,
			Opts:     []survey.AskOpt{survey.WithValidator(validateAzureSubscriptionID)},
		},
	}); err != nil {
		return err
	}

	return nil
}

func azureConfigIsEmpty(config *azure.GenerateAzureTfConfigurationArgs) bool {
	return !config.Config && !config.ActivityLog && config.LaceworkProfile == ""
}

func allSubscriptionsDisabled(config *azure.GenerateAzureTfConfigurationArgs) *bool {
	allSubscriptionsDisabled := !config.AllSubscriptions
	return &allSubscriptionsDisabled
}

func writeAzureGenerationArgsCache(config *azure.GenerateAzureTfConfigurationArgs) {
	if !azureConfigIsEmpty(config) {
		cli.WriteAssetToCache(CachedAzureAssetIacParams, time.Now().Add(time.Hour*1), config)
	}
}

// entry point for launching a survey to build out the required Azure generation parameters
func promptAzureGenerate(
	config *azure.GenerateAzureTfConfigurationArgs, extraState *AzureGenerateCommandExtraState,
) error {
	// Cache for later use if generation is abandoned and in interactive mode
	if cli.InteractiveMode() {
		defer writeAzureGenerationArgsCache(config)
		defer extraState.writeCache()
	}

	// Ask subscription ID first as it's common to all integrations
	if err := askAzureSubscriptionID(config); err != nil {
		return err
	}

	// Ask about subscriptions configuration
	if err := promptAzureSubscriptionQuestions(config); err != nil {
		return err
	}

	// Ask Configuration integration
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Icon:     IconAzureConfig,
				Prompt:   &survey.Confirm{Message: QuestionAzureEnableConfig, Default: config.Config},
				Response: &config.Config,
			},
		}); err != nil {
		return err
	}

	// Ask Configuration questions immediately if enabled
	if config.Config {
		if err := promptAzureConfigQuestions(config); err != nil {
			return err
		}
	}

	// Ask Activity Log integration
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Icon:     IconActivityLog,
				Prompt:   &survey.Confirm{Message: QuestionEnableActivityLog, Default: config.ActivityLog},
				Response: &config.ActivityLog,
			},
		}); err != nil {
		return err
	}

	// Ask Activity Log questions immediately if enabled
	if config.ActivityLog {
		if err := promptAzureActivityLogQuestions(config); err != nil {
			return err
		}
	}

	// Ask Agentless integration
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Icon:     IconAzureAgentless,
				Prompt:   &survey.Confirm{Message: QuestionAzureEnableAgentless, Default: config.Agentless},
				Response: &config.Agentless,
			},
		}); err != nil {
		return err
	}

	// Ask Activity Log questions immediately if enabled
	if config.Agentless {
		if err := promptAzureAgentlessQuestions(config); err != nil {
			return err
		}
	}

	// Ask Entra ID integration
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Icon:     IconEntraID,
				Prompt:   &survey.Confirm{Message: QuestionEnableEntraIdActivityLog, Default: config.EntraIdActivityLog},
				Response: &config.EntraIdActivityLog,
			},
		}); err != nil {
		return err
	}

	// Ask Entra ID questions immediately if enabled
	if config.EntraIdActivityLog {
		if err := promptAzureEntraIdQuestions(config); err != nil {
			return err
		}
	}

	// Validate one of config or activity log was enabled; otherwise error out
	if !config.Config && !config.ActivityLog && !config.Agentless && !config.EntraIdActivityLog {
		return errors.New("must enable at least one of: Configuration, Agentless or Activity Log integrations")
	}

	// Ask AD integration
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Icon:     IconAD,
				Prompt:   &survey.Confirm{Message: QuestionEnableAdIntegration, Default: config.CreateAdIntegration},
				Response: &config.CreateAdIntegration,
			},
		}); err != nil {
		return err
	}

	// If AD integration is not being created, ask for existing AD details immediately
	if !config.CreateAdIntegration {
		if err := promptAzureAdIntegrationQuestions(config); err != nil {
			return err
		}
	}

	// Ask about output location
	if err := promptCustomizeAzureOutputLocation(extraState); err != nil {
		return err
	}

	return nil
}

func promptAzureConfigQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {
	// Ask for config integration name
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconAzureConfig,
			Prompt:   &survey.Input{Message: QuestionAzureConfigName, Default: config.ConfigIntegrationName},
			Response: &config.ConfigIntegrationName,
		},
	}); err != nil {
		return err
	}

	// Ask about management group if config is enabled
	if err := promptAzureManagementGroupQuestions(config); err != nil {
		return err
	}

	return nil
}

func promptAzureActivityLogQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {
	// Ask for activity log integration name
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconActivityLog,
			Prompt:   &survey.Input{Message: QuestionActivityLogName, Default: config.ActivityLogIntegrationName},
			Response: &config.ActivityLogIntegrationName,
		},
	}); err != nil {
		return err
	}

	// Ask about storage account configuration
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconActivityLog,
			Prompt:   &survey.Confirm{Message: QuestionUseExistingStorageAccount, Default: config.ExistingStorageAccount},
			Response: &config.ExistingStorageAccount,
		},
	}); err != nil {
		return err
	}

	// If using existing storage account, ask for its details
	if config.ExistingStorageAccount {
		if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
			{
				Icon:     IconActivityLog,
				Prompt:   &survey.Input{Message: QuestionStorageAccountName, Default: config.StorageAccountName},
				Required: true,
				Response: &config.StorageAccountName,
			},
			{
				Icon:     IconActivityLog,
				Prompt:   &survey.Input{Message: QuestionStorageAccountResourceGroup, Default: config.StorageAccountResourceGroup},
				Required: true,
				Response: &config.StorageAccountResourceGroup,
			},
		}); err != nil {
			return err
		}
	}

	// Ask for storage location with validation
	var region string
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconActivityLog,
			Prompt: &survey.Input{
				Message: QuestionStorageLocation,
				Default: config.StorageLocation,
				Help:    "Enter a valid Azure region (e.g., 'East US', 'West Europe')",
			},
			Response: &region,
			Opts:     []survey.AskOpt{survey.WithValidator(validateAzureLocation)},
		},
	}); err != nil {
		return err
	}
	config.StorageLocation = region

	return nil
}


func promptAzureAgentlessQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {
	// Ask for Agentless integration
	// if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
	// 	{
	// 		Icon:     IconEntraID,
	// 		Prompt:   &survey.Input{Message: QuestionEntraIdActivityLogName, Default: config.EntraIdIntegrationName},
	// 		Response: &config.EntraIdIntegrationName,
	// 	},
	// }); err != nil {
	// 	return err
	// }

	return nil
}