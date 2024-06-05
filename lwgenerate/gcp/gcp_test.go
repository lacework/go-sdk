package gcp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/lacework/go-sdk/lwgenerate"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/lwgenerate/gcp"
)

const (
	projectName = "project1"
)

func TestGenerateGcpTfConfigurationArgs_Generate_AuditLog(t *testing.T) {
	tests := []struct {
		description string
		gcp         *gcp.GenerateGcpTfConfigurationArgs
		expects     string
	}{
		{
			"TestGenerationProjectLevelAuditLogWithoutConfig",
			gcp.NewTerraform(
				false,
				false,
				true,
				false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName)),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelPubSubAuditLogWithoutConfig",
			gcp.NewTerraform(
				false,
				false,
				true,
				true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName)),
			ReqProvider(projectName, moduleImportProjectLevelPubSubAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogWithoutCredentialsAndProject",
			gcp.NewTerraform(false, false, true, false),
			fmt.Sprintf("%s\n%s\n%s", RequiredProviders, gcpProviderWithoutCredentialsAndProject, moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogWithLaceworkProfile",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithLaceworkProfile("test-profile"),
			),
			fmt.Sprintf("%s\n%s\n%s\n%s", RequiredProviders, ProviderWithCredentials(projectName), laceworkProvider, moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogCustomIntegrationName",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithAuditLogIntegrationName("custom_integration_name")),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogCustomIntegrationName),
		},
		{
			"TestGenerationProjectLevelPubSubAuditLogCustomIntegrationName",
			gcp.NewTerraform(false, false, true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithAuditLogIntegrationName("custom_integration_name")),
			ReqProvider(projectName, moduleImportProjectLevelPubSubAuditLogCustomIntegrationName),
		},
		{
			"TestGenerationProjectLevelAuditLogLabels",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithAuditLogLabels(map[string]string{"key": "value"}),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogBucketLabels",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithBucketLabels(map[string]string{"key": "value"}),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogBucketLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogPubSubSubscriptionLabels",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithPubSubSubscriptionLabels(map[string]string{"key": "value"}),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogPubSubSubscriptionLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogPubSubTopicLabels",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithPubSubTopicLabels(map[string]string{"key": "value"}),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogPubSubTopicLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogBucketRegion",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithBucketRegion("us-west"),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogBucketRegion),
		},
		{
			"TestGenerationProjectLevelAuditLogExistingBucketName",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithExistingLogBucketName("foo"),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogExistingBucketName),
		},
		{
			"TestGenerationProjectLevelAuditLogExistingLogSinkName",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithExistingLogSinkName("foo"),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogExistingLogSinkName),
		},
		{
			"TestGenerationProjectLevelPubSubAuditLogExistingLogSinkName",
			gcp.NewTerraform(false, false, true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithExistingLogSinkName("foo"),
			),
			ReqProvider(projectName, moduleImportProjectLevelPubSubAuditLogExistingLogSinkName),
		},
		{"TestGenerationProjectLevelAuditLogEnableUBLA",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithEnableUBLA(true),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogEnableUBLA),
		},
		{
			"TestGenerationProjectLevelAuditLogDisableUBLA",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithEnableUBLA(false),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogDisableUBLA),
		},
		{

			"TestGenerationProjectLevelAuditLogBucketLifecycleRuleAge",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithLogBucketLifecycleRuleAge(420),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogBucketLifecycleRuleAge),
		},
		{
			"TestGenerationOrganizationLevelAuditLogWithoutConfig",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationOrganizationLevelPubSubAuditLogWithoutConfig",
			gcp.NewTerraform(false, false, true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelPubSubAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationOrganizationLevelAuditLogCustomIntegrationName",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithAuditLogIntegrationName("custom_integration_name"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelAuditLogCustomIntegrationName),
		},
		{
			"TestGenerationOrganizationLevelPubSubAuditLogCustomIntegrationName",
			gcp.NewTerraform(false, false, true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithAuditLogIntegrationName("custom_integration_name"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelPubSubAuditLogCustomIntegrationName),
		},
		{
			"TestGenerationProjectLevelAuditLogWithConfiguration",
			gcp.NewTerraform(false, true, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfiguration, moduleImportProjectLevelAuditLogWithConfiguration),
		},
		{
			"TestGenerationProjectLevelPubSubAuditLogWithConfiguration",
			gcp.NewTerraform(false, true, true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfiguration, moduleImportProjectLevelPubSubAuditLogWithConfiguration),
		},
		{
			"TestGenerationOrganizationLevelAuditLogWithConfig",
			gcp.NewTerraform(false, true, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelConfiguration, moduleImportOrganizationLevelAuditLogWithConfiguration),
		},
		{
			"TestGenerationOrganizationLevelPubSubAuditLogWithConfig",
			gcp.NewTerraform(false, true, true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelConfiguration, moduleImportOrganizationLevelPubSubAuditLogWithConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogCustomBucketName",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithCustomBucketName("bucket"),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.4"
  custom_bucket_name = "bucket"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogCustomFilter",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithCustomFilter("custom-filter"),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source        = "lacework/audit-log/gcp"
  version       = "~> 3.4"
  custom_filter = "custom-filter"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogGoogleWorkspaceFilter",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithGoogleWorkspaceFilter(true),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogGoogleWorkspaceFilterFalse",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithGoogleWorkspaceFilter(false),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source                  = "lacework/audit-log/gcp"
  version                 = "~> 3.4"
  google_workspace_filter = false
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogK8sFilter",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithK8sFilter(true),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogK8sFilterFalse",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithK8sFilter(false),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source     = "lacework/audit-log/gcp"
  version    = "~> 3.4"
  k8s_filter = false
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogFoldersToInclude",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToInclude([]string{"abc", "abc", "def", "def"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.4"
  folders_to_include = ["abc", "def"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogFoldersToExclude",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToExclude([]string{"abc", "abc", "def", "def"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.4"
  folders_to_exclude = ["abc", "def"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogIncludeRootProjectsSolo",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source          = "lacework/audit-log/gcp"
  version         = "~> 3.4"
  org_integration = true
  organization_id = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogIncludeRootProjectsFalse",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source                = "lacework/audit-log/gcp"
  version               = "~> 3.4"
  folders_to_exclude    = ["abc"]
  include_root_projects = false
  org_integration       = true
  organization_id       = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogIncludeRootProjectsTrue",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(true),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.4"
  folders_to_exclude = ["abc"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogPrefix",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithPrefix("rar"),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
  prefix  = "rar"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogWaitTime",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithWaitTime("30s"),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source    = "lacework/audit-log/gcp"
  version   = "~> 3.4"
  wait_time = "30s"
}
`),
		},
		{
			"TestGenerationMultipleProject",
			gcp.NewTerraform(false, false, true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithMultipleProject([]string{"project1", "project2", "project3"}),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"

  for_each = {
    project1 = "project1"
    project2 = "project2"
    project3 = "project3"
  }
  project_id = each.key
}
`),
		},
	}

	for _, tc := range tests {
		hcl, err := tc.gcp.Generate()

		if err != nil {
			t.Errorf("Test case `%s` error: %s", tc.description, err)
		}

		if tc.expects != hcl {
			t.Errorf("Test case `%s` HCL error", tc.description)
		}

		assert.Equal(t, tc.expects, hcl)
	}
}

func TestGenerateGcpTfConfigurationArgs_Generate_Configuration(t *testing.T) {
	tests := []struct {
		description string
		gcp         *gcp.GenerateGcpTfConfigurationArgs
		expects     string
	}{
		{
			"TestGenerationProjectLevelConfiguration",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfiguration),
		},
		{
			"TestGenerationProjectLevelConfigurationExistingSA",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithExistingServiceAccount(gcp.NewExistingServiceAccountDetails("foo", "123456789")),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfigurationExistingSA),
		},
		{
			"TestGenerationProjectLevelConfigurationCustomIntegrationName",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithConfigurationIntegrationName("custom_integration_name"),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfigurationCustomIntegrationName),
		},
		{
			"TestGenerationOrganizationLevelConfiguration",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelConfiguration),
		},
		{
			"TestGenerationOrganizationLevelConfigurationExistingSA",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithExistingServiceAccount(gcp.NewExistingServiceAccountDetails("foo", "123456789")),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelConfigurationExistingSA),
		},
		{
			"TestGenerationOrganizationLevelConfigurationCustomIntegrationName",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithConfigurationIntegrationName("custom_integration_name"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelConfigurationCustomIntegrationName),
		},
		{
			"TestGenerationOrganizationLevelConfigurationFoldersToInclude",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToInclude([]string{"abc", "abc", "def", "def"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_config" {
  source             = "lacework/config/gcp"
  version            = "~> 2.3"
  folders_to_include = ["abc", "def"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelConfigurationFoldersToExclude",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToExclude([]string{"abc", "def"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_config" {
  source             = "lacework/config/gcp"
  version            = "~> 2.3"
  folders_to_exclude = ["abc", "def"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelConfigurationIncludeRootProjectsSolo",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
			),
			ReqProvider(projectName, `module "gcp_organization_level_config" {
  source          = "lacework/config/gcp"
  version         = "~> 2.3"
  org_integration = true
  organization_id = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelConfigurationIncludeRootProjectsFalse",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_config" {
  source                = "lacework/config/gcp"
  version               = "~> 2.3"
  folders_to_exclude    = ["abc"]
  include_root_projects = false
  org_integration       = true
  organization_id       = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelConfigurationIncludeRootProjectsTrue",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(true),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_config" {
  source             = "lacework/config/gcp"
  version            = "~> 2.3"
  folders_to_exclude = ["abc"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationProjectLevelConfigurationPrefix",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithPrefix("rar"),
			),
			ReqProvider(projectName, `module "gcp_project_level_config" {
  source  = "lacework/config/gcp"
  version = "~> 2.3"
  prefix  = "rar"
}
`),
		},
		{
			"TestGenerationProjectConfigurationLogWaitTime",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithWaitTime("30s"),
			),
			ReqProvider(projectName, `module "gcp_project_level_config" {
  source    = "lacework/config/gcp"
  version   = "~> 2.3"
  wait_time = "30s"
}
`),
		},
		{
			"TestGenerationMultipleProject",
			gcp.NewTerraform(false, true, false, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithMultipleProject([]string{"project1", "project2", "project3"}),
			),
			ReqProvider(projectName, `module "gcp_project_level_config" {
  source  = "lacework/config/gcp"
  version = "~> 2.3"

  for_each = {
    project1 = "project1"
    project2 = "project2"
    project3 = "project3"
  }
  project_id = each.key
}
`),
		},
	}

	for _, tc := range tests {
		hcl, err := tc.gcp.Generate()

		if err != nil {
			t.Errorf("Test case `%s` error: %s", tc.description, err)
		}

		if tc.expects != hcl {
			t.Errorf("Test case `%s` HCL error", tc.description)
		}
		assert.Equal(t, tc.expects, hcl)
	}
}

func TestGenerateGcpTfConfigurationArgs_Generate_Agentless(t *testing.T) {
	tests := []struct {
		description string
		gcp         *gcp.GenerateGcpTfConfigurationArgs
		expects     string
	}{
		{
			"TestGenerationProjectLevelAgentless",
			gcp.NewTerraform(true, false, false, false,
				gcp.WithProjectId(projectName),
				gcp.WithRegions([]string{"us-east1"}),
			),
			fmt.Sprintf("%s\n%s", RequiredProviders, moduleImportProjectLevelAgentless),
		},
		{
			"TestGenerationProjectLevelAgentlessWithProjectFilterList",
			gcp.NewTerraform(true, false, false, false,
				gcp.WithProjectId(projectName),
				gcp.WithRegions([]string{"us-east1"}),
				gcp.WithProjectFilterList([]string{"p1", "p2"}),
			),
			fmt.Sprintf("%s\n%s", RequiredProviders, moduleImportProjectLevelAgentlessWithProjectFilterList),
		},
		{
			"TestGenerationOrgLevelAgentless",
			gcp.NewTerraform(true, false, false, false,
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithRegions([]string{"us-east1"}),
			),
			fmt.Sprintf("%s\n%s", RequiredProviders, moduleImportOrgLevelAgentless),
		},
	}

	for _, tc := range tests {
		hcl, err := tc.gcp.Generate()

		if err != nil {
			t.Errorf("Test case `%s` error: %s", tc.description, err)
		}

		if tc.expects != hcl {
			t.Errorf("Test case `%s` HCL error", tc.description)
		}
		assert.Equal(t, tc.expects, hcl)
	}
}

func TestGenerationOrganizationLevelAuditLogNoOrgId(t *testing.T) {
	hcl, err := gcp.NewTerraform(false, false, true, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithOrganizationIntegration(true),
	).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: an Organization ID must be provided for an Organization Integration")
}

func TestGenerationOrganizationLevelAuditLogNoOrgIntegrationFlag(t *testing.T) {
	hcl, err := gcp.NewTerraform(false, false, true, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithOrganizationId("123456789"),
	).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: to provide an Organization ID, Organization Integration must be true")
}

func TestGenerationNoIntegration(t *testing.T) {
	hcl, err := gcp.NewTerraform(false, false, false, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
	).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: agentless, audit log or configuration integration must be enabled")
}

var RequiredProviders = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
}
`

var requiredProvidersWithCustomBlock = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 1.0"
    }
  }
  backend "s3" {
  }
}
`

func TestGenerationConfigWithExtraBlocks(t *testing.T) {
	extraBlock, err := lwgenerate.HclCreateGenericBlock("variable", []string{"var_name"}, nil)
	assert.NoError(t, err)

	hcl, err := gcp.NewTerraform(false, false, true, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithExtraBlocks([]*hclwrite.Block{extraBlock}),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RequiredProviders+"\n"+gcpProvider+"\n"+moduleImportProjectLevelAuditLogWithoutConfiguration+"\n"+testVariable, hcl)
}

func TestGenerationConfigWithCustomBackendBlock(t *testing.T) {
	customBlock, err := lwgenerate.HclCreateGenericBlock("backend", []string{"s3"}, nil)
	assert.NoError(t, err)
	hcl, err := gcp.NewTerraform(false, false, true, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithExtraRootBlocks([]*hclwrite.Block{customBlock}),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, requiredProvidersWithCustomBlock+"\n"+gcpProvider+"\n"+moduleImportProjectLevelAuditLogWithoutConfiguration, hcl)
}

func TestGenerationConfigWithCustomProviderAttributes(t *testing.T) {
	hcl, err := gcp.NewTerraform(false, false, true, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithRegions([]string{"us-east1"}),
		gcp.WithExtraProviderArguments(map[string]interface{}{"foo": "bar"}),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RequiredProviders+"\n"+gcpProviderWithExtraArguments+"\n"+moduleImportProjectLevelAuditLogWithoutConfiguration, hcl)
}

func TestGenerationConfigWithOutputs(t *testing.T) {
	hcl, err := gcp.NewTerraform(
		false, false, true, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithCustomOutputs([]lwgenerate.HclOutput{
			*lwgenerate.NewOutput("test", []string{"module", "gcp_config", "lacework_integration_guid"}, "test description"),
		})).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RequiredProviders+"\n"+gcpProvider+"\n"+moduleImportProjectLevelAuditLogWithoutConfiguration+"\n"+customOutput, hcl)
}

func TestGenerationConfigWithDefaultProviderLabels(t *testing.T) {
	hcl, err := gcp.NewTerraform(
		false, false, true, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithRegions([]string{"us-east1"}),
		gcp.WithProviderDefaultLabels(map[string]interface{}{"LABEL_TEST": "foo", "LABEL_TEST1": "bar"})).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, RequiredProviders+"\n"+gcpProviderWithDefaultLabels+"\n"+moduleImportProjectLevelAuditLogWithoutConfiguration, hcl)

}

func ProviderWithCredentials(projectName string) string {
	return fmt.Sprintf(`provider "google" {
  credentials = "/path/to/credentials"
  project     = "%s"
}
`, projectName)
}

// ReqProvider Helper for combining string expected values
func ReqProvider(projectName string, extraInputs ...string) string {
	base := RequiredProviders + "\n" + ProviderWithCredentials(projectName)
	countInputs := len(extraInputs)
	for i, e := range extraInputs {
		if i < countInputs {
			base = base + "\n" + e
		}

		if i >= countInputs {
			base = base + e
		}
	}
	return base
}

var gcpProviderWithoutCredentialsAndProject = `provider "google" {
}
`

var gcpProvider = `provider "google" {
  credentials = "/path/to/credentials"
  project     = "project1"
}
`

var gcpProviderWithExtraArguments = `provider "google" {
  alias       = "us-east1"
  credentials = "/path/to/credentials"
  foo         = "bar"
  project     = "project1"
  region      = "us-east1"
}
`

var gcpProviderWithDefaultLabels = `provider "google" {
  alias       = "us-east1"
  credentials = "/path/to/credentials"
  default_labels = {
    LABEL_TEST  = "foo"
    LABEL_TEST1 = "bar"
  }
  project = "project1"
  region  = "us-east1"
}
`

var testVariable = `variable "var_name" {
}
`

var customOutput = `output "test" {
  description = "test description"
  value       = module.gcp_config.lacework_integration_guid
}
`

var laceworkProvider = `provider "lacework" {
  profile = "test-profile"
}
`

var moduleImportProjectLevelAuditLogWithConfiguration = `module "gcp_project_audit_log" {
  source                       = "lacework/audit-log/gcp"
  version                      = "~> 3.4"
  service_account_name         = module.gcp_project_level_config.service_account_name
  service_account_private_key  = module.gcp_project_level_config.service_account_private_key
  use_existing_service_account = true
}
`
var moduleImportProjectLevelPubSubAuditLogWithConfiguration = `module "gcp_project_audit_log" {
  source                       = "lacework/pub-sub-audit-log/gcp"
  version                      = "~> 0.2"
  service_account_name         = module.gcp_project_level_config.service_account_name
  service_account_private_key  = module.gcp_project_level_config.service_account_private_key
  use_existing_service_account = true
}
`

var moduleImportProjectLevelAuditLogWithoutConfiguration = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
}
`

var moduleImportProjectLevelAuditLogCustomIntegrationName = `module "gcp_project_audit_log" {
  source                    = "lacework/audit-log/gcp"
  version                   = "~> 3.4"
  lacework_integration_name = "custom_integration_name"
}
`
var moduleImportProjectLevelPubSubAuditLogWithoutConfiguration = `module "gcp_project_audit_log" {
  source  = "lacework/pub-sub-audit-log/gcp"
  version = "~> 0.2"
}
`

var moduleImportProjectLevelPubSubAuditLogCustomIntegrationName = `module "gcp_project_audit_log" {
  source                    = "lacework/pub-sub-audit-log/gcp"
  version                   = "~> 0.2"
  lacework_integration_name = "custom_integration_name"
}
`

var moduleImportProjectLevelAuditLogLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
  labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogBucketLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
  bucket_labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogPubSubSubscriptionLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
  pubsub_subscription_labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogPubSubTopicLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
  pubsub_topic_labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogBucketRegion = `module "gcp_project_audit_log" {
  source        = "lacework/audit-log/gcp"
  version       = "~> 3.4"
  bucket_region = "us-west"
}
`

var moduleImportProjectLevelAuditLogExistingBucketName = `module "gcp_project_audit_log" {
  source               = "lacework/audit-log/gcp"
  version              = "~> 3.4"
  existing_bucket_name = "foo"
}
`

var moduleImportProjectLevelAuditLogExistingLogSinkName = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.4"
  existing_sink_name = "foo"
}
`
var moduleImportProjectLevelPubSubAuditLogExistingLogSinkName = `module "gcp_project_audit_log" {
  source             = "lacework/pub-sub-audit-log/gcp"
  version            = "~> 0.2"
  existing_sink_name = "foo"
}
`

var moduleImportProjectLevelAuditLogEnableUBLA = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.4"
}
`

var moduleImportProjectLevelAuditLogDisableUBLA = `module "gcp_project_audit_log" {
  source      = "lacework/audit-log/gcp"
  version     = "~> 3.4"
  enable_ubla = false
}
`

var moduleImportProjectLevelAuditLogBucketLifecycleRuleAge = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.4"
  lifecycle_rule_age = 420
}
`

var moduleImportOrganizationLevelAuditLogWithConfiguration = `module "gcp_organization_level_audit_log" {
  source                       = "lacework/audit-log/gcp"
  version                      = "~> 3.4"
  org_integration              = true
  organization_id              = "123456789"
  service_account_name         = module.gcp_organization_level_config.service_account_name
  service_account_private_key  = module.gcp_organization_level_config.service_account_private_key
  use_existing_service_account = true
}
`

var moduleImportOrganizationLevelPubSubAuditLogWithConfiguration = `module "gcp_organization_level_audit_log" {
  source                       = "lacework/pub-sub-audit-log/gcp"
  version                      = "~> 0.2"
  integration_type             = "ORGANIZATION"
  organization_id              = "123456789"
  service_account_name         = module.gcp_organization_level_config.service_account_name
  service_account_private_key  = module.gcp_organization_level_config.service_account_private_key
  use_existing_service_account = true
}
`

var moduleImportOrganizationLevelAuditLogWithoutConfiguration = `module "gcp_organization_level_audit_log" {
  source          = "lacework/audit-log/gcp"
  version         = "~> 3.4"
  org_integration = true
  organization_id = "123456789"
}
`

var moduleImportOrganizationLevelPubSubAuditLogWithoutConfiguration = `module "gcp_organization_level_audit_log" {
  source           = "lacework/pub-sub-audit-log/gcp"
  version          = "~> 0.2"
  integration_type = "ORGANIZATION"
  organization_id  = "123456789"
}
`
var moduleImportOrganizationLevelAuditLogCustomIntegrationName = `module "gcp_organization_level_audit_log" {
  source                    = "lacework/audit-log/gcp"
  version                   = "~> 3.4"
  lacework_integration_name = "custom_integration_name"
  org_integration           = true
  organization_id           = "123456789"
}
`
var moduleImportOrganizationLevelPubSubAuditLogCustomIntegrationName = `module "gcp_organization_level_audit_log" {
  source                    = "lacework/pub-sub-audit-log/gcp"
  version                   = "~> 0.2"
  integration_type          = "ORGANIZATION"
  lacework_integration_name = "custom_integration_name"
  organization_id           = "123456789"
}
`

var moduleImportProjectLevelConfiguration = `module "gcp_project_level_config" {
  source  = "lacework/config/gcp"
  version = "~> 2.3"
}
`

var moduleImportProjectLevelConfigurationExistingSA = `module "gcp_project_level_config" {
  source                       = "lacework/config/gcp"
  version                      = "~> 2.3"
  service_account_name         = "foo"
  service_account_private_key  = "123456789"
  use_existing_service_account = true
}
`

var moduleImportProjectLevelConfigurationCustomIntegrationName = `module "gcp_project_level_config" {
  source                    = "lacework/config/gcp"
  version                   = "~> 2.3"
  lacework_integration_name = "custom_integration_name"
}
`

var moduleImportOrganizationLevelConfiguration = `module "gcp_organization_level_config" {
  source          = "lacework/config/gcp"
  version         = "~> 2.3"
  org_integration = true
  organization_id = "123456789"
}
`

var moduleImportOrganizationLevelConfigurationExistingSA = `module "gcp_organization_level_config" {
  source                       = "lacework/config/gcp"
  version                      = "~> 2.3"
  org_integration              = true
  organization_id              = "123456789"
  service_account_name         = "foo"
  service_account_private_key  = "123456789"
  use_existing_service_account = true
}
`

var moduleImportOrganizationLevelConfigurationCustomIntegrationName = `module "gcp_organization_level_config" {
  source                    = "lacework/config/gcp"
  version                   = "~> 2.3"
  lacework_integration_name = "custom_integration_name"
  org_integration           = true
  organization_id           = "123456789"
}
`

var moduleImportProjectLevelAgentless = `provider "google" {
  alias   = "us-east1"
  project = "project1"
  region  = "us-east1"
}

module "lacework_gcp_agentless_scanning_global" {
  source   = "lacework/agentless-scanning/gcp"
  version  = "~> 0.1"
  global   = true
  regional = true

  providers = {
    google = google.us-east1
  }
}
`

var moduleImportProjectLevelAgentlessWithProjectFilterList = `provider "google" {
  alias   = "us-east1"
  project = "project1"
  region  = "us-east1"
}

module "lacework_gcp_agentless_scanning_global" {
  source              = "lacework/agentless-scanning/gcp"
  version             = "~> 0.1"
  global              = true
  project_filter_list = ["p1", "p2"]
  regional            = true

  providers = {
    google = google.us-east1
  }
}
`

var moduleImportOrgLevelAgentless = `provider "google" {
  alias   = "us-east1"
  project = "project1"
  region  = "us-east1"
}

module "lacework_gcp_agentless_scanning_global" {
  source           = "lacework/agentless-scanning/gcp"
  version          = "~> 0.1"
  global           = true
  integration_type = "ORGANIZATION"
  organization_id  = "123456789"
  regional         = true

  providers = {
    google = google.us-east1
  }
}
`
