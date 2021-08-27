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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestResourceGroupGcpGet(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetGcpResourceGroup() should be a GET method")
		fmt.Fprintf(w, generateResourceGroupResponse(singleGcpResourceGroup(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ResourceGroups.GetGcpResourceGroup(resourceGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
	assert.Equal(t, "group_name", response.Data.Name)
	assert.Equal(t, "All Organizations and Projects", response.Data.GetProps().Description)
	assert.Equal(t, []string{"*"}, response.Data.GetProps().Projects)
	assert.Equal(t, "*", response.Data.GetProps().Organization)
}

func TestResourceGroupsGcpUpdate(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateGcpResourceGroup() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, resourceGUID, "ResourceGUID missing")
			assert.Contains(t, body, "group_name", "Resource Group  name is missing")
			assert.Contains(t, body, "GCP", "wrong Resource Group  type")
			assert.Contains(t, body, "Updated", "wrong description")
			assert.Contains(t, body, "[\"abc123\",\"cba321\"]", "wrong project ids")
		}

		fmt.Fprintf(w, generateResourceGroupResponse(singleGcpResourceGroup(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	resourceGroup := api.NewResourceGroup("group_name",
		api.GcpResourceGroup,
		api.GcpResourceGroupData{
			Props: api.GcpResourceGroupProps{
				Description:  "Updated",
				Projects:     []string{"abc123", "cba321"},
				Organization: "ORG123",
			},
		},
	)
	assert.Equal(t, "group_name", resourceGroup.Name, "Gcp Resource Group name mismatch")
	assert.Equal(t, "GCP", resourceGroup.Type, "a new Gcp Resource Group should match its type")
	assert.Equal(t, 1, resourceGroup.Enabled, "a new Gcp Resource Group should be enabled")
	resourceGroup.ResourceGuid = resourceGUID

	response, err := c.V2.ResourceGroups.UpdateGcpResourceGroup(resourceGroup)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
}

func singleGcpResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
        "isDefault": "1",
        "props": "{\"DESCRIPTION\":\"All Organizations and Projects\",\"PROJECTS\":[\"*\"],\"ORGANIZATION\":\"*\",\"UPDATED_BY\":null,\"LAST_UPDATED\":1586453993529}",
        "resourceGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "AWS",
        "enabled": 1
	}
	`
}
