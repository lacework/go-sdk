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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestIntegrationsNewSplunkAlertChannel(t *testing.T) {
	subject := api.NewSplunkAlertChannel("integration_name",
		api.SplunkChannelData{
			Channel:  "channel-name",
			HecToken: "AA111111-11AA-1AA1-11AA-11111AA1111A",
			Host:     "localhost",
			Port:     80,
			Ssl:      false,
			EventData: api.SplunkEventData{
				Index:  "test-index",
				Source: "test-source",
			},
		},
	)
	assert.Equal(t, api.SplunkIntegration.String(), subject.Type)
}

func TestIntegrationsCreateSplunkAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateSplunkAlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "SPLUNK_HEC", "wrong integration type")
			assert.Contains(t, body, "channel-name", "wrong splunk channel name type")
			assert.Contains(t, body, "AA111111-11AA-1AA1-11AA-11111AA1111A", "wrong splunk hec token")
			assert.Contains(t, body, "localhost", "wrong splunk host")
			assert.Contains(t, body, "80", "wrong splunk port")
			assert.Contains(t, body, "test-index", "wrong splunk index")
			assert.Contains(t, body, "test-source", "wrong splunk source")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, splunkChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewSplunkAlertChannel("integration_name",
		api.SplunkChannelData{
			Channel:  "channel-name",
			HecToken: "AA111111-11AA-1AA1-11AA-11111AA1111A",
			Host:     "localhost",
			Port:     80,
			Ssl:      false,
			EventData: api.SplunkEventData{
				Index:  "test-index",
				Source: "test-source",
			},
		},
	)

	assert.Equal(t, "integration_name", data.Name, "Splunk integration name mismatch")
	assert.Equal(t, "SPLUNK_HEC", data.Type, "a new Splunk integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new Splunk integration should be enabled")

	response, err := c.Integrations.CreateSplunkAlertChannel(data)
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(b))

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "AA111111-11AA-1AA1-11AA-11111AA1111A", resData.Data.HecToken)
		assert.Equal(t, "channel-name", resData.Data.Channel)
		assert.Equal(t, "localhost", resData.Data.Host)
		assert.Equal(t, "test-index", resData.Data.EventData.Index)
		assert.Equal(t, "test-source", resData.Data.EventData.Source)
		assert.Equal(t, 80, resData.Data.Port)
	}
}

func TestIntegrationsGetSplunkAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetSplunkAlertChannel should be a GET method")
		fmt.Fprintf(w, splunkChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetSplunkAlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "AA111111-11AA-1AA1-11AA-11111AA1111A", resData.Data.HecToken)
		assert.Equal(t, "channel-name", resData.Data.Channel)
		assert.Equal(t, "localhost", resData.Data.Host)
		assert.Equal(t, "test-index", resData.Data.EventData.Index)
		assert.Equal(t, "test-source", resData.Data.EventData.Source)
		assert.Equal(t, 80, resData.Data.Port)
	}
}

func TestIntegrationsUpdateSplunkAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateSplunkAlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "SPLUNK_HEC", "wrong integration type")
			assert.Contains(t, body, "channel-name", "wrong splunk channel name type")
			assert.Contains(t, body, "AA111111-11AA-1AA1-11AA-11111AA1111A", "wrong splunk hec token")
			assert.Contains(t, body, "localhost", "wrong splunk host")
			assert.Contains(t, body, "80", "wrong splunk port")
			assert.Contains(t, body, "test-index", "wrong splunk index")
			assert.Contains(t, body, "test-source", "wrong splunk source")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, splunkChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewSplunkAlertChannel("integration_name",
		api.SplunkChannelData{
			Channel:  "channel-name",
			HecToken: "AA111111-11AA-1AA1-11AA-11111AA1111A",
			Host:     "localhost",
			Port:     80,
			Ssl:      false,
			EventData: api.SplunkEventData{
				Index:  "test-index",
				Source: "test-source",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "Splunk integration name mismatch")
	assert.Equal(t, "SPLUNK_HEC", data.Type, "a new Splunk integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new Splunk integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateSplunkAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListSplunkAlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/SPLUNK_HEC",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListSplunkAlertChannel should be a GET method")
			fmt.Fprintf(w, splunkChanMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListSplunkAlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func splunkChannelIntegrationJsonResponse(intgGUID string) string {
	return `
{
  "data": [` + singleSplunkChanIntegration(intgGUID) + `],
  "ok": true,
  "message": "SUCCESS"
}
`
}

func splunkChanMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleSplunkChanIntegration(guid))
	}
	return `
{
"data": [` + strings.Join(integrations, ", ") + `],
"ok": true,
"message": "SUCCESS"
}
`
}

func singleSplunkChanIntegration(id string) string {
	return `
{
  "INTG_GUID": "` + id + `",
  "CREATED_OR_UPDATED_BY": "user@email.com",
  "CREATED_OR_UPDATED_TIME": "2020-Jul-16 19:59:22 UTC",
  "DATA": {
    "ISSUE_GROUPING": "Events",
    "CHANNEL": "channel-name",
	"HEC_TOKEN": "AA111111-11AA-1AA1-11AA-11111AA1111A",
	"EVENT_DATA": {
		"SOURCE": "test-source",
		"INDEX": "test-index"
	},
    "HOST": "localhost",
    "PORT": 80,
    "SSL": true
  },
  "ENABLED": 1,
  "IS_ORG": 0,
  "NAME": "integration_name",
  "STATE": {
    "lastSuccessfulTime": "2020-Jul-16 18:26:54 UTC",
    "lastUpdatedTime": "2020-Jul-16 18:26:54 UTC",
    "ok": true
  },
  "TYPE": "SPLUNK_HEC",
  "TYPE_NAME": "SPLUNK_HEC"
}
`
}
