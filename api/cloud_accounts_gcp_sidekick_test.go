//
// Author:: Ammar Ekbote (<ammar.ekbote@lacework.net>)
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

// These two objects are used to test Create, Get and Update operations.
var (
	gcpSidekickData = api.GcpSidekickData{
		IDType:                  "PROJECT",
		ID:                      "12345",
		ScanningProjectId:       "54321",
		SharedBucket:            "storageBucket",
		ScanFrequency:           24,
		ScanContainers:          true,
		ScanHostVulnerabilities: true,
		Credentials: api.GcpSidekickCredentials{
			ClientID:     "Client123",
			ClientEmail:  "client@test.com",
			PrivateKeyID: "privateKeyID",
			PrivateKey:   "privateKey",
			TokenUri:     "tokenTest",
		},
		FilterList: "proj1,proj2",
		QueryText:  "queryText",
	}

	gcpUpdatedSidekickData = api.GcpSidekickData{
		IDType:                  "PROJECT",
		ID:                      "12345",
		ScanningProjectId:       "updated-54321",
		SharedBucket:            "updated-storageBucket",
		ScanFrequency:           12,
		ScanContainers:          false,
		ScanHostVulnerabilities: true,
		Credentials: api.GcpSidekickCredentials{
			ClientID:     "updated-Client123",
			ClientEmail:  "updated-client@test.com",
			PrivateKeyID: "updated-privateKeyID",
			PrivateKey:   "updated-privateKey",
			TokenUri:     "updated-tokenTest",
		},
		FilterList: "updated-proj1,proj2",
		QueryText:  "updated-queryText",
	}
)

func TestCloudAccountsGcpSidekickCreate(t *testing.T) {
	integration := api.NewCloudAccount("integration_name", api.GcpSidekickCloudAccount, gcpSidekickData)
	assert.Equal(t, api.GcpSidekickCloudAccount.String(), integration.Type)

	// casting the data interface{} to type GcpSidekickData
	integrationData := integration.Data.(api.GcpSidekickData)

	assert.Equal(t, integrationData.IDType, "PROJECT")
	assert.Equal(t, integrationData.ID, "12345")
	assert.Equal(t, integrationData.ScanningProjectId, "54321")
	assert.Equal(t, integrationData.SharedBucket, "storageBucket")
	assert.Equal(t, integrationData.ScanFrequency, 24)
	assert.Equal(t, integrationData.ScanContainers, true)
	assert.Equal(t, integrationData.ScanHostVulnerabilities, true)

	assert.Equal(t, integrationData.Credentials.ClientID, "Client123")
	assert.Equal(t, integrationData.Credentials.ClientEmail, "client@test.com")
	assert.Equal(t, integrationData.Credentials.PrivateKeyID, "privateKeyID")
	assert.Equal(t, integrationData.Credentials.PrivateKey, "privateKey")
	assert.Equal(t, integrationData.Credentials.TokenUri, "tokenTest")
}

func TestCloudAccountsGcpSidekickGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpSidekick() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(getGcpData(intgGUID, gcpSidekickData)))

	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetGcpSidekick(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	integration := response.Data
	assert.Equal(t, intgGUID, integration.IntgGuid)
	assert.Equal(t, "integration_test", integration.Name)
	assert.True(t, integration.State.Ok)

	integrationData := integration.Data
	assert.Equal(t, "PROJECT", integrationData.IDType)
	assert.Equal(t, "12345", integrationData.ID)
	assert.Equal(t, "54321", integrationData.ScanningProjectId)
	assert.Equal(t, "storageBucket", integrationData.SharedBucket)
	assert.Equal(t, 24, integrationData.ScanFrequency)
	assert.Equal(t, true, integrationData.ScanContainers)
	assert.Equal(t, true, integrationData.ScanHostVulnerabilities)
	assert.Equal(t, "Client123", integrationData.Credentials.ClientID)
	assert.Equal(t, "client@test.com", integrationData.Credentials.ClientEmail)
	assert.Equal(t, "privateKeyID", integrationData.Credentials.PrivateKeyID)
	assert.Equal(t, "privateKey", integrationData.Credentials.PrivateKey)
	assert.Equal(t, "tokenTest", integrationData.Credentials.TokenUri)
	assert.Equal(t, "proj1,proj2", integrationData.FilterList)
	assert.Equal(t, "queryText", integrationData.QueryText)
}

func TestCloudAccountsGcpSidekickUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	// Step 1 - Start Fake Server to return updated data
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGcpSidekick() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_test", "cloud account name is missing")
			assert.Contains(t, body, "GcpSidekick", "wrong cloud account type")
			assert.Contains(t, body, gcpSidekickData.Credentials.ClientID, "wrong client ID")
			assert.Contains(t, body, gcpSidekickData.Credentials.ClientEmail, "wrong client email")
			assert.Contains(t, body, gcpSidekickData.SharedBucket, "wrong client email")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(getGcpData(intgGUID, gcpUpdatedSidekickData)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	// Step 2 - Create New Account
	cloudAccount := api.NewCloudAccount("integration_test",
		api.GcpSidekickCloudAccount,
		gcpSidekickData,
	)

	integrationData := cloudAccount.Data.(api.GcpSidekickData)
	assert.Equal(t, "integration_test", cloudAccount.Name)
	assert.Equal(t, "GcpSidekick", cloudAccount.Type)
	assert.Equal(t, 1, cloudAccount.Enabled)
	assert.Equal(t, "PROJECT", integrationData.IDType)
	assert.Equal(t, "12345", integrationData.ID)
	assert.Equal(t, "54321", integrationData.ScanningProjectId)
	assert.Equal(t, "storageBucket", integrationData.SharedBucket)
	assert.Equal(t, 24, integrationData.ScanFrequency)
	assert.Equal(t, true, integrationData.ScanContainers)
	assert.Equal(t, true, integrationData.ScanHostVulnerabilities)
	assert.Equal(t, "Client123", integrationData.Credentials.ClientID)
	assert.Equal(t, "client@test.com", integrationData.Credentials.ClientEmail)
	assert.Equal(t, "privateKeyID", integrationData.Credentials.PrivateKeyID)
	assert.Equal(t, "privateKey", integrationData.Credentials.PrivateKey)
	assert.Equal(t, "tokenTest", integrationData.Credentials.TokenUri)
	assert.Equal(t, "proj1,proj2", integrationData.FilterList)
	assert.Equal(t, "queryText", integrationData.QueryText)

	// Step 3 - Get Updated data from Fake server
	cloudAccount.IntgGuid = intgGUID
	response, err := c.V2.CloudAccounts.UpdateGcpSidekick(cloudAccount)
	assert.Nil(t, err, "Cannot update integration")
	assert.NotNil(t, response)
	integration := response.Data
	assert.Equal(t, intgGUID, integration.IntgGuid)

	integrationData = integration.Data
	assert.Equal(t, "integration_test", cloudAccount.Name)
	assert.Equal(t, "GcpSidekick", cloudAccount.Type)
	assert.Equal(t, 1, cloudAccount.Enabled)
	assert.Equal(t, "PROJECT", integrationData.IDType)
	assert.Equal(t, "12345", integrationData.ID)
	assert.Equal(t, "updated-54321", integrationData.ScanningProjectId)
	assert.Equal(t, "updated-storageBucket", integrationData.SharedBucket)
	assert.Equal(t, 12, integrationData.ScanFrequency)
	assert.Equal(t, false, integrationData.ScanContainers)
	assert.Equal(t, true, integrationData.ScanHostVulnerabilities)
	assert.Equal(t, "updated-Client123", integrationData.Credentials.ClientID)
	assert.Equal(t, "updated-client@test.com", integrationData.Credentials.ClientEmail)
	assert.Equal(t, "updated-privateKeyID", integrationData.Credentials.PrivateKeyID)
	assert.Equal(t, "updated-privateKey", integrationData.Credentials.PrivateKey)
	assert.Equal(t, "updated-tokenTest", integrationData.Credentials.TokenUri)
	assert.Equal(t, "updated-proj1,proj2", integrationData.FilterList)
	assert.Equal(t, "updated-queryText", integrationData.QueryText)
}

// getGcpData converts integration data to json string
func getGcpData(id string, data api.GcpSidekickData) string {

	scanFrequency := fmt.Sprintf("%d", data.ScanFrequency)
	scanContainers := fmt.Sprintf("%t", data.ScanContainers)
	scanHostVulnerabilities := fmt.Sprintf("%t", data.ScanHostVulnerabilities)

	return `
  {
  	"createdOrUpdatedBy": "ammar.ekbote@lacework.net",
  	"createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
  	"enabled": 1,
  	"intgGuid": "` + id + `",
  	"isOrg": 0,
  	"name": "integration_test",
  	"state": {
  		"details": {},
  		"lastSuccessfulTime": 1624456896915,
  		"lastUpdatedTime": 1624456896915,
  		"ok": true
  	},
  	"type": "GcpSidekick",
  	"data": {
  		"credentials": {
  			"clientId": "` + data.Credentials.ClientID + `",
  			"privateKeyId": "` + data.Credentials.PrivateKeyID + `",
  			"clientEmail": "` + data.Credentials.ClientEmail + `",
  			"privateKey": "` + data.Credentials.PrivateKey + `",
  			"tokenuri": "` + data.Credentials.TokenUri + `"
  		},
  		"idType": "` + data.IDType + `",
  		"id": "` + data.ID + `",
  		"scanningProjectId":  "` + data.ScanningProjectId + `",
  		"sharedBucketName": "` + data.SharedBucket + `",
  		"filterList": "` + data.FilterList + `",
  		"queryText": "` + data.QueryText + `",
  		"scanFrequency": ` + scanFrequency + `,
  		"scanContainers":  ` + scanContainers + `,
  		"scanHostVulnerabilities": ` + scanHostVulnerabilities + `
  	}
  }
  `
}
