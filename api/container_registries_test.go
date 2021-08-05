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

func TestNewContainerRegistry(t *testing.T) {
	ghcr := api.NewContainerRegistry("foo",
		api.GhcrContainerRegistry,
		api.GhcrData{
			Credentials: api.GhcrCredentials{
				Username: "bubulubu",
				Password: "secret",
				Ssl:      true,
			},
		},
	)

	assert.Equal(t, "ContVulnCfg", ghcr.Type)
	assert.Equal(t,
		api.GhcrContainerRegistry, ghcr.ContainerRegistryType(),
		"wrong container registry type",
	)
}

func TestContainerRegistryTypeGcpGar(t *testing.T) {
	assert.Equal(t,
		"GCP_GAR", api.GcpGarContainerRegistry.String(),
		"wrong container registry type",
	)
}

func TestFindContainerRegistryType(t *testing.T) {
	registryFound, found := api.FindContainerRegistryType("SOME_NON_EXISTING_INTEGRATION")
	assert.False(t, found, "container registry type should not be found")
	assert.Equal(t, 0, int(registryFound), "wrong container registry type")
	assert.Equal(t, "None", registryFound.String(), "wrong container registry type")

	registryFound, found = api.FindContainerRegistryType("GCP_GAR")
	assert.True(t, found, "container registry type should exist")
	assert.Equal(t, "GCP_GAR", registryFound.String(), "wrong container registry type")

	registryFound, found = api.FindContainerRegistryType("GHCR")
	assert.True(t, found, "container registry type should exist")
	assert.Equal(t, "GHCR", registryFound.String(), "wrong container registry type")
}

func TestContainerRegistriesGet(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		vanillaInt  = singleVanillaContainerRegistry(intgGUID, vanillaType, "")
		fakeServer  = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateContainerRegistryResponse(vanillaInt))
			}
		},
	)

	fakeServer.MockAPI("ContainerRegistries/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when container registry exists", func(t *testing.T) {
		response, err := c.V2.ContainerRegistries.Get(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.IntgGuid)
			assert.Equal(t, "integration_name", response.Data.Name)
			assert.Equal(t, "VANILLA", response.Data.Type)
		}
	})

	t.Run("when container registry does NOT exist", func(t *testing.T) {
		response, err := c.V2.ContainerRegistries.Get("UNKNOWN_INTG_GUID")
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ContainerRegistries/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestContainerRegistriesDelete(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		vanillaInt  = singleVanillaContainerRegistry(intgGUID, vanillaType, "")
		getResponse = generateContainerRegistryResponse(vanillaInt)
		fakeServer  = lacework.MockServer()
	)
	fakeServer.UseApiV2()
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

	fakeServer.MockAPI("ContainerRegistries/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify container registry exists", func(t *testing.T) {
		response, err := c.V2.ContainerRegistries.Get(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.IntgGuid)
			assert.Equal(t, "integration_name", response.Data.Name)
			assert.Equal(t, "VANILLA", response.Data.Type)
		}
	})

	t.Run("when container registry has been deleted", func(t *testing.T) {
		err := c.V2.ContainerRegistries.Delete(intgGUID)
		assert.Nil(t, err)

		response, err := c.V2.ContainerRegistries.Get(intgGUID)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ContainerRegistries/MOCK_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestContainerRegistriesList(t *testing.T) {
	var (
		gcpIntgGUIDs   = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		v2RegIntgGUIDs = []string{intgguid.New(), intgguid.New()}
		ghcrIntgGUIDs  = []string{
			intgguid.New(), intgguid.New(), intgguid.New(), intgguid.New(),
		}
		allGUIDs    = append(v2RegIntgGUIDs, append(gcpIntgGUIDs, ghcrIntgGUIDs...)...)
		expectedLen = len(allGUIDs)
		fakeServer  = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("ContainerRegistries",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			containerRegistries := []string{
				generateContainerRegistries(gcpIntgGUIDs, "GCP_GAR"),
				// TODO @afiune come back here and update these Container Registries types when they exist
				generateContainerRegistries(ghcrIntgGUIDs, "GHCR"),
				generateContainerRegistries(v2RegIntgGUIDs, "GCP_GAR"), // "V2_REGISTRY"),
			}
			fmt.Fprintf(w,
				generateContainerRegistriesResponse(
					strings.Join(containerRegistries, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ContainerRegistries.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.IntgGuid)
		switch d.ContainerRegistryType() {
		case api.NoneContainerRegistry:
			assert.True(t, false, "verify the container registry, it shouldn't be None")
		}
	}
}

func generateContainerRegistries(guids []string, iType string) string {
	containerRegistries := make([]string, len(guids))
	for i, guid := range guids {
		switch iType {
		case api.GcpGarContainerRegistry.String():
			containerRegistries[i] = singleGcpGarContainerRegistry(guid)
		case api.GhcrContainerRegistry.String():
			containerRegistries[i] = singleGhcrContainerRegistry(guid)
			// TODO @afiune come back here and update these Container Registries types
			// when they exist
			//case api.GcpGcrContainerRegistry.String():
			//containerRegistries[i] = singleGcpGcrContainerRegistry(guid)
		}
	}
	return strings.Join(containerRegistries, ", ")
}

func generateContainerRegistriesResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateContainerRegistryResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

// @afiune move to other test
func singleVanillaContainerRegistry(id string, iType string, data string) string {
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
        "details": {},
        "lastSuccessfulTime": 1624456896915,
        "lastUpdatedTime": 1624456896915,
        "ok": true
      },
      "type": "` + iType + `"
    }
	`
}
