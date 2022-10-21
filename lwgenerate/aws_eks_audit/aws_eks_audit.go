package aws_eks_audit

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/lwgenerate"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

type ExistingCrossAccountIamRoleDetails struct {
	// Existing IAM Role ARN
	Arn string

	// Existing IAM Role External ID
	ExternalId string
}

func (e *ExistingCrossAccountIamRoleDetails) IsPartial() bool {
	// If nil, return false
	if e == nil {
		return false
	}

	// If all values are empty, return false
	if e.Arn == "" && e.ExternalId == "" {
		return false
	}

	// If all values are populated, return false
	if e.Arn != "" && e.ExternalId != "" {
		return false
	}

	return true
}

// NewExistingCrossAccountIamRoleDetails Create new existing IAM role details
func NewExistingCrossAccountIamRoleDetails(arn string, externalId string) *ExistingCrossAccountIamRoleDetails {
	return &ExistingCrossAccountIamRoleDetails{
		Arn:        arn,
		ExternalId: externalId,
	}
}

type GenerateAwsEksAuditTfConfigurationArgs struct {

	// Supply an AWS Profile name
	AwsProfile string

	// Should we require MFA for object deletion?
	BucketEnableMfaDelete bool

	// Should we enable bucket encryption?
	BucketEnableEncryption bool

	// Should we force destroy the bucket if it has stuff in it?
	BucketForceDestroy bool

	// The lifetime, in days, of the bucket objects. The value must be a non-zero positive integer
	BucketLifecycleExpirationDays int

	// The encryption algorithm to use for S3 bucket server-side encryption
	BucketSseAlgorithm string

	// The ARN of the KMS encryption key to be used for S3
	// (Required when bucket_sse_algorithm is aws:kms and using an existing kms key)
	BucketSseKeyArn string

	// Should we enable bucket versioning?
	BucketVersioning bool

	// The name of the AWS EKS Audit Log integration in Lacework. Defaults to "TF AWS EKS Audit Log"
	EksAuditIntegrationName string

	// Optionally supply existing cloudwatch IAM role ARN
	ExistingCloudWatchIamRoleArn string

	// Optionally supply existing cross account IAM role details
	ExistingCrossAccountIamRole *ExistingCrossAccountIamRoleDetails

	// Should we allow the user to configure an existing Firehose IAM role?
	ExistingFirehoseIam bool

	// Optionally supply existing firehose role ARN if ExistingFirehoseIam is true
	ExistingFirehoseIamRoleArn string

	// The Cloudwatch Log Subscription Filter pattern
	FilterPattern string

	// Should encryption be disabled on the created firehose?
	// This will set kinesis_firehose_encryption_enabled which defaults to true
	FirehoseEncryptionDisabled bool

	// The ARN of an existing KMS encryption key to be used for the Kinesis Firehose
	FirehoseEncryptionKeyArn string

	// The waiting period, specified in number of days. Defaults to 30.
	KmsKeyDeletionDays int

	// Whether the KMS key is a multi-region or regional key
	KmsKeyMultiRegion bool

	// Enable KMS automatic key rotation
	KmsKeyRotation bool

	// The prefix that will be used at the beginning of every generated resource. Defaults to "lw-eks-al"
	Prefix string

	// Parsed version of RegionClusterMap
	RegionClusterMap map[string]string

	// Parsed version of RegionClusterMap
	ParsedRegionClusterMap map[string][]string

	// Should encryption be disabled for the sns topic?
	// This will set sns_topic_encryption_enabled which defaults to true
	SnsTopicEncryptionDisabled bool

	// The ARN of an existing KMS encryption key to be used for the SNS topic
	SnsTopicEncryptionKeyArn string

	// Lacework Profile to use
	LaceworkProfile string
}

