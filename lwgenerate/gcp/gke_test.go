package gcp_test

import (
	"testing"

	"github.com/lacework/go-sdk/lwgenerate/gcp"

	"github.com/stretchr/testify/assert"
)

const (
	gkeProjName = "project1"
)

func TestGenerateGKETfConfigurationArgs_Generate(t *testing.T) {
	tests := []struct {
		description string
		gke         *gcp.GenerateGkeTfConfigurationArgs
		expects     string
	}{
		{
			"TestGeneration organization level GKE Audit Log",
			gcp.NewGkeTerraform(gcp.WithGkeServiceAccountCredentials("/path/to/credentials"),
				gcp.WithGkeProjectId(gkeProjName),
				gcp.WithGkeOrganizationIntegration(true),
				gcp.WithGkeOrganizationId("123456789"),
			),
			ReqProvider(gkeProjName, `module "gcp_organization_level_gke_audit_log" {
  source           = "lacework/gke-audit-log/gcp"
  version          = "~> 0.3"
  integration_type = "ORGANIZATION"
  organization_id  = "123456789"
}
`),
		},
		{
			"TestGeneration project level GKE Audit Log existing_sink_name",
			gcp.NewGkeTerraform(gcp.WithGkeServiceAccountCredentials("/path/to/credentials"),
				gcp.WithGkeProjectId("project1"),
				gcp.WithGkeExistingSinkName("sink"),
			),
			ReqProvider(gkeProjName, `module "gcp_project_level_gke_audit_log" {
  source             = "lacework/gke-audit-log/gcp"
  version            = "~> 0.3"
  existing_sink_name = "sink"
  integration_type   = "PROJECT"
}
`),
		},
		{
			"TestGeneration project level GKE Audit Log lacework_integration_name",
			gcp.NewGkeTerraform(
				gcp.WithGkeServiceAccountCredentials("/path/to/credentials"),
				gcp.WithGkeProjectId("project1"),
				gcp.WithGkeIntegrationName("custom-integration"),
			),
			ReqProvider(gkeProjName, `module "gcp_project_level_gke_audit_log" {
  source                    = "lacework/gke-audit-log/gcp"
  version                   = "~> 0.3"
  integration_type          = "PROJECT"
  lacework_integration_name = "custom-integration"
}
`),
		},
		{
			"TestGeneration project level GKE Audit Log prefix",
			gcp.NewGkeTerraform(
				gcp.WithGkeServiceAccountCredentials("/path/to/credentials"),
				gcp.WithGkeProjectId("project1"),
				gcp.WithGkePrefix("custom-prefix"),
			),
			ReqProvider(gkeProjName, `module "gcp_project_level_gke_audit_log" {
  source           = "lacework/gke-audit-log/gcp"
  version          = "~> 0.3"
  integration_type = "PROJECT"
  prefix           = "custom-prefix"
}
`),
		},
		{
			"TestGeneration project level GKE Audit Log service_account_name",
			gcp.NewGkeTerraform(
				gcp.WithGkeServiceAccountCredentials("/path/to/credentials"),
				gcp.WithGkeProjectId("project1"),
				gcp.WithGkeExistingServiceAccount(gcp.NewServiceAccount("foo", "123456789")),
			),
			ReqProvider(gkeProjName, `module "gcp_project_level_gke_audit_log" {
  source                       = "lacework/gke-audit-log/gcp"
  version                      = "~> 0.3"
  integration_type             = "PROJECT"
  service_account_name         = "foo"
  service_account_private_key  = "123456789"
  use_existing_service_account = true
}
`),
		},
		{
			"TestGeneration project level GKE Audit Log WaitTime",
			gcp.NewGkeTerraform(
				gcp.WithGkeServiceAccountCredentials("/path/to/credentials"),
				gcp.WithGkeProjectId("project1"),
				gcp.WithGkeWaitTime("30s"),
			),
			ReqProvider(gkeProjName, `module "gcp_project_level_gke_audit_log" {
  source           = "lacework/gke-audit-log/gcp"
  version          = "~> 0.3"
  integration_type = "PROJECT"
  wait_time        = "30s"
}
`),
		},
	}

	for _, tc := range tests {
		hcl, err := tc.gke.Generate()

		if err != nil {
			t.Errorf("Test case `%s` error: %s", tc.description, err)
		}

		if tc.expects != hcl {
			t.Errorf("Test case `%s` HCL error", tc.description)
		}

		assert.Equal(t, tc.expects, hcl)
	}
}
