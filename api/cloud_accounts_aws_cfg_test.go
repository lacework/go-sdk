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

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestCloudAccountsNewAwsCfgWithCustomTemplateFile(t *testing.T) {
	awsCfgData := api.AwsCfgData{
		Credentials: api.AwsCfgCredentials{
			RoleArn:    "arn:foo:bar",
			ExternalID: "0123456789",
		},
	}

	subject := api.NewCloudAccount("integration_name", api.AwsCfgCloudAccount, awsCfgData)
	assert.Equal(t, api.AwsCfgCloudAccount.String(), subject.Type)

	// casting the data interface{} to type AwsCfgData
	subjectData := subject.Data.(api.AwsCfgData)

	assert.Equal(t, subjectData.Credentials.RoleArn, "arn:foo:bar")
	assert.Equal(t, subjectData.Credentials.ExternalID, "0123456789")
}

func TestCloudAccountsAwsCfgGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsCfg() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsCfgCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAwsCfg(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "arn:foo:bar", response.Data.Data.Credentials.RoleArn)
	assert.Equal(t, "0123456789", response.Data.Data.Credentials.ExternalID)
}

func TestCloudAccountsAwsCfgUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsCfg() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "AwsCfg", "wrong cloud account type")
			assert.Contains(t, body, "arn:incorrect:rolearn", "wrong role arn")
			assert.Contains(t, body, "abc123", "wrong external ID")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsCfgCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.AwsCfgCloudAccount,
		api.AwsCfgData{
			Credentials: api.AwsCfgCredentials{
				RoleArn:    "arn:foo:bar",
				ExternalID: "abc123",
			},
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "AwsCfg cloud account name mismatch")
	assert.Equal(t, "AwsCfg", cloudAccount.Type, "a new AwsCfg cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new AwsCfg cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	response, err := c.V2.CloudAccounts.UpdateAwsCfg(cloudAccount)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "arn:foo:bar", response.Data.Data.Credentials.RoleArn)
	assert.Equal(t, "abc123", response.Data.Data.Credentials.ExternalID)
}

func singleAwsCfgCloudAccount(id string) string {
	return `
  {
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
    "data": {
      "awsAccountId": "123456789000",
      "crossAccountCredentials": {
        "externalId": "0123456789",
        "roleArn": "arn:foo:bar"
      }
    },
    "enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "state": {
      "details": {},
      "lastSuccessfulTime": 1624456896915,
      "lastUpdatedTime": 1624456896915,
      "ok": true
    },
    "type": "AwsCfg"
  }
  `
}
