package aws

import (
	"fmt"
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

func TestGenerationAgentless(t *testing.T) {
	hcl, err := NewTerraform(
		false,
		true,
		false,
		false,
		WithAwsProfile("main"),
		WithAwsRegion("us-east-2"),
		WithAgentlessScanningAccounts(
			NewAwsSubAccount("scanning-1", "us-east-1", "scanning-1-us-east-1"),
		),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, moduleImportAgentless, hcl)
}

func TestGenerationAgentlessOrganization(t *testing.T) {
	hcl, err := NewTerraform(
		true,
		true,
		false,
		false,
		WithAwsProfile("main"),
		WithAwsRegion("us-east-2"),
		WithAgentlessManagementAccountID("123456789000"),
		WithAgentlessMonitoredAccountIDs([]string{"123456789001", "ou-abcd-12345678"}),
		WithAgentlessMonitoredAccounts(
			NewAwsSubAccount("monitored-account-1", "us-west-2", "monitored-account-1-us-west-2"),
		),
		WithAgentlessScanningAccounts(
			NewAwsSubAccount("scanning-1", "us-east-1", "scanning-1-us-east-1"),
			NewAwsSubAccount("scanning-2", "us-east-2", "scanning-2-us-east-2"),
		),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, moduleImportAgentlessOrganization, hcl)
}

func TestGenerationCloudTrail(t *testing.T) {
	hcl, err := NewTerraform(false, false, false, true, WithAwsRegion("us-east-2")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCloudtrail), hcl)
}

func TestGenerationConfig(t *testing.T) {
	hcl, err := NewTerraform(false, false, true, false, WithAwsRegion("us-east-2")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportConfig), hcl)
}

func TestGenerationConfigWithMultipleAccounts(t *testing.T) {
	hcl, err := NewTerraform(false, false, true, false,
		WithAwsProfile("main"),
		WithAwsRegion("us-east-2"),
		WithConfigAdditionalAccounts(
			NewAwsSubAccount("subaccount1", "us-east-1", "subaccount1-us-east-1"),
			NewAwsSubAccount("subaccount2", "us-east-2", "subaccount2-us-east-2"),
		),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, moduleImportConfigWithMultipleAccounts, hcl)
}

func TestGenerationWithAwsProfile(t *testing.T) {
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsProfile("myprofile"),
		WithAwsRegion("us-east-2"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(
		t,
		fmt.Sprintf("%s\n%s\n%s", requiredProviders, awsProviderWithProfile, moduleImportCloudtrail),
		hcl,
	)
}

func TestGenerationWithLaceworkProvider(t *testing.T) {
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithLaceworkProfile("test-profile"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(laceworkProvider, moduleImportCloudtrail), hcl)
}

func TestGenerationWithLaceworkAccountID(t *testing.T) {
	hcl, err := NewTerraform(
		false,
		false,
		true,
		true,
		WithAwsRegion("us-east-2"),
		WithLaceworkAccountID("123456789"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportConfigWithLaceworkAccountID, moduleImportCtWithLaceworkAccountID), hcl)
}

func TestGenerationCloudtrailConsolidated(t *testing.T) {
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
		Cloudtrail:                    true,
		CloudtrailUseExistingSNSTopic: true,
		ExistingSnsTopicArn:           existingSnsTopicArn,
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
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
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
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithSnsTopicName(snsTopicName),
		WithSnsTopicEncryptionEnabled(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSnsNoConfigNoEncryption), hcl)
}

func TestGenerationCloudtrailSnsWithEncrytptionNotSet(t *testing.T) {
	snsTopicName := "sns-topic-name"
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithSnsTopicName(snsTopicName),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSnsNoConfigEncryptionNotSet), hcl)
}

func TestGenerationCloudtrailSqsWithEncryption(t *testing.T) {
	ssqQueueName := "sqs-queue-name"
	sqsEncryptionArn := "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
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
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithSqsQueueName(ssqQueueName),
		WithSqsEncryptionEnabled(false),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSqsNoConfigNoEncryption), hcl)
}

func TestGenerationCloudtrailSqsWithWithEncryptionNotSet(t *testing.T) {
	ssqQueueName := "sqs-queue-name"
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithSqsQueueName(ssqQueueName),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProviderAndRegion(moduleImportCtWithSqsNoConfigEncryptionNotSet), hcl)
}

