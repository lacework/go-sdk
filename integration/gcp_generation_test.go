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
	gcpPath        = "/lacework/gcp/"
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "n"},
				MsgOnly{"ERROR collecting/confirming parameters: must enable audit log or configuration"},
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "n"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, false, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true, false,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test pub sub auditlog only generation
func TestGenerationPubSubAuditlogOnlyGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 0},
				MsgRsp{cmd.QuestionUsePubSubAudit, "y"},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "n"},
				MsgRsp{cmd.QuestionGcpCustomFilter, ""},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true, true,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationPubSubAuditlogOrgGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "y"},
				MsgRsp{cmd.QuestionGcpOrganizationID, organizationId},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 0},
				MsgRsp{cmd.QuestionUsePubSubAudit, "y"},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "n"},
				MsgRsp{cmd.QuestionGcpCustomFilter, ""},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true, true,
		gcp.WithProjectId(projectId),
		gcp.WithOrganizationIntegration(true),
		gcp.WithOrganizationId(organizationId),
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--enable_ubla",
	)

	assert.Contains(t, final, "Terraform code saved in")

	buildTf, _ := gcp.NewTerraform(false, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--enable_ubla=false",
	)

	assert.Contains(t, final, "Terraform code saved in")

	buildTf, _ := gcp.NewTerraform(false, true, false,
		gcp.WithProjectId("project-1"),
		gcp.WithEnableUBLA(false),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test organization integration. configuration & audit log
func TestOrganizationIntegrationAllIntegrationGcp(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "y"},
				MsgRsp{cmd.QuestionGcpOrganizationID, organizationId},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(true, true, false,
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

	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	serviceAccountFilePath := filepath.Join(dir, "service_account_creds.json")

	if err := os.WriteFile(serviceAccountFilePath, serviceAccountCreds, 0644); err != nil {
		panic(err)
	}

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, serviceAccountFilePath},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 0},
				MsgRsp{cmd.QuestionUsePubSubAudit, "n"},
				MsgRsp{cmd.QuestionGcpUseExistingBucket, "y"},
				MsgRsp{cmd.QuestionGcpExistingBucketName, "bucketMcBucketFace"},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "n"},
				MsgRsp{cmd.QuestionGcpCustomFilter, ""},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()

		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 0},
				MsgRsp{cmd.QuestionUsePubSubAudit, "n"},
				MsgRsp{cmd.QuestionGcpUseExistingBucket, "n"},
				MsgRsp{cmd.QuestionGcpConfigureNewBucket, "n"},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "n"},
				MsgRsp{cmd.QuestionGcpCustomFilter, ""},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
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

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 0},
				MsgRsp{cmd.QuestionUsePubSubAudit, "n"},
				MsgRsp{cmd.QuestionGcpUseExistingBucket, "n"},
				MsgRsp{cmd.QuestionGcpConfigureNewBucket, "y"},
				MsgRsp{cmd.QuestionGcpBucketRegion, "us-west1"},
				MsgRsp{cmd.QuestionGcpCustomBucketName, bucketName},
				MsgRsp{cmd.QuestionGcpBucketLifecycle, "420"},
				MsgRsp{cmd.QuestionGcpEnableUBLA, "y"},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "n"},
				MsgRsp{cmd.QuestionGcpCustomFilter, ""},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
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

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 0},
				MsgRsp{cmd.QuestionUsePubSubAudit, "n"},
				MsgRsp{cmd.QuestionGcpUseExistingBucket, "n"},
				MsgRsp{cmd.QuestionGcpConfigureNewBucket, "y"},
				MsgRsp{cmd.QuestionGcpBucketRegion, "us-west1"},
				MsgRsp{cmd.QuestionGcpCustomBucketName, ""},
				MsgRsp{cmd.QuestionGcpBucketLifecycle, "420"},
				MsgRsp{cmd.QuestionGcpEnableUBLA, "y"},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "y"},
				MsgRsp{cmd.QuestionGcpExistingSinkName, "sink"},
				MsgRsp{cmd.QuestionGcpCustomFilter, ""},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
		gcp.WithProjectId(projectId),
		gcp.WithBucketRegion("us-west1"),
		gcp.WithLogBucketLifecycleRuleAge(420),
		gcp.WithEnableUBLA(true),
		gcp.WithExistingLogSinkName("sink"),
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 0},
				MsgRsp{cmd.QuestionUsePubSubAudit, "n"},
				MsgRsp{cmd.QuestionGcpUseExistingBucket, "y"},
				MsgRsp{cmd.QuestionGcpExistingBucketName, "bucketMcBucketFace"},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "n"},
				MsgRsp{cmd.QuestionGcpCustomFilter, filter},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 1},
				MsgRsp{cmd.QuestionExistingServiceAccountName, "SA_1"},
				MsgRsp{cmd.QuestionExistingServiceAccountPrivateKey, "cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ=="},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(true, true, false,
		gcp.WithProjectId(projectId),
		gcp.WithExistingServiceAccount(serviceAccountDetails),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

// Test advanced options. Use existing Service Account details for pub sub audit.
func TestGenerationPubSubUseExistingSA(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 0},
				MsgRsp{cmd.QuestionUsePubSubAudit, "y"},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "n"},
				MsgRsp{cmd.QuestionGcpCustomFilter, ""},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 1},
				MsgRsp{cmd.QuestionExistingServiceAccountName, "SA_1"},
				MsgRsp{cmd.QuestionExistingServiceAccountPrivateKey, "cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ=="},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(true, true, true,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 2},
				MsgRsp{cmd.QuestionGcpConfigurationIntegrationName, "customConfigurationIntegrationName"},
				MsgRsp{cmd.QuestionGcpAuditLogIntegrationName, ""},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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
	buildTf, _ := gcp.NewTerraform(true, true, false,
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

	// Run CLI
	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 2},
				MsgRsp{cmd.QuestionGcpConfigurationIntegrationName, ""},
				MsgRsp{cmd.QuestionGcpAuditLogIntegrationName, "customAuditlogIntegrationName"},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
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

	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 3},
				MsgRsp{cmd.QuestionGcpCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 5},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "n"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptDone, 4},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, false, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptAuditLog, 3},
				MsgRsp{cmd.QuestionGcpCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), "n"},
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "y"},
				MsgRsp{cmd.QuestionGcpOrganizationID, organizationId},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "y"},
				MsgRsp{cmd.QuestionGcpOrganizationID, organizationId},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "y"},
				MsgRsp{cmd.QuestionGcpOrganizationID, organizationId},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "y"},
				MsgRsp{cmd.QuestionGcpOrganizationID, organizationId},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(true, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(false, true, false,
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
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
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

	buildTf, _ := gcp.NewTerraform(false, true, false,
		gcp.WithProjectId(projectId),
		gcp.WithGoogleWorkspaceFilter(false),
		gcp.WithK8sFilter(false),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationGcpInvalidProjectId(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, "1"},
				MsgOnly{cmd.InvalidProjectIDMessage},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(false, true, false,
		gcp.WithProjectId(projectId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationGcpOverwrite(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	os.Setenv("HOME", dir)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s/lacework/gcp", dir))

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{"already exists, overwrite?", "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s/lacework/gcp", dir))
}

func TestGenerationGcpOverwriteOutput(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	os.Setenv("HOME", dir)

	output_dir := createDummyTOMLConfig()
	defer os.RemoveAll(output_dir)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--output",
		output_dir,
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "n"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{"already exists, overwrite?", "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--output",
		output_dir,
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))
}

func TestGenerationGcpLaceworkProfile(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	gcpProfile := "v2"

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--profile",
		gcpProfile,
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
		gcp.WithProjectId(projectId),
		gcp.WithLaceworkProfile(gcpProfile),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationGcpMultipleProjects(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	gcpProjects := []string{"project1", "project2", "project3"}

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
		"--projects",
		"project1",
		"--projects",
		"project2",
		"--projects",
		"project3",
		"--projects",
		"project1",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
		gcp.WithProjectId(projectId),
		gcp.WithMultipleProject(gcpProjects),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationGcpMultipleProjectsInteractive(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	gcpProjects := []string{"project1", "project2", "project3"}

	tfResult := runGcpGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpEnableConfiguration, "y"},
				MsgRsp{cmd.QuestionGcpEnableAuditLog, "y"},
				MsgRsp{cmd.QuestionGcpProjectID, projectId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GcpAdvancedOptProjects, 4},
				MsgRsp{cmd.QuestionGcpCustomizeProjects, "project1, project2  ,project3"},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"cloud-account",
		"gcp",
	)

	assertTerraformSaved(t, final)

	buildTf, _ := gcp.NewTerraform(true, true, false,
		gcp.WithProjectId(projectId),
		gcp.WithMultipleProject(gcpProjects),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func runGcpGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	os.Setenv("HOME", tfPath)

	runFakeTerminalTestFromDir(t, tfPath, conditions, args...)
	out, err := ioutil.ReadFile(filepath.Join(tfPath, gcpPath, "main.tf"))
	if err != nil {
		return fmt.Sprintf("main.tf not found: %s", err)
	}

	t.Cleanup(func() {
		os.Remove(filepath.Join(tfPath, gcpPath, "main.tf"))
	})

	result := terraformValidate(filepath.Join(tfPath, gcpPath))

	assert.True(t, result.Valid)

	return string(out)
}
