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
		Props: api.ReportDefinitionProps{
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
                                                     
                        POLICIES                        
--------------------------------------------------------
         TITLE                    POLICY                
  ------------------+---------------------------------  
    Example Section   lacework-global-22,               
                      lacework-global-78                
                                                        
                                                        
`
)