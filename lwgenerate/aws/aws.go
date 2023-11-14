// A package that generates Lacework deployment code for Amazon Web Services.
package aws

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/lwgenerate"
)

type ExistingIamRoleDetails struct {
	// Existing IAM Role ARN
	Arn string

	// Existing IAM Role Name
	Name string

	// Existing IAM Role External Id
	ExternalId string
}

func (e *ExistingIamRoleDetails) IsEmpty() bool {
	if e == nil {
		return true
	}
	return e.Arn == "" && e.Name == "" && e.ExternalId == ""
}

func (e *ExistingIamRoleDetails) IsPartial() bool {
	// If nil, return false
	if e == nil {
		return false
	}

	// If all values are empty, return false
	if e.Arn == "" && e.Name == "" && e.ExternalId == "" {
		return false
	}

	// If all values are populated, return false
	if e.Arn != "" && e.Name != "" && e.ExternalId != "" {
		return false
	}

	return true
}

// NewExistingIamRoleDetails Create new existing IAM role details
func NewExistingIamRoleDetails(name string, arn string, externalId string) *ExistingIamRoleDetails {
	return &ExistingIamRoleDetails{
		Arn:        arn,
		Name:       name,
		ExternalId: externalId,
	}
}

type OrgAccountMapping struct {
	DefaultLaceworkAccount string          `json:"default_lacework_account"`
	Mapping                []OrgAccountMap `json:"mapping"`
}

func (orgMap *OrgAccountMapping) IsEmpty() bool {
	return len(orgMap.Mapping) == 0 && orgMap.DefaultLaceworkAccount == ""
}

func (orgMap *OrgAccountMapping) ToMap() (map[string]any, error) {
	var mappings []map[string]any
	mappingsJsonString, err := json.Marshal(orgMap.Mapping)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(mappingsJsonString, &mappings)
	if err != nil {
		return nil, err
	}

	orgMap.Mapping = []OrgAccountMap{}

	result := make(map[string]any)
	jsonString, err := json.Marshal(orgMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		return nil, err
	}

	result["mapping"] = mappings
	return result, nil
}

type OrgAccountMap struct {
	LaceworkAccount string   `json:"lacework_account"`
	AwsAccounts     []string `json:"aws_accounts"`
}

type AwsSubAccount struct {
	// The name of the AwsProfile to use (in AWS configuration)
	AwsProfile string

	// The AwsRegion this profile should use if any resources are created
	AwsRegion string

	// The Alias of the provider block
	Alias string
}

// Create a new AWS sub account
//
// A subaccount consists of the profile name (which needs to match the executing machines aws configuration) and a
// region for any new resources to be created in
func NewAwsSubAccount(profile string, region string, alias ...string) AwsSubAccount {
	subaccount := AwsSubAccount{AwsProfile: profile, AwsRegion: region}
	if len(alias) > 0 {
		subaccount.Alias = alias[0]
	}
	return subaccount
}

