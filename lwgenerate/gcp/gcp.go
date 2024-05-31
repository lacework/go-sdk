// A package that generates Lacework deployment code for Google cloud.
package gcp

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
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
	// Should we configure Agentless integration in LW?
	Agentless bool

	// Should we configure AuditLog integration in LW?
	AuditLog bool

	// Should we use the Pub Sub Audit Log or use the Bucket based one
	UsePubSubAudit bool

	// Should we configure CSPM integration in LW?
	Configuration bool

	// A list of GCP project IDs to monitor for Agentless integration
	ProjectFilterList []string

	// A list of regions to deploy for Agentless integration
	Regions []string

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
	// DEPRECATED
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

	Projects []string

	// Default GCP Provider labels
	ProviderDefaultLabels map[string]interface{}

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
func (args *GenerateGcpTfConfigurationArgs) validate() error {
	// Validate one of agentless, config or audit log was enabled; otherwise error out
	if !args.Agentless && !args.AuditLog && !args.Configuration {
		return errors.New("agentless, audit log or configuration integration must be enabled")
	}

	if args.Agentless && len(args.Regions) == 0 {
		return errors.New("regions must be provided for Agentless Integration")
	}

	// Validate if this is an organization integration, verify that the organization id has been provided
	if args.OrganizationIntegration && args.GcpOrganizationId == "" {
		return errors.New("an Organization ID must be provided for an Organization Integration")
	}

	// Validate if an organization id has been provided that this is and organization integration
	if !args.OrganizationIntegration && args.GcpOrganizationId != "" {
		return errors.New("to provide an Organization ID, Organization Integration must be true")
	}

	// Validate existing Service Account values, if set
	if args.ExistingServiceAccount != nil {
		if args.ExistingServiceAccount.Name == "" ||
			args.ExistingServiceAccount.PrivateKey == "" {
			return errors.New("when using an existing Service Account, existing name, and base64 " +
				"encoded JSON Private Key fields all must be set")
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
//
//	           create a string output of the required HCL.
//
//	hcl, err := gcp.NewTerraform(true, true, true, true,
//	  gcp.WithGcpServiceAccountCredentials("/path/to/sa/credentials.json")).Generate()
func NewTerraform(
	enableAgentless, enableConfig bool, enableAuditLog bool, enablePubSubAudit bool, mods ...GcpTerraformModifier,
) *GenerateGcpTfConfigurationArgs {
	config := &GenerateGcpTfConfigurationArgs{
		Agentless:             enableAgentless,
		AuditLog:              enableAuditLog,
		UsePubSubAudit:        enablePubSubAudit,
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

// WithUsePubSubAudit Set wether we use pub sub with the audit log rather than bucket based
func WithUsePubSubAudit(usePubSub bool) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.UsePubSubAudit = usePubSub
	}
}

// WithGcpServiceAccountCredentials Set the path for the GCP Service Account to be utilized by the GCP provider
func WithGcpServiceAccountCredentials(path string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ServiceAccountCredentials = path
	}
}

// WithProviderDefaultLabels adds default_labels to the provider configuration for GCP (if labels are present)
func WithProviderDefaultLabels(labels map[string]interface{}) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ProviderDefaultLabels = labels
	}
}

// WithConfigOutputs Set Custom Terraform Outputs
func WithCustomOutputs(outputs []lwgenerate.HclOutput) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.CustomOutputs = outputs
	}
}

// WithExtraRootBlocks allows adding generic hcl blocks to the root `terraform{}` block
// this enables custom use cases
func WithExtraRootBlocks(blocks []*hclwrite.Block) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExtraBlocksRootTerraform = blocks
	}
}

// WithExtraProviderArguments enables adding additional arguments into the `gcp` provider block
// this enables custom use cases
func WithExtraProviderArguments(arguments map[string]interface{}) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExtraProviderArguments = arguments
	}
}

// WithExtraBlocks enables adding additional arbitrary blocks to the root hcl document
func WithExtraBlocks(blocks []*hclwrite.Block) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ExtraBlocks = blocks
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

func WithMultipleProject(projects []string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.Projects = projects
	}
}

func WithProjectFilterList(projectFilterList []string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.ProjectFilterList = projectFilterList
	}
}

func WithRegions(regions []string) GcpTerraformModifier {
	return func(c *GenerateGcpTfConfigurationArgs) {
		c.Regions = regions
	}
}

