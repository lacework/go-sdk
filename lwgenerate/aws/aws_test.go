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
		WithSubAccounts(
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

var requiredProviders = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 0.12.2"
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
  version               = "~> 0.1"
  iam_role_arn          = module.aws_config.iam_role_arn
  iam_role_external_id  = module.aws_config.external_id
  iam_role_name         = module.aws_config.iam_role_name
  use_existing_iam_role = true
}
`

var moduleImportCtWithoutConfig = `module "main_cloudtrail" {
  source  = "lacework/cloudtrail/aws"
  version = "~> 0.1"
}
`

var moduleImportConfig = `module "aws_config" {
  source  = "lacework/config/aws"
  version = "~> 0.1"
}
`
