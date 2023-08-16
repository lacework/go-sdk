//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/capturer"
)

func TestSplitIDAndAlias(t *testing.T) {
	cases := []struct {
		subjectText   string
		expectedID    string
		expectedAlias string
	}{
		// empty text will return empty id and alias
		{"", "", ""},
		// alias should not be empty
		{"()", "", ""},
		// minimum text that can be splitted
		{"a (b)", "a", "b"},
		// if we couldn't get the alias from the provided text
		// it means that the entire text is the id
		{"1234567890", "1234567890", ""},
		// other common test cases
		{"1234567890 (alias-example)", "1234567890", "alias-example"},
		{"proj-id-with-numbers (alias with spaces)", "proj-id-with-numbers", "alias with spaces"},
		{"only-project-id-123", "only-project-id-123", ""},
		// seriously, we should never have only the alias in the response ;-)
		{"(this should never happen)", "", "this should never happen"},
	}
	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actualID, actualAlias := splitIDAndAlias(kase.subjectText)
			assert.Equalf(t, kase.expectedID, actualID, "wrong id")
			assert.Equalf(t, kase.expectedAlias, actualAlias, "wrong alias")
		})
	}
}

func TestSplitGcpProjectsApiResponse(t *testing.T) {
	cases := []struct {
		subject  api.GcpConfigData
		expected cliComplianceGcpInfo
	}{
		// empty projects will return empty cli info
		{
			api.GcpConfigData{},
			cliComplianceGcpInfo{Projects: make([]cliComplianceIDAlias, 0)},
		},
		// real test case with NO alias
		{
			api.GcpConfigData{
				Organization: "1234567890123",
				Projects:     []string{"project-id-1", "project-id-2", "project-id-3", "project-id-4"},
			},
			cliComplianceGcpInfo{
				Organization: cliComplianceIDAlias{"1234567890123", ""},
				Projects: []cliComplianceIDAlias{
					{"project-id-1", ""},
					{"project-id-2", ""},
					{"project-id-3", ""},
					{"project-id-4", ""},
				},
			},
		},
		// real test case with alias
		{
			api.GcpConfigData{
				Organization: "1234567890123 (cool.org.alias.example.com)",
				Projects: []string{
					"id-1 (a test project)",
					"xmen-project (serious alias)",
					"disney-movies (Maybe Production)",
					"foo (bar)",
				},
			},
			cliComplianceGcpInfo{
				Organization: cliComplianceIDAlias{"1234567890123", "cool.org.alias.example.com"},
				Projects: []cliComplianceIDAlias{
					{"id-1", "a test project"},
					{"xmen-project", "serious alias"},
					{"disney-movies", "Maybe Production"},
					{"foo", "bar"},
				},
			},
		},
	}
	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equalf(t,
				kase.expected, splitGcpProjectsApiResponse(kase.subject),
				"there is a problem with this test case, please check",
			)
		})
	}
}

func TestDuplicateGcpAccountCheck(t *testing.T) {
	gcpOne := gcpProject{OrganizationID: "n/a", ProjectID: "1"}
	gcpTwo := gcpProject{OrganizationID: "n/a", ProjectID: "2"}
	gcpThree := gcpProject{OrganizationID: "n/a", ProjectID: "3"}
	mockGcpAccounts := []gcpProject{gcpOne, gcpTwo, gcpThree}

	duplicate := containsDuplicateProjectID(mockGcpAccounts, "1")
	different := containsDuplicateProjectID(mockGcpAccounts, "4")

	assert.True(t, duplicate)
	assert.False(t, different)
}

func TestCliListGcpListProjectsAndOrgsWithoutData(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, cliListGcpProjectsAndOrgs(api.CloudAccountsResponse{}, mockGcpConfigsResponse()))
	})
	assert.Contains(t, cliOutput, "There are no GCP integrations configured in your account.")

	t.Run("test JSON output", func(t *testing.T) {
		cli.EnableJSONOutput()
		defer cli.EnableHumanOutput()
		cliJSONOutput := capturer.CaptureOutput(func() {
			assert.Nil(t, cliListGcpProjectsAndOrgs(api.CloudAccountsResponse{}, mockGcpConfigsResponse()))
		})
		expectedJSON := `{
  "gcp_projects": []
}
`
		assert.Equal(t, expectedJSON, cliJSONOutput)
	})
}

