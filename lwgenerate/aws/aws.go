// A package that generates Lacework deployment code for Amazon Web Services.
package aws

import (
	"fmt"

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
	// Should we configure Agentless integration in LW?
	Agentless bool

	// Should we configure Cloudtrail integration in LW?
	Cloudtrail bool

	// Optional name for CloudTrail
	CloudtrailName string

	// Should we configure CSPM integration in LW?
	Config bool

	// Optional name for config
	ConfigName string

	// Supply an AWS region for where to find the cloudtrail resources
	// TODO @ipcrm future: support split regions for resources (s3 one place, sns another, etc)
	AwsRegion string

	// Supply an AWS Profile name for the main account, only asked if configuring multiple
	AwsProfile string

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

	S3BucketNotification bool
}

// Ensure all combinations of inputs our valid for supported spec
func (args *GenerateAwsTfConfigurationArgs) validate() error {
	// Validate one of config or cloudtrail was enabled; otherwise error out
	if !args.Agentless && !args.Cloudtrail && !args.Config {
		return errors.New("agentless, cloudtrail or config integration must be enabled")
	}

	// Validate that at least region was set
	if args.AwsRegion == "" {
		return errors.New("AWS region must be set")
	}

	// Validate existing role IAM values, if set
	if args.ExistingIamRole != nil {
		if args.ExistingIamRole.Arn == "" ||
			args.ExistingIamRole.Name == "" ||
			args.ExistingIamRole.ExternalId == "" {
			return errors.New("when using an existing IAM role, existing role ARN, name, and external ID all must be set")
		}
	}

	return nil
}

type AwsTerraformModifier func(c *GenerateAwsTfConfigurationArgs)

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
	region string, enableAgentless bool, enableConfig bool, enableCloudtrail bool, mods ...AwsTerraformModifier,
) *GenerateAwsTfConfigurationArgs {
	config := &GenerateAwsTfConfigurationArgs{
		AwsRegion:  region,
		Agentless:  enableAgentless,
		Cloudtrail: enableCloudtrail,
		Config:     enableConfig,
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

// ExistingCloudtrailBucketArn Set the bucket ARN of an existing Cloudtrail setup
func ExistingCloudtrailBucketArn(arn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ExistingCloudtrailBucketArn = arn
	}
}

// ExistingSnsTopicArn Set the SNS Topic ARN of an existing Cloudtrail setup
func ExistingSnsTopicArn(arn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ExistingSnsTopicArn = arn
	}
}

// UseConsolidatedCloudtrail Enable Consolidated Cloudtrail use
func UseConsolidatedCloudtrail() AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConsolidatedCloudtrail = true
	}
}

// UseExistingIamRole Set an existing IAM role configuration to use with the created Terraform code
func UseExistingIamRole(iamDetails *ExistingIamRoleDetails) AwsTerraformModifier {
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

// WithConfigName add optional name for Config integration
func WithConfigName(configName string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConfigName = configName
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
	if err := args.validate(); err != nil {
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
	if args.AwsRegion != "" || args.AwsProfile != "" || len(args.SubAccounts) > 0 {
		attrs := map[string]interface{}{}
		if args.AwsRegion != "" {
			attrs["region"] = args.AwsRegion
		}

		if args.AwsProfile != "" {
			attrs["profile"] = args.AwsProfile
		}

		if len(args.SubAccounts) > 0 {
			attrs["alias"] = "main"
		}

		provider, err := lwgenerate.NewProvider(
			"aws",
			lwgenerate.HclProviderWithAttributes(attrs),
		).ToBlock()
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, provider)
	}

	if len(args.SubAccounts) > 0 {
		for _, subaccount := range args.SubAccounts {
			attrs := map[string]interface{}{
				"alias":   subaccount.AwsProfile,
				"profile": subaccount.AwsProfile,
				"region":  subaccount.AwsRegion,
			}
			providerBlock, err := lwgenerate.NewProvider(
				"aws",
				lwgenerate.HclProviderWithAttributes(attrs),
			).ToBlock()

			if err != nil {
				return nil, err
			}

			blocks = append(blocks, providerBlock)
		}
	}

	return blocks, nil
}

func createLaceworkProvider(args *GenerateAwsTfConfigurationArgs) (*hclwrite.Block, error) {
	if args.LaceworkProfile != "" {
		return lwgenerate.NewProvider(
			"lacework",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{"profile": args.LaceworkProfile}),
		).ToBlock()
	}
	return nil, nil
}

