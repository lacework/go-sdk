package gcp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGenerationProjectLevelAuditLogWithoutConfig(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogWithoutConfig), hcl)
}

func TestGenerationProjectLevelAuditLogWithoutCredentialsAndProject(t *testing.T) {
	hcl, err := NewTerraform(false, true).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(
		t,
		fmt.Sprintf("%s\n%s\n%s", requiredProviders, gcpProviderWithoutCredentialsAndProject, moduleImportProjectLevelAuditLogWithoutConfig),
		hcl,
	)
}

func TestGenerationProjectLevelAuditLogWithLaceworkProfile(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithLaceworkProfile("test-profile")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(
		t,
		fmt.Sprintf("%s\n%s\n%s\n%s", requiredProviders, gcpProviderWithCredentials, laceworkProvider, moduleImportProjectLevelAuditLogWithoutConfig),
		hcl,
	)
}

func TestGenerationProjectLevelAuditLogCustomIntegrationName(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithAuditLogIntegrationName("custom_integration_name")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogCustomIntegrationName), hcl)
}

func TestGenerationProjectLevelAuditLogLabels(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithAuditLogLabels(map[string]string{"key": "value"})).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogLabels), hcl)
}

func TestGenerationProjectLevelAuditLogBucketLabels(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithBucketLabels(map[string]string{"key": "value"})).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogBucketLabels), hcl)
}

func TestGenerationProjectLevelAuditLogPubSubSubscriptionLabels(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithPubSubSubscriptionLabels(map[string]string{"key": "value"})).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogPubSubSubscriptionLabels), hcl)
}

func TestGenerationProjectLevelAuditLogPubSubTopicLabels(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithPubSubTopicLabels(map[string]string{"key": "value"}),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogPubSubTopicLabels), hcl)
}

func TestGenerationProjectLevelAuditLogBucketRegion(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithBucketRegion("us-west"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogBucketRegion), hcl)
}

func TestGenerationProjectLevelAuditLogBucketLocation(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithBucketLocation("us"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogBucketLocation), hcl)
}

func TestGenerationProjectLevelAuditLogBucketName(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithBucketName("foo"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogBucketName), hcl)
}

func TestGenerationProjectLevelAuditLogExistingBucketName(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithExistingLogBucketName("foo"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogExistingBucketName), hcl)
}

func TestGenerationProjectLevelAuditLogExistingLogSinkName(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithExistingLogSinkName("foo"),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogExistingLogSinkName), hcl)
}

func TestGenerationProjectLevelAuditLogEnableForceDestroyBucket(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithEnableForceDestroyBucket(),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogEnableForceDestroyBucket), hcl)
}

func TestGenerationProjectLevelAuditLogEnableUBLA(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithEnableUBLA(),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogEnableUBLA), hcl)
}

func TestGenerationProjectLevelAuditLogBucketLifecycleRuleAge(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithLogBucketLifecycleRuleAge(420),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogBucketLifecycleRuleAge), hcl)
}

func TestGenerationProjectLevelAuditLogBucketRetentionDays(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithLogBucketRetentionDays(420),
	).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelAuditLogBucketRetentionDays), hcl)
}

func TestGenerationOrganizationLevelAuditLogWithoutConfig(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithOrganizationIntegration(true),
		WithOrganizationId("123456789")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportOrganizationLevelAuditLogWithoutConfig), hcl)
}

func TestGenerationOrganizationLevelAuditLogCustomIntegrationName(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithOrganizationIntegration(true),
		WithOrganizationId("123456789"),
		WithAuditLogIntegrationName("custom_integration_name")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportOrganizationLevelAuditLogCustomIntegrationName), hcl)
}

func TestGenerationOrganizationLevelAuditLogNoOrgId(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithOrganizationIntegration(true)).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: An Organization ID must be provided for an Organization Integration")
}

func TestGenerationOrganizationLevelAuditLogNoOrgIntegrationFlag(t *testing.T) {
	hcl, err := NewTerraform(false, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithOrganizationId("123456789")).Generate()
	assert.Empty(t, hcl)
	assert.EqualError(t, err, "invalid inputs: To provide an Organization ID, Organization Integration must be true")
}

