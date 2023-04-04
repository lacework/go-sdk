//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestCloudAccountsGcpAtSesGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpAtSes() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleGcpAtCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetGcpAtSes(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "test-gcp", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "test@project.iam.gserviceaccount.com", response.Data.Data.Credentials.ClientEmail)
	assert.Equal(t, "123456789", response.Data.Data.Credentials.ClientID)
	assert.Equal(t, "", response.Data.Data.Credentials.PrivateKeyID)
	assert.Equal(t, "", response.Data.Data.Credentials.PrivateKey)
	assert.Equal(t, "GcpAtSes", response.Data.Type)
	assert.Equal(t, "test-project-123", response.Data.Data.ID)
	assert.Equal(t, "PROJECT", response.Data.Data.IDType)
	assert.Equal(t, "projects/test-project/subscriptions/test", response.Data.Data.SubscriptionName)
}

func TestCloudAccountsGcpAtSesUpdateEmptyPrivateKeyAndPrivateKeyID(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGcpAtSes() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "GcpAtSes", "wrong cloud account type")
			assert.Contains(t, body, "test@project.iam.gserviceaccount.com", "wrong client email")
			assert.Contains(t, body, "123456789", "wrong client ID")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
			assert.NotContains(t, body, "privateKey", "field whould not be displayed if empty")
			assert.NotContains(t, body, "privateKeyId", "field whould not be displayed if empty")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(singleGcpAtCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.GcpAtSesCloudAccount,
		api.GcpAtSesData{
			Credentials: api.GcpAtSesCredentials{
				ClientID:     "123456789",
				ClientEmail:  "test@project.iam.gserviceaccount.com",
				PrivateKeyID: "", // testing empty fields for both,
				PrivateKey:   "", // private key and private key id
			},
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "GcpAtSes cloud account name mismatch")
	assert.Equal(t, "GcpAtSes", cloudAccount.Type, "a new GcpAtSes cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new GcpAtSes cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	response, err := c.V2.CloudAccounts.UpdateGcpAtSes(cloudAccount)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "test-gcp", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "test@project.iam.gserviceaccount.com", response.Data.Data.Credentials.ClientEmail)
	assert.Equal(t, "123456789", response.Data.Data.Credentials.ClientID)
	assert.Equal(t, "", response.Data.Data.Credentials.PrivateKeyID)
	assert.Equal(t, "", response.Data.Data.Credentials.PrivateKey)
	assert.Equal(t, "GcpAtSes", response.Data.Type)
}

func singleGcpAtCloudAccount(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "test@lacework.net",
        "createdOrUpdatedTime": "2022-04-29T00:33:16.964Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 0,
        "name": "test-gcp",
        "state": {
            "ok": true,
            "lastUpdatedTime": 1663063753455,
            "lastSuccessfulTime": 1663063753455,
            "details": {
                "queueRx": "OK",
                "decodeNtfn": "OK",
                "logFileGet": "OK",
                "queueDel": "OK",
                "lastMsgRxTime": 1663063753455,
                "noData": false
            }
        },
        "type": "GcpAtSes",
        "data": {
            "credentials": {
                "clientId": "123456789",
                "privateKeyId": "",
                "clientEmail": "test@project.iam.gserviceaccount.com",
                "privateKey": ""
            },
            "idType": "PROJECT",
            "id": "test-project-123",
            "subscriptionName": "projects/test-project/subscriptions/test"
        }
    }`, id)
}
