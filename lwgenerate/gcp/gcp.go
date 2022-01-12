package gcp

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/lwgenerate"
	"github.com/pkg/errors"
)

type ExistingServiceAccountDetails struct {
	// Existing Service Account Name
	Name string

	// Existing Service Account private key in JSON format, base64 encoded
	Base64EncodedPrivateKey string
}

// NewExistingServiceAccountDetails Create new existing Service Account details
func NewExistingServiceAccountDetails(name string, base64EncodedPrivateKey string) *ExistingServiceAccountDetails {
	return &ExistingServiceAccountDetails{
		Name:                    name,
		Base64EncodedPrivateKey: base64EncodedPrivateKey,
	}
}

type GenerateGcpTfConfigurationArgs struct {
	// Should we configure AuditLog integration in LW?
	AuditLog bool

	// Should we configure CSPM integration in LW?
	Config bool

	// Path to service account credentials to be used by Terraform
	ServiceAccountCredentials string

	// Should we configure an Organization wide integration?
	OrganizationIntegration bool

	// Supply a GCP Organization ID, only asked if OrganizationIntegration is True
	GcpOrganizationId string

	// Supply a GCP Project ID, to host the new resources
	GcpProjectId string

	// Optionally supply existing Service Account Details
	ExistingServiceAccount *ExistingServiceAccountDetails

	// If Config is true, give the user the opportunity to name their integration. Defaults to "TF Config"
	ConfigIntegrationName string

	// Set of labels which will be added to the resources managed by the module
	AuditLogLabels map[string]string

	// Set of labels which will be added to the audit log bucket
	BucketLabels map[string]string

	// Set of labels which will be added to the subscription
	PubSubSubscriptionLabels map[string]string

	// Set of labels which will be added to the topic
	PubSubTopicLabels map[string]string

	// Supply a GCP region for the new bucket. EU/US/ASIA
	BucketRegion string

	// Supply a GCP location for the new bucket. Defaults to global
	BucketLocation string

	// Supply a name for the new bucket
	BucketName string

	// Existing Bucket Name
	ExistingLogBucketName string

	// Existing Sink Name
	ExistingLogSinkName string

	// Should we force destroy the bucket if it has stuff in it? (only relevant on new AuditLog creation)
	EnableForceDestroyBucket bool

	// Boolean for enabling Uniform Bucket Level Access on the audit log bucket. Defaults to False
	EnableUBLA bool

	// Number of days to keep audit logs in Lacework GCS bucket before deleting.
	// Defaults to -1. Leave default to keep indefinitely.
	LogBucketLifecycleRuleAge int

	// The number of days to keep logs before deleting. Default is 30
	LogBucketRetentionDays int

	// If AuditLog is true, give the user the opportunity to name their integration. Defaults to "TF audit_log"
	AuditLogIntegrationName string

	// Lacework Profile to use
	LaceworkProfile string
}

// Ensure all combinations of inputs are valid for supported spec
func (args *GenerateGcpTfConfigurationArgs) validate() error {
	// Validate one of config or audit log was enabled; otherwise error out
	if !args.AuditLog && !args.Config {
		return errors.New("audit log or config integration must be enabled")
	}

	// Validate if this is an organization integration, verify that the organization id has been provided
	if args.OrganizationIntegration && args.GcpOrganizationId == "" {
		return errors.New("An Organization ID must be provided for an Organization Integration")
	}

	// Validate existing Service Account values, if set
	if args.ExistingServiceAccount != nil {
		if args.ExistingServiceAccount.Name == "" ||
			args.ExistingServiceAccount.Base64EncodedPrivateKey == "" {
			return errors.New("When using an existing Service Account, existing name, and base64 encoded JSON Private Key fields all must be set")
		}
	}

	return nil
}

type GcpTerraformModifier func(c *GenerateGcpTfConfigurationArgs)

// NewTerraform returns an instance of the GenerateGcpTfConfigurationArgs struct with the provided enabled
// settings (config/audit log).
//
// Note: Additional configuration details may be set using modifiers of the GcpTerraformModifier type
//
// Basic usage: Initialize a new GcpTerraformModifier struct, with GCP service account credentials. Then use generate to
//              create a string output of the required HCL.
//
//   hcl, err := gcp.NewTerraform(true, true,
//     gcp.WithGcpServiceAccountCredentials("/path/to/sa/credentials.json")).Generate()
//
func NewTerraform(enableConfig bool, enableAuditLog bool, mods ...GcpTerraformModifier) *GenerateGcpTfConfigurationArgs {
	config := &GenerateGcpTfConfigurationArgs{AuditLog: enableAuditLog, Config: enableConfig}
	for _, m := range mods {
		m(config)
	}
	return config
}

// WithGcpServiceAccountCredentials Set the path for the GCP Service Account to be utilized by the GCP provider
func WithGcpServiceAccountCredentials(path string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ServiceAccountCredentials = path
	}
}

