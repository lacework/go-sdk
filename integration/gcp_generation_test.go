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

const (
	organizationId = "org-1"
	projectId      = "project-1"
)

func assertTerraformSaved(t *testing.T, message string) {
	assert.Contains(t, message, "Terraform code saved in")
}

// Test failing due to no selection
func TestGenerationErrorOnNoSelectionGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// Run CLI
	runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "n"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "n"),
				msgOnly("ERROR collecting/confirming parameters: must enable audit log or configuration"),
			})
		},
		"generate",
		"cloud-account",
		"gcp",
	)
}

// Test bare-bones generation with no customization
func TestGenerationSimpleGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuration only generation
func TestGenerationConfigOnlyGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "n"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, false,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test auditlog only generation
func TestGenerationAuditlogOnlyGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "n"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true,
		gcp.WithProjectId(projectId),
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
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "n"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
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
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "n"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
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

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "y"),
				msgRsp(cmd.QuestionGcpOrganizationID, organizationId),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithOrganizationIntegration(true),
		gcp.WithOrganizationId(organizationId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGeneratePrefixAndWait(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	prefix := "prefix-"
	waitTime := "30s"

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--prefix",
		prefix,
		"--wait_time",
		waitTime,
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithPrefix(prefix),
		gcp.WithWaitTime(waitTime),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test generation with Service Account JSON file path
func TestGenerationSACredsGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
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
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, serviceAccountFilePath),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithGcpServiceAccountCredentials(serviceAccountFilePath),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Use existing Bucket.
func TestGenerationAdvancedAuditLogOptsExistingBucketGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgRsp(cmd.GcpAdvancedOptAuditLog, ""),
				msgRsp(cmd.QuestionGcpUseExistingBucket, "y"),
				msgRsp(cmd.QuestionGcpExistingBucketName, "bucketMcBucketFace"),
				msgRsp(cmd.QuestionGcpUseExistingSink, "n"),
				msgRsp(cmd.QuestionGcpCustomFilter, ""),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})

			final, _ = c.ExpectEOF()

		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithExistingLogBucketName("bucketMcBucketFace"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Create & don't configure new Bucket.
func TestGenerationAdvancedAuditLogOptsNewBucketNotConfiguredGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgRsp(cmd.GcpAdvancedOptAuditLog, ""),
				msgRsp(cmd.QuestionGcpUseExistingBucket, "n"),
				msgRsp(cmd.QuestionGcpConfigureNewBucket, "n"),
				msgRsp(cmd.QuestionGcpUseExistingSink, "n"),
				msgRsp(cmd.QuestionGcpCustomFilter, ""),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test configuring advanced audit log options. Create & configure new Bucket.
func TestGenerationAdvancedAuditLogOptsNewBucketConfiguredGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	bucketName := "my-new-bucket"

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgRsp(cmd.GcpAdvancedOptAuditLog, ""),
				msgRsp(cmd.QuestionGcpUseExistingBucket, "n"),
				msgRsp(cmd.QuestionGcpConfigureNewBucket, "y"),
				msgRsp(cmd.QuestionGcpBucketRegion, "us-west1"),
				msgRsp(cmd.QuestionGcpCustomBucketName, bucketName),
				msgRsp(cmd.QuestionGcpBucketLifecycle, "420"),
				msgRsp(cmd.QuestionGcpEnableUBLA, "y"),
				msgRsp(cmd.QuestionGcpUseExistingSink, "n"),
				msgRsp(cmd.QuestionGcpCustomFilter, ""),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	// Ensure CLI ran correctly
	assertTerraformSaved(t, final)

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithBucketRegion("us-west1"),
		gcp.WithCustomBucketName(bucketName),
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

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgRsp(cmd.GcpAdvancedOptAuditLog, ""),
				msgRsp(cmd.QuestionGcpUseExistingBucket, "n"),
				msgRsp(cmd.QuestionGcpConfigureNewBucket, "y"),
				msgRsp(cmd.QuestionGcpBucketRegion, "us-west1"),
				msgRsp(cmd.QuestionGcpCustomBucketName, ""),
				msgRsp(cmd.QuestionGcpBucketLifecycle, "420"),
				msgRsp(cmd.QuestionGcpEnableUBLA, "y"),
				msgRsp(cmd.QuestionGcpUseExistingSink, "n"),
				msgRsp(cmd.QuestionGcpCustomFilter, ""),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	// Ensure CLI ran correctly
	assertTerraformSaved(t, final)

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithBucketRegion("us-west1"),
		gcp.WithLogBucketLifecycleRuleAge(420),
		gcp.WithEnableUBLA(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAdvancedAuditLogOpts(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	filter := "filter"

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "n"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgRsp(cmd.GcpAdvancedOptAuditLog, ""),
				msgRsp(cmd.QuestionGcpUseExistingBucket, "y"),
				msgRsp(cmd.QuestionGcpExistingBucketName, "bucketMcBucketFace"),
				msgRsp(cmd.QuestionGcpUseExistingSink, "n"),
				msgRsp(cmd.QuestionGcpCustomFilter, filter),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true,
		gcp.WithProjectId(projectId),
		gcp.WithExistingLogBucketName("bucketMcBucketFace"),
		gcp.WithCustomFilter(filter),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test advanced options. Use existing Service Account details.
func TestGenerationAdvancedOptsUseExistingSA(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgMenu(cmd.GcpAdvancedOptAuditLog, 1),
				msgRsp(cmd.QuestionExistingServiceAccountName, "SA_1"),
				msgRsp(cmd.QuestionExistingServiceAccountPrivateKey, "cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ=="),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	serviceAccountDetails := &gcp.ExistingServiceAccountDetails{}
	serviceAccountDetails.Name = "SA_1"
	serviceAccountDetails.PrivateKey = "cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ=="

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithExistingServiceAccount(serviceAccountDetails),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test custom configuration integration name
func TestGenerationCustomizedConfigurationIntegrationNameGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgMenu(cmd.GcpAdvancedOptAuditLog, 2),
				msgRsp(cmd.QuestionGcpConfigurationIntegrationName, "customConfigurationIntegrationName"),
				msgRsp(cmd.QuestionGcpAuditLogIntegrationName, ""),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	// Ensure CLI ran correctly
	assertTerraformSaved(t, final)

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithConfigurationIntegrationName("customConfigurationIntegrationName"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test custom audit log integration name
func TestGenerationCustomizedAuditlogIntegrationNameGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgMenu(cmd.GcpAdvancedOptAuditLog, 2),
				msgRsp(cmd.QuestionGcpConfigurationIntegrationName, ""),
				msgRsp(cmd.QuestionGcpAuditLogIntegrationName, "customAuditlogIntegrationName"),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithAuditLogIntegrationName("customAuditlogIntegrationName"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test customized output location
func TestGenerationCustomizedOutputLocationGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	// Tempdir for test
	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgMenu(cmd.GcpAdvancedOptAuditLog, 3),
				msgRsp(cmd.QuestionGcpCustomizeOutputLocation, dir),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, string(result))
}

// Test Bailing out of Advanced Options
func TestGenerationAdvancedOptsDoneGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgMenu(cmd.GcpAdvancedOptAuditLog, 4),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAdvancedOptsDoneGcpConfiguration(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "n"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgMenu(cmd.GcpAdvancedOptExistingServiceAccount, 3),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, false,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test existing main.tf prompt
func TestGenerationWithExistingTerraformGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

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

	runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "y"),
				msgMenu(cmd.GcpAdvancedOptAuditLog, 3),
				msgRsp(cmd.QuestionGcpCustomizeOutputLocation, dir),
				msgRsp(cmd.QuestionGcpAnotherAdvancedOpt, "n"),
				msgRsp(fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), "n"),
			})
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	// Ensure CLI ran correctly
	data, err := os.ReadFile(fmt.Sprintf("%s/main.tf", dir))
	if err != nil {
		panic(err)
	}

	assert.Empty(t, data)
}

