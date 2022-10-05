//go:build !windows && generation

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/stretchr/testify/assert"
)

// Test failing due to no selection
func TestGenerationAwsErrorOnNoSelection(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("n")
			expectString(t, c, "ERROR collecting/confirming parameters: must enable cloudtrail or config")
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
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
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionAwsCustomizeOutputLocation)
			c.SendLine(dir)
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Get result
	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, false,
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
		aws.WithSqsEncryptionEnabled(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test enabling consolidated trail and force destroy s3
func TestGenerationAwsAdvancedOptsConsolidatedAndForceDestroy(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionConsolidatedCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionUseExistingCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionCloudtrailName)
			c.SendLine("")
			// S3 Bucket Questions
			expectString(t, c, cmd.QuestionForceDestroyS3Bucket)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionBucketName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionBucketEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionBucketSseKeyArn)
			c.SendLine("")
			// SNS Topic Questions
			expectString(t, c, cmd.QuestionsUseExistingSNSTopic)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionSnsTopicName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionSnsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSnsEncryptionKeyArn)
			c.SendLine("")
			// SQS Questions
			expectString(t, c, cmd.QuestionSqsQueueName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionSqsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSqsEncryptionKeyArn)
			c.SendLine("")
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.UseConsolidatedCloudtrail(),
		aws.EnableForceDestroyS3Bucket(),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionConsolidatedCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionUseExistingCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionCloudtrailExistingBucketArn)
			c.SendLine("notright") // test our validator is working
			expectString(t, c, "invalid arn supplied")
			c.SendLine("arn:aws:s3:::bucket_name")
			// SNS Topic Questions
			expectString(t, c, cmd.QuestionsUseExistingSNSTopic)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionSnsTopicName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionSnsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSnsEncryptionKeyArn)
			c.SendLine("")
			// SQS Questions
			expectString(t, c, cmd.QuestionSqsQueueName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionSqsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSqsEncryptionKeyArn)
			c.SendLine("")
			//
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.ExistingCloudtrailBucketArn("arn:aws:s3:::bucket_name"),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionConsolidatedCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionUseExistingCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionCloudtrailName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionForceDestroyS3Bucket)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionBucketName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionBucketEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionBucketSseKeyArn)
			c.SendLine("")
			// SNS Topic Questions
			expectString(t, c, cmd.QuestionsUseExistingSNSTopic)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionSnsTopicName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionSnsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSnsEncryptionKeyArn)
			c.SendLine("")
			// SQS Questions
			expectString(t, c, cmd.QuestionSqsQueueName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionSqsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSqsEncryptionKeyArn)
			c.SendLine("")
			//
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.Send("\x1B[B")
			c.SendLine("\x1B[B") // Down arrow twice and enter on the submenu to add subaccounts
			expectString(t, c, cmd.QuestionPrimaryAwsAccountProfile)
			c.SendLine("default")
			expectString(t, c, cmd.QuestionSubAccountProfileName)
			c.SendLine("account1")
			expectString(t, c, cmd.QuestionSubAccountRegion)
			c.SendLine("us-east-1")
			expectString(t, c, cmd.QuestionSubAccountAddMore)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSubAccountProfileName)
			c.SendLine("account2")
			expectString(t, c, cmd.QuestionSubAccountRegion)
			c.SendLine("us-east-2")
			expectString(t, c, cmd.QuestionSubAccountAddMore)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.UseConsolidatedCloudtrail(),
		aws.WithAwsProfile("default"),
		aws.WithSubaccounts(aws.NewAwsSubAccount("account1", "us-east-1"), aws.NewAwsSubAccount("account2", "us-east-2")),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionPrimaryAwsAccountProfile)
			c.SendLine("default")
			expectString(t, c, cmd.QuestionSubAccountProfileName)
			c.SendLine("account1")
			expectString(t, c, cmd.QuestionSubAccountRegion)
			c.SendLine("us-east-1")
			expectString(t, c, cmd.QuestionSubAccountAddMore)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSubAccountProfileName)
			c.SendLine("account2")
			expectString(t, c, cmd.QuestionSubAccountRegion)
			c.SendLine("us-east-2")
			expectString(t, c, cmd.QuestionSubAccountAddMore)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, false,
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.Send("\x1B[B")
			c.SendLine("\x1B[B") // Down arrow twice and enter on the submenu to add subaccounts
			expectString(t, c, cmd.QuestionPrimaryAwsAccountProfile)
			c.SendLine("default")
			expectString(t, c, fmt.Sprintf(cmd.QuestionSubAccountReplace, "testaccount:us-east-1, testaccount1:us-east-2"))
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSubAccountProfileName)
			c.SendLine("account1")
			expectString(t, c, cmd.QuestionSubAccountRegion)
			c.SendLine("us-east-1")
			expectString(t, c, cmd.QuestionSubAccountAddMore)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSubAccountProfileName)
			c.SendLine("account2")
			expectString(t, c, cmd.QuestionSubAccountRegion)
			c.SendLine("us-east-2")
			expectString(t, c, cmd.QuestionSubAccountAddMore)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
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
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.UseConsolidatedCloudtrail(),
		aws.WithAwsProfile("default"),
		aws.WithSubaccounts(aws.NewAwsSubAccount("account1", "us-east-1"), aws.NewAwsSubAccount("account2", "us-east-2")),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B") // Down arrow once and return
			expectString(t, c, cmd.QuestionExistingIamRoleName)
			c.SendLine(roleName)
			expectString(t, c, cmd.QuestionExistingIamRoleArn)
			c.SendLine(roleArn)
			expectString(t, c, cmd.QuestionExistingIamRoleExtID)
			c.SendLine(roleExtId)
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.UseExistingIamRole(aws.NewExistingIamRoleDetails(roleName, roleArn, roleExtId)),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithSnsEncryptionEnabled(true),
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.SendLine("\x1B[B") // Down arrow once and return
			expectString(t, c, cmd.QuestionConsolidatedCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionUseExistingCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionCloudtrailExistingBucketArn)
			c.SendLine(bucketArn)
			expectString(t, c, cmd.QuestionsUseExistingSNSTopic)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSnsTopicArn)
			c.SendLine(topicArn)
			expectString(t, c, cmd.QuestionSqsQueueName)
			c.SendLine(queueName)
			expectString(t, c, cmd.QuestionSqsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSqsEncryptionKeyArn)
			c.SendLine("")
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.SendLine("\x1B[B") // Down arrow once and return
			expectString(t, c, cmd.QuestionConsolidatedCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionUseExistingCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionCloudtrailName)
			c.SendLine(trailName)
			// S3 Questions
			expectString(t, c, cmd.QuestionForceDestroyS3Bucket)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionBucketName)
			c.SendLine(bucketName)
			expectString(t, c, cmd.QuestionBucketEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionBucketSseKeyArn)
			c.SendLine(kmsArn)
			// SNS Topic Questions
			expectString(t, c, cmd.QuestionsUseExistingSNSTopic)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionSnsTopicName)
			c.SendLine(topicName)
			expectString(t, c, cmd.QuestionSnsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSnsEncryptionKeyArn)
			c.SendLine(kmsArn)
			// SQS Questions
			expectString(t, c, cmd.QuestionSqsQueueName)
			c.SendLine(queueName)
			expectString(t, c, cmd.QuestionSqsEnableEncryption)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionSqsEncryptionKeyArn)
			c.SendLine(kmsArn)
			//
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.WithCloudtrailName(trailName),
		aws.EnableForceDestroyS3Bucket(),
		aws.WithBucketName(bucketName),
		aws.WithBucketEncryptionEnabled(true),
		aws.WithBucketSSEKeyArn(kmsArn),
		aws.WithSnsTopicName(topicName),
		aws.WithSnsEncryptionEnabled(true),
		aws.WithSnsEncryptionKeyArn(kmsArn),
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
	dir, err := ioutil.TempDir("", "lacework-cli")
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
			expectString(t, c, cmd.QuestionAwsEnableConfig)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionEnableCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsRegion)
			c.SendLine(region)
			expectString(t, c, cmd.QuestionAwsConfigAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionAwsCustomizeOutputLocation)
			c.SendLine(dir)
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, fmt.Sprintf("%s/main.tf already exists, overwrite?", dir))
			c.SendLine("n")
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

func runGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	// create a temporal directory where we will check that the
	// configuration file is deployed (.lacework.toml)
	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	runFakeTerminalTestFromDir(t, dir, conditions, args...)
	out, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/aws/main.tf", dir)))
	if err != nil {
		// Assume couldn't be found
		return ""
	}

	return string(out)
}
