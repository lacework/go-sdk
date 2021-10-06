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

const (
	lqlDatasourceStr = `{
	"name": "CloudTrailRawEvents",
	"description": null,
	"resultSchema": [
		{
			"name": "EVENT",
			"dataType": "JSON",
			"description": null
		}
	]
}`
)

func TestQueryListSourcesMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Datasources",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListSources should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.ListSources()
	assert.Nil(t, err)
}

func TestQueryListSourcesOK(t *testing.T) {
	mockResponse := fmt.Sprintf(`{"data": [ %s ] }`, lqlDatasourceStr)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Datasources",
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

	datasourcesExpected := api.QueryDatasourcesResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &datasourcesExpected)

	var datasourcesActual api.QueryDatasourcesResponse
	datasourcesActual, err = c.V2.Query.ListSources()
	assert.Nil(t, err)
	assert.Equal(t, datasourcesExpected, datasourcesActual)
}

func TestQueryListSourcesError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Datasources",
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

	_, err = c.V2.Query.ListSources()
	assert.NotNil(t, err)
}

func TestQueryGetSourceMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Datasources/CloudTrailRawEvents",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetSource should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.GetSource("CloudTrailRawEvents")
	assert.Nil(t, err)
}

func TestQueryGetSourceOK(t *testing.T) {
	mockResponse := fmt.Sprintf(`{"data": %s }`, lqlDatasourceStr)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Datasources/CloudTrailRawEvents",
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

	datasourceExpected := api.QueryDatasourceResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &datasourceExpected)

	var datasourceActual api.QueryDatasourceResponse
	datasourceActual, err = c.V2.Query.GetSource("CloudTrailRawEvents")
	assert.Nil(t, err)
	assert.Equal(t, datasourceExpected, datasourceActual)
}

func TestQueryGetSourceBadInput(t *testing.T) {
	mockResponse := fmt.Sprintf(`{"data": %s }`, lqlDatasourceStr)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Datasources/CloudTrailRawEvents",
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

	_, err = c.V2.Query.GetSource("")
	assert.Equal(t, "datasource ID must be provided", err.Error())
}

func TestQueryGetSourceError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Datasources/CloudTrailRawEvents",
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

	_, err = c.V2.Query.GetSource("CloudTrailRawEvents")
	assert.NotNil(t, err)
}
