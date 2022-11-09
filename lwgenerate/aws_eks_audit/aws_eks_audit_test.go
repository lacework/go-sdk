package aws_eks_audit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper for combining string expected values
func reqProviderAndRegion(extraInputs ...string) string {
	base := requiredProviders
	countInputs := len(extraInputs)
	for i, e := range extraInputs {
		if i < countInputs {
			base = base + "\n" + e
		}

		if i >= countInputs {
			base = base + e
		}
	}
	return base
}

func TestGenerationEksSingleRegion(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap)).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleSingleRegionBasic), hcl)
}

func TestGenerationEksMultiRegion(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	clusterMap["us-east-2"] = []string{"cluster3"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap)).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(multiRegionBasic), hcl)
}

func TestGenerationEksFailureWithNoOptionsSet(t *testing.T) {
	data := &GenerateAwsEksAuditTfConfigurationArgs{}
	_, err := data.Generate()
	assert.Error(t, err)
	assert.Equal(t, "invalid inputs: At least one region with a list of clusters must be set", err.Error())
}

func TestGenerationEksFailureSingleRegionNoClusters(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{}
	_, err := NewTerraform(WithParsedRegionClusterMap(clusterMap)).Generate()
	assert.Error(t, err)
	assert.Equal(t, "invalid inputs: At least one cluster must be supplied per region", err.Error())
}

func TestGenerationEksEnableBucketForceDestroy(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableBucketForceDestroy(),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "bucket_force_destroy=true")
}

func TestGenerationEksWithValidBucketLifecycleExpirationDays(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithBucketLifecycleExpirationDays(10),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "bucket_lifecycle_expiration_days=10")
}

func TestGenerationEksWithInvalidBucketLifecycleExpirationDays(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithBucketLifecycleExpirationDays(-1),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.NotContains(t, hcl, "bucket_lifecycle_expiration_days")
}

func TestGenerationEksBucketWithNoBucketVersioning(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableBucketVersioning(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "bucket_versioning_enabled=false")
}

func TestGenerationEksBucketWithNoBucketVersioningMfaDeleteEnabled(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableBucketVersioning(false),
		EnableBucketMfaDelete(),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "bucket_versioning_enabled=false")
	assert.NotContains(t, hcl, "bucket_enable_mfa_delete")
}

func TestGenerationEksBucketWithNoEncryption(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableBucketEncryption(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "bucket_encryption_enabled=false")
}

func TestGenerationEksBucketWithNoEncryptionAndKeyArn(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableBucketEncryption(false),
		WithBucketSseKeyArn("arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "bucket_encryption_enabled=false")
	assert.NotContains(t, hcl, "bucket_key_arn")
}

func TestGenerationEksWithBucketEncryptionKeyArnAndWithBucketSseAlgorithm(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithBucketSseKeyArn("arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"),
		WithBucketSseAlgorithm("aws:kms"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl,
		"bucket_key_arn=\"arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903\"")
	assert.Contains(t, strippedHcl, "bucket_sse_algorithm=\"aws:kms\"")
}

func TestGenerationEksFirehoseWithNoEncryption(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableFirehoseEncryption(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "kinesis_firehose_encryption_enabled=false")
}

func TestGenerationEksFirehoseWithNoEncryptionAndKeyArn(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableFirehoseEncryption(false),
		WithFirehoseEncryptionKeyArn("arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "kinesis_firehose_encryption_enabled=false")
	assert.NotContains(t, hcl, "kinesis_firehose_key_arn")
}

func TestGenerationEksWithFirehoseEncryptionKeyArn(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithFirehoseEncryptionKeyArn("arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl,
		"kinesis_firehose_key_arn=\"arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903\"")
}

func TestGenerationEksSnsWithNoEncryption(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableSnsTopicEncryption(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "sns_topic_encryption_enabled=false")
}

func TestGenerationEksSnsWithNoEncryptionAndKeyArn(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableSnsTopicEncryption(false),
		WithSnsTopicEncryptionKeyArn("arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "sns_topic_encryption_enabled=false")
	assert.NotContains(t, hcl, "sns_topic_key_arn")
}

func TestGenerationEksWithSnsTopicEncryptionKeyArn(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithSnsTopicEncryptionKeyArn("arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl,
		"sns_topic_key_arn=\"arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903\"")
}

func TestGenerationEksWithExistingCloudWatchIamRoleArn(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithExistingCloudWatchIamRoleArn("arn:aws:iam::249446771485:role/2537e820-cloudwatch-role"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl,
		"cloudwatch_iam_role_arn=\"arn:aws:iam::249446771485:role/2537e820-cloudwatch-role\"")
	assert.Contains(t, strippedHcl,
		"use_existing_cloudwatch_iam_role=true")
}

func TestGenerationEksWithExistingFirehoseIamRoleArn(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithExistingFirehoseIamRoleArn("arn:aws:iam::249446771485:role/2537e820-firehose-role"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl,
		"firehose_iam_role_arn=\"arn:aws:iam::249446771485:role/2537e820-firehose-role\"")
	assert.Contains(t, strippedHcl,
		"use_existing_firehose_iam_role=true")
}

var iamErrorString = "invalid inputs: when using an existing cross account IAM role, existing role ARN and external ID all must be set"

func TestGenerationFailureWithIncompleteExistingCrossAccountIam(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	_, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithExistingCrossAccountIamRole(&ExistingCrossAccountIamRoleDetails{Arn: "foo"})).Generate()
	assert.Error(t, err)
	assert.Equal(t, iamErrorString, err.Error())

	_, err = NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithExistingCrossAccountIamRole(&ExistingCrossAccountIamRoleDetails{ExternalId: "foo"})).Generate()
	assert.Error(t, err)
	assert.Equal(t, iamErrorString, err.Error())
}

func TestGenerationPartialExistingCrossAccountIamValues(t *testing.T) {
	t.Run("partial existing iam roles should be detected", func(t *testing.T) {
		data := NewExistingCrossAccountIamRoleDetails("test", "")
		assert.True(t, data.IsPartial())
	})
	t.Run("emtpy existing iam roles should not be detected as partial", func(t *testing.T) {
		data := NewExistingCrossAccountIamRoleDetails("", "")
		assert.False(t, data.IsPartial())
	})
	t.Run("nil existing iam roles should not be detected as partial", func(t *testing.T) {
		data := ExistingCrossAccountIamRoleDetails{}
		assert.False(t, data.IsPartial())
	})
	t.Run("completed existing iam roles should not be detected as partial", func(t *testing.T) {
		data := NewExistingCrossAccountIamRoleDetails(
			"arn:partition:service:region:account-id:resource-id", "test")
		assert.False(t, data.IsPartial())
	})
}

func TestGenerationEksWithValidKmsKeyDeletionDays(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithKmsKeyDeletionDays(10),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "kms_key_deletion_days=10")
}