// Ensure all combinations of inputs our valid for supported spec
func (args *GenerateAwsEksAuditTfConfigurationArgs) validate() error {

	// Validate that at least one region with clusters was set
	if len(args.ParsedRegionClusterMap) == 0 {
		return errors.New("At least one region with a list of cluster(s) must be set")
	}

	// Validate, at least 1 cluster must be supplied per region
	for _, clusters := range args.ParsedRegionClusterMap {
		if len(clusters) == 0 {
			return errors.New("At least one cluster must be supplied per region")
		}
	}

	// Validate existing role IAM values, if set
	if args.ExistingCrossAccountIamRole != nil {
		if args.ExistingCrossAccountIamRole.Arn == "" ||
			args.ExistingCrossAccountIamRole.ExternalId == "" {
			return errors.New("when using an existing cross account IAM role, existing role ARN and external ID all must be set")
		}
	}

	return nil
}

type AwsEksAuditTerraformModifier func(c *GenerateAwsEksAuditTfConfigurationArgs)

// NewTerraform returns an instance of the GenerateAwsEksAuditTfConfigurationArgs struct.
//
// Note: Additional configuration details may be set using modifiers of the AwsEksAuditTerraformModifier type
//
// Basic usage: Initialize a new AwsEksAuditTerraformModifier struct, with a non-default AWS profile set. Then use generate to
//
//	           create a string output of the required HCL.
//
//	hcl, err := aws.NewTerraform({"us-east-1": ["cluster1", "cluster2"], "us-east-2": ["cluster3"]}
//	  aws.WithAwsProfile("mycorp-profile")).Generate()
func NewTerraform(mods ...AwsEksAuditTerraformModifier) *GenerateAwsEksAuditTfConfigurationArgs {
	config := &GenerateAwsEksAuditTfConfigurationArgs{}
	for _, m := range mods {
		m(config)
	}
	return config
}

// WithAwsProfile Set the AWS Profile to utilize when integrating
func WithAwsProfile(name string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.AwsProfile = name
	}
}

// EnableBucketMfaDelete Set the S3 MfaDelete parameter to true for newly created buckets
func EnableBucketMfaDelete() AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.BucketEnableMfaDelete = true
	}
}

// EnableBucketEncryption Set the S3 Encryption parameter to true for newly created buckets
func EnableBucketEncryption() AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.BucketEnableEncryption = true
	}
}

// EnableBucketForceDestroy Set the S3 ForceDestroy parameter to true for newly created buckets
func EnableBucketForceDestroy() AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.BucketForceDestroy = true
	}
}

// WithBucketLifecycleExpirationDays Set the S3 Lifecycle Expiration Days parameter for newly created buckets
func WithBucketLifecycleExpirationDays(days int) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.BucketLifecycleExpirationDays = days
	}
}

// WithBucketSseAlgorithm Set the encryption algorithm to use for S3 bucket server-side encryption
func WithBucketSseAlgorithm(algorithm string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.BucketSseAlgorithm = algorithm
	}
}

// WithBucketSseKeyArn Set the ARN of the KMS encryption key to be used for S3
// (Required when bucket_sse_algorithm is aws:kms and using an existing aws_kms_key)
func WithBucketSseKeyArn(arn string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.BucketSseKeyArn = arn
	}
}

// EnableBucketVersioning Set the S3 Bucket versioning parameter to true for newly created buckets
func EnableBucketVersioning() AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.BucketVersioning = true
	}
}

// WithEksAuditIntegrationName Set the name of the EKS audit integration
func WithEksAuditIntegrationName(name string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.EksAuditIntegrationName = name
	}
}

// WithExistingCloudWatchIamRoleArn  Set an existing cloudwatch IAM role ARN
func WithExistingCloudWatchIamRoleArn(arn string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.ExistingCloudWatchIamRoleArn = arn
	}
}

// WithExistingCrossAccountIamRole Set an existing cross account IAM role configuration to use with
// the created Terraform code
func WithExistingCrossAccountIamRole(iamDetails *ExistingCrossAccountIamRoleDetails) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.ExistingCrossAccountIamRole = iamDetails
	}
}

