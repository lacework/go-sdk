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
	assert.Equal(t,
		"Datadog", api.DatadogAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"Webhook", api.WebhookAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"CiscoSparkWebhook", api.CiscoSparkWebhookAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"MicrosoftTeams", api.MicrosoftTeamsAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"GcpPubsub", api.GcpPubSubAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"SplunkHec", api.SplunkHecAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"ServiceNowRest", api.ServiceNowRestAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"NewRelicInsights", api.NewRelicInsightsAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"PagerDutyApi", api.PagerDutyApiAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"IbmQradar", api.IbmQRadarAlertChannelType.String(),
		"wrong alert channel type",
	)
	assert.Equal(t,
		"Jira", api.JiraAlertChannelType.String(),
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

	alertFound, found = api.FindAlertChannelType("Datadog")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "Datadog", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("Webhook")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "Webhook", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("VictorOps")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "VictorOps", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("CiscoSparkWebhook")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "CiscoSparkWebhook", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("MicrosoftTeams")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "MicrosoftTeams", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("GcpPubsub")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "GcpPubsub", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("SplunkHec")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "SplunkHec", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("ServiceNowRest")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "ServiceNowRest", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("NewRelicInsights")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "NewRelicInsights", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("PagerDutyApi")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "PagerDutyApi", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("IbmQradar")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "IbmQradar", alertFound.String(), "wrong alert channel type")

	alertFound, found = api.FindAlertChannelType("Jira")
	assert.True(t, found, "alert channel type should exist")
	assert.Equal(t, "Jira", alertFound.String(), "wrong alert channel type")
}

func TestAlertChannelsGet(t *testing.T) {
	var (
		intgGUID    = intgguid.New()
		vanillaType = "VANILLA"
		apiPath     = fmt.Sprintf("AlertChannels/%s", intgGUID)
		vanillaInt  = singleVanillaAlertChannel(intgGUID, vanillaType, "")
		fakeServer  = lacework.MockServer()
	)
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
		allGUIDs                   []string
		emailAlertChan             = generateGuids(&allGUIDs, 3)
		awsS3AlertChan             = generateGuids(&allGUIDs, 2)
		slackAlertChan             = generateGuids(&allGUIDs, 4)
		cloudwatchAlertChan        = generateGuids(&allGUIDs, 2)
		datadogAlertChan           = generateGuids(&allGUIDs, 2)
		webhookAlertChan           = generateGuids(&allGUIDs, 2)
		victorOpsAlertChan         = generateGuids(&allGUIDs, 2)
		ciscoSparkWebhookAlertChan = generateGuids(&allGUIDs, 2)
		microsoftTeamsAlertChan    = generateGuids(&allGUIDs, 2)
		gcpPubSubAlertChan         = generateGuids(&allGUIDs, 2)
		splunkHecAlertChan         = generateGuids(&allGUIDs, 2)
		serviceNowRestAlertChan    = generateGuids(&allGUIDs, 2)
		newRelicInsightsAlertChan  = generateGuids(&allGUIDs, 2)
		pagerDutyApiAlertChan      = generateGuids(&allGUIDs, 2)
		ibmQradarAlertChan         = generateGuids(&allGUIDs, 2)
		jiraAlertChan              = generateGuids(&allGUIDs, 2)
		expectedLen                = len(allGUIDs)
		fakeServer                 = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("AlertChannels",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			alertChannels := []string{
				generateAlertChannels(emailAlertChan, "EmailUser"),
				generateAlertChannels(slackAlertChan, "SlackChannel"),
				generateAlertChannels(awsS3AlertChan, "AwsS3"),
				generateAlertChannels(cloudwatchAlertChan, "CloudwatchEb"),
				generateAlertChannels(datadogAlertChan, "Datadog"),
				generateAlertChannels(webhookAlertChan, "Webhook"),
				generateAlertChannels(victorOpsAlertChan, "VictorOps"),
				generateAlertChannels(ciscoSparkWebhookAlertChan, "CiscoSparkWebhook"),
				generateAlertChannels(microsoftTeamsAlertChan, "MicrosoftTeams"),
				generateAlertChannels(gcpPubSubAlertChan, "GcpPubsub"),
				generateAlertChannels(splunkHecAlertChan, "SplunkHec"),
				generateAlertChannels(serviceNowRestAlertChan, "ServiceNowRest"),
				generateAlertChannels(pagerDutyApiAlertChan, "PagerDutyApi"),
				generateAlertChannels(newRelicInsightsAlertChan, "NewRelicInsights"),
				generateAlertChannels(ibmQradarAlertChan, "IbmQradar"),
				generateAlertChannels(jiraAlertChan, "Jira"),
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

func generateGuids(allGuids *[]string, guidCount int) []string {
	var channelGuids []string

	for i := 0; i < guidCount; i++ {
		channelGuids = append(channelGuids, intgguid.New())
	}

	*allGuids = append(*allGuids, channelGuids...)
	return channelGuids
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
		case api.DatadogAlertChannelType.String():
			alertChannels[i] = singleDatadogAlertChannel(guid)
		case api.WebhookAlertChannelType.String():
			alertChannels[i] = singleWebhookAlertChannel(guid)
		case api.VictorOpsAlertChannelType.String():
			alertChannels[i] = singleVictorOpsAlertChannel(guid)
		case api.CiscoSparkWebhookAlertChannelType.String():
			alertChannels[i] = singleCiscoSparkWebhookAlertChannel(guid)
		case api.MicrosoftTeamsAlertChannelType.String():
			alertChannels[i] = singleMicrosoftTeamsAlertChannel(guid)
		case api.GcpPubSubAlertChannelType.String():
			alertChannels[i] = singleGcpPubSubAlertChannel(guid)
		case api.SplunkHecAlertChannelType.String():
			alertChannels[i] = singleSplunkAlertChannel(guid)
		case api.ServiceNowRestAlertChannelType.String():
			alertChannels[i] = singleServiceNowRestAlertChannel(guid)
		case api.NewRelicInsightsAlertChannelType.String():
			alertChannels[i] = singleNewRelicAlertChannel(guid)
		case api.PagerDutyApiAlertChannelType.String():
			alertChannels[i] = singlePagerDutyAlertChannel(guid)
		case api.IbmQRadarAlertChannelType.String():
			alertChannels[i] = singleIbmQRadarAlertChannel(guid)
		case api.JiraAlertChannelType.String():
			alertChannels[i] = singleJiraCloudAlertChannel(guid)
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
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "POST", r.Method, "Test() should be a POST method") {

			}
		},
	)

	c, err := api.NewClient("test",
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
