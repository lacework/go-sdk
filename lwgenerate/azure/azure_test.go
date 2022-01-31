package azure_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/azure"
	"github.com/stretchr/testify/assert"
)

func getFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
func TestGenerationActivityLogWithoutConfig(t *testing.T) {
	ActivityLogWithoutConfig, fileErr := getFileContent("activity_log_without_config.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithoutConfig, hcl)
}

func TestGenerationActivityLogWithConfig(t *testing.T) {
	var ActivityLogWithConfig, fileErr = getFileContent("activity_log_with_config.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithConfig, hcl)
}

func TestGenerationConfigWithoutActivityLog(t *testing.T) {
	ConfigWithoutActivityLog, fileErr := getFileContent("config_without_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigWithoutActivityLog, hcl)
}

func TestGenerationWithoutActivityLogOrConfig(t *testing.T) {
	hcl, err := azure.NewTerraform(false, false).Generate()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(errors.Unwrap(err).Error(), "invalid inputs"))
	assert.Empty(t, hcl)
}
func TestGenerationRenamedConfig(t *testing.T) {
	RenamedConfig, fileErr := getFileContent("renamed_config.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false,
		azure.WithConfigIntegrationName("Test Config Rename"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedConfig, hcl)
}

func TestGenerationRenamedActivityLog(t *testing.T) {
	RenamedActivityLog, fileErr := getFileContent("renamed_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true,
		azure.WithAuditLogIntegrationName("Test Activity Log Rename"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedActivityLog, hcl)
}

func TestGenerationRenamedConfigAndActivityLog(t *testing.T) {
	RenamedConfigAndActivityLog, fileErr := getFileContent("renamed_config_and_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true,
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithAuditLogIntegrationName("Test Activity Log Rename"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedConfigAndActivityLog, hcl)
}