type GenerateAwsTfConfigurationArgs struct {
	// Should we enable AWS organization integration?
	AwsOrganization bool

	// Should we configure Agentless integration in LW?
	Agentless bool

	// Agentless management AWS account ID
	AgentlessManagementAccountID string

	// Agentless monitored AWS account IDs, OUs, or the organization root.
	AgentlessMonitoredAccountIDs []string

	// Agentless monitored AWS accounts
	AgentlessMonitoredAccounts []AwsSubAccount

	// Agentless scanning AWS accounts
	AgentlessScanningAccounts []AwsSubAccount

	// Should we configure Cloudtrail integration in LW?
	Cloudtrail bool

	// Optional name for CloudTrail
	CloudtrailName string

	// Should we configure AWS organization mappings?
	AwsOrganizationMappings bool

	// Cloudtrail organization account mappings
	OrgAccountMappings OrgAccountMapping

	// OrgAccountMapping json used for flag input
	OrgAccountMappingsJson string

	// Use exisiting CloudTrail S3 bucket
	CloudtrailUseExistingS3 bool

	// Use exisiting CloudTrail SNS topic
	CloudtrailUseExistingSNSTopic bool

	// Should we configure CSPM integration in LW?
	Config bool

	// Optional name for config
	ConfigName string

	// Config addtional AWS accounts
	ConfigAdditionalAccounts []AwsSubAccount

	// Config Lacework account
	ConfigOrgLWAccount string

	// Config Lacework sub-account
	ConfigOrgLWSubaccount string

	// Config Lacework access key ID
	ConfigOrgLWAccessKeyId string

	// Config Lacework secret key
	ConfigOrgLWSecretKey string

	// Config organization ID
	ConfigOrgId string

	// Config organization unit
	ConfigOrgUnit string

	// Config resource prefix
	ConfigOrgResourcePrefix string

	// Supply an AWS region for where to find the cloudtrail resources
	// TODO @ipcrm future: support split regions for resources (s3 one place, sns another, etc)
	AwsRegion string

	// Supply an AWS Profile name for the main account, only asked if configuring multiple
	AwsProfile string

	// Supply an AWS Assume Role for the main account
	AwsAssumeRole string

	// Existing S3 Bucket ARN (Required when using existing cloudtrail)
	ExistingCloudtrailBucketArn string

	// Optionally supply existing IAM role details
	ExistingIamRole *ExistingIamRoleDetails

	// Existing SNS Topic
	ExistingSnsTopicArn string

	// Consolidated Trail
	ConsolidatedCloudtrail bool

	// Should we force destroy the bucket if it has stuff in it? (only relevant on new Cloudtrail creation)
	// DEPRECATED
	ForceDestroyS3Bucket bool

	// Enable encryption of bucket if it is created
	BucketEncryptionEnabled bool

	// Indicates that the Bucket Encryption flag has been actively set
	// this is needed to show this it was set actively to false, rather
	// than default value for bool
	BucketEncryptionEnabledSet bool

	// Optional name of bucket if creating a new one
	BucketName string

	// Arn of the KMS encryption key for S3, required when bucket encryption in enabled
	BucketSseKeyArn string

	// Enable S3 bucket notification
	S3BucketNotification bool

	// SNS Topic name if creating one and not using an existing one
	SnsTopicName string

	// Enable encryption of SNS if it is created
	SnsTopicEncryptionEnabled bool

	// Indicates that the SNS Encryption flag has been actively set
	// this is needed to show this it was set actively to false, rather
	// than default value for bool
	SnsEncryptionEnabledSet bool

	// Arn of the KMS encryption key for SNS, required when SNS encryption in enabled
	SnsTopicEncryptionKeyArn string

	// SSQ Queue name if creating one and not using an existing one
	SqsQueueName string

	// Enable encryption of SQS if it is created
	SqsEncryptionEnabled bool

	// Indicates that the SQS Encryption flag has been actively set
	// this is needed to show this it was set actively to false, rather
	// than default value for bool
	SqsEncryptionEnabledSet bool

	// Arn of the KMS encryption key for SQS, required when SQS encryption in enabled
	SqsEncryptionKeyArn string

	// For AWS Subaccounts in consolidated CT setups
	// TODO @ipcrm future: what about many individual ct/config integrations together?
	SubAccounts []AwsSubAccount

	// Lacework Profile to use
	LaceworkProfile string

	// The Lacework AWS Root Account ID
	LaceworkAccountID string

	// Lacework Organization
	LaceworkOrganizationLevel bool
}

func (args *GenerateAwsTfConfigurationArgs) IsEmpty() bool {
	if args.AwsProfile == "" && args.AwsRegion == "" && !args.Agentless && !args.Config && !args.Cloudtrail {
		return true
	}
	return false
}

// Ensure all combinations of inputs our valid for supported spec
func (args *GenerateAwsTfConfigurationArgs) Validate() error {
	if !args.Agentless && !args.Cloudtrail && !args.Config {
		return errors.New("Agentless, CloudTrail or Config integration must be enabled")
	}

	if args.AwsRegion == "" {
		return errors.New("Main AWS account region must be set")
	}

	if args.ExistingIamRole.IsPartial() {
		return errors.New("when using an existing IAM role, existing role ARN, name, and external ID all must be set")
	}

	if args.AwsOrganization {
		if args.Agentless {
			if args.AgentlessManagementAccountID == "" {
				return errors.New("must specify a management account ID for Agentless organization integration")
			}
			if len(args.AgentlessMonitoredAccountIDs) == 0 {
				return errors.New("must specify monitored account ID list for Agentless organization integration")
			}
			if len(args.AgentlessMonitoredAccounts) == 0 {
				return errors.New("must specify monitored accounts for Agentless organization integration")
			}
			if len(args.AgentlessScanningAccounts) == 0 {
				return errors.New("must specify scanning accounts for Agentless organization integration")
			}
		}
	}

	return nil
}

type AwsTerraformModifier func(c *GenerateAwsTfConfigurationArgs)

type AwsGenerateCommandExtraState struct {
	CloudtrailAdvanced         bool
	Output                     string
	AwsSubAccounts             []string
	AgentlessMonitoredAccounts []string
	AgentlessScanningAccounts  []string
	TerraformApply             bool
}

func (a *AwsGenerateCommandExtraState) IsEmpty() bool {
	return a.Output == "" && len(a.AwsSubAccounts) == 0 && !a.TerraformApply
}

