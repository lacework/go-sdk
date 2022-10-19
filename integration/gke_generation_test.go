//go:build !windows && generation

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/gcp"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/stretchr/testify/assert"
)

const (
	gkeOrgId  = "org-1"
	gkeProjId = "project-1"
)

func assertGkeTerraformSaved(t *testing.T, message string) {
	assert.Contains(t, message, "Terraform code saved in")
}

func runGkeGenerateTest(t *testing.T, conditions func(*expect.Console), args ...string) string {
	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	runFakeTerminalTestFromDir(t, dir, conditions, args...)
	out, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/lacework/gke/main.tf", dir)))
	if err != nil {
		return ""
	}

	return string(out)
}

func TestGenerationSimple(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assertGkeTerraformSaved(t, final)

	buildTf, _ := gcp.NewGkeTerraform(gcp.WithGkeProjectId(gkeProjId)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationOrg(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "y"},
				MsgRsp{cmd.QuestionGcpOrganizationID, gkeOrgId},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assertGkeTerraformSaved(t, final)

	buildTf, _ := gcp.NewGkeTerraform(
		gcp.WithGkeProjectId(gkeProjId),
		gcp.WithGkeOrganizationIntegration(true),
		gcp.WithGkeOrganizationId(gkeOrgId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationSACredsGke(t *testing.T) {
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

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, serviceAccountFilePath},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assertGkeTerraformSaved(t, final)

	buildTf, _ := gcp.NewGkeTerraform(
		gcp.WithGkeProjectId(gkeProjId),
		gcp.WithGkeServiceAccountCredentials(serviceAccountFilePath),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationExistingSink(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GkeAdvancedOpt, 0},
				MsgRsp{cmd.QuestionGcpUseExistingSink, "y"},
				MsgRsp{cmd.QuestionGcpExistingSinkName, "sink"},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assertGkeTerraformSaved(t, final)

	buildTf, _ := gcp.NewGkeTerraform(
		gcp.WithGkeProjectId(gkeProjId),
		gcp.WithGkeExistingSinkName("sink"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationGkeAdvancedOptsUseExistingSA(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GkeAdvancedOpt, 1},
				MsgRsp{cmd.QuestionExistingServiceAccountName, "SA_1"},
				MsgRsp{cmd.QuestionExistingServiceAccountPrivateKey, "cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ=="},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assertGkeTerraformSaved(t, final)

	serviceAccount := &gcp.ServiceAccount{}
	serviceAccount.Name = "SA_1"
	serviceAccount.PrivateKey = "cGFzc3dvcmRNY1Bhc3N3b3JkRmFjZQ=="

	buildTf, _ := gcp.NewGkeTerraform(
		gcp.WithGkeProjectId(gkeProjId),
		gcp.WithGkeExistingServiceAccount(serviceAccount),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationCustomizedConfigurationIntegrationNameGke(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GkeAdvancedOpt, 2},
				MsgRsp{cmd.QuestionGkeIntegrationName, "customIntegrationName"},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	// Ensure CLI ran correctly
	assertGkeTerraformSaved(t, final)

	// Create the TF directly with lwgenerate and validate same result via CLI
	buildTf, _ := gcp.NewGkeTerraform(
		gcp.WithGkeProjectId(gkeProjId),
		gcp.WithGkeIntegrationName("customIntegrationName"),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationCustomizedOutputLocationGke(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GkeAdvancedOpt, 3},
				MsgRsp{cmd.QuestionGcpCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assertGkeTerraformSaved(t, final)

	result, _ := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)))

	buildTf, _ := gcp.NewGkeTerraform(
		gcp.WithGkeProjectId(gkeProjId),
	).Generate()
	assert.Equal(t, buildTf, string(result))
}

func TestGenerationAdvancedOptsDoneGke(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GkeAdvancedOpt, 4},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assertGkeTerraformSaved(t, final)

	buildTf, _ := gcp.NewGkeTerraform(
		gcp.WithGkeProjectId(gkeProjId),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationWithGkeExistingTerraformGke(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	dir, err := ioutil.TempDir("", "lacework-cli")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/main.tf", dir)), []byte{}, 0644); err != nil {
		panic(err)
	}

	runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "y"},
				MsgMenu{cmd.GkeAdvancedOpt, 3},
				MsgRsp{cmd.QuestionGcpCustomizeOutputLocation, dir},
				MsgRsp{cmd.QuestionGcpAnotherAdvancedOpt, "n"},
				MsgRsp{fmt.Sprintf("%s/main.tf already exists, overwrite?", dir), "n"},
			})
		},
		"generate",
		"k8s",
		"gke",
	)

	// Ensure CLI ran correctly
	data, err := os.ReadFile(fmt.Sprintf("%s/main.tf", dir))
	if err != nil {
		panic(err)
	}

	assert.Empty(t, data)
}

func TestGenerateGkePrefixAndWait(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	prefix := "prefix-"
	waitTime := "30s"

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGcpProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGcpOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGcpServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGcpConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})

			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
		"--prefix",
		prefix,
		"--wait_time",
		waitTime,
	)

	assertGkeTerraformSaved(t, final)

	buildTf, _ := gcp.NewGkeTerraform(
		gcp.WithGkeProjectId(gkeProjId),
		gcp.WithGkePrefix(prefix),
		gcp.WithGkeWaitTime(waitTime),
	).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationGkeInvalidProjectId(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	tfResult := runGkeGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGkeProjectID, "1"},
				MsgOnly{cmd.InvalidProjectIDMessage},
				MsgRsp{cmd.QuestionGkeProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGkeOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGkeServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGkeConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assertGkeTerraformSaved(t, final)

	buildTf, _ := gcp.NewGkeTerraform(gcp.WithGkeProjectId(gkeProjId)).Generate()
	assert.Equal(t, buildTf, tfResult)
}

func TestGenerationGkeOverwrite(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createDummyTOMLConfig()
	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)
	defer os.RemoveAll(dir)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGkeProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGkeOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGkeServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGkeConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s/lacework/gke", dir))

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGkeProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGkeOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGkeServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGkeConfigureAdvanced, "n"},
				MsgRsp{"already exists, overwrite?", "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s/lacework/gke", dir))
}

func TestGenerationGkeOverwriteOutput(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createDummyTOMLConfig()
	defer os.RemoveAll(dir)

	homeCache := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", homeCache)

	output_dir := createDummyTOMLConfig()
	defer os.RemoveAll(output_dir)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGkeProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGkeOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGkeServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGkeConfigureAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
		"--output",
		output_dir,
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionGkeProjectID, gkeProjId},
				MsgRsp{cmd.QuestionGkeOrganizationIntegration, "n"},
				MsgRsp{cmd.QuestionGkeServiceAccountCredsPath, ""},
				MsgRsp{cmd.QuestionGkeConfigureAdvanced, "n"},
				MsgRsp{"already exists, overwrite?", "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"generate",
		"k8s",
		"gke",
		"--output",
		output_dir,
	)

	assert.Contains(t, final, fmt.Sprintf("cd %s", output_dir))
}
