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
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestReportRulesGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ReportRules/%s", intgGUID)
		reportRule = singleMockReportRule(intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateReportRuleResponse(reportRule))
			}
		},
	)

	fakeServer.MockAPI("ReportRules/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when report rule exists", func(t *testing.T) {
		var response api.ReportRuleResponse
		err := c.V2.ReportRules.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.Guid)
			assert.Equal(t, "rule_name", response.Data.Filter.Name)
			assert.Equal(t, "Report", response.Data.Type)
		}
	})

	t.Run("when report rule does NOT exist", func(t *testing.T) {
		var response api.ReportRuleResponse
		err := c.V2.ReportRules.Get("UNKNOWN_INTG_GUID", response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ReportRules/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestReportRulesDelete(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		apiPath     = fmt.Sprintf("ReportRules/%s", intgGUID)
		reportRule  = singleMockReportRule(intgGUID)
		getResponse = generateReportRuleResponse(reportRule)
		fakeServer  = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if getResponse != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, getResponse)
				case "DELETE":
					// once deleted, empty the getResponse so that
					// further GET requests return 404s
					getResponse = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	fakeServer.MockAPI("ReportRules/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify report rule exists", func(t *testing.T) {
		var response api.ReportRuleResponse
		err := c.V2.ReportRules.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.Guid)
			assert.Equal(t, "rule_name", response.Data.Filter.Name)
			assert.Equal(t, "Report", response.Data.Type)
		}
	})

	t.Run("when report rule has been deleted", func(t *testing.T) {
		err := c.V2.ReportRules.Delete(intgGUID)
		assert.Nil(t, err)

		var response api.ReportRuleResponse
		err = c.V2.ReportRules.Get(intgGUID, &response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ReportRules/MOCK_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestReportRulesList(t *testing.T) {
	var (
		allGUIDs    []string
		reportRules = generateGuids(&allGUIDs, 3)
		expectedLen = len(allGUIDs)
		fakeServer  = lacework.MockServer()
	)

	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("ReportRules",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			reportRules := []string{
				generateReportRules(reportRules),
			}
			fmt.Fprintf(w,
				generateReportRulesResponse(
					strings.Join(reportRules, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ReportRules.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.Guid)
	}
}

func TestReportRuleUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ReportRules/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Update() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "GUID missing")
			assert.Contains(t, body, "rule_name", "report rule name is missing")
			assert.Contains(t, body, "Report", "wrong report rule type")
			assert.Contains(t, body, "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA", "missing intgGuidList")
		}

		fmt.Fprintf(w, generateReportRuleResponse(singleMockReportRule(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	reportRule, err := api.NewReportRule("rule_name",
		api.ReportRuleConfig{
			EmailAlertChannels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
			Description:        "This is a test report rule",
			Severities:         api.ReportRuleSeverities{api.ReportRuleSeverityHigh},
			ResourceGroups:     []string{"TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB"},
			NotificationTypes:  []api.ReportRuleNotification{api.GcpReportRuleNotifications{GcpCis: true}},
		},
	)
	assert.Equal(t, "rule_name", reportRule.Filter.Name, "report rule name mismatch")
	assert.Equal(t, "Report", reportRule.Type, "a new report rule should match its type")
	assert.Equal(t, 1, reportRule.Filter.Enabled, "a new report rule should be enabled")
	reportRule.Guid = intgGUID

	response, err := c.V2.ReportRules.Update(reportRule)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.Guid)
		assert.Equal(t, response.Data.NotificationTypes.GcpCis, true)
		assert.Contains(t, response.Data.Filter.ResourceGroups, "TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB")
		assert.Contains(t, response.Data.EmailAlertChannels, "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA")
	}
}

func TestReportRuleSeveritiesFromIntSlice(t *testing.T) {
	var (
		sevInts            = []int{1, 2, 55, 0, -1, 4}
		expectedSeverities = api.ReportRuleSeverities{api.ReportRuleSeverityCritical, api.ReportRuleSeverityHigh, api.ReportRuleSeverityLow}
	)

	severities := api.NewReportRuleSeveritiesFromIntSlice(sevInts)
	assert.Equal(t, expectedSeverities, severities)
}

func TestReportRuleSeveritiesToString(t *testing.T) {
	var (
		sevInts                 = []int{1, 2, 55, 0, -1, 4}
		expectedSeverityStrings = []string{"Critical", "High", "Low"}
	)

	severityStrings := api.NewReportRuleSeveritiesFromIntSlice(sevInts).ToStringSlice()
	assert.Equal(t, expectedSeverityStrings, severityStrings)
}

func generateReportRules(guids []string) string {
	reportRules := make([]string, len(guids))
	for i, guid := range guids {
		reportRules[i] = singleMockReportRule(guid)
	}
	return strings.Join(reportRules, ", ")
}

func generateReportRulesResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateReportRuleResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleMockReportRule(id string) string {
	return fmt.Sprintf(`
    {
      "filters": {
            "name": "rule_name",
            "description": "This is a test report rule",
            "createdOrUpdatedBy": "darren.murray@lacework.net",
            "createdOrUpdatedTime": "1634205829636",
            "enabled": 1,
            "resourceGroups": [
                "TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB"
            ],
            "severity": [
                2
            ]
	  },
      "mcGuid": %q,
        "intgGuidList": [
            "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"
        ],
        "reportNotificationTypes": {
			"gcpCis": true
		},
      "name": "rule_name",
      "type": "Report"
    }
	`, id)
}
