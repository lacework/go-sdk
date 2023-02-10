//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

	"github.com/pkg/errors"
)

// ReportRulesService is the service that interacts with
// the ReportRules schema from the Lacework APIv2 Server
type ReportRulesService struct {
	client *Client
}

type reportRuleSeverity int

type ReportRuleSeverities []reportRuleSeverity

const ReportRuleEventType = "Report"

func (sevs ReportRuleSeverities) toInt() []int {
	var res []int
	for _, i := range sevs {
		res = append(res, int(i))
	}
	return res
}

func (sevs ReportRuleSeverities) ToStringSlice() []string {
	var res []string
	for _, i := range sevs {
		switch i {
		case ReportRuleSeverityCritical:
			res = append(res, "Critical")
		case ReportRuleSeverityHigh:
			res = append(res, "High")
		case ReportRuleSeverityMedium:
			res = append(res, "Medium")
		case ReportRuleSeverityLow:
			res = append(res, "Low")
		case ReportRuleSeverityInfo:
			res = append(res, "Info")
		default:
			continue
		}
	}
	return res
}

func NewReportRuleSeverities(sevSlice []string) ReportRuleSeverities {
	var res ReportRuleSeverities
	for _, i := range sevSlice {
		sev := convertReportRuleSeverity(i)
		if sev != ReportRuleSeverityUnknown {
			res = append(res, sev)
		}
	}
	return res
}

func NewReportRuleSeveritiesFromIntSlice(sevSlice []int) ReportRuleSeverities {
	var res ReportRuleSeverities
	for _, i := range sevSlice {
		sev := convertReportRuleSeverityInt(i)
		if sev != ReportRuleSeverityUnknown {
			res = append(res, sev)
		}
	}
	return res
}

func convertReportRuleSeverity(sev string) reportRuleSeverity {
	switch strings.ToLower(sev) {
	case "critical":
		return ReportRuleSeverityCritical
	case "high":
		return ReportRuleSeverityHigh
	case "medium":
		return ReportRuleSeverityMedium
	case "low":
		return ReportRuleSeverityLow
	case "info":
		return ReportRuleSeverityInfo
	default:
		return ReportRuleSeverityUnknown
	}
}

func convertReportRuleSeverityInt(sev int) reportRuleSeverity {
	switch sev {
	case 1:
		return ReportRuleSeverityCritical
	case 2:
		return ReportRuleSeverityHigh
	case 3:
		return ReportRuleSeverityMedium
	case 4:
		return ReportRuleSeverityLow
	case 5:
		return ReportRuleSeverityInfo
	default:
		return ReportRuleSeverityUnknown
	}
}

const (
	ReportRuleSeverityCritical reportRuleSeverity = 1
	ReportRuleSeverityHigh     reportRuleSeverity = 2
	ReportRuleSeverityMedium   reportRuleSeverity = 3
	ReportRuleSeverityLow      reportRuleSeverity = 4
	ReportRuleSeverityInfo     reportRuleSeverity = 5
	ReportRuleSeverityUnknown  reportRuleSeverity = 0
)

// NewReportRule returns an instance of the ReportRule struct
//
// Basic usage: Initialize a new ReportRule struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   reportRule := api.NewReportRule(
//		"Foo",
//		api.ReportRuleConfig{
//		Description: "My Report Rule"
//		Severities: api.ReportRuleSeverities{api.ReportRuleSeverityHigh,
//		EmailAlertChannels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
//		ResourceGroups: []string{"TECHALLY_111111111111AAAAAAAAAAAAAAAAAAAA"}
//		ReportNotificationTypes: api.WeeklyEventsReportRuleNotifications{TrendReport: true},
//       },
//     },
//   )
//
//   client.V2.ReportRules.Create(reportRule)
//
func NewReportRule(name string, rule ReportRuleConfig) (ReportRule, error) {
	notifications, err := NewReportRuleNotificationTypes(rule.NotificationTypes)
	if err != nil {
		return ReportRule{}, err
	}

	return ReportRule{
		EmailAlertChannels: rule.EmailAlertChannels,
		Type:               ReportRuleEventType,
		Filter: ReportRuleFilter{
			Name:           name,
			Enabled:        1,
			Description:    rule.Description,
			Severity:       rule.Severities.toInt(),
			ResourceGroups: rule.ResourceGroups,
		},
		ReportNotificationTypes: notifications,
	}, nil
}

