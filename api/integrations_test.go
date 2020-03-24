//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

func TestIntegrationTypeAwsCfgIntegration(t *testing.T) {
	assert.Equal(t,
		"AWS_CFG", api.AwsCfgIntegration.String(),
		"wrong integration type",
	)
}

func TestIntegrationTypeAwsCloudTrailIntegration(t *testing.T) {
	assert.Equal(t,
		"AWS_CT_SQS", api.AwsCloudTrailIntegration.String(),
		"wrong integration type",
	)
}

func TestIntegrationTypeGcpCfgIntegration(t *testing.T) {
	assert.Equal(t,
		"GCP_CFG", api.GcpCfgIntegration.String(),
		"wrong integration type",
	)
}

func TestIntegrationTypeGcpAuditLogIntegration(t *testing.T) {
	assert.Equal(t,
		"GCP_AT_SES", api.GcpAuditLogIntegration.String(),
		"wrong integration type",
	)
}

func TestIntegrationTypeAzureCfgIntegration(t *testing.T) {
	assert.Equal(t,
		"AZURE_CFG", api.AzureCfgIntegration.String(),
		"wrong integration type",
	)
}

func TestIntegrationTypeAzureActivityLogIntegration(t *testing.T) {
	assert.Equal(t,
		"AZURE_AL_SEQ", api.AzureActivityLogIntegration.String(),
		"wrong integration type",
	)
}

func TestFindIntegrationType(t *testing.T) {
	typeFound, found := api.FindIntegrationType("SOME_NON_EXISTING_INTEGRATION")
	assert.False(t, found, "integration type should not be found")
	assert.Equal(t, 0, int(typeFound), "wrong integration type")
	assert.Equal(t, "NONE", typeFound.String(), "wrong integration type")

	typeFound, found = api.FindIntegrationType("AZURE_AL_SEQ")
	assert.True(t, found, "integration type should exist")
	assert.Equal(t, "AZURE_AL_SEQ", typeFound.String(), "wrong integration type")

	typeFound, found = api.FindIntegrationType("GCP_CFG")
	assert.True(t, found, "integration type should exist")
	assert.Equal(t, "GCP_CFG", typeFound.String(), "wrong integration type")
}

func TestIntegrationsGetSchema(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI("external/integrations/schema/AWS_CFG",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetSchema should be a GET method")
			fmt.Fprintf(w, "{\"any\":\"data\"}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetSchema(api.AwsCfgIntegration)
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.Equal(t, "data", response["any"])
	}

	response, err = c.Integrations.GetSchema(api.NoneIntegration)
	assert.Nil(t, response)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404 page not found")
	}
}

func TestIntegrationsGet(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("external/integrations/%s", intgGUID)
		vanillaInt  = singleVanillaIntegration(intgGUID, vanillaType, "")
		fakeServer  = lacework.MockServer()
	)
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateIntegrationsResponse(vanillaInt))
			}
		},
	)

	fakeServer.MockAPI("external/integrations/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "Not Found", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when integration exists", func(t *testing.T) {
		response, err := c.Integrations.Get(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			if assert.Equal(t, 1, len(response.Data)) {
				resData := response.Data[0]
				assert.Equal(t, intgGUID, resData.IntgGuid)
				assert.Equal(t, "integration_name", resData.Name)
				assert.Equal(t, "VANILLA", resData.Type)
				assert.Equal(t, "Vanilla Integration", resData.TypeName)
			}
		}
	})

	t.Run("when integration does NOT exist", func(t *testing.T) {
		response, err := c.Integrations.Get("UNKNOWN_INTG_GUID")
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v1/external/integrations/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "404 Not Found")
		}
	})
}

