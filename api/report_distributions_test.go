//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

func TestReportDistributionsGet(t *testing.T) {
	var (
		intgGUID           = intgguid.New()
		apiPath            = fmt.Sprintf("ReportDistributions/%s", intgGUID)
		reportDistribution = singleMockReportDistribution(intgGUID)
		fakeServer         = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateReportDistributionResponse(reportDistribution))
			}
		},
	)

	fakeServer.MockAPI("ReportDistributions/UNKNOWN_INTG_GUID",
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

	t.Run("when report distribution exists", func(t *testing.T) {
		response, err := c.V2.ReportDistributions.Get(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAB", response.Data.ReportDefinitionGuid)
			assert.Equal(t, intgGUID, response.Data.ReportDistributionGuid)
			assert.Equal(t, "distribution_name", response.Data.DistributionName)
			assert.Equal(t, []string{"critical"}, response.Data.Data.Severities)
			assert.Equal(t, []string{"NonCompliant"}, response.Data.Data.Violations)
		}
	})

	t.Run("when report distribution does NOT exist", func(t *testing.T) {
		var response api.ReportDistributionResponse
		response, err := c.V2.ReportDistributions.Get("UNKNOWN_INTG_GUID")
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ReportDistributions/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestReportDistributionDelete(t *testing.T) {
	var (
		intgGUID           = intgguid.New()
		apiPath            = fmt.Sprintf("ReportDistributions/%s", intgGUID)
		reportDistribution = singleMockReportDistribution(intgGUID)
		getResponse        = generateReportDistributionResponse(reportDistribution)
		fakeServer         = lacework.MockServer()
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

	fakeServer.MockAPI("ReportDistributions/UNKNOWN_INTG_GUID",
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

	t.Run("verify report distribution exists", func(t *testing.T) {
		var response api.ReportDistributionResponse
		response, err := c.V2.ReportDistributions.Get(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.ReportDistributionGuid)
			assert.Equal(t, "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAB", response.Data.ReportDefinitionGuid)
			assert.Equal(t, "distribution_name", response.Data.DistributionName)
		}
	})

	t.Run("when report distribution has been deleted", func(t *testing.T) {
		err := c.V2.ReportDistributions.Delete(intgGUID)
		assert.Nil(t, err)

		var response api.ReportDistributionResponse
		response, err = c.V2.ReportDistributions.Get(intgGUID)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ReportDistributions/MOCK_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestReportDistributionList(t *testing.T) {
	var (
		allGUIDs            []string
		reportDistributions = generateGuids(&allGUIDs, 3)
		expectedLen         = len(allGUIDs)
		fakeServer          = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("ReportDistributions",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			reportDistributions := []string{
				generateReportDistributions(reportDistributions),
			}
			fmt.Fprintf(w,
				generateReportDistributionsResponse(
					strings.Join(reportDistributions, ", "),
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

	response, err := c.V2.ReportDistributions.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.ReportDistributionGuid)
	}
}

func TestReportDistributionUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ReportDistributions/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Update() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "distribution_name", "report distribution name is missing")
			assert.Contains(t, body, "daily", "report distribution frequency is missing")
			assert.Contains(t, body, "channel_name", "report distribution channel is missing")
		}

		fmt.Fprintf(w, generateReportDistributionResponse(singleMockReportDistribution(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	reportDistribution := api.ReportDistribution{
		ReportDistributionGuid: intgGUID,
		ReportDefinitionGuid:   "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAB",
		DistributionName:       "distribution_name",
		AlertChannels:          []string{"channel_name"},
		Frequency:              "daily",
	}

	assert.Equal(t, "distribution_name", reportDistribution.DistributionName, "report distribution name mismatch")
	assert.Equal(t, "daily", reportDistribution.Frequency, "report distribution frequency mismatch")
	assert.Equal(t, []string{"channel_name"}, reportDistribution.AlertChannels, "report distribution channel mismatch")

	response, err := c.V2.ReportDistributions.Update(intgGUID, reportDistribution.UpdateConfig())
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.ReportDistributionGuid)
		assert.Equal(t, response.Data.Frequency, "daily")
		assert.Contains(t, response.Data.DistributionName, "distribution_name")
		assert.Contains(t, response.Data.ReportDistributionGuid, intgGUID)
		assert.Contains(t, response.Data.ReportDefinitionGuid, "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAB")
	}
}

func generateReportDistributions(guids []string) string {
	reportDistributions := make([]string, len(guids))
	for i, guid := range guids {
		reportDistributions[i] = singleMockReportDistribution(guid)
	}
	return strings.Join(reportDistributions, ", ")
}

func generateReportDistributionsResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateReportDistributionResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleMockReportDistribution(id string) string {
	return fmt.Sprintf(`
{
            "reportDistributionGuid": %q,
            "reportDefinitionGuid": "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAB",
            "distributionName": "distribution_name",
            "alertChannels": [
                "channel_name"
            ],
            "frequency": "daily",
            "data": {
                "resourceGroups": [],
                "integrations": [
                    {
                        "accountId": "0123456789"
                    }
                ],
                "severities": [
                    "critical"
                ],
                "violations": [
                    "NonCompliant"
                ]
            }
        }
	`, id)
}
