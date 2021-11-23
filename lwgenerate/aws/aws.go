package aws

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/lwgenerate"
	"github.com/pkg/errors"
)

type ExistingIamRoleDetails struct {
	// Existing IAM Role ARN
	Arn string

	// Existing IAM Role Name
	Name string

	// Existing IAM Role External Id
	ExternalId string
}

// Create new existing IAM role details
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
}

// Create a new AWS sub account
//
// A subaccount consists of the profile name (which needs to match the executing machines aws configuration) and a
// region for any new resources to be created in
func NewAwsSubAccount(profile string, region string) AwsSubAccount {
	return AwsSubAccount{AwsProfile: profile, AwsRegion: region}
}

type GenerateAwsTfConfigurationArgs struct {
	// Should we configure Cloudtrail integration in LW?
	Cloudtrail bool

	// Should we configure CSPM integration in LW?
	Config bool

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
	ForceDestroyS3Bucket bool

	// For AWS Subaccounts in consolidated CT setups
	// TODO @ipcrm future: what about many individual ct/config integrations together?
	SubAccounts []AwsSubAccount

	// Lacework Profile to use
	LaceworkProfile string
}

// Ensure all combinations of inputs our valid for supported spec
func (args *GenerateAwsTfConfigurationArgs) validate() error {
	// Validate one of config or cloudtrail was enabled; otherwise error out
	if !args.Cloudtrail && !args.Config {
		return errors.New("cloudtrail or config integration must be enabled")
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
//              create a string output of the required HCL.
//
//   hcl, err := aws.NewTerraform("us-east-1", true, true,
//     aws.WithAwsProfile("mycorp-profile")).Generate()
//
func NewTerraform(region string, enableConfig bool, enableCloudtrail bool, mods ...AwsTerraformModifier) *GenerateAwsTfConfigurationArgs {
	config := &GenerateAwsTfConfigurationArgs{AwsRegion: region, Cloudtrail: enableCloudtrail, Config: enableConfig}
	for _, m := range mods {
		m(config)
	}
	return config
}

// Set the AWS Profile to utilize for the main AWS provider
func WithAwsProfile(name string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.AwsProfile = name
	}
}

// Set the Lacework Profile to utilize when integrating
func WithLaceworkProfile(name string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.LaceworkProfile = name
	}
}

// Set the bucket ARN of an existing Cloudtrail setup
func ExistingCloudtrailBucketArn(arn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ExistingCloudtrailBucketArn = arn
	}
}

// Set the SNS Topic ARN of an existing Cloudtrail setup
func ExistingSnsTopicArn(arn string) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ExistingSnsTopicArn = arn
	}
}

// Enable Consolidated Cloudtrail use
func UseConsolidatedCloudtrail() AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ConsolidatedCloudtrail = true
	}
}

// Set the S3 ForceDestroy parameter to true for newly created buckets
func EnableForceDestroyS3Bucket() AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ForceDestroyS3Bucket = true
	}
}

// Set an existing IAM role configuration to use with the created Terraform code
func UseExistingIamRole(iamDetails *ExistingIamRoleDetails) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.ExistingIamRole = iamDetails
	}
}

// Supply additional AWS Profiles to integrate
func WithSubaccounts(subaccounts ...AwsSubAccount) AwsTerraformModifier {
	return func(c *GenerateAwsTfConfigurationArgs) {
		c.SubAccounts = subaccounts
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

	// Render
	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			awsProvider,
			laceworkProvider,
			configModule,
			cloudTrailModule),
	)
	return hclBlocks, nil
}

func createRequiredProviders() (*hclwrite.Block, error) {
	return lwgenerate.CreateRequiredProviders(
		lwgenerate.NewRequiredProvider("lacework",
			lwgenerate.HclRequiredProviderWithSource("lacework/lacework"),
			lwgenerate.HclRequiredProviderWithVersion("~> 0.12.2")))
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

		provider, err := lwgenerate.NewProvider("aws", lwgenerate.HclProviderWithAttributes(attrs)).ToBlock()
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
			providerBlock, err := lwgenerate.NewProvider("aws", lwgenerate.HclProviderWithAttributes(attrs)).ToBlock()

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
		return lwgenerate.NewProvider("lacework",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{"profile": args.LaceworkProfile}),
		).ToBlock()
	}
	return nil, nil
}

func createConfig(args *GenerateAwsTfConfigurationArgs) ([]*hclwrite.Block, error) {
	source := "lacework/config/aws"
	version := "~> 0.1"

	blocks := []*hclwrite.Block{}
	if args.Config {
		// Add main account
		moduleDetails := []lwgenerate.HclModuleModifier{}
		if len(args.SubAccounts) > 0 {
			moduleDetails = append(moduleDetails,
				lwgenerate.HclModuleWithProviderDetails(map[string]string{"aws": "aws.main"}))
		}

		moduleBlock, err := lwgenerate.NewModule("aws_config", source,
			append(moduleDetails, lwgenerate.HclModuleWithVersion(version))...).ToBlock()

		if err != nil {
			return nil, err
		}
		blocks = append(blocks, moduleBlock)

		// Add sub accounts
		for _, subaccount := range args.SubAccounts {
			configModule, err := lwgenerate.NewModule(fmt.Sprintf("aws_config_%s", subaccount.AwsProfile),
				source,
				lwgenerate.HclModuleWithVersion(version),
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
		modDetails := []lwgenerate.HclModuleModifier{lwgenerate.HclModuleWithVersion("~> 0.1")}

		if args.ForceDestroyS3Bucket && args.ExistingCloudtrailBucketArn == "" {
			attributes["bucket_force_destroy"] = true
		}

		if args.ConsolidatedCloudtrail {
			attributes["consolidated_trail"] = true
		}

		if args.ExistingSnsTopicArn != "" {
			attributes["use_existing_sns_topic"] = true
			attributes["sns_topic_arn"] = args.ExistingSnsTopicArn
		}

		if args.ExistingIamRole == nil && args.Config {
			attributes["use_existing_iam_role"] = true
			attributes["iam_role_name"] = lwgenerate.CreateSimpleTraversal([]string{"module", "aws_config", "iam_role_name"})
			attributes["iam_role_arn"] = lwgenerate.CreateSimpleTraversal([]string{"module", "aws_config", "iam_role_arn"})
			attributes["iam_role_external_id"] = lwgenerate.CreateSimpleTraversal([]string{"module", "aws_config", "external_id"})
		}

		if args.ExistingIamRole != nil {
			attributes["use_existing_iam_role"] = true
			attributes["iam_role_name"] = args.ExistingIamRole.Name
			attributes["iam_role_arn"] = args.ExistingIamRole.Arn
			attributes["iam_role_external_id"] = args.ExistingIamRole.ExternalId
		}

		if args.ExistingCloudtrailBucketArn != "" {
			attributes["use_existing_cloudtrail"] = true
			attributes["bucket_arn"] = args.ExistingCloudtrailBucketArn
		}

		if len(args.SubAccounts) > 0 {
			modDetails = append(modDetails, lwgenerate.HclModuleWithProviderDetails(map[string]string{"aws": "aws.main"}))
		}

		modDetails = append(modDetails,
			lwgenerate.HclModuleWithAttributes(attributes),
		)

		return lwgenerate.NewModule("main_cloudtrail", "lacework/cloudtrail/aws", modDetails...).ToBlock()
	}

	return nil, nil
}
