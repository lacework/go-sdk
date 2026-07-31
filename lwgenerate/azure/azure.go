// A package that generates Lacework deployment code for Azure cloud.
package azure

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/v2/lwgenerate"
	"github.com/pkg/errors"
)

type GenerateAzureTfConfigurationArgs struct {
	// Should we configure Activity Log integration in LW?
	ActivityLog bool

	// Should we add Config integration in LW?
	Config bool

	// Should we add Agentless integration in LW?
	Agentless bool

	// Should we create an Entra ID integration in LW?
	EntraIdActivityLog bool

	// Should we create an Active Directory integration
	CreateAdIntegration bool

	// If Config is true, give the user the opportunity to name their integration. Defaults to "TF Config"
	ConfigIntegrationName string

	// If ActivityLog is true, give the user the opportunity to name their integration. Defaults to "TF activity log"
	ActivityLogIntegrationName string

	// If EntraIdIntegration is true, give the user the opportunity to name their integration.
	// Defaults to "TF Entra ID activity log"
	EntraIdIntegrationName string

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

	// Azure region where the event hub for logging will reside
	EventHubLocation string

	// Number of partitions in the Event Hub for logging
	EventHubPartitionCount int

	// Add custom blocks to the root `terraform{}` block. Can be used for advanced configuration. Things like backend, etc
	ExtraBlocksRootTerraform []*hclwrite.Block

	// ExtraAZRMArguments allows adding more arguments to the provider block as needed (custom use cases)
	ExtraAZRMArguments map[string]interface{}

	// ExtraAZReadArguments allows adding more arguments to the provider block as needed (custom use cases)
	ExtraAZReadArguments map[string]interface{}

	// ExtraBlocks allows adding more hclwrite.Block to the root terraform document (advanced use cases)
	ExtraBlocks []*hclwrite.Block

	// Custom outputs
	CustomOutputs []lwgenerate.HclOutput

	// Integration level for agentless scanning (e.g., "SUBSCRIPTION", "TENANT")
	IntegrationLevel string

	// Should agentless scanning be global?
	Global bool

	// Should we create a Log Analytics Workspace for agentless scanning?
	CreateLogAnalyticsWorkspace bool

	// List of regions to deploy for agentless scanning
	Regions []string

	// List of subscription IDs for agentless scanning
	AgentlessSubscriptionIds []string

	// Should we use storage account network rules for activity log?
	UseStorageAccountNetworkRules bool

	// List of IP addresses to access storage account
	StorageAccountNetworkRuleIpRules []string
}

// check if given IP address is an IPv4 address
func IsIpv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

// Ensure all combinations of inputs are valid for supported spec
func (args *GenerateAzureTfConfigurationArgs) validate() error {
	// Validate one of config, agentless or activity log was enabled; otherwise error out
	if !args.ActivityLog && !args.Agentless && !args.Config && !args.EntraIdActivityLog {
		return errors.New("audit log, agentless or config integration must be enabled")
	}

	if (args.ActivityLog || args.Agentless || args.Config || args.EntraIdActivityLog) && args.SubscriptionID == "" {
		return errors.New("subscription_id must be provided")
	}

	// Validate that active directory settings are correct
	if !args.CreateAdIntegration && (args.Config || args.ActivityLog || args.EntraIdActivityLog) {
		if args.AdApplicationId == "" ||
			args.AdServicePrincipalId == "" || args.AdApplicationPassword == "" {
			return errors.New("Active directory details must be set")
		}
	}

	// Validate the Mangement Group
	if args.ManagementGroup && args.ManagementGroupId == "" {
		return errors.New("When Group Management is enabled, then Group Id must be configured")
	}

	// Validate Storage Account
	if args.ExistingStorageAccount && (args.StorageAccountName == "" || args.StorageAccountResourceGroup == "") {
		return errors.New("When using existing storage account, storage account details must be configured")
	}

	// Validate Agentless Scanning
	if args.Agentless {
		if args.IntegrationLevel == "" {
			return errors.New("integration_level must be set for Agentless Integration")
		}
		if args.IntegrationLevel == "SUBSCRIPTION" && len(args.AgentlessSubscriptionIds) == 0 {
			return errors.New("subscription_ids must be provided for Agentless Integration with SUBSCRIPTION integration level")
		}
	}

	return nil
}