// WithExistingFirehoseIamRoleArn  Set an existing firehose IAM role ARN
func WithExistingFirehoseIamRoleArn(arn string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.ExistingFirehoseIamRoleArn = arn
	}
}

// WithFilterPattern Set the filter pattern for the Cloudwatch subscription filter
func WithFilterPattern(pattern string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.FilterPattern = pattern
	}
}

// DisableFirehoseEncryption Set the firehose encryption parameter to true for newly created firehose
func DisableFirehoseEncryption() AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.FirehoseEncryptionDisabled = true
	}
}

// WithFirehoseEncryptionKeyArn Set the ARN of an existing KMS encryption key to be used
// with the Kinesis Firehose
func WithFirehoseEncryptionKeyArn(arn string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.FirehoseEncryptionKeyArn = arn
	}
}

// WithKmsKeyDeletionDays Set the KMS deletion waiting period, specified in number of days
func WithKmsKeyDeletionDays(days int) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.KmsKeyDeletionDays = days
	}
}

// EnableKmsKeyMultiRegion Set whether the KMS key is a multi-region or regional key
func EnableKmsKeyMultiRegion() AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.KmsKeyMultiRegion = true
	}
}

// EnableKmsKeyRotation Set KMS automatic key rotation to true
func EnableKmsKeyRotation() AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.KmsKeyRotation = true
	}
}

// WithPrefix Set the prefix that will be used at the beginning of every generated resource
func WithPrefix(prefix string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.Prefix = prefix
	}
}

// WithParsedRegionClusterMap Set the region cluster map.
// This is a list of clusters per AWS region
func WithParsedRegionClusterMap(regionClusterMap map[string][]string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.ParsedRegionClusterMap = regionClusterMap
	}
}

// DisableSnsTopicEncryption Set whether encryption should be enabled for the sns topic
func DisableSnsTopicEncryption() AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.SnsTopicEncryptionDisabled = true
	}
}

// WithSnsTopicEncryptionKeyArn Set the ARN of an existing KMS encryption key to be used
// with the SNS Topic
func WithSnsTopicEncryptionKeyArn(arn string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.SnsTopicEncryptionKeyArn = arn
	}
}

// WithLaceworkProfile Set the Lacework Profile to utilize when integrating
func WithLaceworkProfile(name string) AwsEksAuditTerraformModifier {
	return func(c *GenerateAwsEksAuditTfConfigurationArgs) {
		c.LaceworkProfile = name
	}
}

// Generate new Terraform code based on the supplied args.
func (args *GenerateAwsEksAuditTfConfigurationArgs) Generate() (string, error) {
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

	eksAuditModule, err := createEksAudit(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate aws eks audit module & resources")
	}

	// Render
	hclBlocks := lwgenerate.CreateHclStringOutput(
		lwgenerate.CombineHclBlocks(
			requiredProviders,
			awsProvider,
			laceworkProvider,
			eksAuditModule),
	)
	return hclBlocks, nil
}

func createRequiredProviders() (*hclwrite.Block, error) {
	return lwgenerate.CreateRequiredProviders(
		lwgenerate.NewRequiredProvider("lacework",
			lwgenerate.HclRequiredProviderWithSource(lwgenerate.LaceworkProviderSource),
			lwgenerate.HclRequiredProviderWithVersion(lwgenerate.LaceworkProviderVersion)))
}

