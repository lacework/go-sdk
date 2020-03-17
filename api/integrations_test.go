//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestGetIntegrations(t *testing.T) {
	// TODO @afiune implement a mocked Lacework API server
}

func TestCreateGCPConfigIntegration(t *testing.T) {
	intgGUID := "12345"

	fakeAPI := NewLaceworkServer()
	fakeAPI.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, createGCPCFGJson(intgGUID))
	})
	defer fakeAPI.Close()

	c, err := api.NewClient("test", api.WithToken("xxxxxx"), api.WithURL(fakeAPI.URL()))
	assert.Nil(t, err)

	data := api.NewGCPIntegrationData("integration_name", api.GcpProject)
	assert.Equal(t, "GCP_CFG", data.Type, "a new GCP integration should match its type")
	data.Data.ID = "xxxxxxxxxx"
	data.Data.Credentials.ClientId = "xxxxxxxxx"
	data.Data.Credentials.ClientEmail = "xxxxxx@xxxxx.iam.gserviceaccount.com"
	data.Data.Credentials.PrivateKeyId = "xxxxxxxxxxxxxxxx"

	response, err := c.CreateGCPConfigIntegration(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestGetGCPConfigIntegration(t *testing.T) {
	intgGUID := "12345"
	apiPath := fmt.Sprintf("external/integrations/%s", intgGUID)

	fakeAPI := NewLaceworkServer()
	fakeAPI.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, createGCPCFGJson(intgGUID))
	})
	defer fakeAPI.Close()

	c, err := api.NewClient("test", api.WithToken("xxxxxx"), api.WithURL(fakeAPI.URL()))
	assert.Nil(t, err)

	response, err := c.GetGCPConfigIntegration(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestUpdateGCPConfigIntegration(t *testing.T) {
	intgGUID := "12345"
	apiPath := fmt.Sprintf("external/integrations/%s", intgGUID)

	fakeAPI := NewLaceworkServer()
	fakeAPI.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, createGCPCFGJson(intgGUID))
	})
	defer fakeAPI.Close()

	c, err := api.NewClient("test", api.WithToken("xxxxxx"), api.WithURL(fakeAPI.URL()))
	assert.Nil(t, err)

	data := api.NewGCPIntegrationData("integration_name", api.GcpProject)
	assert.Equal(t, "GCP_CFG", data.Type, "a new GCP integration should match its type")
	data.IntgGuid = intgGUID
	data.Data.ID = "xxxxxxxxxx"
	data.Data.Credentials.ClientId = "xxxxxxxxx"
	data.Data.Credentials.ClientEmail = "xxxxxx@xxxxx.iam.gserviceaccount.com"
	data.Data.Credentials.PrivateKeyId = "xxxxxxxxxxxxxxxx"

	response, err := c.UpdateGCPConfigIntegration(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestDeleteGCPConfigIntegration(t *testing.T) {
	intgGUID := "12345"
	apiPath := fmt.Sprintf("external/integrations/%s", intgGUID)

	fakeAPI := NewLaceworkServer()
	fakeAPI.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, createGCPCFGJson(intgGUID))
	})
	defer fakeAPI.Close()

	c, err := api.NewClient("test", api.WithToken("xxxxxx"), api.WithURL(fakeAPI.URL()))
	assert.Nil(t, err)

	response, err := c.DeleteGCPConfigIntegration(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func createGCPCFGJson(intgGUID string) string {
	return `
		{
			"data": [
				{
					"INTG_GUID": "` + intgGUID + `",
					"NAME": "integration_name",
					"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
					"CREATED_OR_UPDATED_BY": "user@email.com",
					"TYPE": "GCP_CFG",
					"ENABLED": 1,
					"STATE": {
						"ok": true,
						"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
						"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
					},
					"IS_ORG": 0,
					"DATA": {
						"CREDENTIALS": {
							"CLIENT_ID": "xxxxxxxxx",
							"CLIENT_EMAIL": "xxxxxx@xxxxx.iam.gserviceaccount.com",
							"PRIVATE_KEY_ID": "xxxxxxxxxxxxxxxx"
						},
						"ID_TYPE": "PROJECT",
						"ID": "xxxxxxxxxx"
					},
					"TYPE_NAME": "GCP Compliance"
				}
			],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}
