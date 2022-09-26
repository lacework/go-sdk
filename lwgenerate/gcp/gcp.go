package gcp

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/internal/array"
	"github.com/lacework/go-sdk/internal/unique"
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

func (e *ExistingServiceAccountDetails) IsPartial() bool {
	// If nil, return false
	if e == nil {
		return false
	}

	// If all values are empty, return false
	if e.Name == "" && e.PrivateKey == "" {
		return false
	}

	// If all values are populated, return false
	if e.Name != "" && e.PrivateKey != "" {
		return false
	}

	return true
}

type GenerateGcpTfConfigurationArgs struct {
	// Should we configure AuditLog integration in LW?
	AuditLog bool

	// Should we configure CSPM integration in LW?
	Configuration bool

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

	// If Configuration is true, give the user the opportunity to name their integration. Defaults to "TF Config"
	ConfigurationIntegrationName string

	// Set of labels which will be added to the resources managed by the module
	AuditLogLabels map[string]string

	// Set of labels which will be added to the audit log bucket
	BucketLabels map[string]string

	// Set of labels which will be added to the subscription
	PubSubSubscriptionLabels map[string]string

	// Set of labels which will be added to the topic
	PubSubTopicLabels map[string]string

	CustomBucketName string

	// Supply a GCP region for the new bucket. EU/US/ASIA
	BucketRegion string

	// Existing Bucket Name
	ExistingLogBucketName string

	// Existing Sink Name
	ExistingLogSinkName string

	// Should we force destroy the bucket if it has stuff in it? (only relevant on new Audit Log creation)
	EnableForceDestroyBucket bool

	// Boolean for enabling Uniform Bucket Level Access on the audit log bucket. Defaults to False
	EnableUBLA bool

	// Number of days to keep audit logs in Lacework GCS bucket before deleting.
	// If left empty the TF will default to -1
	LogBucketLifecycleRuleAge int

	// If AuditLog is true, give the user the opportunity to name their integration. Defaults to "TF audit_log"
	AuditLogIntegrationName string

	// Lacework Profile to use
	LaceworkProfile string

	FoldersToInclude []string

	FoldersToExclude []string

	IncludeRootProjects bool

	CustomFilter string

	GoogleWorkspaceFilter bool

	K8sFilter bool

	Prefix string

	WaitTime string
}

// Ensure all combinations of inputs are valid for supported spec
func (args *GenerateGcpTfConfigurationArgs) validate() error {
	// Validate one of config or audit log was enabled; otherwise error out
	if !args.AuditLog && !args.Configuration {
		return errors.New("audit log or configuration integration must be enabled")
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
// settings (configuration/audit log).
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
	config := &GenerateGcpTfConfigurationArgs{
		AuditLog:              enableAuditLog,
		Configuration:         enableConfig,
		IncludeRootProjects:   true,
		EnableUBLA:            true,
		GoogleWorkspaceFilter: true,
		K8sFilter:             true,
	}
	// default LogBucketLifecycleRuleAge to -1. This helps us determine if the var has been set by the end user
	config.LogBucketLifecycleRuleAge = -1
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

// WithConfigurationIntegrationName Set the Config Integration name to be displayed on the Lacework UI
func WithConfigurationIntegrationName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ConfigurationIntegrationName = name
	}
}

// WithAuditLogLabels set labels to be applied to ALL newly created Audit Log resources
func WithAuditLogLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.AuditLogLabels = labels
	}
}

// WithBucketLabels set labels to be applied to the newly created Audit Log Bucket
func WithBucketLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketLabels = labels
	}
}

// WithPubSubSubscriptionLabels set labels to be applied to the newly created Audit Log PubSub
func WithPubSubSubscriptionLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.PubSubSubscriptionLabels = labels
	}
}

// WithPubSubTopicLabels set labels to be applied to the newly created Audit Log PubSub Topic
func WithPubSubTopicLabels(labels map[string]string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.PubSubTopicLabels = labels
	}
}

func WithCustomBucketName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.CustomBucketName = name
	}
}

// WithBucketRegion Set the Region in which the Bucket should be created
func WithBucketRegion(region string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.BucketRegion = region
	}
}

// WithExistingLogBucketName Set the bucket Name of an existing Audit Log Bucket setup
func WithExistingLogBucketName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExistingLogBucketName = name
	}
}

// WithExistingLogSinkName Set the Topic ARN of an existing Audit Log setup
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
func WithEnableUBLA(enable bool) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.EnableUBLA = enable
	}
}

// WithLogBucketLifecycleRuleAge Set the number of days to keep audit logs in Lacework GCS bucket before deleting
// Defaults to -1. Leave default to keep indefinitely.
func WithLogBucketLifecycleRuleAge(ruleAge int) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.LogBucketLifecycleRuleAge = ruleAge
	}
}

// WithAuditLogIntegrationName Set the Config Integration name to be displayed on the Lacework UI
func WithAuditLogIntegrationName(name string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.AuditLogIntegrationName = name
	}
}

func WithFoldersToInclude(folders []string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.FoldersToInclude = folders
	}
}

func WithFoldersToExclude(folders []string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.FoldersToExclude = folders
	}
}

func WithIncludeRootProjects(include bool) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.IncludeRootProjects = include
	}
}

