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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestAlertChannelsGetSplunk(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetSplunkHec() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleSplunkAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetSplunkHec(intgGUID)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.IntgGuid)
		assert.Equal(t, "integration_name", response.Data.Name)
		assert.True(t, response.Data.State.Ok)
		assert.Equal(t, response.Data.Data.Channel, "channel-name")
		assert.Equal(t, response.Data.Data.HecToken, "AA111111-11AA-1AA1-11AA-11111AA1111A")
		assert.Equal(t, response.Data.Data.Host, "localhost")
		assert.Equal(t, response.Data.Data.Port, 80)
		assert.Equal(t, response.Data.Data.Ssl, true)
		assert.Equal(t, response.Data.Data.EventData.Index, "test-index")
		assert.Equal(t, response.Data.Data.EventData.Source, "test-source")
	}
}

func TestAlertChannelSplunkUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateSplunkHec() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "alert channel name is missing")
			assert.Contains(t, body, "SplunkHec", "wrong alert channel type")
			assert.Contains(t, body, "AA111111-11AA-1AA1-11AA-11111AA1111A", "missing splunk hec token")
			assert.Contains(t, body, "channel-name", "missing splunk channel name")
			assert.Contains(t, body, "localhost", "missing splunk host")
			assert.Contains(t, body, "80", "missing splunk port")
			assert.Contains(t, body, "true", "missing splunk ssl")
			assert.Contains(t, body, "test-index", "missing splunk index")
			assert.Contains(t, body, "test-source", "missing splunk source")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleSplunkAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	splunkAlertChan := api.NewAlertChannel("integration_name",
		api.SplunkHecAlertChannelType,
		api.SplunkHecDataV2{
			Channel:  "channel-name",
			HecToken: "AA111111-11AA-1AA1-11AA-11111AA1111A",
			Host:     "localhost",
			Port:     80,
			Ssl:      true,
			EventData: api.SplunkHecEventDataV2{
				Index:  "test-index",
				Source: "test-source",
			},
		},
	)
	assert.Equal(t, "integration_name", splunkAlertChan.Name, "Splunk alert channel name mismatch")
	assert.Equal(t, "SplunkHec", splunkAlertChan.Type, "a new Splunk alert channel should match its type")
	assert.Equal(t, 1, splunkAlertChan.Enabled, "a new Splunk alert channel should be enabled")
	splunkAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateSplunkHec(splunkAlertChan)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.IntgGuid)
		assert.True(t, response.Data.State.Ok)
		assert.Contains(t, response.Data.Data.Channel, "channel-name")
		assert.Contains(t, response.Data.Data.HecToken, "AA111111-11AA-1AA1-11AA-11111AA1111A")
		assert.Contains(t, response.Data.Data.Host, "localhost")
		assert.Equal(t, response.Data.Data.Port, 80)
		assert.Equal(t, response.Data.Data.Ssl, true)
		assert.Contains(t, response.Data.Data.EventData.Source, "test-source")
		assert.Contains(t, response.Data.Data.EventData.Index, "test-index")
	}
}

func TestMarshallAlertChannelLastUpdatedTime(t *testing.T) {
	var res api.SplunkHecAlertChannelResponseV2
	err := json.Unmarshal([]byte(generateAlertChannelResponse(singleSplunkAlertChannel("test"))), &res)
	if err != nil {
		log.Fatal("Unable to unmarshall splunk string")
	}
	jsonString, err := json.Marshal(res)
	if err != nil {
		log.Fatal("Unable to marshall splunk string")
	}

	assert.Equal(t, res.Data.State.LastUpdatedTime.ToTime().UnixNano()/int64(time.Millisecond), int64(1627895573122))
	assert.Contains(t, string(jsonString), "1627895573122")
}

func singleSplunkAlertChannel(id string) string {
	return `
{
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T18:10:40.745Z",
    "data": {
      "channel": "channel-name",
	  "hecToken": "AA111111-11AA-1AA1-11AA-11111AA1111A",
	  "eventData": {
		   "source": "test-source",
		   "index": "test-index"
	},
      "host": "localhost",
      "port": 80,
      "ssl": true
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
    "type": "SplunkHec"
}
  `
}