// NewTerraform returns an instance of the GenerateAwsTfConfigurationArgs struct with the provided region and enabled
// settings (config/cloudtrail).
//
// Note: Additional configuration details may be set using modifiers of the AwsTerraformModifier type
//
// Basic usage: Initialize a new AwsTerraformModifier struct, with a non-default AWS profile set. Then use generate to
//
//	           create a string output of the required HCL.
//
//	hcl, err := aws.NewTerraform("us-east-1", true, true,
//	  aws.WithAwsProfile("mycorp-profile")).Generate()
func NewTerraform(
	enableAwsOrganization bool,
	enableAgentless bool,
	enableConfig bool,
	enableCloudtrail bool,
	mods ...AwsTerraformModifier,
) *GenerateAwsTfConfigurationArgs {
	config := &GenerateAwsTfConfigurationArgs{
		AwsOrganization: enableAwsOrganization,
		Agentless:       enableAgentless,
		Cloudtrail:      enableCloudtrail,
		Config:          enableConfig,
	}
	for _, m := range mods {
		m(config)
	}
	return config
}

// WithAwsProfile Set the AWS Profile to utilize for the main AWS provider
func WithAwsProfile(name string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.AwsProfile = name
	}
}

// WithAwsRegion Set the AWS region to utilize for the main AWS provider
func WithAwsRegion(region string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.AwsRegion = region
	}
}

// WithAwsAssumeRole Set the AWS Assume Role to utilize for the main AWS provider
func WithAwsAssumeRole(assumeRole string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.AwsAssumeRole = assumeRole
	}
}

// WithLaceworkProfile Set the Lacework Profile to utilize when integrating
func WithLaceworkProfile(name string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.LaceworkProfile = name
	}
}

// WithLaceworkAccountID Set the Lacework AWS root account ID to use
func WithLaceworkAccountID(accountID string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.LaceworkAccountID = accountID
	}
}

// WithAgentlessManagementAccountID Set Agentless management account ID
func WithAgentlessManagementAccountID(accountID string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.AgentlessManagementAccountID = accountID
	}
}

// WithAgentlessMonitoredAccountIDs Set Agentless monitored account IDs
func WithAgentlessMonitoredAccountIDs(accountIDs []string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.AgentlessMonitoredAccountIDs = accountIDs
	}
}

// WithAgentlessMonitoredAccounts Set Agentless monitored accounts
func WithAgentlessMonitoredAccounts(accounts ...AwsSubAccount) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.AgentlessMonitoredAccounts = accounts
	}
}

// WithAgentlessScanningAccounts Set Agentless scanning accounts
func WithAgentlessScanningAccounts(accounts ...AwsSubAccount) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.AgentlessScanningAccounts = accounts
	}
}

// WithConfigAdditionalAccounts Set Config additional accounts
func WithConfigAdditionalAccounts(accounts ...AwsSubAccount) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigAdditionalAccounts = accounts
	}
}

// WithConfigOrgLWAccount Set Config org LW account
func WithConfigOrgLWAccount(account string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigOrgLWAccount = account
	}
}

// WithConfigOrgLWSubaccount Set Config org LW sub-account
func WithConfigOrgLWSubaccount(subaccount string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigOrgLWSubaccount = subaccount
	}
}

// WithConfigOrgLWAccessKeyId Set Config org LW access key ID
func WithConfigOrgLWAccessKeyId(accessKeyId string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigOrgLWAccessKeyId = accessKeyId
	}
}

// WithConfigOrgLWSecretKey Set Config org LW secret key
func WithConfigOrgLWSecretKey(secretKey string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigOrgLWSecretKey = secretKey
	}
}

// WithConfigOrgId Set Config org ID
func WithConfigOrgId(orgId string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigOrgId = orgId
	}
}

// WithConfigOrgUnit Set Config org unit
func WithConfigOrgUnit(orgUnit string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigOrgUnit = orgUnit
	}
}

// WithConfigOrgResourcePrefix Set Config org resource prefix
func WithConfigOrgResourcePrefix(resourcePrefix string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigOrgResourcePrefix = resourcePrefix
	}
}

// WithCloudtrailUseExistingS3 Use the existing Cloudtrail S3 bucket
func WithCloudtrailUseExistingS3(useExistingS3 bool) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.CloudtrailUseExistingS3 = useExistingS3
	}
}

// WithCloudtrailUseExistingSNSTopic Use the existing Cloudtrail SNS topic
func WithCloudtrailUseExistingSNSTopic(useExistingSNSTopic bool) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.CloudtrailUseExistingSNSTopic = useExistingSNSTopic
	}
}