func WithCustomFilter(filter string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.CustomFilter = filter
	}
}

func WithGoogleWorkspaceFilter(filter bool) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.GoogleWorkspaceFilter = filter
	}
}

func WithK8sFilter(filter bool) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.K8sFilter = filter
	}
}

func WithPrefix(prefix string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.Prefix = prefix
	}
}

func WithWaitTime(waitTime string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.WaitTime = waitTime
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

	configurationModule, err := createConfiguration(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate gcp configuration module")
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
			configurationModule,
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

func createConfiguration(args *GenerateGcpTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	if args.Configuration {
		attributes := map[string]interface{}{}
		moduleDetails := []lwgenerate.HclModuleModifier{}

		// default to using the project level module
		configurationModuleName := "gcp_project_level_config"
		if args.OrganizationIntegration {
			// if organization integration is true, override configModuleName to use the organization level module
			configurationModuleName = "gcp_organization_level_config"
			attributes["org_integration"] = args.OrganizationIntegration
			attributes["organization_id"] = args.GcpOrganizationId

			if len(args.FoldersToInclude) > 0 {
				attributes["folders_to_include"] = array.SortStrings(unique.StringSlice(args.FoldersToInclude))
			}

			if len(args.FoldersToExclude) > 0 {
				attributes["folders_to_exclude"] = array.SortStrings(unique.StringSlice(args.FoldersToExclude))

				// Default true in gcp-audit-log TF module
				if args.IncludeRootProjects != true {
					attributes["include_root_projects"] = args.IncludeRootProjects
				}
			}
		}

		if args.ExistingServiceAccount != nil {
			attributes["use_existing_service_account"] = true
			attributes["service_account_name"] = args.ExistingServiceAccount.Name
			attributes["service_account_private_key"] = args.ExistingServiceAccount.PrivateKey
		}

		if args.ConfigurationIntegrationName != "" {
			attributes["lacework_integration_name"] = args.ConfigurationIntegrationName
		}

		if args.Prefix != "" {
			attributes["prefix"] = args.Prefix
		}

		if args.WaitTime != "" {
			attributes["wait_time"] = args.WaitTime
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		moduleBlock, err := lwgenerate.NewModule(
			configurationModuleName,
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
			if args.LogBucketLifecycleRuleAge != -1 {
				attributes["lifecycle_rule_age"] = args.LogBucketLifecycleRuleAge
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

			// Default true in gcp-audit-log TF module
			if args.EnableUBLA != true {
				attributes["enable_ubla"] = args.EnableUBLA
			}

			if args.CustomBucketName != "" {
				attributes["custom_bucket_name"] = args.CustomBucketName
			}

			if args.BucketRegion != "" {
				attributes["bucket_region"] = args.BucketRegion
			}
		}

		if args.ExistingLogSinkName != "" {
			attributes["existing_sink_name"] = args.ExistingLogSinkName
		}

		// default to using the project level module
		auditLogModuleName := "gcp_project_audit_log"
		// default to using the project level module
		configurationModuleName := "gcp_project_level_config"
		if args.OrganizationIntegration {
			// if organization integration is true, override configModuleName to use the organization level module
			configurationModuleName = "gcp_organization_level_config"
			auditLogModuleName = "gcp_organization_level_audit_log"
			attributes["org_integration"] = args.OrganizationIntegration
			attributes["organization_id"] = args.GcpOrganizationId

			if len(args.FoldersToInclude) > 0 {
				attributes["folders_to_include"] = array.SortStrings(unique.StringSlice(args.FoldersToInclude))
			}

			if len(args.FoldersToExclude) > 0 {
				attributes["folders_to_exclude"] = array.SortStrings(unique.StringSlice(args.FoldersToExclude))

				// Default true in gcp-audit-log TF module
				if args.IncludeRootProjects != true {
					attributes["include_root_projects"] = args.IncludeRootProjects
				}
			}
		}

		if args.ExistingServiceAccount == nil && args.Configuration {
			attributes["use_existing_service_account"] = true
			attributes["service_account_name"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", configurationModuleName, "service_account_name"},
			)
			attributes["service_account_private_key"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", configurationModuleName, "service_account_private_key"},
			)
		}

		if args.ExistingServiceAccount != nil {
			attributes["use_existing_service_account"] = true
			attributes["service_account_name"] = args.ExistingServiceAccount.Name
			attributes["service_account_private_key"] = args.ExistingServiceAccount.PrivateKey
		}

		if args.AuditLogIntegrationName != "" {
			attributes["lacework_integration_name"] = args.AuditLogIntegrationName
		}

		if args.CustomFilter != "" {
			attributes["custom_filter"] = args.CustomFilter
		}

		// Default true in gcp-audit-log TF module
		if args.GoogleWorkspaceFilter != true {
			attributes["google_workspace_filter"] = args.GoogleWorkspaceFilter
		}

		// Default true in gcp-audit-log TF module
		if args.K8sFilter != true {
			attributes["k8s_filter"] = args.K8sFilter
		}

		if args.Prefix != "" {
			attributes["prefix"] = args.Prefix
		}

		if args.WaitTime != "" {
			attributes["wait_time"] = args.WaitTime
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
