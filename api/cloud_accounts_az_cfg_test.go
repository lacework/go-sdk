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

func TestCloudAccountsAzCfgGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAzureCfg() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleAzCfgloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAzureCfg(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "azure-cfg-test", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "123456789", response.Data.Data.Credentials.ClientID)
	assert.Equal(t, "", response.Data.Data.Credentials.ClientSecret)
	assert.Equal(t, "AzureCfg", response.Data.Type)
	assert.Equal(t, "abcdefgh", response.Data.Data.TenantID)
}

func singleAzCfgloudAccount(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "test@lacework.net",
        "createdOrUpdatedTime": "2022-03-23T16:41:16.039Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 0,
        "name": "azure-cfg-test",
        "state": {
            "ok": true,
            "lastUpdatedTime": 1663062626237,
            "lastSuccessfulTime": 1663062626237,
            "details": {
                "tenantErrors": {
                    "opsDeniedAccess": []
                },
                "subscriptionErrors": {
                    "12345678-123456678": {
                        "opsDeniedAccess": []
                    }
                }
            }
        },
        "type": "AzureCfg",
        "data": {
            "credentials": {
                "clientId": "123456789",
                "clientSecret": ""
            },
            "tenantId": "abcdefgh"
        }
    }`, id)
}
