// A package that generates Lacework deployment code for Azure cloud.
package azure

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/lwgenerate"
	"github.com/pkg/errors"
)

type GenerateAzureTfConfigurationArgs struct {
	// Should we configure Activity Log integration in LW?
	ActivityLog bool

	// Should we add Config integration in LW?
	Config bool

	// Should we create an Active Directory integration
	CreateAdIntegration bool

	// If Config is true, give the user the opportunity to name their integration. Defaults to "TF Config"
	ConfigIntegrationName string

	// If ActivityLog is true, give the user the opportunity to name their integration. Defaults to "TF activity log"
	ActivityLogIntegrationName string

	// Active Directory application Id
	AdApplicationId string

	// Active Directory password
	AdApplicationPassword string

	// Active Directory Enterprise app object id
	AdServicePrincipalId string

	// Should we use the management group, rather than subscription
	ManagementGroup bool

	// Management Group ID to set
	ManagementGroupId string

	// List of subscription Ids
	SubscriptionIds []string

	// Subscription ID configured in azurerm provider block
	SubscriptionID string

	// Grant read access to ALL subscriptions
	AllSubscriptions bool

	// Storage Account name
	StorageAccountName string

	// Storage Account Resource Group
	StorageAccountResourceGroup string

	// Should we use existing storage account
	ExistingStorageAccount bool

	// Azure region where the storage account for logging resides
	StorageLocation string

	LaceworkProfile string

	// Add custom blocks to the root `terraform{}` block. Can be used for advanced configuration. Things like backend, etc
	ExtraBlocksRootTerraform []*hclwrite.Block

	// ExtraProviderArguments allows adding more arguments to the provider block as needed (custom use cases)
	ExtraProviderArguments map[string]interface{}

	// ExtraBlocks allows adding more hclwrite.Block to the root terraform document (advanced use cases)
	ExtraBlocks []*hclwrite.Block

	// Custom outputs
	CustomOutputs []lwgenerate.HclOutput
}

// Ensure all combinations of inputs are valid for supported spec
func (args *GenerateAzureTfConfigurationArgs) validate() error {
	// Validate one of config or activity log was enabled; otherwise error out
	if !args.ActivityLog && !args.Config {
		return errors.New("audit log or config integration must be enabled")
	}

	// Validate that active directory settings are correct
	if !args.CreateAdIntegration && (args.AdApplicationId == "" ||
		args.AdServicePrincipalId == "" || args.AdApplicationPassword == "") {
		return errors.New("Active directory details must be set")
	}

	// Validate the Mangement Group
	if args.ManagementGroup && args.ManagementGroupId == "" {
		return errors.New("When Group Management is enabled, then Group Id must be configured")
	}

	// Validate Storage Account
	if args.ExistingStorageAccount && (args.StorageAccountName == "" || args.StorageAccountResourceGroup == "") {
		return errors.New("When using existing storage account, storage account details must be configured")
	}

	return nil
}

type AzureTerraformModifier func(c *GenerateAzureTfConfigurationArgs)

// NewTerraform returns an instance of the GenerateAzureTfConfigurationArgs struct with the provided enabled
// settings (config/activity log).
//
// Note: Additional configuration details may be set using modifiers of the AzureTerraformModifier type
func NewTerraform(
	enableConfig bool, enableActivityLog bool, createAdIntegration bool, mods ...AzureTerraformModifier,
) *GenerateAzureTfConfigurationArgs {
	config := &GenerateAzureTfConfigurationArgs{
		ActivityLog:         enableActivityLog,
		Config:              enableConfig,
		CreateAdIntegration: createAdIntegration,
	}
	for _, m := range mods {
		m(config)
	}
	return config
}

// WithConfigIntegrationName Set the Config Integration name to be displayed on the Lacework UI
func WithConfigIntegrationName(name string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ConfigIntegrationName = name
	}
}

// WithConfigOutputs Set Custom Terraform Outputs
func WithCustomOutputs(outputs []lwgenerate.HclOutput) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.CustomOutputs = outputs
	}
}

// WithExtraRootBlocks allows adding generic hcl blocks to the root `terraform{}` block
// this enables custom use cases
func WithExtraRootBlocks(blocks []*hclwrite.Block) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ExtraBlocksRootTerraform = blocks
	}
}

// WithExtraProviderArguments enables adding additional arguments into the `aws` provider block
// this enables custom use cases
func WithExtraProviderArguments(arguments map[string]interface{}) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ExtraProviderArguments = arguments
	}
}

// WithExtraBlocks enables adding additional arbitrary blocks to the root hcl document
func WithExtraBlocks(blocks []*hclwrite.Block) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ExtraBlocks = blocks
	}
}

// WithActivityLogIntegrationName Set the Activity Log Integration name to be displayed on the Lacework UI
func WithActivityLogIntegrationName(name string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ActivityLogIntegrationName = name
	}
}

// WithAdApplicationId Set Active Directory application id
func WithAdApplicationId(AdApplicationId string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.AdApplicationId = AdApplicationId
	}
}

