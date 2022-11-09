package gcp_test

import (
	"fmt"
	"testing"

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
				true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName)),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogWithoutCredentialsAndProject",
			gcp.NewTerraform(false, true),
			fmt.Sprintf("%s\n%s\n%s", RequiredProviders, gcpProviderWithoutCredentialsAndProject, moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogWithLaceworkProfile",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithLaceworkProfile("test-profile"),
			),
			fmt.Sprintf("%s\n%s\n%s\n%s", RequiredProviders, ProviderWithCredentials(projectName), laceworkProvider, moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogCustomIntegrationName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithAuditLogIntegrationName("custom_integration_name")),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogCustomIntegrationName),
		},
		{
			"TestGenerationProjectLevelAuditLogLabels",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithAuditLogLabels(map[string]string{"key": "value"}),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogBucketLabels",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithBucketLabels(map[string]string{"key": "value"}),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogBucketLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogPubSubSubscriptionLabels",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithPubSubSubscriptionLabels(map[string]string{"key": "value"}),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogPubSubSubscriptionLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogPubSubTopicLabels",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithPubSubTopicLabels(map[string]string{"key": "value"}),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogPubSubTopicLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogBucketRegion",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithBucketRegion("us-west"),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogBucketRegion),
		},
		{
			"TestGenerationProjectLevelAuditLogExistingBucketName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithExistingLogBucketName("foo"),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogExistingBucketName),
		},
		{
			"TestGenerationProjectLevelAuditLogExistingLogSinkName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithExistingLogSinkName("foo"),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogExistingLogSinkName),
		},
		{
			"TestGenerationProjectLevelAuditLogEnableForceDestroyBucket",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithEnableForceDestroyBucket(),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogEnableForceDestroyBucket),
		},
		{"TestGenerationProjectLevelAuditLogEnableUBLA",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithEnableUBLA(true),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogEnableUBLA),
		},
		{
			"TestGenerationProjectLevelAuditLogDisableUBLA",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithEnableUBLA(false),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogDisableUBLA),
		},
		{

			"TestGenerationProjectLevelAuditLogBucketLifecycleRuleAge",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithLogBucketLifecycleRuleAge(420),
			),
			ReqProvider(projectName, moduleImportProjectLevelAuditLogBucketLifecycleRuleAge),
		},
		{
			"TestGenerationOrganizationLevelAuditLogWithoutConfig",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationOrganizationLevelAuditLogCustomIntegrationName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithAuditLogIntegrationName("custom_integration_name"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelAuditLogCustomIntegrationName),
		},
		{
			"TestGenerationProjectLevelAuditLogWithConfiguration",
			gcp.NewTerraform(true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfiguration, moduleImportProjectLevelAuditLogWithConfiguration),
		},
		{
			"TestGenerationOrganizationLevelAuditLogWithConfig",
			gcp.NewTerraform(true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelConfiguration, moduleImportOrganizationLevelAuditLogWithConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogCustomBucketName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithCustomBucketName("bucket"),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.0"
  custom_bucket_name = "bucket"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogCustomFilter",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithCustomFilter("custom-filter"),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source        = "lacework/audit-log/gcp"
  version       = "~> 3.0"
  custom_filter = "custom-filter"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogGoogleWorkspaceFilter",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithGoogleWorkspaceFilter(true),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogGoogleWorkspaceFilterFalse",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithGoogleWorkspaceFilter(false),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source                  = "lacework/audit-log/gcp"
  version                 = "~> 3.0"
  google_workspace_filter = false
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogK8sFilter",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithK8sFilter(true),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogK8sFilterFalse",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithK8sFilter(false),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source     = "lacework/audit-log/gcp"
  version    = "~> 3.0"
  k8s_filter = false
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogFoldersToInclude",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToInclude([]string{"abc", "abc", "def", "def"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.0"
  folders_to_include = ["abc", "def"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogFoldersToExclude",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToExclude([]string{"abc", "abc", "def", "def"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.0"
  folders_to_exclude = ["abc", "def"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogIncludeRootProjectsSolo",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source          = "lacework/audit-log/gcp"
  version         = "~> 3.0"
  org_integration = true
  organization_id = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogIncludeRootProjectsFalse",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source                = "lacework/audit-log/gcp"
  version               = "~> 3.0"
  folders_to_exclude    = ["abc"]
  include_root_projects = false
  org_integration       = true
  organization_id       = "123456789"
}
`),
		},
		{
			"TestGenerationOrganizationLevelAuditLogIncludeRootProjectsTrue",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(true),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			ReqProvider(projectName, `module "gcp_organization_level_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.0"
  folders_to_exclude = ["abc"]
  org_integration    = true
  organization_id    = "123456789"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogPrefix",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithPrefix("rar"),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
  prefix  = "rar"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogWaitTime",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithWaitTime("30s"),
			),
			ReqProvider(projectName, `module "gcp_project_audit_log" {
  source    = "lacework/audit-log/gcp"
  version   = "~> 3.0"
  wait_time = "30s"
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
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfiguration),
		},
		{
			"TestGenerationProjectLevelConfigurationExistingSA",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithExistingServiceAccount(gcp.NewExistingServiceAccountDetails("foo", "123456789")),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfigurationExistingSA),
		},
		{
			"TestGenerationProjectLevelConfigurationCustomIntegrationName",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithConfigurationIntegrationName("custom_integration_name"),
			),
			ReqProvider(projectName, moduleImportProjectLevelConfigurationCustomIntegrationName),
		},
		{
			"TestGenerationOrganizationLevelConfiguration",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId(projectName),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			ReqProvider(projectName, moduleImportOrganizationLevelConfiguration),
		},
		{
			"TestGenerationOrganizationLevelConfigurationExistingSA",
			gcp.NewTerraform(true, false,
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
			gcp.NewTerraform(true, false,
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
			gcp.NewTerraform(true, false,
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
			gcp.NewTerraform(true, false,
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
			gcp.NewTerraform(true, false,
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
			gcp.NewTerraform(true, false,
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
			gcp.NewTerraform(true, false,
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
			gcp.NewTerraform(true, false,
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
			gcp.NewTerraform(true, false,
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
	hcl, err := gcp.NewTerraform(false, true,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithOrganizationIntegration(true),
	).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: an Organization ID must be provided for an Organization Integration")
}

func TestGenerationOrganizationLevelAuditLogNoOrgIntegrationFlag(t *testing.T) {
	hcl, err := gcp.NewTerraform(false, true,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
		gcp.WithOrganizationId("123456789"),
	).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: to provide an Organization ID, Organization Integration must be true")
}

func TestGenerationNoIntegration(t *testing.T) {
	hcl, err := gcp.NewTerraform(false, false,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId(projectName),
	).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: audit log or configuration integration must be enabled")
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

var laceworkProvider = `provider "lacework" {
  profile = "test-profile"
}
`

var moduleImportProjectLevelAuditLogWithConfiguration = `module "gcp_project_audit_log" {
  source                       = "lacework/audit-log/gcp"
  version                      = "~> 3.0"
  service_account_name         = module.gcp_project_level_config.service_account_name
  service_account_private_key  = module.gcp_project_level_config.service_account_private_key
  use_existing_service_account = true
}
`

var moduleImportProjectLevelAuditLogWithoutConfiguration = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
}
`

var moduleImportProjectLevelAuditLogCustomIntegrationName = `module "gcp_project_audit_log" {
  source                    = "lacework/audit-log/gcp"
  version                   = "~> 3.0"
  lacework_integration_name = "custom_integration_name"
}
`

var moduleImportProjectLevelAuditLogLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
  labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogBucketLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
  bucket_labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogPubSubSubscriptionLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
  pubsub_subscription_labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogPubSubTopicLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
  pubsub_topic_labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogBucketRegion = `module "gcp_project_audit_log" {
  source        = "lacework/audit-log/gcp"
  version       = "~> 3.0"
  bucket_region = "us-west"
}
`

var moduleImportProjectLevelAuditLogExistingBucketName = `module "gcp_project_audit_log" {
  source               = "lacework/audit-log/gcp"
  version              = "~> 3.0"
  existing_bucket_name = "foo"
}
`

var moduleImportProjectLevelAuditLogExistingLogSinkName = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.0"
  existing_sink_name = "foo"
}
`

var moduleImportProjectLevelAuditLogEnableForceDestroyBucket = `module "gcp_project_audit_log" {
  source               = "lacework/audit-log/gcp"
  version              = "~> 3.0"
  bucket_force_destroy = true
}
`

var moduleImportProjectLevelAuditLogEnableUBLA = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
}
`

var moduleImportProjectLevelAuditLogDisableUBLA = `module "gcp_project_audit_log" {
  source      = "lacework/audit-log/gcp"
  version     = "~> 3.0"
  enable_ubla = false
}
`

var moduleImportProjectLevelAuditLogBucketLifecycleRuleAge = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 3.0"
  lifecycle_rule_age = 420
}
`

var moduleImportOrganizationLevelAuditLogWithConfiguration = `module "gcp_organization_level_audit_log" {
  source                       = "lacework/audit-log/gcp"
  version                      = "~> 3.0"
  org_integration              = true
  organization_id              = "123456789"
  service_account_name         = module.gcp_organization_level_config.service_account_name
  service_account_private_key  = module.gcp_organization_level_config.service_account_private_key
  use_existing_service_account = true
}
`

var moduleImportOrganizationLevelAuditLogWithoutConfiguration = `module "gcp_organization_level_audit_log" {
  source          = "lacework/audit-log/gcp"
  version         = "~> 3.0"
  org_integration = true
  organization_id = "123456789"
}
`
var moduleImportOrganizationLevelAuditLogCustomIntegrationName = `module "gcp_organization_level_audit_log" {
  source                    = "lacework/audit-log/gcp"
  version                   = "~> 3.0"
  lacework_integration_name = "custom_integration_name"
  org_integration           = true
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
