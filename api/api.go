//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

const (
	// Common endpoints
	//
	// There will be times where we will have common endpoints between
	// different versions of our APIs, by default such endpoints will
	// be driven by the global client.apiVersion setting, when we are
	// ready to switch over/upgrade we can do so with the option
	// WithApiV2() at the time that the caller initializes the Go Client
	//
	// Example:
	// ```go
	//   api.NewClient("my-account", api.WithApiV2())
	// ```
	apiTokens = "access/tokens" // Auth

	// API v1 Endpoints
	//
	// These endpoints only exist in APIv1 and therefore we prefix them with 'v1/'
	apiIntegrations        = "v1/external/integrations"
	apiIntegrationsByType  = "v1/external/integrations/type/%s"
	apiIntegrationFromGUID = "v1/external/integrations/%s"
	apiIntegrationSchema   = "v1/external/integrations/schema/%s"

	apiAgentTokens      = "v1/external/tokens"
	apiAgentTokenFromID = "v1/external/tokens/%s"

	apiVulnerabilitiesContainerScan             = "v1/external/vulnerabilities/container/repository/images/scan"
	apiVulnerabilitiesContainerScanStatus       = "v1/external/vulnerabilities/container/reqId/%s"
	apiVulnerabilitiesAssessmentFromImageID     = "v1/external/vulnerabilities/container/imageId/%s"
	apiVulnerabilitiesAssessmentFromImageDigest = "v1/external/vulnerabilities/container/imageDigest/%s"
	apiVulnContainerAssessmentsForDateRange     = "v1/external/vulnerabilities/container/GetAssessmentsForDateRange"

	apiVulnerabilitiesScanPkgManifest             = "v1/external/vulnerabilities/scan"
	apiVulnerabilitiesHostListCves                = "v1/external/vulnerabilities/host"
	apiVulnerabilitiesListHostsWithCveID          = "v1/external/vulnerabilities/host/cveId/%s"
	apiVulnerabilitiesHostAssessmentFromMachineID = "v1/external/vulnerabilities/host/machineId/%s"

	apiComplianceAwsLatestReport        = "v1/external/compliance/aws/GetLatestComplianceReport?AWS_ACCOUNT_ID=%s"
	apiComplianceGcpLatestReport        = "v1/external/compliance/gcp/GetLatestComplianceReport?GCP_ORG_ID=%s&GCP_PROJ_ID=%s"
	apiComplianceGcpListProjects        = "v1/external/compliance/gcp/ListProjectsForOrganization?GCP_ORG_ID=%s"
	apiComplianceAzureLatestReport      = "v1/external/compliance/azure/GetLatestComplianceReport?AZURE_TENANT_ID=%s&AZURE_SUBS_ID=%s"
	apiComplianceAzureListSubscriptions = "v1/external/compliance/azure/ListSubscriptionsForTenant?AZURE_TENANT_ID=%s"

	apiRunReportIntegration = "v1/external/runReport/integration/%s"
	apiRunReportGcp         = "v1/external/runReport/gcp/%s"
	apiRunReportAws         = "v1/external/runReport/aws/%s"
	apiRunReportAzure       = "v1/external/runReport/azure/%s"

	apiEventsDetails   = "v1/external/events/GetEventDetails"
	apiEventsDateRange = "v1/external/events/GetEventsForDateRange"

	apiAccountOrganizationInfo = "v1/external/account/organizationInfo"

	// API v2 Endpoints
	//
	// These endpoints only exist in APIv2 and therefore we prefix them with 'v2/'
	apiV2UserProfile = "v2/UserProfile"

	apiV2ContainerRegistries       = "v2/ContainerRegistries"
	apiV2ContainerRegistryFromGUID = "v2/ContainerRegistries/%s"

	apiV2AlertChannels        = "v2/AlertChannels"
	apiV2AlertChannelFromGUID = "v2/AlertChannels/%s"
	apiV2AlertChannelTest     = "v2/AlertChannels/%s/test"

	apiV2AlertRules        = "v2/AlertRules"
	apiV2AlertRuleFromGUID = "v2/AlertRules/%s"

	apiV2CloudAccounts        = "v2/CloudAccounts"
	apiV2CloudAccountFromGUID = "v2/CloudAccounts/%s"

	apiV2AgentAccessTokens       = "v2/AgentAccessTokens"
	apiV2AgentAccessTokensSearch = "v2/AgentAccessTokens/search"
	apiV2AgentAccessTokenFromID  = "v2/AgentAccessTokens/%s"

	apiV2Policies        = "v2/Policies"
	apiV2Queries         = "v2/Queries"
	apiV2QueriesExecute  = "v2/Queries/execute"
	apiV2QueriesValidate = "v2/Queries/validate"

	apiV2ReportRules        = "v2/ReportRules"
	apiV2ReportRuleFromGUID = "v2/ReportRules/%s"

	apiV2ResourceGroups         = "v2/ResourceGroups"
	apiV2ResourceGroupsFromGUID = "v2/ResourceGroups/%s"

	apiV2Datasources = "v2/Datasources"

	apiV2TeamMembers         = "v2/TeamMembers"
	apiV2TeamMembersFromGUID = "v2/TeamMembers/%s"
	apiV2TeamMembersSearch   = "v2/TeamMembers/search"
)

// WithApiV2 configures the client to use the API version 2 (/api/v2)
// for common API endpoints
func WithApiV2() Option {
	return clientFunc(func(c *Client) error {
		c.log.Debug("setting up client", zap.String("api_version", "v2"))
		c.apiVersion = "v2"
		return nil
	})
}

// ApiVersion returns the API client version
func (c *Client) ApiVersion() string {
	return c.apiVersion
}

// apiPath builds a path by using the current API version
func (c *Client) apiPath(p string) string {
	if strings.HasPrefix(p, "v1") || strings.HasPrefix(p, "v2") {
		return fmt.Sprintf("/api/%s", p)
	}

	return fmt.Sprintf("/api/%s/%s", c.apiVersion, p)
}
