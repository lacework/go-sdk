//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

func TestIntegrationsNewAwsS3AlertChannel(t *testing.T) {
	subject := api.NewAwsS3AlertChannel("integration_name",
		api.AwsS3ChannelData{
			Credentials: api.AwsS3Creds{
				RoleArn:    "arn:aws:iam::1234567890:role/lacework_iam_example_role",
				BucketArn:  "arn:aws:s3:::bucket_name/key_name",
				ExternalID: "0123456789",
			},
		},
	)
	assert.Equal(t, api.AwsS3ChannelIntegration.String(), subject.Type)
}

func TestIntegrationsCreateAwsS3AlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateAwsS3AlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "AWS_S3", "wrong integration type")
			assert.Contains(t, body, "arn:aws:iam::1234567890:role/lacework_iam_example_role", "wrong role arn")
			assert.Contains(t, body, "arn:aws:s3:::bucket_name/key_name", "wrong bucket arn")
			assert.Contains(t, body, "0123456789", "wrong external id")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, awsS3ChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewAwsS3AlertChannel("integration_name",
		api.AwsS3ChannelData{
			Credentials: api.AwsS3Creds{
				RoleArn:    "arn:aws:iam::1234567890:role/lacework_iam_example_role",
				BucketArn:  "arn:aws:s3:::bucket_name/key_name",
				ExternalID: "0123456789",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "AwsS3Channel integration name mismatch")
	assert.Equal(t, "AWS_S3", data.Type, "a new AwsS3Channel integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new AwsS3Channel integration should be enabled")

	response, err := c.Integrations.CreateAwsS3AlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "arn:aws:iam::1234567890:role/lacework_iam_example_role", resData.Data.Credentials.RoleArn)
		assert.Equal(t, "arn:aws:s3:::bucket_name/key_name", resData.Data.Credentials.BucketArn)
		assert.Equal(t, "0123456789", resData.Data.Credentials.ExternalID)
	}
}

func TestIntegrationsGetAwsS3AlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsS3AlertChannel should be a GET method")
		fmt.Fprintf(w, awsS3ChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetAwsS3AlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "arn:aws:iam::1234567890:role/lacework_iam_example_role", resData.Data.Credentials.RoleArn)
		assert.Equal(t, "arn:aws:s3:::bucket_name/key_name", resData.Data.Credentials.BucketArn)
		assert.Equal(t, "0123456789", resData.Data.Credentials.ExternalID)
	}
}

func TestIntegrationsUpdateAwsS3AlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsS3AlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "SLACK_CHANNEL", "wrong integration type")
			assert.Contains(t, body, "arn:aws:iam::1234567890:role/lacework_iam_example_role", "wrong role arn")
			assert.Contains(t, body, "arn:aws:s3:::bucket_name/key_name", "wrong bucket arn")
			assert.Contains(t, body, "0123456789", "wrong external id")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, awsS3ChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewAwsS3AlertChannel("integration_name",
		api.AwsS3ChannelData{
			Credentials: api.AwsS3Creds{
				RoleArn:    "arn:aws:iam::1234567890:role/lacework_iam_example_role",
				BucketArn:  "arn:aws:s3:::bucket_name/key_name",
				ExternalID: "0123456789",
			},
		},
	)
	assert.Equal(t, "integration_name", data.Name, "AwsS3Channel integration name mismatch")
	assert.Equal(t, "AWS_S3", data.Type, "a new AwsS3Channel integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new AwsS3Channel integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateAwsS3AlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListAwsS3AlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/AWS_S3",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListAwsS3AlertChannel should be a GET method")
			fmt.Fprintf(w, awsS3ChanMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListAwsS3AlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func awsS3ChannelIntegrationJsonResponse(intgGUID string) string {
	return `
{
  "data": [` + singleAwsS3ChanIntegration(intgGUID) + `],
  "ok": true,
  "message": "SUCCESS"
}
`
}

func awsS3ChanMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleAwsS3ChanIntegration(guid))
	}
	return `
{
"data": [` + strings.Join(integrations, ", ") + `],
"ok": true,
"message": "SUCCESS"
}
`
}

func singleAwsS3ChanIntegration(id string) string {
	return `
	{
		"INTG_GUID": "` + id + `",
		"CREATED_OR_UPDATED_BY": "user@email.com",
		"CREATED_OR_UPDATED_TIME": "2020-Jul-16 19:59:22 UTC",
		"DATA": {
		  "ISSUE_GROUPING": "Events",
		  "S3_CROSS_ACCOUNT_CREDENTIALS": {
		     "ROLE_ARN": "arn:aws:iam::1234567890:role/lacework_iam_example_role",
		     "BUCKET_ARN": "arn:aws:s3:::bucket_name/key_name",
		     "EXTERNAL_ID": "0123456789"
		  }
		},
		"ENABLED": 1,
		"IS_ORG": 0,
		"NAME": "integration_name",
		"STATE": {
		  "lastSuccessfulTime": "2020-Jul-16 18:26:54 UTC",
		  "lastUpdatedTime": "2020-Jul-16 18:26:54 UTC",
		  "ok": true
		},
		"TYPE": "AWS_S3",
		"TYPE_NAME": "AWS_S3"
	  }
`
}
