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
	PrivateKey string
}

// NewExistingServiceAccountDetails Create new existing Service Account details
func NewExistingServiceAccountDetails(name string, privateKey string) *ExistingServiceAccountDetails {
	return &ExistingServiceAccountDetails{
		Name:       name,
		PrivateKey: privateKey,
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
	// If left empty the TF will default to -1
	// Use pointer *int, so we can verify if the value has been set by the end user
	LogBucketLifecycleRuleAge *int

	// The number of days to keep logs before deleting.
	// If left as 0 the TF will default to 30.
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

	// Validate if an organization id has been provided that this is and organization integration
	if !args.OrganizationIntegration && args.GcpOrganizationId != "" {
		return errors.New("To provide an Organization ID, Organization Integration must be true")
	}

	// Validate existing Service Account values, if set
	if args.ExistingServiceAccount != nil {
		if args.ExistingServiceAccount.Name == "" ||
			args.ExistingServiceAccount.PrivateKey == "" {
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

// WithOrganizationIntegration Set whether we configure as an Organization wide integration
func WithOrganizationIntegration(enabled bool) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.OrganizationIntegration = enabled
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

// WithExistingServiceAccount Set an existing Service Account to be used by the Lacework Integration
func WithExistingServiceAccount(serviceAccountDetails *ExistingServiceAccountDetails) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExistingServiceAccount = serviceAccountDetails
	}
}

// WithConfigIntegrationName Set the Config Integration name to be displayed on the Lacework UI
func WithConfigIntegrationName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ConfigIntegrationName = name
	}
}

// WithAuditLogLabels set labels to be applied to ALL newly created AuditLog resources
func WithAuditLogLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.AuditLogLabels = labels
	}
}

// WithBucketLabels set labels to be applied to the newly created AuditLog Bucket
func WithBucketLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketLabels = labels
	}
}

// WithPubSubSubscriptionLabels set labels to be applied to the newly created AuditLog PubSub
func WithPubSubSubscriptionLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.PubSubSubscriptionLabels = labels
	}
}

// WithPubSubTopicLabels set labels to be applied to the newly created AuditLog PubSub Topic
func WithPubSubTopicLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.PubSubTopicLabels = labels
	}
}

// WithBucketRegion Set the Region in which the Bucket should be created
func WithBucketRegion(region string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketRegion = region
	}
}

// WithBucketLocation Set the name of the bucket that will receive log objects
func WithBucketLocation(location string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketLocation = location
	}
}

// WithBucketName Set the Location in which the Bucket should be created
func WithBucketName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketName = name
	}
}

// WithExistingLogBucketName Set the bucket Name of an existing AuditLog Bucket setup
func WithExistingLogBucketName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExistingLogBucketName = name
	}
}

// WithExistingLogSinkName Set the Topic ARN of an existing AuditLog setup
func WithExistingLogSinkName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExistingLogSinkName = name
	}
}

// WithEnableForceDestroyBucket Enable force destroy of the bucket if it has stuff in it
func WithEnableForceDestroyBucket() GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.EnableForceDestroyBucket = true
	}
}

// WithEnableUBLA Enable force destroy of the bucket if it has stuff in it
func WithEnableUBLA() GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.EnableUBLA = true
	}
}

// WithLogBucketLifecycleRuleAge Set the number of days to keep audit logs in Lacework GCS bucket before deleting
// Defaults to -1. Leave default to keep indefinitely.
func WithLogBucketLifecycleRuleAge(ruleAge int) GcpTerraformModifier {
	age := &ruleAge
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.LogBucketLifecycleRuleAge = age
	}
}

// WithLogBucketRetentionDays Set the number of days to keep logs before deleting. Default is 30
func WithLogBucketRetentionDays(days int) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.LogBucketRetentionDays = days
	}
}

// WithAuditLogIntegrationName Set the Config Integration name to be displayed on the Lacework UI
func WithAuditLogIntegrationName(name string) GcpTerraformModifier {
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
		return "", errors.Wrap(err, "failed to generate gcp config module")
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
		lwgenerate.NewRequiredProvider(
			"lacework",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.LaceworkProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(lwgenerate.LaceworkProviderVersion),
		),
	)
}

func createGcpProvider(args *GenerateGcpTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	attrs := map[string]interface{}{}

	if args.ServiceAccountCredentials != "" {
		attrs["credentials"] = args.ServiceAccountCredentials
	}

	if args.GcpProjectId != "" {
		attrs["project"] = args.GcpProjectId
	}

	provider, err := lwgenerate.NewProvider(
		"google",
		lwgenerate.HclProviderWithAttributes(attrs),
	).ToBlock()
	if err != nil {
		return nil, err
	}

	blocks = append(blocks, provider)

	return blocks, nil
}

func createLaceworkProvider(args *GenerateGcpTfConfigurationArgs) (*hclwrite.Block, error) {
	if args.LaceworkProfile != "" {
		return lwgenerate.NewProvider(
			"lacework",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{"profile": args.LaceworkProfile}),
		).ToBlock()
	}
	return nil, nil
}

func createConfig(args *GenerateGcpTfConfigurationArgs) ([]*hclwrite.Block, error) {
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
			attributes["service_account_private_key"] = args.ExistingServiceAccount.PrivateKey
		}

		if args.ConfigIntegrationName != "" {
			attributes["lacework_integration_name"] = args.ConfigIntegrationName
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		moduleBlock, err := lwgenerate.NewModule(
			configModuleName,
			lwgenerate.GcpConfigSource,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(lwgenerate.GcpConfigVersion))...,
		).ToBlock()

		if err != nil {
			return nil, err
		}
		blocks = append(blocks, moduleBlock)

	}

	return blocks, nil
}

func createAuditLog(args *GenerateGcpTfConfigurationArgs) (*hclwrite.Block, error) {
	if args.AuditLog {
		attributes := map[string]interface{}{}

		moduleDetails := []lwgenerate.HclModuleModifier{}

		if args.PubSubSubscriptionLabels != nil {
			attributes["pubsub_subscription_labels"] = args.PubSubSubscriptionLabels
		}

		if args.PubSubTopicLabels != nil {
			attributes["pubsub_topic_labels"] = args.PubSubTopicLabels
		}

		if args.ExistingLogBucketName != "" {
			attributes["existing_bucket_name"] = args.ExistingLogBucketName
		} else {
			if args.LogBucketLifecycleRuleAge == nil {
				defaultValue := -1
				args.LogBucketLifecycleRuleAge = &defaultValue
			}
			attributes["lifecycle_rule_age"] = *args.LogBucketLifecycleRuleAge

			if args.BucketName != "" {
				attributes["log_bucket"] = args.BucketName
			}

			if args.AuditLogLabels != nil {
				attributes["labels"] = args.AuditLogLabels
			}

			if args.BucketLabels != nil {
				attributes["bucket_labels"] = args.BucketLabels
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
			attributes["service_account_private_key"] = args.ExistingServiceAccount.PrivateKey
		}

		if args.AuditLogIntegrationName != "" {
			attributes["lacework_integration_name"] = args.AuditLogIntegrationName
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		return lwgenerate.NewModule(
			auditLogModuleName,
			lwgenerate.GcpAuditLogSource,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(lwgenerate.GcpAuditLogVersion))...,
		).ToBlock()
	}

	return nil, nil
}