func createConfig(args *GenerateAwsTfConfigurationArgs) ([]*hclwrite.Block, error) {
	blocks := []*hclwrite.Block{}
	if args.Config {
		// Add main account
		moduleDetails := []lwgenerate.HclModuleModifier{}
		if len(args.SubAccounts) > 0 {
			moduleDetails = append(moduleDetails,
				lwgenerate.HclModuleWithProviderDetails(map[string]string{"aws": "aws.main"}))
		}
		if args.LaceworkAccountID != "" {
			moduleDetails = append(moduleDetails, lwgenerate.HclModuleWithAttributes(
				map[string]interface{}{"lacework_aws_account_id": args.LaceworkAccountID}))
		}

		moduleBlock, err := lwgenerate.NewModule(
			"aws_config",
			lwgenerate.AwsConfigSource,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(lwgenerate.AwsConfigVersion))...,
		).ToBlock()

		if err != nil {
			return nil, err
		}
		blocks = append(blocks, moduleBlock)

		// Add sub accounts
		for _, subaccount := range args.SubAccounts {
			configModule, err := lwgenerate.NewModule(fmt.Sprintf(
				"aws_config_%s",
				subaccount.AwsProfile),
				lwgenerate.AwsConfigSource,
				lwgenerate.HclModuleWithVersion(lwgenerate.AwsConfigVersion),
				lwgenerate.HclModuleWithProviderDetails(map[string]string{
					"aws": fmt.Sprintf("aws.%s", subaccount.AwsProfile),
				})).ToBlock()

			if err != nil {
				return nil, err
			}

			blocks = append(blocks, configModule)
		}
	}

	return blocks, nil
}

func createCloudtrail(args *GenerateAwsTfConfigurationArgs) (*hclwrite.Block, error) {
	if args.Cloudtrail {
		attributes := map[string]interface{}{}
		modDetails := []lwgenerate.HclModuleModifier{lwgenerate.HclModuleWithVersion(lwgenerate.AwsCloudTrailVersion)}

		if args.LaceworkAccountID != "" {
			attributes["lacework_aws_account_id"] = args.LaceworkAccountID
		}
		if args.ConsolidatedCloudtrail {
			attributes["consolidated_trail"] = true
		}
		if args.CloudtrailName != "" {
			attributes["cloudtrail_name"] = args.CloudtrailName
		}
		// S3 Bucket attributes
		if args.ExistingCloudtrailBucketArn != "" {
			attributes["use_existing_cloudtrail"] = true
			attributes["bucket_arn"] = args.ExistingCloudtrailBucketArn
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
		// SNS Attributes
		if args.ExistingSnsTopicArn != "" {
			attributes["use_existing_sns_topic"] = true
			attributes["sns_topic_arn"] = args.ExistingSnsTopicArn
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
		if args.ExistingIamRole == nil && args.Config {
			attributes["use_existing_iam_role"] = true
			attributes["iam_role_name"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "aws_config", "iam_role_name"})
			attributes["iam_role_arn"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "aws_config", "iam_role_arn"})
			attributes["iam_role_external_id"] = lwgenerate.CreateSimpleTraversal(
				[]string{"module", "aws_config", "external_id"})
		}

		if args.ExistingIamRole != nil {
			attributes["use_existing_iam_role"] = true
			attributes["iam_role_name"] = args.ExistingIamRole.Name
			attributes["iam_role_arn"] = args.ExistingIamRole.Arn
			attributes["iam_role_external_id"] = args.ExistingIamRole.ExternalId
		}

		if args.S3BucketNotification {
			attributes["use_s3_bucket_notification"] = true
		}

		if len(args.SubAccounts) > 0 {
			modDetails = append(modDetails, lwgenerate.HclModuleWithProviderDetails(map[string]string{"aws": "aws.main"}))
		}

		modDetails = append(modDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		return lwgenerate.NewModule(
			"main_cloudtrail",
			lwgenerate.AwsCloudTrailSource,
			modDetails...,
		).ToBlock()
	}

	return nil, nil
}

func createAgentless(args *GenerateAwsTfConfigurationArgs) ([]*hclwrite.Block, error) {
	if !args.Agentless {
		return nil, nil
	}

	blocks := []*hclwrite.Block{}

	// Add global module
	globalModule, err := lwgenerate.NewModule(
		"lacework_aws_agentless_scanning_global",
		lwgenerate.AwsAgentlessSource,
		lwgenerate.HclModuleWithVersion(lwgenerate.AwsAgentlessVersion),
		lwgenerate.HclModuleWithAttributes(map[string]interface{}{"global": true, "regional": true}),
	).ToBlock()

	if err != nil {
		return nil, err
	}

	blocks = append(blocks, globalModule)

	// Add region modules
	for _, subaccount := range args.SubAccounts {
		regionModule, err := lwgenerate.NewModule(
			fmt.Sprintf("lacework_aws_agentless_scanning_region_%s", subaccount.AwsProfile),
			lwgenerate.AwsAgentlessSource,
			lwgenerate.HclModuleWithVersion(lwgenerate.AwsAgentlessVersion),
			lwgenerate.HclModuleWithProviderDetails(map[string]string{
				"aws": fmt.Sprintf("aws.%s", subaccount.AwsProfile),
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

	return blocks, nil
}
