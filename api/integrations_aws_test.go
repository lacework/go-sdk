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

func TestIntegrationsNewAwsCfgIntegration(t *testing.T) {
	subject := api.NewAwsCfgIntegration("integration_name",
		api.AwsIntegrationData{
			Credentials: &api.AwsCrossAccountCreds{
				RoleArn:    "arn:foo:bar",
				ExternalID: "0123456789",
			},
		},
	)
	assert.Equal(t, api.AwsCfgIntegration.String(), subject.Type)
}

func TestIntegrationsEmptyAwsCredentials(t *testing.T) {
	var awsData  api.AwsIntegrationData
	credentials := awsData.GetCredentials()

	externalID := credentials.ExternalID
	roleArn := credentials.RoleArn
	assert.Empty(t, externalID)
	assert.Empty(t, roleArn)
}

func TestIntegrationsEmptyAwsGovCloudCredentials(t *testing.T) {
	var awsData  api.AwsIntegrationData
	credentials := awsData.GetGovCloudCredentials()
	accountID := awsData.GetAccountID()

	secretKey := credentials.SecretAccessKey
	accessID := credentials.AccessKeyID
	assert.Empty(t, accountID)
	assert.Empty(t, secretKey)
	assert.Empty(t, accessID)
}

func TestIntegrationsNewAwsCfgIntegrationWithCustomTemplateFile(t *testing.T) {
	accountMappingJSON := []byte(`{
    "defaultLaceworkAccountAws": "lw_account_1",
    "integration_mappings": {
      "lw_account_2": {
        "aws_accounts": [
          "234556677",
          "774564564"
        ]
      },
      "lw_account_3": {
        "aws_accounts": [
          "553453453",
          "934534535"
        ]
      }
    }
  }`)
	awsData := api.AwsIntegrationData{
		Credentials: &api.AwsCrossAccountCreds{
			RoleArn:    "arn:foo:bar",
			ExternalID: "0123456789",
		},
	}
	awsData.EncodeAccountMappingFile(accountMappingJSON)

	subject := api.NewAwsCfgIntegration("integration_name", awsData)
	assert.Equal(t, api.AwsCfgIntegration.String(), subject.Type)
	assert.Contains(t,
		subject.Data.AccountMappingFile,
		"data:application/json;name=i.json;base64,",
		"check the custom_template_file encoder",
	)
	accountMapping, err := subject.Data.DecodeAccountMappingFile()
	assert.Nil(t, err)
	assert.Equal(t, accountMappingJSON, accountMapping)

	// When there is no custom account mapping file, this function should
	// return an empty string to match the pattern
	subject.Data.AccountMappingFile = ""
	accountMapping, err = subject.Data.DecodeAccountMappingFile()
	assert.Nil(t, err)
	assert.Empty(t, accountMapping)

}

func TestIntegrationsCreateAws(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)

	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateAws should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "AWS_CFG", "wrong integration type")
			assert.Contains(t, body, "arn:foo:bar", "wrong role arn")
			assert.Contains(t, body, "0123456789", "wrong external id")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, awsIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewAwsIntegration("integration_name",
		api.AwsCfgIntegration,
		api.AwsIntegrationData{
			Credentials: &api.AwsCrossAccountCreds{
				RoleArn:    "arn:foo:bar",
				ExternalID: "0123456789",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "AWS integration name mismatch")
	assert.Equal(t, "AWS_CFG", data.Type, "a new AWS integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new AWS integration should be enabled")

	response, err := c.Integrations.CreateAws(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "arn:foo:bar", resData.Data.Credentials.RoleArn)
		assert.Equal(t, "0123456789", resData.Data.Credentials.ExternalID)
	}
}

func TestIntegrationsGetAws(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAws should be a GET method")
		fmt.Fprintf(w, awsIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetAws(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "arn:foo:bar", resData.Data.Credentials.RoleArn)
		assert.Equal(t, "0123456789", resData.Data.Credentials.ExternalID)
	}
}

func TestIntegrationsUpdateAws(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAws should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "AWS_CT_SQS", "wrong integration type")
			assert.Contains(t, body, "arn:foo:bar", "wrong role arn")
			assert.Contains(t, body, "0123456789", "wrong external ID")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, awsIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewAwsIntegration("integration_name",
		api.AwsCloudTrailIntegration,
		api.AwsIntegrationData{
			Credentials: &api.AwsCrossAccountCreds{
				RoleArn:    "arn:foo:bar",
				ExternalID: "0123456789",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "AWS integration name mismatch")
	assert.Equal(t, "AWS_CT_SQS", data.Type, "a new AWS integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new AWS integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateAws(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsDeleteAws(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "DeleteAws should be a DELETE method")
		fmt.Fprintf(w, awsIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.DeleteAws(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListAwsCfg(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/AWS_CFG",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListAwsCfg should be a GET method")
			fmt.Fprintf(w, awsMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListAwsCfg()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func awsIntegrationJsonResponse(intgGUID string) string {
	return `
		{
			"data": [` + singleAwsIntegration(intgGUID) + `],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func awsMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleAwsIntegration(guid))
	}
	return `
		{
			"data": [` + strings.Join(integrations, ", ") + `],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func singleAwsIntegration(id string) string {
	return `
		{
			"INTG_GUID": "` + id + `",
			"NAME": "integration_name",
			"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
			"CREATED_OR_UPDATED_BY": "user@email.com",
			"TYPE": "AWS_CFG",
			"ENABLED": 1,
			"STATE": {
				"ok": true,
				"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
				"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
			},
			"IS_ORG": 0,
			"DATA": {
				"CROSS_ACCOUNT_CREDENTIALS": {
          "ROLE_ARN": "arn:foo:bar",
					"EXTERNAL_ID": "0123456789"
				}
			},
			"TYPE_NAME": "AWS Compliance"
		}
	`
}