func TestGenerationCloudtrailAllEncryptionElementsSet(t *testing.T) {
	s3BucketName := "s3-bucket-name"
	snsTopicName := "sns-topic-name"
	ssqQueueName := "sqs-queue-name"
	encryptionArn := "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
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
		CloudtrailUseExistingS3:     true,
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

func TestGenerationFailureWithNoOptionsSet(t *testing.T) {
	data := &GenerateAwsTfConfigurationArgs{}
	_, err := data.Generate()
	assert.Error(t, err)
	assert.Equal(t, "invalid inputs: Agentless, CloudTrail or Config integration must be enabled", err.Error())
}

func TestGenerationFailureWithNoRegionSet(t *testing.T) {
	data := &GenerateAwsTfConfigurationArgs{Cloudtrail: true}
	_, err := data.Generate()
	assert.Error(t, err)
	assert.Equal(t, "invalid inputs: Main AWS account region must be set", err.Error())
}

func TestGenerationFailureWithIncompleteExistingIam(t *testing.T) {
	iamErrorString := "invalid inputs: when using an existing IAM role, existing role ARN, name, and external ID all must be set"

	_, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithExistingIamRole(&ExistingIamRoleDetails{Arn: "foo"}),
	).Generate()
	assert.Error(t, err)
	assert.Equal(t, iamErrorString, err.Error())

	_, err = NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithExistingIamRole(&ExistingIamRoleDetails{Name: "foo"}),
	).Generate()
	assert.Error(t, err)
	assert.Equal(t, iamErrorString, err.Error())

	_, err = NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithExistingIamRole(&ExistingIamRoleDetails{ExternalId: "foo"}),
	).Generate()
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
	hcl, err := NewTerraform(false, false, false, true,
		WithAwsRegion("us-east-2"),
		WithS3BucketNotification(true),
	).Generate()
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
  alias  = "main"
  region = "us-east-2"
}
`

var awsProviderWithProfile = `provider "aws" {
  alias   = "main"
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
  version               = "~> 2.7"
  iam_role_arn          = module.aws_config.iam_role_arn
  iam_role_external_id  = module.aws_config.external_id
  iam_role_name         = module.aws_config.iam_role_name
  use_existing_iam_role = true

  providers = {
    aws = aws.main
  }
}
`

var moduleImportCtWithSnsWithoutConfig = `module "main_cloudtrail" {
  source                       = "lacework/cloudtrail/aws"
  version                      = "~> 2.7"
  sns_topic_encryption_key_arn = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sns_topic_name               = "sns-topic-name"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportCtWithSnsNoConfigNoEncryption = `module "main_cloudtrail" {
  source                        = "lacework/cloudtrail/aws"
  version                       = "~> 2.7"
  sns_topic_encryption_enabled  = false
  sns_topic_name                = "sns-topic-name"

  providers = {
    aws = aws.main
  }
}
`
var moduleImportCtWithSnsNoConfigEncryptionNotSet = `module "main_cloudtrail" {
  source         = "lacework/cloudtrail/aws"
  version        = "~> 2.7"
  sns_topic_name = "sns-topic-name"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportCtWithSqsWithoutConfig = `module "main_cloudtrail" {
  source                 = "lacework/cloudtrail/aws"
  version                = "~> 2.7"
  sqs_encryption_key_arn = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sqs_queue_name         = "sqs-queue-name"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportCtWithSqsNoConfigNoEncryption = `module "main_cloudtrail" {
  source                  = "lacework/cloudtrail/aws"
  version                 = "~> 2.7"
  sqs_encryption_enabled  = false
  sqs_queue_name          = "sqs-queue-name"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportCtWithSqsNoConfigEncryptionNotSet = `module "main_cloudtrail" {
  source         = "lacework/cloudtrail/aws"
  version        = "~> 2.7"
  sqs_queue_name = "sqs-queue-name"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportCtWithAllEncryptionSet = `module "main_cloudtrail" {
  source                       = "lacework/cloudtrail/aws"
  version                      = "~> 2.7"
  bucket_name                  = "s3-bucket-name"
  bucket_sse_key_arn           = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sns_topic_encryption_key_arn = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sns_topic_name               = "sns-topic-name"
  sqs_encryption_key_arn       = "arn:aws:kms:us-west-2:249446771485:key/2537e820-be82-4ded-8dca-504e199b0903"
  sqs_queue_name               = "sqs-queue-name"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportAgentless = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
}

provider "aws" {
  alias   = "main"
  profile = "main"
  region  = "us-east-2"
}

provider "aws" {
  alias   = "scanning-1-us-east-1"
  profile = "scanning-1"
  region  = "us-east-1"
}

module "lacework_aws_agentless_scanning_global" {
  source   = "lacework/agentless-scanning/aws"
  version  = "~> 0.6"
  global   = true
  regional = true

  providers = {
    aws = aws.main
  }
}

module "lacework_aws_agentless_scanning_region_scanning-1-us-east-1" {
  source                  = "lacework/agentless-scanning/aws"
  version                 = "~> 0.6"
  global_module_reference = module.lacework_aws_agentless_scanning_global
  regional                = true

  providers = {
    aws = aws.scanning-1-us-east-1
  }
}
`

var moduleImportAgentlessOrganization = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
}

provider "aws" {
  alias   = "main"
  profile = "main"
  region  = "us-east-2"
}

provider "aws" {
  alias   = "monitored-account-1-us-west-2"
  profile = "monitored-account-1"
  region  = "us-west-2"
}

provider "aws" {
  alias   = "scanning-1-us-east-1"
  profile = "scanning-1"
  region  = "us-east-1"
}

provider "aws" {
  alias   = "scanning-2-us-east-2"
  profile = "scanning-2"
  region  = "us-east-2"
}

module "lacework_aws_agentless_management_scanning_role" {
  source                  = "lacework/agentless-scanning/aws"
  version                 = "~> 0.6"
  global_module_reference = module.lacework_aws_agentless_scanning_global
  snapshot_role           = true

  providers = {
    aws = aws.main
  }
}

module "lacework_aws_agentless_scanning_global" {
  source  = "lacework/agentless-scanning/aws"
  version = "~> 0.6"
  global  = true
  organization = {
    management_account = "123456789000"
    monitored_accounts = ["123456789001", "ou-abcd-12345678"]
  }
  regional = true

  providers = {
    aws = aws.scanning-1-us-east-1
  }
}

module "lacework_aws_agentless_scanning_region_scanning-2-us-east-2" {
  source                  = "lacework/agentless-scanning/aws"
  version                 = "~> 0.6"
  global_module_reference = module.lacework_aws_agentless_scanning_global
  regional                = true

  providers = {
    aws = aws.scanning-2-us-east-2
  }
}

module "lacework_aws_agentless_monitored_scanning_role_monitored-account-1-us-west-2" {
  source                  = "lacework/agentless-scanning/aws"
  version                 = "~> 0.6"
  global_module_reference = module.lacework_aws_agentless_scanning_global
  snapshot_role           = true

  providers = {
    aws = aws.monitored-account-1-us-west-2
  }
}

resource "aws_cloudformation_stack_set" "snapshot_role" {
  capabilities = ["CAPABILITY_NAMED_IAM"]
  description  = "Lacework AWS Agentless Workload Scanning Organization Roles"
  name         = "lacework-agentless-scanning-stackset"
  parameters = {
    ECSTaskRoleArn     = module.lacework_aws_agentless_scanning_global.agentless_scan_ecs_task_role_arn
    ExternalId         = module.lacework_aws_agentless_scanning_global.external_id
    ResourceNamePrefix = module.lacework_aws_agentless_scanning_global.prefix
    ResourceNameSuffix = module.lacework_aws_agentless_scanning_global.suffix
  }
  permission_model = "SERVICE_MANAGED"
  template_url     = "https://agentless-workload-scanner.s3.amazonaws.com/cloudformation-lacework/latest/snapshot-role.json"

  provider = aws.main

  auto_deployment {
    enabled                          = true
    retain_stacks_on_account_removal = false
  }

  lifecycle {
    ignore_changes = [administration_role_arn]
  }
}

resource "aws_cloudformation_stack_set_instance" "snapshot_role" {
  stack_set_name = aws_cloudformation_stack_set.snapshot_role.name

  provider = aws.main

  deployment_targets {
    organizational_unit_ids = ["ou-abcd-12345678"]
  }
}
`

var moduleImportCloudtrail = `module "main_cloudtrail" {
  source  = "lacework/cloudtrail/aws"
  version = "~> 2.7"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportConfig = `module "aws_config" {
  source  = "lacework/config/aws"
  version = "~> 0.5"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportConfigWithMultipleAccounts = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
}

provider "aws" {
  alias   = "main"
  profile = "main"
  region  = "us-east-2"
}

provider "aws" {
  alias   = "subaccount1-us-east-1"
  profile = "subaccount1"
  region  = "us-east-1"
}

provider "aws" {
  alias   = "subaccount2-us-east-2"
  profile = "subaccount2"
  region  = "us-east-2"
}

module "aws_config" {
  source  = "lacework/config/aws"
  version = "~> 0.5"

  providers = {
    aws = aws.main
  }
}

module "aws_config_subaccount1-us-east-1" {
  source  = "lacework/config/aws"
  version = "~> 0.5"

  providers = {
    aws = aws.subaccount1-us-east-1
  }
}

module "aws_config_subaccount2-us-east-2" {
  source  = "lacework/config/aws"
  version = "~> 0.5"

  providers = {
    aws = aws.subaccount2-us-east-2
  }
}
`

var moduleImportConfigWithLaceworkAccountID = `module "aws_config" {
  source                  = "lacework/config/aws"
  version                 = "~> 0.5"
  lacework_aws_account_id = "123456789"

  providers = {
    aws = aws.main
  }
}
`

var moduleImportCtWithLaceworkAccountID = `module "main_cloudtrail" {
  source                  = "lacework/cloudtrail/aws"
  version                 = "~> 2.7"
  iam_role_arn            = module.aws_config.iam_role_arn
  iam_role_external_id    = module.aws_config.external_id
  iam_role_name           = module.aws_config.iam_role_name
  lacework_aws_account_id = "123456789"
  use_existing_iam_role   = true

  providers = {
    aws = aws.main
  }
}
`

var moduleImportCtWithS3BucketNotification = `module "main_cloudtrail" {
  source                     = "lacework/cloudtrail/aws"
  version                    = "~> 2.7"
  use_s3_bucket_notification = true

  providers = {
    aws = aws.main
  }
}
`
