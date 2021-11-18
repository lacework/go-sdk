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
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

type reportRuleNotification interface {
	allNotifications() ReportRuleNotificationTypes
	ToMap() map[string]bool
}

type ReportRuleNotifications []reportRuleNotification

// Enable all Gcp report rules
var AllGcpReportRuleNotifications = new(GcpReportRuleNotifications).allNotifications()

// Enable all Aws report rules
var AllAwsReportRuleNotifications = new(AwsReportRuleNotifications).allNotifications()

// Enable all Azure report rules
var AllAzureReportRuleNotifications = new(AzureReportRuleNotifications).allNotifications()

// Enable all Daily report rules
var AllDailyReportRuleNotifications = new(DailyEventsReportRuleNotifications).allNotifications()

// Enable all Weekly report rules
var AllWeeklyReportRuleNotifications = new(WeeklyEventsReportRuleNotifications).allNotifications()

// Enable all report rules
var AllReportRuleNotifications = new(ReportRuleNotificationTypes).allNotifications()

func TransformReportRuleNotification(notificationsMap map[string]bool, notificationType reportRuleNotification) error {
	jsonMap, err := json.Marshal(notificationsMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonMap, &notificationType)
	if err != nil {
		return err
	}

	return nil
}

func NewReportRuleNotificationTypes(types []reportRuleNotification) (ReportRuleNotificationTypes, error) {
	notificationsTypes := ReportRuleNotificationTypes{}
	notificationsMap := make(map[string]bool)

	for _, notificationType := range types {
		m := structs.Map(notificationType)
		for k, v := range m {
			if _, ok := notificationsMap[k]; ok {
				return ReportRuleNotificationTypes{}, errors.New(fmt.Sprintf("notification types contains a duplicate type: %s", k))
			}
			if v.(bool) {
				notificationsMap[k] = true
			} else {
				notificationsMap[k] = false
			}
		}

		jsonMap, err := json.Marshal(notificationsMap)
		if err != nil {
			return ReportRuleNotificationTypes{}, errors.New("unable to set report rule notification types")
		}

		err = json.Unmarshal(jsonMap, &notificationsTypes)
		if err != nil {
			return ReportRuleNotificationTypes{}, errors.New("unable to set report rule notification types")
		}
	}
	return notificationsTypes, nil
}

func reportRuleNotificationToMap(notificationType reportRuleNotification) map[string]bool {
	notificationsMap := make(map[string]bool)
	m := structs.Map(notificationType)
	for k, v := range m {
		if v.(bool) {
			notificationsMap[k] = true
		} else {
			notificationsMap[k] = false
		}
	}
	return notificationsMap
}

type GcpReportRuleNotifications struct {
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
}

func (gcp GcpReportRuleNotifications) allNotifications() ReportRuleNotificationTypes {
	return ReportRuleNotificationTypes{
		GcpCis:       true,
		GcpHipaa:     true,
		GcpHipaaRev2: true,
		GcpIso27001:  true,
		GcpCis12:     true,
		GcpK8s:       true,
		GcpPci:       true,
		GcpPciRev2:   true,
		GcpSoc:       true,
		GcpSocRev2:   true,
	}
}

func (gcp GcpReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(gcp)
}

type AwsReportRuleNotifications struct {
	AwsCisS3          bool `json:"awsCisS3"`
	AwsHipaa          bool `json:"hipaa"`
	AwsIso2700        bool `json:"iso2700"`
	AwsNist80053Rev4  bool `json:"nist800-53Rev4"`
	AwsNist800171Rev2 bool `json:"nist800-171Rev2"`
	AwsPci            bool `json:"pci"`
	AwsSoc            bool `json:"soc"`
	AwsSocRev2        bool `json:"awsSocRev2"`
}

func (aws AwsReportRuleNotifications) allNotifications() ReportRuleNotificationTypes {
	return ReportRuleNotificationTypes{
		AwsCisS3:          true,
		AwsHipaa:          true,
		AwsIso2700:        true,
		AwsNist80053Rev4:  true,
		AwsNist800171Rev2: true,
		AwsPci:            true,
		AwsSoc:            true,
		AwsSocRev2:        true,
	}
}

