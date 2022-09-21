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
	"github.com/lacework/go-sdk/lwgenerate/gcp"
	"github.com/stretchr/testify/assert"
)

// Test failing due to no selection
func TestGenerationErrorOnNoSelectionGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// Run CLI
	runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("n")
			expectString(t, c, "ERROR collecting/confirming parameters: must enable audit log or configuration")
		},
		"cloud",
		"iac",
		"gcp",
	)
}

// Test bare-bones generation with no customization
func TestGenerationSimpleGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuration only generation
func TestGenerationConfigOnlyGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, false,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test auditlog only generation
func TestGenerationAuditlogOnlyGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(false, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAuditlogEnableUBLA(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
		"--enable_ubla",
	)

	assert.Contains(t, final, "Terraform code saved in")

	buildTf, _ := gcp.NewTerraform(false, true,
		gcp.WithProjectId("project-1"),
		gcp.WithEnableUBLA(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAuditlogDisableUBLA(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
		"--enable_ubla=false",
	)

	assert.Contains(t, final, "Terraform code saved in")

	buildTf, _ := gcp.NewTerraform(false, true,
		gcp.WithProjectId("project-1"),
		gcp.WithEnableUBLA(false),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test organization integration. configuration & audit log
func TestOrganizationIntegrationConfigAndAuditLogGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"
	organizationId := "org-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpOrganizationID)
			c.SendLine(organizationId)
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithOrganizationIntegration(true),
		gcp.WithOrganizationId(organizationId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generation with Service Account JSON file path
func TestGenerationSACredsGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"
	serviceAccountCreds := []byte(`{
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvwIBADANBgkqhkiG9\n97zZ15XlqSAd5r7q2DasgMWYriEHSJb2V6xwvji5kYeV6U\nY5PR+mPfVbb4xX3UMzwUEvK0cw==\n-----END PRIVATE KEY-----\n",
			"client_email": "test_email@lacework.iam.gserviceaccount.com"
	}`)

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	serviceAccountFilePath := filepath.Join(dir, "service_account_creds.json")

	if err := os.WriteFile(serviceAccountFilePath, serviceAccountCreds, 0644); err != nil {
		panic(err)
	}

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine(serviceAccountFilePath)
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithGcpServiceAccountCredentials(serviceAccountFilePath),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Use existing Bucket.
func TestGenerationAdvancedAuditLogOptsExistingBucketGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[C")
			expectString(t, c, cmd.QuestionGcpUseExistingBucket)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpExistingBucketName)
			c.SendLine("bucketMcBucketFace")
			expectString(t, c, cmd.QuestionGcpUseExistingSink)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithExistingLogBucketName("bucketMcBucketFace"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Create & don't configure new Bucket.
func TestGenerationAdvancedAuditLogOptsNewBucketNotConfiguredGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[C")
			expectString(t, c, cmd.QuestionGcpUseExistingBucket)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpConfigureNewBucket)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpUseExistingSink)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Create & configure new Bucket.
func TestGenerationAdvancedAuditLogOptsNewBucketConfiguredGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[C")
			expectString(t, c, cmd.QuestionGcpUseExistingBucket)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpConfigureNewBucket)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpBucketName)
			c.SendLine("newBucketMcBucketFace")
			expectString(t, c, cmd.QuestionGcpBucketRegion)
			c.SendLine("us-west1")
			expectString(t, c, cmd.QuestionGcpBucketLocation)
			c.SendLine("us")
			expectString(t, c, cmd.QuestionGcpBucketRetention)
			c.SendLine("10")
			expectString(t, c, cmd.QuestionGcpBucketLifecycle)
			c.SendLine("420")
			expectString(t, c, cmd.QuestionGcpEnableUBLA)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpUseExistingSink)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithBucketName("newBucketMcBucketFace"),
		gcp.WithBucketRegion("us-west1"),
		gcp.WithBucketLocation("us"),
		gcp.WithLogBucketRetentionDays(10),
		gcp.WithLogBucketLifecycleRuleAge(420),
		gcp.WithEnableUBLA(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Use existing sink.
func TestGenerationAdvancedAuditLogOptsExistingSinkGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[C")
			expectString(t, c, cmd.QuestionGcpUseExistingBucket)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpConfigureNewBucket)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpBucketName)
			c.SendLine("newBucketMcBucketFace")
			expectString(t, c, cmd.QuestionGcpBucketRegion)
			c.SendLine("us-west1")
			expectString(t, c, cmd.QuestionGcpBucketLocation)
			c.SendLine("us")
			expectString(t, c, cmd.QuestionGcpBucketRetention)
			c.SendLine("10")
			expectString(t, c, cmd.QuestionGcpBucketLifecycle)
			c.SendLine("420")
			expectString(t, c, cmd.QuestionGcpEnableUBLA)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpUseExistingSink)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithBucketName("newBucketMcBucketFace"),
		gcp.WithBucketRegion("us-west1"),
		gcp.WithBucketLocation("us"),
		gcp.WithLogBucketRetentionDays(10),
		gcp.WithLogBucketLifecycleRuleAge(420),
		gcp.WithEnableUBLA(false),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test advanced options. Use existing Service Account details.
func TestGenerationAdvancedOptsUseExistingSA(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionExistingServiceAccountName)
			c.SendLine("SA_1")
			expectString(t, c, cmd.QuestionExistingServiceAccountPrivateKey)
			c.SendLine("cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ==")
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	serviceAccountDetails := &gcp.ExistingServiceAccountDetails{}
	serviceAccountDetails.Name = "SA_1"
	serviceAccountDetails.PrivateKey = "cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ=="

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithExistingServiceAccount(serviceAccountDetails),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test custom configuration integration name
func TestGenerationCustomizedConfigurationIntegrationNameGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key down x2 in ANSI
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionGcpConfigurationIntegrationName)
			c.SendLine("customConfigurationIntegrationName")
			expectString(t, c, cmd.QuestionGcpAuditLogIntegrationName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithConfigurationIntegrationName("customConfigurationIntegrationName"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test custom audit log integration name
func TestGenerationCustomizedAuditlogIntegrationNameGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key down x2 in ANSI
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionGcpConfigurationIntegrationName)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpAuditLogIntegrationName)
			c.SendLine("customAuditlogIntegrationName")
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithAuditLogIntegrationName("customAuditlogIntegrationName"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test customized output location
func TestGenerationCustomizedOutputLocationGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
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
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key down x3 in ANSI
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionGcpCustomizeOutputLocation)
			c.SendLine(dir)
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Get result
	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, string(result))
}

