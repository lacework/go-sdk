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

func TestContainerRegistriesNewDockerhub(t *testing.T) {
	subject := api.NewContainerRegistry("integration_name",
		api.DockerhubContainerRegistry,
		api.DockerhubData{
			LimitByTag: []string{"foo"},
			LimitByLabel: []map[string]string{
				{"key1": "value1"},
				{"key2": "value2"},
			},
			LimitByRep:  []string{"xyz/name"},
			LimitNumImg: 15,
			Credentials: api.DockerhubCredentials{
				Username: "username",
				Password: "password",
			},
		},
	)

	assert.Equal(t, "ContVulnCfg", subject.Type)
	assert.Equal(t,
		api.DockerhubContainerRegistry.String(), subject.ContainerRegistryType().String(),
		"wrong container registry type",
	)

	hub, ok := subject.Data.(api.DockerhubData)
	if assert.True(t, ok) {
		assert.Equal(t, "index.docker.io",
			hub.RegistryDomain, "wrong container registry domain")
	}
}

func TestContainerRegistriesDockerhubGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetDockerhub() should be a GET method")
		fmt.Fprintf(w, generateContainerRegistryResponse(singleDockerhubContainerRegistry(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ContainerRegistries.GetDockerhub(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "index.docker.io", response.Data.Data.RegistryDomain)
	assert.Equal(t, "bubulubu", response.Data.Data.Credentials.Username)
	assert.Equal(t, "", response.Data.Data.Credentials.Password)
}

func TestContainerRegistriesDockerhubUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateDockerhub() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "container registry name is missing")
			assert.Contains(t, body, "ContVulnCfg", "wrong container registry type")
			assert.Contains(t, body, "DOCKERHUB", "wrong container registry sub-type")
			assert.Contains(t, body, "index.docker.io", "wrong registry domain")
			assert.Contains(t, body, "my_user", "wrong username")
			assert.Contains(t, body, "my_pass", "wrong password")
			assert.Contains(t, body, "enabled\":1", "container registry is not enabled")
		}

		fmt.Fprintf(w, generateContainerRegistryResponse(singleDockerhubContainerRegistry(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	ctrReg := api.NewContainerRegistry("integration_name",
		api.DockerhubContainerRegistry,
		api.DockerhubData{
			Credentials: api.DockerhubCredentials{
				Username: "my_user",
				Password: "my_pass",
			},
		},
	)
	assert.Equal(t, "integration_name", ctrReg.Name,
		"ContVulnCfg container registry name mismatch")
	assert.Equal(t, "ContVulnCfg", ctrReg.Type,
		"a new ContVulnCfg container registry should match its type")
	assert.Equal(t, 1, ctrReg.Enabled,
		"a new ContVulnCfg container registry should be enabled")
	assert.Equal(t,
		api.DockerhubContainerRegistry.String(), ctrReg.ContainerRegistryType().String(),
		"wrong container registry type",
	)

	ctrReg.IntgGuid = intgGUID

	response, err := c.V2.ContainerRegistries.UpdateDockerhub(ctrReg)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "index.docker.io", response.Data.Data.RegistryDomain)
}

func TestContainerRegistriesNewDockerhubJson(t *testing.T) {
	subject := api.NewContainerRegistry("integration_name",
		api.DockerhubContainerRegistry,
		api.DockerhubData{
			LimitByTag: []string{"foo"},
			LimitByLabel: []map[string]string{
				{"key1": "value1"},
				{"key2": "value2"},
			},
			LimitByRep:  []string{"xyz/name"},
			LimitNumImg: 15,
			Credentials: api.DockerhubCredentials{
				Username: "user",
				Password: "pass",
			},
		},
	)
	jsonOut, _ := json.Marshal(subject)
	assert.Contains(t, string(jsonOut), "\"nonOsPackageEval\":false")
}

func singleDockerhubContainerRegistry(id string) string {
	return `
  {
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-12-18T02:55:39.767Z",
    "data": {
      "credentials": {
        "password": "",
        "username": "bubulubu"
      },
      "limitByLabel": [],
      "limitByRep": [],
      "limitByTag": [],
      "limitNumImg": 5,
      "nonOsPackageEval": true,
      "registryDomain": "index.docker.io",
      "registryType": "DOCKERHUB"
    },
    "enabled": 1,
    "intgGuid": "` + id + `",
    "isOrg": 0,
    "name": "integration_name",
    "props": {
      "policyEvaluation": {
        "evaluate": false
      },
      "tags": "DOCKERHUB"
    },
    "state": {
      "details": {
        "errorMap": {}
      },
      "lastSuccessfulTime": 1669164302390,
      "lastUpdatedTime": 1669164302390,
      "ok": true
    },
    "type": "ContVulnCfg"
  }
  `
}
