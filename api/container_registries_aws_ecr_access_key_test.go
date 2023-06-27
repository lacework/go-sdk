//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestContainerRegistriesAwsEcrAccessKeyGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsEcr() should be a GET method")
		fmt.Fprintf(w, generateContainerRegistryResponse(singleAwsEcrAccessKeyContainerRegistry(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ContainerRegistries.GetAwsEcrAccessKey(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "aws-ecr-access-key", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "", response.Data.Data.AccessKeyCredentials.SecretAccessKey)
	assert.Equal(t, "ABCDEFGH", response.Data.Data.AccessKeyCredentials.AccessKeyID)
	assert.Equal(t, "AWS_ACCESS_KEY", response.Data.Data.AwsAuthType)
	assert.Equal(t, "AWS_ECR", response.Data.Data.RegistryType)
	assert.Equal(t, "AWS_ACCESS_KEY", response.Data.Data.AwsAuthType)
	assert.Equal(t, "123456.test.ecr.us-west-1.amazonaws.com", response.Data.Data.RegistryDomain)
}

func singleAwsEcrAccessKeyContainerRegistry(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "test@lacework.net",
        "createdOrUpdatedTime": "2022-07-13T23:59:47.943Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 0,
        "name": "aws-ecr-access-key",
        "props": {
            "tags": "AWS_ECR"
        },
        "state": {
            "ok": true,
            "lastUpdatedTime": 1663066471375,
            "lastSuccessfulTime": 1663066471375,
            "details": {
                "error": {
                    "customerThrottling": "Max-hourly image scan capacity reached"
                },
                "errorMap": {
                    "test-repo": {
                        "errors": []
                    }
                }
            }
        },
        "type": "ContVulnCfg",
        "data": {
            "accessKeyCredentials": {
                "accessKeyId": "ABCDEFGH",
                "secretAccessKey": ""
            },
            "awsAuthType": "AWS_ACCESS_KEY",
            "registryType": "AWS_ECR",
            "registryDomain": "123456.test.ecr.us-west-1.amazonaws.com",
            "limitByTag": [],
            "limitByLabel": [],
            "limitByRep": [],
            "limitNumImg": 5,
            "nonOsPackageEval": true
        }
}`, id)
}
