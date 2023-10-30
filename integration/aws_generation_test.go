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

// Test failing due to no selection
func TestGenerationAwsErrorOnNoSelection(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgOnly{"ERROR collecting/confirming parameters: must enable agentless, cloudtrail or config"},
			})
		},
		"generate",
		"cloud-account",
		"aws",
	)
}

// Test barebones generation with no customization
func TestGenerationAwsSimple(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test customized output location
func TestGenerationAwsCustomizedOutputLocation(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

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
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 5},
				MsgRsp{cmd.QuestionAwsCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, string(result))
}

// Test config only generation
func TestGenerationAwsConfigOnly(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "n"},
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
	buildTf, _ := aws.NewTerraform(region, false, true, false,
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Bailing out of Advanced Options
func TestGenerationAwsAdvancedOptsDone(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 6},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test enabling consolidated trail
func TestGenerationAwsAdvancedOptsConsolidated(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 2},
				MsgRsp{cmd.QuestionConsolidatedCloudtrail, "y"},
				MsgRsp{cmd.QuestionUseExistingCloudtrail, "n"},
				MsgRsp{cmd.QuestionCloudtrailName, ""},
				// S3 Bucket Questions
				MsgRsp{cmd.QuestionBucketName, ""},
				MsgRsp{cmd.QuestionBucketEnableEncryption, "y"},
				MsgRsp{cmd.QuestionBucketSseKeyArn, ""},
				MsgRsp{cmd.QuestionS3BucketNotification, ""},
				// SNS Topic Questions
				MsgRsp{cmd.QuestionsUseExistingSNSTopic, "n"},
				MsgRsp{cmd.QuestionSnsTopicName, ""},
				MsgRsp{cmd.QuestionSnsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSnsEncryptionKeyArn, ""},
				// SQS Questions
				MsgRsp{cmd.QuestionSqsQueueName, ""},
				MsgRsp{cmd.QuestionSqsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSqsEncryptionKeyArn, ""},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.UseConsolidatedCloudtrail(),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test use existing cloudtrail
func TestGenerationAwsAdvancedOptsUseExistingCloudtrail(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 2},
				MsgRsp{cmd.QuestionConsolidatedCloudtrail, "n"},
				MsgRsp{cmd.QuestionUseExistingCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailExistingBucketArn, "notright"},
				MsgRsp{"invalid arn supplied", "arn:aws:s3:::bucket_name"},
				// SNS Topic Questions
				MsgRsp{cmd.QuestionsUseExistingSNSTopic, "n"},
				MsgRsp{cmd.QuestionSnsTopicName, ""},
				MsgRsp{cmd.QuestionSnsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSnsEncryptionKeyArn, ""},
				// SQS Questions
				MsgRsp{cmd.QuestionSqsQueueName, ""},
				MsgRsp{cmd.QuestionSqsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSqsEncryptionKeyArn, ""},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.ExistingCloudtrailBucketArn("arn:aws:s3:::bucket_name"),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test using consolidated cloudtrail with subaccounts
func TestGenerationAwsAdvancedOptsConsolidatedWithSubAccounts(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 2},
				MsgRsp{cmd.QuestionConsolidatedCloudtrail, "y"},
				MsgRsp{cmd.QuestionUseExistingCloudtrail, "n"},
				MsgRsp{cmd.QuestionCloudtrailName, ""},
				MsgRsp{cmd.QuestionBucketName, ""},
				MsgRsp{cmd.QuestionBucketEnableEncryption, "y"},
				MsgRsp{cmd.QuestionBucketSseKeyArn, ""},
				MsgRsp{cmd.QuestionS3BucketNotification, ""},
				// SNS Topic Questions
				MsgRsp{cmd.QuestionsUseExistingSNSTopic, "n"},
				MsgRsp{cmd.QuestionSnsTopicName, ""},
				MsgRsp{cmd.QuestionSnsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSnsEncryptionKeyArn, ""},
				// SQS Questions
				MsgRsp{cmd.QuestionSqsQueueName, ""},
				MsgRsp{cmd.QuestionSqsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSqsEncryptionKeyArn, ""},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 3},
				MsgRsp{cmd.QuestionPrimaryAwsAccountProfile, "default"},
				MsgRsp{cmd.QuestionSubAccountProfileName, "account1"},
				MsgRsp{cmd.QuestionSubAccountRegion, "us-east-1"},
				MsgRsp{cmd.QuestionSubAccountAddMore, "y"},
				MsgRsp{cmd.QuestionSubAccountProfileName, "account2"},
				MsgRsp{cmd.QuestionSubAccountRegion, "us-east-2"},
				MsgRsp{cmd.QuestionSubAccountAddMore, "n"},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.UseConsolidatedCloudtrail(),
		aws.WithAwsProfile("default"),
		aws.WithSubaccounts(aws.NewAwsSubAccount("account1", "us-east-1"), aws.NewAwsSubAccount("account2", "us-east-2")),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test using Config with subaccounts
