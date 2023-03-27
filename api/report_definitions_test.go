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

func TestReportDefinitionsGet(t *testing.T) {
	var (
		intgGUID         = intgguid.New()
		apiPath          = fmt.Sprintf("ReportDefinitions/%s", intgGUID)
		reportDefinition = singleMockReportDefinition(intgGUID)
		fakeServer       = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateReportDefinitionResponse(reportDefinition))
			}
		},
	)

	fakeServer.MockAPI("ReportDefinitions/UNKNOWN_INTG_GUID",
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

	t.Run("when report definition exists", func(t *testing.T) {
		var response api.ReportDefinitionResponse
		response, err := c.V2.ReportDefinitions.Get(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.ReportDefinitionGuid)
			assert.Equal(t, "mockReportDefinition", response.Data.ReportName)
			assert.Equal(t, "mockReportDefinition Display", response.Data.DisplayName)
			assert.Equal(t, "COMPLIANCE", response.Data.ReportType)
			assert.Equal(t, "AWS", response.Data.SubReportType)
			assert.Equal(t, "Test Section", response.Data.ReportDefinitionDetails.Sections[0].Title)
			assert.Equal(t, []string{"lacework-global-1"}, response.Data.ReportDefinitionDetails.Sections[0].Policies)
		}
	})

	t.Run("when report definition does NOT exist", func(t *testing.T) {
		var response api.ReportDefinitionResponse
		response, err := c.V2.ReportDefinitions.Get("UNKNOWN_INTG_GUID")
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ReportDefinitions/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestReportDefinitionsDelete(t *testing.T) {
	var (
		intgGUID         = intgguid.New()
		apiPath          = fmt.Sprintf("ReportDefinitions/%s", intgGUID)
		reportDefinition = singleMockReportDefinition(intgGUID)
		getResponse      = generateReportDefinitionResponse(reportDefinition)
		fakeServer       = lacework.MockServer()
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

	fakeServer.MockAPI("ReportDefinitions/UNKNOWN_INTG_GUID",
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

	t.Run("verify report definition exists", func(t *testing.T) {
		var response api.ReportDefinitionResponse
		response, err := c.V2.ReportDefinitions.Get(intgGUID)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.ReportDefinitionGuid)
			assert.Equal(t, "mockReportDefinition", response.Data.ReportName)
			assert.Equal(t, "mockReportDefinition Display", response.Data.DisplayName)
			assert.Equal(t, "COMPLIANCE", response.Data.ReportType)
			assert.Equal(t, "AWS", response.Data.SubReportType)
			assert.Equal(t, "Test Section", response.Data.ReportDefinitionDetails.Sections[0].Title)
			assert.Equal(t, []string{"lacework-global-1"}, response.Data.ReportDefinitionDetails.Sections[0].Policies)
		}
	})

	t.Run("when report definition has been deleted", func(t *testing.T) {
		err := c.V2.ReportDefinitions.Delete(intgGUID)
		assert.Nil(t, err)

		var response api.ReportDefinitionResponse
		response, err = c.V2.ReportDefinitions.Get(intgGUID)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ReportDefinitions/MOCK_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestReportDefinitionsList(t *testing.T) {
	var (
		allGUIDs          []string
		reportDefinitions = generateGuids(&allGUIDs, 3)
		expectedLen       = len(allGUIDs)
		fakeServer        = lacework.MockServer()
	)

	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("ReportDefinitions",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			reportDefinitions := []string{
				generateReportDefinitions(reportDefinitions),
			}
			fmt.Fprintf(w,
				generateReportDefinitionsResponse(
					strings.Join(reportDefinitions, ", "),
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

	response, err := c.V2.ReportDefinitions.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.ReportDefinitionGuid)
	}
}

func TestReportDefinitionUpdate(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("ReportDefinitions/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Update() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "mockReportDefinitionUpdate", "report definition name is missing")
			assert.Contains(t, body, "mockReportDefinitionUpdate Display", "report definition name is missing")
			assert.Contains(t, body, "lacework-global-1", "missing policies")
		}

		fmt.Fprintf(w, generateReportDefinitionResponse(singleMockReportDefinition(intgGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	reportDefinition := api.ReportDefinitionUpdate{
		ReportName:  "mockReportDefinitionUpdate",
		DisplayName: "mockReportDefinitionUpdate Display",
		ReportDefinitionDetails: &api.ReportDefinitionDetails{
			Sections: []api.ReportDefinitionSection{{
				Title:    "Test Section",
				Policies: []string{"lacework-global-1"},
			}},
		},
	}
	assert.Equal(t, "mockReportDefinitionUpdate", reportDefinition.ReportName, "report definition name mismatch")
	assert.Equal(t, "mockReportDefinitionUpdate Display", reportDefinition.DisplayName, "a new report definition should match its type")
	assert.Equal(t, 1, len(reportDefinition.ReportDefinitionDetails.Sections), "a new report definition should be enabled")

	response, err := c.V2.ReportDefinitions.Update(intgGUID, reportDefinition)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, intgGUID, response.Data.ReportDefinitionGuid)
		assert.Equal(t, "mockReportDefinition", response.Data.ReportName)
		assert.Equal(t, "mockReportDefinition Display", response.Data.DisplayName)
		assert.Equal(t, "COMPLIANCE", response.Data.ReportType)
		assert.Equal(t, "AWS", response.Data.SubReportType)
		assert.Equal(t, "Test Section", response.Data.ReportDefinitionDetails.Sections[0].Title)
		assert.Equal(t, []string{"lacework-global-1"}, response.Data.ReportDefinitionDetails.Sections[0].Policies)
	}
}

func generateReportDefinitions(guids []string) string {
	reportDefinitions := make([]string, len(guids))
	for i, guid := range guids {
		reportDefinitions[i] = singleMockReportDefinition(guid)
	}
	return strings.Join(reportDefinitions, ", ")
}

func generateReportDefinitionsResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateReportDefinitionResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleMockReportDefinition(id string) string {
	return fmt.Sprintf(`{
    "createdBy": "test.user@lacework.net",
    "createdTime": "2023-03-16T14:10:24Z",
    "displayName": "mockReportDefinition Display",
    "reportDefinition": {
      "sections": [
        {
          "category": "",
          "policies": [
            "lacework-global-1"
          ],
          "title": "Test Section"
        }
      ]
    },
    "reportDefinitionGuid": "%s",
    "reportName": "mockReportDefinition",
    "reportType": "COMPLIANCE",
    "subReportType": "AWS",
    "version": 1
}
`, id)
}
