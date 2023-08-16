//
// Author:: Teddy Reed (<teddy.reed@lacework.net>)
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

func TestCloudAccountsNewAwsSidekickWithCustomTemplateFile(t *testing.T) {
	awsCfgData := api.AwsSidekickData{
		CrossAccountCreds: api.AwsSidekickCrossAccountCredentials{
			RoleArn:    "arn:foo:bar",
			ExternalID: "0123456789",
		},
	}

	subject := api.NewCloudAccount("integration_name", api.AwsSidekickCloudAccount, awsCfgData)
	assert.Equal(t, api.AwsSidekickCloudAccount.String(), subject.Type)

	// casting the data interface{} to type AwsCfgData
	subjectData := subject.Data.(api.AwsSidekickData)

	assert.Equal(t, subjectData.CrossAccountCreds.RoleArn, "arn:foo:bar")
	assert.Equal(t, subjectData.CrossAccountCreds.ExternalID, "0123456789")
}

func TestCloudAccountsAwsSidekickGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsSidekick() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsSidekickCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAwsSidekick(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "arn:foo:bar", response.Data.Data.CrossAccountCreds.RoleArn)
	assert.Equal(t, "0123456789", response.Data.Data.CrossAccountCreds.ExternalID)
	assert.Equal(t, 24, response.Data.Data.ScanFrequency)
	assert.True(t, response.Data.Data.ScanContainers)
	assert.True(t, response.Data.Data.ScanHostVulnerabilities)
}

func TestCloudAccountsAwsSidekickUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsSidekick() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "AwsSidekick", "wrong cloud account type")
			assert.Contains(t, body, "arn:foo:bar", "wrong role arn")
			assert.Contains(t, body, "0123456789", "wrong external ID")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsSidekickCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.AwsSidekickCloudAccount,
		api.AwsSidekickData{
			CrossAccountCreds: api.AwsSidekickCrossAccountCredentials{
				RoleArn:    "arn:foo:bar",
				ExternalID: "0123456789",
			},
			ScanFrequency: 24,
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "AwsSidekick cloud account name mismatch")
	assert.Equal(t, "AwsSidekick", cloudAccount.Type, "a new AwsSidekick cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new AwsSidekick cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	response, err := c.V2.CloudAccounts.UpdateAwsSidekick(cloudAccount)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "arn:foo:bar", response.Data.Data.CrossAccountCreds.RoleArn)
	assert.Equal(t, "0123456789", response.Data.Data.CrossAccountCreds.ExternalID)
	assert.Equal(t, 24, response.Data.Data.ScanFrequency)
	assert.True(t, response.Data.Data.ScanContainers)
	assert.True(t, response.Data.Data.ScanHostVulnerabilities)
}

func singleAwsSidekickCloudAccount(id string) string {
	return `
  {
    "createdOrUpdatedBy": "teddy.reed@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
    "data": {
      "awsAccountId": "123456789000",
      "crossAccountCredentials": {
        "externalId": "0123456789",
        "roleArn": "arn:foo:bar"
      },
	  "scanFrequency": 24,
	  "scanContainers": true,
	  "scanHostVulnerabilities": true
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
    "type": "AwsSidekick"
  }
  `
}
