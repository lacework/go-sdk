//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

func TestIntegrationsNewPagerDutyAlertChannel(t *testing.T) {
	subject := api.NewPagerDutyAlertChannel("integration_name",
		api.PagerDutyData{
			IntegrationKey:   "1234567890abcd1234567890",
			MinAlertSeverity: 3,
		},
	)
	assert.Equal(t, api.PagerDutyIntegration.String(), subject.Type)
	assert.Equal(t, api.MediumAlertLevel, subject.Data.MinAlertSeverity)
}

func TestIntegrationsCreatePagerDutyAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreatePagerDutyAlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "PAGER_DUTY_API", "wrong integration type")
			assert.Contains(t, body, "1234567890abcd1234567890", "wrong integration_key")
			assert.Contains(t, body, "MIN_ALERT_SEVERITY\":3", "wrong alert severity")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, pagerDutyIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewPagerDutyAlertChannel("integration_name",
		api.PagerDutyData{
			IntegrationKey:   "1234567890abcd1234567890",
			MinAlertSeverity: 3,
		},
	)
	assert.Equal(t, "integration_name", data.Name, "PagerDuty integration name mismatch")
	assert.Equal(t, "PAGER_DUTY_API", data.Type, "a new PagerDuty integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new PagerDuty integration should be enabled")

	response, err := c.Integrations.CreatePagerDutyAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "1234567890abcd1234567890", resData.Data.IntegrationKey)
		assert.Equal(t, api.AlertLevel(3), resData.Data.MinAlertSeverity)
	}
}

func TestIntegrationsGetPagerDutyAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetPagerDutyAlertChannel should be a GET method")
		fmt.Fprintf(w, pagerDutyIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetPagerDutyAlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "1234567890abcd1234567890", resData.Data.IntegrationKey)
		assert.Equal(t, api.AlertLevel(3), resData.Data.MinAlertSeverity)
	}
}

func TestIntegrationsUpdatePagerDutyAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdatePagerDutyAlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "PAGER_DUTY_API", "wrong integration type")
			assert.Contains(t, body, "1234567890abcd1234567890", "wrong integration_key")
			assert.Contains(t, body, "MIN_ALERT_SEVERITY\":3", "wrong alert severity")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, pagerDutyIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewPagerDutyAlertChannel("integration_name",
		api.PagerDutyData{
			IntegrationKey:   "1234567890abcd1234567890",
			MinAlertSeverity: 3,
		},
	)
	assert.Equal(t, "integration_name", data.Name, "PagerDuty integration name mismatch")
	assert.Equal(t, "PAGER_DUTY_API", data.Type, "a new PagerDuty integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new PagerDuty integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdatePagerDutyAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListPagerDutyAlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/PAGER_DUTY_API",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListPagerDutyAlertChannel should be a GET method")
			fmt.Fprintf(w, pagerDutyMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListPagerDutyAlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func pagerDutyIntegrationJsonResponse(intgGUID string) string {
	return `
{
  "data": [` + singlePagerDutyIntegration(intgGUID) + `],
  "ok": true,
  "message": "SUCCESS"
}
`
}

func pagerDutyMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singlePagerDutyIntegration(guid))
	}
	return `
{
"data": [` + strings.Join(integrations, ", ") + `],
"ok": true,
"message": "SUCCESS"
}
`
}

func singlePagerDutyIntegration(id string) string {
	return `
{
  "INTG_GUID": "` + id + `",
  "CREATED_OR_UPDATED_BY": "user@email.com",
  "CREATED_OR_UPDATED_TIME": "2020-Jul-16 19:59:22 UTC",
  "DATA": {
    "ISSUE_GROUPING": "Events",
    "MIN_ALERT_SEVERITY": 3,
    "API_INTG_KEY": "1234567890abcd1234567890"
  },
  "ENABLED": 1,
  "IS_ORG": 0,
  "NAME": "integration_name",
  "STATE": {
    "lastSuccessfulTime": "2020-Jul-16 18:26:54 UTC",
    "lastUpdatedTime": "2020-Jul-16 18:26:54 UTC",
    "ok": true
  },
  "TYPE": "PAGER_DUTY_API",
  "TYPE_NAME": "PAGER_DUTY_API"
}
`
}