func TestCliListGcpListProjectsAndOrgsWithDataEnabled(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, cliListGcpProjectsAndOrgs(mockGcpIntegrationsResponse(1, 1, 1), mockGcpConfigsResponse()))
	})
	// NOTE (@afiune): We purposely leave trailing spaces in this table, we need them!
	expectedTable := `
  ORGANIZATION ID           PROJECT ID            STATUS   
------------------+-----------------------------+----------
  n/a               gcr-jenkins-sandbox-274317    Enabled  
  n/a               techally-hipstershop-275821   Enabled  
  12345             techally-test                 Enabled  
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestCliListGcpListProjectsAndOrgsWithDataDisabled(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, cliListGcpProjectsAndOrgs(mockGcpIntegrationsResponse(0, 0, 1), mockGcpConfigsResponse()))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
  ORGANIZATION ID           PROJECT ID             STATUS   
------------------+-----------------------------+-----------
  n/a               gcr-jenkins-sandbox-274317    Disabled  
  n/a               techally-hipstershop-275821   Disabled  
  12345             techally-test                 Enabled   
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func mockGcpIntegrationsResponse(proj1Enabled, proj2Enabled, proj3Enabled int) api.CloudAccountsResponse {
	var response = api.CloudAccountsResponse{}
	err := json.Unmarshal([]byte(`
{
    "data": [
        {
            "createdOrUpdatedBy": "test@lacework.net",
            "createdOrUpdatedTime": "2022-05-18T15:34:29.872Z",
            "enabled": `+strconv.Itoa(proj1Enabled)+`,
            "intgGuid": "EXAMPLE_12345",
            "isOrg": 0,
            "name": "test",
            "state": {
                "ok": true,
                "lastUpdatedTime": 1652888077205,
                "lastSuccessfulTime": 1652888077205,
                "details": {}
            },
            "type": "GcpPubsub",
            "data": {
                "issueGrouping": "Events",
                "credentials": {
                    "clientId": "123456",
                    "privateKeyId": "",
                    "clientEmail": "test@test.iam.gserviceaccount.com",
                    "privateKey": ""
                },
                "idType": "PROJECT",
                "id": "gcr-jenkins-sandbox-274317"
            }
        },
        {
            "createdOrUpdatedBy": "test@lacework.net",
            "createdOrUpdatedTime": "2022-05-18T15:34:29.872Z",
            "enabled": `+strconv.Itoa(proj2Enabled)+`,
            "intgGuid": "EXAMPLE_12345",
            "isOrg": 0,
            "name": "test",
            "state": {
                "ok": true,
                "lastUpdatedTime": 1652888077205,
                "lastSuccessfulTime": 1652888077205,
                "details": {}
            },
            "type": "GcpPubsub",
            "data": {
                "issueGrouping": "Events",
                "credentials": {
                    "clientId": "123456",
                    "privateKeyId": "",
                    "clientEmail": "test@test.iam.gserviceaccount.com",
                    "privateKey": ""
                },
                "idType": "PROJECT",
                "id": "techally-hipstershop-275821"
            }
        },
        {
            "createdOrUpdatedBy": "test@lacework.net",
            "createdOrUpdatedTime": "2022-05-18T15:34:29.872Z",
            "enabled": `+strconv.Itoa(proj3Enabled)+`,
            "intgGuid": "EXAMPLE_12345",
            "isOrg": 0,
            "name": "test",
            "state": {
                "ok": true,
                "lastUpdatedTime": 1652888077205,
                "lastSuccessfulTime": 1652888077205,
                "details": {}
            },
            "type": "GcpPubsub",
            "data": {
                "issueGrouping": "Events",
                "credentials": {
                    "clientId": "123456",
                    "privateKeyId": "",
                    "clientEmail": "test@test.iam.gserviceaccount.com",
                    "privateKey": ""
                },
                "idType": "PROJECT",
                "id": "techally-test"
            }
        }
]
}`), &response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func mockGcpConfigsResponse() api.GcpConfigsResponse {
	return api.GcpConfigsResponse{
		Data: []api.GcpConfigData{{
			Organization: "12345",
			Projects:     []string{"techally-test (techally-test)"},
		}},
	}
}
