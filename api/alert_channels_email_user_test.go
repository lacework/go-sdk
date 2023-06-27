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

func TestAlertChannelsGetEmailUser(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetEmailUser() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleEmailUserAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetEmailUser(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.ChannelProps.Recipients, "foo@lacework.net")
	assert.Contains(t, response.Data.Data.ChannelProps.Recipients, "bar@lacework.net")
}

func TestAlertChannelsEmailUserUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateEmailUser() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "EmailUser", "wrong cloud account type")
			assert.Contains(t, body, "foo@example.com", "missing recipient")
			assert.Contains(t, body, "bar@example.com", "missing recipient")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleEmailUserAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	emailAlertChan := api.NewAlertChannel("integration_name",
		api.EmailUserAlertChannelType,
		api.EmailUserData{
			ChannelProps: api.EmailUserChannelProps{
				Recipients: []string{"foo@example.com", "bar@example.com"},
			},
		},
	)
	assert.Equal(t, "integration_name", emailAlertChan.Name, "EmailUser cloud account name mismatch")
	assert.Equal(t, "EmailUser", emailAlertChan.Type, "a new EmailUser cloud account should match its type")
	assert.Equal(t, 1, emailAlertChan.Enabled, "a new EmailUser cloud account should be enabled")
	emailAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateEmailUser(emailAlertChan)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.ChannelProps.Recipients, "foo@lacework.net")
	assert.Contains(t, response.Data.Data.ChannelProps.Recipients, "bar@lacework.net")
}

func singleEmailUserAlertChannel(id string) string {
	return `
{
    "createdOrUpdatedBy": "Lacework",
    "createdOrUpdatedTime": "2020-02-08T00:05:35.996Z",
    "data": {
      "channelProps": {
        "recipients": ["foo@lacework.net","bar@lacework.net"]
      },
      "notificationTypes": {
        "properties": {
          "agentEvents": true,
          "awsCisS3": false,
          "awsCloudtrailEvents": true,
          "awsComplianceEvents": true,
          "azureCis": true,
          "gcpCis": true,
          "noEvents": false,
          "time": "1200 (GMT)",
          "trendReport": true
        }
      }
    },
    "enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "props": {
      "isDefault": 1
    },
    "state": {
      "details": {
        "message": "Unable to send email to channel with intg guid: foo"
      },
      "lastSuccessfulTime": 1617380198077,
      "lastUpdatedTime": 1617380198077,
      "ok": true
    },
    "type": "EmailUser"
}
  `
}

// Workaround from APIv2
// Bug: https://lacework.atlassian.net/browse/RAIN-20070
func TestAlertChannelsGetEmailUserBug(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetEmailUser() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(bugGetEmailUserAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetEmailUser(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.ChannelProps.Recipients, "foo@lacework.net")
	assert.Contains(t, response.Data.Data.ChannelProps.Recipients, "bar@lacework.net")
}

func bugGetEmailUserAlertChannel(id string) string {
	return `
{
    "createdOrUpdatedBy": "Lacework",
    "createdOrUpdatedTime": "2020-02-08T00:05:35.996Z",
    "data": {
      "channelProps": {
        "recipients": "foo@lacework.net,bar@lacework.net"
      },
      "notificationTypes": {
        "properties": {
          "agentEvents": true,
          "awsCisS3": false,
          "awsCloudtrailEvents": true,
          "awsComplianceEvents": true,
          "azureCis": true,
          "gcpCis": true,
          "noEvents": false,
          "time": "1200 (GMT)",
          "trendReport": true
        }
      }
    },
    "enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "props": {
      "isDefault": 1
    },
    "state": {
      "details": {
        "message": "Unable to send email to channel with intg guid: foo"
      },
      "lastSuccessfulTime": 1617380198077,
      "lastUpdatedTime": 1617380198077,
      "ok": true
    },
    "type": "EmailUser"
}
  `
}
