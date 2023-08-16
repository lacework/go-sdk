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

package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestAlertChannelsGetSlackChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetSlackChannel() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleSlackChannelAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetSlackChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.SlackUrl, "https://hooks.slack.com/services/abc/foo/secret123")
}

func TestAlertChannelsSlackChannelUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateSlackChannel() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "SlackChannel", "wrong cloud account type")
			assert.Contains(t, body, "https://hooks.slack.com/services/abc/foo/secret123", "missing slack url")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleSlackChannelAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	emailAlertChan := api.NewAlertChannel("integration_name",
		api.SlackChannelAlertChannelType,
		api.SlackChannelDataV2{
			SlackUrl: "https://hooks.slack.com/services/abc/foo/secret123",
		},
	)
	assert.Equal(t, "integration_name", emailAlertChan.Name, "SlackChannel cloud account name mismatch")
	assert.Equal(t, "SlackChannel", emailAlertChan.Type, "a new SlackChannel cloud account should match its type")
	assert.Equal(t, 1, emailAlertChan.Enabled, "a new SlackChannel cloud account should be enabled")
	emailAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateSlackChannel(emailAlertChan)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.SlackUrl, "https://hooks.slack.com/services/abc/foo/secret123")
}

func singleSlackChannelAlertChannel(id string) string {
	return `
{
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T18:10:40.745Z",
    "data": {
      "slackUrl": "https://hooks.slack.com/services/abc/foo/secret123"
    },
    "enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "state": {
      "details": {},
      "lastSuccessfulTime": 1627895573122,
      "lastUpdatedTime": 1627895573122,
      "ok": true
    },
    "type": "SlackChannel"
}
  `
}
