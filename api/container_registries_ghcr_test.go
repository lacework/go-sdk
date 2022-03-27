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
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestContainerRegistriesNewGhcr(t *testing.T) {
	subject := api.NewContainerRegistry("integration_name",
		api.GhcrContainerRegistry,
		api.GhcrData{
			LimitByTag:   []string{"foo"},
			LimitByLabel: []map[string]string{{"key": "value"}},
			LimitByRep:   []string{"xyz/name"},
			LimitNumImg:  15,
			Credentials: api.GhcrCredentials{
				Username: "user",
				Password: "pass",
				Ssl:      true,
			},
		},
	)

	assert.Equal(t, "ContVulnCfg", subject.Type)
	assert.Equal(t,
		api.GhcrContainerRegistry.String(), subject.ContainerRegistryType().String(),
		"wrong container registry type",
	)
}

func TestContainerRegistriesGhcrGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGhcr() should be a GET method")
		fmt.Fprintf(w, generateContainerRegistryResponse(singleGhcrContainerRegistry(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ContainerRegistries.GetGhcr(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "ghcr.io", response.Data.Data.RegistryDomain)
	assert.Equal(t, "bubulubu", response.Data.Data.Credentials.Username)
	assert.Equal(t, "secret", response.Data.Data.Credentials.Password)
	assert.True(t, response.Data.Data.Credentials.Ssl)
}

func TestContainerRegistriesGhcrUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGhcr() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "container registry name is missing")
			assert.Contains(t, body, "ContVulnCfg", "wrong container registry type")
			assert.Contains(t, body, "GHCR", "wrong container registry sub-type")
			assert.Contains(t, body, "ghcr.io", "wrong registry domain")
			assert.Contains(t, body, "my_user", "wrong username")
			assert.Contains(t, body, "my_pass", "wrong password")
			assert.Contains(t, body, "ssl\":true", "wrong ssl config")
			assert.Contains(t, body, "enabled\":1", "container registry is not enabled")
		}

		fmt.Fprintf(w, generateContainerRegistryResponse(singleGhcrContainerRegistry(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	ctrReg := api.NewContainerRegistry("integration_name",
		api.GhcrContainerRegistry,
		api.GhcrData{
			Credentials: api.GhcrCredentials{
				Username: "my_user",
				Password: "my_pass",
				Ssl:      true,
			},
		},
	)
	assert.Equal(t, "integration_name", ctrReg.Name, "ContVulnCfg container registry name mismatch")
	assert.Equal(t, "ContVulnCfg", ctrReg.Type, "a new ContVulnCfg container registry should match its type")
	assert.Equal(t, 1, ctrReg.Enabled, "a new ContVulnCfg container registry should be enabled")
	assert.Equal(t,
		api.GhcrContainerRegistry.String(), ctrReg.ContainerRegistryType().String(),
		"wrong container registry type",
	)
	ctrReg.IntgGuid = intgGUID

	response, err := c.V2.ContainerRegistries.UpdateGhcr(ctrReg)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "ghcr.io", response.Data.Data.RegistryDomain)
}

func TestContainerRegistriesNewGhcrJson(t *testing.T) {
	subject := api.NewContainerRegistry("integration_name",
		api.GhcrContainerRegistry,
		api.GhcrData{
			LimitByTag:   []string{"foo"},
			LimitByLabel: []map[string]string{{"key": "value"}},
			LimitByRep:   []string{"xyz/name"},
			LimitNumImg:  15,
			Credentials: api.GhcrCredentials{
				Username: "user",
				Password: "pass",
				Ssl:      true,
			},
		},
	)
	jsonOut, _ := json.Marshal(subject)
	assert.Contains(t, string(jsonOut), "\"nonOsPackageEval\":false")
}

func singleGhcrContainerRegistry(id string) string {
	return `
  {
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
    "data": {
      "credentials": {
        "username": "bubulubu",
        "password": "secret",
        "ssl": true
      },
      "limitByLabel": [],
      "limitByRep": [],
      "limitByTag": [],
      "limitNumImg": 15,
      "registryDomain": "ghcr.io",
      "registryType": "GHCR"
    },
    "enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "state": {
      "details": {
          "errorMap": {}
      },
      "lastSuccessfulTime": 1624456896915,
      "lastUpdatedTime": 1624456896915,
      "ok": true
    },
    "type": "ContVulnCfg"
  }
  `
}
