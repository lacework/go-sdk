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
)

const (
	// API v2 Endpoints
	apiTokens = "v2/access/tokens" // Auth

	apiV2UserProfile = "v2/UserProfile"

	apiV2ContainerRegistries       = "v2/ContainerRegistries"
	apiV2ContainerRegistryFromGUID = "v2/ContainerRegistries/%s"

	apiV2AlertChannels        = "v2/AlertChannels"
	apiV2AlertChannelFromGUID = "v2/AlertChannels/%s"
	apiV2AlertChannelTest     = "v2/AlertChannels/%s/test"

	apiV2AlertProfiles        = "v2/AlertProfiles"
	apiV2AlertProfileFromGUID = "v2/AlertProfiles/%s"

	apiV2AlertTemplates         = "v2/AlertProfiles/%s/AlertTemplates"
	apiV2AlertTemplatesFromGUID = "v2/AlertProfiles/%s/AlertTemplates/%s"

	apiV2AlertRules        = "v2/AlertRules"
	apiV2AlertRuleFromGUID = "v2/AlertRules/%s"

	apiV2CloudAccounts          = "v2/CloudAccounts"
	apiV2CloudAccountsWithParam = "v2/CloudAccounts/%s"

	apiV2AgentAccessTokens       = "v2/AgentAccessTokens"
	apiV2AgentAccessTokensSearch = "v2/AgentAccessTokens/search"
	apiV2AgentAccessTokenFromID  = "v2/AgentAccessTokens/%s"

	apiV2AgentInfoSearch = "v2/AgentInfo/search"

	apiV2PolicyExceptions                = "v2/Exceptions?policyId=%s"
	apiV2PolicyExceptionsFromExceptionID = "v2/Exceptions/%s?policyId=%s"

	apiV2InventorySearch  = "v2/Inventory/search"
	apiV2InventoryScanCsp = "v2/Inventory/scan?csp=%s"

	apiV2ComplianceEvaluationsSearch = "v2/Configs/ComplianceEvaluations/search"

	apiV2Components         = "v2/Components?os=%s&arch=%s"
	apiV2ComponentsVersions = "v2/Components/%d?os=%s&arch=%s"
	apiV2ComponentsFetch    = "v2/Components/Artifact/%d?os=%s&arch=%s&version=%s"

	apiV2ComponentDataRequest  = "v2/ComponentData/requestUpload"
	apiV2ComponentDataComplete = "v2/ComponentData/completeUpload"

	apiV2ConfigsAzure              = "v2/Configs/AzureSubscriptions"
	apiV2ConfigsAzureSubscriptions = "v2/Configs/AzureSubscriptions?tenantId=%s"
	apiV2ConfigsGcp                = "v2/Configs/GcpProjects"
	apiV2ConfigsGcpProjects        = "v2/Configs/GcpProjects?orgId=%s"

	apiV2FeatureFlags = "v2/FeatureFlags"

	apiV2Policies        = "v2/Policies"
	apiV2Queries         = "v2/Queries"
	apiV2QueriesExecute  = "v2/Queries/execute"
	apiV2QueriesValidate = "v2/Queries/validate"

	apiV2Reports               = "v2/Reports?primaryQueryId=%s&format=%s&%s=%s"
	apiV2ReportsSecondaryQuery = "v2/Reports?primaryQueryId=%s&secondaryQueryId=%s&format=%s&%s=%s"

	apiV2ReportDefinitions         = "v2/ReportDefinitions"
	apiV2ReportDefinitionsFromGUID = "v2/ReportDefinitions/%s"
	apiV2ReportDefinitionsRevert   = "v2/ReportDefinitions/%s?revertTo=%d"
	apiV2ReportDefinitionsVersions = "v2/ReportDefinitions/%s?allVersions=true"

	apiV2ReportDistributions         = "v2/ReportDistributions"
	apiV2ReportDistributionsFromGUID = "v2/ReportDistributions/%s"

	apiV2ReportRules        = "v2/ReportRules"
	apiV2ReportRuleFromGUID = "v2/ReportRules/%s"

	apiV2ResourceGroups         = "v2/ResourceGroups"
	apiV2ResourceGroupsFromGUID = "v2/ResourceGroups/%s"

	apiV2Datasources = "v2/Datasources"

	apiV2DataExportRules         = "v2/DataExportRules"
	apiV2DataExportRulesFromGUID = "v2/DataExportRules/%s"
	apiV2DataExportRulesSearch   = "v2/DataExportRules/search"

	apiV2VulnerabilitiesContainersSearch     = "v2/Vulnerabilities/Containers/search"
	apiV2VulnerabilitiesContainersScan       = "v2/Vulnerabilities/Containers/scan"
	apiV2VulnerabilitiesContainersScanStatus = "v2/Vulnerabilities/Containers/scan/%s"
	apiV2VulnerabilitiesHostsSearch          = "v2/Vulnerabilities/Hosts/search"
	apiV2VulnerabilitiesSoftwarePackagesScan = "v2/Vulnerabilities/SoftwarePackages/scan"

	apiV2VulnerabilityExceptions        = "v2/VulnerabilityExceptions"
	apiV2VulnerabilityExceptionFromGUID = "v2/VulnerabilityExceptions/%s"

	apiV2EntitiesSearch = "v2/Entities/%s/search"

	apiV2Alerts        = "v2/Alerts"
	apiV2AlertsByTime  = "v2/Alerts?startTime=%s&endTime=%s"
	apiV2AlertsSearch  = "v2/Alerts/search"
	apiV2AlertsDetails = "v2/Alerts/%d?scope=%s"
	apiV2AlertsComment = "v2/Alerts/%d/comment"
	apiV2AlertsClose   = "v2/Alerts/%d/close"

	apiV2OrganizationInfo = "v2/OrganizationInfo"

	apiSuppressions = "v2/suppressions/%s/allExceptions"

	apiRecommendations = "v2/recommendations/%s"

	apiV2MigrateGcpAtSes = "v2/migrateGcpAtSes"
)

// WithApiV2 configures the client to use the API version 2 (/api/v2)
// for common API endpoints
//
// (no-op) DEPRECATED
func WithApiV2() Option {
	return clientFunc(func(c *Client) error {
		c.log.Warn("WithApiV2() has been deprecated, all clients now default to APIv2")
		return nil
	})
}

// ApiVersion returns the API client version
func (c *Client) ApiVersion() string {
	return c.apiVersion
}

// apiPath builds a path by using the current API version
func (c *Client) apiPath(p string) string {
	if strings.HasPrefix(p, "/api/v") {
		return p
	}

	if strings.HasPrefix(p, "v1") || strings.HasPrefix(p, "v2") {
		return fmt.Sprintf("/api/%s", p)
	}

	return fmt.Sprintf("/api/%s/%s", c.apiVersion, p)
}
