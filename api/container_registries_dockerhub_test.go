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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
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
