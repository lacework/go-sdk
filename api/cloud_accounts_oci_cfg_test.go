//
// Author:: Kolbeinn Karlsson (<kolbeinn.karlsson@lacework.net>)
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

func TestCloudAccountsOciCfgGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetOciCfg() should be a GET method")
		fmt.Fprint(w, generateCloudAccountResponse(singleOciCfgCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetOciCfg(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assertResponseMatches(t, response)
}

func TestCloudAccountsOciCfgUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateOciCfg() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "OciCfg", "wrong cloud account type")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
			assert.Contains(t, body, "us-sanjose-1", "wrong home region")
			assert.Contains(t, body, "ocid1.tenancy.oc1..abcdefghijklmnopqrstuvwxyz", "wrong tenant ID")
			assert.Contains(t, body, "integration_tenant", "wrong tenant name")
			assert.Contains(t, body, "ocid1.user.oc1..abcdefghijklmnopqrstuvwxyz", "wrong user OCID")
			assert.Contains(t, body, "01:02:03:04:05:06:07:08:09:0A::0B:0C:0D:0E:0F", "wrong fingerprint")
		}

		fmt.Fprint(w, generateCloudAccountResponse(singleOciCfgCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.OciCfgCloudAccount,
		api.OciCfgData{
			Credentials: api.OciCfgCredentials{
				Fingerprint: "01:02:03:04:05:06:07:08:09:0A::0B:0C:0D:0E:0F",
				PrivateKey:  "",
			},
			HomeRegion: "us-sanjose-1",
			TenantID:   "ocid1.tenancy.oc1..abcdefghijklmnopqrstuvwxyz",
			TenantName: "integration_tenant",
			UserOCID:   "ocid1.user.oc1..abcdefghijklmnopqrstuvwxyz",
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "OciCfg cloud account name mismatch")
	assert.Equal(t, "OciCfg", cloudAccount.Type, "a new OciCfg cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new OciCfg cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	response, err := c.V2.CloudAccounts.UpdateOciCfg(cloudAccount)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assertResponseMatches(t, response)
}

func assertResponseMatches(t *testing.T, response api.OciCfgIntegrationResponse) {
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "us-sanjose-1", response.Data.Data.HomeRegion)
	assert.Equal(t, "ocid1.tenancy.oc1..abcdefghijklmnopqrstuvwxyz", response.Data.Data.TenantID)
	assert.Equal(t, "integration_tenant", response.Data.Data.TenantName)
	assert.Equal(t, "ocid1.user.oc1..abcdefghijklmnopqrstuvwxyz", response.Data.Data.UserOCID)
	assert.Equal(
		t,
		"01:02:03:04:05:06:07:08:09:0A::0B:0C:0D:0E:0F",
		response.Data.Data.Credentials.Fingerprint,
	)
}

func singleOciCfgCloudAccount(id string) string {
	return `{
		"createdOrUpdatedBy": "test@lacework.net",
		"createdOrUpdatedTime": "2023-06-22T09:45:37.766Z",
		"enabled": 1,
		"intgGuid": "` + id + `",
		"isOrg": 0,
		"name": "integration_name",
		"state": {
			"ok": true,
			"lastUpdatedTime": 1683199976188,
			"lastSuccessfulTime": 1683199976188,
			"details": {}
		},
		"type": "OciCfg",
		"data": {
			"homeRegion": "us-sanjose-1",
			"tenantId": "ocid1.tenancy.oc1..abcdefghijklmnopqrstuvwxyz",
			"tenantName": "integration_tenant",
			"userOcid": "ocid1.user.oc1..abcdefghijklmnopqrstuvwxyz",
			"credentials": {
				"fingerprint": "01:02:03:04:05:06:07:08:09:0A::0B:0C:0D:0E:0F",
				"privateKey": ""
			}
		}
	}`
}