func TestGenerationProjectLevelAuditLogWithConfig(t *testing.T) {
	hcl, err := NewTerraform(true, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelConfig, moduleImportProjectLevelAuditLogWithConfig), hcl)
}

func TestGenerationOrganizationLevelAuditLogWithConfig(t *testing.T) {
	hcl, err := NewTerraform(true, true,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithOrganizationIntegration(true),
		WithOrganizationId("123456789")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportOrganizationLevelConfig, moduleImportOrganizationLevelAuditLogWithConfig), hcl)
}

func TestGenerationProjectLevelConfig(t *testing.T) {
	hcl, err := NewTerraform(true, false,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelConfig), hcl)
}

func TestGenerationProjectLevelConfigExistingSA(t *testing.T) {
	hcl, err := NewTerraform(true, false,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithExistingServiceAccount(NewExistingServiceAccountDetails("foo", "123456789"))).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelConfigExistingSA), hcl)
}

func TestGenerationProjectLevelConfigCustomIntegrationName(t *testing.T) {
	hcl, err := NewTerraform(true, false,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithConfigIntegrationName("custom_integration_name")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportProjectLevelConfigCustomIntegrationName), hcl)
}

func TestGenerationOrganizationLevelConfig(t *testing.T) {
	hcl, err := NewTerraform(true, false,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithOrganizationIntegration(true),
		WithOrganizationId("123456789")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportOrganizationLevelConfig), hcl)
}

func TestGenerationOrganizationLevelConfigExistingSA(t *testing.T) {
	hcl, err := NewTerraform(true, false,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithOrganizationIntegration(true),
		WithOrganizationId("123456789"),
		WithExistingServiceAccount(NewExistingServiceAccountDetails("foo", "123456789"))).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportOrganizationLevelConfigExistingSA), hcl)
}

func TestGenerationOrganizationLevelConfigCustomIntegrationName(t *testing.T) {
	hcl, err := NewTerraform(true, false,
		WithGcpServiceAccountCredentials("/path/to/credentials"),
		WithProjectId("project1"),
		WithOrganizationIntegration(true),
		WithOrganizationId("123456789"),
		WithConfigIntegrationName("custom_integration_name")).Generate()
	assert.Nil(t, err)
	assert.NotNil(t, hcl)
	assert.Equal(t, reqProvider(moduleImportOrganizationLevelConfigCustomIntegrationName), hcl)
}

var requiredProviders = `terraform {
  required_providers {
    lacework = {
      source  = "lacework/lacework"
      version = "~> 0.12.2"
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

var moduleImportProjectLevelAuditLogWithConfig = `module "gcp_project_audit_log" {
  source                       = "lacework/audit-log/gcp"
  version                      = "~> 2.0"
  lifecycle_rule_age           = -1
  service_account_name         = module.gcp_project_level_config.service_account_name
  service_account_private_key  = module.gcp_project_level_config.service_account_private_key
  use_existing_service_account = true
}
`

var moduleImportProjectLevelAuditLogWithoutConfig = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  lifecycle_rule_age = -1
}
`

var moduleImportProjectLevelAuditLogCustomIntegrationName = `module "gcp_project_audit_log" {
  source                    = "lacework/audit-log/gcp"
  version                   = "~> 2.0"
  lacework_integration_name = "custom_integration_name"
  lifecycle_rule_age        = -1
}
`

var moduleImportProjectLevelAuditLogLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 2.0"
  labels = {
    key = "value"
  }
  lifecycle_rule_age = -1
}
`

var moduleImportProjectLevelAuditLogBucketLabels = `module "gcp_project_audit_log" {
  source  = "lacework/audit-log/gcp"
  version = "~> 2.0"
  bucket_labels = {
    key = "value"
  }
  lifecycle_rule_age = -1
}
`

var moduleImportProjectLevelAuditLogPubSubSubscriptionLabels = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  lifecycle_rule_age = -1
  pubsub_subscription_labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogPubSubTopicLabels = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  lifecycle_rule_age = -1
  pubsub_topic_labels = {
    key = "value"
  }
}
`