type AzureTerraformModifier func(c *GenerateAzureTfConfigurationArgs)

// NewTerraform returns an instance of the GenerateAzureTfConfigurationArgs struct with the provided enabled
// settings (config/activity log).
//
// Note: Additional configuration details may be set using modifiers of the AzureTerraformModifier type
func NewTerraform(
	enableConfig bool, enableActivityLog bool, enableAgentless bool, enableEntraIdActivityLog, createAdIntegration bool,
	mods ...AzureTerraformModifier,
) *GenerateAzureTfConfigurationArgs {
	config := &GenerateAzureTfConfigurationArgs{
		ActivityLog:         enableActivityLog,
		Config:              enableConfig,
		Agentless:           enableAgentless,
		EntraIdActivityLog:  enableEntraIdActivityLog,
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

// WithExtraAZRMArguments enables adding additional arguments into the `azurerm` provider block
// this enables custom use cases
func WithExtraAZRMArguments(arguments map[string]interface{}) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ExtraAZRMArguments = arguments
	}
}

// WithExtraAZReadArguments enables adding additional arguments into the `azuread` provider block
// this enables custom use cases
func WithExtraAZReadArguments(arguments map[string]interface{}) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ExtraAZReadArguments = arguments
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

// WithEntraIdActivityLogIntegrationName Set the Entra ID Activity Log Integration name
// to be displayed on the Lacework UI
func WithEntraIdActivityLogIntegrationName(name string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.EntraIdIntegrationName = name
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

// WithAgentlessSubscriptionIds List of subscriptions for agentless scanning.
func WithAgentlessSubscriptionIds(agentlessSubscriptionIds []string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.AgentlessSubscriptionIds = agentlessSubscriptionIds
	}
}

func WithRegions(regions []string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.Regions = regions
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

// WithEventHubLocation The Azure region where the event hub for logging resides
func WithEventHubLocation(location string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.EventHubLocation = location
	}
}

// WitthEventHubPartitionCount The number of partitions in the Event Hub for logging
func WithEventHubPartitionCount(partitionCount int) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.EventHubPartitionCount = partitionCount
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

// WithGlobal sets the Global field for agentless scanning
func WithGlobal(global bool) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.Global = global
	}
}

// WithCreateLogAnalyticsWorkspace sets the CreateLogAnalyticsWorkspace field for agentless scanning
func WithCreateLogAnalyticsWorkspace(create bool) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.CreateLogAnalyticsWorkspace = create
	}
}

// WithIntegrationLevel sets the IntegrationLevel field for agentless scanning
func WithIntegrationLevel(level string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.IntegrationLevel = level
	}
}

// WithUseStorageAccountNetworkRules sets the UseStorageAccountNetworkRules field for activity log
func WithUseStorageAccountNetworkRules(useNetworkRules bool) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.UseStorageAccountNetworkRules = useNetworkRules
	}
}

