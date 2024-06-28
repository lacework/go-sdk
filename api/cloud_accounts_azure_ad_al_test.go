//
// Author:: Rubinder Singh (<rubinder.singh@lacework.net>)
// Copyright:: Copyright 2024, Lacework Inc.
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

	"github.com/circleci/fork-lacework-go-sdk/api"
	"github.com/circleci/fork-lacework-go-sdk/internal/intgguid"
	"github.com/circleci/fork-lacework-go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestCloudAccountsAzureAdAlGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAzureAdAl() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(azureAdAlCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAzureAdAl(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "azure_ad_al_integration_test", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "123456777", response.Data.Data.Credentials.ClientID)
	assert.Equal(t, "test-secret-1234", response.Data.Data.Credentials.ClientSecret)
	assert.Equal(t, "AzureAdAl", response.Data.Type)
	assert.Equal(t, "tenant-1", response.Data.Data.TenantID)
	assert.Equal(t, "eventHubNamespace-1", response.Data.Data.EventHubNamespace)
	assert.Equal(t, "eventHubName-1", response.Data.Data.EventHubName)
}

func TestCloudAccountsAzureAdAlUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
		intgData   = api.AzureAdAlData{
			TenantID: "tenant-1",
			Credentials: api.AzureAdAlCredentials{
				ClientID:     "123456777",
				ClientSecret: "test-secret-1234",
			},
			EventHubNamespace: "eventHubNamespace-1",
			EventHubName:      "eventHubName-1",
		}
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	// Step 1 - Start Fake Server to return updated data
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAzureAdAl() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "azure_ad_al_integration_test", "cloud account name is missing")
			assert.Contains(t, body, "AzureAdAl", "wrong cloud account type")
			assert.Contains(t, body, intgData.Credentials.ClientID, "wrong ClientId")
			assert.Contains(t, body, intgData.Credentials.ClientSecret, "wrong ClientSecret")
			assert.Contains(t, body, intgData.TenantID, "wrong TenantId")
			assert.Contains(t, body, intgData.EventHubNamespace, "wrong EventHubNamespace")
			assert.Contains(t, body, intgData.EventHubName, "wrong EventHubName")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(azureAdAlCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	// Step 2 - Get Updated data from Fake server
	cloudAccount := api.NewCloudAccount("azure_ad_al_integration_test",
		api.AzureAdAlCloudAccount,
		intgData,
	)

	cloudAccount.IntgGuid = intgGUID
	response, err := c.V2.CloudAccounts.UpdateAzureAdAl(cloudAccount)
	assert.Nil(t, err, "Cannot update integration")
	assert.NotNil(t, response)
	integration := response.Data
	assert.Equal(t, intgGUID, integration.IntgGuid)

	integrationData := integration.Data
	assert.Equal(t, "azure_ad_al_integration_test", cloudAccount.Name)
	assert.Equal(t, "AzureAdAl", cloudAccount.Type)
	assert.Equal(t, 1, cloudAccount.Enabled)
	assert.Equal(t, "tenant-1", integrationData.TenantID)
	assert.Equal(t, "eventHubNamespace-1", integrationData.EventHubNamespace)
	assert.Equal(t, "eventHubName-1", integrationData.EventHubName)
	assert.Equal(t, "123456777", integrationData.Credentials.ClientID)
	assert.Equal(t, "test-secret-1234", integrationData.Credentials.ClientSecret)
}

func azureAdAlCloudAccount(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "rubinder.singh@lacework.net",
        "createdOrUpdatedTime": "2024-03-11T00:00:00.000Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 0,
        "name": "azure_ad_al_integration_test",
        "state": {
            "ok": true,
            "lastUpdatedTime": 1710104691000,
            "lastSuccessfulTime": 1710104691000,
            "details": {
                "queueRx": "OK",
                "decodeNtfn": "OK",
                "logFileGet": "OK",
                "queueDel": "OK",
                "lastMsgRxTime": 1710104691000,
                "noData": true
            }
        },
        "type": "AzureAdAl",
        "data": {
            "credentials": {
                "clientId": "123456777",
                "clientSecret": "test-secret-1234"
            },
            "tenantId": "tenant-1",
            "eventHubNamespace": "eventHubNamespace-1",
            "eventHubName": "eventHubName-1"
        }
    }`, id)
}
