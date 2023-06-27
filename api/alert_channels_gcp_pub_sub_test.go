//
// Author:: Vatasha White (<vatasha.white@lacework.net>)
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

func TestAlertChannelsService_GetGcpPubSub(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpPubSub() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleGcpPubSubAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetGcpPubSub(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "vatasha.white@lacework.net", response.Data.Data.Credentials.ClientEmail)
	assert.Equal(t, "fake-client-id", response.Data.Data.Credentials.ClientID)
	assert.Equal(t, "fake-private-key", response.Data.Data.Credentials.PrivateKey)
	assert.Equal(t, "fake-private-key-id", response.Data.Data.Credentials.PrivateKeyID)
	assert.Equal(t, "fake-project-id", response.Data.Data.ProjectID)
	assert.Equal(t, "fake-topic-id", response.Data.Data.TopicID)
	assert.Equal(t, "Events", response.Data.Data.IssueGrouping)
}

func TestAlertChannelsService_UpdateGcpPubSub(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGcpPubSub() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "IntgGUID missing")
			assert.Contains(t, body, "name\":\"integration_name", "cloud account name is missing")
			assert.Contains(t, body, "type\":\"GcpPubsub", "wrong cloud account type")
			assert.Contains(t, body, "clientEmail\":\"vatasha.white@lacework.net", "missing client email")
			assert.Contains(t, body, "clientId\":\"fake-client-id", "missing client id")
			assert.Contains(t, body, "privateKey\":\"fake-private-key", "missing private key")
			assert.Contains(t, body, "privateKeyId\":\"fake-private-key-id", "missing private key id")
			assert.Contains(t, body, "projectId\":\"fake-project-id", "missing project id")
			assert.Contains(t, body, "topicId\":\"fake-topic-id", "missing topic id")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleGcpPubSubAlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	gcpPubSubAlertChan := api.NewAlertChannel("integration_name",
		api.GcpPubSubAlertChannelType,
		api.GcpPubSubDataV2{
			Credentials: api.GcpPubSubCredentials{
				ClientEmail:  "vatasha.white@lacework.net",
				ClientID:     "fake-client-id",
				PrivateKey:   "fake-private-key",
				PrivateKeyID: "fake-private-key-id",
			},
			IssueGrouping: "Events",
			ProjectID:     "fake-project-id",
			TopicID:       "fake-topic-id",
		},
	)
	assert.Equal(t, "integration_name", gcpPubSubAlertChan.Name)
	assert.Equal(t, "GcpPubsub", gcpPubSubAlertChan.Type)
	assert.Equal(t, 1, gcpPubSubAlertChan.Enabled)
	gcpPubSubAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateGcpPubSub(gcpPubSubAlertChan)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.IntgGuid)
		assert.True(t, response.Data.State.Ok)
		assert.Equal(t, "integration_name", response.Data.Name)
		assert.Equal(t, "vatasha.white@lacework.net", response.Data.Data.Credentials.ClientEmail)
	}
}

func singleGcpPubSubAlertChannel(id string) string {
	return fmt.Sprintf(`
	{
		"createdOrUpdatedBy": "vatasha.white@lacework.net",
		"createdOrUpdatedTime": "2021-09-29T117:55:47.277316",
		"data": {
			"credentials": {
				"clientEmail": "vatasha.white@lacework.net",
				"clientId": "fake-client-id",
				"privateKey": "fake-private-key",
				"privateKeyId": "fake-private-key-id"
			},
			"issueGrouping": "Events",
			"projectId": "fake-project-id",
			"topicId": "fake-topic-id"
		},
		"enabled": 1,
		"intgGuid": %q,
		"isOrg": 0,
		"name": "integration_name",
		"state": {
		"details": {},
		"lastSuccessfulTime": 1632932665892,
			"lastUpdatedTime": 1632932665892,
			"ok": true
	},
		"type": "GcpPubsub"
	}
	`, id)
}
