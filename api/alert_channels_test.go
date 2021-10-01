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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestAlertChannelTypes(t *testing.T) {
	assert.Equal(t,
		"EmailUser", api.EmailUserAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"SlackChannel", api.SlackChannelAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"AwsS3", api.AwsS3AlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"CloudwatchEb", api.CloudwatchEbAlertChannelType.String(),
		"wrong alert channel type",
	)
}

func TestFindAlertChannelType(t *testing.T) {
	alertFound, found := api.FindAlertChannelType("SOME_NON_EXISTING_INTEGRATION")
	assert.False(t, found, "alert channel type should not be found")
	assert.Equal(t, 0, int(alertFound), "wrong alert channel type")
	assert.Equal(t, "None", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("EmailUser")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "EmailUser", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("SlackChannel")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "SlackChannel", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("AwsS3")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "AwsS3", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("CloudwatchEb")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "CloudwatchEb", alertFound.String(), "wrong alert channel type")
}

func TestAlertChannelsGet(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("AlertChannels/%s", intgGUID)
		vanillaInt  = singleVanillaAlertChannel(intgGUID, vanillaType, "")
		fakeServer  = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateAlertChannelResponse(vanillaInt))
			}
		},
	)

	fakeServer.MockAPI("AlertChannels/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when alert channel exists", func(t *testing.T) {
		var response api.AlertChannelResponse
		err := c.V2.AlertChannels.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.IntgGuid)
			assert.Equal(t, "integration_name", response.Data.Name)
			assert.Equal(t, "VANILLA", response.Data.Type)
		}
	})

	t.Run("when alert channel does NOT exist", func(t *testing.T) {
		var response api.AlertChannelResponse
		err := c.V2.AlertChannels.Get("UNKNOWN_INTG_GUID", response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/AlertChannels/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestAlertChannelsDelete(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("AlertChannels/%s", intgGUID)
		vanillaInt  = singleVanillaAlertChannel(intgGUID, vanillaType, "")
		getResponse = generateAlertChannelResponse(vanillaInt)
		fakeServer  = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if getResponse != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, getResponse)
				case "DELETE":
					// once deleted, empty the getResponse so that
					// further GET requests return 404s
					getResponse = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	fakeServer.MockAPI("AlertChannels/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify alert channel exists", func(t *testing.T) {
		var response api.AlertChannelResponse
		err := c.V2.AlertChannels.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.IntgGuid)
			assert.Equal(t, "integration_name", response.Data.Name)
			assert.Equal(t, "VANILLA", response.Data.Type)
		}
	})

	t.Run("when alert channel has been deleted", func(t *testing.T) {
		err := c.V2.AlertChannels.Delete(intgGUID)
		assert.Nil(t, err)

		var response api.AlertChannelResponse
		err = c.V2.AlertChannels.Get(intgGUID, &response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/AlertChannels/MOCK_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestAlertChannelsList(t *testing.T) {
	var (
		emailAlertChan = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		awsS3AlertChan = []string{intgguid.New(), intgguid.New()}
		slackAlertChan = []string{
			intgguid.New(), intgguid.New(), intgguid.New(), intgguid.New(),
		}
		cloudwatchAlertChan = []string{intgguid.New(), intgguid.New()}
		someGUIDs           = append(awsS3AlertChan, append(slackAlertChan, emailAlertChan...)...)
		allGUIDs            = append(someGUIDs, cloudwatchAlertChan...)
		expectedLen         = len(allGUIDs)
		fakeServer          = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("AlertChannels",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			alertChannels := []string{
				generateAlertChannels(emailAlertChan, "EmailUser"),
				generateAlertChannels(slackAlertChan, "SlackChannel"),
				generateAlertChannels(awsS3AlertChan, "AwsS3"),
				generateAlertChannels(cloudwatchAlertChan, "CloudwatchEb"),
			}
			fmt.Fprintf(w,
				generateAlertChannelsResponse(
					strings.Join(alertChannels, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.AlertChannels.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.IntgGuid)
	}
}

func generateAlertChannels(guids []string, iType string) string {
	alertChannels := make([]string, len(guids))
	for i, guid := range guids {
		switch iType {
		case api.EmailUserAlertChannelType.String():
			alertChannels[i] = singleEmailUserAlertChannel(guid)
		case api.SlackChannelAlertChannelType.String():
			alertChannels[i] = singleSlackChannelAlertChannel(guid)
		case api.AwsS3AlertChannelType.String():
			alertChannels[i] = singleAwsS3AlertChannel(guid)
		case api.CloudwatchEbAlertChannelType.String():
			alertChannels[i] = singleAWSCloudwatchAlertChannel(guid)
		}
	}
	return strings.Join(alertChannels, ", ")
}

func TestAlertChannelsTest(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("AlertChannels/%s/test", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "POST", r.Method, "Test() should be a POST method") {

			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when alert channel exists", func(t *testing.T) {
		err := c.V2.AlertChannels.Test(intgGUID)
		assert.Nil(t, err)
	})

	t.Run("when alert channel id is empty", func(t *testing.T) {
		err := c.V2.AlertChannels.Test("")
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "specify an intgGuid")
		}
	})
}

func generateAlertChannelsResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateAlertChannelResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleVanillaAlertChannel(id string, iType string, data string) string {
	if data == "" {
		data = "{}"
	}
	return `
    {
      "createdOrUpdatedBy": "salim.afiunemaya@lacework.net",
      "createdOrUpdatedTime": "2021-06-01T19:28:00.092Z",
      "data": ` + data + `,
      "enabled": 1,
      "intgGuid": "` + id + `",
      "isOrg": 0,
      "name": "integration_name",
      "state": {
        "details": {},
        "lastSuccessfulTime": 1624456896915,
        "lastUpdatedTime": 1624456896915,
        "ok": true
      },
      "type": "` + iType + `"
    }
	`
}
