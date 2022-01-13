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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestIntegrationsNewAwsEcrAccessKeyIntegration(t *testing.T) {
	subject := api.NewAwsEcrWithAccessKeyIntegration("integration_name",
		api.AwsEcrDataWithAccessKeyCreds{
			Credentials: api.AwsEcrAccessKeyCreds{
				AccessKeyID:     "id",
				SecretAccessKey: "secret",
			},
		},
	)
	assert.Equal(t, api.ContainerRegistryIntegration.String(), subject.Type)
	assert.Equal(t, api.EcrRegistry.String(), subject.Data.RegistryType)
	assert.Equal(t, api.AwsEcrAccessKey.String(), subject.Data.AwsAuthType)
}

func TestIntegrationsNewAwsEcrAccessKeyIntegrationJson(t *testing.T) {
	subject := api.NewAwsEcrWithAccessKeyIntegration("integration_name",
		api.AwsEcrDataWithAccessKeyCreds{
			Credentials: api.AwsEcrAccessKeyCreds{
				AccessKeyID:     "id",
				SecretAccessKey: "secret",
			},
		},
	)
	jsonOut, _ := json.Marshal(subject)
	assert.Contains(t, string(jsonOut), "\"NON_OS_PACKAGE_EVAL\":false")
}
