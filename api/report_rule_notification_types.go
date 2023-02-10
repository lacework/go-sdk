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
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type reportRuleNotification interface {
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

// Enable all Vulnerability report rules
var AllVulnerabilityReportRuleNotifications = new(VulnerabilityReportRuleNotifications).allNotifications()

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
			notificationsMap[k] = v.(bool)
		}
	}

	err := mapstructure.Decode(notificationsMap, &notificationsTypes)
	if err != nil {
		return ReportRuleNotificationTypes{}, errors.New("unable to set report rule notification types")
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
	// DEPRECATED (GROW-1409) --
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

	GcpHipaa2013            bool `json:"gcpHipaa2013"`
	GcpSoc2                 bool `json:"gcpSoc2"`
	GcpPciDss321            bool `json:"gcpPciDss321"`
	GcpCis130NistCsf        bool `json:"gcpCis130NistCsf"`
	GcpCis130Nist80053Rev5  bool `json:"gcpCis130Nist80053Rev5"`
	GcpCmmc102              bool `json:"gcpCmmc102"`
	GcpCis130Nist800171Rev2 bool `json:"gcpCis130Nist800171Rev2"`
	GcpIso270012013         bool `json:"gcpIso270012013"`
	GcpCis13                bool `json:"gcpCis13"`
}

func (gcp GcpReportRuleNotifications) allNotifications() GcpReportRuleNotifications {
	return GcpReportRuleNotifications{
		GcpCis:                  true,
		GcpHipaa:                true,
		GcpHipaaRev2:            true,
		GcpIso27001:             true,
		GcpCis12:                true,
		GcpK8s:                  true,
		GcpPci:                  true,
		GcpPciRev2:              true,
		GcpSoc:                  true,
		GcpSocRev2:              true,
		GcpHipaa2013:            true,
		GcpSoc2:                 true,
		GcpPciDss321:            true,
		GcpCis130NistCsf:        true,
		GcpCis130Nist80053Rev5:  true,
		GcpCmmc102:              true,
		GcpCis130Nist800171Rev2: true,
		GcpIso270012013:         true,
		GcpCis13:                true,
	}
}

func (gcp GcpReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(gcp)
}

type AwsReportRuleNotifications struct {
	// DEPRECATED (GROW-1409) --
	AwsCisS3         bool `json:"awsCisS3"`
	AwsIso2700       bool `json:"iso2700"`
	AwsNist80053Rev4 bool `json:"nist800-53Rev4"`
	AwsPci           bool `json:"pci"`
	AwsSoc           bool `json:"soc"`
	AwsSocRev2       bool `json:"awsSocRev2"`
	// -- DEPRECATED

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
}

func (aws AwsReportRuleNotifications) allNotifications() AwsReportRuleNotifications {
	return AwsReportRuleNotifications{
		AwsCisS3:                true,
		AwsIso2700:              true,
		AwsNist80053Rev4:        true,
		AwsPci:                  true,
		AwsSoc:                  true,
		AwsSocRev2:              true,
		AwsSoc2:                 true,
		AwsCmmc102:              true,
		AwsNistCsf:              true,
		AwsHipaa:                true,
		AwsCsaCcm405:            true,
		AwsCis14:                true,
		AwsIso270012013:         true,
		AwsCyberEssentials22:    true,
		AwsCis14IsoIec270022022: true,
		AwsPciDss321:            true,
		AwsNist800171Rev2:       true,
		LwAwsSecAdd10:           true,
		AwsNist80053Rev5:        true,
	}
}

func (aws AwsReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(aws)
}

type AzureReportRuleNotifications struct {
	// DEPRECATED (GROW-1409) --
	AzureCis    bool `json:"azureCis"`
	AzureCis131 bool `json:"azureCis131"`
	AzurePci    bool `json:"azurePci"`
	AzureSoc    bool `json:"azureSoc"`
	// -- DEPRECATED

	AzureIso270012013Cis15   bool `json:"azureIso27001:2013Cis15"`
	AzureSoc2Cis15           bool `json:"azureSoc2Cis15"`
	AzureCis15               bool `json:"azureCis15"`
	AzureNistCsfCis15        bool `json:"azureNistCsfCis15"`
	AzurePciDss321Cis15      bool `json:"azurePciDss321Cis15"`
	AzureHipaaCis15          bool `json:"azureHipaaCis15"`
	AzureNist80053Rev5Cis15  bool `json:"azureNist800-53Rev5Cis15"`
	AzureNist800171Rev2Cis15 bool `json:"azureNist800-171Rev2Cis15"`
}

