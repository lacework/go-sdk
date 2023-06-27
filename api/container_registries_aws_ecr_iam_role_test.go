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

func TestContainerRegistriesAwsEcrRoleGet(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ContainerRegistries/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsEcr() should be a GET method")
		fmt.Fprintf(w, generateContainerRegistryResponse(singleAwsEcrRoleContainerRegistry(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ContainerRegistries.GetAwsEcrIamRole(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, intgGUID, response.Data.IntgGuid)
	assert.Equal(t, "aws-ecr-iam-role", response.Data.Name)
	assert.True(t, response.Data.State.Ok)
	assert.Equal(t, "arn:aws:iam::12345678:role/test", response.Data.Data.CrossAccountCredentials.RoleArn)
	assert.Equal(t, "ABCD1234567", response.Data.Data.CrossAccountCredentials.ExternalID)
	assert.Equal(t, "AWS_IAM", response.Data.Data.AwsAuthType)
	assert.Equal(t, "AWS_ECR", response.Data.Data.RegistryType)
	assert.Equal(t, "AWS_IAM", response.Data.Data.AwsAuthType)
	assert.Equal(t, "12345678.dkr.ecr.us-east-1.amazonaws.com", response.Data.Data.RegistryDomain)
}

func singleAwsEcrRoleContainerRegistry(id string) string {
	return fmt.Sprintf(`{
        "createdOrUpdatedBy": "test@lacework.net",
        "createdOrUpdatedTime": "2021-09-28T09:16:57.092Z",
        "enabled": 1,
        "intgGuid": %q,
        "isOrg": 0,
        "name": "aws-ecr-iam-role",
        "props": {
            "tags": "AWS_ECR"
        },
        "state": {
            "ok": true,
            "lastUpdatedTime": 1663067362457,
            "lastSuccessfulTime": 1663067362457,
            "details": {
                "errorMap": {
                    "tech-ally": {
                        "errors": []
                    }
                }
            }
        },
        "type": "ContVulnCfg",
        "data": {
            "crossAccountCredentials": {
                "roleArn": "arn:aws:iam::12345678:role/test",
                "externalId": "ABCD1234567"
            },
            "awsAuthType": "AWS_IAM",
            "registryType": "AWS_ECR",
            "registryDomain": "12345678.dkr.ecr.us-east-1.amazonaws.com",
            "limitNumImg": 5,
            "limitByTag": [],
            "limitByLabel": [],
            "limitByRep": [],
            "nonOsPackageEval": true
        }
    }`, id)
}
