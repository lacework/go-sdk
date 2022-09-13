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

func TestCloudAccountsAzAlSeqGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAzureAlSeq() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleAzAlSeqCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAzureAlSeq(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "azure-al-seq-test", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "123456789", response.Data.Data.Credentials.ClientID)
	assert.Equal(t, "", response.Data.Data.Credentials.ClientSecret)
	assert.Equal(t, "AzureAlSeq", response.Data.Type)
	assert.Equal(t, "https://test.queue.core.windows.net/lwtest", response.Data.Data.QueueUrl)
	assert.Equal(t, "abcdefgh", response.Data.Data.TenantID)
}

func singleAzAlSeqCloudAccount(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "test@lacework.net",
        "createdOrUpdatedTime": "2022-06-17T19:08:48.988Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 0,
        "name": "azure-al-seq-test",
        "state": {
            "ok": true,
            "lastUpdatedTime": 1663064333628,
            "lastSuccessfulTime": 1663039070672,
            "details": {
                "queueRx": "OK",
                "decodeNtfn": "OK",
                "logFileGet": "OK",
                "queueDel": "OK",
                "lastMsgRxTime": 1663035465558,
                "noData": true
            }
        },
        "type": "AzureAlSeq",
        "data": {
            "credentials": {
                "clientId": "123456789",
                "clientSecret": ""
            },
            "tenantId": "abcdefgh",
            "queueUrl": "https://test.queue.core.windows.net/lwtest"
        }
    }`, id)
}
