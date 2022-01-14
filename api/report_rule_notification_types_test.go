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

package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestNewNotificationType(t *testing.T) {
	var gcp api.GcpReportRuleNotifications
	var daily api.DailyEventsReportRuleNotifications

	gcpMap := map[string]bool{
		"GcpCis":       true,
		"GcpHipaa":     true,
		"GcpHipaaRev2": true,
		"GcpIso27001":  true,
		"GcpCis12":     false,
		"GcpK8s":       false,
		"GcpPci":       true,
		"GcpPciRev2":   true,
		"GcpSoc":       false,
		"GcpSocRev2":   true}

	dailyMap := map[string]bool{"OpenShiftComplianceEvents": true, "OpenShiftCompliance": true}

	err := api.TransformReportRuleNotification(gcpMap, &gcp)

	if assert.NoError(t, err) {
		assert.True(t, gcp.GcpCis)
		assert.True(t, gcp.GcpHipaa)
		assert.True(t, gcp.GcpHipaaRev2)
		assert.True(t, gcp.GcpIso27001)
		assert.True(t, gcp.GcpPci)
		assert.True(t, gcp.GcpPciRev2)
		assert.False(t, gcp.GcpCis12)
		assert.False(t, gcp.GcpK8s)
		assert.False(t, gcp.GcpSoc)
	}

	err = api.TransformReportRuleNotification(dailyMap, &daily)
	if assert.NoError(t, err) {
		assert.True(t, daily.OpenShiftComplianceEvents)
		assert.True(t, daily.OpenShiftCompliance)
	}
}

func TestNotificationTypeToMap(t *testing.T) {
	allMap := api.AllReportRuleNotifications.ToMap()
	gcpMap := api.GcpReportRuleNotifications{GcpCis: false, GcpCis12: true}.ToMap()
	awsMap := api.AwsReportRuleNotifications{AwsCisS3: true, AwsPci: true, AwsHipaa: false}.ToMap()
	azMap := api.AzureReportRuleNotifications{AzureCis: false, AzurePci: true, AzureSoc: false}.ToMap()
	dailyMap := api.DailyEventsReportRuleNotifications{AgentEvents: false, GcpAuditTrailEvents: true}.ToMap()
	weeklyMap := api.WeeklyEventsReportRuleNotifications{TrendReport: true}.ToMap()

	// allMap
	assert.True(t, allMap["GcpSocRev2"])
	assert.True(t, allMap["AzureCis"])
	assert.True(t, allMap["PlatformEvents"])
	assert.True(t, allMap["GcpSocRev2"])
	assert.True(t, allMap["AwsComplianceEvents"])
	assert.True(t, allMap["TrendReport"])

	// gcpMap
	assert.True(t, gcpMap["GcpCis12"])
	assert.False(t, gcpMap["GcpCis"])

	// awsMap
	assert.True(t, awsMap["AwsCisS3"])
	assert.True(t, awsMap["AwsPci"])
	assert.False(t, awsMap["AwsHipaa"])

	//azMap
	assert.True(t, azMap["AzurePci"])
	assert.False(t, azMap["AzureSoc"])
	assert.False(t, azMap["AzureCis"])

	//dailyMap
	assert.True(t, dailyMap["GcpAuditTrailEvents"])
	assert.False(t, dailyMap["AgentEvents"])

	//weeklyMap
	assert.True(t, weeklyMap["TrendReport"])
}

func TestReportRuleDuplicateTypes(t *testing.T) {
	_, errOne := api.NewReportRule("rule_name",
		api.ReportRuleConfig{
			EmailAlertChannels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
			Description:        "This is a test report rule",
			Severities:         api.ReportRuleSeverities{api.ReportRuleSeverityHigh},
			ResourceGroups:     []string{"TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB"},
			NotificationTypes: api.ReportRuleNotifications{
				api.GcpReportRuleNotifications{GcpCis: true},
				api.GcpReportRuleNotifications{GcpCis: true},
			},
		},
	)

	_, errTwo := api.NewReportRule("rule_name",
		api.ReportRuleConfig{
			EmailAlertChannels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
			Description:        "This is a test report rule",
			Severities:         api.ReportRuleSeverities{api.ReportRuleSeverityHigh},
			ResourceGroups:     []string{"TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB"},
			NotificationTypes: api.ReportRuleNotifications{
				api.GcpReportRuleNotifications{GcpCis: true},
				api.AwsReportRuleNotifications{AwsHipaa: true},
				api.AzureReportRuleNotifications{AzureSoc: false},
				api.WeeklyEventsReportRuleNotifications{TrendReport: true},
				api.DailyEventsReportRuleNotifications{GcpAuditTrailEvents: true},
				api.GcpReportRuleNotifications{GcpCis: true},
			},
		},
	)

	_, errThree := api.NewReportRule("rule_name",
		api.ReportRuleConfig{
			EmailAlertChannels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
			Description:        "This is a test report rule",
			Severities:         api.ReportRuleSeverities{api.ReportRuleSeverityHigh},
			ResourceGroups:     []string{"TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB"},
			NotificationTypes: api.ReportRuleNotifications{
				api.GcpReportRuleNotifications{GcpCis: true},
				api.AwsReportRuleNotifications{AwsHipaa: true},
				api.AzureReportRuleNotifications{AzureSoc: false},
				api.WeeklyEventsReportRuleNotifications{TrendReport: true},
				api.DailyEventsReportRuleNotifications{GcpAuditTrailEvents: true},
				api.ReportRuleNotificationTypes{GcpCis: false},
			},
		},
	)

	reportRule, errFour := api.NewReportRule("rule_name",
		api.ReportRuleConfig{
			EmailAlertChannels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
			Description:        "This is a test report rule",
			Severities:         api.ReportRuleSeverities{api.ReportRuleSeverityHigh},
			ResourceGroups:     []string{"TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB"},
			NotificationTypes: api.ReportRuleNotifications{
				api.GcpReportRuleNotifications{GcpCis: true},
				api.AwsReportRuleNotifications{AwsHipaa: true},
				api.AzureReportRuleNotifications{AzureSoc: false},
				api.WeeklyEventsReportRuleNotifications{TrendReport: true},
				api.DailyEventsReportRuleNotifications{GcpAuditTrailEvents: true},
			},
		},
	)

	assert.Error(t, errOne, "notification types contains a duplicate type")
	assert.Error(t, errTwo, "notification types contains a duplicate type")
	assert.Error(t, errThree, "notification types contains a duplicate type")
	assert.NoError(t, errFour, "report rule should not return error")
	assert.True(t, reportRule.ReportNotificationTypes.AwsHipaa)
	assert.True(t, reportRule.ReportNotificationTypes.GcpCis)
	assert.True(t, reportRule.ReportNotificationTypes.TrendReport)
	assert.True(t, reportRule.ReportNotificationTypes.GcpAuditTrailEvents)
	assert.False(t, reportRule.ReportNotificationTypes.AzureSoc)
}
