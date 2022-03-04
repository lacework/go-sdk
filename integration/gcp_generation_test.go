//go:build !windows && generation

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwgenerate/gcp"
	"github.com/stretchr/testify/assert"
)

func expectStringGcp(c *expect.Console, str string, runError *error) {
	out, err := c.Expect(expect.WithTimeout(time.Second), expect.String(str))
	if err != nil {
		fmt.Println(out) // To see the errored line, you can enable this and update _ above to out
		*runError = err
	}
}

// Test failing due to no selection
func TestGenerationErrorOnNoSelectionGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var runError error

	// Run CLI
	runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectStringGcp(c, cmd.QuestionGcpEnableConfig, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpEnableAuditLog, &runError)
			c.SendLine("n")
			expectStringGcp(c, "ERROR collecting/confirming parameters: must enable audit log or config", &runError)
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI errored properly
	assert.Nil(t, runError)
}

// Test bare-bones generation with no customization
func TestGenerationSimpleGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectStringGcp(c, cmd.QuestionGcpEnableConfig, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpEnableAuditLog, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpProjectID, &runError)
			c.SendLine(projectId)
			expectStringGcp(c, cmd.QuestionGcpOrganizationIntegration, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpServiceAccountCredsPath, &runError)
			c.SendLine("")
			expectStringGcp(c, cmd.QuestionGcpConfigureAdvanced, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test customized output location
func TestGenerationCustomizedOutputLocationGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	projectId := "project-1"

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Run CLI
	runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectStringGcp(c, cmd.QuestionGcpEnableConfig, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpEnableAuditLog, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpProjectID, &runError)
			c.SendLine(projectId)
			expectStringGcp(c, cmd.QuestionGcpOrganizationIntegration, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpServiceAccountCredsPath, &runError)
			c.SendLine("")
			expectStringGcp(c, cmd.QuestionGcpConfigureAdvanced, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.GcpAdvancedOptAuditLog, &runError)
			// This is key down x3 in ANSI
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectStringGcp(c, cmd.QuestionGcpCustomizeOutputLocation, &runError)
			c.SendLine(dir)
			expectStringGcp(c, cmd.QuestionGcpAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Get result
	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/gcp.tf", dir)))

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, string(result))
}

