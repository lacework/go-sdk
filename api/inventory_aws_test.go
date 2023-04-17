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
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestInventoryAwsSearch(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Inventory/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockInventoryAwsResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response := api.InventoryAwsResponse{}
	err = c.V2.Inventory.Search(&response, api.InventorySearch{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "my-example-id", response.Data[0].ResourceId)
	}
}

func TestInventoryAwsScan(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Inventory/scan",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Scan() should be a POST method")
			fmt.Fprintf(w, mockInventoryScanResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Inventory.Scan(api.AwsInventoryType)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Scan has been requested", response.Data.Details)
	assert.Equal(t, "scanning", response.Data.Status)
}

func mockInventoryAwsResponse() string {
	return `
{
  "data": [
        {
            "apiKey": "ABCD",
            "cloudDetails": {
                "accountAlias": "tech-ally",
                "accountID": "0123456789"
            },
            "csp": "AWS",
            "endTime": "2022-08-22T20:00:00.000Z",
            "resourceConfig": {
                "KmsKeyId": "arn:aws:kms:region:012345:example",
                "LogFileValidationEnabled": 1,
                "Name": "my-example-bucket",
                "S3BucketName": "my-example-bucket",
                "SnsTopicARN": "arn:aws:sns:region:012345:example",
                "SnsTopicName": "example-sns",
                "TrailARN": "arn:aws:cloudtrail::example"
            },
            "resourceId": "my-example-id",
            "resourceRegion": "us-east-1",
            "resourceTags": {},
            "resourceType": "cloudtrail:trail",
            "service": "cloudtrail",
            "startTime": "2022-08-22T19:00:00.000Z",
            "status": {
                "formatVersion": 2,
                "props": {},
                "status": "success"
            },
            "urn": "arn:aws:cloudtrail::example"
        }
  ],
  "paging": {
    "rows": 1,
    "totalRows": 1,
    "urls": {
      "nextPage": null
    }
  }
}
	`
}

func mockInventoryScanResponse() string {
	return `
{
  "data": {
	"status": "scanning",
	"details": "Scan has been requested"
	}
}
	`
}
