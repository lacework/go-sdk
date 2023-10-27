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

func TestCloudAccountsNewAwsSidekickOrgWithCustomTemplateFile(t *testing.T) {
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
	awsSidekickOrgData := api.AwsSidekickOrgData{
		CrossAccountCreds: api.AwsSidekickCrossAccountCredentials{
			RoleArn:    "arn:foo:bar",
			ExternalID: "0123456789",
		},
	}
	awsSidekickOrgData.EncodeAccountMappingFile(accountMappingJSON)
	subject := api.NewCloudAccount("integration_name", api.AwsSidekickOrgCloudAccount, awsSidekickOrgData)
	assert.Equal(t, api.AwsSidekickOrgCloudAccount.String(), subject.Type)

	// casting the data interface{} to type AwsCfgData
	subjectData := subject.Data.(api.AwsSidekickOrgData)

	assert.Equal(t, subjectData.CrossAccountCreds.RoleArn, "arn:foo:bar")
	assert.Equal(t, subjectData.CrossAccountCreds.ExternalID, "0123456789")
	assert.Contains(t,
		subjectData.AccountMappingFile,
		"data:application/json;name=i.json;base64,",
		"check the custom_template_file encoder",
	)
	accountMapping, err := subjectData.DecodeAccountMappingFile()
	assert.Nil(t, err)
	assert.Equal(t, accountMappingJSON, accountMapping)

	// When there is no custom account mapping file, this function should
	// return an empty string to match the pattern
	subjectData.AccountMappingFile = ""
	accountMapping, err = subjectData.DecodeAccountMappingFile()
	assert.Nil(t, err)
	assert.Empty(t, accountMapping)
}

func TestCloudAccountsAwsSidekickOrgGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsSidekickOrg() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsSidekickOrgCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAwsSidekickOrg(intgGUID)
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
	assert.Equal(t, "123456789000", response.Data.Data.ScanningAccount)
	assert.Equal(t, "000123456789", response.Data.Data.ManagementAccount)
	assert.Equal(t, "r-1234, ou-0987", response.Data.Data.MonitoredAccounts)

}

func TestCloudAccountsAwsSidekickOrgUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsSidekickOrg() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "AwsSidekickOrg", "wrong cloud account type")
			assert.Contains(t, body, "arn:foo:bar", "wrong role arn")
			assert.Contains(t, body, "0123456789", "wrong external ID")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsSidekickOrgCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.AwsSidekickOrgCloudAccount,
		api.AwsSidekickOrgData{
			CrossAccountCreds: api.AwsSidekickCrossAccountCredentials{
				RoleArn:    "arn:foo:bar",
				ExternalID: "0123456789",
			},
			ScanFrequency:     24,
			MonitoredAccounts: "r-1234, ou-0987",
			ScanningAccount:   "123456789000",
			ManagementAccount: "000123456789",
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "AwsSidekickOrg cloud account name mismatch")
	assert.Equal(t, "AwsSidekickOrg", cloudAccount.Type, "a new AwsSidekickOrg cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new AwsSidekickOrg cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	response, err := c.V2.CloudAccounts.UpdateAwsSidekickOrg(cloudAccount)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "arn:foo:bar", response.Data.Data.CrossAccountCreds.RoleArn)
	assert.Equal(t, "0123456789", response.Data.Data.CrossAccountCreds.ExternalID)
	assert.Equal(t, 24, response.Data.Data.ScanFrequency)
	assert.True(t, response.Data.Data.ScanContainers)
	assert.True(t, response.Data.Data.ScanHostVulnerabilities)
	assert.Equal(t, "123456789000", response.Data.Data.ScanningAccount)
	assert.Equal(t, "000123456789", response.Data.Data.ManagementAccount)
	assert.Equal(t, "r-1234, ou-0987", response.Data.Data.MonitoredAccounts)
}

func singleAwsSidekickOrgCloudAccount(id string) string {
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
	  "scanHostVulnerabilities": true,
	  "scanShortLivedInstances": false,
	  "scanStoppedInstances": true,
	  "scanMultiVolume": false,
	  "managementAccount": "000123456789",
	  "monitoredAccounts": "r-1234, ou-0987",
	  "scanningAccount": "123456789000"
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
    "type": "AwsSidekickOrg"
  }
  `
}
