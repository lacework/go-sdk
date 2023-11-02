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

func TestAlertRulesGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertRules/%s", intgGUID)
		alertRule  = singleMockAlertRule(intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateAlertRuleResponse(alertRule))
			}
		},
	)

	fakeServer.MockAPI("AlertRules/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when alert rule exists", func(t *testing.T) {
		var response api.AlertRuleResponse
		err := c.V2.AlertRules.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.Guid)
			assert.Equal(t, "rule_name", response.Data.Filter.Name)
			assert.Equal(t, "Event", response.Data.Type)
		}
	})

	t.Run("when alert rule does NOT exist", func(t *testing.T) {
		var response api.AlertRuleResponse
		err := c.V2.AlertRules.Get("UNKNOWN_INTG_GUID", response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/AlertRules/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestAlertRulesDelete(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		apiPath     = fmt.Sprintf("AlertRules/%s", intgGUID)
		alertRule   = singleMockAlertRule(intgGUID)
		getResponse = generateAlertRuleResponse(alertRule)
		fakeServer  = lacework.MockServer()
	)
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

	fakeServer.MockAPI("AlertRules/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify alert rule exists", func(t *testing.T) {
		var response api.AlertRuleResponse
		err := c.V2.AlertRules.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.Guid)
			assert.Equal(t, "rule_name", response.Data.Filter.Name)
			assert.Equal(t, "Event", response.Data.Type)
		}
	})

	t.Run("when alert rule has been deleted", func(t *testing.T) {
		err := c.V2.AlertRules.Delete(intgGUID)
		assert.Nil(t, err)

		var response api.AlertRuleResponse
		err = c.V2.AlertRules.Get(intgGUID, &response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/AlertRules/MOCK_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestAlertRulesList(t *testing.T) {
	var (
		allGUIDs    []string
		alertRules  = generateGuids(&allGUIDs, 3)
		expectedLen = len(allGUIDs)
		fakeServer  = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("AlertRules",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			alertRules := []string{
				generateAlertRules(alertRules),
			}
			fmt.Fprintf(w,
				generateAlertRulesResponse(
					strings.Join(alertRules, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertRules.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.Guid)
	}
}

func TestAlertRuleUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertRules/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Update() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "GUID missing")
			assert.Contains(t, body, "rule_name", "alert rule name is missing")
			assert.Contains(t, body, "Event", "wrong alert rule type")
			assert.Contains(t, body, "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA", "missing intgGuidList")
		}

		fmt.Fprintf(w, generateAlertRuleResponse(singleMockAlertRule(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	alertRule := api.NewAlertRule("rule_name",
		api.AlertRuleConfig{
			Channels:           []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
			Description:        "This is a test alert rule",
			Severities:         api.AlertRuleSeverities{api.AlertRuleSeverityHigh},
			ResourceGroups:     []string{"TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB"},
			AlertSubCategories: []string{"Compliance", "SystemCall"},
			AlertSources:       []string{"AWS", "Agent", "K8s"},
			AlertCategories:    []string{"Policy", "Anomaly"},
		},
	)
	assert.Equal(t, "rule_name", alertRule.Filter.Name, "alert rule name mismatch")
	assert.Equal(t, "Event", alertRule.Type, "a new alert rule should match its type")
	assert.Equal(t, 1, alertRule.Filter.Enabled, "a new alert rule should be enabled")
	alertRule.Guid = intgGUID

	response, err := c.V2.AlertRules.Update(alertRule)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.Guid)
		assert.Contains(t, response.Data.Filter.AlertSubCategories, "Compliance", "SystemCall")
		assert.Contains(t, response.Data.Filter.AlertCategories, "Policy", "Anomaly")
		assert.Contains(t, response.Data.Filter.AlertSources, "AWS", "Agent", "K8s")
		assert.Contains(t, response.Data.Filter.ResourceGroups, "TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB")
		assert.Contains(t, response.Data.Channels, "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA")
	}
}

func TestAlertRuleSeveritiesFromIntSlice(t *testing.T) {
	var (
		sevInts            = []int{1, 2, 55, 0, -1, 4}
		expectedSeverities = api.AlertRuleSeverities{api.AlertRuleSeverityCritical, api.AlertRuleSeverityHigh, api.AlertRuleSeverityLow}
	)

	severities := api.NewAlertRuleSeveritiesFromIntSlice(sevInts)
	assert.Equal(t, expectedSeverities, severities)
}

func TestAlertRuleSeveritiesToString(t *testing.T) {
	var (
		sevInts                 = []int{1, 2, 55, 0, -1, 4}
		expectedSeverityStrings = []string{"Critical", "High", "Low"}
	)

	severityStrings := api.NewAlertRuleSeveritiesFromIntSlice(sevInts).ToStringSlice()
	assert.Equal(t, expectedSeverityStrings, severityStrings)
}

func generateAlertRules(guids []string) string {
	alertRules := make([]string, len(guids))
	for i, guid := range guids {
		alertRules[i] = singleMockAlertRule(guid)
	}
	return strings.Join(alertRules, ", ")
}

func generateAlertRulesResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateAlertRuleResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleMockAlertRule(id string) string {
	return fmt.Sprintf(`
    {
      "filters": {
            "name": "rule_name",
            "description": "This is a test alert rule",
            "createdOrUpdatedBy": "darren.murray@lacework.net",
            "createdOrUpdatedTime": "1634205829636",
            "enabled": 1,
            "resourceGroups": [
                "TECHALLY_100000000000AAAAAAAAAAAAAAAAAAAB"
            ],
            "severity": [
                2
            ],
            "subCategory": [
                "Compliance",
                "SystemCall"
            ],
            "category": [
                "Policy",
                "Anomaly"
            ],
            "source": [
                "AWS",
                "Agent",
                "K8s"
            ]
    },
      "mcGuid": %q,
        "intgGuidList": [
            "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"
        ],
      "name": "rule_name",
      "type": "Event"
    }
	`, id)
}
