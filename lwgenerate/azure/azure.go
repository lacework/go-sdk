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

	// If Config is true, give the user the opportunity to name their integration. Defaults to "TF Config"
	ConfigIntegrationName string

	// If ActivityLog is true, give the user the opportunity to name their integration. Defaults to "TF activity log"
	ActivityLogIntegrationName string
}

// Ensure all combinations of inputs are valid for supported spec
func (args *GenerateAzureTfConfigurationArgs) validate() error {
	// Validate one of config or activity log was enabled; otherwise error out
	if !args.ActivityLog && !args.Config {
		return errors.New("audit log or config integration must be enabled")
	}
	return nil
}

type AzureTerraformModifier func(c *GenerateAzureTfConfigurationArgs)

// NewTerraform returns an instance of the GenerateAzureTfConfigurationArgs struct with the provided enabled
// settings (config/activity log).
//
// Note: Additional configuration details may be set using modifiers of the AzureTerraformModifier type
//
func NewTerraform(enableConfig bool, enableActivityLog bool, mods ...AzureTerraformModifier) *GenerateAzureTfConfigurationArgs {
	config := &GenerateAzureTfConfigurationArgs{ActivityLog: enableActivityLog, Config: enableConfig}
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

// WithAuditLogIntegrationName Set the Config Integration name to be displayed on the Lacework UI
func WithAuditLogIntegrationName(name string) AzureTerraformModifier {
	return func(c *GenerateAzureTfConfigurationArgs) {
		c.ActivityLogIntegrationName = name
	}
}

// Generate new Terraform code based on the supplied args.
func (args *GenerateAzureTfConfigurationArgs) Generate() (string, error) {
	// Validate inputs
	if err := args.validate(); err != nil {
		return "", errors.Wrap(err, "invalid inputs")
	}

	// Create blocks
	requiredProviders, err := createRequiredProviders()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate required providers")
	}

	azureADProvider, err := createAzureADProvider()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate AD provider")
	}

	azureRMProvider, err := createAzureRMProvider()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate AM provider")
	}

	laceworkProvider, err := createLaceworkAZADModule()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework provider")
	}

	configModule, err := createConfig(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate azure config module")
	}

	activityLogModule, err := createActivityLog(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate azure activity log module")
	}

	// Render
	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			azureADProvider,
			azureRMProvider,
			laceworkProvider,
			configModule,
			activityLogModule),
	)
	return hclBlocks, nil
}

func createRequiredProviders() (*hclwrite.Block, error) {
	return lwgenerate.CreateRequiredProviders(
		lwgenerate.NewRequiredProvider(
			"azuread",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.HashAzureADProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(lwgenerate.HashAzureADProviderVersion),
		),
		lwgenerate.NewRequiredProvider(
			"azurerm",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.HashAzureRMProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(lwgenerate.HashAzureRMProviderVersion),
		),
	)
}

func createAzureADProvider() ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	//attrs := map[string]interface{}{}
	provider, err := lwgenerate.NewProvider(
		"azuread",
	).ToBlock()

	if err != nil {
		return nil, err
	}

	blocks = append(blocks, provider)
	return blocks, nil
}

func createAzureRMProvider() ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	attrs := map[string]interface{}{}

	provider, err := lwgenerate.NewProvider(
		"azurerm",
	).ToBlock()

	if err != nil {
		return nil, err
	}
	// Create the features block
	featuresBlock, err := lwgenerate.HclCreateGenericBlock("features", []string{}, attrs)
	provider.Body().AppendBlock(featuresBlock)

	if err != nil {
		return nil, err
	}

	blocks = append(blocks, provider)
	return blocks, nil
}

func createLaceworkAZADModule() ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}

	provider, err := lwgenerate.NewModule(
		"az_ad_application",
		lwgenerate.LWAzureADSource,
		lwgenerate.HclModuleWithVersion(lwgenerate.LWAzureADVersion),
	).ToBlock()

	if err != nil {
		return nil, err
	}

	blocks = append(blocks, provider)
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

		// Set default values
		attributes["use_existing_ad_application"] = true
		attributes["application_id"] = "module.az_ad_application.application_id"
		attributes["application_password"] = "module.az_ad_application.application_password"
		attributes["service_principal_id"] = "module.az_ad_application.service_principal_id"

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

		// Set default values
		attributes["use_existing_ad_application"] = true
		attributes["application_id"] = "module.az_ad_application.application_id"
		attributes["application_password"] = "module.az_ad_application.application_password"
		attributes["service_principal_id"] = "module.az_ad_application.service_principal_id"

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