// WithExistingCloudtrailBucketArn Set the bucket ARN of an existing Cloudtrail setup
func WithExistingCloudtrailBucketArn(arn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ExistingCloudtrailBucketArn = arn
	}
}

// WithExistingSnsTopicArn Set the SNS Topic ARN of an existing Cloudtrail setup
func WithExistingSnsTopicArn(arn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ExistingSnsTopicArn = arn
	}
}

// WithConsolidatedCloudtrail Enable Consolidated Cloudtrail use
func WithConsolidatedCloudtrail(consolidatedCloudtrail bool) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConsolidatedCloudtrail = consolidatedCloudtrail
	}
}

// WithExistingIamRole Set an existing IAM role configuration to use with the created Terraform code
func WithExistingIamRole(iamDetails *ExistingIamRoleDetails) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ExistingIamRole = iamDetails
	}
}

// WithSubaccounts Supply additional AWS Profiles to integrate
func WithSubaccounts(subaccounts ...AwsSubAccount) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.SubAccounts = subaccounts
	}
}

// WithCloudtrailName add optional name for CloudTrail integration
func WithCloudtrailName(cloudtrailName string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.CloudtrailName = cloudtrailName
	}
}

// WithOrgAccountMappings add optional name for Organization account mappings
// Sets lacework org level to true
func WithOrgAccountMappings(mapping OrgAccountMapping) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.OrgAccountMappings = mapping
		if !mapping.IsEmpty() {
			c.LaceworkOrganizationLevel = true
		}
	}
}

// WithBucketName add bucket name for CloudTrail integration
func WithBucketName(bucketName string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.BucketName = bucketName
	}
}

// WithBucketEncryptionEnabled Enable encryption on a newly created bucket
func WithBucketEncryptionEnabled(enableBucketEncryption bool) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.BucketEncryptionEnabled = enableBucketEncryption
		c.BucketEncryptionEnabledSet = true
	}
}

// WithBucketSSEKeyArn Set existing KMS encryption key arn for bucket
func WithBucketSSEKeyArn(bucketSseKeyArn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.BucketSseKeyArn = bucketSseKeyArn
	}
}

// WithSnsTopicName Set SNS Topic Name if creating new one
func WithSnsTopicName(snsTopicName string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.SnsTopicName = snsTopicName
	}
}

// WithSnsTopicEncryptionEnabled Enable encryption on SNS Topic when created
func WithSnsTopicEncryptionEnabled(snsTopicEncryptionEnabled bool) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.SnsTopicEncryptionEnabled = snsTopicEncryptionEnabled
		c.SnsEncryptionEnabledSet = true
	}
}

// WithSnsTopicEncryptionKeyArn Set existing KMS encryption key arn for SNS topic
func WithSnsTopicEncryptionKeyArn(snsTopicEncryptionKeyArn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.SnsTopicEncryptionKeyArn = snsTopicEncryptionKeyArn
	}
}

// WithSqsQueueName Set SQS Queue Name if creating new one
func WithSqsQueueName(sqsQueueName string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.SqsQueueName = sqsQueueName
	}
}

// WithSqsEncryptionEnabled Enable encryption on SQS queue when created
func WithSqsEncryptionEnabled(sqsEncryptionEnabled bool) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.SqsEncryptionEnabled = sqsEncryptionEnabled
		c.SqsEncryptionEnabledSet = true
	}
}

// WithSqsEncryptionKeyArn Set existing KMS encryption key arn for SQS queue
func WithSqsEncryptionKeyArn(ssqEncryptionKeyArn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.SqsEncryptionKeyArn = ssqEncryptionKeyArn
	}
}

func WithS3BucketNotification(s3BucketNotifiaction bool) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.S3BucketNotification = s3BucketNotifiaction
	}
}

// Generate new Terraform code based on the supplied args.
func (args *GenerateAwsTfConfigurationArgs) Generate() (string, error) {
	// Validate inputs
	if err := args.Validate(); err != nil {
		return "", errors.Wrap(err, "invalid inputs")
	}

	// Create blocks
	requiredProviders, err := createRequiredProviders()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate required providers")
	}

	awsProvider, err := createAwsProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate aws provider")
	}

	laceworkProvider, err := createLaceworkProvider(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate lacework provider")
	}

	configModule, err := createConfig(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate aws config module")
	}

	cloudTrailModule, err := createCloudtrail(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate aws cloudtrail module")
	}

	agentlessModule, err := createAgentless(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate aws agentless global module")
	}

	// Render
	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			awsProvider,
			laceworkProvider,
			configModule,
			cloudTrailModule,
			agentlessModule),
	)
	return hclBlocks, nil
}

