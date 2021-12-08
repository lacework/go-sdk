// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
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

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestBuildReportRules(t *testing.T) {
	reportDetails := buildReportRuleDetailsTable(mockReportRule)
	assert.Equal(t, reportDetails, reportRuleOutput)
}

var (
	mockReportRule = api.ReportRule{
		Guid:               "TECHALLY_94D51929DE541793243644C755067305AF564A62AD63713",
		Type:               api.ReportRuleEventType,
		EmailAlertChannels: []string{"TECHALLY_2F0C086E17AB64BEC84F4A5FF8A3F068CF2CE15847BCBCA"},
		Filter: api.ReportRuleFilter{
			Name:                 "Report Rule Test",
			Enabled:              1,
			Description:          "Report Rule created by Cli",
			Severity:             []int{2, 3},
			ResourceGroups:       []string{"TECHALLY_8416B4ADCED28565254842AA5906B729174653E1725F107"},
			CreatedOrUpdatedTime: "2021-11-29T16:36:36Z",
			CreatedOrUpdatedBy:   "darren.murray@lacework.net",
		},
		ReportNotificationTypes: api.AllReportRuleNotifications,
	}

	reportRuleOutput = `               ALERT RULE DETAILS                
-------------------------------------------------
    SEVERITIES     High, Medium                  
    DESCRIPTION    Report Rule created by Cli    
    UPDATED BY     darren.murray@lacework.net    
    LAST UPDATED                                 
                                                 
     NOTIFICATION TYPES       ENABLED  
----------------------------+----------
  AgentEvents                 True     
  AwsCisS3                    True     
  AwsCloudtrailEvents         True     
  AwsComplianceEvents         True     
  AwsHipaa                    True     
  AwsIso2700                  True     
  AwsNist800171Rev2           True     
  AwsNist80053Rev4            True     
  AwsPci                      True     
  AwsSoc                      True     
  AwsSocRev2                  True     
  AzureActivityLogEvents      True     
  AzureCis                    True     
  AzureCis131                 True     
  AzureComplianceEvents       True     
  AzurePci                    True     
  AzureSoc                    True     
  GcpAuditTrailEvents         True     
  GcpCis                      True     
  GcpCis12                    True     
  GcpComplianceEvents         True     
  GcpHipaa                    True     
  GcpHipaaRev2                True     
  GcpIso27001                 True     
  GcpK8s                      True     
  GcpPci                      True     
  GcpPciRev2                  True     
  GcpSoc                      True     
  GcpSocRev2                  True     
  OpenShiftCompliance         True     
  OpenShiftComplianceEvents   True     
  PlatformEvents              True     
  TrendReport                 True     

                    EMAIL ALERT CHANNELS                    
------------------------------------------------------------
  TECHALLY_2F0C086E17AB64BEC84F4A5FF8A3F068CF2CE15847BCBCA  

                      RESOURCE GROUPS                       
------------------------------------------------------------
  TECHALLY_8416B4ADCED28565254842AA5906B729174653E1725F107  
`
)