func (aws AwsReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(aws)
}

type AzureReportRuleNotifications struct {
	AzureCis    bool `json:"azureCis"`
	AzureCis131 bool `json:"azureCis131"`
	AzurePci    bool `json:"azurePci"`
	AzureSoc    bool `json:"azureSoc"`
}

func (az AzureReportRuleNotifications) allNotifications() ReportRuleNotificationTypes {
	return ReportRuleNotificationTypes{
		AzureActivityLogEvents: true,
		AzureCis:               true,
		AzureCis131:            true,
		AzurePci:               true,
		AzureSoc:               true,
	}
}

func (az AzureReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(az)
}

type DailyEventsReportRuleNotifications struct {
	AgentEvents               bool `json:"agentEvents"`
	OpenShiftCompliance       bool `json:"openShiftCompliance"`
	OpenShiftComplianceEvents bool `json:"openShiftComplianceEvents"`
	PlatformEvents            bool `json:"platformEvents"`
	AwsCloudtrailEvents       bool `json:"awsCloudtrailEvents"`
	AwsComplianceEvents       bool `json:"awsComplianceEvents"`
	AzureComplianceEvents     bool `json:"azureComplianceEvents"`
	AzureActivityLogEvents    bool `json:"azureActivityLogEvents"`
	GcpAuditTrailEvents       bool `json:"gcpAuditTrailEvents"`
	GcpComplianceEvents       bool `json:"gcpComplianceEvents"`
}

func (daily DailyEventsReportRuleNotifications) allNotifications() ReportRuleNotificationTypes {
	return ReportRuleNotificationTypes{
		AgentEvents:               true,
		OpenShiftCompliance:       true,
		OpenShiftComplianceEvents: true,
		PlatformEvents:            true,
		AwsCloudtrailEvents:       true,
		AwsComplianceEvents:       true,
		AzureComplianceEvents:     true,
		AzureActivityLogEvents:    true,
		GcpAuditTrailEvents:       true,
		GcpComplianceEvents:       true,
	}
}

func (daily DailyEventsReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(daily)
}

type WeeklyEventsReportRuleNotifications struct {
	TrendReport bool `json:"trendReport"`
}

func (weekly WeeklyEventsReportRuleNotifications) allNotifications() ReportRuleNotificationTypes {
	return ReportRuleNotificationTypes{
		TrendReport: true,
	}
}

func (weekly WeeklyEventsReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(weekly)
}

func (all ReportRuleNotificationTypes) allNotifications() ReportRuleNotificationTypes {
	return ReportRuleNotificationTypes{
		AgentEvents:               true,
		AwsCisS3:                  true,
		AwsCloudtrailEvents:       true,
		AwsComplianceEvents:       true,
		AwsHipaa:                  true,
		AwsIso2700:                true,
		AwsNist80053Rev4:          true,
		AwsNist800171Rev2:         true,
		AwsPci:                    true,
		AwsSoc:                    true,
		AwsSocRev2:                true,
		AzureActivityLogEvents:    true,
		AzureCis:                  true,
		AzureCis131:               true,
		AzureComplianceEvents:     true,
		AzurePci:                  true,
		AzureSoc:                  true,
		GcpAuditTrailEvents:       true,
		GcpCis:                    true,
		GcpComplianceEvents:       true,
		GcpHipaa:                  true,
		GcpHipaaRev2:              true,
		GcpIso27001:               true,
		GcpCis12:                  true,
		GcpK8s:                    true,
		GcpPci:                    true,
		GcpPciRev2:                true,
		GcpSoc:                    true,
		GcpSocRev2:                true,
		OpenShiftCompliance:       true,
		OpenShiftComplianceEvents: true,
		PlatformEvents:            true,
		TrendReport:               true,
	}
}

func (all ReportRuleNotificationTypes) ToMap() map[string]bool {
	return reportRuleNotificationToMap(all)
}