// WithAdApplicationPassword Set the Active Directory password
func WithAdApplicationPassword(AdApplicationPassword string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.AdApplicationPassword = AdApplicationPassword
	}
}

// WithAdServicePrincipalId Set Active Directory principal id
func WithAdServicePrincipalId(AdServicePrincipalId string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.AdServicePrincipalId = AdServicePrincipalId
	}
}

// WithManagementGroup Enable the Management Group to allow AD to be reader on management group
// rather then subscription
func WithManagementGroup(enableManagentGroup bool) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ManagementGroup = enableManagentGroup
	}
}

// WithManagementGroupId The Group Id to add reader permissions
func WithManagementGroupId(managementGroupId string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ManagementGroupId = managementGroupId
	}
}

// WithSubscriptionIds List of subscriptions to to enable logging
func WithSubscriptionIds(subscriptionIds []string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.SubscriptionIds = subscriptionIds
	}
}

// WithAllSubscriptions Grant read access to ALL subscriptions within
// the selected Tenant (overrides 'subscription_ids')
func WithAllSubscriptions(allSubscriptions bool) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.AllSubscriptions = allSubscriptions
	}
}

// WithExistingStorageAccount Use an existing Storage Account
func WithExistingStorageAccount(existingStorageAccount bool) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ExistingStorageAccount = existingStorageAccount
	}
}

// WithStorageAccountName The name of the Storage Account
func WithStorageAccountName(storageAccountName string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.StorageAccountName = storageAccountName
	}
}

// WithStorageAccountResourceGroup The Resource Group for the existing Storage Account
func WithStorageAccountResourceGroup(storageAccountResourceGroup string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.StorageAccountResourceGroup = storageAccountResourceGroup
	}
}

// WithStorageLocation The Azure region where storage account for logging is
func WithStorageLocation(location string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.StorageLocation = location
	}
}

func WithLaceworkProfile(name string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.LaceworkProfile = name
	}
}

func WithSubscriptionID(subcriptionID string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.SubscriptionID = subcriptionID
	}
}

// Generate new Terraform code based on the supplied args.
func (args *GenerateAzureTfConfigurationArgs) Generate() (string, error) {
	// Validate inputs
	if err := args.validate(); err != nil {
		return "", errors.Wrap(err, "invalid inputs")
	}

	// Create blocks
	requiredProviders, err := createRequiredProviders(args.ExtraBlocksRootTerraform)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate required providers")
	}

	laceworkProvider, err := createLaceworkProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework provider")
	}

	azureADProvider, err := createAzureADProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate AD provider")
	}

	azureRMProvider, err := createAzureRMProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate AM provider")
	}

	laceworkADProvider, err := createLaceworkAzureADModule(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework Azure AD provider")
	}

	configModule, err := createConfig(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate azure config module")
	}

	activityLogModule, err := createActivityLog(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate azure activity log module")
	}

	outputBlocks := []*hclwrite.Block{}
	for _, output := range args.CustomOutputs {
		outputBlock, err := output.ToBlock()
		if err != nil {
			return "", errors.Wrap(err, "failed to add custom output")
		}
		outputBlocks = append(outputBlocks, outputBlock)
	}

	// Render
	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			laceworkProvider,
			azureADProvider,
			azureRMProvider,
			laceworkADProvider,
			configModule,
			activityLogModule,
			outputBlocks,
			args.ExtraBlocks),
	)
	return hclBlocks, nil
}

func createRequiredProviders(extraBlocks []*hclwrite.Block) (*hclwrite.Block, error) {
	return lwgenerate.CreateRequiredProvidersWithCustomBlocks(
		extraBlocks,
		lwgenerate.NewRequiredProvider(
			"lacework",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.LaceworkProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(lwgenerate.LaceworkProviderVersion),
		),
	)
}

func createLaceworkProvider(args *GenerateAzureTfConfigurationArgs) (*hclwrite.Block, error) {
	if args.LaceworkProfile != "" {
		return lwgenerate.NewProvider(
			"lacework",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{"profile": args.LaceworkProfile}),
		).ToBlock()
	}
	return nil, nil
}

func createAzureADProvider(args *GenerateAzureTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	attrs := map[string]interface{}{}

	// set custom args before the required ones below to ensure expected behavior (i.e., no overrides)
	for k, v := range args.ExtraProviderArguments {
		attrs[k] = v
	}

	provider, err := lwgenerate.NewProvider(
		"azuread",
		lwgenerate.HclProviderWithAttributes(attrs),
	).ToBlock()

	if err != nil {
		return nil, err
	}

	blocks = append(blocks, provider)
	return blocks, nil
}