func TestGenerationAwsAdvancedOptsConfigWithSubAccounts(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 1},
				MsgRsp{cmd.QuestionPrimaryAwsAccountProfile, "default"},
				MsgRsp{cmd.QuestionSubAccountProfileName, "account1"},
				MsgRsp{cmd.QuestionSubAccountRegion, "us-east-1"},
				MsgRsp{cmd.QuestionSubAccountAddMore, "y"},
				MsgRsp{cmd.QuestionSubAccountProfileName, "account2"},
				MsgRsp{cmd.QuestionSubAccountRegion, "us-east-2"},
				MsgRsp{cmd.QuestionSubAccountAddMore, "n"},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, false, true, false,
		aws.WithAwsProfile("default"),
		aws.WithSubaccounts(aws.NewAwsSubAccount("account1", "us-east-1"), aws.NewAwsSubAccount("account2", "us-east-2")),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// for testing the interaction of the CLI prompts when accounts have been supplied
func TestGenerationAwsAdvancedOptsConsolidatedWithSubAccountsPassedByFlag(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 3},
				MsgRsp{cmd.QuestionPrimaryAwsAccountProfile, "default"},
				MsgRsp{fmt.Sprintf(cmd.QuestionSubAccountReplace, "testaccount:us-east-1, testaccount1:us-east-2"), "y"},
				MsgRsp{cmd.QuestionSubAccountProfileName, "account1"},
				MsgRsp{cmd.QuestionSubAccountRegion, "us-east-1"},
				MsgRsp{cmd.QuestionSubAccountAddMore, "y"},
				MsgRsp{cmd.QuestionSubAccountProfileName, "account2"},
				MsgRsp{cmd.QuestionSubAccountRegion, "us-east-2"},
				MsgRsp{cmd.QuestionSubAccountAddMore, "n"},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
		"--consolidated_cloudtrail",
		"--aws_subaccount",
		"testaccount:us-east-1",
		"--aws_subaccount",
		"testaccount1:us-east-2",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.UseConsolidatedCloudtrail(),
		aws.WithAwsProfile("default"),
		aws.WithSubaccounts(aws.NewAwsSubAccount("account1", "us-east-1"), aws.NewAwsSubAccount("account2", "us-east-2")),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test use existing IAM role
func TestGenerationAwsAdvancedOptsUseExistingIAM(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"
	roleName := "test-iamrole"
	roleArn := "arn:aws:iam::123456789012:role/application_abc/component_xyz/abc_role"
	roleExtId := "foo"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 4},
				MsgRsp{cmd.QuestionExistingIamRoleName, roleName},
				MsgRsp{cmd.QuestionExistingIamRoleArn, roleArn},
				MsgRsp{cmd.QuestionExistingIamRoleExtID, roleExtId},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.UseExistingIamRole(aws.NewExistingIamRoleDetails(roleName, roleArn, roleExtId)),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test use of existing Cloudtrail S3 Bucket, SNS topic and SQS Queue
func TestGenerationAwsAdvancedOptsUseExistingElements(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"
	bucketArn := "arn:aws:s3:::bucket-name"
	topicArn := "arn:aws:sns:us-east-2:249446771485:topic-name"
	queueName := "sqs-queue-name"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 2},
				MsgRsp{cmd.QuestionConsolidatedCloudtrail, "n"},
				MsgRsp{cmd.QuestionUseExistingCloudtrail, "y"},
				MsgRsp{cmd.QuestionCloudtrailExistingBucketArn, bucketArn},
				MsgRsp{cmd.QuestionsUseExistingSNSTopic, "y"},
				MsgRsp{cmd.QuestionSnsTopicArn, topicArn},
				MsgRsp{cmd.QuestionSqsQueueName, queueName},
				MsgRsp{cmd.QuestionSqsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSqsEncryptionKeyArn, ""},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.ExistingCloudtrailBucketArn(bucketArn),
		aws.ExistingSnsTopicArn(topicArn),
		aws.WithSqsEncryptionEnabled(true),
		aws.WithSqsQueueName(queueName),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test use of creating S3 Bucket, SNS topic and SQS Queue
