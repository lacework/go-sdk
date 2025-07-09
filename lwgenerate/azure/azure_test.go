package azure_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/v2/lwgenerate"
	"github.com/lacework/go-sdk/v2/lwgenerate/azure"
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
	hcl, err := azure.NewTerraform(false, true, false, false, true, azure.WithSubscriptionID("test-subscription")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithoutConfig, hcl)
}

func TestGenerationActivityLogWithConfig(t *testing.T) {
	var ActivityLogWithConfig, fileErr = getFileContent("test-data/activity_log_with_config.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true, false, false, true, azure.WithSubscriptionID("test-subscription")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithConfig, hcl)

}

func TestGenerationActivityLogWithConfigAndExtraBlocks(t *testing.T) {
	var ActivityLogWithConfig, fileErr = getFileContent("test-data/activity_log_with_config_extra.tf")
	assert.Nil(t, fileErr)
	extraBlock, err := lwgenerate.HclCreateGenericBlock("variable", []string{"var_name"}, nil)
	assert.NoError(t, err)
	hcl, err := azure.NewTerraform(true, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithExtraBlocks([]*hclwrite.Block{extraBlock}),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithConfig, hcl)
}

func TestGenerationActivityLogWithConfigAndExtraAzureRMProviderBlocks(t *testing.T) {
	var ActivityLogWithConfig, fileErr = getFileContent("test-data/activity_log_with_config_provider_args.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithExtraAZRMArguments(map[string]interface{}{"foo": "bar"}),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithConfig, hcl)
}

func TestGenerationActivityLogWithConfigAndExtraAZUReadProviderBlocks(t *testing.T) {
	var ActivityLogWithConfig, fileErr = getFileContent("test-data/activity_log_with_config_azureadprovider_args.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithExtraAZReadArguments(map[string]interface{}{"foo": "bar"}),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithConfig, hcl)
}

func TestGenerationActivityLogWithConfigAndCustomBackendBlock(t *testing.T) {
	customBlock, err := lwgenerate.HclCreateGenericBlock("backend", []string{"s3"}, nil)
	assert.NoError(t, err)
	var ActivityLogWithConfig, fileErr = getFileContent("test-data/activity_log_with_config_root_blocks.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithExtraRootBlocks([]*hclwrite.Block{customBlock}),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogWithConfig, hcl)
}

func TestGenerationConfigWithoutActivityLog(t *testing.T) {
	ConfigWithoutActivityLog, fileErr := getFileContent("test-data/config_without_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false, false, false, true, azure.WithSubscriptionID("test-subscription")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigWithoutActivityLog, hcl)
}

func TestGenerationWithoutActivityLogOrConfig(t *testing.T) {
	hcl, err := azure.NewTerraform(false, false, false, false, true).Generate()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(errors.Unwrap(err).Error(), "invalid inputs"))
	assert.Empty(t, hcl)
}
func TestGenerationRenamedConfig(t *testing.T) {
	RenamedConfig, fileErr := getFileContent("test-data/renamed_config.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithConfigIntegrationName("Test Config Rename"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedConfig, hcl)
}

func TestGenerationRenamedActivityLog(t *testing.T) {
	RenamedActivityLog, fileErr := getFileContent("test-data/renamed_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithActivityLogIntegrationName("Test Activity Log Rename"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedActivityLog, hcl)
}

func TestGenerationRenamedConfigAndActivityLog(t *testing.T) {
	RenamedConfigAndActivityLog, fileErr := getFileContent("test-data/renamed_config_and_activity_log.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithActivityLogIntegrationName("Test Activity Log Rename"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RenamedConfigAndActivityLog, hcl)
}

func TestGenerationNoActiveDirectorySettings(t *testing.T) {
	hcl, err := azure.NewTerraform(true, true, false, false, false,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithActivityLogIntegrationName("Test Activity Log Rename"),
	).Generate()
	assert.True(t, strings.Contains(errors.Unwrap(err).Error(), "invalid inputs"))
	assert.Empty(t, hcl)
}

func TestGenerationCustomActiveDirectory(t *testing.T) {
	CustomADDetails, fileErr := getFileContent("test-data/customer-ad-details.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, true, false, false, false,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithConfigIntegrationName("Test Config Rename"),
		azure.WithActivityLogIntegrationName("Test Activity Log Rename"),
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
	hcl, err := azure.NewTerraform(false, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
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
	hcl, err := azure.NewTerraform(false, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithAllSubscriptions(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogAllSubs, hcl)
}

func TestGenerationConfigWithAllSubscriptions(t *testing.T) {
	ConfigAllSubs, fileErr := getFileContent("test-data/config-with-all-subscriptions.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithAllSubscriptions(true),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigAllSubs, hcl)
}

func TestGenerationConfigWithManagementGroup(t *testing.T) {
	ConfigWithMgmtGroup, fileErr := getFileContent("test-data/config-with-management-group.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(true, false, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithManagementGroup(true),
		azure.WithManagementGroupId("test-management-group-1"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigWithMgmtGroup, hcl)
}

func TestGenerationConfigWithManagementGroupError(t *testing.T) {
	hcl, err := azure.NewTerraform(true, false, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithManagementGroup(true),
	).Generate()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(errors.Unwrap(err).Error(), "invalid inputs"))
	assert.Empty(t, hcl)
}

func TestGenerationActivityLogWithSubscriptionsList(t *testing.T) {
	ActivityLogWithSubscriptions, fileErr := getFileContent("test-data/activity-log-with-list-subscriptions.tf")
	assert.Nil(t, fileErr)
	testIds := []string{"test-id-1", "test-id-2", "test-id-3"}
	hcl, err := azure.NewTerraform(false, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
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
	hcl, err := azure.NewTerraform(true, false, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithSubscriptionIds(testIds),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ConfigWithSubscriptions, hcl)
}

func TestGenerationLocation(t *testing.T) {
	ActivityLogLocation, fileErr := getFileContent("test-data/activity-log-with-location.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithStorageLocation("West US 2"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogLocation, hcl)
}

func TestGenerationWithLaceworkProvider(t *testing.T) {
	laceworkProfile, fileErr := getFileContent("test-data/activity-log-with-lacework-profile.tf")
	assert.Nil(t, fileErr)

	hcl, err := azure.NewTerraform(false, true, false, false, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithLaceworkProfile("test-profile"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, laceworkProfile, hcl)
}

func TestGenerationAzureRmProviderWithSubscriptionID(t *testing.T) {
	configWithSubscription, fileErr := getFileContent("test-data/config-with-azurerm-subscription.tf")
	assert.Nil(t, fileErr)

	hcl, err := azure.NewTerraform(true, false, false, false, true, azure.WithSubscriptionID("test-subscription")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, configWithSubscription, hcl)
}

func TestGenerationEntraIDActivityLog(t *testing.T) {
	ActivityLogEntraID, fileErr := getFileContent("test-data/entra-id-activity-log-no-custom-input.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, false, false, true, true, azure.WithSubscriptionID("test-subscription")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogEntraID, hcl)
}

func TestGenerationEntraIDActivityLogExistingActiveDirectoryApp(t *testing.T) {
	ActivityLogEntraID, fileErr := getFileContent("test-data/entra-id-activity-log-existing-ad-app.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, false, false, true, false,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithAdApplicationId("testID"),
		azure.WithAdApplicationPassword("pass"),
		azure.WithAdServicePrincipalId("principal"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogEntraID, hcl)
}

func TestGenerationEntraIDActivityLogEventHubLocationAndPartition(t *testing.T) {
	ActivityLogEntraID, fileErr := getFileContent("test-data/entra-id-activity-log-event-hub-location-and-partition.tf")
	assert.Nil(t, fileErr)
	hcl, err := azure.NewTerraform(false, false, false, true, true,
		azure.WithSubscriptionID("test-subscription"),
		azure.WithEventHubLocation("West US 2"),
		azure.WithEventHubPartitionCount(2),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, ActivityLogEntraID, hcl)
}
