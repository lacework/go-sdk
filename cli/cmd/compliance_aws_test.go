//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/capturer"
)

func TestCliListAwsAccountsWithNoAccounts(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, cliListAwsAccounts(api.CloudAccountsResponse{}))
	})
	assert.Contains(t, cliOutput, "There are no AWS accounts configured in your account.")

	t.Run("test JSON output", func(t *testing.T) {
		cli.EnableJSONOutput()
		defer cli.EnableHumanOutput()
		cliJSONOutput := capturer.CaptureOutput(func() {
			assert.Nil(t, cliListAwsAccounts(api.CloudAccountsResponse{}))
		})
		expectedJSON := `{
  "aws_accounts": []
}
`
		assert.Equal(t, expectedJSON, cliJSONOutput)
	})
}

func TestCliListAwsAccountsWithAccountsEnabled(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, cliListAwsAccounts(mockAwsIntegrationsResponse(1, 1)))
	})
	// NOTE (@afiune): We purposely leave trailing spaces in this table, we need them!
	expectedTable := `
  AWS ACCOUNT    STATUS   
---------------+----------
  123456789012   Enabled  
  098765432109   Enabled  
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestCliListAwsAccountsWithAccountsDisabled(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, cliListAwsAccounts(mockAwsIntegrationsResponse(0, 0)))
	})
	// NOTE (@afiune): We purposely leave trailing spaces in this table, we need them!
	expectedTable := `
  AWS ACCOUNT     STATUS   
---------------+-----------
  123456789012   Disabled  
  098765432109   Disabled  
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func mockAwsIntegrationsResponse(acc1Enabled, acc2Enabled int) api.CloudAccountsResponse {
	var response api.CloudAccountsResponse
	jsonString := fmt.Sprintf(`{
  "data": [
           {
            "createdOrUpdatedBy": "darren.murray@lacework.net",
            "createdOrUpdatedTime": "2022-10-10T12:08:54.629Z",
            "enabled": %d,
            "intgGuid": "TECHALLY_123ABC",
            "isOrg": 0,
            "name": "TF config",
            "state": {
                "ok": true,
                "lastUpdatedTime": 1669136762967,
                "lastSuccessfulTime": 1669136762967,
                "details": {
                    "complianceOpsDeniedAccess": [
                        "GetBucketLogging"
                    ]
                }
            },
            "type": "AwsCfg",
            "data": {
                "crossAccountCredentials": {
                    "roleArn": "arn:aws:iam::123456789012:role/test",
                    "externalId": "example"
                },
                "awsAccountId": "123456789012"
            }
        },
           {
            "createdOrUpdatedBy": "darren.murray@lacework.net",
            "createdOrUpdatedTime": "2022-10-10T12:08:54.629Z",
            "enabled": %d,
            "intgGuid": "TECHALLY_123ABC",
            "isOrg": 0,
            "name": "TF config",
            "state": {
                "ok": true,
                "lastUpdatedTime": 1669136762967,
                "lastSuccessfulTime": 1669136762967,
                "details": {
                    "complianceOpsDeniedAccess": [
                        "GetBucketLogging"
                    ]
                }
            },
            "type": "AwsCfg",
            "data": {
                "crossAccountCredentials": {
                    "roleArn": "arn:aws:iam::098765432109:role/test",
                    "externalId": "example"
                },
                "awsAccountId": "098765432109"
            }
        }
]
}`, acc1Enabled, acc2Enabled)

	err := json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
