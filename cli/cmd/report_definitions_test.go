// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestBuildReportDefinitions(t *testing.T) {
	reportDetails := buildReportDefinitionDetailsTable(mockReportDefinition)
	assert.Equal(t, reportDetails, reportDefinitionOutput)
}

func TestFilterReportDefinitionsWithSubType(t *testing.T) {
	defer func() { reportDefinitionsCmdState.SubType = "" }()
	var reportDefinitionResponse api.ReportDefinitionsResponse
	reportDefinitionResponse.Data = []api.ReportDefinition{mockReportDefinition, mockReportDefinitionGCP, mockReportDefinitionAzure}

	reportDefinitionsTableTest := []struct {
		Name     string
		Input    api.ReportDefinitionsResponse
		Expected []api.ReportDefinition
		Cloud    string
	}{
		{
			Name:     "Test Filter AWS report Definitions",
			Input:    reportDefinitionResponse,
			Expected: []api.ReportDefinition{mockReportDefinition},
			Cloud:    "AWS",
		},
		{
			Name:     "Test Filter GCP report Definitions",
			Input:    reportDefinitionResponse,
			Expected: []api.ReportDefinition{mockReportDefinitionGCP},
			Cloud:    "GCP",
		},
		{
			Name:     "Test Filter Azure report Definitions",
			Input:    reportDefinitionResponse,
			Expected: []api.ReportDefinition{mockReportDefinitionAzure},
			Cloud:    "Azure",
		},
	}

	for _, rdtt := range reportDefinitionsTableTest {
		t.Run(rdtt.Name, func(t *testing.T) {
			reportDefinitionsCmdState.SubType = rdtt.Cloud
			filterReportDefinitions(&rdtt.Input)
			assert.Len(t, rdtt.Input.Data, 1)
			assert.Equal(t, rdtt.Input.Data, rdtt.Expected)
		})
	}
}

func TestReportDiff(t *testing.T) {
	reportDiffTableTest := []struct {
		Name     string
		InputOne diffCfg
		InputTwo diffCfg
		Expected string
	}{
		{
			Name:     "diff valid report definitions",
			InputOne: diffCfg{name: "mock-1", object: mockReportDefinition},
			InputTwo: diffCfg{name: "mock-2", object: mockReportDefinitionGCP},
			Expected: mockReportDefinitionDiff,
		},
		{
			Name:     "diff empty report definitions",
			InputOne: diffCfg{name: "mock-1", object: ""},
			InputTwo: diffCfg{name: "mock-2", object: ""},
			Expected: "",
		},
		{
			Name:     "diff single word",
			InputOne: diffCfg{name: "mock-1", object: "testOne"},
			InputTwo: diffCfg{name: "mock-2", object: "testTwo"},
			Expected: mockDiffSingleWord,
		},
		{
			Name:     "diff no change",
			InputOne: diffCfg{name: "mock-1", object: mockReportDefinition},
			InputTwo: diffCfg{name: "mock-1", object: mockReportDefinition},
			Expected: "",
		},
	}

	for _, test := range reportDiffTableTest {
		t.Run(test.Name, func(t *testing.T) {
			diff, err := diffAsYamlString(test.InputOne, test.InputTwo)
			assert.NoError(t, err)
			assert.Equal(t, diff, test.Expected)
		})
	}
}

var created, _ = time.Parse(time.RFC3339, "2022-09-09T10:35:16Z")

var (
	mockReportDefinition = api.ReportDefinition{
		ReportDefinitionGuid: "EXAMPLE_GUID",
		ReportName:           "My Custom Report",
		DisplayName:          "My Custom Report Display",
		ReportType:           "Compliance",
		SubReportType:        "AWS",
		ReportDefinitionDetails: api.ReportDefinitionDetails{
			Sections: []api.ReportDefinitionSection{{
				Category: "1.0.0",
				Title:    "Example Section",
				Policies: []string{"lacework-global-22", "lacework-global-78"},
			}},
		},
		Props: &api.ReportDefinitionProps{
			Engine: "lpp",
		},
		DistributionType: "pdf",
		Frequency:        "daily",
		Version:          2,
		CreatedBy:        "SYSTEM",
		CreatedTime:      &created,
		Enabled:          1,
	}
	mockReportDefinitionGCP = api.ReportDefinition{
		ReportDefinitionGuid: "EXAMPLE_GUID",
		ReportName:           "My Custom Report",
		DisplayName:          "My Custom Report Display",
		ReportType:           "Compliance",
		SubReportType:        "GCP",
		ReportDefinitionDetails: api.ReportDefinitionDetails{
			Sections: []api.ReportDefinitionSection{{
				Category: "1.0.0",
				Title:    "Example Section",
				Policies: []string{"lacework-global-22", "lacework-global-78"},
			}},
		},
		Props: &api.ReportDefinitionProps{
			Engine: "lpp",
		},
		DistributionType: "pdf",
		Frequency:        "daily",
		Version:          2,
		CreatedBy:        "SYSTEM",
		CreatedTime:      &created,
		Enabled:          1,
	}
	mockReportDefinitionAzure = api.ReportDefinition{
		ReportDefinitionGuid: "EXAMPLE_GUID",
		ReportName:           "My Custom Report",
		DisplayName:          "My Custom Report Display",
		ReportType:           "Compliance",
		SubReportType:        "Azure",
		ReportDefinitionDetails: api.ReportDefinitionDetails{
			Sections: []api.ReportDefinitionSection{{
				Category: "1.0.0",
				Title:    "Example Section",
				Policies: []string{"lacework-global-22", "lacework-global-78"},
			}},
		},
		Props: &api.ReportDefinitionProps{
			Engine: "lpp",
		},
		DistributionType: "pdf",
		Frequency:        "daily",
		Version:          2,
		CreatedBy:        "SYSTEM",
		CreatedTime:      &created,
		Enabled:          1,
	}

	reportDefinitionOutput = `              REPORT DEFINITION DETAILS              
-----------------------------------------------------
    FREQUENCY       daily                            
    ENGINE          lpp                              
    RELEASE LABEL                                    
    UPDATED BY      SYSTEM                           
    LAST UPDATED    2022-09-09 10:35:16 +0000 UTC    
    VERSION         2                                
                                                     
                        POLICIES                        
--------------------------------------------------------
         TITLE                    POLICY                
  ------------------+---------------------------------  
    Example Section   lacework-global-22,               
                      lacework-global-78                
                                                        
                                                        
`
	mockReportDefinitionDiff = `--- mock-1
+++ mock-2
@@ -2,7 +2,7 @@
 reportName: My Custom Report
 displayName: My Custom Report Display
 reportType: Compliance
-subReportType: AWS
+subReportType: GCP
 reportDefinition:
     sections:
         - category: 1.0.0

`

	mockDiffSingleWord = `--- mock-1
+++ mock-2
@@ -1,2 +1,2 @@
-testOne
+testTwo
 

`
)