// WithLaceworkProfile Set the Lacework Profile to utilize when integrating
func WithLaceworkProfile(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.LaceworkProfile = name
	}
}

// WithOrganizationId Set the Lacework organization ID to integrate with for an organization integration
func WithOrganizationId(id string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.GcpOrganizationId = id
	}
}

// WithProjectId Set the Lacework project ID that new resources should be created in
// (required for both project & org integration)
func WithProjectId(id string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.GcpProjectId = id
	}
}

// UseExistingServiceAccount Set an existing Service Account to be used by the Lacework Integration
func UseExistingServiceAccount(serviceAccountDetails *ExistingServiceAccountDetails) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExistingServiceAccount = serviceAccountDetails
	}
}

// UseConfigIntegrationName Set the Config Integration name to be displayed on the Lacework UI
func UseConfigIntegrationName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ConfigIntegrationName = name
	}
}

// AuditLogLabels set labels to be applied to ALL newly created AuditLog resources
func AuditLogLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.AuditLogLabels = labels
	}
}

// BucketLabels set labels to be applied to the newly created AuditLog Bucket
func BucketLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketLabels = labels
	}
}

// PubSubSubscriptionLabels set labels to be applied to the newly created AuditLog PubSub
func PubSubSubscriptionLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.PubSubSubscriptionLabels = labels
	}
}

// PubSubTopicLabels set labels to be applied to the newly created AuditLog PubSub Topic
func PubSubTopicLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.PubSubTopicLabels = labels
	}
}

// BucketRegion Set the Region in which the Bucket should be created
func BucketRegion(region string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketRegion = region
	}
}

// BucketLocation Set the name of the bucket that will receive log objects
func BucketLocation(location string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketLocation = location
	}
}

// BucketName Set the Location in which the Bucket should be created
func BucketName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketName = name
	}
}

// ExistingLogBucketName Set the bucket Name of an existing AuditLog Bucket setup
func ExistingLogBucketName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExistingLogBucketName = name
	}
}

// ExistingLogSinkName Set the Topic ARN of an existing AuditLog setup
func ExistingLogSinkName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExistingLogSinkName = name
	}
}

// EnableForceDestroyBucket Enable force destroy of the bucket if it has stuff in it
func EnableForceDestroyBucket() GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.EnableForceDestroyBucket = true
	}
}

// EnableUBLA Enable force destroy of the bucket if it has stuff in it
func EnableUBLA() GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.EnableUBLA = true
	}
}

// LogBucketLifecycleRuleAge Set the number of days to keep audit logs in Lacework GCS bucket before deleting
// Defaults to -1. Leave default to keep indefinitely.
func LogBucketLifecycleRuleAge(age int) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.LogBucketLifecycleRuleAge = age
	}
}

// LogBucketRetentionDays Set the number of days to keep logs before deleting. Default is 30
func LogBucketRetentionDays(days int) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.LogBucketRetentionDays = days
	}
}

// UseAuditLogIntegrationName Set the Config Integration name to be displayed on the Lacework UI
func UseAuditLogIntegrationName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.AuditLogIntegrationName = name
	}
}

// Generate new Terraform code based on the supplied args.
func (args *GenerateGcpTfConfigurationArgs) Generate() (string, error) {
	// Validate inputs
	if err := args.validate(); err != nil {
		return "", errors.Wrap(err, "invalid inputs")
	}

	// Create blocks
	requiredProviders, err := createRequiredProviders()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate required providers")
	}

	gcpProvider, err := createGcpProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate gcp provider")
	}

	laceworkProvider, err := createLaceworkProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework provider")
	}

	configModule, err := createConfig(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate aws config module")
	}

	auditLogModule, err := createAuditLog(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate gcp audit log module")
	}

	// Render
	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			gcpProvider,
			laceworkProvider,
			configModule,
			auditLogModule),
	)
	return hclBlocks, nil
}

func createRequiredProviders() (*hclwrite.Block, error) {
	return lwgenerate.CreateRequiredProviders(
		lwgenerate.NewRequiredProvider("lacework",
			lwgenerate.HclRequiredProviderWithSource("lacework/lacework"),
			lwgenerate.HclRequiredProviderWithVersion("~> 0.12.2")))
}

func createGcpProvider(args *GenerateGcpTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	if args.ServiceAccountCredentials != "" || args.GcpProjectId != "" {
		attrs := map[string]interface{}{}

		if args.ServiceAccountCredentials != "" {
			attrs["credentials"] = args.ServiceAccountCredentials
		}

		if args.GcpProjectId != "" {
			attrs["project"] = args.GcpProjectId
		}

		provider, err := lwgenerate.NewProvider("google", lwgenerate.HclProviderWithAttributes(attrs)).ToBlock()
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, provider)
	}

	return blocks, nil
}