func (rule ReportRuleFilter) Status() string {
	if rule.Enabled == 1 {
		return "Enabled"
	}
	return "Disabled"
}

// List returns a list of Report Rules
func (svc *ReportRulesService) List() (response ReportRulesResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2ReportRules, nil, &response)
	return
}

// Create creates a single Report Rule
func (svc *ReportRulesService) Create(rule ReportRule) (
	response ReportRuleResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2ReportRules, rule, &response)
	return
}

// Delete deletes a Report Rule that matches the provided guid
func (svc *ReportRulesService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2ReportRuleFromGUID, guid),
		nil,
		nil,
	)
}

// Update updates a single Report Rule of the provided guid.
func (svc *ReportRulesService) Update(data ReportRule) (
	response ReportRuleResponse,
	err error,
) {
	if data.Guid == "" {
		err = errors.New("specify a Guid")
		return
	}
	apiPath := fmt.Sprintf(apiV2ReportRuleFromGUID, data.Guid)
	err = svc.client.RequestEncoderDecoder("PATCH", apiPath, data, &response)
	return
}

// Get returns a raw response of the Report Rule with the matching guid.
func (svc *ReportRulesService) Get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify a Guid")
	}
	apiPath := fmt.Sprintf(apiV2ReportRuleFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, &response)
}

type ReportRuleConfig struct {
	EmailAlertChannels []string
	Description        string
	Severities         ReportRuleSeverities
	NotificationTypes  []reportRuleNotification
	ResourceGroups     []string
}

type ReportRule struct {
	Guid                    string                      `json:"mcGuid,omitempty"`
	Type                    string                      `json:"type"`
	EmailAlertChannels      []string                    `json:"intgGuidList"`
	Filter                  ReportRuleFilter            `json:"filters"`
	ReportNotificationTypes ReportRuleNotificationTypes `json:"reportNotificationTypes"`
}

