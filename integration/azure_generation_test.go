//go:build !windows && generation

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwgenerate/azure"
	"github.com/stretchr/testify/assert"
)

func expectAzureString(c *expect.Console, str string, runError *error) {
	out, err := c.Expect(expect.WithTimeout(time.Second), expect.String(str))
	if err != nil {
		fmt.Println(out) // To see the errored line, you can enable this and update _ above to out
		*runError = err
	}
}

// Test failing due to no selection
func TestGenerationAzureErrorOnNoSelection(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var runError error

	// Run CLI
	runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("n")
			expectAzureString(c, "ERROR collecting/confirming parameters: must enable activity log or config", &runError)
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI errored properly
	assert.Nil(t, runError)
}

// Test barebones generation with no customization
func TestGenerationAzureSimple(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, true, true).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test customized output location
func TestGenerationAzureCustomizedOutputLocation(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Run CLI
	runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectAzureString(c, cmd.QuestionAzureCustomizeOutputLocation, &runError)
			c.SendLine(dir)
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Get result
	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, true, true).Generate()
	assert.Equal(t, buildTf, string(result))
}

// Test config only generation
func TestGenerationAzureConfigOnly(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, false, true).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test activity log only generation
func TestGenerationAzureActivityLogOnly(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(false, true, true).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test no AD creation generation
func TestGenerationAzureNoADEnabled(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	var pass string = "super-secret-password"
	var principalId string = "test-prinicpal-id"
	var applicationId string = "test-application-id"

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptLocation, &runError)
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")

			expectAzureString(c, cmd.QuestionADApplicationPass, &runError)
			c.SendLine(pass)
			expectAzureString(c, cmd.QuestionADApplicationId, &runError)
			c.SendLine(applicationId)
			expectAzureString(c, cmd.QuestionADServicePrincpleId, &runError)
			c.SendLine(principalId)

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, true, false,
		azure.WithAdApplicationPassword(pass),
		azure.WithAdServicePrincipalId(principalId),
		azure.WithAdApplicationId(applicationId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generation with config name customization
func TestGenerationAzureNamedConfig(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	var configName string = "Test-Config-Rename"

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.SendLine("")

			expectAzureString(c, cmd.QuestionAzureConfigName, &runError)
			c.SendLine(configName)

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, false, true,
		azure.WithConfigIntegrationName(configName),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generation with activity log customization
func TestGenerationAzureNamedActivityLog(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	var activityName string = "Test Activity Log Rename"

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.SendLine("")

			expectAzureString(c, cmd.QuestionActivityLogName, &runError)
			c.SendLine(activityName)

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(false, true, true,
		azure.WithActivityLogIntegrationName(activityName)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test Bailing out of Advanced Options
func TestGenerationAzureAdvancedOptsDone(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, true, true).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test existing main.tf prompt
func TestGenerationAzureWithExistingTerraform(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	//var final string
	var runError error

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	//defer os.RemoveAll(dir)

	// Create fake main.tf
	if err := os.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)), []byte{}, 0644); err != nil {
		panic(err)
	}

	// Run CLI
	runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectAzureString(c, cmd.QuestionAzureCustomizeOutputLocation, &runError)
			c.SendLine(dir)
			expectAzureString(c, fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), &runError)
			c.SendLine("n")

			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			_, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	data, err := os.ReadFile(fmt.Sprintf("%s/main.tf", dir))
	if err != nil {
		panic(err)
	}

	assert.Empty(t, data)
	assert.Nil(t, runError)
}

// Test generate Config with all subscriptions
func TestGenerationAzureConfigAllSubs(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.SendLine("\x1B[B")

			expectAzureString(c, cmd.QuestionEnableAllSubscriptions, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, false, true,
		azure.WithAllSubscriptions(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generate Config with management group
func TestGenerationAzureConfigMgmntGroup(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	var mgmtGrpId string = "test-management-group-1"

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")

			expectAzureString(c, cmd.QuestionEnableManagementGroup, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionManagementGroupId, &runError)
			c.SendLine(mgmtGrpId)

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, false, true,
		azure.WithManagementGroup(true),
		azure.WithManagementGroupId(mgmtGrpId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generate Config log with subscription ids
func TestGenerationAzureConfigSubs(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	testIds := []string{"test-id-1", "test-id-2", "test-id-3"}

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.SendLine("\x1B[B")

			expectAzureString(c, cmd.QuestionEnableAllSubscriptions, &runError)
			c.SendLine("n")

			expectAzureString(c, cmd.QuestionSubscriptionIds, &runError)
			c.SendLine(strings.Join(testIds[:], ","))

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, false, true,
		azure.WithSubscriptionIds(testIds),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generate Activity log with subscription ids
func TestGenerationAzureActivityLogSubs(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	testIds := []string{"test-id-1", "test-id-2", "test-id-3"}

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.SendLine("\x1B[B")

			expectAzureString(c, cmd.QuestionEnableAllSubscriptions, &runError)
			c.SendLine("n")

			expectAzureString(c, cmd.QuestionSubscriptionIds, &runError)
			c.SendLine(strings.Join(testIds[:], ","))

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(false, true, true,
		azure.WithSubscriptionIds(testIds),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generate ActivityLog with exising storage account
func TestGenerationAzureActivityLogStorageAccount(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	var storageAccountName string = "test-storage-account-name"
	var storageResourceGrp string = "test-storage-account-resource-group"

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")

			expectAzureString(c, cmd.QuestionUseExistingStorageAccount, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.QuestionStorageAccountName, &runError)
			c.SendLine(storageAccountName)
			expectAzureString(c, cmd.QuestionStorageAccountResourceGroup, &runError)
			c.SendLine(storageResourceGrp)

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(false, true, true,
		azure.WithExistingStorageAccount(true),
		azure.WithStorageAccountName(storageAccountName),
		azure.WithStorageAccountResourceGroup(storageResourceGrp),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generate ActivityLog with all subscriptions
func TestGenerationAzureActivityLogAllSubs(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.SendLine("\x1B[B")

			expectAzureString(c, cmd.QuestionEnableAllSubscriptions, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(false, true, true,
		azure.WithAllSubscriptions(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generate user supplied location for logging
func TestGenerationAzureActivityLogLocation(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	var region string = "West US 2"

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectAzureString(c, cmd.QuestionAzureEnableConfig, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionEnableActivityLog, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionEnableAdIntegration, &runError)
			c.SendLine("y")
			expectAzureString(c, cmd.QuestionAzureConfigAdvanced, &runError)
			c.SendLine("y")

			expectAzureString(c, cmd.AzureAdvancedOptDone, &runError)
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")

			expectAzureString(c, cmd.QuestionStorageLocation, &runError)
			c.SendLine(region)

			expectAzureString(c, cmd.QuestionAzureAnotherAdvancedOpt, &runError)
			c.SendLine("n")
			expectAzureString(c, cmd.QuestionRunTfPlan, &runError)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	// Ensure CLI ran correctly
	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(false, true, true,
		azure.WithStorageLocation(region),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func runGenerateAzureTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	// create a temporal directory where we will check that the
	// configuration file is deployed (.lacework.toml)
	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	runFakeTerminalTestFromDir(t, dir, conditions, args...)
	out, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/azure/main.tf", dir)))
	if err != nil {
		// Assume couldn't be found
		return ""
	}

	return string(out)
}