func createRequiredProviders() (*hclwrite.Block, error) {
	return lwgenerate.CreateRequiredProviders(
		lwgenerate.NewRequiredProvider("lacework",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.LaceworkProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(lwgenerate.LaceworkProviderVersion)))
}

func createAwsProvider(args *GenerateAwsTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}

	attributes := map[string]interface{}{
		"region": args.AwsRegion,
	}

	if args.AwsProfile != "" {
		attributes["alias"] = "main"
		attributes["profile"] = args.AwsProfile
	}

	modifiers := []lwgenerate.HclProviderModifier{
		lwgenerate.HclProviderWithAttributes(attributes),
	}

	if args.AwsAssumeRole != "" {
		assumeRoleBlock, err := lwgenerate.HclCreateGenericBlock(
			"assume_role",
			nil,
			map[string]interface{}{"role_arn": args.AwsAssumeRole},
		)
		if err != nil {
			return nil, err
		}
		modifiers = append(modifiers, lwgenerate.HclProviderWithGenericBlocks(assumeRoleBlock))
	}

	provider, err := lwgenerate.NewProvider("aws", modifiers...).ToBlock()
	if err != nil {
		return nil, err
	}

	blocks = append(blocks, provider)

	accounts := []AwsSubAccount{}
	accounts = append(accounts, args.AgentlessMonitoredAccounts...)
	accounts = append(accounts, args.AgentlessScanningAccounts...)
	accounts = append(accounts, args.ConfigAdditionalAccounts...)
	seenAccounts := []string{}

	for _, account := range accounts {
		alias := fmt.Sprintf("%s-%s", account.AwsProfile, account.AwsRegion)
		if account.Alias != "" {
			alias = account.Alias
		}
		// Skip duplicate account
		if slices.Contains(seenAccounts, alias) {
			continue
		}
		seenAccounts = append(seenAccounts, alias)
		providerBlock, err := lwgenerate.NewProvider(
			"aws",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{
				"alias":   alias,
				"profile": account.AwsProfile,
				"region":  account.AwsRegion,
			}),
		).ToBlock()
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, providerBlock)
	}

	return blocks, nil
}

func createLaceworkProvider(args *GenerateAwsTfConfigurationArgs) (*hclwrite.Block, error) {
	lwProviderAttributes := map[string]any{}

	if args.LaceworkProfile != "" {
		lwProviderAttributes["profile"] = args.LaceworkProfile
	}

	if args.LaceworkOrganizationLevel {
		lwProviderAttributes["organization"] = true
	}

	if len(lwProviderAttributes) > 0 {
		return lwgenerate.NewProvider(
			"lacework", lwgenerate.HclProviderWithAttributes(lwProviderAttributes),
		).ToBlock()
	}
	return nil, nil
}

func createConfig(args *GenerateAwsTfConfigurationArgs) ([]*hclwrite.Block, error) {
	if !args.Config {
		return nil, nil
	}

	blocks := []*hclwrite.Block{}

	if args.AwsOrganization {
		block, err := lwgenerate.NewModule(
			"aws_config",
			lwgenerate.AwsConfigOrgSource,
			lwgenerate.HclModuleWithVersion(lwgenerate.AwsConfigOrgVersion),
			lwgenerate.HclModuleWithProviderDetails(map[string]string{"aws": "aws.main"}),
			lwgenerate.HclModuleWithAttributes(
				map[string]interface{}{
					"lacework_account":       args.ConfigOrgLWAccount,
					"lacework_sub_account":   args.ConfigOrgLWSubaccount,
					"lacework_access_key_id": args.ConfigOrgLWAccessKeyId,
					"lacework_secret_key":    args.ConfigOrgLWSecretKey,
					"organization_id":        args.ConfigOrgId,
					"organization_unit":      args.ConfigOrgUnit,
					"resource_prefix":        args.ConfigOrgResourcePrefix,
				},
			),
		).ToBlock()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
		return blocks, nil
	}

	moduleDetails := []lwgenerate.HclModuleModifier{
		lwgenerate.HclModuleWithVersion(lwgenerate.AwsConfigVersion),
		lwgenerate.HclModuleWithProviderDetails(map[string]string{"aws": "aws.main"}),
	}
	if args.LaceworkAccountID != "" {
		moduleDetails = append(
			moduleDetails,
			lwgenerate.HclModuleWithAttributes(
				map[string]interface{}{"lacework_aws_account_id": args.LaceworkAccountID},
			),
		)
	}

	moduleBlock, err := lwgenerate.NewModule(
		"aws_config",
		lwgenerate.AwsConfigSource,
		moduleDetails...,
	).ToBlock()
	if err != nil {
		return nil, err
	}
	blocks = append(blocks, moduleBlock)

	for _, account := range args.ConfigAdditionalAccounts {
		configModule, err := lwgenerate.NewModule(
			fmt.Sprintf("aws_config_%s", account.Alias),
			lwgenerate.AwsConfigSource,
			lwgenerate.HclModuleWithVersion(lwgenerate.AwsConfigVersion),
			lwgenerate.HclModuleWithProviderDetails(map[string]string{
				"aws": fmt.Sprintf("aws.%s", account.Alias),
			}),
		).ToBlock()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, configModule)
	}

	return blocks, nil
}

