package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper for combining string expected values
func reqProviderAndRegion(extraInputs ...string) string {
	base := requiredProviders + "\n" + awsProvider
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

func TestGenerationCloudTrail(t *testing.T) {
	hcl, err := NewTerraform("us-east-2", false, true).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithoutConfig), hcl)
}

func TestGenerationConfig(t *testing.T) {
	hcl, err := NewTerraform("us-east-2", true, false).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportConfig), hcl)
}

func TestGenerationWithCustomAwsProfile(t *testing.T) {
	hcl, err := NewTerraform("us-east-2", false, true, WithAwsProfile("myprofile")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(
		t,
		fmt.Sprintf("%s\n%s\n%s", requiredProviders, awsProviderWithProfile, moduleImportCtWithoutConfig),
		hcl,
	)
}

func TestGenerationConfigAndCt(t *testing.T) {
	hcl, err := NewTerraform("us-east-2", true, true).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportConfig, moduleImportCtWithConfig), hcl)
}

func TestGenerationWithLaceworkProvider(t *testing.T) {
	hcl, err := NewTerraform("us-east-2", false, true, WithLaceworkProfile("test-profile")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(laceworkProvider, moduleImportCtWithoutConfig), hcl)
}

func TestGenerationWithLaceworkAccountID(t *testing.T) {
	hcl, err := NewTerraform("us-east-2", true, true, WithLaceworkAccountID("123456789")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportConfigWithLaceworkAccountID, moduleImportCtWithLaceworkAccountID), hcl)
}

func TestGenerationCloudtrailForceDestroyS3(t *testing.T) {
	data, err := createCloudtrail(&GenerateAwsTfConfigurationArgs{
		Cloudtrail:           true,
		ForceDestroyS3Bucket: true,
	})
	assert.Nil(t, err)
	assert.Equal(t,
		"bucket_force_destroy=true\n",
		string(data.Body().GetAttribute("bucket_force_destroy").BuildTokens(nil).Bytes()))
}

func TestGenerationCloudtrailConsolidatedTrail(t *testing.T) {
	data, err := createCloudtrail(&GenerateAwsTfConfigurationArgs{
		Cloudtrail:             true,
		ConsolidatedCloudtrail: true,
	})
	assert.Nil(t, err)
	assert.Equal(t,
		"consolidated_trail=true\n",
		string(data.Body().GetAttribute("consolidated_trail").BuildTokens(nil).Bytes()))
}

func TestGenerationCloudtrailExistingSns(t *testing.T) {
	existingSnsTopicArn := "arn:aws:sns:::foo"
	data, err := createCloudtrail(&GenerateAwsTfConfigurationArgs{
		Cloudtrail:          true,
		ExistingSnsTopicArn: existingSnsTopicArn,
	},
	)
	assert.Nil(t, err)
	assert.Equal(t,
		fmt.Sprintf("sns_topic_arn=\"%s\"\n", existingSnsTopicArn),
		string(data.Body().GetAttribute("sns_topic_arn").BuildTokens(nil).Bytes()))
	assert.Equal(t,
		"use_existing_sns_topic=true\n",
		string(data.Body().GetAttribute("use_existing_sns_topic").BuildTokens(nil).Bytes()))
}

func TestGenerationCloudtrailSnsWithEncryption(t *testing.T) {
	snsTopicName := "sns-topic-name"
	snsEncryptionArn := "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
	hcl, err := NewTerraform("us-east-2", false, true,
		WithSnsTopicName(snsTopicName),
		WithSnsTopicEncryptionEnabled(true),
		WithSnsTopicEncryptionKeyArn(snsEncryptionArn),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSnsWithoutConfig), hcl)
}

func TestGenerationCloudtrailSnsWithNoEncryption(t *testing.T) {
	snsTopicName := "sns-topic-name"
	hcl, err := NewTerraform("us-east-2", false, true,
		WithSnsTopicName(snsTopicName),
		WithSnsTopicEncryptionEnabled(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSnsNoConfigNoEncryption), hcl)
}

