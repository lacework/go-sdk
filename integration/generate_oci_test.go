package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwgenerate/oci"
	"github.com/stretchr/testify/assert"
)

// Test failing due to no selection
func TestGenerateOciErrorOnNoSelection(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")

	// Run CLI
	runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionOciEnableConfig, "n"},
				MsgOnly{"ERROR collecting/confirming parameters: must enable config to continue"},
			})
		},
		withCloud("oci"),
	)
}

// Test OCI TF generation with minimal customization
func TestGenerateOciBasic(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	tenantOcid := "ocid1.tenancy...abc"
	userEmail := "test@example.com"

	// Run CLI
	actual := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionOciEnableConfig, "y"},
				MsgRsp{cmd.QuestionOciTenantOcid, tenantOcid},
				MsgRsp{cmd.QuestionOciUserEmail, userEmail},
				MsgRsp{cmd.QuestionOciConfigAdvanced, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		withCloud("oci"),
	)

	assert.Contains(t, final, "Terraform code saved in")

	expected, _ := oci.NewTerraform(
		true,
		false,
		oci.WithTenantOcid(tenantOcid),
		oci.WithUserEmail(userEmail),
	).Generate()
	assert.Equal(t, expected, actual)
}

// Test OCI TF generation with advanced customization
func TestGenerateOciCustomConfigName(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	tenantOcid := "ocid1.tenancy...abc"
	userEmail := "test@example.com"
	configName := "test_integration_oci"

	actual := runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionOciEnableConfig, "y"},
				MsgRsp{cmd.QuestionOciTenantOcid, tenantOcid},
				MsgRsp{cmd.QuestionOciUserEmail, userEmail},
				MsgRsp{cmd.QuestionOciConfigAdvanced, "y"},
				MsgMenu{cmd.OciAdvancedOptIntegrationName, 0},
				MsgRsp{cmd.QuestionOciConfigName, configName},
				MsgRsp{cmd.QuestionOciAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		withCloud("oci"),
	)

	assert.Contains(t, final, "Terraform code saved in")

	expected, _ := oci.NewTerraform(
		true,
		false,
		oci.WithTenantOcid(tenantOcid),
		oci.WithUserEmail(userEmail),
		oci.WithConfigName(configName),
	).Generate()
	assert.Equal(t, expected, actual)
}

// Test OCI TF generation with minimal customization
func TestGenerateOciCustomLocation(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	location, err := os.MkdirTemp("", "t")
	assert.Nil(t, err, "failed to create temporary directory")

	_ = runGenerateTest(t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionOciEnableConfig, "y"},
				MsgRsp{cmd.QuestionOciTenantOcid, "ocid1.tenancy...abc"},
				MsgRsp{cmd.QuestionOciUserEmail, "test@example.com"},
				MsgRsp{cmd.QuestionOciConfigAdvanced, "y"},
				MsgMenu{cmd.OciAdvancedOptLocation, 1},
				MsgRsp{cmd.QuestionOciCustomizeOutputLocation, location},
				MsgRsp{cmd.QuestionOciAnotherAdvancedOpt, "n"},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		withCloud("oci"),
		withLocation(location),
	)

	assert.Contains(t, final, fmt.Sprintf("Terraform code saved in %s", location))
}

// Test noninteractive with insufficient flags
func TestGenerateOciNoninteractiveNoFlags(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig("generate", "cloud-account", "oci", "--noninteractive")
	assert.Contains(t, err.String(), "ERROR collecting/confirming parameters: must enable config to continue")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

// Test noninteractive with insufficient flags
func TestGenerateOciNoninteractiveOnlyTenancyOcidFlag(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig(
		"generate",
		"cloud-account",
		"oci",
		"--noninteractive",
		"--tenant_ocid",
		"ocid1.tenancy...a",
	)
	assert.Contains(t, err.String(), "ERROR collecting/confirming parameters: must enable config to continue")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

// Test noninteractive with insufficient flags
func TestGenerateOciNoninteractiveOnlyUserEmail(t *testing.T) {
	_, err, exitcode := LaceworkCLIWithTOMLConfig(
		"generate",
		"cloud-account",
		"oci",
		"--noninteractive",
		"--oci_user_email",
		"a@b.c",
	)
	assert.Contains(t, err.String(), "ERROR collecting/confirming parameters: must enable config to continue")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

// test noninteractive with minimal flags
func TestGenerateOciNoninteractiveBasic(t *testing.T) {
	tenantOcid := "ocid1.tenancy...a"
	userEmail := "a@b.c"

	outputLocation, err := os.MkdirTemp("", "t")
	assert.Nil(t, err, "failed to create temporary directory")
	t.Cleanup(func() {
		os.RemoveAll(outputLocation)
	})

	_, stdErr, exitcode := LaceworkCLIWithTOMLConfig(
		"generate",
		"cloud-account",
		"oci",
		"--noninteractive",
		"--config",
		"--tenant_ocid",
		tenantOcid,
		"--oci_user_email",
		userEmail,
		"--output",
		outputLocation,
	)
	assert.Empty(t, stdErr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// tf validate
	tfValidateResult := terraformValidate(outputLocation)
	assert.True(t, tfValidateResult.Valid)

	// compare results to calling lwgenerate directly
	actual, err := os.ReadFile(filepath.Join(outputLocation, "main.tf"))
	assert.Nil(t, err, "error reading Terraform output")
	expected, _ := oci.NewTerraform(
		true,
		false,
		oci.WithTenantOcid(tenantOcid),
		oci.WithUserEmail(userEmail),
	).Generate()
	assert.Equal(t, expected, string(actual))
}

// test noninteractive with custom integration name
func TestGenerateOciNoninteractiveCustomConfigName(t *testing.T) {
	tenantOcid := "ocid1.tenancy...a"
	userEmail := "a@b.c"
	configName := "test_integration_oci"

	outputLocation, err := os.MkdirTemp("", "t")
	assert.Nil(t, err, "failed to create temporary directory")
	t.Cleanup(func() {
		os.RemoveAll(outputLocation)
	})

	_, stdErr, exitcode := LaceworkCLIWithTOMLConfig(
		"generate",
		"cloud-account",
		"oci",
		"--noninteractive",
		"--config",
		"--tenant_ocid",
		tenantOcid,
		"--oci_user_email",
		userEmail,
		"--output",
		outputLocation,
		"--config_name",
		configName,
	)
	assert.Empty(t, stdErr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// tf validate
	tfValidateResult := terraformValidate(outputLocation)
	assert.True(t, tfValidateResult.Valid)

	// compare to calling lwgenerate directly
	actual, err := os.ReadFile(filepath.Join(outputLocation, "main.tf"))
	assert.Nil(t, err, "error reading Terraform output")
	expected, _ := oci.NewTerraform(
		true,
		false,
		oci.WithTenantOcid(tenantOcid),
		oci.WithUserEmail(userEmail),
		oci.WithConfigName(configName),
	).Generate()
	assert.Equal(t, expected, string(actual))
}
