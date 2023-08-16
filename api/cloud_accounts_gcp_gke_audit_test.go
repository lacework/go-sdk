//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

func TestCloudAccountsGcpGkeAuditProjectLevelGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpGkeAudit() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(gcpGkeAuditProjectLevelCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetGcpGkeAudit(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "0123456789", response.Data.Data.Credentials.ClientId)
	assert.Equal(t, "ross.moles@lacework.net", response.Data.Data.Credentials.ClientEmail)
	assert.Empty(t, response.Data.Data.Credentials.PrivateKeyId)
	assert.Empty(t, response.Data.Data.Credentials.PrivateKey)
	assert.Equal(t, "Project", response.Data.Data.IntegrationType)
	assert.Equal(t, "ProjectMcProjectFace", response.Data.Data.ProjectId)
	assert.Equal(
		t,
		"projects/ProjectMcProjectFace/subscriptions/SubscribeyMcSubscribeFace",
		response.Data.Data.SubscriptionName,
	)
}

func TestCloudAccountsGcpGkeAuditOrganizationLevelGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpGkeAudit() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(gcpGkeAuditOrganizationLevelCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetGcpGkeAudit(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "0123456789", response.Data.Data.Credentials.ClientId)
	assert.Equal(t, "ross.moles@lacework.net", response.Data.Data.Credentials.ClientEmail)
	assert.Empty(t, response.Data.Data.Credentials.PrivateKeyId)
	assert.Empty(t, response.Data.Data.Credentials.PrivateKey)
	assert.Equal(t, "Project", response.Data.Data.IntegrationType)
	assert.Equal(t, "OrgMcOrgFace", response.Data.Data.OrganizationId)
	assert.Equal(t, "ProjectMcProjectFace", response.Data.Data.ProjectId)
	assert.Equal(
		t,
		"projects/ProjectMcProjectFace/subscriptions/SubscribeyMcSubscribeFace",
		response.Data.Data.SubscriptionName,
	)
}

func TestCloudAccountsGcpGkeAuditProjectLevelUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGcpGkeAudit() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "GcpGkeAudit", "wrong cloud account type")
			assert.Contains(t, body, "ross.moles@lacework.net", "wrong client email")
			assert.Contains(t, body, "0123456789", "wrong client email")
			assert.Contains(t, body, "\"integrationType\":\"Project", "wrong integration type")
			assert.Contains(t, body, "projectId\":\"ProjectMcProjectFace", "wrong project id")
			assert.Contains(
				t,
				body,
				"projects/ProjectMcProjectFace/subscriptions/SubscribeyMcSubscribeFace",
				"wrong subscription name")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(gcpGkeAuditProjectLevelCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.GcpGkeAuditCloudAccount,
		api.GcpGkeAuditData{
			Credentials: api.GcpGkeAuditCredentials{
				ClientEmail:  "ross.moles@lacework.net",
				ClientId:     "0123456789",
				PrivateKey:   "",
				PrivateKeyId: "",
			},
			IntegrationType: "Project",
			ProjectId:       "ProjectMcProjectFace",
			SubscriptionName: "projects/ProjectMcProjectFace/subscriptions/" +
				"SubscribeyMcSubscribeFace",
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "GcpGkeAudit cloud account name mismatch")
	assert.Equal(t, "GcpGkeAudit", cloudAccount.Type, "a new GcpGkeAudit cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new GcpGkeAudit cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	response, err := c.V2.CloudAccounts.UpdateGcpGkeAudit(cloudAccount)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t,
		"projects/ProjectMcProjectFace/subscriptions/SubscribeyMcSubscribeFace",
		response.Data.Data.SubscriptionName)
}

func gcpGkeAuditProjectLevelCloudAccount(id string) string {
	return `
  {
    "createdOrUpdatedBy": "ross.moles@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
	"enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "state": {
      "details": {
        "decodeNtfn": "OK",
        "lastMsgRxTime": 1655136633387,
        "logFileGet": "OK",
        "noData": false,
        "queueDel": "OK",
        "queueRx": "OK"
      },
      "lastSuccessfulTime": 1655136633387,
      "lastUpdatedTime": 1655136633387,
      "ok": true
    },
	"type": "GcpGkeAudit",
    "data": {
      "credentials": {
        "clientEmail": "ross.moles@lacework.net",
        "clientId": "0123456789",
        "privateKey": "",
        "privateKeyId": ""
      },
      "integrationType": "Project",
      "projectId": "ProjectMcProjectFace",
      "subscriptionName": "projects/ProjectMcProjectFace/subscriptions/SubscribeyMcSubscribeFace"
    }
  }
  `
}

func gcpGkeAuditOrganizationLevelCloudAccount(id string) string {
	return `
  {
    "createdOrUpdatedBy": "ross.moles@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
	"enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "state": {
      "details": {
        "decodeNtfn": "OK",
        "lastMsgRxTime": 1655136633387,
        "logFileGet": "OK",
        "noData": false,
        "queueDel": "OK",
        "queueRx": "OK"
      },
      "lastSuccessfulTime": 1655136633387,
      "lastUpdatedTime": 1655136633387,
      "ok": true
    },
	"type": "GcpGkeAudit",
    "data": {
      "credentials": {
        "clientEmail": "ross.moles@lacework.net",
        "clientId": "0123456789",
        "privateKey": "",
        "privateKeyId": ""
      },
      "integrationType": "Project",
      "organizationId": "OrgMcOrgFace",
      "projectId": "ProjectMcProjectFace",
      "subscriptionName": "projects/ProjectMcProjectFace/subscriptions/SubscribeyMcSubscribeFace"
    }
  }
  `
}