func TestGenerationCloudtrailSnsWithEncrytptionNotSet(t *testing.T) {
	snsTopicName := "sns-topic-name"
	hcl, err := NewTerraform("us-east-2", false, true,
		WithSnsTopicName(snsTopicName),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSnsNoConfigEncryptionNotSet), hcl)
}

func TestGenerationCloudtrailSqsWithEncryption(t *testing.T) {
	ssqQueueName := "sqs-queue-name"
	sqsEncryptionArn := "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
	hcl, err := NewTerraform("us-east-2", false, true,
		WithSqsQueueName(ssqQueueName),
		WithSqsEncryptionEnabled(true),
		WithSqsEncryptionKeyArn(sqsEncryptionArn),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSqsWithoutConfig), hcl)
}

func TestGenerationCloudtrailSqsWithNoEncryption(t *testing.T) {
	ssqQueueName := "sqs-queue-name"
	hcl, err := NewTerraform("us-east-2", false, true,
		WithSqsQueueName(ssqQueueName),
		WithSqsEncryptionEnabled(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSqsNoConfigNoEncryption), hcl)
}

func TestGenerationCloudtrailSqsWithWithEncryptionNotSet(t *testing.T) {
	ssqQueueName := "sqs-queue-name"
	hcl, err := NewTerraform("us-east-2", false, true,
		WithSqsQueueName(ssqQueueName),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSqsNoConfigEncryptionNotSet), hcl)
}

func TestGenerationCloudtrailAllEncryptionElementsSet(t *testing.T) {
	cloudTrailName := "cloudtrail-name"
	s3BucketName := "s3-bucket-name"
	snsTopicName := "sns-topic-name"
	ssqQueueName := "sqs-queue-name"
	encryptionArn := "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
	hcl, err := NewTerraform("us-east-2", false, true,
		WithCloudtrailName(cloudTrailName),
		WithBucketName(s3BucketName),
		WithBucketEncryptionEnabled(true),
		WithBucketSSEKeyArn(encryptionArn),
		WithSnsTopicName(snsTopicName),
		WithSnsTopicEncryptionEnabled(true),
		WithSnsTopicEncryptionKeyArn(encryptionArn),
		WithSqsQueueName(ssqQueueName),
		WithSqsEncryptionEnabled(true),
		WithSqsEncryptionKeyArn(encryptionArn),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithAllEncryptionSet), hcl)
}

func TestGenerationCloudtrailExistingBucket(t *testing.T) {
	existingBucketArn := "arn:aws:s3:::test-bucket-12345"
	data, err := createCloudtrail(&GenerateAwsTfConfigurationArgs{
		Cloudtrail:                  true,
		ExistingCloudtrailBucketArn: existingBucketArn,
	})
	assert.Nil(t, err)
	assert.Equal(t,
		"use_existing_cloudtrail=true\n",
		string(data.Body().GetAttribute("use_existing_cloudtrail").BuildTokens(nil).Bytes()))
	assert.Equal(t,
		fmt.Sprintf("bucket_arn=\"%s\"\n", existingBucketArn),
		string(data.Body().GetAttribute("bucket_arn").BuildTokens(nil).Bytes()))
}

func TestGenerationCloudtrailExistingRole(t *testing.T) {
	iamRoleArn := "arn:aws:iam::123456789012:role/test-role"
	iamRoleName := "test-role"
	extId := "1234567890123456"

	data, err := createCloudtrail(&GenerateAwsTfConfigurationArgs{
		Cloudtrail:      true,
		ExistingIamRole: NewExistingIamRoleDetails(iamRoleName, iamRoleArn, extId),
	})

	assert.Nil(t, err)
	assert.Equal(t,
		"use_existing_iam_role=true\n",
		string(data.Body().GetAttribute("use_existing_iam_role").BuildTokens(nil).Bytes()))
	assert.Equal(t,
		fmt.Sprintf("iam_role_name=\"%s\"\n", iamRoleName),
		string(data.Body().GetAttribute("iam_role_name").BuildTokens(nil).Bytes()))
	assert.Equal(t,
		fmt.Sprintf("iam_role_arn=\"%s\"\n", iamRoleArn),
		string(data.Body().GetAttribute("iam_role_arn").BuildTokens(nil).Bytes()))
	assert.Equal(t,
		fmt.Sprintf("iam_role_external_id=\"%s\"\n", extId),
		string(data.Body().GetAttribute("iam_role_external_id").BuildTokens(nil).Bytes()))
}

