//
// Copyright:: Copyright 2026, Lacework Inc.
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
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/v2/api"
	"github.com/lacework/go-sdk/v2/internal/intgguid"
	"github.com/lacework/go-sdk/v2/internal/lacework"
)

func TestCloudAccountsCreateAwsFortiDspm(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	defer fakeServer.Close()

	fakeServer.MockAPI("CloudAccounts", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		body, _ := io.ReadAll(r.Body)
		var sent map[string]interface{}
		assert.Nil(t, json.Unmarshal(body, &sent))
		assert.Equal(t, "AwsDspm", sent["type"])
		data := sent["data"].(map[string]interface{})
		assert.Equal(t, "459969247747", data["awsAccountId"])
		assert.Equal(t, []interface{}{"us-west-2"}, data["regions"])
		_, hasBucket := data["bucketArn"]
		assert.False(t, hasBucket, "a FortiDSPM-managed account must not carry a bucket")
		fmt.Fprint(w, generateCloudAccountResponse(fortiDspmAwsCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	account := api.NewCloudAccount("aws-dspm-459969247747", api.AwsDspmCloudAccount,
		api.AwsFortiDspmData{AccountID: "459969247747", Regions: []string{"us-west-2"}})
	response, err := c.V2.CloudAccounts.CreateAwsFortiDspm(account)
	assert.Nil(t, err)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "459969247747", response.Data.Data.AccountID)

	deployment := response.Data.FortiDspmDeployment
	assert.NotNil(t, deployment)
	assert.Equal(t, "3c0d5a0e-1b47-4f5b-9b8e-2a5b1f0d7c11", deployment.DeploymentID)
	assert.Equal(t, 86400, deployment.TokenExpiresIn)
	assert.Equal(t, map[string]string{"us-west-2": "eyJhbGciOiJSUzI1NiJ9.token"}, deployment.ActivationTokens())
	assert.Equal(t, map[string]string{"us-west-2": "ami-0db586fa50c952fe9"}, deployment.Images())
	assert.Empty(t, deployment.HypervGenerations())
}

func TestCloudAccountsCreateAzureFortiDspm(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	defer fakeServer.Close()

	fakeServer.MockAPI("CloudAccounts", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		fmt.Fprint(w, generateCloudAccountResponse(fortiDspmAzureCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	account := api.NewCloudAccount("azure-dspm-tenant", api.AzureDspmCloudAccount,
		api.AzureFortiDspmData{
			TenantID:       "fbdd69b5-5210-47e3-9dce-608699edc599",
			SubscriptionID: "51bb4964-95c7-4718-bdf5-355368aaba4a",
			Regions:        []string{"westus2"},
		})
	response, err := c.V2.CloudAccounts.CreateAzureFortiDspm(account)
	assert.Nil(t, err)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "fbdd69b5-5210-47e3-9dce-608699edc599", response.Data.Data.TenantID)

	deployment := response.Data.FortiDspmDeployment
	assert.NotNil(t, deployment)
	assert.Equal(t, 14400, deployment.ImageURLExpiresIn)
	assert.Equal(t, map[string]string{
		"westus2": "https://fortidspmdev.blob.core.windows.net/fortidspm-pkg/x.vhd?sp=r&sig=abc",
	}, deployment.Images())
	assert.Equal(t, map[string]string{"westus2": "V1"}, deployment.HypervGenerations())
	assert.Equal(t, int64(6867124736), deployment.Regions[0].ImageSizeBytes)
}

func TestCloudAccountsGetAwsFortiDspmHasNoDeployment(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, generateCloudAccountResponse(fortiDspmAwsCloudAccountWithoutDeployment(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAwsFortiDspm(intgGUID)
	assert.Nil(t, err)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Nil(t, response.Data.FortiDspmDeployment)
	assert.Equal(t, []string{"us-west-2"}, response.Data.Data.Regions)
}

func fortiDspmAwsCloudAccountWithoutDeployment(id string) string {
	return `
    {
      "createdOrUpdatedBy": "self-deployment@lacework.net",
      "createdOrUpdatedTime": "2026-09-11T20:00:00.000Z",
      "enabled": 1,
      "intgGuid": "` + id + `",
      "isOrg": 0,
      "name": "aws-dspm-459969247747",
      "state": {"ok": true, "lastUpdatedTime": 1757620800000, "lastSuccessfulTime": 1757620800000, "details": {}},
      "type": "AwsDspm",
      "data": {"awsAccountId": "459969247747", "regions": ["us-west-2"]}
    }
	`
}

func fortiDspmAwsCloudAccount(id string) string {
	return `
    {
      "createdOrUpdatedBy": "self-deployment@lacework.net",
      "createdOrUpdatedTime": "2026-09-11T20:00:00.000Z",
      "enabled": 1,
      "intgGuid": "` + id + `",
      "isOrg": 0,
      "name": "aws-dspm-459969247747",
      "state": {"ok": true, "lastUpdatedTime": 1757620800000, "lastSuccessfulTime": 1757620800000, "details": {}},
      "type": "AwsDspm",
      "data": {"awsAccountId": "459969247747", "regions": ["us-west-2"]},
      "fortiDspmDeployment": {
        "deploymentId": "3c0d5a0e-1b47-4f5b-9b8e-2a5b1f0d7c11",
        "deploymentName": "aws-dspm-459969247747",
        "envId": "nancy",
        "tokenExpiresIn": 86400,
        "regions": [
          {"region": "us-west-2", "activationToken": "eyJhbGciOiJSUzI1NiJ9.token", "amiId": "ami-0db586fa50c952fe9"}
        ]
      }
    }
	`
}

func fortiDspmAzureCloudAccount(id string) string {
	return `
    {
      "createdOrUpdatedBy": "self-deployment@lacework.net",
      "createdOrUpdatedTime": "2026-09-11T20:00:00.000Z",
      "enabled": 1,
      "intgGuid": "` + id + `",
      "isOrg": 0,
      "name": "azure-dspm-tenant",
      "state": {"ok": true, "lastUpdatedTime": 1757620800000, "lastSuccessfulTime": 1757620800000, "details": {}},
      "type": "AzureDspm",
      "data": {
        "tenantId": "fbdd69b5-5210-47e3-9dce-608699edc599",
        "subscriptionId": "51bb4964-95c7-4718-bdf5-355368aaba4a",
        "regions": ["westus2"]
      },
      "fortiDspmDeployment": {
        "deploymentId": "7c9e2f10-4a3b-4d2e-8f61-0b1c2d3e4f55",
        "deploymentName": "azure-dspm-tenant",
        "tokenExpiresIn": 86400,
        "imageUrlExpiresIn": 14400,
        "regions": [
          {
            "region": "westus2",
            "activationToken": "eyJhbGciOiJSUzI1NiJ9.token2",
            "imageUrl": "https://fortidspmdev.blob.core.windows.net/fortidspm-pkg/x.vhd?sp=r&sig=abc",
            "hypervGeneration": "V1",
            "imageSizeBytes": 6867124736
          }
        ]
      }
    }
	`
}

func TestCloudAccountsReportFortiDspmDeploymentStatus(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("FortiDspm/deployments/%s/status", intgGUID)
		fakeServer = lacework.MockServer()
	)
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		body, _ := io.ReadAll(r.Body)
		var sent map[string]interface{}
		assert.Nil(t, json.Unmarshal(body, &sent))
		assert.Equal(t, "succeeded", sent["status"])
		assert.Equal(t, "3c0d5a0e-1b47-4f5b-9b8e-2a5b1f0d7c11", sent["deploymentId"])
		region := sent["regions"].([]interface{})[0].(map[string]interface{})
		assert.Equal(t, "us-west-2", region["region"])
		assert.Equal(t, "i-0449750feb0cb12e9", region["instanceId"])
		_, hasErr := sent["error"]
		assert.False(t, hasErr)
		fmt.Fprint(w, `{"data": {"ok": true}}`)
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	err = c.V2.CloudAccounts.ReportFortiDspmDeploymentStatus(intgGUID, api.FortiDspmDeploymentStatus{
		DeploymentID: "3c0d5a0e-1b47-4f5b-9b8e-2a5b1f0d7c11",
		Status:       "succeeded",
		Regions: []api.FortiDspmRegionStatus{{
			Region: "us-west-2", Status: "succeeded", InstanceID: "i-0449750feb0cb12e9",
			NatGatewayPublicIP: "54.12.34.56", IamRoleArn: "arn:aws:iam::459969247747:role/fortidspm-scan-engine",
		}},
	})
	assert.Nil(t, err)
}
