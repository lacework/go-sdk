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

// Test auditlog only generation
func TestGenerationAuditlogOnlyGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	projectId := "project-1"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectStringGcp(c, cmd.QuestionGcpEnableConfig, &runError)
			c.SendLine("n")
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
	buildTf, _ := gcp.NewTerraform(false, true,
		gcp.WithProjectId("project-1"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test organization integration. config & auditlog
func TestOrganizationIntegrationConfigAndAuditLogGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	var runError error
	projectId := "project-1"
	organizationId := "org-1"

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
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpOrganizationID, &runError)
			c.SendLine(organizationId)
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
	var runError error
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
			expectStringGcp(c, cmd.QuestionGcpEnableConfig, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpEnableAuditLog, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpProjectID, &runError)
			c.SendLine(projectId)
			expectStringGcp(c, cmd.QuestionGcpOrganizationIntegration, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpServiceAccountCredsPath, &runError)
			c.SendLine(serviceAccountFilePath)
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
		gcp.WithGcpServiceAccountCredentials(serviceAccountFilePath),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Use existing Bucket.
func TestGenerationAdvancedAuditLogOptsExistingBucketGcp(t *testing.T) {
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
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[C")
			expectStringGcp(c, cmd.QuestionGcpUseExistingBucket, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpExistingBucketName, &runError)
			c.SendLine("bucketMcBucketFace")
			expectStringGcp(c, cmd.QuestionGcpUseExistingSink, &runError)
			c.SendLine("n")
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
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[C")
			expectStringGcp(c, cmd.QuestionGcpUseExistingBucket, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpConfigureNewBucket, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpUseExistingSink, &runError)
			c.SendLine("n")
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
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[C")
			expectStringGcp(c, cmd.QuestionGcpUseExistingBucket, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpConfigureNewBucket, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpBucketName, &runError)
			c.SendLine("newBucketMcBucketFace")
			expectStringGcp(c, cmd.QuestionGcpBucketRegion, &runError)
			c.SendLine("us-west1")
			expectStringGcp(c, cmd.QuestionGcpBucketLocation, &runError)
			c.SendLine("us")
			expectStringGcp(c, cmd.QuestionGcpBucketRetention, &runError)
			c.SendLine("10")
			expectStringGcp(c, cmd.QuestionGcpBucketLifecycle, &runError)
			c.SendLine("420")
			expectStringGcp(c, cmd.QuestionGcpEnableUBLA, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpUseExistingSink, &runError)
			c.SendLine("n")
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

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithBucketName("newBucketMcBucketFace"),
		gcp.WithBucketRegion("us-west1"),
		gcp.WithBucketLocation("us"),
		gcp.WithLogBucketRetentionDays(10),
		gcp.WithLogBucketLifecycleRuleAge(420),
		gcp.WithEnableUBLA(),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Use existing sink.
func TestGenerationAdvancedAuditLogOptsExistingSinkGcp(t *testing.T) {
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
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[C")
			expectStringGcp(c, cmd.QuestionGcpUseExistingBucket, &runError)
			c.SendLine("n")
			expectStringGcp(c, cmd.QuestionGcpConfigureNewBucket, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpBucketName, &runError)
			c.SendLine("newBucketMcBucketFace")
			expectStringGcp(c, cmd.QuestionGcpBucketRegion, &runError)
			c.SendLine("us-west1")
			expectStringGcp(c, cmd.QuestionGcpBucketLocation, &runError)
			c.SendLine("us")
			expectStringGcp(c, cmd.QuestionGcpBucketRetention, &runError)
			c.SendLine("10")
			expectStringGcp(c, cmd.QuestionGcpBucketLifecycle, &runError)
			c.SendLine("420")
			expectStringGcp(c, cmd.QuestionGcpEnableUBLA, &runError)
			c.SendLine("y")
			expectStringGcp(c, cmd.QuestionGcpUseExistingSink, &runError)
			c.SendLine("n")
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

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithBucketName("newBucketMcBucketFace"),
		gcp.WithBucketRegion("us-west1"),
		gcp.WithBucketLocation("us"),
		gcp.WithLogBucketRetentionDays(10),
		gcp.WithLogBucketLifecycleRuleAge(420),
		gcp.WithEnableUBLA(),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test advanced options. Use existing Service Account details.
func TestGenerationAdvancedOptsUseExistingSA(t *testing.T) {
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
			//expectStringGcp(c, cmd.GcpAdvancedOptAuditLog, &runError)
			//// This is key down x1 in ANSI
			//c.SendLine("\x1B[B")
			expectStringGcp(c, cmd.GcpAdvancedOptAuditLog, &runError)
			// This is key forward x1 in ANSI
			c.SendLine("\x1B[B")
			expectStringGcp(c, cmd.QuestionExistingServiceAccountName, &runError)
			c.SendLine("SA_1")
			expectStringGcp(c, cmd.QuestionExistingServiceAccountPrivateKey, &runError)
			c.SendLine("cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ==")
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

// Test custom config integration name
func TestGenerationCustomizedConfigIntegrationNameGcp(t *testing.T) {
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
			// This is key down x2 in ANSI
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectStringGcp(c, cmd.QuestionGcpConfigIntegrationName, &runError)
			c.SendLine("customConfigIntegrationName")
			expectStringGcp(c, cmd.QuestionGcpAuditLogIntegrationName, &runError)
			c.SendLine("")
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

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId("project-1"),
		gcp.WithConfigIntegrationName("customConfigIntegrationName"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test custom audit log integration name
func TestGenerationCustomizedAuditlogIntegrationNameGcp(t *testing.T) {
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
			// This is key down x2 in ANSI
			c.Send("\x1B[B")
			c.SendLine("\x1B[B")
			expectStringGcp(c, cmd.QuestionGcpConfigIntegrationName, &runError)
			c.SendLine("")
			expectStringGcp(c, cmd.QuestionGcpAuditLogIntegrationName, &runError)
			c.SendLine("customAuditlogIntegrationName")
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

// Test existing main.tf prompt
func TestGenerationWithExistingTerraformGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var runError error
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
			expectStringGcp(c, fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), &runError)
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
	assert.Nil(t, runError)
}

func runGcpGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	// create a temporal directory where we will check that the
	// configuration file is deployed (.lacework.toml)
	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	runGenerationTestFromDir(t, dir, conditions, args...)
	out, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/gcp/main.tf", dir)))
	if err != nil {
		// Assume couldn't be found
		return ""
	}

	return string(out)
}