func createCloudtrail(args *GenerateAwsTfConfigurationArgs) (*hclwrite.Block, error) {
	if !args.Cloudtrail {
		return nil, nil
	}

	attributes := map[string]interface{}{}
	modDetails := []lwgenerate.HclModuleModifier{lwgenerate.HclModuleWithVersion(lwgenerate.AwsCloudTrailVersion)}

	if args.LaceworkAccountID != "" {
		attributes["lacework_aws_account_id"] = args.LaceworkAccountID
	}
	if args.ConsolidatedCloudtrail {
		attributes["consolidated_trail"] = true
	}
	// S3 Bucket attributes
	if args.CloudtrailUseExistingS3 {
		attributes["use_existing_cloudtrail"] = true
		if args.CloudtrailName != "" {
			attributes["cloudtrail_name"] = args.CloudtrailName
		}
		if args.ExistingCloudtrailBucketArn != "" {
			attributes["bucket_arn"] = args.ExistingCloudtrailBucketArn
		}
	} else {
		if args.BucketName != "" {
			attributes["bucket_name"] = args.BucketName
		}
		if args.BucketEncryptionEnabledSet {
			if args.BucketEncryptionEnabled {
				if args.BucketSseKeyArn != "" {
					attributes["bucket_sse_key_arn"] = args.BucketSseKeyArn
				}
			} else {
				attributes["bucket_encryption_enabled"] = false
			}
		}
	}
	if args.S3BucketNotification {
		attributes["use_s3_bucket_notification"] = true
	}

	// Aws Organization CloudTrail
	if args.AwsOrganization {
		attributes["is_organization_trail"] = true

		if !args.OrgAccountMappings.IsEmpty() {
			orgAccountMappings, err := args.OrgAccountMappings.ToMap()
			if err != nil {
				return nil, errors.Wrap(err, "unable to parse 'org_account_mappings'")
			}
			attributes["org_account_mappings"] = []map[string]any{orgAccountMappings}
		}
	}

	// SNS Attributes
	if args.CloudtrailUseExistingSNSTopic {
		attributes["use_existing_sns_topic"] = true
		if args.ExistingSnsTopicArn != "" {
			attributes["sns_topic_arn"] = args.ExistingSnsTopicArn
		}
	} else {
		if args.SnsTopicName != "" {
			attributes["sns_topic_name"] = args.SnsTopicName
		}
		if args.SnsEncryptionEnabledSet {
			if args.SnsTopicEncryptionEnabled {
				if args.SnsTopicEncryptionKeyArn != "" {
					attributes["sns_topic_encryption_key_arn"] = args.SnsTopicEncryptionKeyArn
				}
			} else {
				attributes["sns_topic_encryption_enabled "] = false
			}
		}
	}
	// SQS Attributes
	if args.SqsQueueName != "" {
		attributes["sqs_queue_name"] = args.SqsQueueName
	}
	if args.SqsEncryptionEnabledSet {
		if args.SqsEncryptionEnabled {
			if args.SqsEncryptionKeyArn != "" {
				attributes["sqs_encryption_key_arn"] = args.SqsEncryptionKeyArn
			}
		} else {
			attributes["sqs_encryption_enabled "] = false
		}
	}
	if args.ExistingIamRole.IsEmpty() && args.Config {
		attributes["use_existing_iam_role"] = true
		attributes["iam_role_name"] = lwgenerate.CreateSimpleTraversal(
			[]string{"module", "aws_config", "iam_role_name"})
		attributes["iam_role_arn"] = lwgenerate.CreateSimpleTraversal(
			[]string{"module", "aws_config", "iam_role_arn"})
		attributes["iam_role_external_id"] = lwgenerate.CreateSimpleTraversal(
			[]string{"module", "aws_config", "external_id"})
	}

	if !args.ExistingIamRole.IsEmpty() {
		attributes["use_existing_iam_role"] = true
		attributes["iam_role_name"] = args.ExistingIamRole.Name
		attributes["iam_role_arn"] = args.ExistingIamRole.Arn
		attributes["iam_role_external_id"] = args.ExistingIamRole.ExternalId
	}

	modDetails = append(modDetails, lwgenerate.HclModuleWithProviderDetails(map[string]string{"aws": "aws.main"}))
	modDetails = append(modDetails,
		lwgenerate.HclModuleWithAttributes(attributes),
	)

	return lwgenerate.NewModule(
		"main_cloudtrail",
		lwgenerate.AwsCloudTrailSource,
		modDetails...,
	).ToBlock()
}

