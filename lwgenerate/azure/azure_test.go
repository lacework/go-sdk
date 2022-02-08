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
	ActivityLogWithoutConfig, fileErr := getFileContent("test-data/activity_log_without_config.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true,
		azure.WithAdIntegration(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithoutConfig, hcl)
}

func TestGenerationActivityLogWithConfig(t *testing.T) {
	var ActivityLogWithConfig, fileErr = getFileContent("test-data/activity_log_with_config.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true,
		azure.WithAdIntegration(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithConfig, hcl)
}

func TestGenerationConfigWithoutActivityLog(t *testing.T) {
	ConfigWithoutActivityLog, fileErr := getFileContent("test-data/config_without_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false,
		azure.WithAdIntegration(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigWithoutActivityLog, hcl)
}

func TestGenerationWithoutActivityLogOrConfig(t *testing.T) {
	hcl, err := azure.NewTerraform(false, false,
		azure.WithAdIntegration(true),
	).Generate()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(errors.Unwrap(err).Error(), "invalid inputs"))
	assert.Empty(t, hcl)
}
func TestGenerationRenamedConfig(t *testing.T) {
	RenamedConfig, fileErr := getFileContent("test-data/renamed_config.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false,
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithAdIntegration(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedConfig, hcl)
}

func TestGenerationRenamedActivityLog(t *testing.T) {
	RenamedActivityLog, fileErr := getFileContent("test-data/renamed_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true,
		azure.WithAuditLogIntegrationName("Test Activity Log Rename"),
		azure.WithAdIntegration(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedActivityLog, hcl)
}

func TestGenerationRenamedConfigAndActivityLog(t *testing.T) {
	RenamedConfigAndActivityLog, fileErr := getFileContent("test-data/renamed_config_and_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true,
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithAuditLogIntegrationName("Test Activity Log Rename"),
		azure.WithAdIntegration(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedConfigAndActivityLog, hcl)
}

func TestGenerationNoActiveDirectorySettings(t *testing.T) {
	hcl, err := azure.NewTerraform(true, true,
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithAuditLogIntegrationName("Test Activity Log Rename"),
	).Generate()
	assert.True(t, strings.Contains(errors.Unwrap(err).Error(), "invalid inputs"))
	assert.Empty(t, hcl)
}

func TestGenerationCustomActiveDirectory(t *testing.T) {
	CustomADDetails, fileErr := getFileContent("test-data/customer-ad-details.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true,
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithAuditLogIntegrationName("Test Activity Log Rename"),
		azure.WithAdIntegration(false),
		azure.WithAdApplicationPassword("AD-Test-Password"),
		azure.WithAdServicePrincipalId("AD-Test-Principal-ID"),
		azure.WithAdApplicationId("AD-Test-Application-ID"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, CustomADDetails, hcl)
}

func TestGenerationActivityLogWithExistingStorageAccount(t *testing.T) {
	ActivityLogWithStorage, fileErr := getFileContent("test-data/activity-log-with-existing-storage.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true,
		azure.WithAdIntegration(true),
		azure.WithExistingStorageAccount(true),
		azure.WithStorageAccountName("Test-Storage-Account-Name"),
		azure.WithStorageAccountResourceGroup("Test-Storage-Account-Resource-Group"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithStorage, hcl)
}

func TestGenerationActivityLogWithAllSubscriptions(t *testing.T) {
	ActivityLogAllSubs, fileErr := getFileContent("test-data/activity-log-with-all-subscriptions.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true,
		azure.WithAdIntegration(true),
		azure.WithAllSubscription(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogAllSubs, hcl)
}

func TestGenerationConfigWithAllSubscriptions(t *testing.T) {
	ConfigAllSubs, fileErr := getFileContent("test-data/config-with-all-subscriptions.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false,
		azure.WithAdIntegration(true),
		azure.WithAllSubscription(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigAllSubs, hcl)
}

func TestGenerationConfigWithManagementGroup(t *testing.T) {
	ConfigWithMgmtGroup, fileErr := getFileContent("test-data/config-with-management-group.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false,
		azure.WithAdIntegration(true),
		azure.WithManagementGroup(true),
		azure.WithManagementGroupId("test-management-group-1"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigWithMgmtGroup, hcl)
}

func TestGenerationConfigWithManagementGroupError(t *testing.T) {
	hcl, err := azure.NewTerraform(true, false,
		azure.WithAdIntegration(true),
		azure.WithManagementGroup(true),
	).Generate()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(errors.Unwrap(err).Error(), "invalid inputs"))
	assert.Empty(t, hcl)
}

func TestGenerationConfigWithStorageAccount(t *testing.T) {
	ConfigWithStorage, fileErr := getFileContent("test-data/config-with-storage-account.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false,
		azure.WithAdIntegration(true),
		azure.WithExistingStorageAccount(true),
		azure.WithStorageAccountName("test-storage-account-name"),
		azure.WithStorageAccountResourceGroup("test-storage-account-resource-group"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigWithStorage, hcl)
}

func TestGenerationConfigWithStorageAccountError(t *testing.T) {
	hcl, err := azure.NewTerraform(true, false,
		azure.WithAdIntegration(true),
		azure.WithExistingStorageAccount(true),
		azure.WithStorageAccountName("test-storage-account-name"),
	).Generate()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(errors.Unwrap(err).Error(), "invalid inputs"))
	assert.Empty(t, hcl)
}

func TestGenerationActivityLogWithSubscriptionsList(t *testing.T) {
	ActivityLogWithSubscriptions, fileErr := getFileContent("test-data/activity-log-with-list-subscriptions.tf")
	assert.Nil(t, fileErr)
	testIds := []string{"test-id-1", "test-id-2", "test-id-3"}
	hcl, err := azure.NewTerraform(false, true,
		azure.WithAdIntegration(true),
		azure.WithSubscriptionIds(testIds),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithSubscriptions, hcl)
}

func TestGenerationConfigWithSubscriptionsList(t *testing.T) {
	ConfigWithSubscriptions, fileErr := getFileContent("test-data/config-log-with-list-subscriptions.tf")
	assert.Nil(t, fileErr)
	testIds := []string{"test-id-1", "test-id-2", "test-id-3"}
	hcl, err := azure.NewTerraform(true, false,
		azure.WithAdIntegration(true),
		azure.WithSubscriptionIds(testIds),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigWithSubscriptions, hcl)
}

func TestGenerationLocation(t *testing.T) {
	ActivityLogLocation, fileErr := getFileContent("test-data/activity-log-with-location.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true,
		azure.WithAdIntegration(true),
		azure.WithLocation("West US 2"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogLocation, hcl)
}
