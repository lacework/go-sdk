package cmd

import (
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/lwgenerate/azure"
	"github.com/pkg/errors"
)

var (
	// Define question text here so they can be reused in testing
	QuestionAzureEnableConfig = "Enable Azure configuration integration?"
	QuestionAzureConfigName   = "Specify custom configuration integration name: (optional)"
	QuestionEnableActivityLog = "Enable Azure Activity Log Integration?"
	QuestionActivityLogName   = "Specify custom Activity Log integration name: (optional)"

	QuestionAzureAnotherAdvancedOpt      = "Configure another advanced integration option"
	QuestionAzureConfigAdvanced          = "Configure advanced integration options?"
	QuestionAzureCustomizeOutputLocation = "Provide the location for the output to be written:"

	// Active Directory
	QuestionEnableAdIntegration = "Create Active Directory Integration?"
	QuestionADApplicationPass   = "Specify the password of an existing Active Directory application"
	QuestionADApplicationId     = "Specify the ID of an existing Active Directory application"
	QuestionADServicePrincpleId = "Specify the Service Principle ID of an existing Active Directory application"

	// Storage Account
	QuestionUseExistingStorageAccount   = "Use an existing Storage Account?"
	QuestionAzureRegion                 = "Specify the Azure region to be used by Storage Account logging"
	QuestionStorageAccountName          = "Specify existing Storage Account name"
	QuestionStorageAccountResourceGroup = "Specify existing Storage Account Resource Group"

	QuestionStorageLocation = "Specify Azure region where Storage Account for logging resides "

	// Subscriptions
	QuestionEnableAllSubscriptions = "Enable all subscriptions?"
	QuestionSubscriptionIds        = "Specify list of subscription ids to enable logging"

	// Management Group
	QuestionEnableManagementGroup = "Enable Management Group level Integration?"
	QuestionManagementGroupId     = "Specify Management Group ID"

	// Select options
	AzureAdvancedOptDone       = "Done"
	AdvancedAdIntegration      = "Configure Lacework integration with an existing Active Directory (optional)"
	AzureExistingStorageAcount = "Configure Storage Account (optional)"
	AzureSubscriptions         = "Configure Subscriptions (optional)"
	AzureManagmentGroup        = "Configure Management Group (optional)"
	AzureStorageGroup          = "Configure Storage Group (optional)"
	AzureUserIntegrationNames  = "Customize integration name(s)"
	AzureAdvancedOptLocation   = "Customize output location (optional)"
	AzureRegionStorage         = "Customize Azure region for Storage Account (optional)"

	GenerateAzureCommandState      = &azure.GenerateAzureTfConfigurationArgs{}
	GenerateAzureCommandExtraState = &AzureGenerateCommandExtraState{}
	CachedAzureAssetIacParams      = "iac-azure-generate-params"
	CachedAzureAssetExtraState     = "iac-azure-extra-state"

	// List of valid Azure Storage locations
	validStorageLocations = map[string]bool{
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

By default, this command will function interactively, prompting for the required information to setup the new cloud account. In interactive mode, this command will:
		
* Prompt for the required information to setup the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version will be confirmed suitable for use
	* If Terraform is not installed, or the version installed is not suitable, a new version will be installed into a temporary location
	* Once Terraform is detected or installed, Terraform plan will be executed
	* The command will prompt with the outcome of the plan and allow to view more details or continue with Terraform apply
	* If confirmed, Terraform apply will be run, completing the setup of the cloud account
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate TF Code
			cli.StartProgress("Generating Azure Terraform Code...")

			// Setup modifiers for NewTerraform constructor
			mods := []azure.AzureTerraformModifier{
				azure.WithAllSubscriptions(GenerateAzureCommandState.AllSubscriptions),
				azure.WithManagementGroup(GenerateAzureCommandState.ManagementGroup),
				azure.WithExistingStorageAccount(GenerateAzureCommandState.ExistingStorageAccount),
				azure.WithStorageAccountName(GenerateAzureCommandState.StorageAccountName),
				azure.WithStorageLocation(GenerateAzureCommandState.StorageLocation),
				azure.WithActivityLogIntegrationName(GenerateAzureCommandState.ActivityLogIntegrationName),
				azure.WithConfigIntegrationName(GenerateAzureCommandState.ConfigIntegrationName),
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
				mods = append(mods, azure.WithStorageAccountResourceGroup(GenerateAzureCommandState.StorageAccountResourceGroup))
			}

			// Create new struct
			data := azure.NewTerraform(
				GenerateAzureCommandState.Config,
				GenerateAzureCommandState.ActivityLog,
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
			if err := validateStorageLocation(storageLocation); storageLocation != "" && err != nil {
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

func validateStorageLocation(location string) error {
	if !validStorageLocations[location] {
		return errors.New("invalid storage location supplied")
	}
	return nil
}

func initGenerateAzureTfCommandFlags() {
	// Azure sub-command flags
	generateAzureTfCommand.PersistentFlags().BoolVar(
		&GenerateAzureCommandState.ActivityLog,
		"activity_log",
		false,
		"enable active log integration")

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandState.ActivityLogIntegrationName,
		"activity_log_integration_name",
		"",
		"specify a custom activity log integration name")

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

	generateAzureTfCommand.PersistentFlags().StringVar(
		&GenerateAzureCommandExtraState.Output,
		"output",
		"",
		"location to write generated content (default is ~/lacework/azure)",
	)
}

func promptAzureIntegrationNameQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionAzureConfigName, Default: config.ConfigIntegrationName},
			Checks:   []*bool{&config.Config},
			Response: &config.ConfigIntegrationName,
		},
		{
			Prompt:   &survey.Input{Message: QuestionActivityLogName, Default: config.ActivityLogIntegrationName},
			Checks:   []*bool{&config.ActivityLog},
			Response: &config.ActivityLogIntegrationName,
		},
	}); err != nil {
		return err
	}
	return nil
}

func promptAzureStorageAccountQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionUseExistingStorageAccount, Default: config.ExistingStorageAccount},
			Response: &config.ExistingStorageAccount,
		},
		{
			Prompt:   &survey.Input{Message: QuestionStorageAccountName, Default: config.StorageAccountName},
			Required: true,
			Response: &config.StorageAccountName,
		},
		{
			Prompt:   &survey.Input{Message: QuestionStorageAccountResourceGroup, Default: config.StorageAccountResourceGroup},
			Checks:   []*bool{&config.ExistingStorageAccount},
			Required: true,
			Response: &config.StorageAccountResourceGroup,
		},
	}); err != nil {
		return err
	}

	return nil
}

func promptAzureSubscriptionQuestions(config *azure.GenerateAzureTfConfigurationArgs) error {

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
			Prompt:   &survey.Confirm{Message: QuestionEnableManagementGroup, Default: config.ManagementGroup},
			Response: &config.ManagementGroup,
		},
		{
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
			Prompt:   &survey.Input{Message: QuestionADApplicationPass, Default: config.AdApplicationPassword},
			Required: true,
			Response: &config.AdApplicationPassword,
		},
		{
			Prompt:   &survey.Input{Message: QuestionADApplicationId, Default: config.AdApplicationId},
			Required: true,
			Response: &config.AdApplicationId,
		},
		{
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

func promptCustomizeAzureLoggingRegion(config *azure.GenerateAzureTfConfigurationArgs) error {
	var region string
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionStorageLocation, Default: config.StorageLocation},
			Response: &region,
		},
	}); err != nil {
		return err
	}
	if err := validateStorageLocation(region); err != nil {
		return err
	}
	config.StorageLocation = region
	return nil
}

func askAdvancedAzureOptions(config *azure.GenerateAzureTfConfigurationArgs, extraState *AzureGenerateCommandExtraState) error {
	answer := ""

	// Prompt for options
	for answer != AzureAdvancedOptDone {

		// Set the initial options
		options := []string{AzureUserIntegrationNames, AzureSubscriptions}
		// Only ask about Active Directory information if one was requested to be created
		if !config.CreateAdIntegration {
			options = append(options, AdvancedAdIntegration)
		}

		// Only show Region Storage options in the case of Activity Log integration
		if config.ActivityLog {
			options = append(options, AzureRegionStorage)
			options = append(options, AzureExistingStorageAcount)
		}

		// Only show management options in the case of Config integration
		if config.Config {
			options = append(options, AzureManagmentGroup)
		}

		options = append(options, AzureAdvancedOptLocation, AzureAdvancedOptDone)
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt: &survey.Select{
				Message: "Which options would you like to enable?",
				Options: options,
			},
			Response: &answer,
		}); err != nil {
			return err
		}

		// Based on response, prompt for actions
		switch answer {
		case AzureUserIntegrationNames:
			if err := promptAzureIntegrationNameQuestions(config); err != nil {
				return err
			}
		case AzureExistingStorageAcount:
			if err := promptAzureStorageAccountQuestions(config); err != nil {
				return err
			}
		case AzureSubscriptions:
			if err := promptAzureSubscriptionQuestions(config); err != nil {
				return err
			}
		case AzureManagmentGroup:
			if err := promptAzureManagementGroupQuestions(config); err != nil {
				return err
			}
		case AdvancedAdIntegration:
			if err := promptAzureAdIntegrationQuestions(config); err != nil {
				return err
			}
		case AzureRegionStorage:
			if err := promptCustomizeAzureLoggingRegion(config); err != nil {
				return err
			}
		case AzureAdvancedOptLocation:
			if err := promptCustomizeAzureOutputLocation(extraState); err != nil {
				return err
			}
			return nil
		}

		// Re-prompt if not done
		innerAskAgain := true
		if answer == AzureAdvancedOptDone {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: QuestionAzureAnotherAdvancedOpt, Default: false},
			Response: &innerAskAgain,
		}); err != nil {
			return err
		}

		if !innerAskAgain {
			answer = AzureAdvancedOptDone
		}
	}

	return nil
}

func azureConfigIsEmpty(config *azure.GenerateAzureTfConfigurationArgs) bool {
	return !config.Config && !config.ActivityLog
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
func promptAzureGenerate(config *azure.GenerateAzureTfConfigurationArgs, extraState *AzureGenerateCommandExtraState) error {

	// Cache for later use if generation is abandoned and in interactive mode
	if cli.InteractiveMode() {
		defer writeAzureGenerationArgsCache(config)
		defer extraState.writeCache()
	}
	// These are the core questions that should be asked.
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Confirm{Message: QuestionAzureEnableConfig, Default: false},
				Response: &config.Config,
			},
			{
				Prompt:   &survey.Confirm{Message: QuestionEnableActivityLog, Default: false},
				Response: &config.ActivityLog,
			},
			{
				Prompt:   &survey.Confirm{Message: QuestionEnableAdIntegration, Default: false},
				Response: &config.CreateAdIntegration,
			},
		}); err != nil {
		return err
	}

	// Validate one of config or activity log was enabled; otherwise error out
	if !config.Config && !config.ActivityLog {
		return errors.New("must enable activity log or config")
	}

	// Find out if the customer wants to specify more advanced features
	askAdvanced := false
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: QuestionAzureConfigAdvanced, Default: askAdvanced},
		Response: &askAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if askAdvanced {
		if err := askAdvancedAzureOptions(config, extraState); err != nil {
			return err
		}
	}

	return nil
}
