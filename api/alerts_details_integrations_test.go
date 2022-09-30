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
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

var alertIntegrationsJSON = `{
    "data": [
        {
            "alertChannel": {
                "INTG_GUID": "DEV7CE4D_6339B995C84CDBAEEEEEEC8EBA629CFDB5FB6",
                "NAME": "LARS",
                "CREATED_OR_UPDATED_TIME": "2022-05-02T20:44:56.612Z",
                "CREATED_OR_UPDATED_BY": "david.hazekamp@lacework.net",
                "ENV_GUID": "DEV7_A899511115596FC78BEEEEEEE8374DCD52A324B997BBCB97",
                "TYPE": "Webhook",
                "ENABLED": 1,
                "STATE": {
                    "ok": true,
                    "lastUpdatedTime": 1664561511626,
                    "lastSuccessfulTime": 1664561511626,
                    "lastWebhookUpdateTime": 0,
                    "details": {
                        "errorMessage": "<html>\r\n<head><title>503 Service Temporarily Unavailable</title></head>\r\n<body>\r\n<center><h1>503 Service Temporarily Unavailable</h1></center>\r\n</body>\r\n</html>\r\n",
                        "errorSubtitle": "Here is the response returned during the most recent channel failure:\n",
                        "errorTitle": "HTTP 503",
                        "message": "Could not send alert.  HTTP status code = 503 Channel Type: WEBHOOK integration name: LARS INTG_GUID: DEV7CE4D_6339B995C84CEEEEEEE3C8EBA629CFDB5FB6 response: <html>\r\n<head><title>503 Service Temporarily Unavailable</title></head>\r\n<body>\r\n<center><h1>503 Service Temporarily Unavailable</h1></center>\r\n</body>\r\n</html>\r\n",
                        "statusCode": "503"
                    }
                },
                "IS_ORG": 0,
                "DATA": {
                    "WEBHOOK_URL": "https://lars.run/alerts",
                    "MIN_ALERT_SEVERITY": 3
                }
            },
            "alertIntegrationId": "d7b76b0a-a9d6-e953-3asdf-53595515953f",
            "alertId": 168155,
            "integrationType": "WEBHOOK",
            "integrationContext": {
                "id": "",
                "link": ""
            },
            "intgGuid": "DEV7CE4D_6339B995C8EEEEEEEE6C3C8EBA629CFDB5FB6",
            "lastSyncTime": "2022-09-29T20:42:33.846Z",
            "alertIntegrationStatus": "",
            "status": "Open",
            "isBidirectional": false
        }
    ]
}`

func TestAlertsGetIntegrationsMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetIntegrations should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.GetIntegrations(alertID)
	assert.Nil(t, err)
}

func TestAlertsGetIntegrationsOK(t *testing.T) {
	mockResponse := alertInvestigationJSON

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, mockResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	integrationsExpected := api.AlertIntegrationsResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &integrationsExpected)

	var integrationsActual api.AlertIntegrationsResponse
	integrationsActual, err = c.V2.Alerts.GetIntegrations(alertID)
	assert.Nil(t, err)
	assert.Equal(t, integrationsExpected, integrationsActual)
}

func TestAlertsGetIntegrationsError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, lqlErrorReponse, http.StatusInternalServerError)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.GetIntegrations(alertID)
	assert.NotNil(t, err)
}
