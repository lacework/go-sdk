//go:build !windows && generation

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/v2/cli/cmd"
	"github.com/lacework/go-sdk/v2/lwgenerate/azure"
	"github.com/stretchr/testify/assert"
)

const (
	azurePath          = "/lacework/azure/"
	mockSubscriptionID = "111aaa1a-a1a1-11aa-a111-1aaaa1a11a11"
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgOnly{"ERROR collecting/confirming parameters: must enable at least one of: Configuration or Activity Log integration"},
			})
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(true, true, false, true, azure.WithSubscriptionID(mockSubscriptionID)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test customized output location
func TestGenerationAzureCustomizedOutputLocation(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Tempdir for test
	dir, err := os.MkdirTemp("", "t")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Run CLI
	runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	result, _ := os.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := azure.NewTerraform(true, true, false, true, azure.WithSubscriptionID(mockSubscriptionID)).Generate()
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(true, false, false, true, azure.WithSubscriptionID(mockSubscriptionID)).Generate()
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(false, true, false, true, azure.WithSubscriptionID(mockSubscriptionID)).Generate()
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "n"},
				MsgRsp{cmd.QuestionADApplicationPass, pass},
				MsgRsp{cmd.QuestionADApplicationId, applicationId},
				MsgRsp{cmd.QuestionADServicePrincpleId, principalId},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(true, true, false, false,
		azure.WithSubscriptionID(mockSubscriptionID),
		azure.WithAdApplicationPassword(pass),
		azure.WithAdServicePrincipalId(principalId),
		azure.WithAdApplicationId(applicationId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Unstable test disabled as part of GROW-1396
// Test generation with config name customization
func _TestGenerationAzureNamedConfig(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	var configName string = "Test-Config-Rename"

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, configName},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(true, false, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
		azure.WithConfigIntegrationName(configName),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Unstable test disabled as part of GROW-1396
// Test generation with activity log customization
func _TestGenerationAzureNamedActivityLog(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	var activityName string = "Test Activity Log Rename"

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(false, true, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
		azure.WithActivityLogIntegrationName(activityName),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test existing main.tf prompt
func TestGenerationAzureWithExistingTerraform(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	//var final string
	var runError error

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
	runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, dir},
				MsgRsp{fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), "n"},

				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "y"},
				MsgRsp{cmd.QuestionEnableAllSubscriptions, "y"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(true, false, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "y"},
				MsgRsp{cmd.QuestionManagementGroupId, mgmtGrpId},
				MsgRsp{cmd.QuestionEnableActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(true, false, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "y"},
				MsgRsp{cmd.QuestionEnableAllSubscriptions, "n"},
				MsgRsp{cmd.QuestionSubscriptionIds, strings.Join(testIds[:], ",")},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(true, false, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "y"},
				MsgRsp{cmd.QuestionEnableAllSubscriptions, "n"},
				MsgRsp{cmd.QuestionSubscriptionIds, strings.Join(testIds[:], ",")},
				MsgRsp{cmd.QuestionAzureEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(false, true, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "y"},
				MsgRsp{cmd.QuestionStorageAccountName, storageAccountName},
				MsgRsp{cmd.QuestionStorageAccountResourceGroup, storageResourceGrp},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(false, true, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "y"},
				MsgRsp{cmd.QuestionEnableAllSubscriptions, "y"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(false, true, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, region},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(false, true, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
		azure.WithStorageLocation(region),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAzureOverwrite(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", dir))

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, dir},
				MsgRsp{fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", dir))
}

func TestGenerationAzureOverwriteOutput(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	output_dir := createDummyTOMLConfig()
	defer os.RemoveAll(output_dir)

	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)

	// First run - generate initial code
	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, output_dir},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))

	// Second run - attempt to overwrite
	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, output_dir},
				MsgRsp{fmt.Sprintf("%s/main.tf already exists, overwrite?", output_dir), "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))
}

func TestGenerationAzureLaceworkProfile(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	azProfile := "v2"

	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"az",
		"--profile",
		azProfile,
	)

	assert.Nil(t, runError)
	assert.Contains(t, final, "Terraform code saved in")

	buildTf, _ := azure.NewTerraform(true, true, false, true,
		azure.WithSubscriptionID(mockSubscriptionID),
		azure.WithLaceworkProfile(azProfile),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAzureWithSubscriptionID(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error

	// Run CLI
	tfResult := runGenerateAzureTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionAzureSubscriptionID, mockSubscriptionID},
				MsgRsp{cmd.AzureSubscriptions, "n"},
				MsgRsp{cmd.QuestionAzureEnableConfig, "y"},
				MsgRsp{cmd.QuestionAzureConfigName, ""},
				MsgRsp{cmd.QuestionEnableManagementGroup, "n"},
				MsgRsp{cmd.QuestionEnableActivityLog, "y"},
				MsgRsp{cmd.QuestionActivityLogName, ""},
				MsgRsp{cmd.QuestionUseExistingStorageAccount, "n"},
				MsgRsp{cmd.QuestionStorageLocation, ""},
				MsgRsp{cmd.QuestionEnableEntraIdActivityLog, "n"},
				MsgRsp{cmd.QuestionEnableAdIntegration, "y"},
				MsgRsp{cmd.QuestionAzureCustomizeOutputLocation, ""},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
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
	buildTf, _ := azure.NewTerraform(true, true, false, true, azure.WithSubscriptionID(mockSubscriptionID)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func runGenerateAzureTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	os.Setenv("HOME", tfPath)

	hcl_path := filepath.Join(tfPath, azurePath, "main.tf")

	runFakeTerminalTestFromDir(t, tfPath, conditions, args...)
	out, err := os.ReadFile(hcl_path)
	if err != nil {
		return fmt.Sprintf("main.tf not found: %s", err)
	}

	t.Cleanup(func() {
		os.Remove(hcl_path)
	})

	result := terraformValidate(filepath.Join(tfPath, azurePath))

	assert.True(t, result.Valid)

	return string(out)
}
