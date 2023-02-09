//
// Author:: David McTavish (<david.mctavish@lacework.net>)
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

func TestCloudAccountsGcpAlPubSubProjectGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpAlPubSub() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleGcpAlPubSubProjectCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetGcpAlPubSub(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "test-gcp-project", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "test@project.iam.gserviceaccount.com", response.Data.Data.Credentials.ClientEmail)
	assert.Equal(t, "123456789", response.Data.Data.Credentials.ClientID)
	assert.Equal(t, "", response.Data.Data.Credentials.PrivateKeyID)
	assert.Equal(t, "", response.Data.Data.Credentials.PrivateKey)
	assert.Equal(t, "GcpAlPubSub", response.Data.Type)
	assert.Equal(t, "test-project-123", response.Data.Data.ProjectID)
	assert.Equal(t, "PROJECT", response.Data.Data.IntegrationType)
	assert.Equal(t, "projects/test-project/subscriptions/test", response.Data.Data.SubscriptionName)
	assert.Equal(t, "projects/test-project/subscriptions/test-topic", response.Data.Data.TopicID)
}

func TestCloudAccountsGcpAlPubSubOrgGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpAlPubSub() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleGcpAlPubSubOrgCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetGcpAlPubSub(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "test-gcp-org", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "test@project.iam.gserviceaccount.com", response.Data.Data.Credentials.ClientEmail)
	assert.Equal(t, "123456789", response.Data.Data.Credentials.ClientID)
	assert.Equal(t, "", response.Data.Data.Credentials.PrivateKeyID)
	assert.Equal(t, "", response.Data.Data.Credentials.PrivateKey)
	assert.Equal(t, "GcpAlPubSub", response.Data.Type)
	assert.Equal(t, "test-org-123", response.Data.Data.OrganizationID)
	assert.Equal(t, "ORGANIZATION", response.Data.Data.IntegrationType)
	assert.Equal(t, "projects/test-project/subscriptions/test", response.Data.Data.SubscriptionName)
	assert.Equal(t, "projects/test-project/subscriptions/test-topic", response.Data.Data.TopicID)
}

func singleGcpAlPubSubProjectCloudAccount(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "test@lacework.net",
        "createdOrUpdatedTime": "2022-04-29T00:33:16.964Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 0,
        "name": "test-gcp-project",
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
        "type": "GcpAlPubSub",
        "data": {
            "credentials": {
                "clientId": "123456789",
                "privateKeyId": "",
                "clientEmail": "test@project.iam.gserviceaccount.com",
                "privateKey": ""
            },
            "integrationType": "PROJECT",
            "projectId": "test-project-123",
            "subscriptionName": "projects/test-project/subscriptions/test",
			"topicId":"projects/test-project/subscriptions/test-topic"
        }
    }`, id)
}

func singleGcpAlPubSubOrgCloudAccount(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "test@lacework.net",
        "createdOrUpdatedTime": "2022-04-29T00:33:16.964Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 1,
        "name": "test-gcp-org",
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
        "type": "GcpAlPubSub",
        "data": {
            "credentials": {
                "clientId": "123456789",
                "privateKeyId": "",
                "clientEmail": "test@project.iam.gserviceaccount.com",
                "privateKey": ""
            },
            "integrationType": "ORGANIZATION",
            "organizationId": "test-org-123",
            "subscriptionName": "projects/test-project/subscriptions/test",
			"topicId":"projects/test-project/subscriptions/test-topic"
        }
    }`, id)
}
