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
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestSplitAzureSubscriptionsApiResponse(t *testing.T) {
	cases := []struct {
		subject  api.CompAzureSubscriptions
		expected cliComplianceAzureInfo
	}{
		// empty subscriptions will return empty cli info
		{
			api.CompAzureSubscriptions{},
			cliComplianceAzureInfo{Subscriptions: make([]cliComplianceIDAlias, 0)},
		},
		// real test case with NO alias
		{
			api.CompAzureSubscriptions{
				Tenant:        "ABCCC123-abc-123-AB12-XYZ987",
				Subscriptions: []string{"subscription-id-1", "subscription-id-2", "subscription-id-3", "subscription-id-4"},
			},
			cliComplianceAzureInfo{
				Tenant: cliComplianceIDAlias{"ABCCC123-abc-123-AB12-XYZ987", ""},
				Subscriptions: []cliComplianceIDAlias{
					cliComplianceIDAlias{"subscription-id-1", ""},
					cliComplianceIDAlias{"subscription-id-2", ""},
					cliComplianceIDAlias{"subscription-id-3", ""},
					cliComplianceIDAlias{"subscription-id-4", ""},
				},
			},
		},
		// real test case with alias
		{
			api.CompAzureSubscriptions{
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
					cliComplianceIDAlias{"id-1", "a test subscription"},
					cliComplianceIDAlias{"xmen-subscription", "serious alias"},
					cliComplianceIDAlias{"disney-movies", "Maybe Production"},
					cliComplianceIDAlias{"foo", "bar"},
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
	cliOutput := captureOutput(func() {
		assert.Nil(t, cliListTenantsAndSubscriptions(new(api.AzureIntegrationsResponse)))
	})
	assert.Contains(t, cliOutput, "There are no Azure Tenants configured in your account.")

	t.Run("test JSON output", func(t *testing.T) {
		cli.EnableJSONOutput()
		defer cli.EnableHumanOutput()
		cliJSONOutput := captureOutput(func() {
			assert.Nil(t, cliListTenantsAndSubscriptions(new(api.AzureIntegrationsResponse)))
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
		"external/compliance/azure/ListSubscriptionsForTenant",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, mockAzureSubsResponse(tenantID))
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
		cliOutput := captureOutput(func() {
			assert.Nil(t, cliListTenantsAndSubscriptions(mockAzureIntegrationsResponse(tenantID, 1)))
		})
		// NOTE (@afiune): We purposly leave trailing spaces in this table, we need them!
		expectedTable := `
              AZURE TENANT                        AZURE SUBSCRIPTION            STATUS   
---------------------------------------+--------------------------------------+----------
  abc123xy-1234-abcd-a1b2-09876zxy1234   ABC123XX-1234-ABCD-1234-ABCD1234XYZZ   Enabled  
`
		assert.Equal(t, strings.TrimPrefix(expectedTable, "\n"), cliOutput)
	})

	t.Run("disabled", func(t *testing.T) {
		cliOutput := captureOutput(func() {
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

func mockAzureIntegrationsResponse(tenantID string, enabled int) *api.AzureIntegrationsResponse {
	response := &api.AzureIntegrationsResponse{}
	err := json.Unmarshal([]byte(`{
  "data": [
	    {
      "CREATED_OR_UPDATED_BY": "salim.afiune-maya@lacework.net",
      "CREATED_OR_UPDATED_TIME": "2021-08-02T17:53:24.116Z",
      "DATA": {
        "TENANT_ID": "`+tenantID+`"
      },
      "ENABLED": `+strconv.Itoa(enabled)+`,
      "INTG_GUID": "MOCK_1234",
      "IS_ORG": 0,
      "NAME": "TF Config",
      "STATE": {
        "lastSuccessfulTime": "2021-Jun-04 09:40:39 UTC",
        "lastUpdatedTime": "2022-Jan-31 11:49:09 UTC",
        "ok": false
      },
      "TYPE": "AZURE_CFG",
      "TYPE_NAME": "Azure Compliance"
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
