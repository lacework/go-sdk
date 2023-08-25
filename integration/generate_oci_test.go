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

// run a test for lacework generate cloud-account oci
func runOciGenerateTest(t *testing.T, conditions func(*expect.Console)) string {
	outputLocation := filepath.Join(tfPath, "lacework/oci")

	os.Setenv("HOME", tfPath)
	runFakeTerminalTestFromDir(t, tfPath, conditions, "generate", "cloud-account", "oci")
	out, err := os.ReadFile(filepath.Join(outputLocation, "main.tf"))
	if err != nil {
		return fmt.Sprintf("main.tf not found: %s", err)
	}

	t.Cleanup(func() {
		os.Remove(filepath.Join(outputLocation, "main.tf"))
	})

	result := terraformValidate(outputLocation)

	assert.True(t, result.Valid)

	return string(out)
}

// Test OCI TF generation with minimal customization
func TestGenerateOciBasic(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string
	tenantOcid := "ocid1.tenancy...abc"
	userEmail := "test@example.com"

	// Run CLI
	actual := runOciGenerateTest(
		t,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.QuestionOciTenantOcid, tenantOcid},
				MsgRsp{cmd.QuestionOciUserEmail, userEmail},
				MsgRsp{cmd.QuestionRunTfPlan, "n"},
			})
			final, _ = c.ExpectEOF()
		},
	)

	assert.Contains(t, final, "Terraform code saved in")

	expected, _ := oci.NewTerraform(
		true,
		oci.WithTenantOcid(tenantOcid),
		oci.WithUserEmail(userEmail),
	).Generate()
	assert.Equal(t, expected, actual)
}