// Test config only generation
func TestGenerationConfigOnlyGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectStringGcp(c, cmd.QuestionGcpEnableConfig, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpEnableAuditLog, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpProjectID, &runError)
			c.SendLine(projectId)
			expectStringGcp(c, cmd.QuestionGcpOrganizationIntegration, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpServiceAccountCredsPath, &runError)
			c.SendLine("")
			expectStringGcp(c, cmd.QuestionGcpConfigureAdvanced, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, false,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Bailing out of Advanced Options
func TestGenerationAdvancedOptsDoneGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectStringGcp(c, cmd.QuestionGcpEnableConfig, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpEnableAuditLog, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpProjectID, &runError)
			c.SendLine(projectId)
			expectStringGcp(c, cmd.QuestionGcpOrganizationIntegration, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpServiceAccountCredsPath, &runError)
			c.SendLine("")
			expectStringGcp(c, cmd.QuestionGcpConfigureAdvanced, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.GcpAdvancedOptAuditLog, &runError)
			// This is key down x3 in ANSI
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

//// Test enabling consolidated trail and force destroy s3
//func TestGenerationAdvancedOptsConsolidatedAndForceDestroy(t *testing.T) {
//	os.Setenv("LW_NOCACHE", "true")
//	defer os.Setenv("LW_NOCACHE", "")
//	var final string
//	var runError error
//	region := "us-east-2"
//
//	// Run CLI
//	tfResult := runGcpGenerateTest(t,
//		func(c *expect.Console) {
//			expectStringGcp(c, cmd.QuestionAwsEnableConfig, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionEnableCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionAwsRegion, &runError)
//			c.SendLine(region)
//			expectStringGcp(c, cmd.QuestionAwsConfigAdvanced, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.AwsAdvancedOptDone, &runError)
//			c.SendLine("")
//			expectStringGcp(c, cmd.QuestionConsolidatedCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionUseExistingCloudtrail, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionForceDestroyS3Bucket, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionAwsAnotherAdvancedOpt, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
//			c.SendLine("n")
//			final, _ = c.ExpectEOF()
//		},
//		"cloud",
//		"iac",
//		"aws",
//	)
//
//	// Ensure CLI ran correctly
//	assert.Nil(t, runError)
//	assert.Contains(t, final, "Terraform code saved in")
//
//	// Create the TF directly with lwgenerate and validate same result via CLI
//	buildTf, _ := aws.NewTerraform(region, true, true,
//		aws.UseConsolidatedCloudtrail(), aws.EnableForceDestroyS3Bucket(), aws.WithAwsProfile("default")).Generate()
//	assert.Equal(t, buildTf, tfResult)
//}
//
//// Test use existing cloudtrail
//func TestGenerationAdvancedOptsUseExistingCloudtrail(t *testing.T) {
//	os.Setenv("LW_NOCACHE", "true")
//	defer os.Setenv("LW_NOCACHE", "")
//	var final string
//	var runError error
//	region := "us-east-2"
//
//	// Run CLI
//	tfResult := runGcpGenerateTest(t,
//		func(c *expect.Console) {
//			expectStringGcp(c, cmd.QuestionAwsEnableConfig, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionEnableCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionAwsRegion, &runError)
//			c.SendLine(region)
//			expectStringGcp(c, cmd.QuestionAwsConfigAdvanced, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.AwsAdvancedOptDone, &runError)
//			c.SendLine("")
//			expectStringGcp(c, cmd.QuestionConsolidatedCloudtrail, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionUseExistingCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionCloudtrailExistingBucketArn, &runError)
//			c.SendLine("notright") // test our validator is working
//			expectStringGcp(c, "invalid arn supplied", &runError)
//			c.SendLine("arn:aws:s3:::bucket_name")
//			expectStringGcp(c, cmd.QuestionAwsAnotherAdvancedOpt, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
//			c.SendLine("n")
//			final, _ = c.ExpectEOF()
//		},
//		"cloud",
//		"iac",
//		"aws",
//	)
//
//	// Ensure CLI ran correctly
//	assert.Nil(t, runError)
//	assert.Contains(t, final, "Terraform code saved in")
//
//	// Create the TF directly with lwgenerate and validate same result via CLI
//	buildTf, _ := aws.NewTerraform(region, true, true,
//		aws.ExistingCloudtrailBucketArn("arn:aws:s3:::bucket_name"), aws.WithAwsProfile("default")).Generate()
//	assert.Equal(t, buildTf, tfResult)
//}
//
//// Test using consolidated cloudtrail with subaccounts
//func TestGenerationAdvancedOptsConsolidatedWithSubAccounts(t *testing.T) {
//	os.Setenv("LW_NOCACHE", "true")
//	defer os.Setenv("LW_NOCACHE", "")
//	var final string
//	var runError error
//	region := "us-east-2"
//
//	// Run CLI
//	tfResult := runGcpGenerateTest(t,
//		func(c *expect.Console) {
//			expectStringGcp(c, cmd.QuestionAwsEnableConfig, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionEnableCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionAwsRegion, &runError)
//			c.SendLine(region)
//			expectStringGcp(c, cmd.QuestionAwsConfigAdvanced, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.AwsAdvancedOptDone, &runError)
//			c.SendLine("")
//			expectStringGcp(c, cmd.QuestionConsolidatedCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionUseExistingCloudtrail, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionForceDestroyS3Bucket, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionAwsAnotherAdvancedOpt, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.AwsAdvancedOptDone, &runError)
//			c.Send("\x1B[B") // Down arrow twice and enter on the submenu to add subaccounts
//			c.SendLine("\x1B[B")
//			expectStringGcp(c, cmd.QuestionPrimaryAwsAccountProfile, &runError)
//			c.SendLine("default")
//			expectStringGcp(c, cmd.QuestionSubAccountProfileName, &runError)
//			c.SendLine("account1")
//			expectStringGcp(c, cmd.QuestionSubAccountRegion, &runError)
//			c.SendLine("us-east-1")
//			expectStringGcp(c, cmd.QuestionSubAccountAddMore, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionSubAccountProfileName, &runError)
//			c.SendLine("account2")
//			expectStringGcp(c, cmd.QuestionSubAccountRegion, &runError)
//			c.SendLine("us-east-2")
//			expectStringGcp(c, cmd.QuestionSubAccountAddMore, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionAwsAnotherAdvancedOpt, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
//			c.SendLine("n")
//			final, _ = c.ExpectEOF()
//		},
//		"cloud",
//		"iac",
//		"aws",
//	)
//
//	// Ensure CLI ran correctly
//	assert.Nil(t, runError)
//	assert.Contains(t, final, "Terraform code saved in")
//
//	// Create the TF directly with lwgenerate and validate same result via CLI
//	buildTf, _ := aws.NewTerraform(region, true, true,
//		aws.UseConsolidatedCloudtrail(),
//		aws.WithAwsProfile("default"),
//		aws.WithSubaccounts(aws.NewAwsSubAccount("account1", "us-east-1"), aws.NewAwsSubAccount("account2", "us-east-2")),
//	).Generate()
//	assert.Equal(t, buildTf, tfResult)
//}
//
//// for testing the interaction of the CLI prompts when accounts have been supplied
//func TestGenerationAdvancedOptsConsolidatedWithSubAccountsPassedByFlag(t *testing.T) {
//	os.Setenv("LW_NOCACHE", "true")
//	defer os.Setenv("LW_NOCACHE", "")
//	var final string
//	var runError error
//	region := "us-east-2"
//
//	// Run CLI
//	tfResult := runGcpGenerateTest(t,
//		func(c *expect.Console) {
//			expectStringGcp(c, cmd.QuestionAwsEnableConfig, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionEnableCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionAwsRegion, &runError)
//			c.SendLine(region)
//			expectStringGcp(c, cmd.QuestionAwsConfigAdvanced, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.AwsAdvancedOptDone, &runError)
//			c.Send("\x1B[B") // Down arrow twice and enter on the submenu to add subaccounts
//			c.SendLine("\x1B[B")
//			expectStringGcp(c, cmd.QuestionPrimaryAwsAccountProfile, &runError)
//			c.SendLine("default")
//			expectStringGcp(c, fmt.Sprintf(cmd.QuestionSubAccountReplace, "testaccount:us-east-1, testaccount1:us-east-2"), &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionSubAccountProfileName, &runError)
//			c.SendLine("account1")
//			expectStringGcp(c, cmd.QuestionSubAccountRegion, &runError)
//			c.SendLine("us-east-1")
//			expectStringGcp(c, cmd.QuestionSubAccountAddMore, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionSubAccountProfileName, &runError)
//			c.SendLine("account2")
//			expectStringGcp(c, cmd.QuestionSubAccountRegion, &runError)
//			c.SendLine("us-east-2")
//			expectStringGcp(c, cmd.QuestionSubAccountAddMore, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionAwsAnotherAdvancedOpt, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
//			c.SendLine("n")
//			final, _ = c.ExpectEOF()
//		},
//		"cloud",
//		"iac",
//		"aws",
//		"--consolidated_cloudtrail",
//		"--aws_subaccount",
//		"testaccount:us-east-1",
//		"--aws_subaccount",
//		"testaccount1:us-east-2",
//	)
//
//	// Ensure CLI ran correctly
//	assert.Nil(t, runError)
//	assert.Contains(t, final, "Terraform code saved in")
//
//	// Create the TF directly with lwgenerate and validate same result via CLI
//	buildTf, _ := aws.NewTerraform(region, true, true,
//		aws.UseConsolidatedCloudtrail(),
//		aws.WithAwsProfile("default"),
//		aws.WithSubaccounts(aws.NewAwsSubAccount("account1", "us-east-1"), aws.NewAwsSubAccount("account2", "us-east-2")),
//	).Generate()
//	assert.Equal(t, buildTf, tfResult)
//}
//
//// Test use existing IAM rolej
//func TestGenerationAdvancedOptsUseExistingIAM(t *testing.T) {
//	os.Setenv("LW_NOCACHE", "true")
//	defer os.Setenv("LW_NOCACHE", "")
//	var final string
//	var runError error
//	region := "us-east-2"
//	roleName := "test-iamrole"
//	roleArn := "arn:aws:iam::123456789012:role/application_abc/component_xyz/abc_role"
//	roleExtId := "foo"
//
//	// Run CLI
//	tfResult := runGcpGenerateTest(t,
//		func(c *expect.Console) {
//			expectStringGcp(c, cmd.QuestionAwsEnableConfig, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionEnableCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionAwsRegion, &runError)
//			c.SendLine(region)
//			expectStringGcp(c, cmd.QuestionAwsConfigAdvanced, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.AwsAdvancedOptDone, &runError)
//			c.SendLine("\x1B[B") // Down arrow once and return
//			expectStringGcp(c, cmd.QuestionExistingIamRoleName, &runError)
//			c.SendLine(roleName)
//			expectStringGcp(c, cmd.QuestionExistingIamRoleArn, &runError)
//			c.SendLine(roleArn)
//			expectStringGcp(c, cmd.QuestionExistingIamRoleExtID, &runError)
//			c.SendLine(roleExtId)
//			expectStringGcp(c, cmd.QuestionAwsAnotherAdvancedOpt, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, cmd.QuestionRunTfPlan, &runError)
//			c.SendLine("n")
//			final, _ = c.ExpectEOF()
//		},
//		"cloud-account",
//		"iac-generate",
//		"aws",
//	)
//
//	// Ensure CLI ran correctly
//	assert.Nil(t, runError)
//	assert.Contains(t, final, "Terraform code saved in")
//
//	// Create the TF directly with lwgenerate and validate same result via CLI
//	buildTf, _ := aws.NewTerraform(region, true, true,
//		aws.UseExistingIamRole(aws.NewExistingIamRoleDetails(roleName, roleArn, roleExtId)), aws.WithAwsProfile("default"),
//	).Generate()
//	assert.Equal(t, buildTf, tfResult)
//}
//
//// Test existing aws.tf prompt
//func TestGenerationWithExistingTerraform(t *testing.T) {
//	os.Setenv("LW_NOCACHE", "true")
//	defer os.Setenv("LW_NOCACHE", "")
//	var runError error
//	region := "us-east-2"
//
//	// Tempdir for test
//	dir, err := ioutil.TempDir("", "lacework-cli")
//	if err != nil {
//		panic(err)
//	}
//	defer os.RemoveAll(dir)
//
//	// Create fake aws.tf
//	if err := os.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/aws.tf", dir)), []byte{}, 0644); err != nil {
//		panic(err)
//	}
//
//	// Run CLI
//	runGcpGenerateTest(t,
//		func(c *expect.Console) {
//			expectStringGcp(c, cmd.QuestionAwsEnableConfig, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionEnableCloudtrail, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.QuestionAwsRegion, &runError)
//			c.SendLine(region)
//			expectStringGcp(c, cmd.QuestionAwsConfigAdvanced, &runError)
//			c.SendLine("y")
//			expectStringGcp(c, cmd.AwsAdvancedOptDone, &runError)
//			c.Send("\x1B[B")
//			c.SendLine("\x1B[B")
//			expectStringGcp(c, cmd.QuestionAwsCustomizeOutputLocation, &runError)
//			c.SendLine(dir)
//			expectStringGcp(c, cmd.QuestionAwsAnotherAdvancedOpt, &runError)
//			c.SendLine("n")
//			expectStringGcp(c, fmt.Sprintf("%s/aws.tf already exists, overwrite?", dir), &runError)
//			c.SendLine("n")
//		},
//		"cloud",
//		"iac",
//		"aws",
//	)
//
//	// Ensure CLI ran correctly
//	data, err := os.ReadFile(fmt.Sprintf("%s/aws.tf", dir))
//	if err != nil {
//		panic(err)
//	}
//
//	assert.Empty(t, data)
//	assert.Nil(t, runError)
//}
//
func runGcpGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	// create a temporal directory where we will check that the
	// configuration file is deployed (.lacework.toml)
	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	runGenerationTestFromDir(t, dir, conditions, args...)
	out, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/gcp.tf", dir)))
	if err != nil {
		// Assume couldn't be found
		return ""
	}

	return string(out)
}
