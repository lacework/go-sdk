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

func TestIntegrationsNewDatadogAlertChannel(t *testing.T) {
	subject := api.NewDatadogAlertChannel("integration_name",
		api.DatadogChannelData{
			DatadogSite: "eu",
			DatadogType: "Events Summary",
			ApiKey:      "datadog-key",
		},
	)
	assert.Equal(t, api.DatadogChannelIntegration.String(), subject.Type)
}

func TestIntegrationsCreateDatadogAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateDatadogAlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "DATADOG", "wrong integration type")
			assert.Contains(t, body, "eu", "wrong datadog site")
			assert.Contains(t, body, "Events Summary", "wrong datadog type")
			assert.Contains(t, body, "datadog-key", "wrong datadog api key")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, datadogChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewDatadogAlertChannel("integration_name",
		api.DatadogChannelData{
			DatadogSite: "eu",
			DatadogType: "Events Summary",
			ApiKey:      "datadog-key",
		},
	)

	assert.Equal(t, "integration_name", data.Name, "Datadog integration name mismatch")
	assert.Equal(t, "DATADOG", data.Type, "a new Datadog integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new Datadog integration should be enabled")

	response, err := c.Integrations.CreateDatadogAlertChannel(data)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "eu", resData.Data.DatadogSite)
		assert.Equal(t, "Events Summary", resData.Data.DatadogType)
		assert.Equal(t, "datadog-key", resData.Data.ApiKey)
	}
}

func TestIntegrationsGetDatadogAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetDatadogAlertChannel should be a GET method")
		fmt.Fprintf(w, datadogChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetDatadogAlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "eu", resData.Data.DatadogSite)
		assert.Equal(t, "Events Summary", resData.Data.DatadogType)
		assert.Equal(t, "datadog-key", resData.Data.ApiKey)
	}
}

func TestIntegrationsUpdateDatadogAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateDatadogAlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "DATADOG", "wrong integration type")
			assert.Contains(t, body, "eu", "wrong datadog site")
			assert.Contains(t, body, "Events Summary", "wrong datadog type")
			assert.Contains(t, body, "datadog-key", "wrong datadog api key")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, datadogChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewDatadogAlertChannel("integration_name",
		api.DatadogChannelData{
			DatadogSite: "eu",
			DatadogType: "Events Summary",
			ApiKey:      "datadog-key",
		},
	)
	assert.Equal(t, "integration_name", data.Name, "Datadog integration name mismatch")
	assert.Equal(t, "DATADOG", data.Type, "a new Datadog integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new Datadog integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateDatadogAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListDatadogAlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/DATADOG",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListDatadogAlertChannel should be a GET method")
			fmt.Fprintf(w, datadogChanMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListDatadogAlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func datadogChannelIntegrationJsonResponse(intgGUID string) string {
	return `
	{
		"data": [` + singleDatadogChanIntegration(intgGUID) + `],
		"ok": true,
		"message": "SUCCESS"
	}
	`
}

func datadogChanMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleDatadogChanIntegration(guid))
	}
	return `
	{
		"data": [` + strings.Join(integrations, ", ") + `],
		"ok": true,
		"message": "SUCCESS"
	}
	`
}

func singleDatadogChanIntegration(id string) string {
	return `
		{
			"INTG_GUID": "` + id + `",
			"NAME": "integration_name",
			"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
			"CREATED_OR_UPDATED_BY": "user@email.com",
			"TYPE": "DATADOG",
			"ENABLED": 1,
			"STATE": {
				"ok": true,
				"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
				"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
			},
			"IS_ORG": 0,
			"DATA": {
				"DATADOG_SITE": "eu",
				"DATADOG_TYPE": "Events Summary",
				"API_KEY": "datadog-key"
			},
			"TYPE_NAME": "DATADOG"
		}
	`
}
