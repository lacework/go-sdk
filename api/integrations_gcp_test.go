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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestIntegrationsCreateGcp(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	// WORKAROUND (@afiune) The backend is currently not triggering an initial
	// report automatically after creation of Cloud Account (CFG) Integrations,
	// we are implementing this trigger here until we implement it in the backend
	// with RAIN-13422
	fakeServer.MockAPI("external/runReport/integration/"+intgGUID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "RunReport should be a POST method")
	})
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateGcp should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "GCP_AT_SES", "wrong integration type")
			assert.Contains(t, body, "data_id", "missing data id")
			assert.Contains(t, body, "client_id", "missing client id")
			assert.Contains(t, body, "id_type", "missing id type")
			assert.Contains(t, body, "subscription_name", "missing subscription name")
			assert.Contains(t, body,
				"foo@example.iam.gserviceaccount.com", "missing client email",
			)
			assert.Contains(t, body, "priv_key", "wrong private key")
			assert.Contains(t, body, "p_key_id", "wrong private key id")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, gcpIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewGcpAuditLogIntegration("integration_name",
		api.GcpIntegrationData{
			ID:               "data_id",
			IDType:           "id_type",
			SubscriptionName: "subscription_name",
			Credentials: api.GcpCredentials{
				ClientID:     "client_id",
				ClientEmail:  "foo@example.iam.gserviceaccount.com",
				PrivateKey:   "priv_key",
				PrivateKeyID: "p_key_id",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "GCP integration name mismatch")
	assert.Equal(t, "GCP_AT_SES", data.Type, "a new GCP integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new GCP integration should be enabled")

	response, err := c.Integrations.CreateGcp(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "data_id", resData.Data.ID)
		assert.Equal(t, "PROJECT", resData.Data.IDType)
		assert.Equal(t,
			"foo@example.iam.gserviceaccount.com",
			resData.Data.Credentials.ClientEmail,
		)
	}
}

func TestIntegrationsGetGcp(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcp should be a GET method")
		fmt.Fprintf(w, gcpIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetGcp(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "data_id", resData.Data.ID)
		assert.Equal(t, "PROJECT", resData.Data.IDType)
		assert.Equal(t,
			"foo@example.iam.gserviceaccount.com",
			resData.Data.Credentials.ClientEmail,
		)
	}
}

func TestIntegrationsUpdateGcp(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGcp should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "GCP_CFG", "wrong integration type")
			assert.Contains(t, body, "data_id", "missing data id")
			assert.Contains(t, body, "client_id", "missing client id")
			assert.Contains(t, body, "id_type", "missing id type")
			assert.NotContains(t, body, "subscription_name", "sholuld not have subscription_name")
			assert.Contains(t, body,
				"foo@example.iam.gserviceaccount.com", "missing client email",
			)
			assert.Contains(t, body, "priv_key", "wrong private key")
			assert.Contains(t, body, "p_key_id", "wrong private key id")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, gcpIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewGcpIntegration("integration_name",
		api.GcpCfgIntegration,
		api.GcpIntegrationData{
			ID:     "data_id",
			IDType: "id_type",
			Credentials: api.GcpCredentials{
				ClientID:     "client_id",
				ClientEmail:  "foo@example.iam.gserviceaccount.com",
				PrivateKey:   "priv_key",
				PrivateKeyID: "p_key_id",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "GCP integration name mismatch")
	assert.Equal(t, "GCP_CFG", data.Type, "a new GCP integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new GCP integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateGcp(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsDeleteGcp(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "DeleteGcp should be a DELETE method")
		fmt.Fprintf(w, gcpIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.DeleteGcp(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListGcpCfg(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/GCP_CFG",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListGcpCfg should be a GET method")
			fmt.Fprintf(w, gcpMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListGcpCfg()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func gcpIntegrationJsonResponse(intgGUID string) string {
	return `
		{
			"data": [` + singleGcpIntegration(intgGUID) + `],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func gcpMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleGcpIntegration(guid))
	}
	return `
		{
			"data": [` + strings.Join(integrations, ", ") + `],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func singleGcpIntegration(id string) string {
	return `
		{
			"INTG_GUID": "` + id + `",
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
					"CLIENT_ID": "client_id",
					"CLIENT_EMAIL": "foo@example.iam.gserviceaccount.com",
					"PRIVATE_KEY": "priv_key",
					"PRIVATE_KEY_ID": "p_key_id"
				},
				"ID_TYPE": "PROJECT",
				"ID": "data_id"
			},
			"TYPE_NAME": "GCP Compliance"
		}
	`
}
