//
// Author:: Ao Zhang (<ao.zhang@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	azureSidekickData = api.AzureSidekickData{
		IntegrationType:         "SUBSCRIPTION",
		SubscriptionId:          "54321",
		TenantId:                "98765",
		BlobContainerName:       "blobContainer",
		ScanFrequency:           24,
		ScanContainers:          true,
		ScanHostVulnerabilities: true,
		Credentials: api.AzureSidekickCredentials{
			ClientID:       "Client123",
			CredentialType: "ShareCredentials",
			ClientSecret:   "Secret",
		},
		SubscriptionList: "sub1,sub2",
		QueryText:        "queryText",
	}

	azureUpdatedSidekickData = api.AzureSidekickData{
		IntegrationType:         "SUBSCRIPTION",
		SubscriptionId:          "12345",
		TenantId:                "87654",
		SubscriptionId:          "updated-54321",
		BlobContainerName:       "updated-blobContainer",
		ScanFrequency:           12,
		ScanContainers:          false,
		ScanHostVulnerabilities: true,
		Credentials: api.AzureSidekickCredentials{
			ClientID:       "updated-Client123",
			CredentialType: "updated-SharedCredentials",
			ClientSecret:   "updated-secret",
		},
		SubscriptionList: "updated-proj1,proj2",
		QueryText:        "updated-queryText",
	}
)

// TODO: update this
func TestCloudAccountsAzureSidekickCreate(t *testing.T) {
	accountMappingJSON := []byte(`{
		"defaultLaceworkAccountAws": "lw_account_1",
		"integration_mappings": {
		  "lw_account_2": {
			"aws_accounts": [
			  "234556677",
			  "774564564"
			]
		  },
		  "lw_account_3": {
			"aws_accounts": [
			  "553453453",
			  "934534535"
			]
		  }
		}
	  }`)
	integration := api.NewCloudAccount("integration_name", api.AzureSidekickCloudAccount, azureSidekickData)
	assert.Equal(t, api.AzureSidekickCloudAccount.String(), integration.Type)

	// casting the data interface{} to type AzureSidekickData
	integrationData := integration.Data.(api.AzureSidekickData)
	integrationData.EncodeAccountMappingFile(accountMappingJSON)

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
	assert.Contains(t,
		integrationData.AccountMappingFile,
		"data:application/json;name=i.json;base64,",
		"check the custom_template_file encoder",
	)
	accountMapping, err := integrationData.DecodeAccountMappingFile()
	assert.Nil(t, err)
	assert.Equal(t, accountMappingJSON, accountMapping)

	// When there is no custom account mapping file, this function should
	// return an empty string to match the pattern
	integrationData.AccountMappingFile = ""
	accountMapping, err = integrationData.DecodeAccountMappingFile()
	assert.Nil(t, err)
	assert.Empty(t, accountMapping)
}

func TestCloudAccountsAzureSidekickGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAzureSidekick() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(getAzureData(intgGUID, azureSidekickData)))

	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAzureSidekick(intgGUID)
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
	assert.Equal(t, "token_"+integration.IntgGuid, integration.ServerToken)
}

func TestCloudAccountsAzureSidekickUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	// Step 1 - Start Fake Server to return updated data
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAzureSidekick() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_test", "cloud account name is missing")
			assert.Contains(t, body, "AzureSidekick", "wrong cloud account type")
			assert.Contains(t, body, azureSidekickData.Credentials.ClientID, "wrong client ID")
			assert.Contains(t, body, azureSidekickData.Credentials.ClientEmail, "wrong client email")
			assert.Contains(t, body, azureSidekickData.SharedBucket, "wrong client email")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(getAzureData(intgGUID, azureUpdatedSidekickData)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	// Step 2 - Create New Account
	cloudAccount := api.NewCloudAccount("integration_test",
		api.AzureSidekickCloudAccount,
		azureSidekickData,
	)

	integrationData := cloudAccount.Data.(api.AzureSidekickData)
	assert.Equal(t, "integration_test", cloudAccount.Name)
	assert.Equal(t, "AzureSidekick", cloudAccount.Type)
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
	response, err := c.V2.CloudAccounts.UpdateAzureSidekick(cloudAccount)
	assert.Nil(t, err, "Cannot update integration")
	assert.NotNil(t, response)
	integration := response.Data
	assert.Equal(t, intgGUID, integration.IntgGuid)

	integrationData = integration.Data
	assert.Equal(t, "integration_test", cloudAccount.Name)
	assert.Equal(t, "AzureSidekick", cloudAccount.Type)
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

// getAzureData converts integration data to json string
func getAzureData(id string, data api.AzureSidekickData) string {

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
  	"type": "AzureSidekick",
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
	},
	"serverToken": {
		"serverToken": "token_` + id + `"
	}
  }
  `
}