func TestConsolidatedCtWithMultipleAccounts(t *testing.T) {
	data, err := NewTerraform("us-east-2",
		true,
		true,
		WithAwsProfile("main"),
		WithSubaccounts(
			NewAwsSubAccount("subaccount1", "us-east-1"),
			NewAwsSubAccount("subaccount2", "us-east-2"),
		),
	).Generate()

	strippedData := strings.ReplaceAll(strings.ReplaceAll(data, "\n", ""), " ", "")
	assert.Nil(t, err)
	assert.Contains(t, strippedData, "provider\"aws\"{alias=\"main\"profile=\"main\"region=\"us-east-2\"}")
	assert.Contains(t, strippedData, "providers={aws=aws.main}")
	assert.Contains(t, strippedData, "module\"aws_config_subaccount1\"")
	assert.Contains(t, strippedData, "providers={aws=aws.subaccount1}")
	assert.Contains(t, strippedData, "provider\"aws\"{alias=\"subaccount1\"profile=\"subaccount1\"region=\"us-east-1\"}")
	assert.Contains(t, strippedData, "module\"aws_config_subaccount2\"")
	assert.Contains(t, strippedData, "providers={aws=aws.subaccount2}")
	assert.Contains(t, strippedData, "provider\"aws\"{alias=\"subaccount2\"profile=\"subaccount2\"region=\"us-east-2\"}")
}

func TestGenerationFailureWithNoOptionsSet(t *testing.T) {
	data := &GenerateAwsTfConfigurationArgs{}
	_, err := data.Generate()
	assert.Error(t, err)
	assert.Equal(t, "invalid inputs: cloudtrail or config integration must be enabled", err.Error())
}

func TestGenerationFailureWithNoRegionSet(t *testing.T) {
	data := &GenerateAwsTfConfigurationArgs{Cloudtrail: true}
	_, err := data.Generate()
	assert.Error(t, err)
	assert.Equal(t, "invalid inputs: AWS region must be set", err.Error())
}

var iamErrorString = "invalid inputs: when using an existing IAM role, existing role ARN, name, and external ID all must be set"

func TestGenerationFailureWithIncompleteExistingIam(t *testing.T) {
	_, err := NewTerraform("us-east-2", false, true,
		UseExistingIamRole(&ExistingIamRoleDetails{Arn: "foo"})).Generate()
	assert.Error(t, err)
	assert.Equal(t, iamErrorString, err.Error())

	_, err = NewTerraform("us-east-2", false, true,
		UseExistingIamRole(&ExistingIamRoleDetails{Name: "foo"})).Generate()
	assert.Error(t, err)
	assert.Equal(t, iamErrorString, err.Error())

	_, err = NewTerraform("us-east-2", false, true,
		UseExistingIamRole(&ExistingIamRoleDetails{ExternalId: "foo"})).Generate()
	assert.Error(t, err)
	assert.Equal(t, iamErrorString, err.Error())
}

func TestGenerationPartialExistingIamValues(t *testing.T) {
	t.Run("partial existing iam roles should be detected", func(t *testing.T) {
		data := NewExistingIamRoleDetails("test", "", "foo")
		assert.True(t, data.IsPartial())
	})
	t.Run("emtpy existing iam roles should not be detected as partial", func(t *testing.T) {
		data := NewExistingIamRoleDetails("", "", "")
		assert.False(t, data.IsPartial())
	})
	t.Run("nil existing iam roles should not be detected as partial", func(t *testing.T) {
		data := ExistingIamRoleDetails{}
		assert.False(t, data.IsPartial())
	})
	t.Run("completed existing iam roles should not be detected as partial", func(t *testing.T) {
		data := NewExistingIamRoleDetails("test", "arn:partition:service:region:account-id:resource-id", "foo")
		assert.False(t, data.IsPartial())
	})
}

