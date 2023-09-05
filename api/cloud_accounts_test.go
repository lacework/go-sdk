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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestCloudAccountTypeAwsCtsSqs(t *testing.T) {
	assert.Equal(t,
		"AwsCtSqs", api.AwsCtSqsCloudAccount.String(),
		"wrong cloud account type",
	)
}

func TestFindCloudAccountType(t *testing.T) {
	cloudFound, found := api.FindCloudAccountType("SOME_NON_EXISTING_INTEGRATION")
	assert.False(t, found, "cloud account type should not be found")
	assert.Equal(t, 0, int(cloudFound), "wrong cloud account type")
	assert.Equal(t, "None", cloudFound.String(), "wrong cloud account type")

	cloudFound, found = api.FindCloudAccountType("AwsCtSqs")
	assert.True(t, found, "cloud account type should exist")
	assert.Equal(t, "AwsCtSqs", cloudFound.String(), "wrong cloud account type")

	//cloudFound, found = api.FindCloudAccountType("GcpCfg")
	//assert.True(t, found, "cloud account type should exist")
	//assert.Equal(t, "GcpCfg", cloudFound.String(), "wrong cloud account type")
}

func TestCloudAccountsGet(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		vanillaInt  = singleVanillaCloudAccount(intgGUID, vanillaType, "")
		fakeServer  = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateCloudAccountResponse(vanillaInt))
			}
		},
	)

	fakeServer.MockAPI("CloudAccounts/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when cloud account exists", func(t *testing.T) {
		var response api.CloudAccountResponse
		err := c.V2.CloudAccounts.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.IntgGuid)
			assert.Equal(t, "integration_name", response.Data.Name)
			assert.Equal(t, "VANILLA", response.Data.Type)
		}
	})

	t.Run("when cloud account does NOT exist", func(t *testing.T) {
		var response api.CloudAccountResponse
		err := c.V2.CloudAccounts.Get("UNKNOWN_INTG_GUID", response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/CloudAccounts/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestCloudAccountsDelete(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("CloudAccounts/%s", intgGUID)
		vanillaInt  = singleVanillaCloudAccount(intgGUID, vanillaType, "")
		getResponse = generateCloudAccountResponse(vanillaInt)
		fakeServer  = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if getResponse != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, getResponse)
				case "DELETE":
					// once deleted, empty the getResponse so that
					// further GET requests return 404s
					getResponse = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	fakeServer.MockAPI("CloudAccounts/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify cloud account exists", func(t *testing.T) {
		var response api.CloudAccountResponse
		err := c.V2.CloudAccounts.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.IntgGuid)
			assert.Equal(t, "integration_name", response.Data.Name)
			assert.Equal(t, "VANILLA", response.Data.Type)
		}
	})

	t.Run("when cloud account has been deleted", func(t *testing.T) {
		err := c.V2.CloudAccounts.Delete(intgGUID)
		assert.Nil(t, err)

		var response api.CloudAccountResponse
		err = c.V2.CloudAccounts.Get(intgGUID, response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/CloudAccounts/MOCK_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestCloudAccountsList(t *testing.T) {
	var (
		awsIntgGUIDs        = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		awsCfgGUIDs         = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		awsEksAuditLogGUIDs = []string{intgguid.New()}
		azureIntgGUIDs      = []string{intgguid.New(), intgguid.New()}
		gcpIntgGUIDs        = []string{
			intgguid.New(), intgguid.New(), intgguid.New(), intgguid.New(),
		}
		allGUIDs    = append(awsEksAuditLogGUIDs, append(azureIntgGUIDs, append(awsCfgGUIDs, append(gcpIntgGUIDs, awsIntgGUIDs...)...)...)...)
		expectedLen = len(allGUIDs)
		fakeServer  = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("CloudAccounts",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			cloudAccounts := []string{
				generateCloudAccounts(awsIntgGUIDs, "AwsCtSqs"),
				generateCloudAccounts(awsEksAuditLogGUIDs, "AwsEksAudit"),
				// TODO @afiune come back here and update these Cloud Accounts types when they exist
				generateCloudAccounts(awsCfgGUIDs, "AwsCfg"),
				generateCloudAccounts(gcpIntgGUIDs, "AwsCtSqs"),   // "GcpCfg"),
				generateCloudAccounts(azureIntgGUIDs, "AwsCtSqs"), // "AzureAlSeq"),
			}
			fmt.Fprintf(w,
				generateCloudAccountsResponse(
					strings.Join(cloudAccounts, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.CloudAccounts.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.IntgGuid)
	}
}

func TestCloudAccountsListByType(t *testing.T) {
	var (
		awsIntgGUIDs = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		expectedLen  = len(awsIntgGUIDs)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("CloudAccounts/AwsCtSqs",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			cloudAccounts := []string{
				generateCloudAccounts(awsIntgGUIDs, "AwsCtSqs"),
			}
			fmt.Fprintf(w,
				generateCloudAccountsResponse(
					strings.Join(cloudAccounts, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	caType, _ := api.FindCloudAccountType("AwsCtSqs")
	response, err := c.V2.CloudAccounts.ListByType(caType)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, awsIntgGUIDs, d.IntgGuid)
	}
}

func TestCloudAccountMigrate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = "migrateGcpAtSes"
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "cloudAccountMigrate() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "props", "migration props are missing")
			assert.Contains(t, body, "migrate\":true",
				"migrate field is missing or it is set to false")
			assert.Contains(t, body, "migrationTimestamp", "migration timestamp is missing")
		}
		fmt.Fprintf(w, generateCloudAccountResponse(singleGcpAtCloudAccount(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	cloudAccount := api.NewCloudAccount("integration_name",
		api.GcpAtSesCloudAccount,
		api.GcpAtSesData{
			Credentials: api.GcpAtSesCredentials{
				ClientID:     "123456789",
				ClientEmail:  "test@project.iam.gserviceaccount.com",
				PrivateKeyID: "",
				PrivateKey:   "",
			},
		},
	)
	assert.Equal(t, "integration_name", cloudAccount.Name, "GcpAtSes cloud account name mismatch")
	assert.Equal(t, "GcpAtSes", cloudAccount.Type,
		"a new GcpAtSes cloud account should match its type")
	assert.Equal(t, 1, cloudAccount.Enabled, "a new GcpAtSes cloud account should be enabled")
	cloudAccount.IntgGuid = intgGUID

	err = c.V2.CloudAccounts.Migrate(intgGUID)
	assert.Nil(t, err)
}

func generateCloudAccounts(guids []string, iType string) string {
	cloudAccounts := make([]string, len(guids))
	for i, guid := range guids {
		switch iType {
		case api.AwsCtSqsCloudAccount.String():
			cloudAccounts[i] = singleAwsCtSqsCloudAccount(guid)
		case api.AwsEksAuditCloudAccount.String():
			cloudAccounts[i] = singleAwsEksAuditCloudAccount(guid)
		case api.AwsCfgCloudAccount.String():
			cloudAccounts[i] = singleAwsCfgCloudAccount(guid)
		}
	}
	return strings.Join(cloudAccounts, ", ")
}

func generateCloudAccountsResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateCloudAccountResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleVanillaCloudAccount(id string, iType string, data string) string {
	if data == "" {
		data = "{}"
	}
	return `
    {
      "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
      "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
      "data": ` + data + `,
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
      "type": "` + iType + `"
    }
	`
}
