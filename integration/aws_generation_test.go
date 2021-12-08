//go:build generation || (!windows && generation)

package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/stretchr/testify/assert"
)

func expectString(c *expect.Console, str string, runError *error) {
	out, err := c.Expect(expect.WithTimeout(time.Second), expect.String(str))
	if err != nil {
		fmt.Println(out) // To see the errored line, you can enable this and update _ above to out
		*runError = err
	}
}

// Test failing due to no selection
func TestGenerationErrorOnNoSelection(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var runError error

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("n")
			expectString(c, "ERROR error when collecting/confirming parameters: must enable cloudtrail or config", &runError)
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI errored properly
	assert.Nil(t, runError)
}

// Test barebones generation with no customization
func TestGenerationSimple(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
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
	var runError error
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
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(c, cmd.QuestionCustomizeOutputLocation, &runError)
			c.SendLine(dir)
			expectString(c, cmd.QuestionAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Get result
	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/main.tf", dir)))

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true, aws.WithAwsProfile("default")).Generate()
	assert.Equal(t, buildTf, string(result))
}

// Test config only generation
func TestGenerationConfigOnly(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
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
	var runError error
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
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
	var runError error
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.SendLine("")
			expectString(c, cmd.QuestionConsolidatedCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionUseExistingCloudtrail, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionForceDestroyS3Bucket, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
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
	var runError error
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.SendLine("")
			expectString(c, cmd.QuestionConsolidatedCloudtrail, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionUseExistingCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionCloudtrailExistingBucketArn, &runError)
			c.SendLine("notright") // test our validator is working
			expectString(c, "invalid arn supplied", &runError)
			c.SendLine("arn:aws:s3:::bucket_name")
			expectString(c, cmd.QuestionAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
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
	var runError error
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.SendLine("")
			expectString(c, cmd.QuestionConsolidatedCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionUseExistingCloudtrail, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionForceDestroyS3Bucket, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionAnotherAdvancedOpt, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.Send("\x1B[B") // Down arrow twice and enter on the submenu to add subaccounts
			c.SendLine("\x1B[B")
			expectString(c, cmd.QuestionPrimaryAwsAccountProfile, &runError)
			c.SendLine("default")
			expectString(c, cmd.QuestionSubAccountProfileName, &runError)
			c.SendLine("account1")
			expectString(c, cmd.QuestionSubAccountRegion, &runError)
			c.SendLine("us-east-1")
			expectString(c, cmd.QuestionSubAccountAddMore, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionSubAccountProfileName, &runError)
			c.SendLine("account2")
			expectString(c, cmd.QuestionSubAccountRegion, &runError)
			c.SendLine("us-east-2")
			expectString(c, cmd.QuestionSubAccountAddMore, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
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
	var runError error
	region := "us-east-2"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.Send("\x1B[B") // Down arrow twice and enter on the submenu to add subaccounts
			c.SendLine("\x1B[B")
			expectString(c, cmd.QuestionPrimaryAwsAccountProfile, &runError)
			c.SendLine("default")
			expectString(c, fmt.Sprintf(cmd.QuestionSubAccountReplace, "testaccount:us-east-1, testaccount1:us-east-2"), &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionSubAccountProfileName, &runError)
			c.SendLine("account1")
			expectString(c, cmd.QuestionSubAccountRegion, &runError)
			c.SendLine("us-east-1")
			expectString(c, cmd.QuestionSubAccountAddMore, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionSubAccountProfileName, &runError)
			c.SendLine("account2")
			expectString(c, cmd.QuestionSubAccountRegion, &runError)
			c.SendLine("us-east-2")
			expectString(c, cmd.QuestionSubAccountAddMore, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
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
	assert.Nil(t, runError)
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
	var runError error
	region := "us-east-2"
	roleName := "test-iamrole"
	roleArn := "arn:aws:iam::123456789012:role/application_abc/component_xyz/abc_role"
	roleExtId := "foo"

	// Run CLI
	tfResult := runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.SendLine("\x1B[B") // Down arrow once and return
			expectString(c, cmd.QuestionExistingIamRoleName, &runError)
			c.SendLine(roleName)
			expectString(c, cmd.QuestionExistingIamRoleArn, &runError)
			c.SendLine(roleArn)
			expectString(c, cmd.QuestionExistingIamRoleExtId, &runError)
			c.SendLine(roleExtId)
			expectString(c, cmd.QuestionAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud-account",
		"iac-generate",
		"aws",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := aws.NewTerraform(region, true, true,
		aws.UseExistingIamRole(aws.NewExistingIamRoleDetails(roleName, roleArn, roleExtId)), aws.WithAwsProfile("default"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test existing main.tf prompt
func TestGenerationWithExistingTerraform(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var runError error
	region := "us-east-2"

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Create fake main.tf
	if err := os.Mkdir(filepath.FromSlash(fmt.Sprintf("%s/lacework", dir)), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/main.tf", dir)), []byte{}, 0644); err != nil {
		panic(err)
	}

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectString(c, cmd.QuestionEnableConfig, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionEnableCloudtrail, &runError)
			c.SendLine("y")
			expectString(c, cmd.QuestionAwsRegion, &runError)
			c.SendLine(region)
			expectString(c, cmd.QuestionConfigAdvanced, &runError)
			c.SendLine("y")
			expectString(c, cmd.AdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(c, cmd.QuestionCustomizeOutputLocation, &runError)
			c.SendLine(dir)
			expectString(c, cmd.QuestionAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectString(c, fmt.Sprintf("%s/lacework/main.tf already exists, overwrite?", dir), &runError)
			c.SendLine("n")
		},
		"cloud",
		"iac",
		"aws",
	)

	// Ensure CLI ran correctly
	data, err := os.ReadFile(fmt.Sprintf("%s/lacework/main.tf", dir))
	if err != nil {
		panic(err)
	}

	assert.Empty(t, data)
	assert.Nil(t, runError)
}

func runGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	// create a temporal directory where we will check that the
	// configuration file is deployed (.lacework.toml)
	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)
	runGenerationTestFromDir(t, dir, conditions, append(args, "--output", dir)...)

	out, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/main.tf", dir)))
	if err != nil {
		// Assume couldn't be found
		return ""
	}

	return string(out)
}

func runGenerationTestFromDir(t *testing.T, dir string, conditions func(*expect.Console), args ...string) {
	console, state, err := vt10x.NewVT10XConsole()
	if err != nil {
		panic(err)
	}
	defer console.Close()

	if os.Getenv("DEBUG") != "" {
		state.DebugLogger = log.Default()
	}

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		conditions(console)
	}()

	cmd := NewLaceworkCLI(dir, nil, args...)
	cmd.Stdin = console.Tty()
	cmd.Stdout = console.Tty()
	cmd.Stderr = console.Tty()
	err = cmd.Start()
	assert.Nil(t, err)

	// read the remaining bytes
	console.Tty().Close()
	<-donec
}
