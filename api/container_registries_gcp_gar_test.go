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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestContainerRegistriesNewGcpGar(t *testing.T) {
	subject := api.NewContainerRegistry("integration_name",
		api.GcpGarContainerRegistry,
		api.GcpGarData{
			RegistryDomain: "southamerica-east1-docker.pkg.dev",
			LimitByTag:     []string{"foo"},
			LimitByLabel:   []map[string]string{{"key": "value"}},
			LimitByRep:     []string{"xyz/name"},
			LimitNumImg:    15,
			Credentials: api.GcpCredentialsV2{
				ClientEmail:  "email",
				ClientID:     "client_id",
				PrivateKey:   "priv_key",
				PrivateKeyID: "priv_key_id",
			},
		},
	)

	assert.Equal(t, "ContVulnCfg", subject.Type)
	assert.Equal(t,
		api.GcpGarContainerRegistry.String(), subject.ContainerRegistryType().String(),
		"wrong container registry type",
	)
}

func TestContainerRegistriesGcpGarGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpGar() should be a GET method")
		fmt.Fprintf(w, generateContainerRegistryResponse(singleGcpGarContainerRegistry(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ContainerRegistries.GetGcpGar(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "integration_name", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "us-central1-docker.pkg.dev", response.Data.Data.RegistryDomain)
	assert.Contains(t, response.Data.Data.LimitByLabel, map[string]string{"key": "value"})
	assert.Equal(t, "techally-team@techally-275821.iam.gserviceaccount.com", response.Data.Data.Credentials.ClientEmail)
	assert.Equal(t, "123456789012345678901", response.Data.Data.Credentials.ClientID)
}

func TestContainerRegistriesGcpGarUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGcpGar() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "container registry name is missing")
			assert.Contains(t, body, "ContVulnCfg", "wrong container registry type")
			assert.Contains(t, body, "GCP_GAR", "wrong container registry sub-type")
			assert.Contains(t, body, "us-central1-docker.pkg.dev", "wrong registry domain")
			assert.Contains(t, body, "my_client_id", "wrong client ID")
			assert.Contains(t, body, "my_priv_key_id", "wrong private key ID")
			assert.Contains(t, body, "my_priv_key", "wrong private key")
			assert.Contains(t, body, "my_email", "wrong email")
			assert.Contains(t, body, "foo", "wrong limit by tag")
			assert.Contains(t, body, "xyz/name", "wrong limit by repo")
			assert.Contains(t, body, "15", "wrong limit num images")
			assert.Contains(t, body, "enabled\":1", "container registry is not enabled")
		}

		fmt.Fprintf(w, generateContainerRegistryResponse(singleGcpGarContainerRegistry(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	ctrReg := api.NewContainerRegistry("integration_name",
		api.GcpGarContainerRegistry,
		api.GcpGarData{
			RegistryDomain: "us-central1-docker.pkg.dev",
			LimitByTag:     []string{"foo"},
			LimitByLabel:   []map[string]string{{"key": "value"}},
			LimitByRep:     []string{"xyz/name"},
			LimitNumImg:    15,
			Credentials: api.GcpCredentialsV2{
				ClientEmail:  "my_email",
				ClientID:     "my_client_id",
				PrivateKey:   "my_priv_key",
				PrivateKeyID: "my_priv_key_id",
			},
		},
	)
	assert.Equal(t, "integration_name", ctrReg.Name, "ContVulnCfg container registry name mismatch")
	assert.Equal(t, "ContVulnCfg", ctrReg.Type, "a new ContVulnCfg container registry should match its type")
	assert.Equal(t, 1, ctrReg.Enabled, "a new ContVulnCfg container registry should be enabled")
	assert.Equal(t,
		api.GcpGarContainerRegistry.String(), ctrReg.ContainerRegistryType().String(),
		"wrong container registry type",
	)
	ctrReg.IntgGuid = intgGUID

	response, err := c.V2.ContainerRegistries.UpdateGcpGar(ctrReg)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t,
		"us-central1-docker.pkg.dev",
		response.Data.Data.RegistryDomain)
}

func singleGcpGarContainerRegistry(id string) string {
	return `
  {
    "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
    "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
    "data": {
      "credentials": {
        "clientEmail": "techally-team@techally-275821.iam.gserviceaccount.com",
        "clientId": "123456789012345678901",
        "privateKeyId": "1a2s3d4f5g6h7j8k9l01234567890abcdefghijk"
      },
      "limitByLabel": [
        {
          "key": "value"
        }
      ],
      "limitByRep": [
        "foo/bar"
      ],
      "limitByTag": [
        "owner"
      ],
      "limitNumImg": 5,
      "registryDomain": "us-central1-docker.pkg.dev",
      "registryType": "GCP_GAR"
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
