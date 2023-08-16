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

func TestCloudAccountsAwsUsGovCtSqsGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsUsGovCtSqs() should be a GET method")
		fmt.Fprintf(w, generateCloudAccountResponse(singleAwsUsGovCtSqsCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.GetAwsUsGovCtSqs(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "awsgov-config", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "123456789", response.Data.Data.Credentials.AwsAccountID)
	assert.Equal(t, "ABCDEFGHIJ", response.Data.Data.Credentials.AccessKeyID)
	assert.Equal(t, "AwsUsGovCtSqs", response.Data.Type)
	assert.Equal(t, "https://sqs.us-gov-east-1.amazonaws.com/1234567/test", response.Data.Data.QueueUrl)
}

func singleAwsUsGovCtSqsCloudAccount(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "test@lacework.net",
        "createdOrUpdatedTime": "2022-06-27T19:28:45.600Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 0,
        "name": "awsgov-config",
        "state": {
            "ok": true,
            "lastUpdatedTime": 1663061236782,
            "lastSuccessfulTime": 1663061236782,
            "details": {
                "sqsRx": "AccessDenied",
                "decodeNtfn": "OK",
                "s3Get": "OK",
                "sqsDel": "OK",
                "lastMsgRxTime": 1656339585273,
                "noData": false
            }
        },
        "type": "AwsUsGovCtSqs",
        "data": {
            "accessKeyCredentials": {
                "accountId": "123456789",
                "accessKeyId": "ABCDEFGHIJ",
                "secretAccessKey": ""
            },
		"queueUrl": "https://sqs.us-gov-east-1.amazonaws.com/1234567/test"
        }
}`, id)
}
