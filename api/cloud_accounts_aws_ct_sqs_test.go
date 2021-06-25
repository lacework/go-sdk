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

func TestCloudAccountsNewAwsCtSqsWithCustomTemplateFile(t *testing.T) {
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
	awsCtSqsData := api.AwsCtSqsData{
		QueueUrl: "https://sqs.us-west-2.amazonaws.com/123456789000/lw",
		Credentials: api.AwsCtSqsCredentials{
			RoleArn:    "arn:foo:bar",
			ExternalID: "0123456789",
		},
	}
	awsCtSqsData.EncodeAccountMappingFile(accountMappingJSON)

	subject := api.NewCloudAccount("integration_name", api.AwsCtSqsCloudAccount, awsCtSqsData)
	assert.Equal(t, api.AwsCtSqsCloudAccount.String(), subject.Type)

	// casting the data interface{} to type AwsCtSqsData
	subjectData := subject.Data.(api.AwsCtSqsData)

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

func TestCloudAccountsAwsCtSqsGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsCtSqs() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsCtSqsCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAwsCtSqs(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "arn:foo:bar", response.Data.Data.Credentials.RoleArn)
	assert.Equal(t, "0123456789", response.Data.Data.Credentials.ExternalID)
	assert.Equal(t, "https://sqs.us-west-2.amazonaws.com/123456789000/lw", response.Data.Data.QueueUrl)
}

func TestCloudAccountsAwsCtSqsUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsCtSqs() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "AwsCtSqs", "wrong cloud account type")
			assert.Contains(t, body, "arn:bubu:lubu", "wrong role arn")
			assert.Contains(t, body, "abc123", "wrong external ID")
			assert.Contains(t, body, "https://sqs.us-west-2.amazonaws.com/123456789000/lw", "wrong queue url")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsCtSqsCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.AwsCtSqsCloudAccount,
		api.AwsCtSqsData{
			QueueUrl: "https://sqs.us-west-2.amazonaws.com/123456789000/lw",
			Credentials: api.AwsCtSqsCredentials{
				RoleArn:    "arn:bubu:lubu",
				ExternalID: "abc123",
			},
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "AwsCtSqs cloud account name mismatch")
	assert.Equal(t, "AwsCtSqs", cloudAccount.Type, "a new AwsCtSqs cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new AwsCtSqs cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	response, err := c.V2.CloudAccounts.UpdateAwsCtSqs(cloudAccount)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t,
		"https://sqs.us-west-2.amazonaws.com/123456789000/lw",
		response.Data.Data.QueueUrl)
}

func singleAwsCtSqsCloudAccount(id string) string {
	return `
  {
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
    "data": {
      "awsAccountId": "123456789000",
      "queueUrl": "https://sqs.us-west-2.amazonaws.com/123456789000/lw",
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
      "details": {
        "complianceOpsDeniedAccess": [
          "GetBucketAcl",
          "GetBucketLogging"
        ]
      },
      "lastSuccessfulTime": 1624456896915,
      "lastUpdatedTime": 1624456896915,
      "ok": true
    },
    "type": "AwsCtSqs"
  }
  `
}
