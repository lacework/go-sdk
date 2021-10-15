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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestAlertChannelsGetNewRelic(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetNewRelicInsights() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleNewRelicAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetNewRelicInsights(intgGUID)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.IntgGuid)
		assert.Equal(t, "integration_name", response.Data.Name)
		assert.True(t, response.Data.State.Ok)
		assert.Equal(t, response.Data.Data.AccountID, 2338053)
		assert.Equal(t, response.Data.Data.InsertKey, "x-xx-xxxxxxxxxxxxxxxxxx")
	}
}

func TestAlertChannelNewRelicUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateNewRelicInsights() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "alert channel name is missing")
			assert.Contains(t, body, "NewRelicInsights", "wrong alert channel type")
			assert.Contains(t, body, "\"accountId\":2338053", "missing account id")
			assert.Contains(t, body, "\"insertKey\":\"x-xx-xxxxxxxxxxxxxxxxxx\"", "missing insert key")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleNewRelicAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	newRelicAlertChan := api.NewAlertChannel("integration_name",
		api.NewRelicInsightsAlertChannelType,
		api.NewRelicInsightsDataV2{
			AccountID: 2338053,
			InsertKey: "x-xx-xxxxxxxxxxxxxxxxxx",
		},
	)
	assert.Equal(t, "integration_name", newRelicAlertChan.Name, "NewRelic alert channel name mismatch")
	assert.Equal(t, "NewRelicInsights", newRelicAlertChan.Type, "a new NewRelic alert channel should match its type")
	assert.Equal(t, 1, newRelicAlertChan.Enabled, "a new NewRelic alert channel should be enabled")
	newRelicAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateNewRelicInsights(newRelicAlertChan)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.IntgGuid)
		assert.True(t, response.Data.State.Ok)
		assert.Equal(t, response.Data.Data.InsertKey, "x-xx-xxxxxxxxxxxxxxxxxx")
		assert.Equal(t, response.Data.Data.AccountID, 2338053)
	}
}

func singleNewRelicAlertChannel(id string) string {
	return fmt.Sprintf(`
{
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T18:10:40.745Z",
    "data": {
		"accountId": 2338053,
		"insertKey": "x-xx-xxxxxxxxxxxxxxxxxx"
    },
    "enabled": 1,
    "intgGuid": %q,
    "isOrg": 0,
    "name": "integration_name",
    "state": {
      "details": {},
      "lastSuccessfulTime": 1627895573122,
      "lastUpdatedTime": 1627895573122,
      "ok": true
    },
    "type": "NewRelicInsights"
}
  `, id)
}
