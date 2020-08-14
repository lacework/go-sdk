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

	"go.uber.org/zap"
)

const (
	apiIntegrations        = "external/integrations"
	apiIntegrationsByType  = "external/integrations/type/%s"
	apiIntegrationFromGUID = "external/integrations/%s"
	apiIntegrationSchema   = "external/integrations/schema/%s"
	apiTokens              = "access/tokens"

	apiVulnerabilitiesContainerScan             = "external/vulnerabilities/container/repository/images/scan"
	apiVulnerabilitiesContainerScanStatus       = "external/vulnerabilities/container/reqId/%s"
	apiVulnerabilitiesAssessmentFromImageID     = "external/vulnerabilities/container/imageId/%s"
	apiVulnerabilitiesAssessmentFromImageDigest = "external/vulnerabilities/container/imageDigest/%s"
	apiVulnContainerAssessmentsForDateRange     = "external/vulnerabilities/container/GetAssessmentsForDateRange"

	apiVulnerabilitiesScanPkgManifest             = "external/vulnerabilities/scan"
	apiVulnerabilitiesHostListCves                = "external/vulnerabilities/host"
	apiVulnerabilitiesListHostsWithCveID          = "external/vulnerabilities/host/cveId/%s"
	apiVulnerabilitiesHostAssessmentFromMachineID = "external/vulnerabilities/host/machineId/%s"

	apiComplianceAwsLatestReport        = "external/compliance/aws/GetLatestComplianceReport?AWS_ACCOUNT_ID=%s"
	apiComplianceGcpLatestReport        = "external/compliance/gcp/GetLatestComplianceReport?GCP_ORG_ID=%s&GCP_PROJ_ID=%s"
	apiComplianceGcpListProjects        = "external/compliance/gcp/ListProjectsForOrganization?GCP_ORG_ID=%s"
	apiComplianceAzureLatestReport      = "external/compliance/azure/GetLatestComplianceReport?AZURE_TENANT_ID=%s&AZURE_SUBS_ID=%s"
	apiComplianceAzureListSubscriptions = "external/compliance/azure/ListSubscriptionsForTenant?AZURE_TENANT_ID=%s"

	apiRunReportGcp   = "external/runReport/gcp/%s"
	apiRunReportAws   = "external/runReport/aws/%s"
	apiRunReportAzure = "external/runReport/azure/%s"

	apiEventsDetails   = "external/events/GetEventDetails"
	apiEventsDateRange = "external/events/GetEventsForDateRange"

	// Alpha
	apiLQLQuery = "external/lql/query"
)

// WithApiV2 configures the client to use the API version 2 (/api/v2)
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
	return fmt.Sprintf("/api/%s/%s", c.apiVersion, p)
}
