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

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/capturer"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestSplitAzureSubscriptionsApiResponse(t *testing.T) {
	cases := []struct {
		subject  api.AzureConfigData
		expected cliComplianceAzureInfo
	}{
		// empty subscriptions will return empty cli info
		{
			api.AzureConfigData{},
			cliComplianceAzureInfo{Subscriptions: make([]cliComplianceIDAlias, 0)},
		},
		// real test case with NO alias
		{
			api.AzureConfigData{
				Tenant:        "ABCCC123-abc-123-AB12-XYZ987",
				Subscriptions: []string{"subscription-id-1", "subscription-id-2", "subscription-id-3", "subscription-id-4"},
			},
			cliComplianceAzureInfo{
				Tenant: cliComplianceIDAlias{"ABCCC123-abc-123-AB12-XYZ987", ""},
				Subscriptions: []cliComplianceIDAlias{
					{"subscription-id-1", ""},
					{"subscription-id-2", ""},
					{"subscription-id-3", ""},
					{"subscription-id-4", ""},
				},
			},
		},
		// real test case with alias
		{
			api.AzureConfigData{
				Tenant: "ABCCC123-abc-123-AB12-XYZ987 (cool.org.alias.example.com)",
				Subscriptions: []string{
					"id-1 (a test subscription)",
					"xmen-subscription (serious alias)",
					"disney-movies (Maybe Production)",
					"foo (bar)",
				},
			},
			cliComplianceAzureInfo{
				Tenant: cliComplianceIDAlias{"ABCCC123-abc-123-AB12-XYZ987", "cool.org.alias.example.com"},
				Subscriptions: []cliComplianceIDAlias{
					{"id-1", "a test subscription"},
					{"xmen-subscription", "serious alias"},
					{"disney-movies", "Maybe Production"},
					{"foo", "bar"},
				},
			},
		},
	}
	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equalf(t,
				kase.expected, splitAzureSubscriptionsApiResponse(kase.subject),
				"there is a problem with this test case, please check",
			)
		})
	}
}

func TestCliListAzureTenantsAndSubscriptionsWithoutData(t *testing.T) {
	cliOutput := capturer.CaptureOutput(func() {
		assert.Nil(t, cliListTenantsAndSubscriptions(api.CloudAccountsResponse{}))
	})
	assert.Contains(t, cliOutput, "There are no Azure Tenants configured in your account.")

	t.Run("test JSON output", func(t *testing.T) {
		cli.EnableJSONOutput()
		defer cli.EnableHumanOutput()
		cliJSONOutput := capturer.CaptureOutput(func() {
			assert.Nil(t, cliListTenantsAndSubscriptions(api.CloudAccountsResponse{}))
		})
		expectedJSON := `{
  "azure_subscriptions": []
}
`
		assert.Equal(t, expectedJSON, cliJSONOutput)
	})
}

func TestCliListAzureTenantsAndSubscriptionsWithData(t *testing.T) {
	var (
		fakeServer = lacework.MockServer()
		tenantID   = "abc123xy-1234-abcd-a1b2-09876zxy1234"
	)
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI(
		"CloudAccounts/AzureCfg",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, mockAzureSubsResponse(tenantID))
		})
	fakeServer.MockAPI(
		"Configs/AzureSubscriptions",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, mockAzureConfigsResponse())
		})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)

	assert.Nil(t, err)

	cli.LwApi = c
	cli.NonInteractive()
	defer func() {
		cli.LwApi = nil
		cli.Interactive()
	}()

	t.Run("enabled", func(t *testing.T) {
		cliOutput := capturer.CaptureOutput(func() {
			assert.Nil(t, cliListTenantsAndSubscriptions(mockAzureIntegrationsResponse(tenantID, 1)))
		})
		// NOTE (@afiune): We purposely leave trailing spaces in this table, we need them!
		expectedTable := `
              AZURE TENANT                        AZURE SUBSCRIPTION            STATUS   
---------------------------------------+--------------------------------------+----------
  abc123xy-1234-abcd-a1b2-09876zxy1234   ABC123XX-1234-ABCD-1234-ABCD1234XYZZ   Enabled  
`
		assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
	})

	t.Run("disabled", func(t *testing.T) {
		cliOutput := capturer.CaptureOutput(func() {
			assert.Nil(t, cliListTenantsAndSubscriptions(mockAzureIntegrationsResponse(tenantID, 0)))
		})
		// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
		expectedTable := `
              AZURE TENANT                        AZURE SUBSCRIPTION             STATUS   
---------------------------------------+--------------------------------------+-----------
  abc123xy-1234-abcd-a1b2-09876zxy1234   ABC123XX-1234-ABCD-1234-ABCD1234XYZZ   Disabled  
`
		assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
	})
}

func mockAzureIntegrationsResponse(tenantID string, enabled int) api.CloudAccountsResponse {
	var response = api.CloudAccountsResponse{}
	err := json.Unmarshal([]byte(`
{
    "data": [
        {
            "createdOrUpdatedBy": "test@lacework.net",
            "createdOrUpdatedTime": "2022-04-20T16:16:01.448Z",
            "enabled": `+strconv.Itoa(enabled)+`,
            "intgGuid": "EXAMPLE_12345",
            "isOrg": 0,
            "name": "Azure config",
            "state": {
                "ok": true,
                "lastUpdatedTime": 1643137889612,
                "lastSuccessfulTime": 1643137889612,
                "details": {
                    "tenantErrors": {
                        "opsDeniedAccess": []
                    },
                    "subscriptionErrors": {
                        "ABC123XX-1234-ABCD-1234-ABCD1234XYZZ": {
                            "opsDeniedAccess": []
                        }
                    }
                }
            },
            "type": "AzureCfg",
            "data": {
                "credentials": {
                    "clientId": "ABC123XX-1234-ABCD-1234-ABCD1234XYZZ",
                    "clientSecret": ""
                },
                "tenantId": "`+tenantID+`"
            }
        }
]
}
`), &response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func mockAzureSubsResponse(tenantID string) string {
	return `{
  "data": [
    {
      "subscriptions": [
        "ABC123XX-1234-ABCD-1234-ABCD1234XYZZ (Default-account)"
      ],
      "tenant": "` + tenantID + ` (Default Directory)"
    }
  ],
  "message": "SUCCESS",
  "ok": true
}`
}

func mockAzureConfigsResponse() string {
	return `{
    "data": [
        {
            "tenant": "abc123xy-1234-abcd-a1b2-09876zxy1234",
            "subscriptions": [
                "ABC123XX-1234-ABCD-1234-ABCD1234XYZZ (Default-account)"
            ]
        }
    ]
}`
}
