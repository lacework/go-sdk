//go:build !windows && generation

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/stretchr/testify/assert"
)

const (
	awsPath = "/lacework/aws/"
)

// Test all intergrations in non-interactive mode
func TestGenerationAwsNoninteractive(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	assumeRole := "arn:aws:iam::123456789012:role/S3Access"
	cloudtrailName := "cloudtrail-integration-name"
	s3BucketArn := "arn:aws:s3:::bucket-name"
	snsTopicArn := "arn:aws:sns:us-east-2:249446771485:topic-name"
	sqsQueueName := "sqs-queue-name"
	kmsArn := "arn:aws:kms:us-west-2:249446771485:key/203a7566-41eb-42dc-8cc3-51800a87defe"
	iamRoleName := "test-iamrole"
	iamRoleArn := "arn:aws:iam::123456789012:role/application_abc/component_xyz/abc_role"
	iamRoleExtId := "123456"

	laceworkAccount := "lw-account"
	laceworkSubaccount := "lw-subaccount"
	laceworkAccessKeyId := "lw-access-key-id"
	laceworkSecretKey := "lw-secret-key"
	organizationId := "aws-org-id"
	organizationUnit := "aws-org-unit"
	cfResourcePrefix := "cd-resource-prefix"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
		"--noninteractive",
		"--aws_organization",
		"--agentless",
		"--config",
		"--cloudtrail",
		"--aws_profile",
		"main",
		"--aws_region",
		"us-east-2",
		"--aws_assume_role",
		assumeRole,
		"--aws_subaccount",
		"subaccount-1:us-west-1,subaccount-2:us-west-2",
		"--lacework_aws_account_id",
		"123456789",
		"--agentless_management_account_id",
		"123456789000",
		"--agentless_monitored_account_ids",
		"123456789000,ou-abcd-12345678,r-abcd",
		"--agentless_monitored_accounts",
		"monitored-1:us-west-1,monitored-2:us-west-2",
		"--agentless_scanning_accounts",
		"scanning-1:us-east-1,scanning-2:us-east-2",
		"--config_lacework_account",
		laceworkAccount,
		"--config_lacework_sub_account",
		laceworkSubaccount,
		"--config_lacework_access_key_id",
		laceworkAccessKeyId,
		"--config_lacework_secret_key",
		laceworkSecretKey,
		"--config_organization_id",
		organizationId,
		"--config_organization_unit",
		organizationUnit,
		"--config_cf_resource_prefix",
		cfResourcePrefix,
		"--consolidated_cloudtrail",
		"--cloudtrail_org_account_mapping",
		"{\"default_lacework_account\":\"main\", \"mapping\": [{ \"aws_accounts\": [\"123456789011\"], \"lacework_account\": \"subaccount-1\"}]}",
		"--cloudtrail_name",
		cloudtrailName,
		"--existing_bucket_arn",
		s3BucketArn,
		"--bucket_encryption_enabled",
		"--bucket_sse_key_arn",
		kmsArn,
		"--existing_sns_topic_arn",
		snsTopicArn,
		"--sns_topic_encryption_enabled",
		"--sns_topic_encryption_key_arn",
		kmsArn,
		"--sqs_queue_name",
		sqsQueueName,
		"--sqs_encryption_enabled",
		"--sqs_encryption_key_arn",
		kmsArn,
		"--existing_iam_role_name",
		iamRoleName,
		"--existing_iam_role_arn",
		iamRoleArn,
		"--existing_iam_role_externalid",
		iamRoleExtId,
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	orgAccountMappings := aws.OrgAccountMapping{
		DefaultLaceworkAccount: "main",
		Mapping: []aws.OrgAccountMap{
			{
				LaceworkAccount: "subaccount-1",
				AwsAccounts:     []string{"123456789011"},
			},
		},
	}

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(true, true, true, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithAwsAssumeRole(assumeRole),
		aws.WithLaceworkAccountID("123456789"),
		aws.WithAgentlessManagementAccountID("123456789000"),
		aws.WithAgentlessMonitoredAccountIDs([]string{"123456789000", "ou-abcd-12345678", "r-abcd"}),
		aws.WithAgentlessMonitoredAccounts(
			aws.NewAwsSubAccount("monitored-1", "us-west-1", "monitored-1-us-west-1"),
			aws.NewAwsSubAccount("monitored-2", "us-west-2", "monitored-2-us-west-2"),
		),
		aws.WithAgentlessScanningAccounts(
			aws.NewAwsSubAccount("scanning-1", "us-east-1", "scanning-1-us-east-1"),
			aws.NewAwsSubAccount("scanning-2", "us-east-2", "scanning-2-us-east-2"),
		),
		aws.WithConfigAdditionalAccounts(
			aws.NewAwsSubAccount("subaccount-1", "us-west-1", "subaccount-1-us-west-1"),
			aws.NewAwsSubAccount("subaccount-2", "us-west-2", "subaccount-2-us-west-2"),
		),
		aws.WithConfigOrgLWAccount(laceworkAccount),
		aws.WithConfigOrgLWSubaccount(laceworkSubaccount),
		aws.WithConfigOrgLWAccessKeyId(laceworkAccessKeyId),
		aws.WithConfigOrgLWSecretKey(laceworkSecretKey),
		aws.WithConfigOrgId(organizationId),
		aws.WithConfigOrgUnit(organizationUnit),
		aws.WithConfigOrgCfResourcePrefix(cfResourcePrefix),
		aws.WithConsolidatedCloudtrail(true),
		aws.WithOrgAccountMappings(orgAccountMappings),
		aws.WithCloudtrailUseExistingS3(true),
		aws.WithCloudtrailName(cloudtrailName),
		aws.WithExistingCloudtrailBucketArn(s3BucketArn),
		aws.WithCloudtrailUseExistingSNSTopic(true),
		aws.WithExistingSnsTopicArn(snsTopicArn),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionKeyArn(kmsArn),
		aws.WithSqsQueueName(sqsQueueName),
		aws.WithSqsEncryptionEnabled(true),
		aws.WithSqsEncryptionKeyArn(kmsArn),
		aws.WithExistingIamRole(
			aws.NewExistingIamRoleDetails(iamRoleName, iamRoleArn, iamRoleExtId),
		),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Agentless single account integration
func TestGenerationAwsAgentless(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountAddMore, "y"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountProfile, "scanning-1"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountRegion, "us-east-1"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountAddMore, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, true, false, false,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithAgentlessScanningAccounts(
			aws.NewAwsSubAccount("scanning-1", "us-east-1", "scanning-1-us-east-1"),
		),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Agentless organization integration
func TestGenerationAwsAgentlessOrganization(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "y"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAgentlessManagementAccountID, "123456789000"},
				MsgRsp{cmd.QuestionAgentlessMonitoredAccountIDs, "123456789000,ou-abcd-12345678,r-abcd"},
				MsgRsp{cmd.QuestionAgentlessMonitoredAccountProfile, "monitored-1"},
				MsgRsp{cmd.QuestionAgentlessMonitoredAccountRegion, "us-west-1"},
				MsgRsp{cmd.QuestionAgentlessMonitoredAccountAddMore, "y"},
				MsgRsp{cmd.QuestionAgentlessMonitoredAccountProfile, "monitored-2"},
				MsgRsp{cmd.QuestionAgentlessMonitoredAccountRegion, "us-west-2"},
				MsgRsp{cmd.QuestionAgentlessMonitoredAccountAddMore, "n"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountProfile, "scanning-1"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountRegion, "us-east-1"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountAddMore, "y"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountProfile, "scanning-2"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountRegion, "us-east-2"},
				MsgRsp{cmd.QuestionAgentlessScanningAccountAddMore, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(true, true, false, false,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithAgentlessManagementAccountID("123456789000"),
		aws.WithAgentlessMonitoredAccountIDs([]string{"123456789000", "ou-abcd-12345678", "r-abcd"}),
		aws.WithAgentlessMonitoredAccounts(
			aws.NewAwsSubAccount("monitored-1", "us-west-1", "monitored-1-us-west-1"),
			aws.NewAwsSubAccount("monitored-2", "us-west-2", "monitored-2-us-west-2"),
		),
		aws.WithAgentlessScanningAccounts(
			aws.NewAwsSubAccount("scanning-1", "us-east-1", "scanning-1-us-east-1"),
			aws.NewAwsSubAccount("scanning-2", "us-east-2", "scanning-2-us-east-2"),
		),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Config integration
func TestGenerationAwsConfig(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "y"},
				MsgRsp{cmd.QuestionConfigAdditionalAccountAddMore, "y"},
				MsgRsp{cmd.QuestionConfigAdditionalAccountProfile, "subaccount-1"},
				MsgRsp{cmd.QuestionConfigAdditionalAccountRegion, "us-west-1"},
				MsgRsp{cmd.QuestionConfigAdditionalAccountAddMore, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, true, false,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithConfigAdditionalAccounts(
			aws.NewAwsSubAccount("subaccount-1", "us-west-1", "subaccount-1-us-west-1"),
		),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Config organization integration
func TestGenerationAwsConfigOrganization(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	laceworkAccount := "lw-account"
	laceworkSubaccount := "lw-subaccount"
	laceworkAccessKeyId := "lw-access-key-id"
	laceworkSecretKey := "lw-secret-key"
	organizationId := "aws-org-id"
	organizationUnit := "aws-org-unit"
	cfResourcePrefix := "cd-resource-prefix"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "y"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "y"},
				MsgRsp{cmd.QuestionConfigOrgLWAccount, laceworkAccount},
				MsgRsp{cmd.QuestionConfigOrgLWSubaccount, laceworkSubaccount},
				MsgRsp{cmd.QuestionConfigOrgLWAccessKeyId, laceworkAccessKeyId},
				MsgRsp{cmd.QuestionConfigOrgLWSecretKey, laceworkSecretKey},
				MsgRsp{cmd.QuestionConfigOrgId, organizationId},
				MsgRsp{cmd.QuestionConfigOrgUnit, organizationUnit},
				MsgRsp{cmd.QuestionConfigOrgCfResourcePrefix, cfResourcePrefix},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(true, false, true, false,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithConfigOrgLWAccount(laceworkAccount),
		aws.WithConfigOrgLWSubaccount(laceworkSubaccount),
		aws.WithConfigOrgLWAccessKeyId(laceworkAccessKeyId),
		aws.WithConfigOrgLWSecretKey(laceworkSecretKey),
		aws.WithConfigOrgId(organizationId),
		aws.WithConfigOrgUnit(organizationUnit),
		aws.WithConfigOrgCfResourcePrefix(cfResourcePrefix),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test CloudTrail single account integration
func TestGenerationAwsCloudtrail(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailUseConsolidated, "y"},
				MsgRsp{cmd.QuestionCloudtrailAdvanced, "n"},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, false, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithConsolidatedCloudtrail(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test CloudTrail organization integration
func TestGenerationAwsCloudtrailOrganization(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "y"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailUseConsolidated, "y"},
				MsgRsp{cmd.QuestionCloudtrailAdvanced, "y"},
				MsgMenu{cmd.OptCloudtrailOrg, 0},
				MsgRsp{cmd.QuestionCloudtrailOrgAccountMappingsDefaultLWAccount, "main"},
				MsgRsp{cmd.QuestionCloudtrailOrgAccountMappingsLWAccount, "subaccount-1"},
				MsgMultilineRsp{cmd.QuestionCloudtrailOrgAccountMappingsAwsAccounts, []string{"123456789001"}},
				MsgRsp{cmd.QuestionCloudtrailOrgAccountMappingsAnotherAddMore, "y"},
				MsgRsp{cmd.QuestionCloudtrailOrgAccountMappingsLWAccount, "subaccount-2"},
				MsgMultilineRsp{cmd.QuestionCloudtrailOrgAccountMappingsAwsAccounts, []string{"123456789002"}},
				MsgRsp{cmd.QuestionCloudtrailOrgAccountMappingsAnotherAddMore, "n"},
				MsgMenu{cmd.OptCloudtrailDone, 5},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	orgAccountMappings := aws.OrgAccountMapping{
		DefaultLaceworkAccount: "main",
		Mapping: []aws.OrgAccountMap{
			{
				LaceworkAccount: "subaccount-1",
				AwsAccounts:     []string{"123456789001"},
			},
			{
				LaceworkAccount: "subaccount-2",
				AwsAccounts:     []string{"123456789002"},
			},
		},
	}

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(true, false, false, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithConsolidatedCloudtrail(true),
		aws.WithOrgAccountMappings(orgAccountMappings),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test CloutTrail integration with existing S3 bucket
func TestGenerationAwsCloudtrailWithExistingS3Bucket(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	cloudtrailName := "cloudtrail-integration-name"
	s3BucketArn := "arn:aws:s3:::bucket-name"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailUseConsolidated, "n"},
				MsgRsp{cmd.QuestionCloudtrailAdvanced, "y"},
				MsgMenu{cmd.OptCloudtrailS3, 0},
				MsgRsp{cmd.QuestionCloudtrailUseExistingS3, "y"},
				MsgRsp{cmd.QuestionCloudtrailName, cloudtrailName},
				MsgRsp{cmd.QuestionCloudtrailS3ExistingBucketArn, s3BucketArn},
				MsgMenu{cmd.OptCloudtrailDone, 4},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, false, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithCloudtrailUseExistingS3(true),
		aws.WithCloudtrailName(cloudtrailName),
		aws.WithExistingCloudtrailBucketArn(s3BucketArn),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test CloutTrail integration with new S3 bucket
func TestGenerationAwsCloudtrailWithNewS3Bucket(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	kmsArn := "arn:aws:kms:us-west-2:249446771485:key/203a7566-41eb-42dc-8cc3-51800a87defe"
	s3BucketName := "s3-bucket-name"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailUseConsolidated, "n"},
				MsgRsp{cmd.QuestionCloudtrailAdvanced, "y"},
				MsgMenu{cmd.OptCloudtrailS3, 0},
				MsgRsp{cmd.QuestionCloudtrailUseExistingS3, "n"},
				MsgRsp{cmd.QuestionCloudtrailS3BucketName, s3BucketName},
				MsgRsp{cmd.QuestionCloudtrailS3BucketEnableEncryption, "y"},
				MsgRsp{cmd.QuestionCloudtrailS3BucketSseKeyArn, kmsArn},
				MsgRsp{cmd.QuestionCloudtrailS3BucketNotification, "y"},
				MsgMenu{cmd.OptCloudtrailDone, 4},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, false, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithCloudtrailUseExistingS3(false),
		aws.WithBucketName(s3BucketName),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithBucketSSEKeyArn(kmsArn),
		aws.WithS3BucketNotification(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test CloutTrail integration with existing SNS topic
func TestGenerationAwsCloudtrailWithExistingSnsTopic(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	snsTopicArn := "arn:aws:sns:us-east-2:249446771485:topic-name"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailUseConsolidated, "n"},
				MsgRsp{cmd.QuestionCloudtrailAdvanced, "y"},
				MsgMenu{cmd.OptCloudtrailSNS, 1},
				MsgRsp{cmd.QuestionCloudtrailUseExistingSNSTopic, "y"},
				MsgRsp{cmd.QuestionCloudtrailSnsExistingTopicArn, snsTopicArn},
				MsgMenu{cmd.OptCloudtrailDone, 4},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, false, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithCloudtrailUseExistingSNSTopic(true),
		aws.WithExistingSnsTopicArn(snsTopicArn),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test CloutTrail integration with new SNS topic
func TestGenerationAwsCloudtrailWithNewSNSTopic(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	kmsArn := "arn:aws:kms:us-west-2:249446771485:key/203a7566-41eb-42dc-8cc3-51800a87defe"
	snsTopicName := "sns-topic-name"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailUseConsolidated, "n"},
				MsgRsp{cmd.QuestionCloudtrailAdvanced, "y"},
				MsgMenu{cmd.OptCloudtrailSNS, 1},
				MsgRsp{cmd.QuestionCloudtrailUseExistingSNSTopic, "n"},
				MsgRsp{cmd.QuestionCloudtrailSnsTopicName, snsTopicName},
				MsgRsp{cmd.QuestionCloudtrailSnsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionCloudtrailSnsEncryptionKeyArn, kmsArn},
				MsgMenu{cmd.OptCloudtrailDone, 4},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, false, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithCloudtrailUseExistingSNSTopic(false),
		aws.WithSnsTopicName(snsTopicName),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionKeyArn(kmsArn),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test CloutTrail integration with new SQS queue
func TestGenerationAwsCloudtrailWithNewSQSQueue(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	sqsQueueName := "sqs-queue-name"
	kmsArn := "arn:aws:kms:us-west-2:249446771485:key/203a7566-41eb-42dc-8cc3-51800a87defe"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailUseConsolidated, "n"},
				MsgRsp{cmd.QuestionCloudtrailAdvanced, "y"},
				MsgMenu{cmd.OptCloudtrailSQS, 2},
				MsgRsp{cmd.QuestionCloudtrailSqsQueueName, sqsQueueName},
				MsgRsp{cmd.QuestionCloudtrailSqsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionCloudtrailSqsEncryptionKeyArn, kmsArn},
				MsgMenu{cmd.OptCloudtrailDone, 4},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, false, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithSqsQueueName(sqsQueueName),
		aws.WithSqsEncryptionKeyArn(kmsArn),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test CloutTrail integration with existing IAM role
func TestGenerationAwsCloudtrailWithExistingIamRole(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	roleName := "test-iamrole"
	roleArn := "arn:aws:iam::123456789012:role/application_abc/component_xyz/abc_role"
	roleExtId := "123456"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailUseConsolidated, "n"},
				MsgRsp{cmd.QuestionCloudtrailAdvanced, "y"},
				MsgMenu{cmd.OptCloudtrailIAM, 3},
				MsgRsp{cmd.QuestionCloudtrailExistingIamRoleName, roleName},
				MsgRsp{cmd.QuestionCloudtrailExistingIamRoleArn, roleArn},
				MsgRsp{cmd.QuestionCloudtrailExistingIamRoleExtID, roleExtId},
				MsgMenu{cmd.OptCloudtrailDone, 4},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, false, true,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
		aws.WithExistingIamRole(
			aws.NewExistingIamRoleDetails(roleName, roleArn, roleExtId),
		),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test overwriting existing main.tf
func TestGenerationAwsWithExistingTerraform(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// Tempdir for test
	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Create fake main.tf
	if err := os.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)), []byte{}, 0644); err != nil {
		panic(err)
	}

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "y"},
				MsgRsp{cmd.QuestionConfigAdditionalAccountAddMore, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsOutputLocation, dir},
				MsgRsp{fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), "y"},
			})
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	data, err := os.ReadFile(fmt.Sprintf("%s/main.tf", dir))
	if err != nil {
		panic(err)
	}

	assert.Empty(t, data)
}

// Test custom output location
func TestGenerationAwsOutputLocation(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	// Tempdir for test
	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "y"},
				MsgRsp{cmd.QuestionConfigAdditionalAccountAddMore, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsOutputLocation, dir},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Get result
	result, _ := os.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(false, false, true, false,
		aws.WithAwsProfile("main"),
		aws.WithAwsRegion("us-east-2"),
	).Generate()
	assert.Equal(t, buildTf, string(result))
}

// Test error on no integration selected
func TestGenerationAwsErrorOnNoSelection(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAwsOrganization, "n"},
				MsgRsp{cmd.QuestionMainAwsProfile, "main"},
				MsgRsp{cmd.QuestionMainAwsRegion, "us-east-2"},
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsOutputLocation, ""},
				MsgOnly{"ERROR failed to generate terraform code: invalid inputs: Agentless, CloudTrail or Config integration must be enabled"},
			})
		},
		"generate",
		"cloud-account",
		"aws",
	)
}

func runGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	os.Setenv("HOME", tfPath)

	hcl_path := filepath.Join(tfPath, awsPath, "main.tf")

	runFakeTerminalTestFromDir(t, tfPath, conditions, args...)
	out, err := os.ReadFile(hcl_path)
	if err != nil {
		return fmt.Sprintf("main.tf not found: %s", err)
	}

	t.Cleanup(func() {
		os.Remove(hcl_path)
	})

	result := terraformValidate(filepath.Join(tfPath, awsPath))

	assert.True(t, result.Valid)

	return string(out)
}