func TestGenerationFolders(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "y"),
				msgRsp(cmd.QuestionGcpOrganizationID, organizationId),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--folders_to_include", "folder/abc",
		"--folders_to_include", "folder/def",
		"--folders_to_include", "folder/abc",
		"--folders_to_exclude", "folder/abc",
		"--folders_to_exclude", "folder/def",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithOrganizationIntegration(true),
		gcp.WithOrganizationId(organizationId),
		gcp.WithFoldersToExclude([]string{"folder/abc", "folder/def"}),
		gcp.WithFoldersToInclude([]string{"folder/abc", "folder/abc", "folder/def"}),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationFoldersShorthand(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "y"),
				msgRsp(cmd.QuestionGcpOrganizationID, organizationId),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"-i", "folder/abc",
		"-i", "folder/abc",
		"-i", "folder/def",
		"-e", "folder/abc",
		"-e", "folder/def",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithOrganizationIntegration(true),
		gcp.WithOrganizationId(organizationId),
		gcp.WithFoldersToExclude([]string{"folder/abc", "folder/def"}),
		gcp.WithFoldersToInclude([]string{"folder/abc", "folder/abc", "folder/def"}),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationIncludeRootProjects(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "y"),
				msgRsp(cmd.QuestionGcpOrganizationID, organizationId),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--folders_to_exclude",
		"folder/abc",
		"--include_root_projects",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithOrganizationIntegration(true),
		gcp.WithOrganizationId(organizationId),
		gcp.WithFoldersToExclude([]string{"folder/abc"}),
		gcp.WithIncludeRootProjects(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationIncludeRootProjectsFalse(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "y"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "y"),
				msgRsp(cmd.QuestionGcpOrganizationID, organizationId),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--folders_to_exclude",
		"folder/abc",
		"--include_root_projects=false",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true,
		gcp.WithProjectId(projectId),
		gcp.WithOrganizationIntegration(true),
		gcp.WithOrganizationId(organizationId),
		gcp.WithFoldersToExclude([]string{"folder/abc"}),
		gcp.WithIncludeRootProjects(false),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAuditLogFiltersTrue(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "n"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--google_workspace_filter",
		"--k8s_filter",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true,
		gcp.WithProjectId(projectId),
		gcp.WithGoogleWorkspaceFilter(true),
		gcp.WithK8sFilter(true),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationAuditlogFiltersFalse(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRsp{
				msgRsp(cmd.QuestionGcpEnableConfiguration, "n"),
				msgRsp(cmd.QuestionGcpEnableAuditLog, "y"),
				msgRsp(cmd.QuestionGcpProjectID, projectId),
				msgRsp(cmd.QuestionGcpOrganizationIntegration, "n"),
				msgRsp(cmd.QuestionGcpServiceAccountCredsPath, ""),
				msgRsp(cmd.QuestionGcpConfigureAdvanced, "n"),
				msgRsp(cmd.QuestionRunTfPlan, "n"),
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--google_workspace_filter=false",
		"--k8s_filter=false",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true,
		gcp.WithProjectId(projectId),
		gcp.WithGoogleWorkspaceFilter(false),
		gcp.WithK8sFilter(false),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
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
