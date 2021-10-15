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

func TestAlertChannelsGetIbmQRadar(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetIbmQRadar() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleIbmQRadarAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetIbmQRadar(intgGUID)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.IntgGuid)
		assert.Equal(t, "integration_name", response.Data.Name)
		assert.True(t, response.Data.State.Ok)
		assert.Equal(t, response.Data.Data.HostURL, "https://qradar-lacework.com")
		assert.Equal(t, response.Data.Data.HostPort, 443)
		assert.Equal(t, response.Data.Data.QRadarCommType, api.QRadarCommHttps)
	}
}

func TestAlertChannelIbmQRadarUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateIbmQRadar() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "alert channel name is missing")
			assert.Contains(t, body, "IbmQradar", "wrong alert channel type")
			assert.Contains(t, body, "https://qradar-lacework.com", "missing host url")
			assert.Contains(t, body, "443", "missing host port")
			assert.Contains(t, body, "HTTPS", "missing comm type")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleIbmQRadarAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	qradarAlertChan := api.NewAlertChannel("integration_name",
		api.IbmQRadarAlertChannelType,
		api.IbmQRadarDataV2{
			HostURL:        "https://qradar-lacework.com",
			HostPort:       443,
			QRadarCommType: "HTTPS",
		},
	)
	assert.Equal(t, "integration_name", qradarAlertChan.Name, "IbmQRadar alert channel name mismatch")
	assert.Equal(t, "IbmQradar", qradarAlertChan.Type, "a new IbmQRadar alert channel should match its type")
	assert.Equal(t, 1, qradarAlertChan.Enabled, "a new IbmQRadar alert channel should be enabled")
	qradarAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateIbmQRadar(qradarAlertChan)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.IntgGuid)
		assert.True(t, response.Data.State.Ok)
		assert.Equal(t, response.Data.Data.QRadarCommType, api.QRadarCommHttps)
		assert.Equal(t, response.Data.Data.HostPort, 443)
		assert.Equal(t, response.Data.Data.HostURL, "https://qradar-lacework.com")
	}
}

func singleIbmQRadarAlertChannel(id string) string {
	return `
{
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T18:10:40.745Z",
    "data": {
	  "qradarHostUrl": "https://qradar-lacework.com",
	  "qradarHostPort": 443,
	  "qradarCommType": "HTTPS"
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
    "type": "IbmQradar"
}
  `
}