func TestGenerationEksWithInvalidKmsKeyDeletionDays(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		WithKmsKeyDeletionDays(6),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.NotContains(t, hcl, "kms_key_deletion_days")
}

func TestGenerationEksEnableKmsKeyMultiRegionTrueWithSingleRegion(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableKmsKeyMultiRegion(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "kms_key_multi_region=false")
}

func TestGenerationEksEnableKmsKeyMultiRegionTrueWithMultiRegion(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	clusterMap["us-east-2"] = []string{"cluster3"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableKmsKeyMultiRegion(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.NotContains(t, hcl, "kms_key_multi_region")
}

func TestGenerationEksEnableKmsKeyMultiRegionFalse(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableKmsKeyMultiRegion(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "kms_key_multi_region=false")
}

func TestGenerationEksEnableKmsKeyRotationTrue(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableKmsKeyRotation(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.NotContains(t, hcl, "kms_key_rotation")
}

func TestGenerationEksEnableKmsKeyRotationFalse(t *testing.T) {
	clusterMap := make(map[string][]string)
	clusterMap["us-east-1"] = []string{"cluster1", "cluster2"}
	hcl, err := NewTerraform(WithParsedRegionClusterMap(clusterMap),
		EnableKmsKeyRotation(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	strippedHcl := strings.ReplaceAll(hcl, " ", "")
	assert.Contains(t, strippedHcl, "kms_key_rotation=false")
}

var requiredProviders = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
}
`

var moduleSingleRegionBasic = `provider "aws" {
  region = "us-east-1"
}

module "aws_eks_audit_log" {
  source                    = "lacework/eks-audit-log/aws"
  version                   = "~> 0.4"
  cloudwatch_regions        = ["us-east-1"]
  cluster_names             = ["cluster1", "cluster2"]
  kms_key_multi_region      = false
  no_cw_subscription_filter = false
}
`

var multiRegionBasic = `provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"
}

provider "aws" {
  alias  = "us-east-2"
  region = "us-east-2"
}

resource "aws_cloudwatch_log_subscription_filter" "lw_cw_subscription_filter_us-east-1" {
  depends_on      = [module.aws_eks_audit_log]
  destination_arn = module.aws_eks_audit_log.firehose_arn
  filter_pattern  = module.aws_eks_audit_log.filter_pattern
  for_each        = toset(["cluster1", "cluster2"])
  log_group_name  = "/aws/eks/${each.value}/cluster"
  name            = "${module.aws_eks_audit_log.filter_prefix}-${each.value}"
  role_arn        = module.aws_eks_audit_log.cloudwatch_iam_role_arn

  provider = aws.us-east-1
}

resource "aws_cloudwatch_log_subscription_filter" "lw_cw_subscription_filter_us-east-2" {
  depends_on      = [module.aws_eks_audit_log]
  destination_arn = module.aws_eks_audit_log.firehose_arn
  filter_pattern  = module.aws_eks_audit_log.filter_pattern
  for_each        = toset(["cluster3"])
  log_group_name  = "/aws/eks/${each.value}/cluster"
  name            = "${module.aws_eks_audit_log.filter_prefix}-${each.value}"
  role_arn        = module.aws_eks_audit_log.cloudwatch_iam_role_arn

  provider = aws.us-east-2
}

module "aws_eks_audit_log" {
  source                    = "lacework/eks-audit-log/aws"
  version                   = "~> 0.4"
  cloudwatch_regions        = ["us-east-1", "us-east-2"]
  no_cw_subscription_filter = true
}
`