type ReportRuleNotificationTypes struct {
	// DEPRECATED (GROW-1409) --
	AwsCisS3         bool `json:"awsCisS3"`
	AwsIso2700       bool `json:"iso2700"`
	AwsNist80053Rev4 bool `json:"nist800-53Rev4"`
	AwsPci           bool `json:"pci"`
	AwsSoc           bool `json:"soc"`
	AwsSocRev2       bool `json:"awsSocRev2"`

	AzureCis    bool `json:"azureCis"`
	AzureCis131 bool `json:"azureCis131"`
	AzurePci    bool `json:"azurePci"`
	AzureSoc    bool `json:"azureSoc"`

	GcpCis       bool `json:"gcpCis"`
	GcpHipaa     bool `json:"gcpHipaa"`
	GcpHipaaRev2 bool `json:"gcpHipaaRev2"`
	GcpIso27001  bool `json:"gcpIso27001"`
	GcpCis12     bool `json:"gcpCis12"`
	GcpK8s       bool `json:"gcpK8s"`
	GcpPci       bool `json:"gcpPci"`
	GcpPciRev2   bool `json:"gcpPciRev2"`
	GcpSoc       bool `json:"gcpSoc"`
	GcpSocRev2   bool `json:"gcpSocRev2"`
	// -- DEPRECATED

	// AWS
	AwsSoc2                 bool `json:"awsSoc2"`
	AwsCmmc102              bool `json:"awsCmmc1.02"`
	AwsNistCsf              bool `json:"awsNistCsf"`
	AwsHipaa                bool `json:"awsHipaa"`
	AwsCsaCcm405            bool `json:"awsCsaCcm405"`
	AwsCis14                bool `json:"awsCis14"`
	AwsIso270012013         bool `json:"awsIso27001:2013"`
	AwsCyberEssentials22    bool `json:"awsCyberEssentials22"`
	AwsCis14IsoIec270022022 bool `json:"awsCis14IsoIec270022022"`
	AwsPciDss321            bool `json:"awsPciDss3.2.1"`
	AwsNist800171Rev2       bool `json:"awsNist800-171Rev2"`
	LwAwsSecAdd10           bool `json:"lwAwsSecAdd10"`
	AwsNist80053Rev5        bool `json:"awsNist800-53Rev5"`

	// AZURE
	AzureIso270012013Cis15   bool `json:"azureIso27001:2013Cis15"`
	AzureSoc2Cis15           bool `json:"azureSoc2Cis15"`
	AzureCis15               bool `json:"azureCis15"`
	AzureNistCsfCis15        bool `json:"azureNistCsfCis15"`
	AzurePciDss321Cis15      bool `json:"azurePciDss321Cis15"`
	AzureHipaaCis15          bool `json:"azureHipaaCis15"`
	AzureNist80053Rev5Cis15  bool `json:"azureNist800-53Rev5Cis15"`
	AzureNist800171Rev2Cis15 bool `json:"azureNist800-171Rev2Cis15"`

	// GCP
	GcpHipaa2013            bool `json:"gcpHipaa2013"`
	GcpSoc2                 bool `json:"gcpSoc2"`
	GcpPciDss321            bool `json:"gcpPciDss321"`
	GcpCis130NistCsf        bool `json:"gcpCis130NistCsf"`
	GcpCis130Nist80053Rev5  bool `json:"gcpCis130Nist80053Rev5"`
	GcpCmmc102              bool `json:"gcpCmmc102"`
	GcpCis130Nist800171Rev2 bool `json:"gcpCis130Nist800171Rev2"`
	GcpIso270012013         bool `json:"gcpIso270012013"`
	GcpCis13                bool `json:"gcpCis13"`

	// Daily
	OpenShiftComplianceEvents bool `json:"openShiftComplianceEvents"`
	K8SAuditLogEvents         bool `json:"k8sAuditLogEvents"`
	GcpComplianceEvents       bool `json:"gcpComplianceEvents"`
	AgentEvents               bool `json:"agentEvents"`
	AwsComplianceEvents       bool `json:"awsComplianceEvents"`
	AzureComplianceEvents     bool `json:"azureComplianceEvents"`
	AzureActivityLogEvents    bool `json:"azureActivityLogEvents"`
	GcpAuditTrailEvents       bool `json:"gcpAuditTrailEvents"`
	PlatformEvents            bool `json:"platformEvents"`
	AwsCloudtrailEvents       bool `json:"awsCloudtrailEvents"`
	IncidentEvents            bool `json:"incidentEvents"`
	OpenShiftCompliance       bool `json:"openShiftCompliance"`

	// Weekly
	TrendReport bool `json:"trendReport"`

	// Vulnerability
	ContainerVulnerabilityReport bool `json:"containerVulnerabilityReport"`
	HostVulnerabilityReport      bool `json:"hostVulnerabilityReport"`
}

type ReportRuleFilter struct {
	Name                 string   `json:"name"`
	Enabled              int      `json:"enabled"`
	Description          string   `json:"description,omitempty"`
	Severity             []int    `json:"severity"`
	ResourceGroups       []string `json:"resourceGroups,omitempty"`
	CreatedOrUpdatedTime string   `json:"createdOrUpdatedTime,omitempty"`
	CreatedOrUpdatedBy   string   `json:"createdOrUpdatedBy,omitempty"`
}

type ReportRuleResponse struct {
	Data ReportRule `json:"data"`
}

type ReportRulesResponse struct {
	Data []ReportRule `json:"data"`
}