var moduleImportProjectLevelAuditLogBucketRegion = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  bucket_region      = "us-west"
  lifecycle_rule_age = -1
}
`

var moduleImportProjectLevelAuditLogBucketLocation = `module "gcp_project_audit_log" {
  source              = "lacework/audit-log/gcp"
  version             = "~> 2.0"
  lifecycle_rule_age  = -1
  log_bucket_location = "us"
}
`

var moduleImportProjectLevelAuditLogBucketName = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  lifecycle_rule_age = -1
  log_bucket         = "foo"
}
`

var moduleImportProjectLevelAuditLogExistingBucketName = `module "gcp_project_audit_log" {
  source               = "lacework/audit-log/gcp"
  version              = "~> 2.0"
  existing_bucket_name = "foo"
}
`

var moduleImportProjectLevelAuditLogExistingLogSinkName = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  existing_sink_name = "foo"
  lifecycle_rule_age = -1
}
`

var moduleImportProjectLevelAuditLogEnableForceDestroyBucket = `module "gcp_project_audit_log" {
  source               = "lacework/audit-log/gcp"
  version              = "~> 2.0"
  bucket_force_destroy = true
  lifecycle_rule_age   = -1
}
`

var moduleImportProjectLevelAuditLogEnableUBLA = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  enable_ubla        = true
  lifecycle_rule_age = -1
}
`

var moduleImportProjectLevelAuditLogBucketLifecycleRuleAge = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  lifecycle_rule_age = 420
}
`

var moduleImportProjectLevelAuditLogBucketRetentionDays = `module "gcp_project_audit_log" {
  source                    = "lacework/audit-log/gcp"
  version                   = "~> 2.0"
  lifecycle_rule_age        = -1
  log_bucket_retention_days = 420
}
`

var moduleImportOrganizationLevelAuditLogWithConfig = `module "gcp_organization_level_audit_log" {
  source                       = "lacework/audit-log/gcp"
  version                      = "~> 2.0"
  lifecycle_rule_age           = -1
  org_integration              = true
  organization_id              = "123456789"
  service_account_name         = module.gcp_organization_level_config.service_account_name
  service_account_private_key  = module.gcp_organization_level_config.service_account_private_key
  use_existing_service_account = true
}
`

var moduleImportOrganizationLevelAuditLogWithoutConfig = `module "gcp_organization_level_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  lifecycle_rule_age = -1
  org_integration    = true
  organization_id    = "123456789"
}
`

var moduleImportOrganizationLevelAuditLogCustomIntegrationName = `module "gcp_organization_level_audit_log" {
  source                    = "lacework/audit-log/gcp"
  version                   = "~> 2.0"
  lacework_integration_name = "custom_integration_name"
  lifecycle_rule_age        = -1
  org_integration           = true
  organization_id           = "123456789"
}
`

var moduleImportProjectLevelConfig = `module "gcp_project_level_config" {
  source  = "lacework/config/gcp"
  version = "~> 1.0"
}
`

var moduleImportProjectLevelConfigExistingSA = `module "gcp_project_level_config" {
  source                       = "lacework/config/gcp"
  version                      = "~> 1.0"
  service_account_name         = "foo"
  service_account_private_key  = "123456789"
  use_existing_service_account = true
}
`

var moduleImportProjectLevelConfigCustomIntegrationName = `module "gcp_project_level_config" {
  source                    = "lacework/config/gcp"
  version                   = "~> 1.0"
  lacework_integration_name = "custom_integration_name"
}
`

var moduleImportOrganizationLevelConfig = `module "gcp_organization_level_config" {
  source          = "lacework/config/gcp"
  version         = "~> 1.0"
  org_integration = true
  organization_id = "123456789"
}
`

var moduleImportOrganizationLevelConfigExistingSA = `module "gcp_organization_level_config" {
  source                       = "lacework/config/gcp"
  version                      = "~> 1.0"
  org_integration              = true
  organization_id              = "123456789"
  service_account_name         = "foo"
  service_account_private_key  = "123456789"
  use_existing_service_account = true
}
`

var moduleImportOrganizationLevelConfigCustomIntegrationName = `module "gcp_organization_level_config" {
  source                    = "lacework/config/gcp"
  version                   = "~> 1.0"
  lacework_integration_name = "custom_integration_name"
  org_integration           = true
  organization_id           = "123456789"
}
`