func TestIntegrationsDelete(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("external/integrations/%s", intgGUID)
		vanillaInt  = singleVanillaIntegration(intgGUID, vanillaType, "")
		// TODO @afiune revisit this with https://github.com/lacework/go-sdk/issues/23
		// this will change when the test hits DELETE
		getResponse = generateIntegrationsResponse(vanillaInt)
		fakeServer  = lacework.MockServer()
	)
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			//if getResponse != "" {
			switch r.Method {
			case "GET":
				fmt.Fprintf(w, getResponse)
			case "DELETE":
				fmt.Fprintf(w, getResponse)
				// once deleted, empty the getResponse so that
				// further GET requests return 404s
				//getResponse = ""
			}
			//} else {
			//http.Error(w, "Not Found", 404)
			//}
		},
	)

	fakeServer.MockAPI("external/integrations/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "Not Found", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify integration exists", func(t *testing.T) {
		response, err := c.Integrations.Get(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			if assert.Equal(t, 1, len(response.Data)) {
				resData := response.Data[0]
				assert.Equal(t, intgGUID, resData.IntgGuid)
				assert.Equal(t, "integration_name", resData.Name)
				assert.Equal(t, "VANILLA", resData.Type)
				assert.Equal(t, "Vanilla Integration", resData.TypeName)
			}
		}
	})

	t.Run("when integration has been deleted", func(t *testing.T) {
		response, err := c.Integrations.Delete(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			if assert.Equal(t, 1, len(response.Data)) {
				resData := response.Data[0]
				assert.Equal(t, intgGUID, resData.IntgGuid)
				assert.Equal(t, "integration_name", resData.Name)
				assert.Equal(t, "VANILLA", resData.Type)
				assert.Equal(t, "Vanilla Integration", resData.TypeName)
			}
		}
		//response, err = c.Integrations.Get(intgGUID)
		//assert.Empty(t, response)
		//if assert.NotNil(t, err) {
		//assert.Contains(t, err.Error(), "api/v1/external/integrations/MOCK_")
		//assert.Contains(t, err.Error(), "404 Not Found")
		//}
	})
}

func TestIntegrationsList(t *testing.T) {
	var (
		awsIntgGUIDs   = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		azureIntgGUIDs = []string{intgguid.New(), intgguid.New()}
		gcpIntgGUIDs   = []string{
			intgguid.New(), intgguid.New(), intgguid.New(), intgguid.New(),
		}
		allGUIDs    = append(azureIntgGUIDs, append(gcpIntgGUIDs, awsIntgGUIDs...)...)
		expectedLen = len(allGUIDs)
		fakeServer  = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			integrations := []string{
				generateIntegrations(awsIntgGUIDs, "AWS_CFG"),
				generateIntegrations(gcpIntgGUIDs, "GCP_CFG"),
				generateIntegrations(azureIntgGUIDs, "AZURE_CFG"),
			}
			fmt.Fprintf(w,
				generateIntegrationsResponse(
					strings.Join(integrations, ", "),
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

	response, err := c.Integrations.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.IntgGuid)
	}
}

func generateIntegrations(guids []string, iType string) string {
	integrations := make([]string, len(guids))
	for i, guid := range guids {
		switch iType {
		case api.AwsCfgIntegration.String():
			integrations[i] = singleAwsIntegration(guid)
		case api.AzureCfgIntegration.String():
			integrations[i] = singleAzureIntegration(guid)
		case api.GcpCfgIntegration.String():
			integrations[i] = singleGcpIntegration(guid)
		}
	}
	return strings.Join(integrations, ", ")
}

func generateIntegrationsResponse(data string) string {
	return `
		{
			"data": [` + data + `],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func singleVanillaIntegration(id string, iType string, data string) string {
	if data == "" {
		data = "{}"
	}
	return `
		{
			"INTG_GUID": "` + id + `",
			"NAME": "integration_name",
			"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
			"CREATED_OR_UPDATED_BY": "user@email.com",
			"TYPE": "` + iType + `",
			"STATE": {
				"ok": true,
				"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
				"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
			},
			"IS_ORG": 0,
			"TYPE_NAME": "Vanilla Integration",
			"DATA": ` + data + `
		}
	`
}
