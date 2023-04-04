// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestParseNewReportDefinition(t *testing.T) {
	reportDefinitionsTests := []struct {
		Name     string
		Input    string
		ErrMsg   string
		Expected api.ReportDefinitionConfig
	}{
		{
			Name:     "test parse report definition response json",
			Input:    mockReportDefinitionResponseJson,
			Expected: expectedReport,
			ErrMsg:   "",
		},
		{
			Name:     "test parse report definition config json",
			Input:    mockReportDefinitionJson,
			Expected: expectedReport,
			ErrMsg:   "",
		},
		{
			Name:     "test parse report definition config yaml",
			Input:    mockReportDefinitionYaml,
			Expected: expectedReport,
			ErrMsg:   "",
		},
		{
			Name:     "test invalid",
			Input:    "invalid",
			Expected: api.ReportDefinitionConfig{},
			ErrMsg:   "unable to parse report definition file",
		},
	}

	for _, test := range reportDefinitionsTests {
		t.Run(test.Name, func(t *testing.T) {
			report, err := parseNewReportDefinition(test.Input)

			if err != nil {
				assert.Error(t, err, test.ErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, report.ReportName, test.Expected.ReportName)
				assert.Equal(t, report.DisplayName, test.Expected.DisplayName)
				assert.Equal(t, report.SubReportType, test.Expected.SubReportType)
				assert.Equal(t, report.ReportType, test.Expected.ReportType)
				assert.Equal(t, report.Sections[0].Policies, test.Expected.Sections[0].Policies)
				assert.Equal(t, report.Sections[0].Title, test.Expected.Sections[0].Title)
			}

		},
		)
	}

}

var (
	mockReportDefinitionResponseJson = `{
  "data": {
    "createdBy": "test.user@lacework.net",
    "createdTime": "2023-03-16T14:10:24Z",
    "displayName": "Mock Report Display Name",
    "reportDefinition": {
      "sections": [
        {
          "category": "",
          "policies": [
            "lacework-global-1"
          ],
          "title": "Test CLI"
        }
      ]
    },
    "reportDefinitionGuid": "ABCD123ABC4F4C9E3547DCC6E62A8C6D1234E45E9610999FE5FA51ABC",
    "reportName": "mockReportDefinition",
    "reportType": "COMPLIANCE",
    "subReportType": "AWS",
    "version": 1
  }
}
`

	mockReportDefinitionJson = `{
	"displayName": "Mock Report Display Name",
    "reportDefinition": {
      "sections": [
        {
          "policies": [
            "lacework-global-1"
          ],
          "title": "Test CLI"
        }
      ]
    },
    "reportName": "mockReportDefinition",
    "reportType": "COMPLIANCE",
    "subReportType": "AWS"
}
`

	mockReportDefinitionYaml = `
---
displayName: Mock Report Display Name
reportDefinition:
  sections:
  - policies:
    - lacework-global-1
    title: Test CLI
reportName: mockReportDefinition
reportType: COMPLIANCE
subReportType: AWS
`

	expectedReport = api.ReportDefinitionConfig{
		ReportName:    "mockReportDefinition",
		DisplayName:   "Mock Report Display Name",
		ReportType:    api.ReportDefinitionTypeCompliance.String(),
		SubReportType: api.ReportDefinitionSubTypeAws.String(),
		Sections:      []api.ReportDefinitionSection{{Policies: []string{"lacework-global-1"}, Title: "Test CLI"}},
	}
)
