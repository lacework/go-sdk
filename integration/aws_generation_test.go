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
func TestGenerationErrorOnNoSelection(t *testing.T) {
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
		"cloud",
		"iac",
		"aws",
	)
}

// Test barebones generation with no customization
func TestGenerationSimple(t *testing.T) {
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
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true, aws.WithAwsProfile("default")).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test customized output location
func TestGenerationCustomizedOutputLocation(t *testing.T) {
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
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionAwsCustomizeOutputLocation)
			c.SendLine(dir)
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Get result
	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/aws.tf", dir)))

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true, aws.WithAwsProfile("default")).Generate()
	assert.Equal(t, buildTf, string(result))
}

// Test config only generation
func TestGenerationConfigOnly(t *testing.T) {
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
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, false, aws.WithAwsProfile("default")).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Bailing out of Advanced Options
func TestGenerationAdvancedOptsDone(t *testing.T) {
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
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true, aws.WithAwsProfile("default")).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test enabling consolidated trail and force destroy s3
func TestGenerationAdvancedOptsConsolidatedAndForceDestroy(t *testing.T) {
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
			c.SendLine("")
			expectString(t, c, cmd.QuestionConsolidatedCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionUseExistingCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionForceDestroyS3Bucket)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.UseConsolidatedCloudtrail(), aws.EnableForceDestroyS3Bucket(), aws.WithAwsProfile("default")).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test use existing cloudtrail
func TestGenerationAdvancedOptsUseExistingCloudtrail(t *testing.T) {
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
			c.SendLine("")
			expectString(t, c, cmd.QuestionConsolidatedCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionUseExistingCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionCloudtrailExistingBucketArn)
			c.SendLine("notright") // test our validator is working
			expectString(t, c, "invalid arn supplied")
			c.SendLine("arn:aws:s3:::bucket_name")
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.ExistingCloudtrailBucketArn("arn:aws:s3:::bucket_name"), aws.WithAwsProfile("default")).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test using consolidated cloudtrail with subaccounts
func TestGenerationAdvancedOptsConsolidatedWithSubAccounts(t *testing.T) {
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
			c.SendLine("")
			expectString(t, c, cmd.QuestionConsolidatedCloudtrail)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionUseExistingCloudtrail)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionForceDestroyS3Bucket)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("y")
			expectString(t, c, cmd.AwsAdvancedOptDone)
			c.Send("\x1B[B") // Down arrow twice and enter on the submenu to add subaccounts
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
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.UseConsolidatedCloudtrail(),
		aws.WithAwsProfile("default"),
		aws.WithSubaccounts(aws.NewAwsSubAccount("account1", "us-east-1"), aws.NewAwsSubAccount("account2", "us-east-2")),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// for testing the interaction of the CLI prompts when accounts have been supplied
func TestGenerationAdvancedOptsConsolidatedWithSubAccountsPassedByFlag(t *testing.T) {
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
			c.Send("\x1B[B") // Down arrow twice and enter on the submenu to add subaccounts
			c.SendLine("\x1B[B")
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
		"cloud",
		"iac",
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
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test use existing IAM rolej
func TestGenerationAdvancedOptsUseExistingIAM(t *testing.T) {
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
		"cloud-account",
		"iac-generate",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.UseExistingIamRole(aws.NewExistingIamRoleDetails(roleName, roleArn, roleExtId)), aws.WithAwsProfile("default"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test existing aws.tf prompt
func TestGenerationWithExistingTerraform(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	region := "us-east-2"

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Create fake aws.tf
	if err := os.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/aws.tf", dir)), []byte{}, 0644); err != nil {
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
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionAwsCustomizeOutputLocation)
			c.SendLine(dir)
			expectString(t, c, cmd.QuestionAwsAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, fmt.Sprintf("%s/aws.tf already exists, overwrite?", dir))
			c.SendLine("n")
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	data, err := os.ReadFile(fmt.Sprintf("%s/aws.tf", dir))
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
	out, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/aws.tf", dir)))
	if err != nil {
		// Assume couldn't be found
		return ""
	}

	return string(out)
}
