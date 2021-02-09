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

func TestIntegrationsNewQRadarAlertChannel(t *testing.T) {
	subject := api.NewQRadarAlertChannel("integration_name",
		api.QRadarChannelData{
			HostURL:           "https://qradar-lacework.com",
			HostPort:          8080,
			CommunicationType: api.QRadarCommHttps,
		},
	)
	assert.Equal(t, api.QRadarChannelIntegration.String(), subject.Type)
}

func TestIntegrationsCreateQRadarAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateQRadarAlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "IBM_QRADAR", "wrong integration type")
			assert.Contains(t, body, "HTTPS", "wrong communication type")
			assert.Contains(t, body, "https://qradar-lacework.com", "wrong host url")
			assert.Contains(t, body, "8080", "wrong port")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, qradarChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewQRadarAlertChannel("integration_name",
		api.QRadarChannelData{
			HostURL:           "https://qradar-lacework.com",
			HostPort:          8080,
			CommunicationType: api.QRadarCommHttps,
		},
	)
	assert.Equal(t, "integration_name", data.Name, "QRadarChannel integration name mismatch")
	assert.Equal(t, "IBM_QRADAR", data.Type, "a new QRadarChannel integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new QRadarChannel integration should be enabled")

	response, err := c.Integrations.CreateQRadarAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "https://qradar-lacework.com", resData.Data.HostURL)
		assert.Equal(t, api.QRadarCommHttps, resData.Data.CommunicationType)
		assert.Equal(t, 8080, resData.Data.HostPort)
	}
}

func TestIntegrationsGetQRadarAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetQRadarAlertChannel should be a GET method")
		fmt.Fprintf(w, qradarChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetQRadarAlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "https://qradar-lacework.com", resData.Data.HostURL)
		assert.Equal(t, api.QRadarCommHttps, resData.Data.CommunicationType)
		assert.Equal(t, 8080, resData.Data.HostPort)
	}
}

func TestIntegrationsUpdateQRadarAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateQRadarAlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "IBM_QRADAR", "wrong integration type")
			assert.Contains(t, body, "HTTPS", "wrong communication type")
			assert.Contains(t, body, "https://qradar-lacework.com", "wrong host url")
			assert.Contains(t, body, "8080", "wrong port")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, qradarChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewQRadarAlertChannel("integration_name",
		api.QRadarChannelData{
			HostURL:           "https://qradar-lacework.com",
			HostPort:          8080,
			CommunicationType: api.QRadarCommHttps,
		},
	)
	assert.Equal(t, "integration_name", data.Name, "QRadarChannel integration name mismatch")
	assert.Equal(t, "IBM_QRADAR", data.Type, "a new QRadarChannel integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new QRadarChannel integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateQRadarAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListQRadarAlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/IBM_QRADAR",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListQRadarAlertChannel should be a GET method")
			fmt.Fprintf(w, qradarChanMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListQRadarAlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func qradarChannelIntegrationJsonResponse(intgGUID string) string {
	return `
{
  "data": [` + singleQRadarChanIntegration(intgGUID) + `],
  "ok": true,
  "message": "SUCCESS"
}
`
}

func qradarChanMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleQRadarChanIntegration(guid))
	}
	return `
{
"data": [` + strings.Join(integrations, ", ") + `],
"ok": true,
"message": "SUCCESS"
}
`
}

func singleQRadarChanIntegration(id string) string {
	return `
		{
			"INTG_GUID": "` + id + `",
			"NAME": "integration_name",
			"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
			"CREATED_OR_UPDATED_BY": "user@email.com",
			"TYPE": "IBM_QRADAR",
			"ENABLED": 1,
			"STATE": {
				"ok": true,
				"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
				"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
			},
			"IS_ORG": 0,
			"DATA": {
				"QRADAR_HOST_URL": "https://qradar-lacework.com",
				"QRADAR_COMM_TYPE": "HTTPS",
				"QRADAR_HOST_PORT": 8080,
			},
			"TYPE_NAME": "IBM_QRADAR"
		}
	`
}