func TestGenerationCloudTrailS3BucketNotification(t *testing.T) {
	hcl, err := NewTerraform("us-east-2", false, true, WithS3BucketNotification(true)).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(
		t,
		reqProviderAndRegion(moduleImportCtWithS3BucketNotification),
		hcl,
	)
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

var awsProvider = `provider "aws" {
  region = "us-east-2"
}
`

var awsProviderWithProfile = `provider "aws" {
  profile = "myprofile"
  region  = "us-east-2"
}
`

var laceworkProvider = `provider "lacework" {
  profile = "test-profile"
}
`

var moduleImportCtWithConfig = `module "main_cloudtrail" {
  source                = "lacework/cloudtrail/aws"
  version               = "~> 2.0"
  iam_role_arn          = module.aws_config.iam_role_arn
  iam_role_external_id  = module.aws_config.external_id
  iam_role_name         = module.aws_config.iam_role_name
  use_existing_iam_role = true
}
`

var moduleImportCtWithSnsWithoutConfig = `module "main_cloudtrail" {
  source                       = "lacework/cloudtrail/aws"
  version                      = "~> 2.0"
  sns_topic_encryption_key_arn = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sns_topic_name               = "sns-topic-name"
}
`

var moduleImportCtWithSnsNoConfigNoEncryption = `module "main_cloudtrail" {
  source                        = "lacework/cloudtrail/aws"
  version                       = "~> 2.0"
  sns_topic_encryption_enabled  = false
  sns_topic_name                = "sns-topic-name"
}
`
var moduleImportCtWithSnsNoConfigEncryptionNotSet = `module "main_cloudtrail" {
  source         = "lacework/cloudtrail/aws"
  version        = "~> 2.0"
  sns_topic_name = "sns-topic-name"
}
`

var moduleImportCtWithSqsWithoutConfig = `module "main_cloudtrail" {
  source                 = "lacework/cloudtrail/aws"
  version                = "~> 2.0"
  sqs_encryption_key_arn = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sqs_queue_name         = "sqs-queue-name"
}
`

var moduleImportCtWithSqsNoConfigNoEncryption = `module "main_cloudtrail" {
  source                  = "lacework/cloudtrail/aws"
  version                 = "~> 2.0"
  sqs_encryption_enabled  = false
  sqs_queue_name          = "sqs-queue-name"
}
`

var moduleImportCtWithSqsNoConfigEncryptionNotSet = `module "main_cloudtrail" {
  source         = "lacework/cloudtrail/aws"
  version        = "~> 2.0"
  sqs_queue_name = "sqs-queue-name"
}
`

var moduleImportCtWithAllEncryptionSet = `module "main_cloudtrail" {
  source                       = "lacework/cloudtrail/aws"
  version                      = "~> 2.0"
  bucket_name                  = "s3-bucket-name"
  bucket_sse_key_arn           = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  cloudtrail_name              = "cloudtrail-name"
  sns_topic_encryption_key_arn = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sns_topic_name               = "sns-topic-name"
  sqs_encryption_key_arn       = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sqs_queue_name               = "sqs-queue-name"
}
`

var moduleImportCtWithoutConfig = `module "main_cloudtrail" {
  source  = "lacework/cloudtrail/aws"
  version = "~> 2.0"
}
`

var moduleImportConfig = `module "aws_config" {
  source  = "lacework/config/aws"
  version = "~> 0.5"
}
`

var moduleImportConfigWithLaceworkAccountID = `module "aws_config" {
  source                  = "lacework/config/aws"
  version                 = "~> 0.5"
  lacework_aws_account_id = "123456789"
}
`

var moduleImportCtWithLaceworkAccountID = `module "main_cloudtrail" {
  source                  = "lacework/cloudtrail/aws"
  version                 = "~> 2.0"
  iam_role_arn            = module.aws_config.iam_role_arn
  iam_role_external_id    = module.aws_config.external_id
  iam_role_name           = module.aws_config.iam_role_name
  lacework_aws_account_id = "123456789"
  use_existing_iam_role   = true
}
`

var moduleImportCtWithS3BucketNotification = `module "main_cloudtrail" {
  source                     = "lacework/cloudtrail/aws"
  version                    = "~> 2.0"
  use_s3_bucket_notification = true
}
`
