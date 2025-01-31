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

	"github.com/lacework/go-sdk/v2/api"
	"github.com/lacework/go-sdk/v2/internal/intgguid"
	"github.com/lacework/go-sdk/v2/internal/lacework"
)

func TestAlertChannelsGetMicrosoftTeams(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetMicrosoftTeams() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleMicrosoftTeamsAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetMicrosoftTeams(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.TeamsURL, "https://example.webhook.office.com/webwook123/abc")
}

func TestAlertChannelsMicrosoftTeamsUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateMicrosoftTeams() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "MicrosoftTeams", "wrong cloud account type")
			assert.Contains(t, body, "https://example.webhook.office.com/webwook123/abc", "missing slack url")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleMicrosoftTeamsAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	emailAlertChan := api.NewAlertChannel("integration_name",
		api.MicrosoftTeamsAlertChannelType,
		api.MicrosoftTeamsData{
			TeamsURL: "https://example.webhook.office.com/webwook123/abc",
		},
	)
	assert.Equal(t, "integration_name", emailAlertChan.Name, "MicrosoftTeams cloud account name mismatch")
	assert.Equal(t, "MicrosoftTeams", emailAlertChan.Type, "a new MicrosoftTeams cloud account should match its type")
	assert.Equal(t, 1, emailAlertChan.Enabled, "a new MicrosoftTeams cloud account should be enabled")
	emailAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateMicrosoftTeams(emailAlertChan)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.TeamsURL, "https://example.webhook.office.com/webwook123/abc")
}

func singleMicrosoftTeamsAlertChannel(id string) string {
	return `
{
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T18:10:40.745Z",
    "data": {
      "teamsUrl": "https://example.webhook.office.com/webwook123/abc"
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
    "type": "MicrosoftTeams"
}
  `
}