func createAgentless(args *GenerateAwsTfConfigurationArgs) ([]*hclwrite.Block, error) {
	if !args.Agentless {
		return nil, nil
	}

	blocks := []*hclwrite.Block{}

	if args.AwsOrganization {
		// Create Agenetless integration for organization

		// Add management module
		managementModule, err := lwgenerate.NewModule(
			"lacework_aws_agentless_management_scanning_role",
			lwgenerate.AwsAgentlessSource,
			lwgenerate.HclModuleWithVersion(lwgenerate.AwsAgentlessVersion),
			lwgenerate.HclModuleWithProviderDetails(map[string]string{
				"aws": "aws.main",
			}),
			lwgenerate.HclModuleWithAttributes(map[string]interface{}{
				"snapshot_role": true,
				"global_module_reference": lwgenerate.CreateSimpleTraversal(
					[]string{"module", "lacework_aws_agentless_scanning_global"},
				),
			}),
		).ToBlock()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, managementModule)

		// Add global scanning module
		monitoredAccountIDs := []string{}
		for _, accountID := range args.AgentlessMonitoredAccountIDs {
			monitoredAccountIDs = append(monitoredAccountIDs, fmt.Sprintf("\"%s\"", accountID))
		}
		globalModule, err := lwgenerate.NewModule(
			"lacework_aws_agentless_scanning_global",
			lwgenerate.AwsAgentlessSource,
			lwgenerate.HclModuleWithVersion(lwgenerate.AwsAgentlessVersion),
			lwgenerate.HclModuleWithAttributes(map[string]interface{}{
				"global":   true,
				"regional": true,
				"organization": lwgenerate.CreateMapTraversalTokens(map[string]string{
					"management_account": fmt.Sprintf("\"%s\"", args.AgentlessManagementAccountID),
					"monitored_accounts": fmt.Sprintf("[%s]", strings.Join(monitoredAccountIDs, ", ")),
				}),
			}),
			lwgenerate.HclModuleWithProviderDetails(
				map[string]string{"aws": fmt.Sprintf("aws.%s", args.AgentlessScanningAccounts[0].Alias)},
			),
		).ToBlock()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, globalModule)

		// Add regional scanning modules
		for _, account := range args.AgentlessScanningAccounts[1:] {
			regionModule, err := lwgenerate.NewModule(
				fmt.Sprintf("lacework_aws_agentless_scanning_region_%s", account.Alias),
				lwgenerate.AwsAgentlessSource,
				lwgenerate.HclModuleWithVersion(lwgenerate.AwsAgentlessVersion),
				lwgenerate.HclModuleWithProviderDetails(map[string]string{
					"aws": fmt.Sprintf("aws.%s", account.Alias),
				}),
				lwgenerate.HclModuleWithAttributes(
					map[string]interface{}{
						"regional": true,
						"global_module_reference": lwgenerate.CreateSimpleTraversal(
							[]string{"module", "lacework_aws_agentless_scanning_global"},
						),
					},
				),
			).ToBlock()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, regionModule)
		}

		// Add monitored modules
		for _, monitoredAccount := range args.AgentlessMonitoredAccounts {
			monitoredModule, err := lwgenerate.NewModule(
				fmt.Sprintf("lacework_aws_agentless_monitored_scanning_role_%s", monitoredAccount.Alias),
				lwgenerate.AwsAgentlessSource,
				lwgenerate.HclModuleWithVersion(lwgenerate.AwsAgentlessVersion),
				lwgenerate.HclModuleWithProviderDetails(map[string]string{
					"aws": fmt.Sprintf("aws.%s", monitoredAccount.Alias),
				}),
				lwgenerate.HclModuleWithAttributes(
					map[string]interface{}{
						"snapshot_role": true,
						"global_module_reference": lwgenerate.CreateSimpleTraversal(
							[]string{"module", "lacework_aws_agentless_scanning_global"},
						),
					},
				),
			).ToBlock()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, monitoredModule)
		}

		autoDeploymentBlock, err := lwgenerate.HclCreateGenericBlock(
			"auto_deployment",
			nil,
			map[string]interface{}{"enabled": true, "retain_stacks_on_account_removal": false},
		)
		if err != nil {
			return nil, err
		}
		lifecycleBlock, err := lwgenerate.HclCreateGenericBlock(
			"lifecycle",
			nil,
			map[string]interface{}{
				"ignore_changes": lwgenerate.CreateSimpleTraversal([]string{"[administration_role_arn]"}),
			},
		)
		if err != nil {
			return nil, err
		}

		stacksetResource, err := lwgenerate.NewResource(
			"aws_cloudformation_stack_set",
			"snapshot_role",
			lwgenerate.HclResourceWithAttributesAndProviderDetails(
				map[string]interface{}{
					"capabilities":     lwgenerate.CreateSimpleTraversal([]string{"[\"CAPABILITY_NAMED_IAM\"]"}),
					"description":      "Lacework AWS Agentless Workload Scanning Organization Roles",
					"name":             "lacework-agentless-scanning-stackset",
					"permission_model": "SERVICE_MANAGED",
					"template_url": "https://agentless-workload-scanner.s3.amazonaws.com" +
						"/cloudformation-lacework/latest/snapshot-role.json",
					"parameters": lwgenerate.CreateMapTraversalTokens(map[string]string{
						"ExternalId":         "module.lacework_aws_agentless_scanning_global.external_id",
						"ECSTaskRoleArn":     "module.lacework_aws_agentless_scanning_global.agentless_scan_ecs_task_role_arn",
						"ResourceNamePrefix": "module.lacework_aws_agentless_scanning_global.prefix",
						"ResourceNameSuffix": "module.lacework_aws_agentless_scanning_global.suffix",
					}),
				},
				[]string{"aws.main"},
			),
			lwgenerate.HclResourceWithGenericBlocks(autoDeploymentBlock, lifecycleBlock),
		).ToResourceBlock()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, stacksetResource)

		// Get OU IDs for the organizational_unit_ids attribute
		OUIDs := []string{}
		for _, accountID := range args.AgentlessMonitoredAccountIDs {
			if strings.HasPrefix(accountID, "ou-") {
				OUIDs = append(OUIDs, fmt.Sprintf("\"%s\"", accountID))
			}
		}

		deploymentTargetsBlock, err := lwgenerate.HclCreateGenericBlock(
			"deployment_targets",
			nil,
			map[string]interface{}{"organizational_unit_ids": lwgenerate.CreateSimpleTraversal(
				[]string{fmt.Sprintf("[%s]", strings.Join(OUIDs, ","))},
			)},
		)
		if err != nil {
			return nil, err
		}
		stacksetInstanceResource, err := lwgenerate.NewResource(
			"aws_cloudformation_stack_set_instance",
			"snapshot_role",
			lwgenerate.HclResourceWithAttributesAndProviderDetails(
				map[string]interface{}{
					"stack_set_name": lwgenerate.CreateSimpleTraversal(
						[]string{"aws_cloudformation_stack_set", "snapshot_role", "name"},
					),
				},
				[]string{"aws.main"},
			),
			lwgenerate.HclResourceWithGenericBlocks(deploymentTargetsBlock),
		).ToResourceBlock()

		if err != nil {
			return nil, err
		}

		blocks = append(blocks, stacksetInstanceResource)
	} else {
		// Create Agenetless integration for single account
		globalModule, err := lwgenerate.NewModule(
			"lacework_aws_agentless_scanning_global",
			lwgenerate.AwsAgentlessSource,
			lwgenerate.HclModuleWithVersion(lwgenerate.AwsAgentlessVersion),
			lwgenerate.HclModuleWithAttributes(map[string]interface{}{"global": true, "regional": true}),
			lwgenerate.HclModuleWithProviderDetails(
				map[string]string{"aws": "aws.main"},
			),
		).ToBlock()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, globalModule)

		for _, account := range args.AgentlessScanningAccounts {
			regionModule, err := lwgenerate.NewModule(
				fmt.Sprintf("lacework_aws_agentless_scanning_region_%s", account.Alias),
				lwgenerate.AwsAgentlessSource,
				lwgenerate.HclModuleWithVersion(lwgenerate.AwsAgentlessVersion),
				lwgenerate.HclModuleWithProviderDetails(map[string]string{
					"aws": fmt.Sprintf("aws.%s", account.Alias),
				}),
				lwgenerate.HclModuleWithAttributes(
					map[string]interface{}{
						"regional": true,
						"global_module_reference": lwgenerate.CreateSimpleTraversal(
							[]string{"module", "lacework_aws_agentless_scanning_global"},
						),
					},
				),
			).ToBlock()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, regionModule)
		}
	}

	return blocks, nil
}
