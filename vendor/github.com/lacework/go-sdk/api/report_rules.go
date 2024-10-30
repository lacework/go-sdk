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
//
//	             use the new instance to do CRUD operations
//
//	  client, err := api.NewClient("account")
//	  if err != nil {
//	    return err
//	  }
//
//	  reportRule := api.NewReportRule(
//			"Foo",
//			api.ReportRuleConfig{
//			Description: "My Report Rule"
//			Severities: api.ReportRuleSeverities{api.ReportRuleSeverityHigh,
//			EmailAlertChannels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
//			ResourceGroups: []string{"TECHALLY_111111111111AAAAAAAAAAAAAAAAAAAA"}
//			ReportNotificationTypes: api.WeeklyEventsReportRuleNotifications{TrendReport: true},
//	      },
//	    },
//	  )
//
//	  client.V2.ReportRules.Create(reportRule)
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
	AgentEvents               bool `json:"agentEvents"`
	AwsCisS3                  bool `json:"awsCisS3"`
	AwsCloudtrailEvents       bool `json:"awsCloudtrailEvents"`
	AwsComplianceEvents       bool `json:"awsComplianceEvents"`
	AwsHipaa                  bool `json:"hipaa"`
	AwsIso2700                bool `json:"iso2700"`
	AwsNist80053Rev4          bool `json:"nist800-53Rev4"`
	AwsNist800171Rev2         bool `json:"nist800-171Rev2"`
	AwsPci                    bool `json:"pci"`
	AwsSoc                    bool `json:"soc"`
	AwsSocRev2                bool `json:"awsSocRev2"`
	AzureActivityLogEvents    bool `json:"azureActivityLogEvents"`
	AzureCis                  bool `json:"azureCis"`
	AzureCis131               bool `json:"azureCis131"`
	AzureComplianceEvents     bool `json:"azureComplianceEvents"`
	AzurePci                  bool `json:"azurePci"`
	AzureSoc                  bool `json:"azureSoc"`
	GcpAuditTrailEvents       bool `json:"gcpAuditTrailEvents"`
	GcpCis                    bool `json:"gcpCis"`
	GcpComplianceEvents       bool `json:"gcpComplianceEvents"`
	GcpHipaa                  bool `json:"gcpHipaa"`
	GcpHipaaRev2              bool `json:"gcpHipaaRev2"`
	GcpIso27001               bool `json:"gcpIso27001"`
	GcpCis12                  bool `json:"gcpCis12"`
	GcpK8s                    bool `json:"gcpK8s"`
	GcpPci                    bool `json:"gcpPci"`
	GcpPciRev2                bool `json:"gcpPciRev2"`
	GcpSoc                    bool `json:"gcpSoc"`
	GcpSocRev2                bool `json:"gcpSocRev2"`
	OpenShiftCompliance       bool `json:"openShiftCompliance"`
	OpenShiftComplianceEvents bool `json:"openShiftComplianceEvents"`
	PlatformEvents            bool `json:"platformEvents"`
	TrendReport               bool `json:"trendReport"`
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
