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

func TestCloudAccountsAwsEksAuditGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsEksAudit() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsEksAuditCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAwsEksAudit(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "arn:foo:bar", response.Data.Data.Credentials.RoleArn)
	assert.Equal(t, "0123456789", response.Data.Data.Credentials.ExternalID)
	assert.Equal(
		t,
		"arn:aws:sns:us-west-2:0123456789:foo-lacework-eks:00777777-ab77-1234-a123-a12ab1d12c1d",
		response.Data.Data.SnsArn,
	)
	assert.Equal(t, "arn:aws:s3:::example-bucket-name", response.Data.Data.S3BucketArn)
}

func TestCloudAccountsAwsEksAuditUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsEksAudit() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "AwsEksAudit", "wrong cloud account type")
			assert.Contains(t, body, "arn:bubu:lubu", "wrong role arn")
			assert.Contains(t, body, "abc123", "wrong external ID")
			assert.Contains(
				t,
				body,
				"arn:aws:sns:us-west-2:0123456789:foo-lacework-eks:00777777-ab77-1234-a123-a12ab1d12c1d",
				"wrong sns arn")

			assert.Contains(t, body, "arn:aws:s3:::example-bucket-name", "wrong s3 bucket arn")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsEksAuditCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.AwsEksAuditCloudAccount,
		api.AwsEksAuditData{
			SnsArn:      "arn:aws:sns:us-west-2:0123456789:foo-lacework-eks:00777777-ab77-1234-a123-a12ab1d12c1d",
			S3BucketArn: "arn:aws:s3:::example-bucket-name",
			Credentials: api.AwsEksAuditCredentials{
				RoleArn:    "arn:bubu:lubu",
				ExternalID: "abc123",
			},
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "AwsEksAudit cloud account name mismatch")
	assert.Equal(t, "AwsEksAudit", cloudAccount.Type, "a new AwsEksAudit cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new AwsEksAudit cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	response, err := c.V2.CloudAccounts.UpdateAwsEksAudit(cloudAccount)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t,
		"arn:aws:sns:us-west-2:0123456789:foo-lacework-eks:00777777-ab77-1234-a123-a12ab1d12c1d",
		response.Data.Data.SnsArn)
	assert.Equal(t, "arn:aws:s3:::example-bucket-name", response.Data.Data.S3BucketArn)
}

func singleAwsEksAuditCloudAccount(id string) string {
	return `
  {
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
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
	"type": "AwsEksAudit",
    "data": {
      "snsArn": "arn:aws:sns:us-west-2:0123456789:foo-lacework-eks:00777777-ab77-1234-a123-a12ab1d12c1d",
	  "s3BucketArn": "arn:aws:s3:::example-bucket-name",
      "crossAccountCredentials": {
        "externalId": "0123456789",
        "roleArn": "arn:foo:bar"
      }
    }
  }
  `
}