// In this we need to create a provider block with a  features
// configuration but with nothing set,  this is as per the
// Azure examples and is of the format
//
//	provider "azurerm" {
//	   features = {}
//	}
func createAzureRMProvider(args *GenerateAzureTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	attrs := map[string]interface{}{}
	featureAttrs := map[string]interface{}{}

	// set custom args before the required ones below to ensure expected behavior (i.e., no overrides)
	for k, v := range args.ExtraProviderArguments {
		attrs[k] = v
	}

	if args.SubscriptionID != "" {
		attrs["subscription_id"] = args.SubscriptionID
	}

	provider, err := lwgenerate.NewProvider(
		"azurerm",
		lwgenerate.HclProviderWithAttributes(attrs),
	).ToBlock()

	if err != nil {
		return nil, err
	}
	// Create the features block
	featuresBlock, err := lwgenerate.HclCreateGenericBlock("features", []string{}, featureAttrs)
	provider.Body().AppendBlock(featuresBlock)

	if err != nil {
		return nil, err
	}

	blocks = append(blocks, provider)
	return blocks, nil
}

func createLaceworkAzureADModule(args *GenerateAzureTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}

	if args.CreateAdIntegration {
		provider, err := lwgenerate.NewModule(
			"az_ad_application",
			lwgenerate.LWAzureADSource,
			lwgenerate.HclModuleWithVersion(lwgenerate.LWAzureADVersion),
		).ToBlock()

		if err != nil {
			return nil, err
		}

		blocks = append(blocks, provider)
	}
	return blocks, nil
}

func createConfig(args *GenerateAzureTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	if args.Config {
		attributes := map[string]interface{}{}
		moduleDetails := []lwgenerate.HclModuleModifier{}

		if args.ConfigIntegrationName != "" {
			attributes["lacework_integration_name"] = args.ConfigIntegrationName
		}

		// Check if we have created an Active Directory app
		if args.CreateAdIntegration {
			attributes["use_existing_ad_application"] = true
			attributes["application_id"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "az_ad_application", "application_id"})
			attributes["application_password"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "az_ad_application", "application_password"})
			attributes["service_principal_id"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "az_ad_application", "service_principal_id"})
		} else {
			attributes["use_existing_ad_application"] = true
			attributes["application_id"] = args.AdApplicationId
			attributes["application_password"] = args.AdApplicationPassword
			attributes["service_principal_id"] = args.AdServicePrincipalId
		}

		// Only set subscription ids if all subscriptions flag is not set
		if !args.AllSubscriptions {
			if len(args.SubscriptionIds) > 0 {
				attributes["subscription_ids"] = args.SubscriptionIds
			}
		} else {
			// Set Subscription information
			attributes["all_subscriptions"] = args.AllSubscriptions
		}

		// Set Management Group details
		if args.ManagementGroup {
			attributes["use_management_group"] = args.ManagementGroup
			attributes["management_group_id"] = args.ManagementGroupId
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		moduleBlock, err := lwgenerate.NewModule(
			"az_config",
			lwgenerate.LWAzureConfigSource,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(lwgenerate.LWAzureConfigVersion))...,
		).ToBlock()

		if err != nil {
			return nil, err
		}
		blocks = append(blocks, moduleBlock)
	}

	return blocks, nil
}

func createActivityLog(args *GenerateAzureTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	if args.ActivityLog {
		attributes := map[string]interface{}{}
		moduleDetails := []lwgenerate.HclModuleModifier{}

		if args.ActivityLogIntegrationName != "" {
			attributes["lacework_integration_name"] = args.ActivityLogIntegrationName
		}

		// Check if we have created an Active Directory integration
		if args.CreateAdIntegration {
			attributes["use_existing_ad_application"] = true
			attributes["application_id"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "az_ad_application", "application_id"})
			attributes["application_password"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "az_ad_application", "application_password"})
			attributes["service_principal_id"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "az_ad_application", "service_principal_id"})
		} else {
			attributes["use_existing_ad_application"] = true
			attributes["application_id"] = args.AdApplicationId
			attributes["application_password"] = args.AdApplicationPassword
			attributes["service_principal_id"] = args.AdServicePrincipalId
		}

		// Only set subscription ids if all subscriptions flag is not set
		if !args.AllSubscriptions {
			if len(args.SubscriptionIds) > 0 {
				attributes["subscription_ids"] = args.SubscriptionIds
			}
		} else {
			// Set Subscription information
			attributes["all_subscriptions"] = args.AllSubscriptions
		}

		// Set storage account name, if set
		if args.StorageAccountName != "" {
			attributes["storage_account_name"] = args.StorageAccountName
		}

		// Set storage info if existing storage flag is set
		if args.ExistingStorageAccount {
			attributes["use_existing_storage_account"] = args.ExistingStorageAccount
			attributes["storage_account_resource_group"] = args.StorageAccountResourceGroup
		}

		// if a new storage account is being created (i.e., ExistingStorageAccount is false), enable infrastructure
		// encryption
		if !args.ExistingStorageAccount {
			attributes["infrastructure_encryption_enabled"] = true
		}

		// Set the location if needed
		if args.StorageLocation != "" {
			attributes["location"] = args.StorageLocation
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		moduleBlock, err := lwgenerate.NewModule(
			"az_activity_log",
			lwgenerate.LWAzureActivityLogSource,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(lwgenerate.LWAzureActivityLogVersion))...,
		).ToBlock()

		if err != nil {
			return nil, err
		}
		blocks = append(blocks, moduleBlock)

	}
	return blocks, nil
}