// Generate new Terraform code based on the supplied args.
func (args *GenerateGcpTfConfigurationArgs) Generate() (string, error) {
	// Validate inputs
	if err := args.validate(); err != nil {
		return "", errors.Wrap(err, "invalid inputs")
	}

	// Create blocks
	requiredProviders, err := createRequiredProviders(false, args.ExtraBlocksRootTerraform)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate required providers")
	}

	gcpProvider, err := createGcpProvider(args.ExtraProviderArguments,
		args.ServiceAccountCredentials, args.GcpProjectId, args.Regions, "", args.ProviderDefaultLabels)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate gcp provider")
	}

	laceworkProvider, err := createLaceworkProvider(args.LaceworkProfile)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework provider")
	}

	agentlessModule, err := createAgentless(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate gcp agentless module")
	}

	configurationModule, err := createConfiguration(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate gcp configuration module")
	}

	auditLogModule, err := createAuditLog(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate gcp audit log module")
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
			gcpProvider,
			laceworkProvider,
			agentlessModule,
			configurationModule,
			auditLogModule,
			outputBlocks,
			args.ExtraBlocks,
		),
	)
	return hclBlocks, nil
}

func createRequiredProviders(useExistingRequiredProviders bool,
	extraBlocks []*hclwrite.Block) (*hclwrite.Block, error) {
	if useExistingRequiredProviders {
		return nil, nil
	}
	return lwgenerate.CreateRequiredProvidersWithCustomBlocks(
		extraBlocks,
		lwgenerate.NewRequiredProvider(
			"lacework",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.LaceworkProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(lwgenerate.LaceworkProviderVersion),
		),
	)
}

func createLaceworkProvider(laceworkProfile string) (*hclwrite.Block, error) {
	if laceworkProfile != "" {
		return lwgenerate.NewProvider(
			"lacework",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{"profile": laceworkProfile}),
		).ToBlock()
	}
	return nil, nil
}

func createGcpProvider(
	extraProviderArguments map[string]interface{},
	serviceAccountCredentials string,
	projectId string,
	regionsArg []string,
	alias string,
	providerDefaultLabels map[string]interface{},
) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}

	regions := append([]string{}, regionsArg...)
	if len(regions) == 0 {
		regions = append(regions, "")
	}

	for _, region := range regions {
		attrs := map[string]interface{}{}

		// set custom args before the required ones below to ensure expected behavior (i.e., no overrides)
		for k, v := range extraProviderArguments {
			attrs[k] = v
		}
		if serviceAccountCredentials != "" {
			attrs["credentials"] = serviceAccountCredentials
		}

		if projectId != "" {
			attrs["project"] = projectId
		}

		if alias != "" {
			attrs["alias"] = alias
		}

		if region != "" {
			attrs["alias"] = region
			attrs["region"] = region
		}

		modifiers := []lwgenerate.HclProviderModifier{
			lwgenerate.HclProviderWithAttributes(attrs),
		}

		if len(providerDefaultLabels) != 0 {
			defaultLabelsBlock, err := lwgenerate.HclCreateGenericBlock(
				"default_labels",
				nil,
				providerDefaultLabels,
			)
			if err != nil {
				return nil, err
			}
			modifiers = append(modifiers, lwgenerate.HclProviderWithGenericBlocks(defaultLabelsBlock))
		}

		provider, err := lwgenerate.NewProvider(
			"google",
			modifiers...).ToBlock()
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, provider)
	}

	return blocks, nil
}

