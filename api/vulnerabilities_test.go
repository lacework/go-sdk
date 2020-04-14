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
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestVulnerabilitiesScan(t *testing.T) {
	expectedStatus := "Scanning"
	expectedRequestID := "efd151c8-abcd-1234-5678-13e8cca93584"
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI("external/vulnerabilities/container/repository/images/scan",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Scan should be a POST method")

			if assert.NotNil(t, r.Body) {
				body := httpBodySniffer(r)
				assert.Contains(t, body, "gcr.io", "container registry missing")
				assert.Contains(t, body, "example/repo", "wrong repository")
				assert.Contains(t, body, "v0.1.0-dev", "missing tag")
			}

			fmt.Fprintf(w, vulScanJsonResponse(expectedRequestID, expectedStatus))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.Scan(
		"gcr.io",
		"example/repo",
		"v0.1.0-dev",
	)
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.Equal(t, expectedStatus, response.Data.Status)
		assert.Equal(t, expectedRequestID, response.Data.RequestID)
	}
}

func TestVulnerabilitiesScanLaceworkError(t *testing.T) {
	expectedError := "Container registry not found"
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI("external/vulnerabilities/container/repository/images/scan",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Scan should be a POST method")

			if assert.NotNil(t, r.Body) {
				body := httpBodySniffer(r)
				assert.Contains(t, body, "example.com", "container registry missing")
				assert.Contains(t, body, "example/repo", "wrong repository")
				assert.Contains(t, body, "v0.1.0-dev", "missing tag")
			}

			fmt.Fprintf(w, vulScanErrorJsonResponse(expectedError))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.Scan(
		"example.com",
		"example/repo",
		"v0.1.0-dev",
	)
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.False(t, response.Ok)
		assert.Equal(t, expectedError, response.Message)
	}
}

func TestVulnerabilitiesScan404Error(t *testing.T) {
	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Vulnerabilities.Scan(
		"example.com",
		"example/repo",
		"v0.1.0-dev",
	)
	assert.Empty(t, response)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404 page not found")
	}
}

func vulScanJsonResponse(reqID, status string) string {
	return `
		{
			"data": { "requestId": "` + reqID + `", "Status": "` + status + `" },
			"ok": true,
			"message": "SUCCESS"
		}
	`
}
func vulScanErrorJsonResponse(message string) string {
	return `
		{
			"data": {},
			"ok": false,
			"message": "` + message + `"
		}
	`
}
