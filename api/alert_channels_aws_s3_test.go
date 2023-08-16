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

func TestAlertChannelsGetAwsS3(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsS3() should be a GET method")
		fmt.Fprintf(w, generateAlertChannelResponse(singleAwsS3AlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.GetAwsS3(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.Credentials.BucketArn, "arn:aws:s3:::data-export")
	assert.Contains(t, response.Data.Data.Credentials.ExternalID, "abc123")
	assert.Contains(t, response.Data.Data.Credentials.RoleArn, "arn:aws:iam::123456789012:role/lw-s3-export")
}

func TestAlertChannelsAwsS3Update(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsS3() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "cloud account name is missing")
			assert.Contains(t, body, "AwsS3", "wrong cloud account type")
			assert.Contains(t, body, "arn:aws:s3:::data-export", "missing bucket arn")
			assert.Contains(t, body, "abc123", "missing external id")
			assert.Contains(t, body, "arn:aws:iam::123456789012:role/lw-s3-expor", "missing role arn")
			assert.Contains(t, body, "enabled\":1", "cloud account is not enabled")
		}

		fmt.Fprintf(w, generateAlertChannelResponse(singleAwsS3AlertChannel(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	emailAlertChan := api.NewAlertChannel("integration_name",
		api.AwsS3AlertChannelType,
		api.AwsS3DataV2{
			Credentials: api.AwsS3Credentials{
				RoleArn:    "arn:aws:iam::123456789012:role/lw-s3-export",
				ExternalID: "abc123",
				BucketArn:  "arn:aws:s3:::data-export",
			},
		},
	)
	assert.Equal(t, "integration_name", emailAlertChan.Name, "AwsS3 cloud account name mismatch")
	assert.Equal(t, "AwsS3", emailAlertChan.Type, "a new AwsS3 cloud account should match its type")
	assert.Equal(t, 1, emailAlertChan.Enabled, "a new AwsS3 cloud account should be enabled")
	emailAlertChan.IntgGuid = intgGUID

	response, err := c.V2.AlertChannels.UpdateAwsS3(emailAlertChan)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.True(t, response.Data.State.Ok)
	assert.Contains(t, response.Data.Data.Credentials.BucketArn, "arn:aws:s3:::data-export")
	assert.Contains(t, response.Data.Data.Credentials.ExternalID, "abc123")
	assert.Contains(t, response.Data.Data.Credentials.RoleArn, "arn:aws:iam::123456789012:role/lw-s3-export")
}

func singleAwsS3AlertChannel(id string) string {
	return `
{
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T18:10:40.745Z",
    "data": {
      "s3CrossAccountCredentials": {
        "bucketArn": "arn:aws:s3:::data-export",
        "externalId": "abc123",
        "roleArn": "arn:aws:iam::123456789012:role/lw-s3-export"
      }
    },
    "enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "state": {
      "details": {},
      "lastSuccessfulTime": 1627895573122,
      "lastUpdatedTime": 1627895573122,
      "ok": true
    },
    "type": "AwsS3"
}
  `
}
