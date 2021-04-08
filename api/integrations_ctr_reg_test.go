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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestIntegrationsNewContainerRegIntegration(t *testing.T) {
	subject := api.NewContainerRegIntegration("integration_name",
		api.ContainerRegData{
			Credentials: api.ContainerRegCreds{
				Username: "techally",
				Password: "secret",
			},
			RegistryType:   api.DockerHubRegistry.String(),
			RegistryDomain: "index.docker.io",
			LimitByTag:     "*",
			LimitByLabel:   "*",
			LimitNumImg:    5,
		},
	)
	assert.Equal(t, api.ContainerRegistryIntegration.String(), subject.Type)
}

func TestIntegrationsNewDockerHubRegistryIntegration(t *testing.T) {
	subject := api.NewDockerHubRegistryIntegration("integration_name",
		api.ContainerRegData{
			Credentials: api.ContainerRegCreds{
				Username: "techally",
				Password: "secret",
			},
		},
	)
	assert.Equal(t, api.ContainerRegistryIntegration.String(), subject.Type)
	assert.Equal(t, api.DockerHubRegistry.String(), subject.Data.RegistryType)
}

func TestIntegrationsNewDockerV2RegistryIntegration(t *testing.T) {
	subject := api.NewDockerV2RegistryIntegration("integration_name",
		api.ContainerRegData{
			Credentials: api.ContainerRegCreds{
				Username: "techally",
				Password: "secret",
				SSL:      true,
			},
		},
	)
	assert.Equal(t, api.ContainerRegistryIntegration.String(), subject.Type)
	assert.Equal(t, api.DockerV2Registry.String(), subject.Data.RegistryType)
}

func TestIntegrationsNewGcrRegistryIntegration(t *testing.T) {
	subject := api.NewGcrRegistryIntegration("integration_name",
		api.ContainerRegData{
			Credentials: api.ContainerRegCreds{
				ClientEmail:  "my@email.com",
				ClientID:     "123abc-id",
				PrivateKeyID: "aaa-key-id",
				PrivateKey:   "key",
			},
			RegistryDomain: "gcr.io",
		},
	)
	assert.Equal(t, api.ContainerRegistryIntegration.String(), subject.Type)
	assert.Equal(t, api.GcrRegistry.String(), subject.Data.RegistryType)
}