func createAwsProvider(args *GenerateAwsEksAuditTfConfigurationArgs) ([]*hclwrite.Block, error) {
	var blocks []*hclwrite.Block
	if len(args.ParsedRegionClusterMap) > 1 {
		for region := range args.ParsedRegionClusterMap {
			attrs := map[string]interface{}{
				"alias":  region,
				"region": region,
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
	} else if len(args.ParsedRegionClusterMap) == 1 {
		for region := range args.ParsedRegionClusterMap {
			attrs := map[string]interface{}{
				"region": region,
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

func createLaceworkProvider(args *GenerateAwsEksAuditTfConfigurationArgs) (*hclwrite.Block, error) {
	if args.LaceworkProfile != "" {
		return lwgenerate.NewProvider(
			"lacework",
			lwgenerate.HclProviderWithAttributes(map[string]interface{}{"profile": args.LaceworkProfile}),
		).ToBlock()
	}
	return nil, nil
}

func createEksAudit(args *GenerateAwsEksAuditTfConfigurationArgs) ([]*hclwrite.Block, error) {
	var blocks []*hclwrite.Block
	moduleAttrs := map[string]interface{}{}
	resourceAttrs := map[string]interface{}{}
	moduleDetails := []lwgenerate.HclModuleModifier{lwgenerate.HclModuleWithVersion(lwgenerate.AwsEksAuditVersion)}

	if args.BucketEnableMfaDelete {
		moduleAttrs["bucket_enable_mfa_delete"] = true
	}

	if args.BucketEnableEncryption {
		moduleAttrs["bucket_encryption_enabled"] = true
	}

	if args.BucketForceDestroy {
		moduleAttrs["bucket_force_destroy"] = true
	}

	if args.BucketLifecycleExpirationDays > 0 {
		moduleAttrs["bucket_lifecycle_expiration_days"] = args.BucketLifecycleExpirationDays
	}

	if args.BucketSseAlgorithm != "" && args.BucketEnableEncryption {
		moduleAttrs["bucket_sse_algorithm"] = args.BucketSseAlgorithm
	}

	if args.BucketSseKeyArn != "" && args.BucketEnableEncryption {
		moduleAttrs["bucket_key_arn"] = args.BucketSseAlgorithm
	}

	if args.ExistingCloudWatchIamRoleArn != "" {
		moduleAttrs["use_existing_cloudwatch_iam_role"] = true
		moduleAttrs["cloudwatch_iam_role_arn"] = args.ExistingCloudWatchIamRoleArn
	}

	if args.ExistingCrossAccountIamRole != nil {
		moduleAttrs["use_existing_cross_account_iam_role"] = true
		moduleAttrs["iam_role_arn"] = args.ExistingCrossAccountIamRole.Arn
		moduleAttrs["iam_role_external_id"] = args.ExistingCrossAccountIamRole.ExternalId
	}

	if args.ExistingFirehoseIamRoleArn != "" {
		moduleAttrs["use_existing_firehose_iam_role"] = true
		moduleAttrs["firehose_iam_role_arn"] = args.ExistingFirehoseIamRoleArn
	}

	if args.FilterPattern != "" {
		moduleAttrs["filter_pattern"] = args.FilterPattern
	}

	if args.FirehoseEncryptionDisabled {
		moduleAttrs["kinesis_firehose_encryption_enabled"] = false
	}

	if args.FirehoseEncryptionKeyArn != "" && !args.FirehoseEncryptionDisabled {
		moduleAttrs["kinesis_firehose_key_arn"] = args.FirehoseEncryptionKeyArn
	}

	if args.KmsKeyDeletionDays >= 7 && args.KmsKeyDeletionDays <= 30 {
		moduleAttrs["kms_key_deletion_days"] = args.KmsKeyDeletionDays
	}

	if args.KmsKeyMultiRegion {
		moduleAttrs["kms_key_multi_region"] = true
	}

	if args.KmsKeyRotation {
		moduleAttrs["kms_key_rotation"] = true
	}

	if args.SnsTopicEncryptionDisabled {
		moduleAttrs["sns_topic_encryption_enabled"] = false
	}

	if args.SnsTopicEncryptionKeyArn != "" && !args.SnsTopicEncryptionDisabled {
		moduleAttrs["sns_topic_key_arn"] = args.SnsTopicEncryptionKeyArn
	}

	var regionList []string

	if len(args.ParsedRegionClusterMap) > 1 {
		// set no_cw_subscription_filter if we have more than 1 region in the ParsedRegionClusterMap
		moduleAttrs["no_cw_subscription_filter"] = true

		// Add aws_cloudwatch_log_subscription_filter(s) resource per region
		for region, clusters := range args.ParsedRegionClusterMap {
			regionList = append(regionList, region)

			// create hcl tokens for each cluster and create a token array to be added to our hcl
			//tuple. (we are unable to add the for loop inside the call to TokensForTuple)
			var tokens []hclwrite.Tokens
			for _, cluster := range clusters {
				tokens = append(tokens, hclwrite.TokensForValue(cty.StringVal(cluster)))
			}

			// the for_each input must be in the following format toset(["", ""])
			// In order to achieve this we can use TokensForTuple to build the tuple `[]`
			// then TokensForFunctionCall to wrap this with our call to the `toset()` function
			resourceAttrs["for_each"] = hclwrite.TokensForFunctionCall(
				"toset",
				hclwrite.TokensForTuple(tokens),
			)
			// Using hclwrite.Tokens as $ is not supported as part of string expression.
			// Adding a single "$" would result in "$$"
			resourceAttrs["name"] = hclwrite.Tokens{
				{Type: hclsyntax.TokenOQuote, Bytes: []byte(`"`)},
				{Type: hclsyntax.TokenIdent, Bytes: []byte(`${module.aws_eks_audit_log.filter_prefix}-${each.value}`)},
				{Type: hclsyntax.TokenCQuote, Bytes: []byte(`"`)},
			}
			resourceAttrs["log_group_name"] = hclwrite.Tokens{
				{Type: hclsyntax.TokenOQuote, Bytes: []byte(`"`)},
				{Type: hclsyntax.TokenIdent, Bytes: []byte(`/aws/eks/${each.value}/cluster`)},
				{Type: hclsyntax.TokenCQuote, Bytes: []byte(`"`)},
			}

			resourceAttrs["role_arn"] = lwgenerate.CreateSimpleTraversal([]string{"module", "aws_eks_audit_log", "cloudwatch_iam_role_arn"})
			resourceAttrs["filter_pattern"] = lwgenerate.CreateSimpleTraversal([]string{"module", "aws_eks_audit_log", "filter_pattern"})
			resourceAttrs["destination_arn"] = lwgenerate.CreateSimpleTraversal([]string{"module", "aws_eks_audit_log", "firehose_arn"})

			// the depends_on input must be in the following format [""]
			// In order to achieve this we can use TokensForTuple to build the tuple `[]`
			resourceAttrs["depends_on"] = hclwrite.TokensForTuple([]hclwrite.Tokens{
				hclwrite.TokensForTraversal(
					lwgenerate.CreateSimpleTraversal([]string{"module", "aws_eks_audit_log"}),
				),
			})

			lwCwSubscriptionFilter, err := lwgenerate.NewResource(
				"aws_cloudwatch_log_subscription_filter",
				fmt.Sprintf(
					"lw_cw_subscription_filter_%s",
					region),
				lwgenerate.HclResourceWithAttributesAndProviderDetails(
					resourceAttrs,
					[]string{fmt.Sprintf("aws.%s", region)},
				),
			).ToResourceBlock()

			if err != nil {
				return nil, err
			}

			blocks = append(blocks, lwCwSubscriptionFilter)
		}
	} else if len(args.ParsedRegionClusterMap) > 0 {
		// set no_cw_subscription_filter to false if we have only 1 region in the ParsedRegionClusterMap
		moduleAttrs["no_cw_subscription_filter"] = false
		for region, clusters := range args.ParsedRegionClusterMap {
			regionList = append(regionList, region)
			moduleAttrs["cluster_names"] = clusters
		}
	}

	moduleAttrs["cloudwatch_regions"] = regionList

	moduleDetails = append(moduleDetails,
		lwgenerate.HclModuleWithAttributes(moduleAttrs),
	)

	modDetails, err := lwgenerate.NewModule(
		"aws_eks_audit_log",
		lwgenerate.AwsEksAuditSource,
		moduleDetails...,
	).ToBlock()

	if err != nil {
		return nil, err
	}
	blocks = append(blocks, modDetails)

	return blocks, nil
}
