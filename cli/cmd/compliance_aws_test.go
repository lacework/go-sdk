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
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestCliListAwsAccountsWithNoAccounts(t *testing.T) {
	cliOutput := captureOutput(func() {
		assert.Nil(t, cliListAwsAccounts(new(api.AwsIntegrationsResponse)))
	})
	assert.Contains(t, cliOutput, "There are no AWS accounts configured in your account.")

	t.Run("test JSON output", func(t *testing.T) {
		cli.EnableJSONOutput()
		defer cli.EnableHumanOutput()
		cliJSONOutput := captureOutput(func() {
			assert.Nil(t, cliListAwsAccounts(new(api.AwsIntegrationsResponse)))
		})
		expectedJSON := `{
  "aws_accounts": []
}
`
		assert.Equal(t, expectedJSON, cliJSONOutput)
	})
}

func TestCliListAwsAccountsWithAccountsEnabled(t *testing.T) {
	cliOutput := captureOutput(func() {
		assert.Nil(t, cliListAwsAccounts(mockAwsIntegrationsResponse(1, 1)))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
  AWS ACCOUNT    STATUS   
---------------+----------
  123456789012   Enabled  
  098765432109   Enabled  
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func TestCliListAwsAccountsWithAccountsDisabled(t *testing.T) {
	cliOutput := captureOutput(func() {
		assert.Nil(t, cliListAwsAccounts(mockAwsIntegrationsResponse(0, 0)))
	})
	// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
	expectedTable := `
  AWS ACCOUNT     STATUS   
---------------+-----------
  123456789012   Disabled  
  098765432109   Disabled  
`
	assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
}

func mockAwsIntegrationsResponse(acc1Enabled, acc2Enabled int) *api.AwsIntegrationsResponse {
	response := &api.AwsIntegrationsResponse{}
	err := json.Unmarshal([]byte(`{
  "data": [
    {
      "CREATED_OR_UPDATED_BY": "salim.afiunemaya@lacework.net",
      "CREATED_OR_UPDATED_TIME": "2021-10-25T15:18:47.945Z",
      "DATA": {
        "AWS_ACCOUNT_ID": "123456789012"
      },
      "ENABLED": `+strconv.Itoa(acc1Enabled)+`,
      "INTG_GUID": "MOCK_1233",
      "IS_ORG": 0,
      "NAME": "TF config",
      "STATE": {
        "lastSuccessfulTime": "2022-Jan-31 16:08:54 UTC",
        "lastUpdatedTime": "2022-Jan-31 16:08:54 UTC",
        "ok": true
      },
      "TYPE": "AWS_CFG",
      "TYPE_NAME": "AWS Config"
    },
    {
      "CREATED_OR_UPDATED_BY": "vatasha.white@lacework.net",
      "CREATED_OR_UPDATED_TIME": "2022-01-13T20:32:59.954Z",
      "DATA": {
        "AWS_ACCOUNT_ID": "098765432109"
      },
      "ENABLED": `+strconv.Itoa(acc2Enabled)+`,
      "INTG_GUID": "MOCK_1234",
      "IS_ORG": 0,
      "NAME": "TF config",
      "STATE": {
        "lastSuccessfulTime": "2022-Jan-31 16:47:04 UTC",
        "lastUpdatedTime": "2022-Jan-31 16:47:04 UTC",
        "ok": true
      },
      "TYPE": "AWS_CFG",
      "TYPE_NAME": "AWS Config"
    }
  ],
  "message": "SUCCESS",
  "ok": true
}
`), response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
