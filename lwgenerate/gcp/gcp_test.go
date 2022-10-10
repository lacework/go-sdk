package gcp_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/lwgenerate/gcp"
)

// Helper for combining string expected values
func reqProvider(extraInputs ...string) string {
	base := requiredProviders + "\n" + gcpProviderWithCredentials
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

func TestGenerateGcpTfConfigurationArgs_Generate(t *testing.T) {
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
				gcp.WithProjectId("project1")),
			reqProvider(moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogWithoutCredentialsAndProject",
			gcp.NewTerraform(false, true),
			fmt.Sprintf("%s\n%s\n%s", requiredProviders, gcpProviderWithoutCredentialsAndProject, moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogWithLaceworkProfile",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithLaceworkProfile("test-profile"),
			),
			fmt.Sprintf("%s\n%s\n%s\n%s", requiredProviders, gcpProviderWithCredentials, laceworkProvider, moduleImportProjectLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationProjectLevelAuditLogCustomIntegrationName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithAuditLogIntegrationName("custom_integration_name")),
			reqProvider(moduleImportProjectLevelAuditLogCustomIntegrationName),
		},
		{
			"TestGenerationProjectLevelAuditLogLabels",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithAuditLogLabels(map[string]string{"key": "value"}),
			),
			reqProvider(moduleImportProjectLevelAuditLogLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogBucketLabels",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithBucketLabels(map[string]string{"key": "value"}),
			),
			reqProvider(moduleImportProjectLevelAuditLogBucketLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogPubSubSubscriptionLabels",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithPubSubSubscriptionLabels(map[string]string{"key": "value"}),
			),
			reqProvider(moduleImportProjectLevelAuditLogPubSubSubscriptionLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogPubSubTopicLabels",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithPubSubTopicLabels(map[string]string{"key": "value"}),
			),
			reqProvider(moduleImportProjectLevelAuditLogPubSubTopicLabels),
		},
		{
			"TestGenerationProjectLevelAuditLogBucketRegion",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithBucketRegion("us-west"),
			),
			reqProvider(moduleImportProjectLevelAuditLogBucketRegion),
		},
		{
			"TestGenerationProjectLevelAuditLogExistingBucketName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithExistingLogBucketName("foo"),
			),
			reqProvider(moduleImportProjectLevelAuditLogExistingBucketName),
		},
		{
			"TestGenerationProjectLevelAuditLogExistingLogSinkName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithExistingLogSinkName("foo"),
			),
			reqProvider(moduleImportProjectLevelAuditLogExistingLogSinkName),
		},
		{
			"TestGenerationProjectLevelAuditLogEnableForceDestroyBucket",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithEnableForceDestroyBucket(),
			),
			reqProvider(moduleImportProjectLevelAuditLogEnableForceDestroyBucket),
		},
		{"TestGenerationProjectLevelAuditLogEnableUBLA",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithEnableUBLA(true),
			),
			reqProvider(moduleImportProjectLevelAuditLogEnableUBLA),
		},
		{
			"TestGenerationProjectLevelAuditLogDisableUBLA",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithEnableUBLA(false),
			),
			reqProvider(moduleImportProjectLevelAuditLogDisableUBLA),
		},
		{

			"TestGenerationProjectLevelAuditLogBucketLifecycleRuleAge",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithLogBucketLifecycleRuleAge(420),
			),
			reqProvider(moduleImportProjectLevelAuditLogBucketLifecycleRuleAge),
		},
		{
			"TestGenerationOrganizationLevelAuditLogWithoutConfig",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			reqProvider(moduleImportOrganizationLevelAuditLogWithoutConfiguration),
		},
		{
			"TestGenerationOrganizationLevelAuditLogCustomIntegrationName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithAuditLogIntegrationName("custom_integration_name"),
			),
			reqProvider(moduleImportOrganizationLevelAuditLogCustomIntegrationName),
		},
		{
			"TestGenerationProjectLevelAuditLogWithConfiguration",
			gcp.NewTerraform(true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
			),
			reqProvider(moduleImportProjectLevelConfiguration, moduleImportProjectLevelAuditLogWithConfiguration),
		},
		{
			"TestGenerationOrganizationLevelAuditLogWithConfig",
			gcp.NewTerraform(true, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			reqProvider(moduleImportOrganizationLevelConfiguration, moduleImportOrganizationLevelAuditLogWithConfiguration),
		},
		{
			"TestGenerationProjectLevelConfiguration",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
			),
			reqProvider(moduleImportProjectLevelConfiguration),
		},
		{
			"TestGenerationProjectLevelConfigurationExistingSA",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithExistingServiceAccount(gcp.NewExistingServiceAccountDetails("foo", "123456789")),
			),
			reqProvider(moduleImportProjectLevelConfigurationExistingSA),
		},
		{
			"TestGenerationProjectLevelConfigurationCustomIntegrationName",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithConfigurationIntegrationName("custom_integration_name"),
			),
			reqProvider(moduleImportProjectLevelConfigurationCustomIntegrationName),
		},
		{
			"TestGenerationOrganizationLevelConfiguration",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
			),
			reqProvider(moduleImportOrganizationLevelConfiguration),
		},
		{
			"TestGenerationOrganizationLevelConfigurationExistingSA",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithExistingServiceAccount(gcp.NewExistingServiceAccountDetails("foo", "123456789")),
			),
			reqProvider(moduleImportOrganizationLevelConfigurationExistingSA),
		},
		{
			"TestGenerationOrganizationLevelConfigurationCustomIntegrationName",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithConfigurationIntegrationName("custom_integration_name"),
			),
			reqProvider(moduleImportOrganizationLevelConfigurationCustomIntegrationName),
		},
		{
			"TestGenerationProjectLevelAuditLogCustomBucketName",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithCustomBucketName("bucket"),
			),
			reqProvider(`module "gcp_project_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithCustomFilter("custom-filter"),
			),
			reqProvider(`module "gcp_project_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithGoogleWorkspaceFilter(true),
			),
			reqProvider(`module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogGoogleWorkspaceFilterFalse",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithGoogleWorkspaceFilter(false),
			),
			reqProvider(`module "gcp_project_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithK8sFilter(true),
			),
			reqProvider(`module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 3.0"
}
`),
		},
		{
			"TestGenerationProjectLevelAuditLogK8sFilterFalse",
			gcp.NewTerraform(false, true,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithK8sFilter(false),
			),
			reqProvider(`module "gcp_project_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToInclude([]string{"abc", "abc", "def", "def"}),
			),
			reqProvider(`module "gcp_organization_level_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToExclude([]string{"abc", "abc", "def", "def"}),
			),
			reqProvider(`module "gcp_organization_level_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
			),
			reqProvider(`module "gcp_organization_level_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			reqProvider(`module "gcp_organization_level_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(true),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			reqProvider(`module "gcp_organization_level_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithPrefix("rar"),
			),
			reqProvider(`module "gcp_project_audit_log" {
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
				gcp.WithProjectId("project1"),
				gcp.WithWaitTime("30s"),
			),
			reqProvider(`module "gcp_project_audit_log" {
  source    = "lacework/audit-log/gcp"
  version   = "~> 3.0"
  wait_time = "30s"
}
`),
		},
		{
			"TestGenerationOrganizationLevelConfigurationFoldersToInclude",
			gcp.NewTerraform(true, false,
				gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToInclude([]string{"abc", "abc", "def", "def"}),
			),
			reqProvider(`module "gcp_organization_level_config" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithFoldersToExclude([]string{"abc", "def"}),
			),
			reqProvider(`module "gcp_organization_level_config" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
			),
			reqProvider(`module "gcp_organization_level_config" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(false),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			reqProvider(`module "gcp_organization_level_config" {
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
				gcp.WithProjectId("project1"),
				gcp.WithOrganizationIntegration(true),
				gcp.WithOrganizationId("123456789"),
				gcp.WithIncludeRootProjects(true),
				gcp.WithFoldersToExclude([]string{"abc"}),
			),
			reqProvider(`module "gcp_organization_level_config" {
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
				gcp.WithProjectId("project1"),
				gcp.WithPrefix("rar"),
			),
			reqProvider(`module "gcp_project_level_config" {
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
				gcp.WithProjectId("project1"),
				gcp.WithWaitTime("30s"),
			),
			reqProvider(`module "gcp_project_level_config" {
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
	}
}

func TestGenerationOrganizationLevelAuditLogNoOrgId(t *testing.T) {
	hcl, err := gcp.NewTerraform(false, true,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId("project1"),
		gcp.WithOrganizationIntegration(true),
	).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: An Organization ID must be provided for an Organization Integration")
}

func TestGenerationOrganizationLevelAuditLogNoOrgIntegrationFlag(t *testing.T) {
	hcl, err := gcp.NewTerraform(false, true,
		gcp.WithGcpServiceAccountCredentials("/path/to/credentials"),
		gcp.WithProjectId("project1"),
		gcp.WithOrganizationId("123456789"),
	).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: To provide an Organization ID, Organization Integration must be true")
}

var requiredProviders = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 0.16"
    }
  }
}
`

var gcpProviderWithCredentials = `provider "google" {
  credentials = "/path/to/credentials"
  project     = "project1"
}
`

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
