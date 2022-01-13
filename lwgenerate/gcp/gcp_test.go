package gcp

import (
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

var laceworkProvider = `provider "lacework" {
  profile = "test-profile"
}
`

var moduleImportProjectLevelAuditLogWithConfig = `module "gcp_project_audit_log" {
  source                       = "lacework/audit-log/gcp"
  version                      = "~> 2.0"
  lifecycle_rule_age           = 0
  service_account_name         = module.gcp_project_level_config.service_account_name
  service_account_private_key  = module.gcp_project_level_config.service_account_private_key
  use_existing_service_account = true
}
`

var moduleImportProjectLevelAuditLogWithoutConfig = `module "gcp_project_audit_log" {
  source             = "lacework/audit-log/gcp"
  version            = "~> 2.0"
  lifecycle_rule_age = 0
}
`

var moduleImportOrganizationLevelAuditLogWithConfig = `module "gcp_organization_level_audit_log" {
  source                       = "lacework/audit-log/gcp"
  version                      = "~> 2.0"
  lifecycle_rule_age           = 0
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
  lifecycle_rule_age = 0
  org_integration    = true
  organization_id    = "123456789"
}
`

var moduleImportProjectLevelConfig = `module "gcp_project_level_config" {
  source  = "lacework/config/gcp"
  version = "~> 1.0"
}
`

var moduleImportOrganizationLevelConfig = `module "gcp_organization_level_config" {
  source          = "lacework/config/gcp"
  version         = "~> 1.0"
  org_integration = true
  organization_id = "123456789"
}
`
