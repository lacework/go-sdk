//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestListComponents(t *testing.T) {
	var (
		os         = "linux"
		arch       = "amd64"
		apiPath    = "Components"
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components() should be a GET method")

		query := r.URL.Query()

		assert.Equal(t, os, query.Get("os"), "Require `os` URL query parameter")
		assert.Equal(t, arch, query.Get("arch"), "Require `arch` URL query parameter")

		fmt.Fprint(w, generateComponentsResponse())
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Components.ListComponents(os, arch)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
	assert.Len(t, response.Data[0].Components, 2)
}

func generateComponentsResponse() string {
	return `
		{
			"data": [
				{
					"components": [
						{
							"id": 1,
							"name": "component-example",
							"description": "component description",
							"version": "2.1.0",
							"size": 1024,
							"type": "CLI_COMMAND"
						},
						{
							"id": 2,
							"name": "iac",
							"description": "IAC component",
							"version": "5.0.0",
							"size": 1024,
							"type": "CLI_COMMAND"
						}
					]
				}
			],
			"ok": true,
			"message": "SUCCESS"
		}
	`
}

func TestListComponentVersions(t *testing.T) {
	var (
		id         int32 = 1
		os               = "linux"
		arch             = "amd64"
		apiPath          = fmt.Sprintf("Components/%d", id)
		fakeServer       = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components() should be a GET method")

		query := r.URL.Query()

		assert.Equal(t, os, query.Get("os"), "Require `os` URL query parameter")
		assert.Equal(t, arch, query.Get("arch"), "Require `arch` URL query parameter")

		fmt.Fprint(w, generateComponentVersionsResponse())
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Components.ListComponentVersions(id, os, arch)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
}

func generateComponentVersionsResponse() string {
	return `
		{
			"data": [
				{
					"id": 1,
					"name":"component-example",
					"description":"CDK example component",
					"type":"CLI_COMMAND",
					"versions": [
						"0.7.22",
						"0.7.28",
						"0.7.29",
						"0.7.30"
					]
				}
			],
			"ok":true,
			"message":"SUCCESS"
		}
	`
}

func TestFetchComponentArtifact(t *testing.T) {
	var (
		id         int32 = 1
		os               = "linux"
		arch             = "amd64"
		version          = "0.7.30"
		apiPath          = fmt.Sprintf("Components/Artifact/%d", id)
		fakeServer       = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components() should be a GET method")

		query := r.URL.Query()

		assert.Equal(t, os, query.Get("os"), "The `os` URL query parameter is required")
		assert.Equal(t, arch, query.Get("arch"), "The `arch` URL query parameter is required")
		assert.Equal(t, version, query.Get("version"), "The `version` URL query parameter is required")

		fmt.Fprint(w, generateFetchComponentResponse())
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Components.FetchComponentArtifact(id, os, arch, version)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Data, 1)
}

func generateFetchComponentResponse() string {
	return `
		{
			"data": [
				{
					"id": 1,
					"name": "component-example",
					"version":"0.7.30",
					"size":472,
					"installMessage":"CDK example component",
					"updateMessage":"This message will be shown in the CLI when the component is updated",
					"artifact_url":"https://lw-cdk-store.s3.us-west-2.amazonaws.com/catalog/component-example/0.7.30/component-example-linux-amd64.tar.gz"
				}
			],
			"ok":true,
			"message":"SUCCESS"
		}
	`
}
