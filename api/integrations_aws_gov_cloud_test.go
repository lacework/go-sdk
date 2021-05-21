//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

func TestIntegrationsNewAwsGovCloudCfgIntegration(t *testing.T) {
	subject := api.NewAwsGovCloudCfgIntegration("integration_name",
		api.AwsGovCloudCfgIntegration,
		api.AwsGovCloudIntegrationData{
			Credentials: api.AwsGovCloudCreds{
				AccessKeyID:     "AWS123abcAccessKeyID",
				SecretAccessKey: "AWS123abc123abcSecretAccessKey0000000000",
				AccountID:       "0123456789",
			},
		},
	)
	assert.Equal(t, api.AwsGovCloudCfgIntegration.String(), subject.Type)
}

func TestIntegrationsCreateAwsGovCloudCfg(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)

	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateAws should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "AWS_US_GOV_CFG", "wrong integration type")
			assert.Contains(t, body, "AWS123abcAccessKeyID", "wrong access key id")
			assert.Contains(t, body, "AWS123abc123abcSecretAccessKey0000000000", "wrong secret key")
			assert.Contains(t, body, "0123456789", "wrong account id")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, awsGovCloudIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewAwsGovCloudCfgIntegration("integration_name",
		api.AwsGovCloudCfgIntegration,
		api.AwsGovCloudIntegrationData{
			Credentials: api.AwsGovCloudCreds{
				AccessKeyID:     "AWS123abcAccessKeyID",
				SecretAccessKey: "AWS123abc123abcSecretAccessKey0000000000",
				AccountID:       "0123456789",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "AWS Gov Cloud integration name mismatch")
	assert.Equal(t, "AWS_US_GOV_CFG", data.Type, "a new AWS Gov Cloud integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new AWS Gov Cloud integration should be enabled")

	response, err := c.Integrations.CreateAwsGovCloudCfg(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "0123456789", resData.Data.Credentials.AccountID)
		assert.Equal(t, "AWS123abcAccessKeyID", resData.Data.Credentials.AccessKeyID)
		assert.Equal(t, "AWS123abc123abcSecretAccessKey0000000000", resData.Data.Credentials.SecretAccessKey)
	}
}

func TestIntegrationsGetAwsGovCloudCfg(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsGovCloud should be a GET method")
		fmt.Fprintf(w, awsGovCloudIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetAwsGovCloudCfg(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "0123456789", resData.Data.Credentials.AccountID)
		assert.Equal(t, "AWS123abcAccessKeyID", resData.Data.Credentials.AccessKeyID)
		assert.Equal(t, "AWS123abc123abcSecretAccessKey0000000000", resData.Data.Credentials.SecretAccessKey)
	}
}

func TestIntegrationsUpdateAwsGovCloudCfg(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsGovCloud should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "AWS_US_GOV_CFG", "wrong integration type")
			assert.Contains(t, body, "AWS123abcAccessKeyID", "wrong access key id")
			assert.Contains(t, body, "AWS123abc123abcSecretAccessKey0000000000", "wrong secret access key")
			assert.Contains(t, body, "0123456789", "wrong account ID")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, awsGovCloudIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewAwsGovCloudCfgIntegration("integration_name",
		api.AwsGovCloudCfgIntegration,
		api.AwsGovCloudIntegrationData{
			Credentials: api.AwsGovCloudCreds{
				AccessKeyID:     "AWS123abcAccessKeyID",
				SecretAccessKey: "AWS123abc123abcSecretAccessKey0000000000",
				AccountID:       "0123456789",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "AWS Gov Cloud integration name mismatch")
	assert.Equal(t, "AWS_US_GOV_CFG", data.Type, "a new AWS Gov Cloud integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new AWS Gov Cloud integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateAwsGovCloudCfg(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsDeleteAwsGovCloudCfg(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method, "DeleteAws should be a DELETE method")
		fmt.Fprintf(w, awsGovCloudIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.DeleteAwsGovCloudCfg(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListAwsGovCloudCfg(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/AWS_US_GOV_CFG",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListAwsGovCloudCfg should be a GET method")
			fmt.Fprintf(w, awsGovCloudMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListAwsGovCloudCfg()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func awsGovCloudIntegrationJsonResponse(intgGUID string) string {
	return `
		{
			"data": [` + singleAwsGovCloudIntegration(intgGUID) + `],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func awsGovCloudMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleAwsGovCloudIntegration(guid))
	}
	return `
		{
			"data": [` + strings.Join(integrations, ", ") + `],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func singleAwsGovCloudIntegration(id string) string {
	return `
		{
			"INTG_GUID": "` + id + `",
			"NAME": "integration_name",
			"CREATED_OR_UPDATED_TIME": "2020-Mar-10 01:00:00 UTC",
			"CREATED_OR_UPDATED_BY": "user@email.com",
			"TYPE": "AWS_US_GOV_CFG",
			"ENABLED": 1,
			"STATE": {
				"ok": true,
				"lastUpdatedTime": "2020-Mar-10 01:00:00 UTC",
				"lastSuccessfulTime": "2020-Mar-10 01:00:00 UTC"
			},
			"IS_ORG": 0,
			"DATA": {
				"ACCESS_KEY_CREDENTIALS": {
          			"ACCOUNT_ID": "0123456789",
					"ACCESS_KEY_ID": "AWS123abcAccessKeyID",
					"SECRET_ACCESS_KEY": "AWS123abc123abcSecretAccessKey0000000000"
				}
			},
			"TYPE_NAME": "AWS Compliance"
		}
	`
}