// Test Bailing out of Advanced Options
func TestGenerationAdvancedOptsDoneGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key down x3 in ANSI
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionRunTfPlan)
			c.SendLine("n")
			final, _ = c.ExpectEOF()
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	assert.Contains(t, final, "Terraform code saved in")

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test existing main.tf prompt
func TestGenerationWithExistingTerraformGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	projectId := "project-1"

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
	runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectString(t, c, cmd.QuestionGcpEnableConfiguration)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpEnableAuditLog)
			c.SendLine("y")
			expectString(t, c, cmd.QuestionGcpProjectID)
			c.SendLine(projectId)
			expectString(t, c, cmd.QuestionGcpOrganizationIntegration)
			c.SendLine("n")
			expectString(t, c, cmd.QuestionGcpServiceAccountCredsPath)
			c.SendLine("")
			expectString(t, c, cmd.QuestionGcpConfigureAdvanced)
			c.SendLine("y")
			expectString(t, c, cmd.GcpAdvancedOptAuditLog)
			// This is key down x3 in ANSI
			c.Send("\x1B[B")
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectString(t, c, cmd.QuestionGcpCustomizeOutputLocation)
			c.SendLine(dir)
			expectString(t, c, cmd.QuestionGcpAnotherAdvancedOpt)
			c.SendLine("n")
			expectString(t, c, fmt.Sprintf("%s/main.tf already exists, overwrite?", dir))
			c.SendLine("n")
		},
		"cloud",
		"iac",
		"gcp",
	)

	// Ensure CLI ran correctly
	data, err := os.ReadFile(fmt.Sprintf("%s/main.tf", dir))
	if err != nil {
		panic(err)
	}

	assert.Empty(t, data)
}

func runGcpGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	// create a temporal directory where we will check that the
	// configuration file is deployed (.lacework.toml)
	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	runFakeTerminalTestFromDir(t, dir, conditions, args...)
	out, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/gcp/main.tf", dir)))
	if err != nil {
		// Assume couldn't be found
		return ""
	}

	return string(out)
}