func TestGenerationAwsAdvancedOptsCreateNewElements(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"
	kmsArn := "arn:aws:kms:us-west-2:249446771485:key/203a7566-41eb-42dc-8cc3-51800a87defe"
	trailName := "cloudtrail-integration-name"
	bucketName := "s3-bucket-name"
	topicName := "sns-topic-name"
	queueName := "sqs-queue-name"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 2},
				MsgRsp{cmd.QuestionConsolidatedCloudtrail, "n"},
				MsgRsp{cmd.QuestionUseExistingCloudtrail, "n"},
				MsgRsp{cmd.QuestionCloudtrailName, trailName},
				// S3 Questions
				MsgRsp{cmd.QuestionBucketName, bucketName},
				MsgRsp{cmd.QuestionBucketEnableEncryption, "y"},
				MsgRsp{cmd.QuestionBucketSseKeyArn, kmsArn},
				MsgRsp{cmd.QuestionS3BucketNotification, ""},
				// SNS Topic Questions
				MsgRsp{cmd.QuestionsUseExistingSNSTopic, "n"},
				MsgRsp{cmd.QuestionSnsTopicName, topicName},
				MsgRsp{cmd.QuestionSnsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSnsEncryptionKeyArn, kmsArn},
				// SQS Questions
				MsgRsp{cmd.QuestionSqsQueueName, queueName},
				MsgRsp{cmd.QuestionSqsEnableEncryption, "y"},
				MsgRsp{cmd.QuestionSqsEncryptionKeyArn, kmsArn},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.WithCloudtrailName(trailName),
		aws.WithBucketName(bucketName),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithBucketSSEKeyArn(kmsArn),
		aws.WithSnsTopicName(topicName),
		aws.WithSnsTopicEncryptionEnabled(true),
		aws.WithSnsTopicEncryptionKeyArn(kmsArn),
		aws.WithSqsQueueName(queueName),
		aws.WithSqsEncryptionEnabled(true),
		aws.WithSqsEncryptionKeyArn(kmsArn),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test existing main.tf prompt
func TestGenerationAwsWithExistingTerraform(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	region := "us-east-2"

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
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 5},
				MsgRsp{cmd.QuestionAwsCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
				MsgRsp{fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), "n"},
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

func TestGenerationAwsOverwrite(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s/lacework/aws", dir))

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "n"},
				MsgRsp{"already exists, overwrite?", "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s/lacework/aws", dir))
}

func TestGenerationAwsOverwriteOutput(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)

	output_dir := createDummyTOMLConfig()
	defer os.RemoveAll(output_dir)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
		"--output",
		output_dir,
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "n"},
				MsgRsp{"already exists, overwrite?", "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
		"--output",
		output_dir,
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))
}

func TestGenerationAwsLaceworkProfile(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	region := "us-west-2"
	awsProfile := "v2"

	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "y"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
		"--profile",
		awsProfile,
	)

	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	buildTf, _ := aws.NewTerraform(region, true, true, true,
		aws.WithLaceworkProfile(awsProfile),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAwsS3BucketNotification(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	region := "us-west-2"

	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
		"--use_s3_bucket_notification",
	)

	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	buildTf, _ := aws.NewTerraform(region, false, false, true,
		aws.WithS3BucketNotification(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAwsS3BucketNotificationInteractive(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	region := "us-west-2"

	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "n"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "y"},
				MsgRsp{cmd.QuestionAwsRegion, region},

				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 0},

				MsgRsp{cmd.QuestionConsolidatedCloudtrail, ""},
				MsgRsp{cmd.QuestionUseExistingCloudtrail, ""},
				MsgRsp{cmd.QuestionCloudtrailName, ""},
				// S3 Questions
				MsgRsp{cmd.QuestionBucketName, ""},
				MsgRsp{cmd.QuestionBucketEnableEncryption, ""},
				MsgRsp{cmd.QuestionBucketSseKeyArn, ""},
				MsgRsp{cmd.QuestionS3BucketNotification, "y"},
				// SNS Topic Questions
				MsgRsp{cmd.QuestionsUseExistingSNSTopic, ""},
				MsgRsp{cmd.QuestionSnsTopicName, ""},
				MsgRsp{cmd.QuestionSnsEnableEncryption, ""},
				MsgRsp{cmd.QuestionSnsEncryptionKeyArn, ""},
				// SQS Questions
				MsgRsp{cmd.QuestionSqsQueueName, ""},
				MsgRsp{cmd.QuestionSqsEnableEncryption, ""},
				MsgRsp{cmd.QuestionSqsEncryptionKeyArn, ""},

				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	buildTf, _ := aws.NewTerraform(region, false, false, true,
		aws.WithS3BucketNotification(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Agentless organization integration
func TestGenerationAgentlessOrganization(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionEnableAgentless, "y"},
				MsgRsp{cmd.QuestionAwsEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableCloudtrail, "n"},
				MsgRsp{cmd.QuestionAwsRegion, region},
				MsgRsp{cmd.QuestionAwsConfigAdvanced, "y"},
				MsgMenu{cmd.AwsAdvancedOptDone, 0},
				MsgRsp{cmd.QuestionEnableAgentlessOrganization, "y"},
				MsgRsp{cmd.QuestionPrimaryAwsAccountProfile, "default-profile"},
				MsgRsp{cmd.QuestionAgentlessManagementAccountID, "123456789000"},
				MsgRsp{cmd.QuestionAgentlessMonitoredAccountIDs, "123456789000,ou-abcd-12345678,r-abcd"},
				MsgRsp{cmd.QuestionAwsAnotherAdvancedOpt, "n"},
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
	buildTf, _ := aws.NewTerraform(region, true, false, false,
		aws.UseConsolidatedCloudtrail(),
		aws.WithAwsProfile("default-profile"),
		aws.WithAgentlessManagementAccountID("123456789000"),
		aws.WithAgentlessMonitoredAccountIDs([]string{"123456789000", "ou-abcd-12345678", "r-abcd"}),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
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
