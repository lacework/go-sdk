//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

func TestIntegrationsNewNewRelicAlertChannel(t *testing.T) {
	subject := api.NewNewRelicAlertChannel("integration_name",
		api.NewRelicChannelData{
			AccountID: 2338053,
			InsertKey: "x-xx-xxxxxxxxxxxxxxxxxx",
		},
	)
	assert.Equal(t, api.NewRelicChannelIntegration.String(), subject.Type)
}

func TestIntegrationsCreateNewRelicAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateNewRelicAlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "NEW_RELIC_INSIGHTS", "wrong integration type")
			assert.Contains(t, body, "2338053", "wrong account id")
			assert.Contains(t, body, "x-xx-xxxxxxxxxxxxxxxxxx", "wrong insert key")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, newRelicChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewNewRelicAlertChannel("integration_name",
		api.NewRelicChannelData{
			AccountID: 2338053,
			InsertKey: "x-xx-xxxxxxxxxxxxxxxxxx",
		},
	)
	assert.Equal(t, "integration_name", data.Name, "NewRelicChannel integration name mismatch")
	assert.Equal(t, "NEW_RELIC_INSIGHTS", data.Type, "a new NewRelicChannel integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new NewRelicChannel integration should be enabled")

	response, err := c.Integrations.CreateNewRelicAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, 2338053, resData.Data.AccountID)
		assert.Equal(t, "x-xx-xxxxxxxxxxxxxxxxxx", resData.Data.InsertKey)
	}
}

func TestIntegrationsGetNewRelicAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetNewRelicAlertChannel should be a GET method")
		fmt.Fprintf(w, newRelicChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetNewRelicAlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, 2338053, resData.Data.AccountID)
		assert.Equal(t, "x-xx-xxxxxxxxxxxxxxxxxx", resData.Data.InsertKey)
	}
}

func TestIntegrationsUpdateNewRelicAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateNewRelicAlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "NEW_RELIC_INSIGHTS", "wrong integration type")
			assert.Contains(t, body, "2338053", "wrong account id")
			assert.Contains(t, body, "x-xx-xxxxxxxxxxxxxxxxxx", "wrong insert key")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, newRelicChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewNewRelicAlertChannel("integration_name",
		api.NewRelicChannelData{
			AccountID: 2338053,
			InsertKey: "x-xx-xxxxxxxxxxxxxxxxxx",
		},
	)
	assert.Equal(t, "integration_name", data.Name, "NewRelicChannel integration name mismatch")
	assert.Equal(t, "NEW_RELIC_INSIGHTS", data.Type, "a new NewRelicChannel integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new NewRelicChannel integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateNewRelicAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListNewRelicAlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/NEW_RELIC_INSIGHTS",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListNewRelicAlertChannel should be a GET method")
			fmt.Fprintf(w, newRelicChanMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListNewRelicAlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func newRelicChannelIntegrationJsonResponse(intgGUID string) string {
	return `
{
  "data": [` + singleNewRelicChanIntegration(intgGUID) + `],
  "ok": true,
  "message": "SUCCESS"
}
`
}

func newRelicChanMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleNewRelicChanIntegration(guid))
	}
	return `
{
"data": [` + strings.Join(integrations, ", ") + `],
"ok": true,
"message": "SUCCESS"
}
`
}

func singleNewRelicChanIntegration(id string) string {
	return `
		{
			"INTG_GUID": "` + id + `",
			"NAME": "integration_name",
			"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
			"CREATED_OR_UPDATED_BY": "user@email.com",
			"TYPE": "NEW_RELIC_INSIGHTS",
			"ENABLED": 1,
			"STATE": {
				"ok": true,
				"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
				"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
			},
			"IS_ORG": 0,
			"DATA": {
				"ACCOUNT_ID": 2338053,
				"INSERT_KEY": "x-xx-xxxxxxxxxxxxxxxxxx",
			},
			"TYPE_NAME": "NEW_RELIC_INSIGHTS"
		}
	`
}