func createLaceworkProvider(args *GenerateGcpTfConfigurationArgs) (*hclwrite.Block, error) {
	if args.LaceworkProfile != "" {
		return lwgenerate.NewProvider("lacework",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{"profile": args.LaceworkProfile}),
		).ToBlock()
	}
	return nil, nil
}

func createConfig(args *GenerateGcpTfConfigurationArgs) ([]*hclwrite.Block, error) {
	source := "lacework/config/gcp"
	version := "~> 1.0"

	blocks := []*hclwrite.Block{}
	if args.Config {
		attributes := map[string]interface{}{}
		moduleDetails := []lwgenerate.HclModuleModifier{}

		// default to using the project level module
		configModuleName := "gcp_project_level_config"
		if args.OrganizationIntegration {
			// if organization integration is true, override configModuleName to use the organization level module
			configModuleName = "gcp_organization_level_config"
			attributes["org_integration"] = args.OrganizationIntegration
			attributes["organization_id"] = args.GcpOrganizationId
		}

		if args.ExistingServiceAccount != nil {
			attributes["use_existing_service_account"] = true
			attributes["service_account_name"] = args.ExistingServiceAccount.Name
			attributes["service_account_private_key"] = args.ExistingServiceAccount.Base64EncodedPrivateKey
		}

		if args.ConfigIntegrationName != "" {
			attributes["lacework_integration_name"] = args.ConfigIntegrationName
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		moduleBlock, err := lwgenerate.NewModule(configModuleName, source,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(version))...).ToBlock()

		if err != nil {
			return nil, err
		}
		blocks = append(blocks, moduleBlock)

	}

	return blocks, nil
}

func createAuditLog(args *GenerateGcpTfConfigurationArgs) (*hclwrite.Block, error) {
	source := "lacework/config/gcp"
	version := "~> 2.0"

	if args.AuditLog {
		attributes := map[string]interface{}{}

		moduleDetails := []lwgenerate.HclModuleModifier{}

		if args.PubSubSubscriptionLabels != nil {
			attributes["pubsub_subscription_labels"] = args.PubSubSubscriptionLabels
		}

		if args.PubSubTopicLabels != nil {
			attributes["pubsub_topic_labels"] = args.BucketLabels
		}

		if args.ExistingLogBucketName != "" {
			attributes["existing_bucket_name"] = args.ExistingLogBucketName
		} else {
			attributes["lifecycle_rule_age"] = args.LogBucketLifecycleRuleAge

			if args.BucketName != "" {
				attributes["log_bucket"] = args.BucketName
			}

			if args.BucketLabels != nil {
				attributes["labels"] = args.BucketLabels
			}

			if args.EnableForceDestroyBucket {
				attributes["bucket_force_destroy"] = true
			}

			if args.EnableUBLA {
				attributes["enable_ubla"] = true
			}

			if args.BucketRegion != "" {
				attributes["bucket_region"] = args.BucketRegion
			}

			if args.BucketLocation != "" {
				attributes["log_bucket_location"] = args.BucketLocation
			}

			if args.LogBucketRetentionDays != 0 {
				attributes["log_bucket_retention_days"] = args.LogBucketRetentionDays
			}
		}

		if args.ExistingLogSinkName != "" {
			attributes["existing_sink_name"] = args.ExistingLogSinkName
		}

		// default to using the project level module
		auditLogModuleName := "gcp_project_audit_log"
		// default to using the project level module
		configModuleName := "gcp_project_level_config"
		if args.OrganizationIntegration {
			// if organization integration is true, override configModuleName to use the organization level module
			configModuleName = "gcp_organization_level_config"
			auditLogModuleName = "gcp_organization_level_audit_log"
			attributes["org_integration"] = args.OrganizationIntegration
			attributes["organization_id"] = args.GcpOrganizationId
		}

		if args.ExistingServiceAccount == nil && args.Config {
			attributes["use_existing_service_account"] = true
			attributes["service_account_name"] = lwgenerate.CreateSimpleTraversal([]string{"module", configModuleName, "service_account_name"})
			attributes["service_account_private_key"] = lwgenerate.CreateSimpleTraversal([]string{"module", configModuleName, "service_account_private_key"})
		}

		if args.ExistingServiceAccount != nil {
			attributes["use_existing_service_account"] = true
			attributes["service_account_name"] = args.ExistingServiceAccount.Name
			attributes["service_account_private_key"] = args.ExistingServiceAccount.Base64EncodedPrivateKey
		}

		if args.AuditLogIntegrationName != "" {
			attributes["lacework_integration_name"] = args.AuditLogIntegrationName
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		return lwgenerate.NewModule(auditLogModuleName, source,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(version))...).ToBlock()
	}

	return nil, nil
}