func (az AzureReportRuleNotifications) allNotifications() AzureReportRuleNotifications {
	return AzureReportRuleNotifications{
		AzureCis:    true,
		AzureCis131: true,
		AzurePci:    true,
		AzureSoc:    true,
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

func (daily DailyEventsReportRuleNotifications) allNotifications() DailyEventsReportRuleNotifications {
	return DailyEventsReportRuleNotifications{
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

func (weekly WeeklyEventsReportRuleNotifications) allNotifications() WeeklyEventsReportRuleNotifications {
	return WeeklyEventsReportRuleNotifications{
		TrendReport: true,
	}
}

func (weekly WeeklyEventsReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(weekly)
}

type VulnerabilityReportRuleNotifications struct {
	ContainerVulnerabilityReport bool `json:"containerVulnerabilityReport"`
	HostVulnerabilityReport      bool `json:"hostVulnerabilityReport"`
}

func (vuln VulnerabilityReportRuleNotifications) allNotifications() VulnerabilityReportRuleNotifications {
	return VulnerabilityReportRuleNotifications{
		ContainerVulnerabilityReport: true,
		HostVulnerabilityReport:      true,
	}
}

func (vuln VulnerabilityReportRuleNotifications) ToMap() map[string]bool {
	return reportRuleNotificationToMap(vuln)
}

func (all ReportRuleNotificationTypes) allNotifications() ReportRuleNotificationTypes {
	return ReportRuleNotificationTypes{
		AwsCisS3:                     true,
		AwsIso2700:                   true,
		AwsNist80053Rev4:             true,
		AwsPci:                       true,
		AwsSoc:                       true,
		AwsSocRev2:                   true,
		AzureCis:                     true,
		AzureCis131:                  true,
		AzurePci:                     true,
		AzureSoc:                     true,
		GcpCis:                       true,
		GcpHipaa:                     true,
		GcpHipaaRev2:                 true,
		GcpIso27001:                  true,
		GcpCis12:                     true,
		GcpK8s:                       true,
		GcpPci:                       true,
		GcpPciRev2:                   true,
		GcpSoc:                       true,
		GcpSocRev2:                   true,
		AwsSoc2:                      true,
		AwsCmmc102:                   true,
		AwsNistCsf:                   true,
		AwsHipaa:                     true,
		AwsCsaCcm405:                 true,
		AwsCis14:                     true,
		AwsIso270012013:              true,
		AwsCyberEssentials22:         true,
		AwsCis14IsoIec270022022:      true,
		AwsPciDss321:                 true,
		AwsNist800171Rev2:            true,
		LwAwsSecAdd10:                true,
		AwsNist80053Rev5:             true,
		AzureIso270012013Cis15:       true,
		AzureSoc2Cis15:               true,
		AzureCis15:                   true,
		AzureNistCsfCis15:            true,
		AzurePciDss321Cis15:          true,
		AzureHipaaCis15:              true,
		AzureNist80053Rev5Cis15:      true,
		AzureNist800171Rev2Cis15:     true,
		GcpHipaa2013:                 true,
		GcpSoc2:                      true,
		GcpPciDss321:                 true,
		GcpCis130NistCsf:             true,
		GcpCis130Nist80053Rev5:       true,
		GcpCmmc102:                   true,
		GcpCis130Nist800171Rev2:      true,
		GcpIso270012013:              true,
		GcpCis13:                     true,
		OpenShiftComplianceEvents:    true,
		K8SAuditLogEvents:            true,
		GcpComplianceEvents:          true,
		AgentEvents:                  true,
		AwsComplianceEvents:          true,
		AzureComplianceEvents:        true,
		AzureActivityLogEvents:       true,
		GcpAuditTrailEvents:          true,
		PlatformEvents:               true,
		AwsCloudtrailEvents:          true,
		IncidentEvents:               true,
		OpenShiftCompliance:          true,
		TrendReport:                  true,
		ContainerVulnerabilityReport: true,
		HostVulnerabilityReport:      true,
	}
}

func (all ReportRuleNotificationTypes) ToMap() map[string]bool {
	return reportRuleNotificationToMap(all)
}