// WithUseStorageAccountNetworkRuleIpRules sets the StorageAccountNetworkRuleIpRules field for activity log
func WithUseStorageAccountNetworkRuleIpRules(ipRules []string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.StorageAccountNetworkRuleIpRules = ipRules
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

	agentlessLogModule, err := createAgentless(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate azure agentless module")
	}

	entraIdActivityLogModule, err := createEntraIdActivityLog(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate azure Entra ID activity log module")
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
			agentlessLogModule,
			entraIdActivityLogModule,
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
	for k, v := range args.ExtraAZReadArguments {
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
	for k, v := range args.ExtraAZRMArguments {
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

		// Add storage account network rules if enabled
		if args.UseStorageAccountNetworkRules && !args.ExistingStorageAccount {
			attributes["use_storage_account_network_rules"] = args.UseStorageAccountNetworkRules

			if len(args.StorageAccountNetworkRuleIpRules) > 0 {
				attributes["storage_account_network_rule_ip_rules"] = args.StorageAccountNetworkRuleIpRules
			}
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

func createAgentless(args *GenerateAzureTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}

	if !args.Agentless {
		return blocks, nil
	}

	// Helper function to format region names for module naming
	formatRegionForModuleName := func(region string) string {
		return strings.ToLower(strings.ReplaceAll(region, " ", "_"))
	}

	// Determine regions to process
	regions := args.Regions
	if len(regions) == 0 {
		regions = []string{"West US"} // Default to West US if no regions specified
	}

	isTenant := strings.EqualFold(args.IntegrationLevel, "TENANT")
	isSubscription := strings.EqualFold(args.IntegrationLevel, "SUBSCRIPTION")

	var firstModuleName string
	for i, region := range regions {
		isFirstRegion := (i == 0)

		// Build module name
		var moduleName string
		if isTenant {
			moduleName = fmt.Sprintf("lacework_azure_agentless_scanning_tenant_%s", formatRegionForModuleName(region))
		} else {
			moduleName = fmt.Sprintf("lacework_azure_agentless_scanning_subscription_%s", formatRegionForModuleName(region))
		}

		// Build attributes
		attrs := map[string]interface{}{
			"integration_level":              args.IntegrationLevel,
			"region":                         region,
			"global":                         args.Global && isFirstRegion,
			"create_log_analytics_workspace": args.CreateLogAnalyticsWorkspace,
		}

		// set the first module name for global reference
		if isFirstRegion {
			firstModuleName = moduleName
		} else {
			attrs["global_module_reference"] = lwgenerate.CreateSimpleTraversal([]string{"module", firstModuleName})
		}

		if args.SubscriptionID != "" {
			attrs["scanning_subscription_id"] = args.SubscriptionID
		}
		if isSubscription && isFirstRegion && len(args.AgentlessSubscriptionIds) > 0 {
			subs := make([]string, len(args.AgentlessSubscriptionIds))
			for j, id := range args.AgentlessSubscriptionIds {
				subs[j] = fmt.Sprintf("/subscriptions/%s", id)
			}
			attrs["included_subscriptions"] = subs
		}

		// Create module details
		moduleDetails := []lwgenerate.HclModuleModifier{
			lwgenerate.HclModuleWithAttributes(attrs),
		}
		block, err := lwgenerate.NewModule(
			moduleName,
			lwgenerate.LWAzureAgentlessSource,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(lwgenerate.LWAzureAgentlessVersion))...,
		).ToBlock()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func createEntraIdActivityLog(args *GenerateAzureTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	if args.EntraIdActivityLog {
		attributes := map[string]interface{}{}
		moduleDetails := []lwgenerate.HclModuleModifier{}

		if args.EntraIdIntegrationName != "" {
			attributes["lacework_integration_name"] = args.EntraIdIntegrationName
		}

		// Check if we have created an Active Directory integration
		if args.CreateAdIntegration {
			attributes["use_existing_ad_application"] = false
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

		if args.EventHubLocation != "" {
			attributes["location"] = args.EventHubLocation
		}

		if args.EventHubPartitionCount > 0 {
			attributes["num_partitions"] = args.EventHubPartitionCount
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		moduleBlock, err := lwgenerate.NewModule(
			"microsoft-entra-id-activity-log",
			lwgenerate.LWAzureEntraIdActivityLogSource,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(lwgenerate.LWAzureEntraIdActivityLogVersion))...,
		).ToBlock()

		if err != nil {
			return nil, err
		}
		blocks = append(blocks, moduleBlock)
	}
	return blocks, nil
}