func createAgentless(args *GenerateGcpTfConfigurationArgs) ([]*hclwrite.Block, error) {
	if !args.Agentless {
		return nil, nil
	}

	blocks := []*hclwrite.Block{}

	for i, region := range args.Regions {
		moduleName := "lacework_gcp_agentless_scanning_global"
		moduleDetails := []lwgenerate.HclModuleModifier{
			lwgenerate.HclModuleWithVersion(lwgenerate.GcpAgentlessVersion),
		}

		attributes := map[string]interface{}{"regional": true}
		if i == 0 {
			attributes["global"] = true
			if len(args.ProjectFilterList) > 0 {
				attributes["project_filter_list"] = args.ProjectFilterList
			}
			if args.OrganizationIntegration {
				attributes["integration_type"] = "ORGANIZATION"
				attributes["organization_id"] = args.GcpOrganizationId
			}
		}
		if i > 0 {
			moduleName = "lacework_gcp_agentless_scanning_region_" + region
			attributes["global_module_reference"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "lacework_gcp_agentless_scanning_global"},
			)
		}

		moduleDetails = append(
			moduleDetails,
			lwgenerate.HclModuleWithProviderDetails(
				map[string]string{"google": fmt.Sprintf("google.%s", region)},
			),
		)

		moduleDetails = append(
			moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		module, err := lwgenerate.NewModule(
			moduleName,
			lwgenerate.GcpAgentlessSource,
			moduleDetails...,
		).ToBlock()
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, module)
	}

	return blocks, nil
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
				set := unique.StringSlice(args.FoldersToInclude)
				sort.Strings(set)
				attributes["folders_to_include"] = set
			}

			if len(args.FoldersToExclude) > 0 {
				set := unique.StringSlice(args.FoldersToExclude)
				sort.Strings(set)
				attributes["folders_to_exclude"] = set

				// Default true in gcp-audit-log TF module
				if !args.IncludeRootProjects {
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

		if len(args.Projects) > 0 {
			value := make(map[string]string, len(args.Projects))
			for _, p := range args.Projects {
				value[p] = p
			}
			moduleDetails = append(moduleDetails, lwgenerate.HclModuleWithForEach("project_id", value))
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		// Regions is required when Agentless integration is enabled
		// Use the first region name as the alias for Google provider if multiple regions are provided
		if args.Agentless && len(args.Regions) > 0 {
			moduleDetails = append(
				moduleDetails,
				lwgenerate.HclModuleWithProviderDetails(
					map[string]string{"google": fmt.Sprintf("google.%s", args.Regions[0])},
				),
			)
		}

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

			// Default true in gcp-audit-log TF module
			if !args.EnableUBLA {
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

		configurationModuleName := "gcp_project_level_config"
		if args.OrganizationIntegration {
			// if organization integration is true, override configModuleName to use the organization level module
			configurationModuleName = "gcp_organization_level_config"
			auditLogModuleName = "gcp_organization_level_audit_log"
			// Determine if this is the a pub-sub audit log
			if args.UsePubSubAudit {
				attributes["integration_type"] = "ORGANIZATION"
			} else {
				attributes["org_integration"] = args.OrganizationIntegration
			}
			attributes["organization_id"] = args.GcpOrganizationId

			if len(args.FoldersToInclude) > 0 {
				set := unique.StringSlice(args.FoldersToInclude)
				sort.Strings(set)
				attributes["folders_to_include"] = set
			}

			if len(args.FoldersToExclude) > 0 {
				set := unique.StringSlice(args.FoldersToExclude)
				sort.Strings(set)
				attributes["folders_to_exclude"] = set

				// Default true in gcp-audit-log TF module
				if !args.IncludeRootProjects {
					attributes["include_root_projects"] = args.IncludeRootProjects
				}
			}
		}

		if args.ExistingServiceAccount == nil && args.Configuration {
			attributes["use_existing_service_account"] = true

			cfgModuleName := configurationModuleName

			if len(args.Projects) > 0 {
				cfgModuleName = fmt.Sprintf("%s[each.key]", cfgModuleName)
			}

			attributes["service_account_name"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", cfgModuleName, "service_account_name"},
			)
			attributes["service_account_private_key"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", cfgModuleName, "service_account_private_key"},
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
		if !args.GoogleWorkspaceFilter {
			attributes["google_workspace_filter"] = args.GoogleWorkspaceFilter
		}

		// Default true in gcp-audit-log TF module
		if !args.K8sFilter {
			attributes["k8s_filter"] = args.K8sFilter
		}

		if args.Prefix != "" {
			attributes["prefix"] = args.Prefix
		}

		if args.WaitTime != "" {
			attributes["wait_time"] = args.WaitTime
		}

		if len(args.Projects) > 0 {
			value := make(map[string]string)
			for _, p := range args.Projects {
				value[p] = p
			}
			moduleDetails = append(moduleDetails, lwgenerate.HclModuleWithForEach("project_id", value))
		}

		moduleDetails = append(moduleDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		// Regions is required when Agentless integration is enabled
		// Use the first region name as the alias for Google provider if multiple regions are provided
		if args.Agentless && len(args.Regions) > 0 {
			moduleDetails = append(
				moduleDetails,
				lwgenerate.HclModuleWithProviderDetails(
					map[string]string{"google": fmt.Sprintf("google.%s", args.Regions[0])},
				),
			)
		}

		return lwgenerate.NewModule(
			auditLogModuleName,
			getAuditLogModule(args.UsePubSubAudit),
			append(moduleDetails, lwgenerate.HclModuleWithVersion(getAuditLogVersion(args.UsePubSubAudit)))...,
		).ToBlock()
	}

	return nil, nil
}

func getAuditLogModule(isPubSub bool) string {
	if isPubSub {
		return lwgenerate.GcpPubSubAuditLog
	}
	return lwgenerate.GcpAuditLogSource
}

func getAuditLogVersion(isPubSub bool) string {
	if isPubSub {
		return lwgenerate.GcpPubSubAuditLogVersion
	}
	return lwgenerate.GcpAuditLogVersion
}
