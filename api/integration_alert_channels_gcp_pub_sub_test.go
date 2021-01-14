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

func TestIntegrationsNewGcpPubSubAlertChannel(t *testing.T) {
	subject := api.NewGcpPubSubAlertChannel("integration_name",
		api.GcpPubSubChannelData{
			ProjectID: "my-sample-project-191923",
			TopicID:   "mytopic",
			Credentials: api.GcpCredentials{
				ClientID:     "client_id",
				ClientEmail:  "foo@example.iam.gserviceaccount.com",
				PrivateKey:   "priv_key",
				PrivateKeyID: "p_key_id",
			},
		},
	)
	assert.Equal(t, api.GcpPubSubChannelIntegration.String(), subject.Type)
}

func TestIntegrationsCreateGcpPubSubAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateGcpPubSubAlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "GCP_PUBSUB", "wrong integration type")
			assert.Contains(t, body, "mytopic", "wrong topic id")
			assert.Contains(t, body, "my-sample-project-191923", "wrong project id")
			assert.Contains(t, body, "client_id", "wrong client id")
			assert.Contains(t, body, "foo@example.iam.gserviceaccount.com", "wrong client email")
			assert.Contains(t, body, "priv_key", "wrong private key")
			assert.Contains(t, body, "p_key_id", "wrong private key id")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, gcpPubSubChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewGcpPubSubAlertChannel("integration_name",
		api.GcpPubSubChannelData{
			ProjectID: "my-sample-project-191923",
			TopicID:   "mytopic",
			Credentials: api.GcpCredentials{
				ClientID:     "client_id",
				ClientEmail:  "foo@example.iam.gserviceaccount.com",
				PrivateKey:   "priv_key",
				PrivateKeyID: "p_key_id",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "GcpPubSubChannel integration name mismatch")
	assert.Equal(t, "GCP_PUBSUB", data.Type, "a new GcpPubSubChannel integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new GcpPubSubChannel integration should be enabled")

	response, err := c.Integrations.CreateGcpPubSubAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "my-sample-project-191923", resData.Data.ProjectID)
		assert.Equal(t, "mytopic", resData.Data.TopicID)
		assert.Equal(t, "client_id", resData.Data.Credentials.ClientID)
		assert.Equal(t, "foo@example.iam.gserviceaccount.com", resData.Data.Credentials.ClientEmail)
		assert.Equal(t, "priv_key", resData.Data.Credentials.PrivateKey)
		assert.Equal(t, "p_key_id", resData.Data.Credentials.PrivateKeyID)
	}
}

func TestIntegrationsGetGcpPubSubAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpPubSubAlertChannel should be a GET method")
		fmt.Fprintf(w, gcpPubSubChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetGcpPubSubAlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "my-sample-project-191923", resData.Data.ProjectID)
		assert.Equal(t, "mytopic", resData.Data.TopicID)
		assert.Equal(t, "client_id", resData.Data.Credentials.ClientID)
		assert.Equal(t, "foo@example.iam.gserviceaccount.com", resData.Data.Credentials.ClientEmail)
		assert.Equal(t, "priv_key", resData.Data.Credentials.PrivateKey)
		assert.Equal(t, "p_key_id", resData.Data.Credentials.PrivateKeyID)
	}
}

func TestIntegrationsUpdateGcpPubSubAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGcpPubSubAlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "GCP_PUBSUB", "wrong integration type")
			assert.Contains(t, body, "mytopic", "wrong topic id")
			assert.Contains(t, body, "my-sample-project-191923", "wrong project id")
			assert.Contains(t, body, "client_id", "wrong client id")
			assert.Contains(t, body, "foo@example.iam.gserviceaccount.com", "wrong client email")
			assert.Contains(t, body, "priv_key", "wrong private key")
			assert.Contains(t, body, "p_key_id", "wrong private key id")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, gcpPubSubChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewGcpPubSubAlertChannel("integration_name",
		api.GcpPubSubChannelData{
			ProjectID: "my-sample-project-191923",
			TopicID:   "mytopic",
			Credentials: api.GcpCredentials{
				ClientID:     "client_id",
				ClientEmail:  "foo@example.iam.gserviceaccount.com",
				PrivateKey:   "priv_key",
				PrivateKeyID: "p_key_id",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "GcpPubSubChannel integration name mismatch")
	assert.Equal(t, "GCP_PUBSUB", data.Type, "a new GcpPubSubChannel integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new GcpPubSubChannel integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateGcpPubSubAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListGcpPubSubAlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/GCP_PUBSUB",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListGcpPubSubAlertChannel should be a GET method")
			fmt.Fprintf(w, gcpPubSubChanMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListGcpPubSubAlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func gcpPubSubChannelIntegrationJsonResponse(intgGUID string) string {
	return `
{
  "data": [` + singleGcpPubSubChanIntegration(intgGUID) + `],
  "ok": true,
  "message": "SUCCESS"
}
`
}

func gcpPubSubChanMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleGcpPubSubChanIntegration(guid))
	}
	return `
{
"data": [` + strings.Join(integrations, ", ") + `],
"ok": true,
"message": "SUCCESS"
}
`
}

func singleGcpPubSubChanIntegration(id string) string {
	return `
		{
			"INTG_GUID": "` + id + `",
			"NAME": "integration_name",
			"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
			"CREATED_OR_UPDATED_BY": "user@email.com",
			"TYPE": "GCP_PUBSUB",
			"ENABLED": 1,
			"STATE": {
				"ok": true,
				"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
				"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
			},
			"IS_ORG": 0,
			"DATA": {
				"CREDENTIALS": {
					"CLIENT_ID": "client_id",
					"CLIENT_EMAIL": "foo@example.iam.gserviceaccount.com",
					"PRIVATE_KEY": "priv_key",
					"PRIVATE_KEY_ID": "p_key_id"
				},
				"PROJECT_ID": "my-sample-project-191923",
				"TOPIC_ID": "mytopic"
			},
			"TYPE_NAME": "GCP PUBSUB"
		}
	`
}
